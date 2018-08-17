package meta

import (
	"fmt"
	"io"
)

/* http://www.somascape.org/midi/tech/mfile.html
Cue Point

FF 07 length text

This optional event is used to describe something that happens within a film, video or stage production at that point in the musical score. E.g. 'Car crashes', 'Door opens', etc.

For a format 1 MIDI file, Cue Point Meta events should only occur within the first MTrk chunk.
*/

// Cuepoint represents a MIDI cue point message
type Cuepoint string

// Text returns the text of the cue point
func (m Cuepoint) Text() string {
	return string(m)
}

// Raw returns the raw MIDI bytes
func (m Cuepoint) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteCuepoint),
		Data: []byte(m),
	}).Bytes()
}

// String represents the cue point MIDI message as a string (for debugging)
func (m Cuepoint) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

func (m Cuepoint) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return Cuepoint(text), nil
}

func (m Cuepoint) meta() {}
