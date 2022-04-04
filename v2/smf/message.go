package smf

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
)

// Message is a MIDI message that might appear in a SMF file, i.e. channel messages, sysex messages and meta messages.
type Message []byte

// Bytes return the underlying bytes of the message.
func (m Message) Bytes() []byte {
	return []byte(m)
}

// IsPlayable returns true, if the message can be send to an instrument.
func (m Message) IsPlayable() bool {
	return midi.Message(m).IsPlayable()
}

// IsMeta returns true, if the message is a meta message.
func (m Message) IsMeta() bool {
	if len(m) == 0 {
		return false
	}
	return m[0] == 0xFF
}

// Type returns the type of the message.
func (m Message) Type() midi.Type {
	return getType(m)
}

func getType(msg []byte) midi.Type {
	if len(msg) == 0 {
		return midi.UnknownMsg
	}
	if Message(msg).IsMeta() {
		return getMetaType(msg[1])
	} else {
		return midi.Message(msg).Type()
	}
}

// Is returns true, if the message is of the given type.
func (m Message) Is(t midi.Type) bool {
	return m.Type().Is(t)
}

// IsOneOf returns true, if the message is one of the given types.
func (m Message) IsOneOf(checkers ...midi.Type) bool {
	for _, checker := range checkers {
		if m.Is(checker) {
			return true
		}
	}
	return false
}

// GetNoteOn returns true if (and only if) the message is a NoteOnMsg.
// Then it also extracts the data to the given arguments
func (m Message) GetNoteOn(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).GetNoteOn(channel, key, velocity)
}

// GetNoteStart returns true if (and only if) the message is a NoteOnMsg with a velocity > 0.
// Then it also extracts the data to the given arguments
func (m Message) GetNoteStart(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).GetNoteStart(channel, key, velocity)
}

// GetNoteOff returns true if (and only if) the message is a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) GetNoteOff(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).GetNoteOff(channel, key, velocity)
}

// GetChannel returns true if (and only if) the message is a ChannelMsg.
// Then it also extracts the data to the given arguments
func (m Message) GetChannel(channel *uint8) (is bool) {
	return midi.Message(m).GetChannel(channel)
}

// GetNoteEnd returns true if (and only if) the message is a NoteOnMsg with a velocity == 0 or a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) GetNoteEnd(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).GetNoteEnd(channel, key, velocity)
}

// GetPolyAfterTouch returns true if (and only if) the message is a PolyAfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) GetPolyAfterTouch(channel, key, pressure *uint8) (is bool) {
	return midi.Message(m).GetPolyAfterTouch(channel, key, pressure)
}

// GetAfterTouch returns true if (and only if) the message is a AfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) GetAfterTouch(channel, pressure *uint8) (is bool) {
	return midi.Message(m).GetAfterTouch(channel, pressure)
}

// GetProgramChange returns true if (and only if) the message is a ProgramChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) GetProgramChange(channel, program *uint8) (is bool) {
	return midi.Message(m).GetProgramChange(channel, program)
}

// GetPitchBend returns true if (and only if) the message is a PitchBendMsg.
// Then it also extracts the data to the given arguments
// Either relative or absolute may be nil, if not needed.
func (m Message) GetPitchBend(channel *uint8, relative *int16, absolute *uint16) (is bool) {
	return midi.Message(m).GetPitchBend(channel, relative, absolute)
}

// GetControlChange returns true if (and only if) the message is a ControlChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) GetControlChange(channel, controller, value *uint8) (is bool) {
	return midi.Message(m).GetControlChange(channel, controller, value)
}

// String represents the Message as a string that contains the Type and its properties.
func (m Message) String() string {

	if m.IsMeta() {
		var bf bytes.Buffer
		fmt.Fprintf(&bf, m.Type().String())

		var val1 uint8
		var val2 uint8
		var val3 uint8
		var val4 uint8
		var val5 uint8
		var val16 uint16
		var bl1 bool
		var bl2 bool
		var text string
		var bpm float64
		var bt []byte

		switch {
		case m.GetMetaTempo(&bpm):
			fmt.Fprintf(&bf, " bpm: %0.2f", bpm)
		case m.GetMetaMeter(&val1, &val2):
			fmt.Fprintf(&bf, " meter: %v/%v", val1, val2)
		case m.GetMetaChannel(&val1):
			fmt.Fprintf(&bf, " channel: %v", val1)
		case m.GetMetaPort(&val1):
			fmt.Fprintf(&bf, " port: %v", val1)
		case m.GetMetaSeqNumber(&val16):
			fmt.Fprintf(&bf, " number: %v", val16)
		case m.GetMetaSMPTEOffsetMsg(&val1, &val2, &val3, &val4, &val5):
			fmt.Fprintf(&bf, " hour: %v minute: %v second: %v frame: %v fractframe: %v", val1, val2, val3, val4, val5)
		case m.GetMetaSeqData(&bt):
			fmt.Fprintf(&bf, " bytes: % X", bt)
		case m.GetMetaKeySig(&val1, &val2, &bl1, &bl2):
			fmt.Fprintf(&bf, " key: %v num: %v ismajor: %v isflat: %v", val1, val2, bl1, bl2)
		default:
			switch m.Type() {
			case MetaLyricMsg, MetaMarkerMsg, MetaCopyrightMsg, MetaTextMsg, MetaCuepointMsg, MetaDeviceMsg, MetaInstrumentMsg, MetaProgramNameMsg, MetaTrackNameMsg:
				m.text(&text)
				fmt.Fprintf(&bf, " text: %q", text)
			}
		}

		return bf.String()
	} else {
		return midi.Message(m).String()
	}

}
