package smf

import (
	//	"bytes"
	//"encoding/binary"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"

	//"math"
	//"math/big"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// MetaMessage represents a MIDI meta message. It can be created from the MIDI bytes of a message, by calling NewMetaMessage.
type MetaMessage struct {

	// MetaMsgType represents the message type of the MIDI meta message
	MetaMsgType

	// Data contains the bytes of the MiDI meta message
	Data []byte
}

func (m MetaMessage) Type() midi.MessageType {
	return m.MetaMsgType
}

func (m MetaMessage) Bytes() []byte {
	return m.Data
}

func (m MetaMessage) Is(t MetaMsgType) bool {
	return m.MetaMsgType == t
}

func NewMetaMessage(typ byte, data []byte) MetaMessage {
	return _MetaMessage(typ, data)
}

func _MetaMessage(typ byte, data []byte) MetaMessage {
	//fmt.Printf("NewMetaMessage %X % X\n", typ, data)
	b := []byte{byte(0xFF), typ}
	b = append(b, utils.VlqEncode(uint32(len(data)))...)
	if len(data) != 0 {
		b = append(b, data...)
	}

	var m MetaMessage
	m.MetaMsgType = GetMetaMsgType(typ)
	m.Data = b
	return m
}

// Meter returns the meter of a MetaTimeSigMsg.
// For other messages, it returns 0,0.
func (m MetaMessage) Meter(num, denom *uint8) (is bool) {
	return m.TimeSig(num, denom, nil, nil)
}

// metaData strips away the meta byte and the metatype byte and the varlength byte
func (m MetaMessage) metaDataWithoutVarlength() []byte {
	//fmt.Printf("original data: % X\n", m.Data)
	return m.Data[3:]
}

// TimeSig returns the numerator, denominator, clocksPerClick and demiSemiQuaverPerQuarter of a
// MetaTimeSigMsg. For other messages, it returns 0,0,0,0.
func (m MetaMessage) TimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter *uint8) (is bool) {
	if m.MetaMsgType != MetaTimeSig {
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

// String represents the Message as a string that contains the MsgType and its properties.
func (m MetaMessage) String() string {
	var bf bytes.Buffer
	fmt.Fprintf(&bf, m.MetaMsgType.String())

	var val1, val2 uint8
	//	var pitchabs uint16
	//	var pitchrel int16
	var text string
	var bpm float64
	//	var spp uint16

	switch {
	case m.Tempo(&bpm):
		fmt.Fprintf(&bf, " bpm: %0.2f", bpm)
	case m.Meter(&val1, &val2):
		fmt.Fprintf(&bf, " meter: %v/%v", val1, val2)
	default:
		switch m.MetaMsgType {
		case MetaLyric, MetaMarker, MetaCopyright, MetaText, MetaCuepoint, MetaDevice, MetaInstrument, MetaProgramName, MetaTrackName:
			m.text(&text)
			fmt.Fprintf(&bf, " text: %q", text)
		}
	}

	return bf.String()
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
func (m MetaMessage) Tempo(bpm *float64) (is bool) {
	if m.MetaMsgType != MetaTempo {
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
func (m MetaMessage) Lyric(text *string) (is bool) {
	if m.MetaMsgType != MetaLyric {
		return false
	}
	m.text(text)
	return true
}

// Copyright returns true if (and only if) the message is a MetaCopyrightMsg.
// Then it also extracts the data to the given arguments
func (m MetaMessage) Copyright(text *string) (is bool) {
	if !m.Is(MetaCopyright) {
		return false
	}
	m.text(text)
	return true
}

// Cuepoint returns true if (and only if) the message is a MetaCuepointMsg.
// Then it also extracts the data to the given arguments
func (m MetaMessage) Cuepoint(text *string) (is bool) {
	if !m.Is(MetaCuepoint) {
		return false
	}
	m.text(text)
	return true
}

// Device returns true if (and only if) the message is a MetaDeviceMsg.
// Then it also extracts the data to the given arguments
func (m MetaMessage) Device(text *string) (is bool) {
	if !m.Is(MetaDevice) {
		return false
	}
	m.text(text)
	return true
}

// Instrument returns true if (and only if) the message is a MetaInstrumentMsg.
// Then it also extracts the data to the given arguments
func (m MetaMessage) Instrument(text *string) (is bool) {
	if !m.Is(MetaInstrument) {
		return false
	}
	m.text(text)
	return true
}

// Marker returns true if (and only if) the message is a MetaMarkerMsg.
// Then it also extracts the data to the given arguments
func (m MetaMessage) Marker(text *string) (is bool) {
	if !m.Is(MetaMarker) {
		return false
	}
	m.text(text)
	return true
}

// ProgramName returns true if (and only if) the message is a MetaProgramNameMsg.
// Then it also extracts the data to the given arguments
func (m MetaMessage) ProgramName(text *string) (is bool) {
	if !m.Is(MetaProgramName) {
		return false
	}
	m.text(text)
	return true
}

// Text returns true if (and only if) the message is a MetaTextMsg.
// Then it also extracts the data to the given arguments
func (m MetaMessage) Text(text *string) (is bool) {
	switch m.MetaMsgType {
	case MetaLyric, MetaMarker, MetaCopyright, MetaText, MetaCuepoint, MetaDevice, MetaInstrument, MetaProgramName, MetaTrackName:
		m.text(text)
		return true
	default:
		return false
	}
}

// TrackName returns true if (and only if) the message is a MetaTrackNameMsg.
// Then it also extracts the data to the given arguments
func (m MetaMessage) TrackName(text *string) (is bool) {
	if !m.Is(MetaTrackName) {
		return false
	}
	m.text(text)
	return true
}

func (m MetaMessage) text(text *string) {
	*text, _ = utils.ReadText(bytes.NewReader(m.Data[2:]))
	return
}

func init() {
	var _ midi.Message = NewMetaCuepoint("test")
}

// Is returns if the given message type is part of the combination of message types
func Is(mt1, mt2 midi.MessageType) bool {
	return mt1.Category() == mt2.Category() && (mt1.Val()&mt2.Val()) != 0
}

func ReadMetaData(tp midi.MsgType, rd io.Reader) (data []byte, err error) {
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

var metaMessages = map[byte]MetaMsgType{
	byteEndOfTrack:        MetaEndOfTrack,
	byteSequenceNumber:    MetaSeqNumber,
	byteText:              MetaText,
	byteCopyright:         MetaCopyright,
	byteTrackSequenceName: MetaTrackName,
	byteInstrument:        MetaInstrument,
	byteLyric:             MetaLyric,
	byteMarker:            MetaMarker,
	byteCuepoint:          MetaCuepoint,
	byteMIDIChannel:       MetaChannel,
	byteDevicePort:        MetaDevice,
	byteMIDIPort:          MetaPort,
	byteTempo:             MetaTempo,
	byteTimeSignature:     MetaTimeSig,
	byteKeySignature:      MetaKeySig,
	byteSMPTEOffset:       MetaSMPTEOffset,
	byteSequencerSpecific: MetaSeqData,
	byteProgramName:       MetaProgramName,
}

// GetMetaMsgType returns the MsgType of a meta message. It should not be used by the end consumer.
func GetMetaMsgType(b byte) MetaMsgType {
	return metaMessages[b]
}

const bpmFac = 60000000

/*
// MetaMessage represents a SMF meta message.
type MetaMessage struct {

	// MsgType represents the message type of the MIDI message
	midi.MsgType

	// Data contains the bytes of the meta message
	Data []byte
}
*/

// newMetaMessage returns a new Message from the bytes of the message, by finding the correct type.
// If the type could not be found, the MsgType of the Message is UnknownMsg.

// MetaLyric returns the bytes of a lyric meta message
func NewMetaLyric(text string) MetaMessage {
	return _MetaMessage(byteLyric, []byte(text))
}

// MetaCopyright returns the bytes of a copyright meta message
func NewMetaCopyright(text string) MetaMessage {
	return _MetaMessage(byteCopyright, []byte(text))
}

// MetaChannel returns the bytes of a channel meta message
func NewMetaChannel(ch uint8) MetaMessage {
	return _MetaMessage(byteMIDIChannel, []byte{byte(ch)})
}

// MetaCuepoint returns the bytes of a cuepoint meta message
func NewMetaCuepoint(text string) MetaMessage {
	return _MetaMessage(byteCuepoint, []byte(text))
}

// MetaDevice returns the bytes of a device meta message
func NewMetaDevice(text string) MetaMessage {
	return _MetaMessage(byteDevicePort, []byte(text))
}

// EOT are the bytes of an End Of Track meta message. Don't use it directly.
var EOT = _MetaMessage(byteEndOfTrack, nil)

// MetaInstrument returns the bytes of a instrument meta message
func NewMetaInstrument(text string) MetaMessage {
	return _MetaMessage(byteInstrument, []byte(text))
}

// MetaMarker returns the bytes of a marker meta message
func NewMetaMarker(text string) MetaMessage {
	return _MetaMessage(byteMarker, []byte(text))
}

// MetaPort returns the bytes of a port meta message
func NewMetaPort(p uint8) MetaMessage {
	return _MetaMessage(byteMIDIPort, []byte{byte(p)})
}

// MetaProgram returns the bytes of a program meta message
func NewMetaProgram(text string) MetaMessage {
	return _MetaMessage(byteProgramName, []byte(text))
}

// MetaSequenceNo returns the bytes of a sequence number meta message
func NewMetaSequenceNo(no uint16) MetaMessage {
	var bf bytes.Buffer
	binary.Write(&bf, binary.BigEndian, no)
	return _MetaMessage(byteSequenceNumber, bf.Bytes())
}

// MetaSequencerData returns the bytes of a sequencer data meta message
func NewMetaSequencerData(data []byte) MetaMessage {
	return _MetaMessage(byteSequencerSpecific, data)
}

// MetaSMPTE returns the bytes of a SMPTE meta message
func NewMetaSMPTE(hour, minute, second, frame, fractionalFrame byte) MetaMessage {
	return _MetaMessage(byteSMPTEOffset, []byte{hour, minute, second, frame, fractionalFrame})
}

// MetaTempo returns the bytes of a tempo meta message for the given beats per minute.
func NewMetaTempo(bpm float64) MetaMessage {
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

// MetaText returns the bytes of a text meta message.
func NewMetaText(text string) MetaMessage {
	return _MetaMessage(byteText, []byte(text))
}

// MetaTrackSequenceName returns the bytes of a track sequence name meta message.
func NewMetaTrackSequenceName(text string) MetaMessage {
	return _MetaMessage(byteTrackSequenceName, []byte(text))
}

// MetaUndefined returns the bytes of an undefined meta message.
func NewMetaUndefined(typ byte, data []byte) MetaMessage {
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

// MetaKey returns the bytes of a key meta message.
func NewMetaKey(key uint8, isMajor bool, num uint8, isFlat bool) MetaMessage {
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

// MetaMeter returns the bytes of a time signature meta message.
func NewMetaMeter(num, denom uint8) MetaMessage {
	if denom == 0 {
		denom = 1
	}

	return NewMetaTimeSig(num, denom, 8, 8)
}

// MetaTimeSig returns the bytes of a time signature meta message.
func NewMetaTimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) MetaMessage {
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
