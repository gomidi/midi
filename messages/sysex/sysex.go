package sysex

import (
	"bytes"
	// "bytes"
	// "encoding/binary"
	"fmt"
	"github.com/gomidi/midi/internal/lib"
	"github.com/gomidi/midi/messages/realtime"
	"io"
)

// if canary >= 0xF0 && canary <= 0xF7 {
const (
	byteSysExStart = byte(0xF0)
	byteSysExEnd   = byte(0xF7)
)

type Message interface {
	String() string
	Raw() []byte
	// readFrom(io.Reader) (Message, error)
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

// ReadLive reads a sysex "over the wire", "in live mode", "as a stream" - you name it -
// opposed to reading a sysex from a SMF standard midi file
// the sysex has already been started (0xF0 has been read)
// we need a realtime.Reader here, since realtime messages must be handled (or ignored from the viewpoit of sysex)
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
		if lib.IsStatusByte(b) {
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

/*
	F0 <length> <bytes to be transmitted after F0>

	The length is stored as a variable-length quantity. It specifies the number of bytes which follow it, not
	including the F0 or the length itself. For instance, the transmitted message F0 43 12 00 07 F7 would be stored
	in a MIDI File as F0 05 43 12 00 07 F7. It is required to include the F7 at the end so that the reader of the
	MIDI File knows that it has read the entire message.
*/

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
func ReadSMF(startcode byte, rd io.Reader) (sys SysEx, err error) {
	/*
		what this means to us is relatively simple:
		we read the data after the startcode based of the following length
		and return the sysex chunk with the start code.
		If it ends with F7 or not, is not our business (the device has to deal with it).
		Also, if there are multiple sysexes belonging to each other yada yada.
	*/

	switch startcode {
	case 0xF0, 0xF7:
		var data []byte
		data, err = lib.ReadVarLengthData(rd)

		if err != nil {
			return nil, err
		}

		sys = append(sys, startcode)
		sys = append(sys, data...)

	default:
		panic("sysex in SMF must start with F0 or F7")
	}

	return

}

var _ Message = SysEx([]byte{})

type SysEx []byte

func (m SysEx) Bytes() []byte {
	return []byte(m)
}

/*
// TODO: implement
func (m SysEx) readFrom(rd io.Reader) (Message, error) {
	return m, nil
}
*/

func (m SysEx) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

func (m SysEx) Len() int {
	return len(m)
}

func (m SysEx) Raw() []byte {
	var b = []byte{0xF0}
	b = append(b, []byte(m)...)
	b = append(b, 0xF7)
	return b
}
