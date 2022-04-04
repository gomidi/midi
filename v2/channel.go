package midi

import (
	"encoding/binary"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

const (
	// PitchReset is the pitch bend value to reset the pitch wheel to zero
	PitchReset = 0

	// PitchLowest is the lowest possible value of the pitch bending
	PitchLowest = -8192

	// PitchHighest is the highest possible value of the pitch bending
	PitchHighest = 8191
)

// Pitchbend returns a pitch bend message.
// If value is > 8191 (max), it will be set to 8191. If value is < -8192, it will be set to -8192.
// A value of 0 is considered as neutral position.
func Pitchbend(channel uint8, value int16) Message {
	if channel > 15 {
		channel = 15
	}

	if value > PitchHighest {
		value = PitchHighest
	}

	if value < PitchLowest {
		value = PitchLowest
	}

	r := utils.MsbLsbSigned(value)

	var b = make([]byte, 2)

	binary.BigEndian.PutUint16(b, r)
	return channelMessage2(channel, 14, b[0], b[1])
}

// PolyAfterTouch returns a polyphonic aftertouch message.
func PolyAfterTouch(channel, key, pressure uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if key > 127 {
		key = 127
	}
	if pressure > 127 {
		pressure = 127
	}
	return channelMessage2(channel, 10, key, pressure)
}

// NoteOn returns a note on message.
func NoteOn(channel, key, velocity uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if key > 127 {
		key = 127
	}
	if velocity > 127 {
		velocity = 127
	}
	return channelMessage2(channel, 9, key, velocity)
}

// NoteOffVelocity returns a note off message with velocity.
func NoteOffVelocity(channel, key, velocity uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if key > 127 {
		key = 127
	}
	if velocity > 127 {
		velocity = 127
	}
	return channelMessage2(channel, 8, key, velocity)
}

// NoteOff returns a note off message.
func NoteOff(channel, key uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if key > 127 {
		key = 127
	}
	return channelMessage2(channel, 8, key, 0)
}

// ProgramChange returns a program change message.
func ProgramChange(channel, program uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if program > 127 {
		program = 127
	}
	return channelMessage1(channel, 12, program)
}

// AfterTouch returns an aftertouch message.
func AfterTouch(channel, pressure uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if pressure > 127 {
		pressure = 127
	}
	return channelMessage1(channel, 13, pressure)
}

// ControlChange returns a control change message.
func ControlChange(channel, controller, value uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if controller > 127 {
		controller = 127
	}
	if value > 127 {
		value = 127
	}
	return channelMessage2(channel, 11, controller, value)
}
