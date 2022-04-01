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

// NewPitchbend returns the  bytes of a pitch bend message on the MIDI channel.
// If value is > 8191 (max), it will be set to 8191. If value is < -8192, it will be set to -8192.
// A value of 0 is considered as neutral position.
func NewPitchbend(channel uint8, value int16) Message {
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

// NewPolyAfterTouch returns the bytes of the polyphonic aftertouch message on the MIDI channel.
func NewPolyAfterTouch(channel, key, pressure uint8) Message {
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

// NewNoteOn returns the bytes of a note on message on the MIDI channel.
func NewNoteOn(channel, key, velocity uint8) Message {
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

// NewNoteOffVelocity returns the bytes of a note off message with velocity on the MIDI channel.
func NewNoteOffVelocity(channel, key, velocity uint8) Message {
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

// NewNoteOff returns the bytes of a note off message on the MIDI channel.
func NewNoteOff(channel, key uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if key > 127 {
		key = 127
	}
	return channelMessage2(channel, 8, key, 0)
}

// NewProgramChange returns the bytes of a program change message on the MIDI channel.
func NewProgramChange(channel, program uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if program > 127 {
		program = 127
	}
	return channelMessage1(channel, 12, program)
}

// NewAfterTouch returns the bytes of an aftertouch message on the MIDI channel.
func NewAfterTouch(channel, pressure uint8) Message {
	if channel > 15 {
		channel = 15
	}

	if pressure > 127 {
		pressure = 127
	}
	return channelMessage1(channel, 13, pressure)
}

// NewControlChange returns the bytes of a control change message on the MIDI channel.
func NewControlChange(channel, controller, value uint8) Message {
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
