package midi

import (
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

type MsgType uint64

// UnknownMsg represents every MIDI message that is invalid or unknown.
// There is no further data associated with messages of this type.
const UnknownMsg MsgType = 0

const (

	// ChannelMsg is a MIDI channel message. It can be used in SMF and over the wire.
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	ChannelMsg MsgType = 1 << iota

	// MetaMsg is a MIDI meta message (used in SMF = Simple MIDI Files)
	MetaMsg

	// RealTimeMsg is a MIDI realtime message. It can only be used over the wire.
	RealTimeMsg

	// SysCommonMsg is a MIDI system common message. It can only be used over the wire.
	SysCommonMsg

	// SysExMsg is a MIDI system exclusive message. It can be used in SMF and over the wire.
	SysExMsg

	// NoteOnMsg is a MIDI note on message (which is a ChannelMsg).
	// The channel of a concrete Message of this type can be retrieved via the Channel method of the Message.
	// The velocity of a concrete Message of this type can be retrieved via the Velocity method of the Message.
	// The key of a concrete Message of this type can be retrieved via the Key method of the Message.
	NoteOnMsg

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

	// MetaChannelMsg is a MIDI channel meta message (which is a MetaMsg).
	// TODO add method to Message to get the channel number and document it.
	MetaChannelMsg

	// MetaCopyrightMsg is a MIDI copyright meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaCopyrightMsg

	// MetaCuepointMsg is a MIDI cuepoint meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaCuepointMsg

	// MetaDeviceMsg is a MIDI device meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaDeviceMsg

	// MetaEndOfTrackMsg is a MIDI end of track meta message (which is a MetaMsg).
	MetaEndOfTrackMsg

	// MetaInstrumentMsg is a MIDI instrument meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaInstrumentMsg

	// MetaKeySigMsg is a MIDI key signature meta message (which is a MetaMsg).
	// TODO add method to Message to get the key signature and document it.
	MetaKeySigMsg

	// MetaLyricMsg is a MIDI lyrics meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaLyricMsg

	// MetaTextMsg is a MIDI text meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaTextMsg

	// MetaMarkerMsg is a MIDI marker meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaMarkerMsg

	// MetaPortMsg is a MIDI port meta message (which is a MetaMsg).
	// TODO add method to Message to get the port number and document it.
	MetaPortMsg

	// MetaSeqNumberMsg is a MIDI sequencer number meta message (which is a MetaMsg).
	// TODO add method to Message to get the sequence number and document it.
	MetaSeqNumberMsg

	// MetaSeqDataMsg is a MIDI sequencer data meta message (which is a MetaMsg).
	// TODO add method to Message to get the sequencer data and document it.
	MetaSeqDataMsg

	// MetaTempoMsg is a MIDI tempo meta message (which is a MetaMsg).
	// The tempo in beats per minute of a concrete Message of this type can be retrieved via the BPM method of the Message.
	MetaTempoMsg

	// MetaTimeSigMsg is a MIDI time signature meta message (which is a MetaMsg).
	// The numerator, denominator, clocksPerClick and demiSemiQuaverPerQuarter of a concrete Message of this type can be retrieved via the TimeSig method of the Message.
	// A more comfortable way to get the meter is to use the Meter method of the Message.
	MetaTimeSigMsg

	// MetaTrackNameMsg is a MIDI track name meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaTrackNameMsg

	// MetaSMPTEOffsetMsg is a MIDI smpte offset meta message (which is a MetaMsg).
	// TODO add method to Message to get the smpte offset and document it.
	MetaSMPTEOffsetMsg

	// MetaUndefinedMsg is an undefined MIDI meta message (which is a MetaMsg).
	MetaUndefinedMsg

	// MetaProgramNameMsg is a MIDI program name meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaProgramNameMsg

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

	// Channel0Msg is a MIDI channel message on channel 0 (which is a ChannelMsg).
	Channel0Msg

	// Channel1Msg is a MIDI channel message on channel 1 (which is a ChannelMsg).
	Channel1Msg

	// Channel2Msg is a MIDI channel message on channel 2 (which is a ChannelMsg).
	Channel2Msg

	// Channel3Msg is a MIDI channel message on channel 3 (which is a ChannelMsg).
	Channel3Msg

	// Channel4Msg is a MIDI channel message on channel 4 (which is a ChannelMsg).
	Channel4Msg

	// Channel5Msg is a MIDI channel message on channel 5 (which is a ChannelMsg).
	Channel5Msg

	// Channel6Msg is a MIDI channel message on channel 6 (which is a ChannelMsg).
	Channel6Msg

	// Channel7Msg is a MIDI channel message on channel 7 (which is a ChannelMsg).
	Channel7Msg

	// Channel8Msg is a MIDI channel message on channel 8 (which is a ChannelMsg).
	Channel8Msg

	// Channel9Msg is a MIDI channel message on channel 9 (which is a ChannelMsg).
	Channel9Msg

	// Channel10Msg is a MIDI channel message on channel 10 (which is a ChannelMsg).
	Channel10Msg

	// Channel11Msg is a MIDI channel message on channel 11 (which is a ChannelMsg).
	Channel11Msg

	// Channel12Msg is a MIDI channel message on channel 12 (which is a ChannelMsg).
	Channel12Msg

	// Channel13Msg is a MIDI channel message on channel 13 (which is a ChannelMsg).
	Channel13Msg

	// Channel14Msg is a MIDI channel message on channel 14 (which is a ChannelMsg).
	Channel14Msg

	// Channel15Msg is a MIDI channel message on channel 15 (which is a ChannelMsg).
	Channel15Msg

	// UndefinedMsg is an undefined MIDI message.
	UndefinedMsg
)

// NoteMsg is either a NoteOnMsg or a NoteOffMsg.
const NoteMsg = NoteOnMsg | NoteOffMsg

var channelType = map[uint8]MsgType{
	0:  Channel0Msg,
	1:  Channel1Msg,
	2:  Channel2Msg,
	3:  Channel3Msg,
	4:  Channel4Msg,
	5:  Channel5Msg,
	6:  Channel6Msg,
	7:  Channel7Msg,
	8:  Channel8Msg,
	9:  Channel9Msg,
	10: Channel10Msg,
	11: Channel11Msg,
	12: Channel12Msg,
	13: Channel13Msg,
	14: Channel14Msg,
	15: Channel15Msg,
}

var msgTypeString = map[MsgType]string{
	ChannelMsg:         "ChannelMsg",
	MetaMsg:            "MetaMsg",
	RealTimeMsg:        "RealTimeMsg",
	SysCommonMsg:       "SysCommonMsg",
	SysExMsg:           "SysExMsg",
	NoteOnMsg:          "NoteOnMsg",
	NoteOffMsg:         "NoteOffMsg",
	ControlChangeMsg:   "ControlChangeMsg",
	PitchBendMsg:       "PitchBendMsg",
	AfterTouchMsg:      "AfterTouchMsg",
	PolyAfterTouchMsg:  "PolyAfterTouchMsg",
	ProgramChangeMsg:   "ProgramChangeMsg",
	MetaChannelMsg:     "MetaChannelMsg",
	MetaCopyrightMsg:   "MetaCopyrightMsg",
	MetaCuepointMsg:    "MetaCuepointMsg",
	MetaDeviceMsg:      "MetaDeviceMsg",
	MetaEndOfTrackMsg:  "MetaEndOfTrackMsg",
	MetaInstrumentMsg:  "MetaInstrumentMsg",
	MetaKeySigMsg:      "MetaKeySigMsg",
	MetaLyricMsg:       "MetaLyricMsg",
	MetaTextMsg:        "MetaTextMsg",
	MetaMarkerMsg:      "MetaMarkerMsg",
	MetaPortMsg:        "MetaPortMsg",
	MetaSeqNumberMsg:   "MetaSeqNumberMsg",
	MetaSeqDataMsg:     "MetaSeqDataMsg",
	MetaTempoMsg:       "MetaTempoMsg",
	MetaTimeSigMsg:     "MetaTimeSigMsg",
	MetaTrackNameMsg:   "MetaTrackNameMsg",
	MetaSMPTEOffsetMsg: "MetaSMPTEOffsetMsg",
	MetaUndefinedMsg:   "MetaUndefinedMsg",
	MetaProgramNameMsg: "MetaProgramNameMsg",
	TimingClockMsg:     "TimingClockMsg",
	TickMsg:            "TickMsg",
	StartMsg:           "StartMsg",
	ContinueMsg:        "ContinueMsg",
	StopMsg:            "StopMsg",
	ActiveSenseMsg:     "ActiveSenseMsg",
	ResetMsg:           "ResetMsg",
	/*
		SysExStartMsg:      "SysExStartMsg",
		SysExEndMsg:        "SysExEndMsg",
		SysExCompleteMsg:   "SysExCompleteMsg",
		SysExEscapeMsg:     "SysExEscapeMsg",
		SysExContinueMsg:   "SysExContinueMsg",
	*/
	MTCMsg:        "MTCMsg",
	SongSelectMsg: "SongSelectMsg",
	SPPMsg:        "SPPMsg",
	UndefinedMsg:  "UndefinedMsg",
	TuneMsg:       "TuneMsg",
	UnknownMsg:    "UnknownMsg",
	Channel0Msg:   "Channel0Msg",
	Channel1Msg:   "Channel1Msg",
	Channel2Msg:   "Channel2Msg",
	Channel3Msg:   "Channel3Msg",
	Channel4Msg:   "Channel4Msg",
	Channel5Msg:   "Channel5Msg",
	Channel6Msg:   "Channel6Msg",
	Channel7Msg:   "Channel7Msg",
	Channel8Msg:   "Channel8Msg",
	Channel9Msg:   "Channel9Msg",
	Channel10Msg:  "Channel10Msg",
	Channel11Msg:  "Channel11Msg",
	Channel12Msg:  "Channel12Msg",
	Channel13Msg:  "Channel13Msg",
	Channel14Msg:  "Channel14Msg",
	Channel15Msg:  "Channel15Msg",
}

/*
GetMsgType returns the message type for the given message (bytes that must include a status byte - no running status).

The returned MsgType will be a combination of message types, if appropriate (binary flags). For example:
A note on message on channel 0 will have a message type that is a combination of a ChannelMsg, a Channel0Msg, and a NoteOnMsg.
A tempo meta message of a SMF file will have a message type that is a combination of a MetaMsg, and a MetaTempoMsg.
*/
func GetMsgType(msg []byte) (mType MsgType) {
	if len(msg) == 0 {
		return UnknownMsg
	}

	var canary = msg[0]

	switch {
	// channel/Voice Category Status
	case canary >= 0x80 && canary <= 0xEF:
		return GetChannelMsgType(canary)
	case canary == 0xF0, canary == 0xF7:
		// TODO what about sysex start stop etc.
		return SysExMsg
	case canary == 0xFF:
		if len(msg) > 1 {
			return GetMetaMsgType(msg[1])
		}
		return GetRealtimeMsgType(canary)
	case canary < 0xF7:
		return GetSysCommonMsgType(canary)
	case canary > 0xF7:
		return GetRealtimeMsgType(canary)
	default:
		return UnknownMsg
	}
}

// Set adds the given message type to the existing message type and returns a combination (via binary flags)
func (m MsgType) Set(flag MsgType) MsgType { return m | flag }

// Clear removes the given message type from the combination of messages types (via binary flags)
func (m MsgType) Clear(flag MsgType) MsgType { return m &^ flag }

// Toggle toggles wether or not the given message type is set (via binary flags)
func (m MsgType) Toggle(flag MsgType) MsgType { return m ^ flag }

// Is returns if the given message type is part of the combination of message types
func (m MsgType) Is(flag MsgType) bool { return m&flag != 0 }

// IsOneOf returns true if one of the given message types is set.
func (m MsgType) IsOneOf(flags ...MsgType) bool {
	for _, fl := range flags {
		if m&fl != 0 {
			return true
		}
	}
	return false
}

// IsAllOf returns true if all of the given message types are set.
func (m MsgType) IsAllOf(flags ...MsgType) bool {
	for _, fl := range flags {
		if m&fl == 0 {
			return false
		}
	}
	return true
}

// String returns a string that represents the message type
func (m MsgType) String() string {
	//return msgTypeString[m]
	if m.Is(SysExMsg) {
		return msgTypeString[SysExMsg]
	}

	if m.Is(MetaMsg) {
		return msgTypeString[m.Clear(MetaMsg)]
	}

	if m.Is(SysCommonMsg) {
		return msgTypeString[m.Clear(SysCommonMsg)]
	}

	if m.Is(RealTimeMsg) {
		return msgTypeString[m.Clear(RealTimeMsg)]
	}

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
}

// GetChannelMsgType returns the MsgType of a channel message. It should not be used by the end consumer.
func GetChannelMsgType(canary byte) (mType MsgType) {
	var sType MsgType
	//r.status = canary
	tp, ch := utils.ParseStatus(canary)
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

	switch tp {
	case 0xC:
		sType = ProgramChangeMsg
	case 0xD:
		sType = AfterTouchMsg
	case 0x8:
		sType = NoteOffMsg
	case 0x9:
		sType = NoteOnMsg
	case 0xA:
		sType = PolyAfterTouchMsg
	case 0xB:
		sType = ControlChangeMsg
	case 0xE:
		sType = PitchBendMsg
	default:
		return UnknownMsg
	}
	mType = mType.Set(sType)
	return mType
}

// GetRealtimeMsgType returns the MsgType of a realtime message. It should not be used by the end consumer.
func GetRealtimeMsgType(b byte) MsgType {
	ty, has := rtMessages[b]
	if !has {
		return UnknownMsg
	}
	return ty
}

// GetSysCommonMsgType returns the MsgType of a sys common message. It should not be used by the end consumer.
func GetSysCommonMsgType(b byte) MsgType {
	ty, has := syscommMessages[b]
	if !has {
		return UnknownMsg
	}
	return ty
}

// GetMetaMsgType returns the MsgType of a meta message. It should not be used by the end consumer.
func GetMetaMsgType(b byte) MsgType {
	return metaMessages[b]
}
