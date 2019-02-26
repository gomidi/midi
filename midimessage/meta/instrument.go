package meta

import (
	"fmt"
	"io"
)

// Instrument represents the MIDI instrument name message
type Instrument string

// String repesents the MIDI track message as a string (for debugging)
func (m Instrument) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

// Raw returns the raw MIDI data
func (m Instrument) Raw() []byte {
	return (&metaMessage{
		Typ:  byteInstrument,
		Data: []byte(m),
	}).Bytes()
}

func (m Instrument) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return Instrument(text), nil
}

// Text returns the name of the track
func (m Instrument) Text() string {
	return string(m)
}

func (m Instrument) meta() {}
