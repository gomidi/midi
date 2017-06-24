package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

type Copyright string

func (m Copyright) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}
func (m Copyright) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return Copyright(text), nil
}

func (m Copyright) Raw() []byte {
	return (&metaMessage{
		Typ:  byteCopyright,
		Data: []byte(m),
	}).Bytes()
}

func (m Copyright) Text() string {
	return string(m)
}

func (m Copyright) meta() {}
