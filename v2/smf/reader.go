package smf

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/internal/runningstatus"
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

type Logger interface {
	Printf(format string, vals ...interface{})
}

// ReadFile opens file, creates the SMF and closes file
func ReadFile(file string) (*SMF, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer func() {
		f.Close()
	}()

	return ReadAll(f)
}

// ReadAll
func ReadAll(f io.Reader) (*SMF, error) {

	rd := NewReader(f)

	err := rd.ReadHeader()

	if err != nil {
		return nil, err
	}

	//fmt.Printf("SMF: %#v\n", *rd.SMF)

	err = rd.ReadTracks()
	if rd.tracksMissing() {
		return nil, ErrMissing
	}

	if err == ErrFinished || err == io.EOF {
		return rd.SMF, nil
	}

	if err != nil {
		return nil, err
	}

	return rd.SMF, nil
}

// New returns a smf.Reader
func NewReader(src io.Reader) *Reader {
	rd := &Reader{
		input:           src,
		processedTracks: -1,
		runningStatus:   runningstatus.NewSMFReader(),
		SMF:             &SMF{},
	}

	return rd
}

// Close closes the internal reader if it is an io.ReadCloser
func (r *Reader) Close() error {
	if cl, is := r.input.(io.ReadCloser); is {
		return cl.Close()
	}
	return nil
}

func (r *Reader) ReadHeader() error {
	if r.input == nil {
		return fmt.Errorf("no input defined")
	}
	if r.headerIsRead {
		return r.error
	}
	r.error = r.readMThd()
	r.headerIsRead = true

	if r.error != nil {
		return r.error
	}

	for i := 0; i < int(r.numTracks); i++ {
		r.tracks = append(r.tracks, NewTrack())
	}

	return r.error
}

type Reader struct {
	*SMF
	Logger Logger

	input               io.Reader
	isDone              bool
	expectChunk         bool
	expectedChunkLength uint32
	runningStatus       runningstatus.Reader
	processedTracks     int16
	deltatime           uint32
	headerIsRead        bool
	error               error
}

// Delta returns the delta time in ticks for the last MIDI message
func (r *Reader) Delta() uint32 {
	return r.deltatime
}

// Track returns the track for the last MIDI message
func (r *Reader) Track() int16 {
	return r.processedTracks
}

func (r *Reader) tracksMissing() bool {
	// allow the last track to skip the endoftrack message
	//return r.processedTracks+1 < int16(r.numTracks)
	return int(r.numTracks) > int(r.processedTracks)+1
}

func (r *Reader) ReadTracks() (err error) {
	var m midi.Message
	var absDelta int64

	for {
		m, err = r.Read()
		if err != nil {
			break
		}
		r.log("message %v", m)
		//fmt.Println(m)
		tr := int(r.Track())

		if m.Is(midi.MetaEndOfTrackMsg) {
			r.log("end of track")
			r.tracks[tr].Close(r.deltatime)
			absDelta = 0
			continue
		}

		absDelta += int64(r.deltatime)

		if m.Is(midi.MetaTempoMsg) {
			tc := TempoChange{
				AbsDelta: absDelta,
				BPM:      m.BPM(),
			}
			r.SMF.tempoChanges = append(r.SMF.tempoChanges, tc)
		}

		r.log("add message %v to track %v", m, tr)
		r.tracks[tr].Add(r.deltatime, m.Data)
	}

	sort.Sort(r.SMF.tempoChanges)

	return err
}

// Read reads the next midi message
// If the file has been read completely, ErrFinished is returned as error.
func (r *Reader) Read() (m midi.Message, err error) {
	msg, err := r.read()
	if err == io.EOF && r.tracksMissing() {
		return m, ErrMissing
	}
	return msg, err
}

func (r *Reader) read() (m midi.Message, err error) {
	if r.isDone {
		return m, ErrFinished
	}

	if !r.headerIsRead {
		r.error = r.ReadHeader()
	}

	if r.error != nil {
		return m, r.error
	}

	//fmt.Println("expectChunk", r.expectChunk)

	if r.expectChunk {
		r.readChunk()
	}

	if r.error != nil {
		return m, r.error
	}

	// now we are inside a track
	r.deltatime = 0
	m, r.error = r.readEvent()
	return m, r.error
}

func (r *Reader) log(format string, vals ...interface{}) {
	if r.Logger != nil {
		r.Logger.Printf(format+"\n", vals...)
	}
}

func (r *Reader) readMThd() (err error) {

	// after the header a chunk should come
	r.expectChunk = true

	var chunk chunk

	_, err = chunk.ReadHeader(r.input)
	r.log("reading header of chunk, error: %v", err)

	if err != nil {
		return
	}

	if chunk.Type() != "MThd" {
		r.log("wrong chunker type: %v", chunk.Type())
		err = errExpectedMthd
		return
	}

	err = r.parseHeaderData(r.input)
	r.log("reading body of header type: %v", err)

	return // leave at the end
}

func (r *Reader) readChunk() {

	if r.error != nil {
		return
	}

	var (
		// define the variables here that are shared along the for loop
		// err error
		//head chunkHeader
		chunk chunk
	)

	r.expectedChunkLength, r.error = chunk.ReadHeader(r.input)
	r.log("reading header of chunk: %v", r.error)

	if r.error != nil {
		// if we are here, not all tracks have been read, so io.EOF would be an error,
		// so return errors here in each case
		return
	}

	r.log("got chunk type: %v", chunk.Type())
	// We have a MTrk
	if chunk.Type() == "MTrk" {
		r.log("is track chunk")
		r.processedTracks++
		r.expectChunk = false
		//p.state = stateExpectTrackEvent
		// we are done, lets go to the track events
		return
	}

	/*
		if p.failOnUnknownChunks {
			err = fmt.Errorf("unknown chunk of type %#v", chunk.Type())
			break
		}
	*/

	// The header is of an unknown type, skip over it.
	_, r.error = io.CopyN(ioutil.Discard, r.input, int64(r.expectedChunkLength))
	r.log("skipping chunk: %v", r.error)
	if r.error != nil {
		return
	}

	r.expectChunk = true
}

func (r *Reader) _readEvent(canary byte) (m midi.Message, err error) {
	r.log("_readEvent, canary: % X", canary)
	m.MsgType = midi.UnknownMsg

	status, changed := r.runningStatus.Read(canary)
	r.log("got status: % X, changed: %v", status, changed)

	// a non-channel message has reset the status
	if status == 0 {

		switch canary {

		// both 0xF0 and 0xF7 may start a sysex in SMF files
		case 0xF0, 0xF7:
			r.log("found sysex")
			var ln uint32
			ln, err = utils.ReadVarLength(r.input)
			if err != nil {
				return m, err
			}
			bt, err := utils.ReadNBytes(int(ln), r.input)
			if err != nil {
				return m, err
			}
			m.Data = midi.SysEx(bt)
			m.MsgType = midi.MetaMsg.Set(midi.SysExMsg)
			return m, nil
		// meta event
		case 0xFF:
			var typ byte
			typ, err = utils.ReadByte(r.input)
			r.log("read system common type: % X, err: %v", typ, err)

			if err != nil {
				return m, err
			}

			var ln uint32
			ln, err = utils.ReadVarLength(r.input)
			if err != nil {
				return m, err
			}
			var bt []byte
			bt, err = utils.ReadNBytes(int(ln), r.input)
			if err != nil {
				return m, err
			}
			m.MsgType = midi.GetMetaMsgType(typ)
			//m.Data = bt
			m.Data = midi.MetaMessage(typ, bt)

			// since System Common messages are not allowed within smf files, there could only be meta messages
			// all (event unknown) meta messages must be handled by the meta dispatcher
			//m, err = newMetaReader(r.input, typ).Read()
			r.log("got meta: %T data: % X", m.MsgType, m.Data)
			//fmt.Printf("got meta: %s data: % X\n", m.Type, m.Data)
		default:
			panic(fmt.Sprintf("must not happen: invalid canary % X", canary))
		}

		// on a voice/channel category message with status either given or cached (running status)
	} else {
		var arg1 = canary // assume running status - we already got arg1

		// was no running status, we have to read arg1
		if changed {
			arg1, err = utils.ReadByte(r.input)
			if err != nil {
				return
			}
		}

		m, err = midi.ReadChannelMessage(status, arg1, r.input)

		// since every possible status is covered by a voice message type, m can't be nil
		//m, err = r.channelReader.Read(status, arg1)
		r.log("got channel message: %#v, err: %v", m, err)
	}

	if err != nil {
		r.log("got err: %v", err)
		return m, err
	}

	if m.Is(midi.MetaEndOfTrackMsg) {
		r.log("got end of track")
		// p.absTrackTime = 0
		//p.deltatime = 0
		// Expect the next chunk midi.

		/*
			if p.processedTracks > -1 && uint16(p.processedTracks) == p.header.NumTracks {
				p.log("last track has been read")
				p.isDone = true
				return nil, ErrFinished
			}
		*/

		// TODO check the read length of the track against the length thas has been read
		// return ErrTruncatedTrack if meta.EndOfTrack comes to early or ErrOverflowingTrack it it comes too late
		if uint16(r.processedTracks+1) == r.numTracks {
			r.log("last track has been read")
			r.isDone = true
		} else {
			r.expectChunk = true
		}

		// p.state = stateExpectChunk
	}

	r.log("returning: %v", m)
	return m, nil
}

func (r *Reader) readEvent() (m midi.Message, err error) {
	if r.error != nil {
		return m, r.error
	}

	//fmt.Println("readevent called")

	var deltatime uint32

	deltatime, err = utils.ReadVarLength(r.input)
	r.log("read delta: %v, err: %v", deltatime, err)
	if err != nil {
		return
	}

	r.deltatime = deltatime

	// read the canary in the coal mine to see, if we have a running status byte or a given one
	var canary byte
	canary, err = utils.ReadByte(r.input)
	r.log("read canary: %v, err: %v", canary, err)

	//fmt.Printf("read canary: %v, err: %v", canary, err)

	if err != nil {
		return
	}

	return r._readEvent(canary)
}

// parseHeaderData parses SMF-header chunk header data.
func (r *Reader) parseHeaderData(reader io.Reader) error {

	format, err := utils.ReadUint16(reader)

	if err != nil {
		return err
	}

	switch format {
	case 0:
		r.format = 0
	case 1:
		r.format = 1
	case 2:
		r.format = 2
	default:
		return errUnsupportedSMFFormat
	}

	r.numTracks, err = utils.ReadUint16(reader)

	if err != nil {
		return err
	}

	var division uint16
	division, err = utils.ReadUint16(reader)

	if err != nil {
		return err
	}

	// "If bit 15 of <division> is zero, the bits 14 thru 0 represent the number
	// of delta time "ticks" which make up a quarter-note."
	if division&0x8000 == 0x0000 {
		r.TimeFormat = MetricTicks(division & 0x7FFF)
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
func parseTimeCode(raw uint16) (t TimeCode) {
	// bit shifting first byte to second inverting sign
	t.FramesPerSecond = uint8(int8(byte(raw>>8)) * (-1))

	// taking the second byte
	t.SubFrames = byte(raw & uint16(255))
	return
}
