package smf

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
)

// Message is a MIDI message that might appear in a SMF file, i.e. channel messages, sysex messages and meta messages.
type Message []byte

func (m Message) Bytes() []byte {
	return []byte(m)
}

func (m Message) IsPlayable() bool {
	return m.Type().IsPlayable()
}

func (m Message) Type() midi.Type {
	return GetType(m)
}

func GetType(msg []byte) midi.Type {
	if len(msg) == 0 {
		return midi.UnknownMsg
	}
	if Message(msg).IsMeta() {
		return GetMetaType(msg[1])
	} else {
		return midi.GetType(msg)
	}
}

func (m Message) Is(t midi.Type) bool {
	return m.Type().Is(t)
}

// ScanNoteOn returns true if (and only if) the message is a NoteOnMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanNoteOn(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).ScanNoteOn(channel, key, velocity)
}

// ScanNoteStart returns true if (and only if) the message is a NoteOnMsg with a velocity > 0.
// Then it also extracts the data to the given arguments
func (m Message) ScanNoteStart(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).ScanNoteStart(channel, key, velocity)
}

// ScanNoteOff returns true if (and only if) the message is a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanNoteOff(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).ScanNoteOff(channel, key, velocity)
}

// ScanChannel returns true if (and only if) the message is a ChannelMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanChannel(channel *uint8) (is bool) {
	return midi.Message(m).ScanChannel(channel)
}

// ScanNoteEnd returns true if (and only if) the message is a NoteOnMsg with a velocity == 0 or a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanNoteEnd(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).ScanNoteEnd(channel, key, velocity)
}

// PolyAfterTouch returns true if (and only if) the message is a PolyAfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanPolyAfterTouch(channel, key, pressure *uint8) (is bool) {
	return midi.Message(m).ScanPolyAfterTouch(channel, key, pressure)
}

// AfterTouch returns true if (and only if) the message is a AfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanAfterTouch(channel, pressure *uint8) (is bool) {
	return midi.Message(m).ScanAfterTouch(channel, pressure)
}

// ProgramChange returns true if (and only if) the message is a ProgramChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanProgramChange(channel, program *uint8) (is bool) {
	return midi.Message(m).ScanProgramChange(channel, program)
}

// PitchBend returns true if (and only if) the message is a PitchBendMsg.
// Then it also extracts the data to the given arguments
// Either relative or absolute may be nil, if not needed.
func (m Message) ScanPitchBend(channel *uint8, relative *int16, absolute *uint16) (is bool) {
	return midi.Message(m).ScanPitchBend(channel, relative, absolute)
}

// ControlChange returns true if (and only if) the message is a ControlChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanControlChange(channel, controller, value *uint8) (is bool) {
	return midi.Message(m).ScanControlChange(channel, controller, value)
}

// String represents the Message as a string that contains the MsgType and its properties.
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
		case m.ScanMetaTempo(&bpm):
			fmt.Fprintf(&bf, " bpm: %0.2f", bpm)
		case m.ScanMetaMeter(&val1, &val2):
			fmt.Fprintf(&bf, " meter: %v/%v", val1, val2)
		case m.ScanMetaChannel(&val1):
			fmt.Fprintf(&bf, " channel: %v", val1)
		case m.ScanMetaPort(&val1):
			fmt.Fprintf(&bf, " port: %v", val1)
		case m.ScanMetaSeqNumber(&val16):
			fmt.Fprintf(&bf, " number: %v", val16)
		case m.ScanMetaSMPTEOffsetMsg(&val1, &val2, &val3, &val4, &val5):
			fmt.Fprintf(&bf, " hour: %v minute: %v second: %v frame: %v fractframe: %v", val1, val2, val3, val4, val5)
		case m.ScanMetaSeqData(&bt):
			fmt.Fprintf(&bf, " bytes: % X", bt)
		case m.ScanMetaKeySig(&val1, &val2, &bl1, &bl2):
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
