package channel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gomidi/midi/internal/lib"
)

type PitchWheel struct {
	channel  uint8
	value    int16
	absValue uint16
}

func (p PitchWheel) Value() int16 {
	return p.value
}

func (p PitchWheel) AbsValue() uint16 {
	return p.absValue
}

func (p PitchWheel) Channel() uint8 {
	return p.channel
}

func (p PitchWheel) Raw() []byte {
	r := lib.MsbLsbSigned(p.value)

	var bf bytes.Buffer
	//	binary.Write(&bf, binary.BigEndian, uint16(change))
	binary.Write(&bf, binary.BigEndian, r)
	b := bf.Bytes()
	return channelMessage2(p.channel, 14, b[0], b[1])
}

func (p PitchWheel) String() string {
	return fmt.Sprintf("%T channel %v value %v absValue %v", p, p.channel, p.value, p.absValue)
}

func (PitchWheel) set(channel uint8, firstArg, secondArg uint8) setter2 {
	var m PitchWheel
	m.channel = channel
	// The value is a signed int (relative to centre), and absoluteValue is the actual value in the file.
	m.value, m.absValue = lib.ParsePitchWheelVals(firstArg, secondArg)
	return m
}
