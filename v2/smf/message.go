package smf

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// Message is a MIDI message that might appear in a SMF file, i.e. channel messages, sysex messages and meta messages.
type Message []byte

// Bytes return the underlying bytes of the message.
func (m Message) Bytes() []byte {
	return []byte(m)
}

// IsPlayable returns true, if the message can be send to an instrument.
func (m Message) IsPlayable() bool {
	if m.IsMeta() {
		return false
	}

	if m.Type() <= midi.UnknownMsg {
		return false
	}
	return true
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
		if len(msg) == 1 {
			return midi.UnknownMsg
		}
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

// GetSysEx returns true, if the message is a sysex message.
// Then it extracts the inner bytes to the given slice.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetSysEx(bt *[]byte) bool {
	return midi.Message(m).GetSysEx(bt)
}

// GetNoteOn returns true if (and only if) the message is a NoteOnMsg.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
func (m Message) GetNoteOn(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).GetNoteOn(channel, key, velocity)
}

// GetNoteStart returns true if (and only if) the message is a NoteOnMsg with a velocity > 0.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
func (m Message) GetNoteStart(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).GetNoteStart(channel, key, velocity)
}

// GetNoteOff returns true if (and only if) the message is a NoteOffMsg.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
func (m Message) GetNoteOff(channel, key, velocity *uint8) (is bool) {
	return midi.Message(m).GetNoteOff(channel, key, velocity)
}

// GetChannel returns true if (and only if) the message is a ChannelMsg.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
func (m Message) GetChannel(channel *uint8) (is bool) {
	return midi.Message(m).GetChannel(channel)
}

// GetNoteEnd returns true if (and only if) the message is a NoteOnMsg with a velocity == 0 or a NoteOffMsg.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
func (m Message) GetNoteEnd(channel, key *uint8) (is bool) {
	return midi.Message(m).GetNoteEnd(channel, key)
}

// GetPolyAfterTouch returns true if (and only if) the message is a PolyAfterTouchMsg.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
func (m Message) GetPolyAfterTouch(channel, key, pressure *uint8) (is bool) {
	return midi.Message(m).GetPolyAfterTouch(channel, key, pressure)
}

// GetAfterTouch returns true if (and only if) the message is a AfterTouchMsg.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
func (m Message) GetAfterTouch(channel, pressure *uint8) (is bool) {
	return midi.Message(m).GetAfterTouch(channel, pressure)
}

// GetProgramChange returns true if (and only if) the message is a ProgramChangeMsg.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
func (m Message) GetProgramChange(channel, program *uint8) (is bool) {
	return midi.Message(m).GetProgramChange(channel, program)
}

// GetPitchBend returns true if (and only if) the message is a PitchBendMsg.
// Then it also extracts the data to the given arguments
// Either relative or absolute may be nil, if not needed.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetPitchBend(channel *uint8, relative *int16, absolute *uint16) (is bool) {
	return midi.Message(m).GetPitchBend(channel, relative, absolute)
}

// GetControlChange returns true if (and only if) the message is a ControlChangeMsg.
// Then it also extracts the data to the given arguments
// Only arguments that are not nil are parsed and filled.
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
		//var bl1 bool
		//var bl2 bool
		var text string
		var bpm float64
		var bt []byte
		var k Key

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
		case m.GetMetaKey(&k):
			fmt.Fprintf(&bf, " key: %s", k.String())
		//case m.GetMetaKeySig(&val1, &val2, &bl1, &bl2):
		//	fmt.Fprintf(&bf, " key: %v num: %v ismajor: %v isflat: %v", val1, val2, bl1, bl2)
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

func _MetaMessage(typ byte, data []byte) Message {
	b := []byte{byte(0xFF), typ}
	b = append(b, utils.VlqEncode(uint32(len(data)))...)
	if len(data) != 0 {
		b = append(b, data...)
	}
	return b
}

// GetMetaMeter is a handier wrapper around GetMetaTimeSig.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaMeter(num, denom *uint8) (is bool) {
	return m.GetMetaTimeSig(num, denom, nil, nil)
}

// metaData strips away the meta byte and the metatype byte and the varlength byte
func (m Message) metaDataWithoutVarlength() []byte {
	//fmt.Printf("original data: % X\n", m.Data)
	return m[3:]
}

// GetMetaChannel return true, if (and only if) the message is a MetaChannelMsg.
// Then it also extracts the channel to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaChannel(channel *uint8) bool {
	if !m.Is(MetaChannelMsg) {
		return false
	}

	if len(m) != 4 {
		return false
	}

	if channel != nil {
		data := m.metaDataWithoutVarlength()
		*channel = data[0]
	}

	return true
}

// GetMetaPort return true, if (and only if) the message is a MetaPortMsg.
// Then it also extracts the port to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaPort(port *uint8) bool {
	if !m.Is(MetaPortMsg) {
		return false
	}

	if len(m) != 4 {
		return false
	}

	if port != nil {
		data := m.metaDataWithoutVarlength()

		*port = data[0]
	}

	return true
}

// GetMetaSeqNumber return true, if (and only if) the message is a MetaSeqNumberMsg.
// Then it also extracts the sequenceNumber to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaSeqNumber(sequenceNumber *uint16) bool {
	if !m.Is(MetaSeqNumberMsg) {
		return false
	}

	if len(m) != 2 && len(m) < 5 {
		return false
	}

	if sequenceNumber != nil {
		// Zero length sequences allowed according to http://home.roadrunner.com/~jgglatt/tech/midifile/seq.htm
		if len(m) == 2 {
			*sequenceNumber = 0
			return true
		}
		//fmt.Printf("% X\n", []byte{m[3], m[4]})
		*sequenceNumber = utils.ParseUint16(m[3], m[4])
	}

	return true

}

// GetMetaSeqData return true, if (and only if) the message is a MetaSeqDataMsg.
// Then it also extracts the data to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaSeqData(bt *[]byte) bool {
	if !m.Is(MetaSeqDataMsg) {
		return false
	}

	if len(m) < 4 {
		return false
	}

	if bt != nil {
		data := m.metaDataWithoutVarlength()
		*bt = data
	}
	return true
}

// GetMetaKey is a handier wrapper around GetMetaKeySig. It returns nil if the message is no MetaKeySigMsg.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaKey(key *Key) bool {
	var k Key
	if m.GetMetaKeySig(&k.Key, &k.Num, &k.IsMajor, &k.IsFlat) {
		if key != nil {
			*key = k
		}
		return true
	}
	return false
}

// GetMetaKeySig return true, if (and only if) the message is a MetaKeySigMsg.
// Then it also extracts the data to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaKeySig(key, num *uint8, isMajor *bool, isFlat *bool) bool {
	if !m.Is(MetaKeySigMsg) {
		return false
	}

	if len(m) != 5 {
		return false
	}

	data := m.metaDataWithoutVarlength()

	if len(data) != 2 {
		//err = unexpectedMessageLengthError("KeySignature expected length 2")
		//return nil, err
		return false
	}

	sharpsOrFlats := int8(data[0])

	// Mode is Major or Minor.
	mode := data[1]

	_num := sharpsOrFlats
	if _num < 0 {
		_num = _num * (-1)
	}

	if key != nil {
		*key = utils.KeyFromSharpsOrFlats(sharpsOrFlats, mode)
	}

	if num != nil {
		*num = uint8(_num)
	}

	if isMajor != nil {
		*isMajor = mode == majorMode
	}

	if isFlat != nil {
		*isFlat = sharpsOrFlats < 0
	}

	return true
}

// GetMetaSMPTEOffsetMsg return true, if (and only if) the message is a MetaSMPTEOffsetMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaSMPTEOffsetMsg(hour, minute, second, frame, fractframe *uint8) bool {
	if !m.Is(MetaSMPTEOffsetMsg) {
		return false
	}

	if len(m) != 8 {
		//err = unexpectedMessageLengthError("KeySignature expected length 2")
		//return nil, err
		return false
	}

	data := m.metaDataWithoutVarlength()

	if len(data) != 5 {
		//err = unexpectedMessageLengthError("SMPTEOffset expected length 5")
		//return nil, err
		return false
	}

	if hour != nil {
		*hour = data[0]
	}

	if minute != nil {
		*minute = data[1]
	}

	if second != nil {
		*second = data[2]
	}

	if frame != nil {
		*frame = data[3]
	}

	if fractframe != nil {
		*fractframe = data[4]
	}

	return true
}

// GetMetaTimeSig return true, if (and only if) the message is a MetaTimeSigMsg.
// Then it also extracts the data to the given arguments.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaTimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter *uint8) (is bool) {
	if !m.Is(MetaTimeSigMsg) {
		//fmt.Println("not timesig message")
		return false
	}

	if len(m) != 7 {
		return false
	}

	data := m.metaDataWithoutVarlength()

	if len(data) != 4 {
		//fmt.Printf("not correct data lenght: % X \n", data)
		//err = unexpectedMessageLengthError("TimeSignature expected length 4")
		return false
	}

	//fmt.Printf("TimeSigData: % X\n", data)

	if numerator != nil {
		*numerator = data[0]
	}

	if clocksPerClick != nil {
		*clocksPerClick = data[2]
	}

	if demiSemiQuaverPerQuarter != nil {
		*demiSemiQuaverPerQuarter = data[3]
	}

	if denominator != nil {
		*denominator = bin2decDenom(data[1])
	}

	return true
}

// GetMetaTempo return true, if (and only if) the message is a MetaTempoMsg.
// Then it also extracts the BPM to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaTempo(bpm *float64) (is bool) {
	if !m.Is(MetaTempoMsg) {
		return false
	}

	if len(m) < 4 {
		return false
	}

	if bpm != nil {
		//fmt.Printf("tempo pure bytes: % X\n", m.metaDataWithoutVarlength())
		rd := bytes.NewReader(m.metaDataWithoutVarlength())
		microsecondsPerCrotchet, err := utils.ReadUint24(rd)
		if err != nil {
			return false
		}

		*bpm = float64(60000000) / float64(microsecondsPerCrotchet)
	}

	return true
}

// GetMetaLyric return true, if (and only if) the message is a MetaLyricMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaLyric(text *string) (is bool) {
	if !m.Is(MetaLyricMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}

	return true
}

// GetMetaCopyright return true, if (and only if) the message is a MetaCopyrightMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaCopyright(text *string) (is bool) {
	if !m.Is(MetaCopyrightMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}
	return true
}

// GetMetaCuepoint return true, if (and only if) the message is a MetaCuepointMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaCuepoint(text *string) (is bool) {
	if !m.Is(MetaCuepointMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}
	return true
}

// GetMetaDevice return true, if (and only if) the message is a MetaDeviceMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaDevice(text *string) (is bool) {
	if !m.Is(MetaDeviceMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}
	return true
}

// GetMetaInstrument return true, if (and only if) the message is a MetaInstrumentMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaInstrument(text *string) (is bool) {
	if !m.Is(MetaInstrumentMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}
	return true
}

// GetMetaMarker return true, if (and only if) the message is a MetaMarkerMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaMarker(text *string) (is bool) {
	if !m.Is(MetaMarkerMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}
	return true
}

// GetMetaProgramName return true, if (and only if) the message is a MetaProgramNameMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaProgramName(text *string) (is bool) {
	if !m.Is(MetaProgramNameMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}
	return true
}

// GetMetaText return true, if (and only if) the message is a MetaTextMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaText(text *string) (is bool) {
	if !m.Is(MetaTextMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}
	return true
}

// GetMetaTrackName return true, if (and only if) the message is a MetaTrackNameMsg.
// Then it also extracts the text to the given argument.
// Only arguments that are not nil are parsed and filled.
func (m Message) GetMetaTrackName(text *string) (is bool) {
	if !m.Is(MetaTrackNameMsg) {
		return false
	}

	if len(m) < 3 {
		return false
	}

	if text != nil {
		m.text(text)
	}
	return true
}

// Only arguments that are not nil are parsed and filled.
func (m Message) text(text *string) {
	if text != nil {
		*text, _ = utils.ReadText(bytes.NewReader(m[2:]))
	}
	return
}
