package sysex

import (
	"bytes"
	"github.com/gomidi/midi/internal/lib"
	"github.com/gomidi/midi/internal/runningstatus"
	"github.com/gomidi/midi/messages/realtime"
)

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

// ReadLive reads a sysex "over the wire", "in live mode", "as a stream" - you name it -
// opposed to reading a sysex from a SMF standard midi file
// the sysex has already been started (0xF0 has been read)
// we need a realtime.Reader here, since realtime messages must be handled (or ignored from the viewpoit of sysex)
// here we can ignore incomplete casio style messages (since they are only interrupted in time)
func ReadLive(rd realtime.Reader) (sys SysEx, status byte, err error) {
	var b byte
	var bf bytes.Buffer
	// read byte by byte
	for {
		b, err = lib.ReadByte(rd)
		if err != nil {
			break
		}

		// the normal way to terminate
		if b == byte(0xF7) {
			sys = SysEx(bf.Bytes())
			return
		}

		// not so elegant way to terminate by sending a new status
		if runningstatus.IsStatusByte(b) {
			sys = SysEx(bf.Bytes())
			status = b
			return
		}

		bf.WriteByte(b)
	}

	// any error, especially io.EOF is considered a failure.
	// however return the sysex that had been received so far back to the user
	// and leave him to decide what to do.
	sys = SysEx(bf.Bytes())
	return
}
