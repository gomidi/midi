package meta

import (
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
	"io"
)

type Undefined struct {
	Typ  byte
	Data []byte
}

func (m Undefined) String() string {
	return fmt.Sprintf("%T type: % X", m, m.Typ)
}

func (m Undefined) Raw() []byte {
	return (&metaMessage{
		Typ:  m.Typ,
		Data: m.Data,
	}).Bytes()
}

func (m Undefined) readFrom(rd io.Reader) (Message, error) {
	data, err := midilib.ReadVarLengthData(rd)

	if err != nil {
		return nil, err
	}

	return Undefined{m.Typ, data}, nil
}

func (m Undefined) meta() {}
