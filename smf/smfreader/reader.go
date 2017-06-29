package smfreader

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gomidi/midi/internal/runningstatus"

	"errors"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/internal/midilib"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
)

// ReadFile opens file, calls callback with a reader and closes file
func ReadFile(file string, callback func(smf.Reader), options ...Option) error {
	f, err := os.Open(file)

	if err != nil {
		return err
	}

	defer func() {
		f.Close()
	}()

	callback(New(f, options...))

	return nil
}

// NewReader returns a smf.Reader
func New(src io.Reader, opts ...Option) smf.Reader {
	rd := &reader{
		input:           src,
		state:           stateExpectHeader,
		processedTracks: -1,
		runningStatus:   runningstatus.NewSMFReader(),
		sysexreader:     newSysexReader(),
	}

	for _, opt := range opts {
		opt(rd)
	}

	if rd.readNoteOffPedantic {
		rd.channelReader = channel.NewReader(rd.input, channel.ReadNoteOffPedantic())
	} else {
		rd.channelReader = channel.NewReader(rd.input)
	}

	return rd
}

type reader struct {
	input  io.Reader
	logger Logger

	state           state
	runningStatus   runningstatus.Reader
	processedTracks int16
	deltatime       uint32
	smf.Header

	sysexreader   *sysexReader
	channelReader channel.Reader

	// options
	failOnUnknownChunks bool
	headerIsRead        bool
	headerError         error
	readNoteOffPedantic bool
}

func (p *reader) Delta() uint32 {
	return p.deltatime
}

func (p *reader) Track() int16 {
	return p.processedTracks
}

// ReadHeader reads the header of SMF file
// If it is not called, the first call to Read will implicitely read the header.
// However to get the header information, ReadHeader must be called (which may also happen after the first message read)
func (p *reader) ReadHeader() (smf.Header, error) {
	err := p.readMThd()
	return p.Header, err
}

// Read reads the next midi message
func (p *reader) Read() (m midi.Message, err error) {

	for {
		switch p.state {
		case stateExpectHeader:
			err = p.readMThd()
		case stateExpectChunk:
			err = p.readChunk()
		case stateExpectTrackEvent:
			p.deltatime = 0
			return p.readEvent()
		case stateDone:
			return nil, io.EOF
		default:
			panic("unreachable")
		}

		if err != nil {
			return nil, err
		}
	}

	return nil, io.EOF
}

func (p *reader) log(format string, vals ...interface{}) {
	if p.logger != nil {
		p.logger.Printf(format+"\n", vals...)
	}
}

func (p *reader) readMThd() error {
	if p.headerIsRead {
		p.log("header already read, error: %v", p.headerError)
		return p.headerError
	}

	defer func() {
		p.headerIsRead = true
	}()

	var head chunkHeader
	p.headerError = head.readFrom(p.input)
	p.log("reading chunkHeader of header, error: %v", p.headerError)

	if p.headerError != nil {
		return p.headerError
	}

	if head.typ != "MThd" {
		p.log("wrong header type: %v", head.typ)
		return ErrExpectedMthd
	}

	p.headerError = p.parseHeaderData(p.input)
	p.log("reading body of header type: %v", p.headerError)

	if p.headerError != nil {
		return p.headerError
	}

	p.state = stateExpectChunk

	return nil
}

func (p *reader) readChunk() (err error) {
	var head chunkHeader
	err = head.readFrom(p.input)
	p.log("reading header of chunk: %v", err)

	if err != nil {
		// If we expect a chunk and we hit the end of the file, that's not so unexpected after all.
		// The file has to end some time, and this is the correct boundary upon which to end it.
		if err == errUnexpectedEOF {
			p.state = stateDone
			return io.EOF
		}
		return
	}

	p.log("got chunk type: %v", head.typ)
	// We have a MTrk
	if head.typ == "MTrk" {
		p.processedTracks++
		p.state = stateExpectTrackEvent
		// we are done, lets go to the track events
		return
	}

	if p.failOnUnknownChunks {
		return fmt.Errorf("unknown chunk of type %#v", head.typ)
	}

	// The header is of an unknown type, skip over it.
	_, err = io.CopyN(ioutil.Discard, p.input, int64(head.length))
	p.log("skipping chunk: %v", err)
	if err != nil {
		return
	}

	// Then we expect another chunk.
	p.state = stateExpectChunk
	return
}

func (p *reader) _readEvent(canary byte) (m midi.Message, err error) {
	p.log("_readEvent, canary: % X", canary)

	status, changed := p.runningStatus.Read(canary)
	p.log("got status: % X, changed: %v", status, changed)

	// a non-channel message has reset the status
	if status == 0 {

		switch canary {

		// both 0xF0 and 0xF7 may start a sysex in SMF files
		case 0xF0, 0xF7:
			p.log("found sysex")
			return p.sysexreader.Read(canary, p.input)

		// meta event
		case 0xFF:
			var typ byte
			typ, err = midilib.ReadByte(p.input)
			p.log("read system common type: % X, err: %v", typ, err)

			if err != nil {
				return nil, nil
			}

			// since System Common messages are not allowed within smf files, there could only be meta messages
			// all (event unknown) meta messages must be handled by the meta dispatcher
			m, err = meta.NewReader(p.input, typ).Read()
			p.log("got meta: %T", m)
		default:
			panic(fmt.Sprintf("must not happen: invalid canary % X", canary))
		}

		// on a voice/channel category message with status either given or cached (running status)
	} else {
		var arg1 = canary // assume running status - we already got arg1

		// was no running status, we have to read arg1
		if changed {
			arg1, err = midilib.ReadByte(p.input)
			if err != nil {
				return
			}
		}

		// since every possible status is covered by a voice message type, m can't be nil
		m, err = p.channelReader.Read(status, arg1)
		p.log("got channel message: %#v, err: %v", m, err)
	}

	if err != nil {
		return nil, err
	}

	if m == nil {
		panic("must not happen: unknown event should be handled inside meta.Reader")
	}

	if m == meta.EndOfTrack {
		p.log("got end of track")
		// p.absTrackTime = 0
		//p.deltatime = 0
		// Expect the next chunk midi.
		p.state = stateExpectChunk
	}

	return m, nil
}

func (p *reader) readEvent() (m midi.Message, err error) {
	if p.processedTracks > -1 && uint16(p.processedTracks) == p.NumTracks {
		p.log("last track has been read")
		p.state = stateDone
		return nil, io.EOF
	}

	var deltatime uint32

	deltatime, err = midilib.ReadVarLength(p.input)
	p.log("read delta: %v, err: %v", deltatime, err)
	if err != nil {
		return
	}

	p.deltatime = deltatime

	// read the canary in the coal mine to see, if we have a running status byte or a given one
	var canary byte
	canary, err = midilib.ReadByte(p.input)
	p.log("read canary: %v, err: %v", canary, err)

	if err != nil {
		return
	}

	return p._readEvent(canary)
}

// parseHeaderData parses SMF-header chunk header data.
func (r *reader) parseHeaderData(reader io.Reader) error {
	format, err := midilib.ReadUint16(reader)

	if err != nil {
		return err
	}

	switch format {
	case 0:
		r.Format = smf.SMF0
	case 1:
		r.Format = smf.SMF1
	case 2:
		r.Format = smf.SMF2
	default:
		return ErrUnsupportedSMFFormat
	}

	r.NumTracks, err = midilib.ReadUint16(reader)

	if err != nil {
		return err
	}

	var division uint16
	division, err = midilib.ReadUint16(reader)

	if err != nil {
		return err
	}

	// "If bit 15 of <division> is zero, the bits 14 thru 0 represent the number
	// of delta time "ticks" which make up a quarter-note."
	if division&0x8000 == 0x0000 {
		r.TimeFormat = smf.MetricTicks(division & 0x7FFF)
	} else {
		r.TimeFormat = parseTimeCode(division)
	}

	/*
			The last two bytes indicate how many Pulses (i.e. clocks) Per Quarter Note
			(abbreviated as PPQN) resolution the time-stamps are based upon, Division.
			For example, if your sequencer has 96 ppqn, this field would be (in hex):

		00 60

		Alternately, if the first byte of Division is negative, then this represents
		the division of a second that the time-stamps are based upon. The first byte
		will be -24, -25, -29, or -30, corresponding to the 4 SMPTE standards
		representing frames per second. The second byte (a positive number)
		is the resolution within a frame (ie, subframe). Typical values may
		be 4 (MIDI Time Code), 8, 10, 80 (SMPTE bit resolution), or 100.

		You can specify millisecond-based timing by the data bytes of -25 and 40 subframes.
	*/

	/* http://www.somascape.org/midi/tech/mfile.html

	tickdiv : specifies the timing interval to be used, and whether timecode (Hrs.Mins.Secs.Frames) or metrical (Bar.Beat) timing is to be used. With metrical timing, the timing interval is tempo related, whereas with timecode the timing interval is in absolute time, and hence not related to tempo.

	    Bit 15 (the top bit of the first byte) is a flag indicating the timing scheme in use :

	    Bit 15 = 0 : metrical timing
	    Bits 0 - 14 are a 15-bit number indicating the number of sub-divisions of a quarter note (aka pulses per quarter note, ppqn). A common value is 96, which would be represented in hex as 00 60. You will notice that 96 is a nice number for dividing by 2 or 3 (with further repeated halving), so using this value for tickdiv allows triplets and dotted notes right down to hemi-demi-semiquavers to be represented.

	    Bit 15 = 1 : timecode
	    Bits 8 - 15 (i.e. the first byte) specifies the number of frames per second (fps),
	    and will be one of the four SMPTE standards - 24, 25, 29 or 30, though expressed as a negative value
	    (using 2's complement notation), as follows :
	    fps	Representation (hex)
	    24 E8
	    25 E7
	    29 E3
	    30 E2


	    Bits 0 - 7 (the second byte) specifies the sub-frame resolution, i.e. the number of sub-divisions of a frame.
	    Typical values are 4 (corresponding to MIDI Time Code), 8, 10, 80 (corresponding to SMPTE bit resolution), or 100.

	    A timing resolution of 1 ms can be achieved by specifying 25 fps and 40 sub-frames, which would be encoded in hex as  E7 28.

	A complete MThd chunk thus contains 14 bytes (including the 8 byte header).
	Example
	Data (hex)	Interpretation
	4D 54 68 64 	identifier, the ascii chars 'MThd'
	00 00 00 06 	chunklen, 6 bytes of data follow . . .
	00 01 	format = 1
	00 11 	ntracks = 17
	00 60 	tickdiv = 96 ppqn, metrical time

	*/

	return nil
}

// Parse parses the timecode from the raw value returned from Header.TimeFormat if the format is TimeCode
// It returns SMPTE frames per second (29 corresponds to 30 drop frame) and the subframes.
func parseTimeCode(raw uint16) (t smf.TimeCode) {
	// bit shifting first byte to second inverting sign
	t.FramesPerSecond = uint8(int8(byte(raw>>8)) * (-1))

	// taking the second byte
	t.SubFrames = byte(raw & uint16(255))
	return
}

var errUnexpectedEOF = errors.New("Unexpected End of File found.")
