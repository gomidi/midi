package midi

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// Message is a complete midi message (not including meta messages)
type Message []byte

// Bytes returns the underlying bytes of the message.
func (m Message) Bytes() []byte {
	return []byte(m)
}

// IsPlayable returns, if the message can be send to an instrument.
func (m Message) IsPlayable() bool {
	if m.Type() <= UnknownMsg {
		return false
	}

	return m.Type() < firstMetaMsg
}

// IsOneOf returns true, if the message has one of the given types.
func (m Message) IsOneOf(checkers ...Type) bool {
	for _, checker := range checkers {
		if m.Is(checker) {
			return true
		}
	}
	return false
}

// Type returns the type of the message.
func (m Message) Type() Type {
	return getType(m)
}

// Is returns true, if the message is of the given type.
func (m Message) Is(t Type) bool {
	return m.Type().Is(t)
}

// GetNoteOn returns true if (and only if) the message is a NoteOnMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetNoteOn(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOnMsg) {
		return false
	}

	if len(m) != 3 {
		return false
	}

	if channel != nil {
		_, *channel = utils.ParseStatus(m[0])
	}

	if key != nil || velocity != nil {
		_key, _velocity := utils.ParseTwoUint7(m[1], m[2])

		if key != nil {
			*key = _key
		}

		if velocity != nil {
			*velocity = _velocity
		}
	}

	return true
}

// GetNoteStart returns true if (and only if) the message is a NoteOnMsg with a velocity > 0.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetNoteStart(channel, key, velocity *uint8) (is bool) {
	var vel uint8

	if !m.GetNoteOn(channel, key, &vel) || vel == 0 {
		return false
	}

	if velocity != nil {
		*velocity = vel
	}
	return true
}

// GetNoteOff returns true if (and only if) the message is a NoteOffMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetNoteOff(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOffMsg) {
		return false
	}

	if len(m) != 3 {
		return false
	}

	if channel != nil {
		_, *channel = utils.ParseStatus(m[0])
	}

	if key != nil || velocity != nil {
		_key, _velocity := utils.ParseTwoUint7(m[1], m[2])

		if key != nil {
			*key = _key
		}

		if velocity != nil {
			*velocity = _velocity
		}
	}

	return true
}

// GetChannel returns true if (and only if) the message is a ChannelMsg.
// Then it also extracts the channel to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetChannel(channel *uint8) (is bool) {
	if !m.Is(ChannelMsg) {
		return false
	}

	if len(m) < 1 {
		return false
	}

	if channel != nil {
		_, *channel = utils.ParseStatus(m[0])
	}
	return true
}

// GetNoteEnd returns true if (and only if) the message is a NoteOnMsg with a velocity == 0 or a NoteOffMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetNoteEnd(channel, key *uint8) (is bool) {
	if !m.Is(NoteOnMsg) && !m.Is(NoteOffMsg) {
		return false
	}

	var vel uint8
	var ch uint8
	var k uint8

	is = false

	switch {
	case m.GetNoteOn(&ch, &k, &vel):
		is = vel == 0
	case m.GetNoteOff(&ch, &k, &vel):
		is = true
	}

	if !is {
		return false
	}

	if channel != nil {
		*channel = ch
	}

	if key != nil {
		*key = k
	}

	return true
}

// GetPolyAfterTouch returns true if (and only if) the message is a PolyAfterTouchMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetPolyAfterTouch(channel, key, pressure *uint8) (is bool) {
	if !m.Is(PolyAfterTouchMsg) {
		return false
	}

	if len(m) != 3 {
		return false
	}

	if channel != nil {
		_, *channel = utils.ParseStatus(m[0])
	}

	if key != nil || pressure != nil {
		var _key, _pressure = utils.ParseTwoUint7(m[1], m[2])

		if key != nil {
			*key = _key
		}

		if pressure != nil {
			*pressure = _pressure
		}
	}
	return true
}

// GetAfterTouch returns true if (and only if) the message is a AfterTouchMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetAfterTouch(channel, pressure *uint8) (is bool) {
	if !m.Is(AfterTouchMsg) {
		return false
	}

	if len(m) != 2 {
		return false
	}

	if channel != nil {
		_, *channel = utils.ParseStatus(m[0])
	}

	if pressure != nil {
		*pressure = utils.ParseUint7(m[1])
	}
	return true
}

// GetProgramChange returns true if (and only if) the message is a ProgramChangeMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetProgramChange(channel, program *uint8) (is bool) {
	if !m.Is(ProgramChangeMsg) {
		return false
	}

	if len(m) != 2 {
		return false
	}

	if channel != nil {
		_, *channel = utils.ParseStatus(m[0])
	}

	if program != nil {
		*program = utils.ParseUint7(m[1])
	}
	return true
}

// GetPitchBend returns true if (and only if) the message is a PitchBendMsg.
// Then it also extracts the data to the given arguments.
// Either relative or absolute may be nil, if not needed.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetPitchBend(channel *uint8, relative *int16, absolute *uint16) (is bool) {
	if !m.Is(PitchBendMsg) {
		return false
	}

	if len(m) != 3 {
		return false
	}

	if channel != nil {
		_, *channel = utils.ParseStatus(m[0])
	}

	rel, abs := utils.ParsePitchWheelVals(m[1], m[2])
	if relative != nil {
		*relative = rel
	}
	if absolute != nil {
		*absolute = abs
	}
	return true
}

// GetControlChange returns true if (and only if) the message is a ControlChangeMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetControlChange(channel, controller, value *uint8) (is bool) {
	if !m.Is(ControlChangeMsg) {
		return false
	}

	if len(m) != 3 {
		return false
	}

	if channel != nil {
		_, *channel = utils.ParseStatus(m[0])
	}

	if controller != nil || value != nil {
		var _controller, _value uint8

		_controller, _value = utils.ParseTwoUint7(m[1], m[2])

		if controller != nil {
			*controller = _controller
		}

		if value != nil {
			*value = _value
		}
	}

	return true
}

/*
MTC Quarter Frame

These are the MTC (i.e. SMPTE based) equivalent of the F8 Timing Clock messages,
though offer much higher resolution, as they are sent at a rate of 96 to 120 times
a second (depending on the SMPTE frame rate). Each Quarter Frame message provides
partial timecode information, 8 sequential messages being required to fully
describe a timecode instant. The reconstituted timecode refers to when the first
partial was received. The most significant nibble of the data byte indicates the
partial (aka Message Type).

Partial	Data byte	Usage
1	0000 bcde	Frame number LSBs 	abcde = Frame number (0 to frameRate-1)
2	0001 000a	Frame number MSB
3	0010 cdef	Seconds LSBs 	abcdef = Seconds (0-59)
4	0011 00ab	Seconds MSBs
5	0100 cdef	Minutes LSBs 	abcdef = Minutes (0-59)
6	0101 00ab	Minutes MSBs
7	0110 defg	Hours LSBs 	ab = Frame Rate (00 = 24, 01 = 25, 10 = 30drop, 11 = 30nondrop)
cdefg = Hours (0-23)
8	0111 0abc	Frame Rate, and Hours MSB
*/

// GetMTC returns true if (and only if) the message is a MTCMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMTC(quarterframe *uint8) (is bool) {
	if !m.Is(MTCMsg) {
		return false
	}

	if len(m) != 2 {
		return false
	}

	if quarterframe != nil {
		*quarterframe = utils.ParseUint7(m[1])
	}

	return true
}

// GetSongSelect returns true if (and only if) the message is a SongSelectMsg.
// Then it also extracts the song number to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetSongSelect(song *uint8) (is bool) {
	if !m.Is(SongSelectMsg) {
		return false
	}

	if len(m) != 2 {
		return false
	}

	if song != nil {
		*song = utils.ParseUint7(m[1])
	}

	return true
}

// GetSPP returns true if (and only if) the message is a SPPMsg.
// Then it also extracts the spp to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetSPP(spp *uint16) (is bool) {
	if !m.Is(SPPMsg) {
		return false
	}

	if len(m) != 3 {
		return false
	}

	if spp != nil {
		_, *spp = utils.ParsePitchWheelVals(m[2], m[1])
	}

	return true
}

// String represents the Message as a string that contains the Type and its properties.
func (m Message) String() string {
	var bf bytes.Buffer
	fmt.Fprintf(&bf, m.Type().String())

	var channel, val1, val2 uint8
	var pitchabs uint16
	var pitchrel int16
	var spp uint16
	var sysex []byte

	switch {
	case m.GetNoteOn(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " channel: %v key: %v velocity: %v", channel, val1, val2)
	case m.GetNoteOff(&channel, &val1, &val2):
		if val2 > 0 {
			fmt.Fprintf(&bf, " channel: %v key: %v velocity: %v", channel, val1, val2)
		} else {
			fmt.Fprintf(&bf, " channel: %v key: %v", channel, val1)
		}
	case m.GetPolyAfterTouch(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " channel: %v key: %v pressure: %v", channel, val1, val2)
	case m.GetAfterTouch(&channel, &val1):
		fmt.Fprintf(&bf, " channel: %v pressure: %v", channel, val1)
	case m.GetControlChange(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " channel: %v controller: %v value: %v", channel, val1, val2)
	case m.GetProgramChange(&channel, &val1):
		fmt.Fprintf(&bf, " channel: %v program: %v", channel, val1)
	case m.GetPitchBend(&channel, &pitchrel, &pitchabs):
		fmt.Fprintf(&bf, " channel: %v pitch: %v (%v)", channel, pitchrel, pitchabs)
	case m.GetMTC(&val1):
		fmt.Fprintf(&bf, " mtc: %v", val1)
	case m.GetSPP(&spp):
		fmt.Fprintf(&bf, " spp: %v", spp)
	case m.GetSongSelect(&val1):
		fmt.Fprintf(&bf, " song: %v", val1)
	case m.GetSysEx(&sysex):
		fmt.Fprintf(&bf, " data: % X", sysex)
	default:
	}

	return bf.String()
}

// GetSysEx returns true, if the message is a sysex message.
// Then it extracts the inner bytes to the given slice.
func (m Message) GetSysEx(bt *[]byte) bool {
	if len(m) < 3 {
		return false
	}

	if !m.Is(SysExMsg) {
		return false
	}

	if m[0] == 0xF0 && m[len(m)-1] == 0xF7 {
		*bt = m[1 : len(m)-1]
		return true
	}

	return false
}
