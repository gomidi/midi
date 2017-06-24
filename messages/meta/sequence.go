package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

type Sequence string

func (m Sequence) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}
func (m Sequence) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return Sequence(text), nil

}

func (m Sequence) Text() string {
	return string(m)
}

func (m Sequence) meta() {}

func (m Sequence) Raw() []byte {
	return (&metaMessage{
		Typ:  byteSequence,
		Data: []byte(m),
	}).Bytes()
}
