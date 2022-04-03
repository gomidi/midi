package smf

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"math/big"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// MetaMessage represents a MIDI meta message. It can be created from the MIDI bytes of a message, by calling NewMetaMessage.
//type MetaMessage midi.Message

/*
func newMetaMessage(typ byte, data []byte) []byte {
	return _MetaMessage(typ, data)
}
*/

func (msg Message) IsMeta() bool {
	if len(msg) == 0 {
		return false
	}
	return msg[0] == 0xFF
}

// TODO or return a MetaMessage here??
func _MetaMessage(typ byte, data []byte) Message {
	//func _MetaMessage(typ byte, data []byte) MetaMessage {
	//fmt.Printf("NewMetaMessage %X % X\n", typ, data)
	b := []byte{byte(0xFF), typ}
	b = append(b, utils.VlqEncode(uint32(len(data)))...)
	if len(data) != 0 {
		b = append(b, data...)
	}

	/*
		var m midi.Message
		m.Type = GetMetaType(typ)
		m.Data = b
		return m
	*/
	//return midi.Message(b)
	return b
}

// Meter returns the meter of a MetaTimeSigMsg.
// For other messages, it returns 0,0.
func (m Message) ScanMetaMeter(num, denom *uint8) (is bool) {
	return m.ScanMetaTimeSig(num, denom, nil, nil)
}

// metaData strips away the meta byte and the metatype byte and the varlength byte
func (m Message) metaDataWithoutVarlength() []byte {
	//fmt.Printf("original data: % X\n", m.Data)
	return m[3:]
}

func (m Message) ScanMetaChannel(channel *uint8) bool {
	if !m.Type().Is(MetaChannelMsg) {
		return false
	}

	data := m.metaDataWithoutVarlength()

	*channel = data[0]
	return true
}

func (m Message) ScanMetaPort(port *uint8) bool {
	if !m.Type().Is(MetaPortMsg) {
		return false
	}

	data := m.metaDataWithoutVarlength()

	*port = data[0]
	return true
}

func (m Message) ScanMetaSeqNumber(sequenceNumber *uint16) bool {
	if !m.Type().Is(MetaSeqNumberMsg) {
		return false
	}

	// Zero length sequences allowed according to http://home.roadrunner.com/~jgglatt/tech/midifile/seq.htm
	if len(m) == 2 {
		*sequenceNumber = 0
		return true
	}

	*sequenceNumber = utils.ParseUint16(m[2], m[3])

	return true

}

func (m Message) ScanMetaSeqData(bt *[]byte) bool {
	if !m.Is(MetaSeqDataMsg) {
		return false
	}

	data := m.metaDataWithoutVarlength()
	*bt = data
	return true
}

func (m Message) ScanMetaKeySig(key, num *uint8, isMajor *bool, isFlat *bool) bool {
	if !m.Is(MetaKeySigMsg) {
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

	*key = utils.KeyFromSharpsOrFlats(sharpsOrFlats, mode)
	*num = uint8(_num)
	*isMajor = mode == majorMode
	*isFlat = sharpsOrFlats < 0

	return true
}

func (m Message) ScanMetaSMPTEOffsetMsg(hour, minute, second, frame, fractframe *uint8) bool {
	if !m.Is(MetaSMPTEOffsetMsg) {
		return false
	}

	data := m.metaDataWithoutVarlength()

	if len(data) != 5 {
		//err = unexpectedMessageLengthError("SMPTEOffset expected length 5")
		//return nil, err
		return false
	}

	*hour = data[0]
	*minute = data[1]
	*second = data[2]
	*frame = data[3]
	*fractframe = data[4]

	return true
}

// TimeSig returns the numerator, denominator, clocksPerClick and demiSemiQuaverPerQuarter of a
// MetaTimeSigMsg. For other messages, it returns 0,0,0,0.
func (m Message) ScanMetaTimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter *uint8) (is bool) {
	if !m.Type().Is(MetaTimeSigMsg) {
		//fmt.Println("not timesig message")
		return false
	}

	data := m.metaDataWithoutVarlength()

	if len(data) != 4 {
		//fmt.Printf("not correct data lenght: % X \n", data)
		//err = unexpectedMessageLengthError("TimeSignature expected length 4")
		return false
	}

	//fmt.Printf("TimeSigData: % X\n", data)

	*numerator = data[0]
	*denominator = data[1]
	if clocksPerClick != nil {
		*clocksPerClick = data[2]
	}
	if demiSemiQuaverPerQuarter != nil {
		*demiSemiQuaverPerQuarter = data[3]
	}
	*denominator = bin2decDenom(*denominator)
	return true
}

// bin2decDenom converts the binary denominator to the decimal
func bin2decDenom(bin uint8) uint8 {
	if bin == 0 {
		return 1
	}
	return 2 << (bin - 1)
}

// Tempo returns true if (and only if) the message is a MetaTempoMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaTempo(bpm *float64) (is bool) {
	if !m.Type().Is(MetaTempoMsg) {
		return false
	}

	//fmt.Printf("tempo pure bytes: % X\n", m.metaDataWithoutVarlength())
	rd := bytes.NewReader(m.metaDataWithoutVarlength())
	microsecondsPerCrotchet, err := utils.ReadUint24(rd)
	if err != nil {
		return false
	}

	*bpm = float64(60000000) / float64(microsecondsPerCrotchet)

	return true
}

// Lyric returns true if (and only if) the message is a MetaLyricMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaLyric(text *string) (is bool) {
	if !m.Type().Is(MetaLyricMsg) {
		return false
	}
	m.text(text)
	return true
}

/*
func (m Message) Is(t midi.Type) bool {
	return m.Type().Is(t)
}
*/

// Copyright returns true if (and only if) the message is a MetaCopyrightMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaCopyright(text *string) (is bool) {
	if !m.Is(MetaCopyrightMsg) {
		return false
	}
	m.text(text)
	return true
}

// Cuepoint returns true if (and only if) the message is a MetaCuepointMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaCuepoint(text *string) (is bool) {
	if !m.Is(MetaCuepointMsg) {
		return false
	}
	m.text(text)
	return true
}

// Device returns true if (and only if) the message is a MetaDeviceMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaDevice(text *string) (is bool) {
	if !m.Is(MetaDeviceMsg) {
		return false
	}
	m.text(text)
	return true
}

// Instrument returns true if (and only if) the message is a MetaInstrumentMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaInstrument(text *string) (is bool) {
	if !m.Is(MetaInstrumentMsg) {
		return false
	}
	m.text(text)
	return true
}

// Marker returns true if (and only if) the message is a MetaMarkerMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaMarker(text *string) (is bool) {
	if !m.Is(MetaMarkerMsg) {
		return false
	}
	m.text(text)
	return true
}

// ProgramName returns true if (and only if) the message is a MetaProgramNameMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaProgramName(text *string) (is bool) {
	if !m.Is(MetaProgramNameMsg) {
		return false
	}
	m.text(text)
	return true
}

// Text returns true if (and only if) the message is a MetaTextMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaText(text *string) (is bool) {
	switch m.Type() {
	case MetaLyricMsg, MetaMarkerMsg, MetaCopyrightMsg, MetaTextMsg, MetaCuepointMsg, MetaDeviceMsg, MetaInstrumentMsg, MetaProgramNameMsg, MetaTrackNameMsg:
		m.text(text)
		return true
	default:
		return false
	}
}

// TrackName returns true if (and only if) the message is a MetaTrackNameMsg.
// Then it also extracts the data to the given arguments
func (m Message) ScanMetaTrackName(text *string) (is bool) {
	if !m.Is(MetaTrackNameMsg) {
		return false
	}
	m.text(text)
	return true
}

func (m Message) text(text *string) {
	*text, _ = utils.ReadText(bytes.NewReader(m[2:]))
	return
}

/*
func init() {
	var _ midi.Message = NewMetaCuepoint("test")
}
*/

/*
// Is returns if the given message type is part of the combination of message types
func Is(mt1, mt2 midi.MessageType) bool {
	return mt1.Category() == mt2.Category() && (mt1.Val()&mt2.Val()) != 0
}
*/

func ReadMetaData(tp midi.Type, rd io.Reader) (data []byte, err error) {
	return utils.ReadVarLengthData(rd)
}

const (
	// End of track
	// the handler is supposed to keep track of the current track

	byteEndOfTrack        = byte(0x2F)
	byteSequenceNumber    = byte(0x00)
	byteText              = byte(0x01)
	byteCopyright         = byte(0x02)
	byteTrackSequenceName = byte(0x03)
	byteInstrument        = byte(0x04)
	byteLyric             = byte(0x05)
	byteMarker            = byte(0x06)
	byteCuepoint          = byte(0x07)
	byteMIDIChannel       = byte(0x20)
	byteDevicePort        = byte(0x9)
	byteMIDIPort          = byte(0x21)
	byteTempo             = byte(0x51)
	byteTimeSignature     = byte(0x58)
	byteKeySignature      = byte(0x59)
	byteSequencerSpecific = byte(0x7F)
	byteSMPTEOffset       = byte(0x54)
	byteProgramName       = byte(0x8)
)

var metaMessages = map[byte]midi.Type{
	byteEndOfTrack:        MetaEndOfTrackMsg,
	byteSequenceNumber:    MetaSeqNumberMsg,
	byteText:              MetaTextMsg,
	byteCopyright:         MetaCopyrightMsg,
	byteTrackSequenceName: MetaTrackNameMsg,
	byteInstrument:        MetaInstrumentMsg,
	byteLyric:             MetaLyricMsg,
	byteMarker:            MetaMarkerMsg,
	byteCuepoint:          MetaCuepointMsg,
	byteMIDIChannel:       MetaChannelMsg,
	byteDevicePort:        MetaDeviceMsg,
	byteMIDIPort:          MetaPortMsg,
	byteTempo:             MetaTempoMsg,
	byteTimeSignature:     MetaTimeSigMsg,
	byteKeySignature:      MetaKeySigMsg,
	byteSMPTEOffset:       MetaSMPTEOffsetMsg,
	byteSequencerSpecific: MetaSeqDataMsg,
	byteProgramName:       MetaProgramNameMsg,
}

// GetMetaType returns the MetaType of a meta message. It should not be used by the end consumer.
func GetMetaType(b byte) midi.Type {
	return metaMessages[b]
}

const bpmFac = 60000000

// NewMetaLyric returns the bytes of a lyric meta message
func MetaLyric(text string) []byte {
	return _MetaMessage(byteLyric, []byte(text))
}

// NewMetaCopyright returns the bytes of a copyright meta message
func MetaCopyright(text string) []byte {
	return _MetaMessage(byteCopyright, []byte(text))
}

// NewMetaChannel returns the bytes of a channel meta message
func MetaChannel(ch uint8) Message {
	return _MetaMessage(byteMIDIChannel, []byte{byte(ch)})
}

// NewMetaCuepoint returns the bytes of a cuepoint meta message
func MetaCuepoint(text string) Message {
	return _MetaMessage(byteCuepoint, []byte(text))
}

// NewMetaDevice returns the bytes of a device meta message
func MetaDevice(text string) Message {
	return _MetaMessage(byteDevicePort, []byte(text))
}

// EOT are the bytes of an End Of Track meta message. Don't use it directly.
var EOT = _MetaMessage(byteEndOfTrack, nil)

// NewMetaInstrument returns the bytes of a instrument meta message
func MetaInstrument(text string) Message {
	return _MetaMessage(byteInstrument, []byte(text))
}

// NewMetaMarker returns the bytes of a marker meta message
func MetaMarker(text string) Message {
	return _MetaMessage(byteMarker, []byte(text))
}

// NewMetaPort returns the bytes of a port meta message
func MetaPort(p uint8) Message {
	return _MetaMessage(byteMIDIPort, []byte{byte(p)})
}

// NewMetaProgram returns the bytes of a program meta message
func MetaProgram(text string) Message {
	return _MetaMessage(byteProgramName, []byte(text))
}

// NewMetaSequenceNo returns the bytes of a sequence number meta message
func MetaSequenceNo(no uint16) Message {
	var bf bytes.Buffer
	binary.Write(&bf, binary.BigEndian, no)
	return _MetaMessage(byteSequenceNumber, bf.Bytes())
}

// NewMetaSequencerData returns the bytes of a sequencer data meta message
func MetaSequencerData(data []byte) Message {
	return _MetaMessage(byteSequencerSpecific, data)
}

// NewMetaSMPTE returns the bytes of a SMPTE meta message
func MetaSMPTE(hour, minute, second, frame, fractionalFrame byte) Message {
	return _MetaMessage(byteSMPTEOffset, []byte{hour, minute, second, frame, fractionalFrame})
}

// NewMetaTempo returns the bytes of a tempo meta message for the given beats per minute.
func MetaTempo(bpm float64) Message {
	r := uint32(math.Round(bpmFac / bpm))
	if r > 0x0FFFFFFF {
		r = 0x0FFFFFFF
	}

	b4 := big.NewInt(int64(r)).Bytes()

	var b = []byte{0, 0, 0}
	switch len(b4) {
	case 0:
	case 1:
		b[2] = b4[0]
	case 2:
		b[2] = b4[1]
		b[1] = b4[0]
	case 3:
		b[2] = b4[2]
		b[1] = b4[1]
		b[0] = b4[0]
	}

	return _MetaMessage(byteTempo, b)
}

// NewMetaText returns the bytes of a text meta message.
func MetaText(text string) Message {
	return _MetaMessage(byteText, []byte(text))
}

// NewMetaTrackSequenceName returns the bytes of a track sequence name meta message.
func MetaTrackSequenceName(text string) Message {
	return _MetaMessage(byteTrackSequenceName, []byte(text))
}

// NewMetaUndefined returns the bytes of an undefined meta message.
func MetaUndefined(typ byte, data []byte) Message {
	return _MetaMessage(typ, data)
}

const (
	degreeC  = 0
	degreeCs = 1
	degreeDf = degreeCs
	degreeD  = 2
	degreeDs = 3
	degreeEf = degreeDs
	degreeE  = 4
	degreeF  = 5
	degreeFs = 6
	degreeGf = degreeFs
	degreeG  = 7
	degreeGs = 8
	degreeAf = degreeGs
	degreeA  = 9
	degreeAs = 10
	degreeBf = degreeAs
	degreeB  = 11
	degreeCf = degreeB
)

// Supplied to KeySignature
const (
	majorMode = 0
	minorMode = 1
)

// NewMetaKey returns the bytes of a key meta message.
func MetaKey(key uint8, isMajor bool, num uint8, isFlat bool) Message {
	mi := int8(0)
	if !isMajor {
		mi = 1
	}
	sf := int8(num)

	if isFlat {
		sf = sf * (-1)
	}

	return _MetaMessage(byteKeySignature, []byte{byte(sf), byte(mi)})
}

// NewMetaMeter returns the bytes of a time signature meta message.
func MetaMeter(num, denom uint8) Message {
	if denom == 0 {
		denom = 1
	}

	return MetaTimeSig(num, denom, 8, 8)
}

// NewMetaTimeSig returns the bytes of a time signature meta message.
func MetaTimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) Message {
	cpcl := clocksPerClick
	if cpcl == 0 {
		cpcl = byte(8)
	}

	dsqpq := demiSemiQuaverPerQuarter
	if dsqpq == 0 {
		dsqpq = byte(8)
	}

	var denom = dec2binDenom(denominator)

	return _MetaMessage(byteTimeSignature, []byte{numerator, denom, cpcl, dsqpq})

}

// dec2binDenom converts the decimal denominator to the binary one
// it works, use it!
func dec2binDenom(dec uint8) (bin uint8) {
	if dec <= 1 {
		return 0
	}
	for dec > 2 {
		bin++
		dec = dec >> 1

	}
	return bin + 1
}
