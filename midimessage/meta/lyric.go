package meta

import (
	"fmt"
	"io"
)

// Lyric represents a lyric MIDI message
type Lyric string

// String represents the lyric message as a string (for debugging)
func (m Lyric) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

func (m Lyric) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return Lyric(text), nil
}

// Raw returns the raw MIDI data
func (m Lyric) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteLyric),
		Data: []byte(m),
	}).Bytes()
}

// Text returns the text of the lyric
func (m Lyric) Text() string {
	return string(m)
}

func (m Lyric) meta() {}
