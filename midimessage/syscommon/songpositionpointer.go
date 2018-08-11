package syscommon

import (
	"encoding/binary"
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
	"io"
)

func clearBitU16(n uint16, pos uint16) uint16 {
	mask := ^(uint16(1) << pos)
	n &= mask
	return n
}

// takes a 14bit uint and pads it to 16 bit like in the specs for e.g. pitchbend
func msbLsbUnsigned(n uint16) uint16 {
	if n > 16383 {
		panic("n must not overflow 14bits (max 16383)")
	}

	lsb := n << 8
	lsb = clearBitU16(lsb, 15)
	lsb = clearBitU16(lsb, 7)

	// 0x7f = 127 = 0000000001111111
	msb := 0x7f & (n >> 7)
	return lsb | msb
}

func (m SongPositionPointer) readFrom(rd io.Reader) (Message, error) {
	bt, err := midilib.ReadNBytes(2, rd)
	if err != nil {
		return nil, err
	}

	_, abs := midilib.ParsePitchWheelVals(bt[1], bt[0])
	return SongPositionPointer(abs), nil

	/*

			// TODO: check if it is correct
			val := uint16((bt[1])&0x7f) << 7
			val |= uint16(bt[0]) & 0x7f

		return SongPositionPointer(val), nil
	*/
}

type SongPositionPointer uint16

func (m SongPositionPointer) Number() uint16 {
	return uint16(m)
}

func (m SongPositionPointer) String() string {
	return fmt.Sprintf("%T: %v", m, m.Number())
}

func (m SongPositionPointer) Raw() []byte {
	r := msbLsbUnsigned(uint16(m))
	var b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, r)

	return []byte{0xF2, b[1], b[0]}
}
func (m SongPositionPointer) sysCommon() {}
