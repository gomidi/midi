package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

/* http://www.somascape.org/midi/tech/mfile.html
Marker

FF 06 length text

This optional event is used to label points within a sequence, e.g. rehearsal letters, loop points, or section names (such as 'First verse').

For a format 1 MIDI file, Marker Meta events should only occur within the first MTrk chunk.
*/

type Marker string

func (m Marker) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m Marker) Text() string {
	return string(m)
}

func (m Marker) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteMarker),
		Data: []byte(m),
	}).Bytes()
}

func (m Marker) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return Marker(text), nil
}

func (m Marker) meta() {}
