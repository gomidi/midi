package channel

import (
	"encoding/binary"
	"fmt"
)

/* http://www.somascape.org/midi/tech/mfile.html#sysex
Pitch Bend

3 bytes : En lsb msb

Apply pitch bend to all notes currently sounding on MIDI channel n.

lsb (0 - 127) and msb (0 - 127) together form a 14-bit number, allowing fine adjustment to pitch.
Using hex, 00 40 is the central (no bend) setting. 00 00 gives the maximum downwards bend, and 7F 7F the maximum upwards bend.

The amount of pitch bend produced by these minimum and maximum settings is determined by the receiving device's Pitch Bend Sensitivity, which can be set using RPN 00 00.
*/

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
	r := msbLsbSigned(p.value)

	var b = make([]byte, 2)

	binary.BigEndian.PutUint16(b, r)
	return channelMessage2(p.channel, 14, b[0], b[1])
}

func (p PitchWheel) String() string {
	return fmt.Sprintf("%T channel %v value %v absValue %v", p, p.channel, p.value, p.absValue)
}

func (PitchWheel) set(channel uint8, firstArg, secondArg uint8) setter2 {
	var m PitchWheel
	m.channel = channel
	// The value is a signed int (relative to centre), and absoluteValue is the actual value in the file.
	m.value, m.absValue = parsePitchWheelVals(firstArg, secondArg)
	return m
}

func parsePitchWheelVals(b1 byte, b2 byte) (relative int16, absolute uint16) {
	var val uint16 = 0

	val = uint16((b2)&0x7f) << 7
	val |= uint16(b1) & 0x7f

	// Turn into a signed value relative to the centre.
	relative = int16(val) - 0x2000

	return relative, val
}
