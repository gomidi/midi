package channel

import (
	"encoding/binary"
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
)

const (
	// PitchReset is the pitch bend value to reset the pitch wheel to zero
	PitchReset = 0

	// PitchLowest is the lowest possible value of the pitch bending
	PitchLowest = -8191

	// PitchHighest is the highest possible value of the pitch bending
	PitchHighest = 8191
)

/* http://www.somascape.org/midi/tech/mfile.html#sysex
Pitch Bend

3 bytes : En lsb msb

Apply pitch bend to all notes currently sounding on MIDI channel n.

lsb (0 - 127) and msb (0 - 127) together form a 14-bit number, allowing fine adjustment to pitch.
Using hex, 00 40 is the central (no bend) setting. 00 00 gives the maximum downwards bend,
and 7F 7F the maximum upwards bend.

The amount of pitch bend produced by these minimum and maximum settings is determined by the
receiving device's Pitch Bend Sensitivity, which can be set using RPN 00 00.
*/

// PitchBend represents a pitch bend message (aka "Portamento").
type PitchBend struct {
	channel  uint8
	value    int16
	absValue uint16
}

// Value returns the relative value of the pitch bending in relation
// to the middle (zero) point at 0 (-8191 to 8191)
func (p PitchBend) Value() int16 {
	return p.value
}

// AbsValue returns the absolute value (14bit) of the pitch bending (unsigned)
func (p PitchBend) AbsValue() uint16 {
	return p.absValue
}

// Channel returns the MIDI channel (starting with 0)
func (p PitchBend) Channel() uint8 {
	return p.channel
}

// Raw returns the raw bytes for the message
func (p PitchBend) Raw() []byte {
	r := midilib.MsbLsbSigned(p.value)

	var b = make([]byte, 2)

	binary.BigEndian.PutUint16(b, r)
	return channelMessage2(p.channel, 14, b[0], b[1])
}

// String represents the MIDI pitch bend message as a string (for debugging)
func (p PitchBend) String() string {
	return fmt.Sprintf("%T (\"Portamento\") channel %v value %v absValue %v", p, p.Channel(), p.Value(), p.AbsValue())
}

func (PitchBend) set(channel uint8, firstArg, secondArg uint8) setter2 {
	var m PitchBend
	m.channel = channel
	// The value is a signed int (relative to centre), and absoluteValue is the actual value in the file.
	m.value, m.absValue = midilib.ParsePitchWheelVals(firstArg, secondArg)
	return m
}
