package syscommon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

func (m SongPositionPointer) readFrom(rd io.Reader) (Message, error) {

	bt, err := lib.ReadN(2, rd)

	if err != nil {
		return nil, err
	}

	if len(bt) != 2 {
		err = lib.UnexpectedMessageLengthError("SongPositionPointer expected length 2")
		return nil, err
	}

	// TODO: check if it is correct
	val := uint16((bt[1])&0x7f) << 7
	val |= uint16(bt[0]) & 0x7f

	return SongPositionPointer(val), nil
}

type SongPositionPointer uint16

func (m SongPositionPointer) Number() uint16 {
	return uint16(m)
}

func (m SongPositionPointer) String() string {
	return fmt.Sprintf("%T: %v", m, uint16(m))
}

// TODO test
func (m SongPositionPointer) Raw() []byte {
	// TODO check - it is totally a guess at the moment

	r := lib.MsbLsbUnsigned(uint16(m))

	var bf bytes.Buffer
	//	binary.Write(&bf, binary.BigEndian, uint16(change))
	binary.Write(&bf, binary.BigEndian, 0xF2)

	binary.Write(&bf, binary.BigEndian, r)
	return bf.Bytes()
}
func (m SongPositionPointer) sysCommon() {}
