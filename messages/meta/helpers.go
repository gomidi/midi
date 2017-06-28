package meta

import (
	"errors"
	"github.com/gomidi/midi/internal/midilib"
	"github.com/gomidi/midi/internal/vlq"
	"io"
)

func unexpectedMessageLengthError(s string) error {
	return errors.New(s)
}

type metaMessage struct {
	Typ  byte
	Data []byte
}

func (m *metaMessage) Bytes() []byte {
	b := []byte{byte(0xFF), m.Typ}
	b = append(b, vlq.Encode(uint32(len(m.Data)))...)
	if len(m.Data) != 0 {
		b = append(b, m.Data...)
	}
	return b
}

func readText(rd io.Reader) (string, error) {
	b, err := midilib.ReadVarLengthData(rd)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
