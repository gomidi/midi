package midireader

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/internal/midilib"
	"github.com/gomidi/midi/internal/runningstatus"
	"github.com/gomidi/midi/messages/syscommon"
	"github.com/gomidi/midi/messages/sysex"
	"io"

	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/realtime"
)

// New returns a new reader for reading "live", "streaming", "over the wire", "realtime" midi messages (you name it).
// When calling Read, any intermediate System Realtime Message will be ignored (if rthandler is nil) or passed to rthandler (if not)
// and other midi message will be returned normally.
//
// The Reader does no buffering and makes no attempt to close src.
// If src.Read returns an io.EOF, the reader stops reading.
func New(src io.Reader, rthandler func(realtime.Message), options ...Option) midi.Reader {
	rd := &reader{
		input:         realtime.NewReader(src, rthandler),
		runningStatus: runningstatus.NewLiveReader(),
	}

	for _, opt := range options {
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
	input               realtime.Reader
	runningStatus       runningstatus.Reader
	channelReader       channel.Reader
	readNoteOffPedantic bool
}

// read starts the reading.
func (p *reader) Read() (ev midi.Message, err error) {
	// read the canary in the coal mine to see, if we have a running status byte or a given one
	var canary byte
	canary, err = midilib.ReadByte(p.input)

	if err != nil {
		return
	}

	return p.readMsg(canary)
}

func (p *reader) discardUntilNextStatus() (canary byte, err error) {
	/*
		A device should be able to "ignore" all MIDI messages that it doesn't use, including currently undefined MIDI messages
		(ie Status is 0xF4, 0xF5, or 0xFD). In other words, a device is expected to be able to deal with all MIDI messages that it
		could possibly be sent, even if it simply ignores those messages that aren't applicable to the device's functions.

		If a MIDI message is not a RealTime Category message, then the way to ignore the message is to throw away its Status and
		all data bytes (ie, bit #7 clear) up to the next received, non-RealTime Status byte.
	*/

	for {
		canary, err = midilib.ReadByte(p.input)

		if err != nil {
			return
		}

		if midilib.IsStatusByte(canary) {
			return
		}
	}

	return
}

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

// readSysEx reads a sysex "over the wire", "in live mode", "as a stream" - you name it -
// opposed to reading a sysex from a SMF standard midi file
// the sysex has already been started (0xF0 has been read)
// we need a realtime.Reader here, since realtime messages must be handled (or ignored from the viewpoit of sysex)
// here we can ignore incomplete casio style messages (since they are only interrupted in time)
func (p *reader) readSysEx() (sys sysex.SysEx, status byte, err error) {
	var b byte
	var bf []byte
	// read byte by byte
	for {
		b, err = midilib.ReadByte(p.input)
		if err != nil {
			break
		}

		// the normal way to terminate
		if b == byte(0xF7) {
			sys = sysex.SysEx(bf)
			return
		}

		// not so elegant way to terminate by sending a new status
		if midilib.IsStatusByte(b) {
			sys = sysex.SysEx(bf)
			status = b
			return
		}

		bf = append(bf, b)
	}

	// any error, especially io.EOF is considered a failure.
	// however return the sysex that had been received so far back to the user
	// and leave him to decide what to do.
	sys = sysex.SysEx(bf)
	return
}

func (p *reader) readMsg(canary byte) (m midi.Message, err error) {
	status, changed := p.runningStatus.Read(canary)

	//	fmt.Printf("canary: % X, status: % X\n", canary, status)

	// the cached running status has been reset, because a status byte
	// came in from a non channel message
	if status == 0 {

		// on a system common message
		switch canary {

		/* start sysex */
		case 0xF0:
			m, status, err = p.readSysEx()

			// TODO check if that works
			/*
				the idea is:
				1. sysex was aborted/finished because a status byte came. it returns the status that it has been read
				2. p.runningStatus.Handle(status) is buffering the status that has been read from sysex
				3. on the next read, the status is missing in the source (since it already has been read). but since it is inside the running status buffer, the correct status should be found
			*/
			if status != 0 {
				p.runningStatus.Read(status)
			}

		case 0xF7:
			// we should never have a 0xF7 since sysex must already have consumed it, but ignore it gracefully and go to the next message
			return p.Read()

		default:
			// must be a system common message, but no sysex (0xF0 < canary < 0xF7)
			m, err = syscommon.NewReader(p.input, canary).Read()
		}

	} else {
		// on a voice/channel message, status came directly or from running status

		var arg1 = canary // assume running status - we already got arg1

		// was no running status, we have to read arg1
		if changed {
			arg1, err = midilib.ReadByte(p.input)
			if err != nil {
				return
			}
		}

		// fmt.Printf("read channel message\n")
		m, err = p.channelReader.Read(status, arg1)
	}

	if err != nil {
		return
	}

	// unknown event: ignore all until next status byte
	if m == nil {
		canary, err = p.discardUntilNextStatus()
		if err != nil {
			return
		}
		// handle events for the next status
		// what I don't understand: what happens to deltatimes (as they come before an status byte)
		return p.readMsg(canary)
	}

	return
}
