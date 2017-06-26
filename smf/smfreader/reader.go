package smfreader

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gomidi/midi/internal/runningstatus"

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
		input:         src,
		state:         stateExpectHeader,
		runningStatus: runningstatus.NewSMFReader(),
		sysexreader:   newSysexReader(),
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
	processedTracks uint16
	deltatime       uint32
	mthd            mThdData

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

func (p *reader) Track() uint16 {
	return p.processedTracks
}

// ReadHeader reads the header of SMF file
// If it is not called, the first call to Read will implicitely read the header.
// However to get the header information, ReadHeader must be called (which may also happen after the first message read)
func (p *reader) ReadHeader() (smf.Header, error) {
	err := p.readMThd()
	return p.mthd, err
}

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
		p.log("header already read: %v", p.headerError)
		return p.headerError
	}

	defer func() {
		p.headerIsRead = true
	}()

	var head chunkHeader
	p.headerError = head.readFrom(p.input)
	p.log("reading chunkHeader of header: %v", p.headerError)

	if p.headerError != nil {
		return p.headerError
	}

	if head.typ != "MThd" {
		p.log("wrong header type: %v", head.typ)
		return ErrExpectedMthd
	}

	p.headerError = p.mthd.readFrom(p.input)
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
		if err == smf.ErrUnexpectedEOF {
			p.state = stateDone
			return io.EOF
		}
		return
	}

	p.log("got chunk type: %v", head.typ)
	// We have a MTrk
	if head.typ == "MTrk" {
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

	// system common category status
	if status == 0 {

		switch canary {

		// both 0xF0 and 0xF7 may start a sysex in SMF files
		case 0xF0, 0xF7:
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
			panic(fmt.Sprintf("must not happen: invalid status % X", canary))
		}

		// on a voice/channel category status
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
		p.processedTracks++
		// p.absTrackTime = 0
		p.deltatime = 0
		// Expect the next chunk midi.
		p.state = stateExpectChunk
	}

	return m, nil
}

func (p *reader) readEvent() (m midi.Message, err error) {
	if p.processedTracks == p.mthd.numTracks {
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
