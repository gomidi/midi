package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

type Lyric string

func (m Lyric) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m Lyric) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return Lyric(text), nil
}

func (m Lyric) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteLyric),
		Data: []byte(m),
	}).Bytes()
}

func (m Lyric) Text() string {
	return string(m)
}

func (m Lyric) meta() {}
