package midi

import (
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

type Type int8

func (t Type) IsOneOf(checkers ...Type) bool {
	for _, checker := range checkers {
		if t.Is(checker) {
			return true
		}
	}
	return false
}

func (t Type) IsAllOf(checkers ...Type) bool {
	for _, checker := range checkers {
		if !t.Is(checker) {
			return false
		}
	}
	return true
}

// t must not be a message kind (exception: sysex), but a concrete type
func (t Type) Is(checker Type) bool {

	switch {
	case t == UnknownType:
		return checker == UnknownType
	case t == SysExType:
		return checker == SysExType
	case t < UnknownType:
		return false
	case checker == UnknownType:
		return false
	case checker > UnknownType:
		return t == checker
	default:
		switch checker {
		case RealTimeType:
			return t <= reservedRealTime14
		case SysCommonType:
			return t >= MTC && t <= reservedSysCommon10
		case ChannelType:
			return t >= NoteOn && t <= reservedChannelMsg16
		case metaType:
			return t >= firstMeta
		default:
			return false
		}
	}
}

// AddTypeName adds names for new types that do not yet have a type.
// Returns false, if the type already has been named, and true on success.
func AddTypeName(m Type, name string) bool {
	if _, has := typeNames[m]; has {
		return false
	}
	typeNames[m] = name
	return true
}

var typeNames = map[Type]string{

	UnknownType:   "UnknownType",
	RealTimeType:  "RealTimeType",
	SysCommonType: "SysCommonType",
	ChannelType:   "ChannelType",
	SysExType:     "SysExType",
	//metaMsg:      "MetaType",

	Tick:        "Tick",
	TimingClock: "TimingClock",
	Start:       "Start",
	Continue:    "Continue",
	Stop:        "Stop",
	ActiveSense: "ActiveSense",
	Reset:       "Reset",

	reservedRealTime8:  "reservedRealTime8",
	reservedRealTime9:  "reservedRealTime9",
	reservedRealTime10: "reservedRealTime10",
	reservedRealTime11: "reservedRealTime11",
	reservedRealTime12: "reservedRealTime12",
	reservedRealTime13: "reservedRealTime13",
	reservedRealTime14: "reservedRealTime14",

	NoteOn:         "NoteOn",
	NoteOff:        "NoteOff",
	ControlChange:  "ControlChange",
	PitchBend:      "PitchBend",
	AfterTouch:     "AfterTouch",
	PolyAfterTouch: "PolyAfterTouch",
	ProgramChange:  "ProgramChange",

	reservedChannelMsg8:  "reservedChannelMsg8",
	reservedChannelMsg9:  "reservedChannelMsg9",
	reservedChannelMsg10: "reservedChannelMsg10",
	reservedChannelMsg11: "reservedChannelMsg11",
	reservedChannelMsg12: "reservedChannelMsg12",
	reservedChannelMsg13: "reservedChannelMsg13",
	reservedChannelMsg14: "reservedChannelMsg14",
	reservedChannelMsg15: "reservedChannelMsg15",
	reservedChannelMsg16: "reservedChannelMsg16",

	MTC:        "MTC",
	SongSelect: "SongSelect",
	SPP:        "SPP",
	Tune:       "Tune",

	reservedSysCommon5:  "reservedSysCommon5",
	reservedSysCommon6:  "reservedSysCommon6",
	reservedSysCommon7:  "reservedSysCommon7",
	reservedSysCommon8:  "reservedSysCommon8",
	reservedSysCommon9:  "reservedSysCommon9",
	reservedSysCommon10: "reservedSysCommon10",
}

func (t Type) IsPlayable() bool {
	if t <= UnknownType {
		return false
	}

	return t < firstMeta
}

func (t Type) String() string {
	if s, has := typeNames[t]; has {
		return s
	}

	if t >= firstMeta {
		return typeNames[metaType]
	}

	return "user defined"
}

/*
func init() {
	fmt.Printf("RT1: %v SysC10: %v Mt30: %v\n", RT1, SysC10, Mt30)
}
*/

const (
	// UnknownMsg is an invalid or unknown MIDI message
	UnknownType Type = 0

	// RealTimeMsg is a MIDI realtime message. It can only be used over the wire.
	RealTimeType Type = -1

	// SysCommonMsg is a MIDI system common message. It can only be used over the wire.
	SysCommonType Type = -2

	// ChannelMsg is a MIDI channel message. It can be used in SMF and over the wire.
	ChannelType Type = -3

	// SysExMsg is a MIDI system exclusive message. It can be used in SMF and over the wire.
	SysExType Type = -4

	// metaMsg is a MIDI meta message (used in SMF = Simple MIDI Files)
	metaType Type = -5
)

const (

	// Tick is a MIDI tick realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	Tick Type = 1 + iota

	// TimingClock is a MIDI timing clock realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	TimingClock

	// Start is a MIDI start realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	Start

	// Continue is a MIDI continue realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	Continue

	// Stop is a MIDI stop realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	Stop

	// ActiveSense is a MIDI active sense realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	ActiveSense

	// Reset is a MIDI reset realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	Reset

	reservedRealTime8
	reservedRealTime9
	reservedRealTime10
	reservedRealTime11
	reservedRealTime12
	reservedRealTime13
	reservedRealTime14

	// NoteOn is a MIDI note on message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The velocity of a concrete Message of this type can be retrieved via the Velocity method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	NoteOn

	// NoteOff is a MIDI note off message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The velocity of a concrete Message of this type can be retrieved via the Velocity method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	NoteOff

	// ControlChange is a MIDI control change message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The controller of a concrete Message of this type can be retrieved via the Controller method of the Message.
	// The change of a concrete Message of this type can be retrieved via the Change method of the Message.
	ControlChange

	// PitchBend is a MIDI pitch bend message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The absolute and releative pitch of a concrete Message of this type can be retrieved via the Pitch method of the Message.
	PitchBend

	// AfterTouch is a MIDI after touch message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The pressure of a concrete Message of this type can be retrieved via the Pressure method of the Message.
	AfterTouch

	// PolyAfterTouch is a polyphonic MIDI after touch message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	// The pressure of a concrete Message of this type can be retrieved via the Pressure method of the Message.
	PolyAfterTouch

	// ProgramChange is a MIDI program change message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The program number of a concrete Message of this type can be retrieved via the Program method of the Message.
	ProgramChange

	reservedChannelMsg8
	reservedChannelMsg9
	reservedChannelMsg10
	reservedChannelMsg11
	reservedChannelMsg12
	reservedChannelMsg13
	reservedChannelMsg14
	reservedChannelMsg15
	reservedChannelMsg16

	// MTC is a MIDI MTC system common message.
	// TODO add method to Message to get the quarter frame and document it.
	MTC

	// SongSelect is a MIDI song select system common message.
	// TODO add method to Message to get the song number and document it.
	SongSelect

	// SPP is a MIDI song position pointer (SPP) system common message.
	// TODO add method to Message to get the song position pointer and document it.
	SPP

	// Tune is a MIDI tune request system common message.
	// There is no further data associated with messages of this type.
	Tune

	reservedSysCommon5
	reservedSysCommon6
	reservedSysCommon7
	reservedSysCommon8
	reservedSysCommon9
	reservedSysCommon10
)

const (
	// everything >= firstMeta are meta messages
	firstMeta Type = 70
)

/*
GetMsgType returns the message type for the given message (bytes that must include a status byte - no running status).

The returned MsgType will be a combination of message types, if appropriate (binary flags). For example:
A note on message on channel 0 will have a message type that is a combination of a ChannelMsg, a Channel0Msg, and a NoteOnMsg.
A tempo meta message of a SMF file will have a message type that is a combination of a MetaMsg, and a MetaTempoMsg.
*/
func GetType(bt []byte) (mType Type) {
	//fmt.Printf("GetMsgType % X\n", msg)

	byte1 := bt[0]

	switch {
	// channel/Voice Category Status
	case byte1 >= 0x80 && byte1 <= 0xEF:
		return GetChannelType(byte1)
	case byte1 == 0xF0, byte1 == 0xF7:
		// TODO what about sysex start stop etc.
		return SysExType
	case byte1 == 0xFF:
		/*
			if byte2 > 0 {
				return MetaMsgType
			}
		*/
		return GetRealtimeType(byte1)
	case byte1 < 0xF7:
		return GetSysCommonType(byte1)
	case byte1 > 0xF7:
		return GetRealtimeType(byte1)
	default:
		return UnknownType
	}
}

// GetChannelMsgType returns the MsgType of a channel message. It should not be used by the end consumer.
func GetChannelType(canary byte) (mType Type) {
	tp, _ := utils.ParseStatus(canary)

	switch tp {
	case 0xC:
		return ProgramChange
	case 0xD:
		return AfterTouch
	case 0x8:
		return NoteOff
	case 0x9:
		return NoteOn
	case 0xA:
		return PolyAfterTouch
	case 0xB:
		return ControlChange
	case 0xE:
		return PitchBend
	default:
		return UnknownType
	}
}

// GetRealtimeMsgType returns the MsgType of a realtime message. It should not be used by the end consumer.
func GetRealtimeType(b byte) Type {
	ty, has := rtMessages[b]
	if !has {
		return UnknownType
	}
	return ty
}

// GetSysCommonMsgType returns the MsgType of a sys common message. It should not be used by the end consumer.
func GetSysCommonType(b byte) Type {
	ty, has := syscommMessages[b]
	if !has {
		return UnknownType
	}
	return ty
}
