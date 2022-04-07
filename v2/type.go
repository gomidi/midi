package midi

import (
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// Type is the type of a midi message
type Type int8

// Is returns true, if the type correspond to the given type.
func (t Type) Is(checker Type) bool {

	switch {
	case t == UnknownMsg:
		return checker == UnknownMsg
	case t == SysExMsg:
		return checker == SysExMsg
	case t < UnknownMsg:
		return false
	case checker == UnknownMsg:
		return false
	case checker > UnknownMsg:
		return t == checker
	default:
		switch checker {
		case RealTimeMsg:
			return t <= reservedRealTimeMsg14
		case SysCommonMsg:
			return t >= MTCMsg && t <= reservedSysCommonMsg10
		case ChannelMsg:
			return t >= NoteOnMsg && t <= reservedChannelMsg16
		case metaMsg:
			return t >= firstMetaMsg
		default:
			return false
		}
	}
}

// AddTypeName adds names for new types that are not part of this package (e.g. meta types from the smf package).
// Don't use this function as a user, it is only internal to the library.
// Returns false, if the type already has been named, and true on success.
func AddTypeName(m Type, name string) bool {
	if _, has := typeNames[m]; has {
		return false
	}
	typeNames[m] = name
	return true
}

var typeNames = map[Type]string{

	UnknownMsg:   "UnknownType",
	RealTimeMsg:  "RealTimeType",
	SysCommonMsg: "SysCommonType",
	ChannelMsg:   "ChannelType",
	SysExMsg:     "SysExType",
	//metaMsg:      "MetaType",

	TickMsg:        "Tick",
	TimingClockMsg: "TimingClock",
	StartMsg:       "Start",
	ContinueMsg:    "Continue",
	StopMsg:        "Stop",
	ActiveSenseMsg: "ActiveSense",
	ResetMsg:       "Reset",

	reservedRealTimeMsg8:  "reservedRealTime8",
	reservedRealTimeMsg9:  "reservedRealTime9",
	reservedRealTimeMsg10: "reservedRealTime10",
	reservedRealTimeMsg11: "reservedRealTime11",
	reservedRealTimeMsg12: "reservedRealTime12",
	reservedRealTimeMsg13: "reservedRealTime13",
	reservedRealTimeMsg14: "reservedRealTime14",

	NoteOnMsg:         "NoteOn",
	NoteOffMsg:        "NoteOff",
	ControlChangeMsg:  "ControlChange",
	PitchBendMsg:      "PitchBend",
	AfterTouchMsg:     "AfterTouch",
	PolyAfterTouchMsg: "PolyAfterTouch",
	ProgramChangeMsg:  "ProgramChange",

	reservedChannelMsg8:  "reservedChannelMsg8",
	reservedChannelMsg9:  "reservedChannelMsg9",
	reservedChannelMsg10: "reservedChannelMsg10",
	reservedChannelMsg11: "reservedChannelMsg11",
	reservedChannelMsg12: "reservedChannelMsg12",
	reservedChannelMsg13: "reservedChannelMsg13",
	reservedChannelMsg14: "reservedChannelMsg14",
	reservedChannelMsg15: "reservedChannelMsg15",
	reservedChannelMsg16: "reservedChannelMsg16",

	MTCMsg:        "MTC",
	SongSelectMsg: "SongSelect",
	SPPMsg:        "SPP",
	TuneMsg:       "Tune",

	reservedSysCommonMsg5:  "reservedSysCommon5",
	reservedSysCommonMsg6:  "reservedSysCommon6",
	reservedSysCommonMsg7:  "reservedSysCommon7",
	reservedSysCommonMsg8:  "reservedSysCommon8",
	reservedSysCommonMsg9:  "reservedSysCommon9",
	reservedSysCommonMsg10: "reservedSysCommon10",
}

// String returns the name of the type.
func (t Type) String() string {
	if s, has := typeNames[t]; has {
		return s
	}

	if t >= firstMetaMsg {
		return typeNames[metaMsg]
	}

	return "user defined"
}

const (
	// UnknownMsg is an invalid or unknown MIDI message
	UnknownMsg Type = 0

	// RealTimeMsg is a MIDI realtime message. It can only be used over the wire.
	RealTimeMsg Type = -1

	// SysCommonMsg is a MIDI system common message. It can only be used over the wire.
	SysCommonMsg Type = -2

	// ChannelMsg is a MIDI channel message. It can be used in SMF and over the wire.
	ChannelMsg Type = -3

	// SysExMsg is a MIDI system exclusive message. It can be used in SMF and over the wire.
	SysExMsg Type = -4

	// metaMsg is a MIDI meta message (used in SMF = Simple MIDI Files)
	metaMsg Type = -5
)

const (

	// Tick is a MIDI tick realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	TickMsg Type = 1 + iota

	// TimingClock is a MIDI timing clock realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	TimingClockMsg

	// Start is a MIDI start realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	StartMsg

	// Continue is a MIDI continue realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	ContinueMsg

	// Stop is a MIDI stop realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	StopMsg

	// ActiveSense is a MIDI active sense realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	ActiveSenseMsg

	// Reset is a MIDI reset realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	ResetMsg

	reservedRealTimeMsg8
	reservedRealTimeMsg9
	reservedRealTimeMsg10
	reservedRealTimeMsg11
	reservedRealTimeMsg12
	reservedRealTimeMsg13
	reservedRealTimeMsg14

	// NoteOn is a MIDI note on message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The velocity of a concrete Message of this type can be retrieved via the Velocity method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	NoteOnMsg

	// NoteOff is a MIDI note off message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The velocity of a concrete Message of this type can be retrieved via the Velocity method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	NoteOffMsg

	// ControlChange is a MIDI control change message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The controller of a concrete Message of this type can be retrieved via the Controller method of the Message.
	// The change of a concrete Message of this type can be retrieved via the Change method of the Message.
	ControlChangeMsg

	// PitchBend is a MIDI pitch bend message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The absolute and releative pitch of a concrete Message of this type can be retrieved via the Pitch method of the Message.
	PitchBendMsg

	// AfterTouch is a MIDI after touch message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The pressure of a concrete Message of this type can be retrieved via the Pressure method of the Message.
	AfterTouchMsg

	// PolyAfterTouch is a polyphonic MIDI after touch message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	// The pressure of a concrete Message of this type can be retrieved via the Pressure method of the Message.
	PolyAfterTouchMsg

	// ProgramChange is a MIDI program change message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The program number of a concrete Message of this type can be retrieved via the Program method of the Message.
	ProgramChangeMsg

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
	MTCMsg

	// SongSelect is a MIDI song select system common message.
	// TODO add method to Message to get the song number and document it.
	SongSelectMsg

	// SPP is a MIDI song position pointer (SPP) system common message.
	// TODO add method to Message to get the song position pointer and document it.
	SPPMsg

	// Tune is a MIDI tune request system common message.
	// There is no further data associated with messages of this type.
	TuneMsg

	reservedSysCommonMsg5
	reservedSysCommonMsg6
	reservedSysCommonMsg7
	reservedSysCommonMsg8
	reservedSysCommonMsg9
	reservedSysCommonMsg10
)

const (
	// everything >= firstMeta are meta messages
	firstMetaMsg Type = 70
)

/*
getType returns the message type for the given message (bytes that must include a status byte - no running status).
*/
func getType(bt []byte) (mType Type) {
	//fmt.Printf("GetMsgType % X\n", msg)
	if len(bt) == 0 {
		return UnknownMsg
	}
	byte1 := bt[0]

	switch {
	// channel/Voice Category Status
	case byte1 >= 0x80 && byte1 <= 0xEF:
		return getChannelType(byte1)
	case byte1 == 0xF0, byte1 == 0xF7:
		// TODO what about sysex start stop etc.
		return SysExMsg
	case byte1 == 0xFF:
		/*
			if byte2 > 0 {
				return MetaMsgType
			}
		*/
		return getRealtimeType(byte1)
	case byte1 < 0xF7:
		return getSysCommonType(byte1)
	case byte1 > 0xF7:
		return getRealtimeType(byte1)
	default:
		return UnknownMsg
	}
}

// getChannelType returns the MsgType of a channel message. It should not be used by the end consumer.
func getChannelType(canary byte) (mType Type) {
	tp, _ := utils.ParseStatus(canary)

	switch tp {
	case 0xC:
		return ProgramChangeMsg
	case 0xD:
		return AfterTouchMsg
	case 0x8:
		return NoteOffMsg
	case 0x9:
		return NoteOnMsg
	case 0xA:
		return PolyAfterTouchMsg
	case 0xB:
		return ControlChangeMsg
	case 0xE:
		return PitchBendMsg
	default:
		return UnknownMsg
	}
}

// getRealtimeMsgType returns the MsgType of a realtime message. It should not be used by the end consumer.
func getRealtimeType(b byte) Type {
	ty, has := rtMessages[b]
	if !has {
		return UnknownMsg
	}
	return ty
}

// getSysCommonMsgType returns the MsgType of a sys common message. It should not be used by the end consumer.
func getSysCommonType(b byte) Type {
	ty, has := syscommMessages[b]
	if !has {
		return UnknownMsg
	}
	return ty
}
