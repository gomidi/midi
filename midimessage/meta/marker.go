package meta

import (
	"fmt"
	"io"
)

/* http://www.somascape.org/midi/tech/mfile.html
Marker

FF 06 length text

This optional event is used to label points within a sequence, e.g. rehearsal letters, loop points, or section names (such as 'First verse').

For a format 1 MIDI file, Marker Meta events should only occur within the first MTrk chunk.
*/

// Marker represents a MIDI marker message
type Marker string

// String represents the marker MIDI message as a string (for debugging)
func (m Marker) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

// Text returns the text of the marker
func (m Marker) Text() string {
	return string(m)
}

// Raw returns the raw MIDI data
func (m Marker) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteMarker),
		Data: []byte(m),
	}).Bytes()
}

func (m Marker) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return Marker(text), nil
}

func (m Marker) meta() {}
