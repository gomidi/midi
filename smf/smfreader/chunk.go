package smfreader

import (
	"github.com/gomidi/midi/internal/midilib"
	"io"
)

// A chunk header
type chunkHeader struct {
	typ    string
	length uint32
}

func (c *chunkHeader) readFrom(rd io.Reader) error {
	b, err := midilib.ReadNBytes(4, rd)

	if err != nil {
		return err
	}

	c.length, err = midilib.ReadUint32(rd)
	c.typ = string(b)

	return err
}
