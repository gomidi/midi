package meta

import (
	"fmt"
	"io"
)

// Track represents the MIDI track message
type Track string

// String repesents the MIDI track message as a string (for debugging)
func (m Track) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

// Raw returns the raw MIDI data
func (m Track) Raw() []byte {
	return (&metaMessage{
		Typ:  byteTrack,
		Data: []byte(m),
	}).Bytes()
}

func (m Track) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return Track(text), nil
}

// Text returns the name of the track
func (m Track) Text() string {
	return string(m)
}

func (m Track) meta() {}
