package midi

import (
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

type MsgKind uint8

const (
	// UnknownMsg represents every MIDI message that is invalid or unknown.
	// There is no further data associated with messages of this type.
	UnknownMsg MsgKind = 0

	ChannelMsg MsgKind = 1

	// RealTimeMsg is a MIDI realtime message. It can only be used over the wire.
	RealTimeMsg MsgKind = 2

	// SysCommonMsg is a MIDI system common message. It can only be used over the wire.
	SysCommonMsg MsgKind = 3

	// SysExMsg is a MIDI system exclusive message. It can be used in SMF and over the wire.
	SysExMsg MsgKind = 4

	// MetaMsg is a MIDI meta message (used in SMF = Simple MIDI Files)
	MetaMsg MsgKind = 5

	// a way for the user to define his own message types (based on the application)
	UserDefinedMsg = 6
)

func (k MsgKind) String() string {
	switch k {
	case UnknownMsg:
		return "UnknownMsg"
	case ChannelMsg:
		return "ChannelMsg"
	case RealTimeMsg:
		return "RealTimeMsg"
	case SysCommonMsg:
		return "SysCommonMsg"
	case SysExMsg:
		return "SysExMsg"
	case MetaMsg:
		return "MetaMsg"
	default:
		return "UserDefinedMsg"
	}
}

// 39 MessageTypes: uint8 genügt (256); im grunde genügen 6 bits (64)
/*
1 unknown message (0)
7 realtime messages (1-8)
4 syscommon messages (9-13)
1 sysex message (14)
7 channel messages
------------------------
20

da kämen wir mit 32, also 5 bits aus

+19 Meta messages

wir sind großzuegig: 256 typen

1 unknown message (0)
15 realtime messages (1-15)
8 syscommon messages (16-23)
16 channel messages (24-39)
32 meta messages (40-71)
64 sysex typen (72-135)
121 eigene freie typen (136-255)
*/

//type MsgType uint64
//type MsgType uint8
type MsgType uint32

func (m MsgType) Kind() MsgKind {
	switch {
	case m == 0:
		return UnknownMsg
	case m <= ProgramChangeMsg:
		return ChannelMsg
	case m <= ResetMsg:
		return RealTimeMsg
	case m <= TuneMsg:
		return SysCommonMsg
		//	case m == MetaMsgType:
	//	return MetaMsg
	case m == SysExMsgType:
		return SysExMsg
	default:
		return UserDefinedMsg
	}
}

/*
func (m MessageType) IsOneOf(ts ...MessageType) bool {
	for _, t := range ts {
		if m.Is(t) {
			return true
		}
	}
	return false
}
*/

//const UnknownMsg MsgType = 0

const (
	// NoteOnMsg is a MIDI note on message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The velocity of a concrete Message of this type can be retrieved via the Velocity method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	NoteOnMsg MsgType = 1 << iota

	// NoteOffMsg is a MIDI note off message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The velocity of a concrete Message of this type can be retrieved via the Velocity method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	NoteOffMsg

	// ControlChangeMsg is a MIDI control change message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The controller of a concrete Message of this type can be retrieved via the Controller method of the Message.
	// The change of a concrete Message of this type can be retrieved via the Change method of the Message.
	ControlChangeMsg

	// PitchBendMsg is a MIDI pitch bend message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The absolute and releative pitch of a concrete Message of this type can be retrieved via the Pitch method of the Message.
	PitchBendMsg

	// AfterTouchMsg is a MIDI after touch message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The pressure of a concrete Message of this type can be retrieved via the Pressure method of the Message.
	AfterTouchMsg

	// PolyAfterTouchMsg is a polyphonic MIDI after touch message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	// The pressure of a concrete Message of this type can be retrieved via the Pressure method of the Message.
	PolyAfterTouchMsg

	// ProgramChangeMsg is a MIDI program change message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The program number of a concrete Message of this type can be retrieved via the Program method of the Message.
	ProgramChangeMsg

	// TimingClockMsg is a MIDI timing clock realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	TimingClockMsg

	// TickMsg is a MIDI tick realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	TickMsg

	// StartMsg is a MIDI start realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	StartMsg

	// ContinueMsg is a MIDI continue realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	ContinueMsg

	// StopMsg is a MIDI stop realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	StopMsg

	// ActiveSenseMsg is a MIDI active sense realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	ActiveSenseMsg

	// ResetMsg is a MIDI reset realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	ResetMsg

	/*
		SysExStartMsg
		SysExEndMsg
		SysExCompleteMsg
		SysExEscapeMsg
		SysExContinueMsg
	*/

	// MTCMsg is a MIDI MTC system common message (which is a SysCommonMsg).
	// TODO add method to Message to get the quarter frame and document it.
	MTCMsg

	// SongSelectMsg is a MIDI song select system common message (which is a SysCommonMsg).
	// TODO add method to Message to get the song number and document it.
	SongSelectMsg

	// SPPMsg is a MIDI song position pointer (SPP) system common message (which is a SysCommonMsg).
	// TODO add method to Message to get the song position pointer and document it.
	SPPMsg

	// SPPMsg is a MIDI tune request system common message (which is a SysCommonMsg).
	// There is no further data associated with messages of this type.
	TuneMsg

	SysExMsgType

	// UndefinedMsg is an undefined MIDI message.
	UndefinedMsgType

	// MetaMsg is a meta message
	//MetaMsgType
)

// NoteMsg is either a NoteOnMsg or a NoteOffMsg.
const NoteMsg = NoteOnMsg | NoteOffMsg

var msgTypeString = map[MsgType]string{
	//	MetaMsgType:       "MetaMsgType",
	SysExMsgType:      "SysExMsgType",
	NoteOnMsg:         "NoteOnMsg",
	NoteOffMsg:        "NoteOffMsg",
	ControlChangeMsg:  "ControlChangeMsg",
	PitchBendMsg:      "PitchBendMsg",
	AfterTouchMsg:     "AfterTouchMsg",
	PolyAfterTouchMsg: "PolyAfterTouchMsg",
	ProgramChangeMsg:  "ProgramChangeMsg",
	TimingClockMsg:    "TimingClockMsg",
	TickMsg:           "TickMsg",
	StartMsg:          "StartMsg",
	ContinueMsg:       "ContinueMsg",
	StopMsg:           "StopMsg",
	ActiveSenseMsg:    "ActiveSenseMsg",
	ResetMsg:          "ResetMsg",
	MTCMsg:            "MTCMsg",
	SongSelectMsg:     "SongSelectMsg",
	SPPMsg:            "SPPMsg",
	UndefinedMsgType:  "UndefinedMsgType",
	TuneMsg:           "TuneMsg",
	//	UnknownMsg:    "UnknownMsg",
}

/*
GetMsgType returns the message type for the given message (bytes that must include a status byte - no running status).

The returned MsgType will be a combination of message types, if appropriate (binary flags). For example:
A note on message on channel 0 will have a message type that is a combination of a ChannelMsg, a Channel0Msg, and a NoteOnMsg.
A tempo meta message of a SMF file will have a message type that is a combination of a MetaMsg, and a MetaTempoMsg.
*/
func GetMsgType(bt []byte) (mType MsgType) {
	//fmt.Printf("GetMsgType % X\n", msg)

	byte1 := bt[0]

	switch {
	// channel/Voice Category Status
	case byte1 >= 0x80 && byte1 <= 0xEF:
		return GetChannelMsgType(byte1)
	case byte1 == 0xF0, byte1 == 0xF7:
		// TODO what about sysex start stop etc.
		return SysExMsgType
	case byte1 == 0xFF:
		/*
			if byte2 > 0 {
				return MetaMsgType
			}
		*/
		return GetRealtimeMsgType(byte1)
	case byte1 < 0xF7:
		return GetSysCommonMsgType(byte1)
	case byte1 > 0xF7:
		return GetRealtimeMsgType(byte1)
	default:
		return UndefinedMsgType
	}
}

// Set adds the given message type to the existing message type and returns a combination (via binary flags)
//func (m MsgType) Set(flag MsgType) MsgType { return m | flag }

// Clear removes the given message type from the combination of messages types (via binary flags)
//func (m MsgType) Clear(flag MsgType) MsgType { return m &^ flag }

// Toggle toggles wether or not the given message type is set (via binary flags)
//func (m MsgType) Toggle(flag MsgType) MsgType { return m ^ flag }

// Is returns if the given message type is part of the combination of message types
func Is(mt1, mt2 MessageType) bool {
	return mt1.Kind() == mt2.Kind() && mt1.Val()&mt2.Val() != 0
}

/*
// IsOneOf returns true if one of the given message types is set.
func (m MsgType) IsOneOf(flags ...MsgType) bool {
	for _, fl := range flags {
		if fl.Kind() == m.Kind() && m&fl != 0 {
			return true
		}
	}
	return false
}
*/

/*
// IsAllOf returns true if all of the given message types are set.
func (m MsgType) IsAllOf(flags ...MsgType) bool {
	for _, fl := range flags {
		if m&fl == 0 {
			return false
		}
	}
	return true
}
*/

func (m MsgType) Val() uint32 {
	return uint32(m)
}

// String returns a string that represents the message type
func (m MsgType) String() string {
	s, has := msgTypeString[m]
	if !has {
		return "UndefinedMsgType"
	}

	return s

	/*
		switch m.Kind() {
		case SysExMsg:
			return "SysExMsgType"
		case MetaMsg:
			return "MetaMsgType"
		case SysCommonMsg:
			return "SysCommonMsg"
		case RealTimeMsg:
			return "RealTimeMsg"
		case ChannelMsg:
			return "ChannelMsg"
		default:
			return "UndefinedMsgType"
		}
	*/

	/*
		if m.Is(ChannelMsg) {
			var clCh MsgType

			if m.Is(Channel0Msg) {
				clCh = Channel0Msg
			}

			if m.Is(Channel1Msg) {
				clCh = Channel1Msg
			}

			if m.Is(Channel2Msg) {
				clCh = Channel2Msg
			}

			if m.Is(Channel3Msg) {
				clCh = Channel3Msg
			}

			if m.Is(Channel4Msg) {
				clCh = Channel4Msg
			}

			if m.Is(Channel5Msg) {
				clCh = Channel5Msg
			}

			if m.Is(Channel6Msg) {
				clCh = Channel6Msg
			}

			if m.Is(Channel7Msg) {
				clCh = Channel7Msg
			}

			if m.Is(Channel8Msg) {
				clCh = Channel8Msg
			}

			if m.Is(Channel9Msg) {
				clCh = Channel9Msg
			}

			if m.Is(Channel10Msg) {
				clCh = Channel10Msg
			}

			if m.Is(Channel11Msg) {
				clCh = Channel11Msg
			}

			if m.Is(Channel12Msg) {
				clCh = Channel12Msg
			}

			if m.Is(Channel13Msg) {
				clCh = Channel13Msg
			}

			if m.Is(Channel14Msg) {
				clCh = Channel14Msg
			}

			if m.Is(Channel15Msg) {
				clCh = Channel15Msg
			}

			return msgTypeString[clCh] + " & " + msgTypeString[m.Clear(ChannelMsg).Clear(clCh)]
		}

		return "Unknown"
	*/
}

// GetChannelMsgType returns the MsgType of a channel message. It should not be used by the end consumer.
func GetChannelMsgType(canary byte) (mType MsgType) {
	tp, _ := utils.ParseStatus(canary)
	/*
		var sType MsgType
		//r.status = canary
		mType = mType.Set(ChannelMsg)
		var ctype MsgType

		switch ch {
		case 0:
			ctype = Channel0Msg
		case 1:
			ctype = Channel1Msg
		case 2:
			ctype = Channel2Msg
		case 3:
			ctype = Channel3Msg
		case 4:
			ctype = Channel4Msg
		case 5:
			ctype = Channel5Msg
		case 6:
			ctype = Channel6Msg
		case 7:
			ctype = Channel7Msg
		case 8:
			ctype = Channel8Msg
		case 9:
			ctype = Channel9Msg
		case 10:
			ctype = Channel10Msg
		case 11:
			ctype = Channel11Msg
		case 12:
			ctype = Channel12Msg
		case 13:
			ctype = Channel13Msg
		case 14:
			ctype = Channel14Msg
		case 15:
			ctype = Channel15Msg
		}

		mType = mType.Set(ctype)
	*/

	//	fmt.Printf("status: % X\n", tp)

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
		return UndefinedMsgType
	}

	//mType = mType.Set(sType)
	//return mType
}

// GetRealtimeMsgType returns the MsgType of a realtime message. It should not be used by the end consumer.
func GetRealtimeMsgType(b byte) MsgType {
	ty, has := rtMessages[b]
	if !has {
		return UndefinedMsgType
	}
	return ty
}

// GetSysCommonMsgType returns the MsgType of a sys common message. It should not be used by the end consumer.
func GetSysCommonMsgType(b byte) MsgType {
	ty, has := syscommMessages[b]
	if !has {
		return UndefinedMsgType
	}
	return ty
}

/*
// GetMetaMsgType returns the MsgType of a meta message. It should not be used by the end consumer.
func GetMetaMsgType(b byte) MsgType {
	return metaMessages[b]
}
*/
