package meta

import (
	"fmt"
	"io"
)

type TrackInstrument string

func (m TrackInstrument) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

func (m TrackInstrument) Raw() []byte {
	return (&metaMessage{
		Typ:  byteTrackInstrument,
		Data: []byte(m),
	}).Bytes()
}

func (m TrackInstrument) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return TrackInstrument(text), nil
}

func (m TrackInstrument) Text() string {
	return string(m)
}

func (m TrackInstrument) meta() {}
