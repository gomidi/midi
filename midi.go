package midi

import (
	"errors"
)

// Message is a MIDI message
type Message interface {
	// String inspects the MIDI message in an informative way
	String() string

	// Raw returns the raw bytes of the MIDI message
	Raw() []byte
}

// Writer writes MIDI messages
type Writer interface {
	// Write writes the given MIDI message and returns any error
	Write(Message) error
}

// Reader reads MIDI messages
type Reader interface {
	// Read reads a MIDI message
	Read() (Message, error)
}

// WriteCloser is a Writer that must be closed at the end of writing.
type WriteCloser interface {
	Writer
	Close() error
}

// ReadCloser is a Reader that must be closed at the end of reading.
type ReadCloser interface {
	Reader
	Close() error
}

// ErrUnexpectedEOF is returned, when an unexspected end of file is reached.
var ErrUnexpectedEOF = errors.New("Unexpected End of File found.")

/*
   A MIDI message is made up of an eight-bit status byte which is generally followed by one or two data bytes.

   MIDI message (status byte + 1-2 data bytes)
      |
      -------- Channel Message (channel number included in status byte) 1000 | 1001 | 1010 | 1011 | 1100 | 1101 | 1110
      |            |
      |            ------ Channel Voice Message (musical performance)
      |            |
      |            ------ Mode Message  (how to does instr. respond to Channel Voice message)
      |                   (1011nnnn Channel Mode Message)
      |
      ---------System Message (no channel number in status byte), all beginning with 1111
                   |
                   ------ System Common Messages
                   |
                   ------ System Real Time Messages
                   |
                   ------ System Exclusive Messages (F0, F7)

   There are a number of different types of MIDI messages. At the highest level, MIDI messages are classified
   as being either Channel Messages or System Messages.

   Channel messages are those which apply to a specific
   Channel, and the Channel number is included in the status byte for these messages. System messages are not
   Channel specific, and no Channel number is indicated in their status bytes.

   Channel Messages may be further classified as being either Channel Voice Messages, or Mode Messages.
   Channel Voice Messages carry musical performance data, and these messages comprise most of the traffic in
   a typical MIDI data stream. Channel Mode messages affect the way a receiving instrument will respond to the
   Channel Voice messages.

   MIDI System Messages are classified as being System Common Messages, System Real Time Messages, or
   System Exclusive Messages. System Common messages are intended for all receivers in the system. System
   Real Time messages are used for synchronisation between clock-based MIDI components. System Exclusive
   messages include a Manufacturer's Identification (ID) code, and are used to transfer any number of data bytes
   in a format specified by the referenced manufacturer.
*/

/*
   read in the next byte (uint8)

   if it is FF -> meta event
   if it is F0 or F7 -> sysex event
   else ->
      System Common Message F0-FF

      F0 1111 0000 sysex event
      F7 1111 0111 sysex event
      FF 1111 1111 meta event

      channel voice message D7-D0

      D0 1101 0000
      D7 1101 0111


*/
