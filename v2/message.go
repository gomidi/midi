package midi

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

type Message struct {
	MsgType
	Data []byte
}

func NewMessage(data []byte) (m Message) {
	m.MsgType = GetMsgType(data)
	m.Data = data
	return
}

// Key returns the MIDI key - a number from 0 to 127 or
// -1, if it is no noteOn / noteOff / PolyAfterTouch message or an invalid key
func (m Message) Key() int8 {
	if m.MsgType.IsOneOf(NoteOnMsg, NoteOffMsg, PolyAfterTouchMsg) {
		k, _ := utils.ParseTwoUint7(m.Data[1], m.Data[2])
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

func (m Message) String() string {
	switch {
	case m.Is(ChannelMsg):
		var bf bytes.Buffer
		fmt.Fprintf(&bf, m.MsgType.String())
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
			return fmt.Sprintf("%s bpm: %v", m.MsgType.String(), m.BPM())
		case m.Is(MetaTimeSigMsg):
			num, denom := m.Meter()
			return fmt.Sprintf("%s meter: %v/%v", m.MsgType.String(), num, denom)
		case m.IsOneOf(MetaLyricMsg, MetaMarkerMsg, MetaCopyrightMsg, MetaTextMsg, MetaCuepointMsg, MetaDeviceMsg, MetaInstrumentMsg, MetaProgramNameMsg, MetaTrackNameMsg):
			return fmt.Sprintf("%s text: %q", m.MsgType.String(), m.Text())
		default:
			return m.MsgType.String()
		}
	case m.Is(SysExMsg):
		return m.MsgType.String()
	case m.Is(SysCommonMsg):
		return m.MsgType.String()
	case m.Is(RealTimeMsg):
		return m.MsgType.String()
	}

	return m.MsgType.String()
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
	if !m.Is(MetaTimeSigMsg) {
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

func (m Message) Pitch() (relative int16, absolute int16) {
	if !m.MsgType.Is(PitchBendMsg) {
		return -1, -1
	}

	rel, abs := utils.ParsePitchWheelVals(m.Data[1], m.Data[2])
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
func (m Message) Text() string {
	rd := bytes.NewReader(m.Data[2:])
	text, _ := utils.ReadText(rd)
	return text
}

func (m Message) Pressure() int8 {
	t := m.MsgType

	if t.Is(PolyAfterTouchMsg) {
		_, v := utils.ParseTwoUint7(m.Data[1], m.Data[2])
		return int8(v)
	}

	if t.Is(AfterTouchMsg) {
		return int8(utils.ParseUint7(m.Data[1]))
	}

	return -1
}

func (m Message) Program() int8 {
	t := m.MsgType

	if t.Is(ProgramChangeMsg) {
		return int8(utils.ParseUint7(m.Data[1]))
	}

	return -1
}

// Change returns the MIDI controllchange a number from 0 to 127 or
// -1, if it is no controller message
func (m Message) Change() int8 {
	if m.MsgType.Is(ControlChangeMsg) {
		_, v := utils.ParseTwoUint7(m.Data[1], m.Data[2])
		return int8(v)
	}

	return -1
}

// Channel returns the MIDI channel - a number from 0 to 15 or
// -1, if it is no channel message or an invalid channel number
func (m Message) Channel() int8 {
	if !m.MsgType.Is(ChannelMsg) {
		return -1
	}

	_, ch := utils.ParseStatus(m.Data[0])
	return int8(ch)
}

// Velocity returns the MIDI velocity - a number from 0 to 127 or
// -1, if it is no channel / noteOn / noteOff message or an invalid velocity
func (m Message) Velocity() int8 {
	if m.MsgType.IsOneOf(NoteOnMsg, NoteOffMsg) {
		_, v := utils.ParseTwoUint7(m.Data[1], m.Data[2])
		return int8(v)
	}

	return -1
}

// Controller returns the MIDI controller - a number from 0 to 127 or
// -1, if it is no controller message
func (m Message) Controller() int8 {
	if m.MsgType.Is(ControlChangeMsg) {
		c, _ := utils.ParseTwoUint7(m.Data[1], m.Data[2])
		return int8(c)
	}

	return -1
}

func GetMetaMsgType(b byte) MsgType {
	return metaMessages[b]
}
