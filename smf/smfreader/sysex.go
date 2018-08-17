package smfreader

import (
	// "fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
	"github.com/gomidi/midi/midimessage/sysex"
)

func newSysexReader() *sysexReader {
	return &sysexReader{}
}

type sysexReader struct {
	inSequence bool
}

func (s *sysexReader) Read(startcode byte, rd io.Reader) (sys sysex.Message, err error) {
	/*
		what this means to us is relatively simple:
		we read the data after the startcode based of the following length
		and return the sysex chunk with the start code.
		If it ends with F7 or not, is not our business (the device has to deal with it).
		Also, if there are multiple sysexes belonging to each other yada yada.
	*/

	switch startcode {
	case 0xF0:
		// fmt.Println("read sysex with startcode 0xF0")
		var data []byte
		data, err = midilib.ReadVarLengthData(rd)

		if err != nil {
			return nil, err
		}

		// complete sysex
		if data[len(data)-1] == 0xF7 {
			s.inSequence = false
			return sysex.SysEx(data[0 : len(data)-1]), nil
		}

		// casio style
		s.inSequence = true
		return sysex.Start(data), nil

	case 0xF7:
		var data []byte
		data, err = midilib.ReadVarLengthData(rd)

		if err != nil {
			return nil, err
		}

		// End of sysex sequence
		if data[len(data)-1] == 0xF7 {
			// casio style
			if s.inSequence {
				s.inSequence = false
				return sysex.End(data[0 : len(data)-1]), nil
			}
			return sysex.Escape(data), nil

		}
		// casio style
		if s.inSequence {
			return sysex.Continue(data), nil
		}
		return sysex.Escape(data), nil

	default:
		panic("sysex in SMF must start with F0 or F7")
	}

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

// even better readable: (from http://www.somascape.org/midi/tech/mfile.html#sysex)

/*
			SysEx events

There are a couple of ways of encoding System Exclusive messages. The normal method is to encode them as a single event, though it is also possible to split messages into separate packets (continuation events). A third form (an escape sequence) is used to wrap up arbitrary bytes that could not otherwise be included in a MIDI file.
Single (complete) SysEx messages

F0 length message

length is a variable length quantity (as used to represent delta-times) which specifies the number of bytes in the following message.
message is the remainder of the system exclusive message, minus the initial 0xF0 status byte.

Thus, it is just like a normal system exclusive message, though with the additional length parameter.

Note that although the terminal 0xF7 is redundant (strictly speaking, due to the use of a length parameter) it must be included.
Example

The system exclusive message :
F0 7E 00 09 01 F7

would be encoded (without the preceding delta-time) as :
F0 05 7E 00 09 01 F7

(In case you're wondering, this is a General MIDI Enable message.)
SysEx messages sent as packets - Continuation events

Some older MIDI devices, with slow onboard processors, cannot cope with receiving a large amount of data en-masse, and require large system exclusive messages to be broken into smaller packets, interspersed with a pause to allow the receiving device to process a packet and be ready for the next one.

This approach can of course be used with the method described above, i.e. with each packet being a self-contained system exclusive message (i.e. each starting with 0xF0 and ending with 0xF7).

Unfortunately, some manufacturers (notably Casio) have chosen to bend the standard, and rather than sending the packets as self-contained system exclusive messages, they act as though running status applied to system exclusive messages (which it doesn't - or at least it shouldn't).

What Casio do is this : the first packet has an initial 0xF0 byte but doesn't have a terminal 0xF7. The last packet doesn't have an initial 0xF0 but does have a terminal 0xF7. All intermediary packets have neither. No unrelated events should occur between these packets. The idea is that all the packets can be stitched together at the receiving device to create a single system exclusive message.

Putting this into a MIDI file, the first packet uses the 0xF0 status, whereas the second and subsequent packets use the 0xF7 status. This use of the 0xF7 status is referred to as a continuation event.
Example

A 3-packet message :
F0 43 12 00
43 12 00 43 12 00
43 12 00 F7

with a 200-tick delay between the first two, and a 100-tick delay between the final two, would be encoded (without the initial delta-time, before the first packet) :
F0 03 43 12 00 	first packet (the 4 bytes F0,43,12,00 are transmitted)
81 48 	200-tick delta-time
F7 06 43 12 00 43 12 00 	second packet (the 6 bytes 43,12,00,43,12,00 are transmitted)
64 	100-tick delta-time
F7 04 43 12 00 F7 	third packet (the 4 bytes 43,12,00,F7 are transmitted)

See the note below regarding distinguishing packets and escape sequences (which both use the 0xF7 status).
Escape sequences

F7 length bytes

length is a variable length quantity which specifies the number of bytes in bytes.

This has nothing to do with System Exclusive messages as such, though it does use the 0xF7 status. It provides a way of including bytes that could not otherwise be included within a MIDI file, e.g. System Common and System Real Time messages (Song Position Pointer, MTC, etc).

Note that Escape sequences do not have a terminal 0xF7 byte.
Example

The Song Select System Common message :
F3 01

would be encoded (without the preceding delta-time) as :
F7 02 F3 01

You are not restricted to single messages per escape sequence - any arbitrary collection of bytes may be included in a single sequence.

Note Parsing the 0xF7 status byte

When an event with an 0xF7 status byte is encountered whilst reading a MIDI file, its interpretation (SysEx packet or escape sequence) is determined as follows :

    When an event with 0xF0 status but lacking a terminal 0xF7 is encountered, then this is the first of a Casio-style multi-packet message, and a flag (boolean variable) should be set to indicate this.

    If an event with 0xF7 status is encountered whilst this flag is set, then this is a continuation event (a system exclusive packet, one of many).
    If this event has a terminal 0xF7, then it is the last packet and flag should be cleared.

    If an event with 0xF7 status is encountered whilst flag is clear, then this event is an escape sequence.

Naturally, the flag should be initialised clear prior to reading each track of a MIDI file.

*/
