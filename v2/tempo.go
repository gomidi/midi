package midi

import (
	"bytes"
	"fmt"
	"io"
)

// ReadUint24 parse a 3-byte 24 bit integer from a Reader.
// It returns the 32-bit value and an error.
// This is a slightly modified variant of the parseUint24 function
// from Joe Wass. See the file midi_functions.go for the original.
func ReadUint24(rd io.Reader) (uint32, error) {
	b, err := ReadNBytes(3, rd)

	if err != nil {
		return 0, err
	}

	var val uint32 = 0x00
	val |= uint32(b[2]) << 0
	val |= uint32(b[1]) << 8
	val |= uint32(b[0]) << 16

	return val, nil
}

func (m Message) BPM() float64 {
	if m.Type.IsNot(MetaTempoMsg) {
		fmt.Println("not tempo message")
		return -1
	}

	rd := bytes.NewReader(m.metaDataWithoutVarlength())
	microsecondsPerCrotchet, err := ReadUint24(rd)
	if err != nil {
		fmt.Println("cant read")
		return -1
	}

	return float64(60000000) / float64(microsecondsPerCrotchet)
}
