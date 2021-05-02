package midi

import (
	"encoding/binary"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// Channel represents a MIDI channel (0-15).
type Channel uint8

// Index returns the index (number) of the MIDI channel (0-15).
func (ch Channel) Index() uint8 {
	return uint8(ch)
}

const (
	// PitchReset is the pitch bend value to reset the pitch wheel to zero
	PitchReset = 0

	// PitchLowest is the lowest possible value of the pitch bending
	PitchLowest = -8192

	// PitchHighest is the highest possible value of the pitch bending
	PitchHighest = 8191
)

// Pitchbend returns the  bytes of a pitch bend message on the MIDI channel.
// If value is > 8191 (max), it will be set to 8191. If value is < -8192, it will be set to -8192.
// A value of 0 is considered as neutral position.
func (ch Channel) Pitchbend(value int16) []byte {
	if value > PitchHighest {
		value = PitchHighest
	}

	if value < PitchLowest {
		value = PitchLowest
	}

	r := utils.MsbLsbSigned(value)

	var b = make([]byte, 2)

	binary.BigEndian.PutUint16(b, r)
	return channelMessage2(ch.Index(), 14, b[0], b[1])
}

// PolyAftertouch returns the bytes of the polyphonic aftertouch message on the MIDI channel.
func (ch Channel) PolyAftertouch(key, pressure uint8) []byte {
	return channelMessage2(ch.Index(), 10, key, pressure)
}

// NoteOn returns the bytes of a note on message on the MIDI channel.
func (ch Channel) NoteOn(key, velocity uint8) []byte {
	return channelMessage2(ch.Index(), 9, key, velocity)
}

// NoteOffVelocity returns the bytes of a note off message with velocity on the MIDI channel.
func (ch Channel) NoteOffVelocity(key, velocity uint8) []byte {
	return channelMessage2(ch.Index(), 8, key, velocity)
}

// NoteOff returns the bytes of a note off message on the MIDI channel.
func (ch Channel) NoteOff(key uint8) []byte {
	return channelMessage2(ch.Index(), 8, key, 0)
}

// ProgramChange returns the bytes of a program change message on the MIDI channel.
func (ch Channel) ProgramChange(program uint8) []byte {
	return channelMessage1(ch.Index(), 12, program)
}

// Aftertouch returns the bytes of an aftertouch message on the MIDI channel.
func (ch Channel) Aftertouch(pressure uint8) []byte {
	return channelMessage1(ch.Index(), 13, pressure)
}

// ControlChange returns the bytes of a control change message on the MIDI channel.
func (ch Channel) ControlChange(controller, value uint8) []byte {
	return channelMessage2(ch.Index(), 11, controller, value)
}
