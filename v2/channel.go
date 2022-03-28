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

// NewPitchbend returns the  bytes of a pitch bend message on the MIDI channel.
// If value is > 8191 (max), it will be set to 8191. If value is < -8192, it will be set to -8192.
// A value of 0 is considered as neutral position.
func (ch Channel) NewPitchbend(value int16) Msg {
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

// NewPolyAfterTouch returns the bytes of the polyphonic aftertouch message on the MIDI channel.
func (ch Channel) NewPolyAfterTouch(key, pressure uint8) Msg {
	if key > 127 {
		key = 127
	}
	if pressure > 127 {
		pressure = 127
	}
	return channelMessage2(ch.Index(), 10, key, pressure)
}

// NewNoteOn returns the bytes of a note on message on the MIDI channel.
func (ch Channel) NewNoteOn(key, velocity uint8) Msg {
	if key > 127 {
		key = 127
	}
	if velocity > 127 {
		velocity = 127
	}
	return channelMessage2(ch.Index(), 9, key, velocity)
}

// NewNoteOffVelocity returns the bytes of a note off message with velocity on the MIDI channel.
func (ch Channel) NewNoteOffVelocity(key, velocity uint8) Msg {
	if key > 127 {
		key = 127
	}
	if velocity > 127 {
		velocity = 127
	}
	return channelMessage2(ch.Index(), 8, key, velocity)
}

// NewNoteOff returns the bytes of a note off message on the MIDI channel.
func (ch Channel) NewNoteOff(key uint8) Msg {
	if key > 127 {
		key = 127
	}
	return channelMessage2(ch.Index(), 8, key, 0)
}

// NewProgramChange returns the bytes of a program change message on the MIDI channel.
func (ch Channel) NewProgramChange(program uint8) Msg {
	if program > 127 {
		program = 127
	}
	return channelMessage1(ch.Index(), 12, program)
}

// NewAfterTouch returns the bytes of an aftertouch message on the MIDI channel.
func (ch Channel) NewAfterTouch(pressure uint8) Msg {
	if pressure > 127 {
		pressure = 127
	}
	return channelMessage1(ch.Index(), 13, pressure)
}

// NewControlChange returns the bytes of a control change message on the MIDI channel.
func (ch Channel) NewControlChange(controller, value uint8) Msg {
	if controller > 127 {
		controller = 127
	}
	if value > 127 {
		value = 127
	}
	return channelMessage2(ch.Index(), 11, controller, value)
}
