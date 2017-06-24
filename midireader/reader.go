package midireader

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/internal/lib"
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
func New(src io.Reader, rthandler func(realtime.Message)) midi.Reader {
	return &reader{
		input: realtime.NewReader(src, rthandler),
	}
}

type reader struct {
	input         realtime.Reader
	runningStatus lib.RunningStatus
}

// read starts the reading.
func (p *reader) Read() (ev midi.Message, err error) {
	// read the canary in the coal mine to see, if we have a running status byte or a given one
	var canary byte
	canary, err = lib.ReadByte(p.input)

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
		canary, err = lib.ReadByte(p.input)

		if err != nil {
			return
		}

		if lib.IsStatusByte(canary) {
			return
		}
	}

	return
}

func (p *reader) readMsg(canary byte) (ev midi.Message, err error) {
	status, _ := p.runningStatus.HandleLive(canary)

	if status != 0 {
		// on a voice/channel message
		ev, err = channel.NewReader(p.input, status).Read()

	} else {
		// on a system common message
		switch canary {

		/* start sysex */
		case 0xF0:
			ev, status, err = sysex.ReadLive(p.input)

			// TODO check if that works
			/*
				the idea is:
				1. sysex was aborted/finished because a status byte came. it returns the status that it has been read
				2. p.runningStatus.Handle(status) is buffering the status that has been read from sysex
				3. on the next read, the status is missing in the source (since it already has been read). but since it is inside the running status buffer, the correct status should be found
			*/
			if status != 0 {
				p.runningStatus.HandleLive(status)
			}

		case 0xF7:
			// we should never have a 0xF7 since sysex must already have consumed it, but ignore it gracefully and go to the next message
			return p.Read()

		default:
			// must be a system common message, but no sysex (0xF0 < canary < 0xF7)
			ev, err = syscommon.NewReader(p.input, canary).Read()
		}
	}

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
		return p.readMsg(canary)
	}

	return
}
