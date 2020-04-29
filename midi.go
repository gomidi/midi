package midi

import (
	"fmt"
)

// Message is a MIDI message.
type Message interface {

	// String inspects the MIDI message in an informative way.
	String() string

	// Raw returns the raw bytes of the MIDI message.
	Raw() []byte
}

// Writer writes MIDI messages.
type Writer interface {

	// Write writes the given MIDI message and returns any error.
	Write(Message) error
}

// Reader reads MIDI messages.
type Reader interface {
	
	// Read reads a MIDI message.
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
var ErrUnexpectedEOF = fmt.Errorf("Unexpected End of File found.")
