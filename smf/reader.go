package smf

import (
	"fmt"
	"io"
	"io/ioutil"
	"lib"

	midi "github.com/gomidi/midi"
	"github.com/gomidi/midi/channel"
	"github.com/gomidi/midi/meta"
)

type ReadOption func(*reader)

type Reader interface {
	midi.Reader
	Track() uint16
	ReadHeader() (Header, error)
	Delta() uint32
}

// NewReader returns a Reader
func NewReader(src io.Reader, opts ...ReadOption) Reader {
	rd := &reader{input: src, state: stateExpectHeader}

	for _, opt := range opts {
		opt(rd)
	}

	return rd
}

func FailOnUnknownChunks() ReadOption {
	return func(r *reader) {
		r.failOnUnknownChunks = true
	}
}

// PostHeader tells the reader that next read is after the smf header
// remainingtracks are the number of tracks that are going to be parsed (must be > 0)
func PostHeader(remainingtracks uint16) ReadOption {
	if remainingtracks == 0 {
		panic("remainingtracks must be at least 1")
	}
	return func(r *reader) {
		r.mthd.numTracks = remainingtracks
		r.state = stateExpectChunk
	}
}

// InsideTrack tells the reader that next read is inside a track (after the track header)
// remainingtracks are the number of tracks that are going to be parsed (must be > 0)
func InsideTrack(remainingtracks uint16) ReadOption {
	if remainingtracks == 0 {
		panic("remainingtracks must be at least 1")
	}
	return func(r *reader) {
		r.mthd.numTracks = remainingtracks
		r.state = stateExpectTrackEvent
	}
}

type Header interface {
	Format() format
	NumTracks() uint16
	TimeFormat() (timeformat, uint16)
}

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

func (p *reader) Delta() uint32 {
	return p.deltatime
}

func (p *reader) Track() uint16 {
	return p.processedTracks
}

func (p *reader) ReadHeader() (Header, error) {
	err := p.readMThd()
	return p.mthd, err
}

func (p *reader) Read() (ev midi.Event, err error) {

	for {
		switch p.state {
		case stateExpectHeader:
			err = p.readMThd()
		case stateExpectChunk:
			err = p.readChunk()
		case stateExpectTrackEvent:
			//err = p.readEvent()
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

// filereader is a Standard Midi File reader.
// Pass this a ReadSeeker to a MIDI file and EventHandler
// and it'll run over the file, EventHandlers HandleEvent method for each midi.
type reader struct {
	input io.Reader

	// State of the parser, as per the above constants.
	state state

	runningStatusBuffer byte

	processedTracks uint16

	absTrackTime uint64

	deltatime uint32

	sysexBuffer []byte
	inSysEx     bool

	mthd mThdData

	failOnUnknownChunks bool
}

func (p *reader) readMThd() error {

	var head chunkHeader
	err := head.readFrom(p.input)

	if err != nil {
		return err
	}

	if head.typ != "MThd" {
		return ErrExpectedMthd
	}

	err = p.mthd.readFrom(p.input)

	if err != nil {
		return err
	}

	p.state = stateExpectChunk

	return nil
}

func (p *reader) readChunk() (err error) {
	var head chunkHeader
	err = head.readFrom(p.input)

	if err != nil {
		// If we expect a chunk and we hit the end of the file, that's not so unexpected after all.
		// The file has to end some time, and this is the correct boundary upon which to end it.
		if err == lib.ErrUnexpectedEOF {
			p.state = stateDone
			return io.EOF
		}
		return
	}

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
	if err != nil {
		return
	}

	// Then we expect another chunk.
	p.state = stateExpectChunk
	return
}

func (p *reader) readMetaEvent(command byte) (ev midi.Event, err error) {

	var met meta.Event = meta.Dispatch(command)

	// fmt.Printf("could not find meta command %X\n", command)

	if met == nil {
		return nil, nil
	}

	return meta.ReadFrom(met, p.input)
}

/*
his (http://midi.teragonaudio.com/tech/midispec.htm) take on running status buffer
A recommended approach for a receiving device is to maintain its "running status buffer" as so:

    Buffer is cleared (ie, set to 0) at power up.
    Buffer stores the status when a Voice Category Status (ie, 0x80 to 0xEF) is received.
    Buffer is cleared when a System Common Category Status (ie, 0xF0 to 0xF7) is received.
    Nothing is done to the buffer when a RealTime Category message is received.
    Any data bytes are ignored when the buffer is 0.
*/

/*
    Each RealTime Category message (ie, Status of 0xF8 to 0xFF) consists of only 1 byte, the Status. These messages are primarily concerned with timing/syncing functions which means that they must be sent and received at specific times without any delays. Because of this, MIDI allows a RealTime message to be sent at any time, even interspersed within some other MIDI message. For example, a RealTime message could be sent inbetween the two data bytes of a Note On message. A device should always be prepared to handle such a situation; processing the 1 byte RealTime message, and then subsequently resume processing the previously interrupted message as if the RealTime message had never occurred.

For more information about RealTime, read the sections Running Status, Ignoring MIDI Messages, and Syncing Sequence Playback.
*/

/*
   Furthermore, although the 0xF7 is supposed to mark the end of a SysEx message, in fact, any status
   (except for Realtime Category messages) will cause a SysEx message to be
   considered "done" (ie, actually "aborted" is a better description since such a scenario
   indicates an abnormal MIDI condition). For example, if a 0x90 happened to be sent sometime
   after a 0xF0 (but before the 0xF7), then the SysEx message would be considered
   aborted at that point. It should be noted that, like all System Common messages,
   SysEx cancels any current running status. In other words, the next Voice Category
   message (after the SysEx message) must begin with a Status.
*/

/*
func (p *filereader) finishSysex() (err error) {
	p.inSysEx = false
	ev := sysEx(p.sysexBuffer)
	p.sysexBuffer = nil
	continueReading := p.handler.OnTrackEvent(p.processedTracks, p.absTrackTime, ev)

	if !continueReading {
		p.state = stateDone
		err = io.EOF
	}

	return
}
*/

func (p *reader) _readEvent(canary byte) (ev midi.Event, err error) {
	//var rawevent, channel, canary, firstArg uint8

	var rawevent, ch, firstArg uint8

	/*
	   his (http://midi.teragonaudio.com/tech/midispec.htm) take on running status buffer
	   A recommended approach for a receiving device is to maintain its "running status buffer" as so:

	       Buffer is cleared (ie, set to 0) at power up.
	       Buffer stores the status when a Voice Category Status (ie, 0x80 to 0xEF) is received.
	       Buffer is cleared when a System Common Category Status (ie, 0xF0 to 0xF7) is received.
	       Nothing is done to the buffer when a RealTime Category message is received.
	       Any data bytes are ignored when the buffer is 0.
	*/

	// on a voice/channel category status: store the runningStatusBuffer
	if canary >= 0x80 && canary <= 0xEF {
		p.runningStatusBuffer = canary
	}

	// on a system common category status: clear the runningStatusBuffer
	if canary >= 0xF0 && canary <= 0xF7 {
		p.runningStatusBuffer = 0
	}

	if p.inSysEx && lib.IsStatusByte(canary) {
		/*
			err = p.finishSysex2()
			if err != nil {
				return err
			}
		*/
		return nil, nil
	}

	// system common category status
	if p.runningStatusBuffer == 0 {

		if p.inSysEx {
			var b byte
			b, err = lib.ReadByte(p.input)
			p.sysexBuffer = append(p.sysexBuffer, b)
			// TODO do further scanning/reading
			return nil, nil
		}

		switch canary {
		/* start sysex */
		case 0xF0:
			p.inSysEx = true
			// TODO do further scanning/reading
			return nil, nil

		/* end sysex */
		case 0xF7:
			if !p.inSysEx {
				panic("must not happen: finishing sysex that never started, severe error")
			}
			//return p.finishSysex()
			return nil, nil
		/*
			case 0xFF:
				firstArg, err = readByte(p.input)

				if err != nil {
					return
				}
		*/
		default:
			return p.readMetaEvent(canary)
		}

		// on a voice/channel category status
	} else {
		rawevent, ch = lib.ParseStatus(canary)

		firstArg, err = lib.ReadByte(p.input)

		if err != nil {
			return
		}

		switch rawevent {

		// one argument only
		case lib.CodeProgramChange, lib.CodeChannelPressure:
			ev = channel.New(ch).Dispatch1(rawevent, firstArg)
			//ev = eventhelper.GetChannelEvent1(rawevent, channel, firstArg)

		// two Arguments needed
		default:
			ev, err = channel.New(ch).Dispatch2(rawevent, firstArg, p.input)
			//ev, err = eventhelper.GetChannelEvent2(rawevent, channel, firstArg, p.input)
		}
	}

	if err != nil {
		return nil, err
	}

	if ev == meta.EndOfTrack {
		p.processedTracks++
		p.absTrackTime = 0
		p.deltatime = 0
		// Expect the next chunk midi.
		p.state = stateExpectChunk
		return ev, nil
	}

	// fallback for unsupported events
	if ev == nil {
		ev = midi.UnknownEvent([]byte{ch, rawevent, firstArg})
	}

	return ev, nil
}

func (p *reader) readEvent() (ev midi.Event, err error) {
	if p.processedTracks == p.mthd.numTracks {
		p.state = stateDone
		return nil, io.EOF
	}

	var deltatime uint32

	deltatime, err = lib.ReadVarLength(p.input)
	if err != nil {
		return
	}

	p.deltatime = deltatime

	// we have to set the absTrackTime in any case, so lets do it early on
	p.absTrackTime += uint64(deltatime)

	// read the canary in the coal mine to see, if we have a running status byte or a given one
	var canary byte
	canary, err = lib.ReadByte(p.input)

	if err != nil {
		return
	}

	return p._readEvent(canary)
}
