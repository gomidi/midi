package meta

import (
	"fmt"
	"io"
)

// Copyright represents the MIDI copyright message
type Copyright string

// String returns the copyright message as a string (for debugging)
func (m Copyright) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}
func (m Copyright) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return Copyright(text), nil
}

// Raw returns the raw MIDI data
func (m Copyright) Raw() []byte {
	return (&metaMessage{
		Typ:  byteCopyright,
		Data: []byte(m),
	}).Bytes()
}

// Text returns the copyright text
func (m Copyright) Text() string {
	return string(m)
}

func (m Copyright) meta() {}
