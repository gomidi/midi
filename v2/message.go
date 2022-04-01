package midi

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// Message represents a live MIDI message. It can be created from the MIDI bytes of a message, by calling NewMessage.
type Message struct {

	// Type represents the message type of the MIDI message
	Type

	// Data contains the bytes of the MiDI message
	Data []byte
}

// NewMessage returns a new Message from the bytes of the message, by finding the correct type.
// If the type could not be found, the Type of the Message is UnknownType.
func NewMessage(bt []byte) (m Message) {
	m.Type = GetType(bt)
	m.Data = bt
	return
}

// NoteOn returns true if (and only if) the message is a NoteOnMsg.
// Then it also extracts the data to the given arguments
func (m Message) NoteOn(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOn) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *velocity = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return true
}

// NoteStart returns true if (and only if) the message is a NoteOnMsg with a velocity > 0.
// Then it also extracts the data to the given arguments
func (m Message) NoteStart(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOn) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *velocity = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	if *velocity == 0 {
		return false
	}
	return true
}

// NoteOff returns true if (and only if) the message is a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) NoteOff(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOff) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *velocity = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return true
}

// Channel returns true if (and only if) the message is a ChannelMsg.
// Then it also extracts the data to the given arguments
func (m Message) Channel(channel *uint8) (is bool) {
	if !m.Is(ChannelType) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	return true
}

// NoteEnd returns true if (and only if) the message is a NoteOnMsg with a velocity == 0 or a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) NoteEnd(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOn) && !m.Is(NoteOff) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *velocity = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return m.Is(NoteOff) || *velocity == 0
}

// PolyAfterTouch returns true if (and only if) the message is a PolyAfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) PolyAfterTouch(channel, key, pressure *uint8) (is bool) {
	if !m.Is(PolyAfterTouch) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *pressure = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return true
}

// AfterTouch returns true if (and only if) the message is a AfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) AfterTouch(channel, pressure *uint8) (is bool) {
	if !m.Is(AfterTouch) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*pressure = utils.ParseUint7(m.Data[1])
	return true
}

// ProgramChange returns true if (and only if) the message is a ProgramChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) ProgramChange(channel, program *uint8) (is bool) {
	if !m.Is(ProgramChange) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*program = utils.ParseUint7(m.Data[1])
	return true
}

// PitchBend returns true if (and only if) the message is a PitchBendMsg.
// Then it also extracts the data to the given arguments
// Either relative or absolute may be nil, if not needed.
func (m Message) PitchBend(channel *uint8, relative *int16, absolute *uint16) (is bool) {
	if !m.Is(PitchBend) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])

	rel, abs := utils.ParsePitchWheelVals(m.Data[1], m.Data[2])
	if relative != nil {
		*relative = rel
	}
	if absolute != nil {
		*absolute = abs
	}
	return true
}

// ControlChange returns true if (and only if) the message is a ControlChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) ControlChange(channel, controller, value *uint8) (is bool) {
	if !m.Is(ControlChange) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*controller, *value = utils.ParseTwoUint7(m.Data[1], m.Data[2])
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

// MTC represents a MIDI timing code message (quarter frame)
func (m Message) MTC(quarterframe *uint8) (is bool) {
	if !m.Is(MTC) {
		return false
	}

	*quarterframe = utils.ParseUint7(m.Data[1])
	return true
}

// Song returns the song number of a MIDI song select system message
func (m Message) SongSelect(song *uint8) (is bool) {
	if !m.Is(SongSelect) {
		return false
	}

	*song = utils.ParseUint7(m.Data[1])
	return true
}

// SPP returns the song position pointer of a MIDI song position pointer system message
func (m Message) SPP(spp *uint16) (is bool) {
	if !m.Is(SPP) {
		return false
	}

	_, *spp = utils.ParsePitchWheelVals(m.Data[2], m.Data[1])
	return true
}

// String represents the Message as a string that contains the MsgType and its properties.
func (m Message) String() string {
	var bf bytes.Buffer
	fmt.Fprintf(&bf, m.Type.String())

	var channel, val1, val2 uint8
	var pitchabs uint16
	var pitchrel int16
	//	var text string
	//	var bpm float64
	var spp uint16

	switch {
	case m.NoteOn(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " channel: %v key: %v velocity: %v", channel, val1, val2)
	case m.NoteOff(&channel, &val1, &val2):
		if val2 > 0 {
			fmt.Fprintf(&bf, " channel: %v key: %v velocity: %v", channel, val1, val2)
		} else {
			fmt.Fprintf(&bf, " channel: %v key: %v", channel, val1)
		}
	case m.PolyAfterTouch(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " channel: %v key: %v pressure: %v", channel, val1, val2)
	case m.AfterTouch(&channel, &val1):
		fmt.Fprintf(&bf, " channel: %v pressure: %v", channel, val1)
	case m.ControlChange(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " channel: %v controller: %v value: %v", channel, val1, val2)
	case m.ProgramChange(&channel, &val1):
		fmt.Fprintf(&bf, " channel: %v program: %v", channel, val1)
	case m.PitchBend(&channel, &pitchrel, &pitchabs):
		fmt.Fprintf(&bf, " channel: %v pitch: %v (%v)", channel, pitchrel, pitchabs)
		/*
			case m.Tempo(&bpm):
				fmt.Fprintf(&bf, " bpm: %0.2f", bpm)
			case m.Meter(&val1, &val2):
				fmt.Fprintf(&bf, " meter: %v/%v", val1, val2)
			case m.IsOneOf(MetaLyricMsg, MetaMarkerMsg, MetaCopyrightMsg, MetaTextMsg, MetaCuepointMsg, MetaDeviceMsg, MetaInstrumentMsg, MetaProgramNameMsg, MetaTrackNameMsg):
				m.text(&text)
				fmt.Fprintf(&bf, " text: %q", text)
		*/
	case m.MTC(&val1):
		fmt.Fprintf(&bf, " mtc: %v", val1)
	case m.SPP(&spp):
		fmt.Fprintf(&bf, " spp: %v", spp)
	case m.SongSelect(&val1):
		fmt.Fprintf(&bf, " song: %v", val1)
	default:
	}

	return bf.String()
}
