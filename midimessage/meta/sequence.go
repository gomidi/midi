package meta

import (
	"fmt"
	"io"
)

// Sequence represents a MIDI sequence message (name of a sequence)
type Sequence string

// String represents the MIDI sequence message as a string (for debugging)
func (m Sequence) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}
func (m Sequence) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return Sequence(text), nil

}

// Text returns the name of the sequence
func (m Sequence) Text() string {
	return string(m)
}

func (m Sequence) meta() {}

// Raw returns the raw bytes for the message
func (m Sequence) Raw() []byte {
	return (&metaMessage{
		Typ:  byteSequence,
		Data: []byte(m),
	}).Bytes()
}
