package midi

import (
	//"eventhelper"
	"fmt"
	"io"
	"lib"

	"github.com/gomidi/midi/channel"
	"github.com/gomidi/midi/meta"
	"github.com/gomidi/midi/realtime"
)

// NewReader returns a new reader. Realtime events will be passed asynchronously
// to the realtimeReader channel when reading.
//
// The Reader does no buffering and makes no attempt to close neither src nor realtimeReader.
// If a meta.EndOfTrack message is received or src returns an io.EOF, the reader stops reading.
func NewReader(src io.Reader, realtimeReader chan<- Event) Reader {
	p := &reader{
		realInput:      src,
		realtimeStream: realtimeReader,
	}

	p.input = &realTimeReader{
		midiReader: p,
		input:      p.realInput,
	}

	return p
}

type Reader interface {
	// Read reads the next midi event from the stream
	// if end of file is reached io.EOF is returned as err and event is nil
	// if an error is not recoverable, the next call to read will return nil, io.EOF
	Read() (event Event, err error)
}

// reader is a Standard Midi File reader.
// Pass this a ReadSeeker to a MIDI file and EventHandler
// and it'll run over the file, EventHandlers HandleEvent method for each
type reader struct {
	realInput io.Reader
	input     *realTimeReader

	runningStatusBuffer byte

	sysexBuffer    []byte
	inSysEx        bool
	realtimeStream chan<- Event
}

// read starts the reading.
func (p *reader) Read() (ev Event, err error) {
	return p.readEvent()
}

func (p *reader) readMetaEvent(command byte) (ev Event, err error) {

	var met meta.Event = meta.Dispatch(command)

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

// a realTimeReader is needed, since realtime events may come in at any time
type realTimeReader struct {
	midiReader *reader
	input      io.Reader
}

func (r *realTimeReader) Read(target []byte) (n int, err error) {
	var bf []byte
	var one int

	for {
		if n == len(target) {
			return
		}
		bf = make([]byte, 1)

		one, err = r.input.Read(bf)

		if err != nil {
			return
		}

		if one != 1 {
			err = fmt.Errorf("could not read %v byte(s)", len(target))
			return
		}

		// => no realtime message
		if bf[0] < 0xF8 && bf[0] > 0xFF {
			// return bf[0], nil
			target[n] = bf[0]
			n++
			continue
		}

		// error needed here to be able to interrupt the reading from the callback (handler)
		// then an io.EOF error is returned and propagated to midireader.read()
		r.midiReader.handleRealtimeMessage(bf[0])

		if err != nil {
			return
		}
	}
}

func (p *reader) handleRealtimeMessage(b byte) {
	ev := realtime.Dispatch(b)

	if ev != nil {
		p.realtimeStream <- ev
	}
}

func (p *reader) finishSysex() (ev Event) {
	p.inSysEx = false
	ev = meta.SysEx(p.sysexBuffer)
	p.sysexBuffer = nil
	return
}

func (p *reader) discardUntilNextStatus() (canary byte, err error) {
	/*
		A device should be able to "ignore" all MIDI messages that it doesn't use, including currently undefined MIDI messages
		(ie Status is 0xF4, 0xF5, or 0xFD). In other words, a device is expected to be able to deal with all MIDI messages that it
		could possibly be sent, even if it simply ignores those messages that aren't applicable to the device's functions.

		If a MIDI message is not a RealTime Category message, then the way to ignore the message is to throw away its Status and
		all data bytes (ie, bit #7 clear) up to the next received, non-RealTime Status byte.
	*/
	// isStatusByte(canary)
	var isStatus bool

	for isStatus == false {
		canary, err = lib.ReadByte(p.input)

		if err != nil {
			return
		}

		isStatus = lib.IsStatusByte(canary)
	}

	return
}

func (p *reader) _readEvent(canary byte) (ev Event, err error) {
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
		ev = p.finishSysex()
		return
	}

	// var ev _Event

	// system common category status
	if p.runningStatusBuffer == 0 {

		if p.inSysEx {
			var b byte
			b, err = lib.ReadByte(p.input)
			p.sysexBuffer = append(p.sysexBuffer, b)
			return
		}

		switch canary {
		/* start sysex */
		case 0xF0:
			p.inSysEx = true
			return

		/* end sysex */
		case 0xF7:
			if !p.inSysEx {
				panic("must not happen: finishing sysex that never started, severe error")
			}
			ev = p.finishSysex()
			return

		default:
			ev, err = p.readMetaEvent(canary)

			if err != nil {
				return
			}

			// unknown event: ignore all until next status byte
			if ev == nil {
				canary, err = p.discardUntilNextStatus()
				if err != nil {
					return
				}
				// handle events for the next status
				// what I don't understand: what happens to deltatimes (as they come before an status byte)
				return p._readEvent(canary)
			}

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

		// two Arguments needed
		default:
			ev, err = channel.New(ch).Dispatch2(rawevent, firstArg, p.input)
		}
	}

	if err != nil {
		return
	}

	// fallback for unsupported events
	if ev == nil {
		ev = UnknownEvent([]byte{ch, rawevent, firstArg})
	}

	return
}

func (p *reader) readEvent() (ev Event, err error) {

	// read the canary in the coal mine to see, if we have a running status byte or a given one
	var canary byte
	canary, err = lib.ReadByte(p.input)

	if err != nil {
		return
	}

	return p._readEvent(canary)
}
