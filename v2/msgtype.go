package midi

import (
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

type MessageCategory uint8

const (
	// UnknownMessages are invalid or unknown MIDI messages.
	UnknownMessages MessageCategory = 0

	// ChannelMessages are MIDI channel messages. They can be used in SMF and over the wire.
	ChannelMessages MessageCategory = 1

	// RealTimeMessages are MIDI realtime messages. They can only be used over the wire.
	RealTimeMessages MessageCategory = 2

	// SysCommonMessages are MIDI system common messages. They can only be used over the wire.
	SysCommonMessages MessageCategory = 3

	// SysExMessages are MIDI system exclusive messages. They can be used in SMF and over the wire.
	SysExMessages MessageCategory = 4

	// MetaMessages are MIDI meta messages (used in SMF = Simple MIDI Files)
	MetaMessages MessageCategory = 5

	// UserDefinedMessages can be defined by the user (based on the application)
	UserDefinedMessages MessageCategory = 6
)

func (k MessageCategory) String() string {
	switch k {
	case UnknownMessages:
		return "UnknownMessages"
	case ChannelMessages:
		return "ChannelMessages"
	case RealTimeMessages:
		return "RealTimeMessages"
	case SysCommonMessages:
		return "SysCommonMessages"
	case SysExMessages:
		return "SysExMessages"
	case MetaMessages:
		return "MetaMessages"
	case UserDefinedMessages:
		return "UserDefinedMessages"
	default:
		return "UnknownMessages"
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

/*
unknown message (0)
realtime message kind (1)
14 realtime messages (2-15)
syscommon message kind (16)
7 syscommon messages (17-23)
channel message kind (24)
16 channel messages (25-39)
meta message kind (40)
31 meta messages (41-71)
sysex message kind (72)
63 sysex typen (73-135)
user defined kind (136)
120 eigene freie typen (137-255)
*/

//type MsgType uint64
//type MsgType uint8
type MsgType uint32

func (m MsgType) Category() MessageCategory {
	switch {
	case m == 0:
		return UnknownMessages
	case m <= ProgramChange:
		return ChannelMessages
	case m <= Reset:
		return RealTimeMessages
	case m <= Tune:
		return SysCommonMessages
		//	case m == MetaMsgType:
	//	return MetaMsg
	case m == SysEx:
		return SysExMessages
	default:
		return UserDefinedMessages
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
	// NoteOn is a MIDI note on message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The velocity of a concrete Message of this type can be retrieved via the Velocity method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	NoteOn MsgType = 1 << iota

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

	// TimingClock is a MIDI timing clock realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	TimingClock

	// Tick is a MIDI tick realtime message (which is a RealTimeMsg).
	// There is no further data associated with messages of this type.
	Tick

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

	/*
		SysExStartMsg
		SysExEndMsg
		SysExCompleteMsg
		SysExEscapeMsg
		SysExContinueMsg
	*/

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

	// SysEx is a MIDI system exclusive message
	SysEx

	// Undefined is an undefined MIDI message.
	Undefined

	// MetaMsg is a meta message
	//MetaMsgType
)

// Note is either a NoteOn or a NoteOff.
const Note = NoteOn | NoteOff

var msgTypeString = map[MsgType]string{
	//	MetaMsgType:       "MetaMsgType",
	SysEx:          "SysEx",
	NoteOn:         "NoteOn",
	NoteOff:        "NoteOff",
	ControlChange:  "ControlChange",
	PitchBend:      "PitchBend",
	AfterTouch:     "AfterTouch",
	PolyAfterTouch: "PolyAfterTouch",
	ProgramChange:  "ProgramChange",
	TimingClock:    "TimingClock",
	Tick:           "Tick",
	Start:          "Start",
	Continue:       "Continue",
	Stop:           "Stop",
	ActiveSense:    "ActiveSense",
	Reset:          "Reset",
	MTC:            "MTC",
	SongSelect:     "SongSelect",
	SPP:            "SPP",
	Undefined:      "Undefined",
	Tune:           "Tune",
	//	UnknownMsg:    "Unknown",
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
		return SysEx
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
		return Undefined
	}
}

// Set adds the given message type to the existing message type and returns a combination (via binary flags)
//func (m MsgType) Set(flag MsgType) MsgType { return m | flag }

// Clear removes the given message type from the combination of messages types (via binary flags)
//func (m MsgType) Clear(flag MsgType) MsgType { return m &^ flag }

// Toggle toggles wether or not the given message type is set (via binary flags)
//func (m MsgType) Toggle(flag MsgType) MsgType { return m ^ flag }

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
		return Undefined
	}

	//mType = mType.Set(sType)
	//return mType
}

// GetRealtimeMsgType returns the MsgType of a realtime message. It should not be used by the end consumer.
func GetRealtimeMsgType(b byte) MsgType {
	ty, has := rtMessages[b]
	if !has {
		return Undefined
	}
	return ty
}

// GetSysCommonMsgType returns the MsgType of a sys common message. It should not be used by the end consumer.
func GetSysCommonMsgType(b byte) MsgType {
	ty, has := syscommMessages[b]
	if !has {
		return Undefined
	}
	return ty
}

/*
// GetMetaMsgType returns the MsgType of a meta message. It should not be used by the end consumer.
func GetMetaMsgType(b byte) MsgType {
	return metaMessages[b]
}
*/
