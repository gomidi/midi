package smfreader

import (
	"fmt"
	"io"
	"io/ioutil"
	"lib"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
)

type state int

const (
	// At the start of the MIDI file.
	// Expect SMF Header chunk.
	stateExpectHeader state = 0

	// Expect a chunk. Any kind of chunk. Except MThd.
	// But really, anything other than MTrk would be weird.
	stateExpectChunk state = 1

	// We're in a Track, expect a track midi.
	stateExpectTrackEvent state = 2

	// This has to happen sooner or later.
	stateDone state = 3
)

type Logger interface {
	Printf(format string, vals ...interface{})
}

// NewReader returns a smf.Reader
func New(src io.Reader, opts ...Option) smf.Reader {
	rd := &reader{input: src, state: stateExpectHeader}

	for _, opt := range opts {
		opt(rd)
	}

	return rd
}

// filereader is a Standard Midi File reader.
// Pass this a ReadSeeker to a MIDI file and EventHandler
// and it'll run over the file, EventHandlers HandleEvent method for each midi.
type reader struct {
	input  io.Reader
	logger Logger
	// State of the parser, as per the above constants.
	state               state
	runningStatus       lib.RunningStatus
	processedTracks     uint16
	absTrackTime        uint64
	deltatime           uint32
	mthd                mThdData
	failOnUnknownChunks bool
	headerIsRead        bool
	headerError         error
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
		if err == lib.ErrUnexpectedEOF {
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

	status, changed := p.runningStatus.Handle(canary)
	p.log("got status: % X, changed: %v", status, changed)

	// system common category status
	if status == 0 {
		var typ byte
		typ, err = lib.ReadByte(p.input)
		p.log("read system common type: % X, err: %v", typ, err)

		if err != nil {
			return nil, nil
		}

		switch typ {
		/* start sysex */
		//
		/*
			F0 <length> <bytes to be transmitted after F0>

			The length is stored as a variable-length quantity. It specifies the number of bytes which follow it, not
			including the F0 or the length itself. For instance, the transmitted message F0 43 12 00 07 F7 would be stored
			in a MIDI File as F0 05 43 12 00 07 F7. It is required to include the F7 at the end so that the reader of the
			MIDI File knows that it has read the entire message.
		*/
		case 0xF0, 0xF7:
			p.log("got sysex")
			// TODO do further scanning/reading
			return nil, nil

			/*
				   Another form of sysex event is provided which does not imply that an F0 should be transmitted. This may be
				   used as an "escape" to provide for the transmission of things which would not otherwise be legal, including
				   system realtime messages, song pointer or select, MIDI Time Code, etc. This uses the F7 code:

				   F7 <length> <all bytes to be transmitted>

				   Unfortunately, some synthesiser manufacturers specify that their system exclusive messages are to be
				   transmitted as little packets. Each packet is only part of an entire syntactical system exclusive message, but
				   the times they are transmitted are important. Examples of this are the bytes sent in a CZ patch dump, or the
				   FB-01's "system exclusive mode" in which microtonal data can be transmitted. The F0 and F7 sysex events
				   may be used together to break up syntactically complete system exclusive messages into timed packets.
				   An F0 sysex event is used for the first packet in a series -- it is a message in which the F0 should be
				   transmitted. An F7 sysex event is used for the remainder of the packets, which do not begin with F0. (Of
				   course, the F7 is not considered part of the system exclusive message).
				   A syntactic system exclusive message must always end with an F7, even if the real-life device didn't send one,
				   so that you know when you've reached the end of an entire sysex message without looking ahead to the next
				   event in the MIDI File. If it's stored in one complete F0 sysex event, the last byte must be an F7. There also
				   must not be any transmittable MIDI events in between the packets of a multi-packet system exclusive
				   message. This principle is illustrated in the paragraph below.

						Here is a MIDI File of a multi-packet system exclusive message: suppose the bytes F0 43 12 00 were to be
						sent, followed by a 200-tick delay, followed by the bytes 43 12 00 43 12 00, followed by a 100-tick delay,
						followed by the bytes 43 12 00 F7, this would be in the MIDI File:

						F0 03 43 12 00						|
						81 48											| 200-tick delta time
						F7 06 43 12 00 43 12 00   |
						64												| 100-tick delta time
						F7 04 43 12 00 F7         |

						When reading a MIDI File, and an F7 sysex event is encountered without a preceding F0 sysex event to start a
						multi-packet system exclusive message sequence, it should be presumed that the F7 event is being used as an
						"escape". In this case, it is not necessary that it end with an F7, unless it is desired that the F7 be transmitted.
			*/

		default:
			mt := meta.Dispatch(typ)
			p.log("found meta: %#v", mt)
			// since System Common messages are not allowed within smf files, there could only be meta messages
			// all (event unknown) meta messages must be handled by the meta dispatcher
			// fmt.Printf("canary: %#v input: %#v\n", canary, p.input)
			m, err = meta.ReadFrom(mt, p.input)
			p.log("got meta: %T", m)
		}

		// on a voice/channel category status
	} else {
		m, err = channel.NewReader(p.input, status).Read()
		p.log("got channel message: %#v, err: %v", m, err)
	}

	if err != nil {
		return nil, err
	}

	if m == nil {
		panic("must not happen: unknown event should be handled inside meta.Reader or channel.Reader")
	}

	if m == meta.EndOfTrack {
		p.log("got end of track")
		p.processedTracks++
		p.absTrackTime = 0
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

	deltatime, err = lib.ReadVarLength(p.input)
	p.log("read delta: %v, err: %v", deltatime, err)
	if err != nil {
		return
	}

	p.deltatime = deltatime

	// we have to set the absTrackTime in any case, so lets do it early on
	p.absTrackTime += uint64(deltatime)

	// read the canary in the coal mine to see, if we have a running status byte or a given one
	var canary byte
	canary, err = lib.ReadByte(p.input)
	p.log("read canary: %v, err: %v", canary, err)

	if err != nil {
		return
	}

	return p._readEvent(canary)
}
