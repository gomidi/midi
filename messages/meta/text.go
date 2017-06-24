package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

type Text string

func (m Text) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m Text) meta() {}

func (m Text) Raw() []byte {
	return (&metaMessage{
		Typ:  byteText,
		Data: []byte(m),
	}).Bytes()
}

func (m Text) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)
	if err != nil {
		return nil, err
	}

	return Text(text), nil
}

func (m Text) Text() string {
	return string(m)
}
