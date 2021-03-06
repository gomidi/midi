package midi

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

type Message struct {
	Type MessageType
	Data []byte
}

func NewMessage(data []byte) (m Message) {
	m.Type = GetMessageType(data)
	m.Data = data
	return
}

type MessageType uint64

const UnknownMsg MessageType = 0

const (
	ChannelMsg MessageType = 1 << iota
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

var channelType = map[uint8]MessageType{
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

var msgTypeString = map[MessageType]string{
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

func (m MessageType) String() string {
	//return msgTypeString[m]
	if m.Is(SysExMsg) {
		return msgTypeString[SysExMsg]
	}

	if m.Is(MetaMsg) {
		return msgTypeString[Clear(m, MetaMsg)]
	}

	if m.Is(SysCommonMsg) {
		return msgTypeString[Clear(m, SysCommonMsg)]
	}

	if m.Is(RealTimeMsg) {
		return msgTypeString[Clear(m, RealTimeMsg)]
	}

	if m.Is(ChannelMsg) {
		var clCh MessageType

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

		return msgTypeString[clCh] + " & " + msgTypeString[Clear(Clear(m, ChannelMsg), clCh)]
	}

	return "Unknown"
}

// Key returns the MIDI key - a number from 0 to 127 or
// -1, if it is no noteOn / noteOff / PolyAfterTouch message or an invalid key
func (msg Message) Key() int8 {
	if msg.Type.IsOneOf(NoteOnMsg, NoteOffMsg, PolyAfterTouchMsg) {
		k, _ := utils.ParseTwoUint7(msg.Data[1], msg.Data[2])
		return int8(k)
	}

	return -1
}

// IsNoteStart checks, if we have a de facto note start, i.e. a NoteOnMsg  with velocity > 0
func (m Message) IsNoteStart() bool {
	if m.Is(NoteOnMsg) && m.Velocity() > 0 {
		return true
	}
	return false
}

// IsNoteEnd checks, if we have a de facto note end, i.e.  a NoteoffMsg or a NoteOnMsg with velocity == 0
func (m Message) IsNoteEnd() bool {
	if m.Is(NoteOffMsg) {
		return true
	}
	if m.Is(NoteOnMsg) && m.Velocity() == 0 {
		return true
	}
	return false
}

func (m Message) Is(t MessageType) bool {
	return m.Type.Is(t)
}

func (m Message) String() string {
	switch {
	case m.Is(ChannelMsg):
		var bf bytes.Buffer
		fmt.Fprintf(&bf, m.Type.String())
		//fmt.Fprintf(&bf, " channel: %v ", m.Channel())
		switch {
		case m.Is(NoteOnMsg):
			fmt.Fprintf(&bf, " key: %v velocity: %v", m.Key(), m.Velocity())
		case m.Is(NoteOffMsg):
			fmt.Fprintf(&bf, " key: %v velocity: %v", m.Key(), m.Velocity())
		case m.Is(PolyAfterTouchMsg):
			fmt.Fprintf(&bf, " key: %v pressure: %v", m.Key(), m.Pressure())
		case m.Is(AfterTouchMsg):
			fmt.Fprintf(&bf, " pressure: %v", m.Pressure())
		case m.Is(ProgramChangeMsg):
			fmt.Fprintf(&bf, " program: %v", m.Program())
		case m.Is(PitchBendMsg):
			rel, abs := m.Pitch()
			fmt.Fprintf(&bf, " pitch: %v / %v", rel, abs)
		case m.Is(ControlChangeMsg):
			fmt.Fprintf(&bf, " controller: %v change: %v", m.Controller(), m.Change())
		default:
		}
		return bf.String()
	case m.Is(MetaMsg):
		switch {
		case m.Is(MetaTempoMsg):
			return fmt.Sprintf("%s bpm: %v", m.Type.String(), m.BPM())
		case m.Is(MetaTimeSigMsg):
			num, denom := m.Meter()
			return fmt.Sprintf("%s meter: %v/%v", m.Type.String(), num, denom)
		case m.IsOneOf(MetaLyricMsg, MetaMarkerMsg, MetaCopyrightMsg, MetaTextMsg, MetaCuepointMsg, MetaDeviceMsg, MetaInstrumentMsg, MetaProgramNameMsg, MetaTrackNameMsg):
			return fmt.Sprintf("%s text: %q", m.Type.String(), m.Text())
		default:
			return m.Type.String()
		}
	case m.Is(SysExMsg):
		return m.Type.String()
	case m.Is(SysCommonMsg):
		return m.Type.String()
	case m.Is(RealTimeMsg):
		return m.Type.String()
	}

	return m.Type.String()
}

func (m Message) Meter() (num, denom uint8) {
	num, denom, _, _ = m.TimeSig()
	return
}

// metaData strips away the meta byte and the metatype byte and the varlength byte
func (m Message) metaDataWithoutVarlength() []byte {
	//fmt.Printf("original data: % X\n", m.Data)
	return m.Data[3:]
}

func (m Message) TimeSig() (numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) {
	if m.IsNot(MetaTimeSigMsg) {
		//fmt.Println("not timesig message")
		return 0, 0, 0, 0
	}

	data := m.metaDataWithoutVarlength()

	if len(data) != 4 {
		//fmt.Printf("not correct data lenght: % X \n", data)
		//err = unexpectedMessageLengthError("TimeSignature expected length 4")
		return 0, 0, 0, 0
	}

	//fmt.Printf("TimeSigData: % X\n", data)

	numerator = data[0]
	denominator = data[1]
	clocksPerClick = data[2]
	demiSemiQuaverPerQuarter = data[3]
	denominator = bin2decDenom(denominator)
	return
}

func (m Message) IsNot(t MessageType) bool {
	return m.Type.IsNot(t)
}
func (m Message) IsOneOf(ts ...MessageType) bool {
	return m.Type.IsOneOf(ts...)
}
func (m Message) IsAllOf(ts ...MessageType) bool {
	return m.Type.IsAllOf(ts...)
}

func Set(b, flag MessageType) MessageType         { return b | flag }
func Clear(b, flag MessageType) MessageType       { return b &^ flag }
func Toggle(b, flag MessageType) MessageType      { return b ^ flag }
func (b MessageType) Is(flag MessageType) bool    { return b&flag != 0 }
func (b MessageType) IsNot(flag MessageType) bool { return b&flag == 0 }
func (b MessageType) IsOneOf(flags ...MessageType) bool {
	for _, fl := range flags {
		if b&fl != 0 {
			return true
		}
	}
	return false
}
func (b MessageType) IsAllOf(flags ...MessageType) bool {
	for _, fl := range flags {
		if b&fl == 0 {
			return false
		}
	}
	return true
}

func (msg Message) Pitch() (relative int16, absolute int16) {
	if msg.Type.IsNot(PitchBendMsg) {
		return -1, -1
	}

	rel, abs := utils.ParsePitchWheelVals(msg.Data[1], msg.Data[2])
	return rel, int16(abs)
}

/*
Text returns the text for the meta messages

Lyric
Copyright
Cuepoint
Device
Instrument
Marker
Program
Text
TrackSequenceName
*/
func (msg Message) Text() string {
	rd := bytes.NewReader(msg.Data[2:])
	text, _ := utils.ReadText(rd)
	return text
}

/*
missing meta messages:
func Channel(ch uint8) Message {
func Port(p uint8) Message {
func SequenceNo(no uint16) Message {
func SequencerData(data []byte) Message {
func SMPTE(hour, minute, second, frame, fractionalFrame byte) Message {
func Tempo(bpm float64) Message {
func Undefined(typ byte, data []byte) Message {
*/

/*
also TODO
SysExMsg
Type parsing of
RealTimeMsg
SysCommonMsg
*/

func (msg Message) Pressure() int8 {
	t := msg.Type

	if t.Is(PolyAfterTouchMsg) {
		_, v := utils.ParseTwoUint7(msg.Data[1], msg.Data[2])
		return int8(v)
	}

	if t.Is(AfterTouchMsg) {
		return int8(utils.ParseUint7(msg.Data[1]))
	}

	return -1
}

func (msg Message) Program() int8 {
	t := msg.Type

	if t.Is(ProgramChangeMsg) {
		return int8(utils.ParseUint7(msg.Data[1]))
	}

	return -1
}

// Change returns the MIDI controllchange a number from 0 to 127 or
// -1, if it is no controller message
func (msg Message) Change() int8 {
	if msg.Type.Is(ControlChangeMsg) {
		_, v := utils.ParseTwoUint7(msg.Data[1], msg.Data[2])
		return int8(v)
	}

	return -1
}

// Channel returns the MIDI channel - a number from 0 to 15 or
// -1, if it is no channel message or an invalid channel number
func (msg Message) Channel() int8 {
	if msg.Type.IsNot(ChannelMsg) {
		return -1
	}

	_, ch := utils.ParseStatus(msg.Data[0])
	return int8(ch)
}

// Velocity returns the MIDI velocity - a number from 0 to 127 or
// -1, if it is no channel / noteOn / noteOff message or an invalid velocity
func (msg Message) Velocity() int8 {
	if msg.Type.IsOneOf(NoteOnMsg, NoteOffMsg) {
		_, v := utils.ParseTwoUint7(msg.Data[1], msg.Data[2])
		return int8(v)
	}

	return -1
}

// Controller returns the MIDI controller - a number from 0 to 127 or
// -1, if it is no controller message
func (msg Message) Controller() int8 {
	if msg.Type.Is(ControlChangeMsg) {
		c, _ := utils.ParseTwoUint7(msg.Data[1], msg.Data[2])
		return int8(c)
	}

	return -1
}

func GetMessageType(msg []byte) (mType MessageType) {
	if len(msg) == 0 {
		return UnknownMsg
	}

	var canary = msg[0]

	// channel/Voice Category Status
	if canary >= 0x80 && canary <= 0xEF {
		var sType MessageType
		//r.status = canary
		tp, ch := utils.ParseStatus(canary)
		mType = Set(mType, ChannelMsg)
		var ctype MessageType

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

		mType = Set(mType, ctype)

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
		mType = Set(mType, sType)
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

func GetMetaMessage(b byte) MessageType {
	return metaMessages[b]
}
