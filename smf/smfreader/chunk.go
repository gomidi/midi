package smfreader

import (
	// "fmt"
	"io"
	// "io/ioutil"
	"github.com/gomidi/midi/internal/lib"
	// "github.com/gomidi/midi"
	// "github.com/gomidi/midi/messages/channel"
	// "github.com/gomidi/midi/messages/meta"
	// "github.com/gomidi/midi/smf"
)

// A chunk header
type chunkHeader struct {
	typ    string
	length uint32
}

func (c *chunkHeader) readFrom(rd io.Reader) error {
	b, err := lib.ReadN(4, rd)

	if err != nil {
		return err
	}

	c.length, err = lib.ReadUint32(rd)
	c.typ = string(b)

	return err
}
