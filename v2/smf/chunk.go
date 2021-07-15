package smf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// chunk is a chunk of a SMF file.
type chunk struct {
	typ  []byte // must always be 4 bytes long, to avoid conversions everytime, we take []byte here instead of [4]byte
	data []byte
}

// Len returns the length of the chunk body.
func (c *chunk) Len() int {
	return len(c.data)
}

// SetType sets the type of the chunk.
func (c *chunk) SetType(typ [4]byte) {
	c.typ = make([]byte, 4)
	c.typ[0] = typ[0]
	c.typ[1] = typ[1]
	c.typ[2] = typ[2]
	c.typ[3] = typ[3]
}

// Type returns the type of the chunk (from the header).
func (c *chunk) Type() string {
	var bf bytes.Buffer
	bf.Write(c.typ)
	return bf.String()
}

// Clear removes all data but keeps the type.
func (c *chunk) Clear() {
	c.data = nil
}

// WriteTo writes the content of the chunk to the given writer.
func (c *chunk) WriteTo(wr io.Writer) (int64, error) {
	if len(c.typ) != 4 {
		return 0, fmt.Errorf("chunk header not set properly")
	}

	var bf bytes.Buffer
	bf.Write(c.typ)
	binary.Write(&bf, binary.BigEndian, int32(c.Len()))
	bf.Write(c.data)
	n, err := wr.Write(bf.Bytes())
	if err != nil {
		return int64(n), fmt.Errorf("could not write chunk: %v", err)
	}
	return int64(n), nil
}

// ReadHeader reads the header from the given reader
// and returns the length of the following body.
// For errors, length of 0 is returned.
func (c *chunk) ReadHeader(rd io.Reader) (length uint32, err error) {
	c.typ, err = utils.ReadNBytes(4, rd)

	if err != nil {
		c.typ = nil
		return
	}

	return utils.ReadUint32(rd)
}

// Write writes the given bytes to the body of the chunk.
func (c *chunk) Write(b []byte) (int, error) {
	c.data = append(c.data, b...)
	return len(b), nil
}
