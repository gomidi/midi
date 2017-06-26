package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
	"github.com/gomidi/midi/smf"
)

type endOfTrack bool

const (
	EndOfTrack = endOfTrack(true)
)

func (m endOfTrack) String() string {
	return fmt.Sprintf("%T", m)
}

func (m endOfTrack) Raw() []byte {
	return (&metaMessage{
		Typ: byte(byteEndOfTrack),
	}).Bytes()
}

func (m endOfTrack) meta() {}

func (m endOfTrack) readFrom(rd io.Reader) (Message, error) {

	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 0 {
		err = smf.UnexpectedMessageLengthError("EndOfTrack expected length 0")
		return nil, err
	}

	return m, nil
}
