package meta

import (
	"fmt"
	"io"
)

type Track string

func (m Track) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

func (m Track) Raw() []byte {
	return (&metaMessage{
		Typ:  byteTrack,
		Data: []byte(m),
	}).Bytes()
}

func (m Track) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return Track(text), nil
}

func (m Track) Text() string {
	return string(m)
}

func (m Track) meta() {}
