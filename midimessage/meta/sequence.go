package meta

import (
	"fmt"
	"io"
)

type Sequence string

func (m Sequence) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}
func (m Sequence) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

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
