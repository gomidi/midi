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

type CuePoint string

func (m CuePoint) Text() string {
	return string(m)
}

func (m CuePoint) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteCuePoint),
		Data: []byte(m),
	}).Bytes()
}

func (m CuePoint) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m CuePoint) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return CuePoint(text), nil
}

func (m CuePoint) meta() {}
