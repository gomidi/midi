package smfwriter

import (
	"bytes"
	"encoding/binary"
	"io"
)

type chunk struct {
	typ  [4]byte
	data []byte
}

/*
func (c *chunk) Type() string {
	var bf bytes.Buffer
	bf.WriteByte(c.typ[0])
	bf.WriteByte(c.typ[1])
	bf.WriteByte(c.typ[2])
	bf.WriteByte(c.typ[3])
	return bf.String()
}
*/

func (c *chunk) writeTo(wr io.Writer) (int, error) {
	length := int32(len(c.data))
	var bf bytes.Buffer
	bf.WriteByte(c.typ[0])
	bf.WriteByte(c.typ[1])
	bf.WriteByte(c.typ[2])
	bf.WriteByte(c.typ[3])
	binary.Write(&bf, binary.BigEndian, length)
	bf.Write(c.data)
	return wr.Write(bf.Bytes())
}
