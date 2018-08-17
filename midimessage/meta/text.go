package meta

import (
	"fmt"
	"io"
)

// Text is a MIDI text meta message
type Text string

// String represents the MIDI text message as a string (for debugging)
func (m Text) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

func (m Text) meta() {}

// Raw returns the raw bytes for the message
func (m Text) Raw() []byte {
	return (&metaMessage{
		Typ:  byteText,
		Data: []byte(m),
	}).Bytes()
}

func (m Text) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)
	if err != nil {
		return nil, err
	}

	return Text(text), nil
}

// Text returns the text within the message
func (m Text) Text() string {
	return string(m)
}
