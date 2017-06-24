package meta

import (
	"github.com/gomidi/midi/internal/lib"
)

type metaMessage struct {
	Typ  byte
	Data []byte
}

func (m *metaMessage) Bytes() []byte {
	b := []byte{byte(0xFF), m.Typ}
	b = append(b, lib.VlqEncode(uint32(len(m.Data)))...)
	if len(m.Data) != 0 {
		b = append(b, m.Data...)
	}
	return b
}
