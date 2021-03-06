package midi

import (
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

type MsgType uint64

const UnknownMsg MsgType = 0

const (
	ChannelMsg MsgType = 1 << iota
	MetaMsg
	RealTimeMsg
	SysCommonMsg
	SysExMsg
	NoteOnMsg
	NoteOffMsg
	ControlChangeMsg
	PitchBendMsg
	AfterTouchMsg
	PolyAfterTouchMsg
	ProgramChangeMsg
	MetaChannelMsg
	MetaCopyrightMsg
	MetaCuepointMsg
	MetaDeviceMsg
	MetaEndOfTrackMsg
	MetaInstrumentMsg
	MetaKeySigMsg
	MetaLyricMsg
	MetaTextMsg
	MetaMarkerMsg
	MetaPortMsg
	MetaSeqNumberMsg
	MetaSeqDataMsg
	MetaTempoMsg
	MetaTimeSigMsg
	MetaTrackNameMsg
	MetaSMPTEOffsetMsg
	MetaUndefinedMsg
	MetaProgramNameMsg
	TimingClockMsg
	TickMsg
	StartMsg
	ContinueMsg
	StopMsg
	ActiveSenseMsg
	ResetMsg
	SysExStartMsg
	SysExEndMsg
	SysExCompleteMsg
	SysExEscapeMsg
	SysExContinueMsg
	MTCMsg
	SongSelectMsg
	SPPMsg
	UndefinedMsg
	TuneMsg
	Channel0Msg
	Channel1Msg
	Channel2Msg
	Channel3Msg
	Channel4Msg
	Channel5Msg
	Channel6Msg
	Channel7Msg
	Channel8Msg
	Channel9Msg
	Channel10Msg
	Channel11Msg
	Channel12Msg
	Channel13Msg
	Channel14Msg
	Channel15Msg
)

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
	SysExStartMsg:      "SysExStartMsg",
	SysExEndMsg:        "SysExEndMsg",
	SysExCompleteMsg:   "SysExCompleteMsg",
	SysExEscapeMsg:     "SysExEscapeMsg",
	SysExContinueMsg:   "SysExContinueMsg",
	MTCMsg:             "MTCMsg",
	SongSelectMsg:      "SongSelectMsg",
	SPPMsg:             "SPPMsg",
	UndefinedMsg:       "UndefinedMsg",
	TuneMsg:            "TuneMsg",
	UnknownMsg:         "UnknownMsg",
	Channel0Msg:        "Channel0Msg",
	Channel1Msg:        "Channel1Msg",
	Channel2Msg:        "Channel2Msg",
	Channel3Msg:        "Channel3Msg",
	Channel4Msg:        "Channel4Msg",
	Channel5Msg:        "Channel5Msg",
	Channel6Msg:        "Channel6Msg",
	Channel7Msg:        "Channel7Msg",
	Channel8Msg:        "Channel8Msg",
	Channel9Msg:        "Channel9Msg",
	Channel10Msg:       "Channel10Msg",
	Channel11Msg:       "Channel11Msg",
	Channel12Msg:       "Channel12Msg",
	Channel13Msg:       "Channel13Msg",
	Channel14Msg:       "Channel14Msg",
	Channel15Msg:       "Channel15Msg",
}

func GetMsgType(msg []byte) (mType MsgType) {
	if len(msg) == 0 {
		return UnknownMsg
	}

	var canary = msg[0]

	// channel/Voice Category Status
	if canary >= 0x80 && canary <= 0xEF {
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
	} else {
		switch canary {
		case 0xF0, 0xF7:
			return SysExMsg
		// meta event
		case 0xFF:
			return GetMetaMessage(msg[1])
		default:
			return UnknownMsg
		}
	}
}

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

func (m MsgType) Set(flag MsgType) MsgType    { return m | flag }
func (m MsgType) Clear(flag MsgType) MsgType  { return m &^ flag }
func (m MsgType) Toggle(flag MsgType) MsgType { return m ^ flag }

func (m MsgType) Is(flag MsgType) bool    { return m&flag != 0 }
func (m MsgType) IsNot(flag MsgType) bool { return m&flag == 0 }
func (m MsgType) IsOneOf(flags ...MsgType) bool {
	for _, fl := range flags {
		if m&fl != 0 {
			return true
		}
	}
	return false
}
func (m MsgType) IsAllOf(flags ...MsgType) bool {
	for _, fl := range flags {
		if m&fl == 0 {
			return false
		}
	}
	return true
}
