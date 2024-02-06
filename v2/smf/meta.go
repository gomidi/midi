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

const (
	MetaMsg midi.Type = -5
)

const (

	// MetaChannelMsg is a MIDI channel meta message
	MetaChannelMsg midi.Type = 70 + iota

	// MetaCopyrightMsg is a MIDI copyright meta message
	MetaCopyrightMsg

	// MetaCuepointMsg is a MIDI cuepoint meta message
	MetaCuepointMsg

	// MetaDeviceMsg is a MIDI device meta message
	MetaDeviceMsg

	// MetaEndOfTrackMsg is a MIDI end of track meta message
	MetaEndOfTrackMsg

	// MetaInstrumentMsg is a MIDI instrument meta message
	MetaInstrumentMsg

	// MetaKeySigMsg is a MIDI key signature meta message
	MetaKeySigMsg

	// MetaLyricMsg is a MIDI lyrics meta message
	MetaLyricMsg

	// MetaTextMsg is a MIDI text meta message
	MetaTextMsg

	// MetaMarkerMsg is a MIDI marker meta message
	MetaMarkerMsg

	// MetaPortMsg is a MIDI port meta message
	MetaPortMsg

	// MetaSeqNumberMsg is a MIDI sequencer number meta message
	MetaSeqNumberMsg

	// MetaSeqDataMsg is a MIDI sequencer data meta message
	MetaSeqDataMsg

	// MetaTempoMsg is a MIDI tempo meta message
	MetaTempoMsg

	// MetaTimeSigMsg is a MIDI time signature meta message
	MetaTimeSigMsg

	// MetaTrackNameMsg is a MIDI track name meta message
	MetaTrackNameMsg

	// MetaSMPTEOffsetMsg is a MIDI smpte offset meta message
	MetaSMPTEOffsetMsg

	// MetaUndefinedMsg is an undefined MIDI meta message
	MetaUndefinedMsg

	// MetaProgramNameMsg is a MIDI program name meta message
	MetaProgramNameMsg
)

var msgTypeString = map[midi.Type]string{
	MetaMsg:            "Meta",
	MetaChannelMsg:     "MetaChannel",
	MetaCopyrightMsg:   "MetaCopyright",
	MetaCuepointMsg:    "MetaCuepoint",
	MetaDeviceMsg:      "MetaDevice",
	MetaEndOfTrackMsg:  "MetaEndOfTrack",
	MetaInstrumentMsg:  "MetaInstrument",
	MetaKeySigMsg:      "MetaKeySig",
	MetaLyricMsg:       "MetaLyric",
	MetaTextMsg:        "MetaText",
	MetaMarkerMsg:      "MetaMarker",
	MetaPortMsg:        "MetaPort",
	MetaSeqNumberMsg:   "MetaSeqNumber",
	MetaSeqDataMsg:     "MetaSeqData",
	MetaTempoMsg:       "MetaTempo",
	MetaTimeSigMsg:     "MetaTimeSig",
	MetaTrackNameMsg:   "MetaTrackName",
	MetaSMPTEOffsetMsg: "MetaSMPTEOffset",
	MetaUndefinedMsg:   "MetaUndefined",
	MetaProgramNameMsg: "MetaProgramName",
}

func init() {
	for ty, name := range msgTypeString {
		midi.AddTypeName(ty, name)
	}
}

func readMetaData(tp midi.Type, rd io.Reader) (data []byte, err error) {
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
func getMetaType(b byte) midi.Type {
	return metaMessages[b]
}

const bpmFac = 60000000

// MetaLyric returns a lyric meta message
func MetaLyric(text string) Message {
	return _MetaMessage(byteLyric, []byte(text))
}

// MetaCopyright returns a copyright meta message
func MetaCopyright(text string) Message {
	return _MetaMessage(byteCopyright, []byte(text))
}

// MetaChannel returns a channel meta message
func MetaChannel(ch uint8) Message {
	return _MetaMessage(byteMIDIChannel, []byte{byte(ch)})
}

// MetaCuepoint returns a cuepoint meta message
func MetaCuepoint(text string) Message {
	return _MetaMessage(byteCuepoint, []byte(text))
}

// MetaDevice returns a device meta message
func MetaDevice(text string) Message {
	return _MetaMessage(byteDevicePort, []byte(text))
}

// EOT is an End Of Track meta message. Don't use it directly.
var EOT = _MetaMessage(byteEndOfTrack, nil)

// MetaInstrument returns an instrument meta message
func MetaInstrument(text string) Message {
	return _MetaMessage(byteInstrument, []byte(text))
}

// MetaMarker returns a marker meta message
func MetaMarker(text string) Message {
	return _MetaMessage(byteMarker, []byte(text))
}

// MetaPort returns a port meta message
func MetaPort(p uint8) Message {
	return _MetaMessage(byteMIDIPort, []byte{byte(p)})
}

// MetaProgram returns a program meta message
func MetaProgram(text string) Message {
	return _MetaMessage(byteProgramName, []byte(text))
}

// MetaSequenceNo returns a sequence number meta message
func MetaSequenceNo(no uint16) Message {
	var bf bytes.Buffer
	binary.Write(&bf, binary.BigEndian, no)
	//fmt.Printf("%X\n", bf.Bytes())
	m := _MetaMessage(byteSequenceNumber, bf.Bytes())
	//fmt.Printf("%X\n", m.Bytes())
	return m
}

// MetaSequencerData returns a sequencer data meta message
func MetaSequencerData(data []byte) Message {
	return _MetaMessage(byteSequencerSpecific, data)
}

// MetaSMPTE returns a SMPTE meta message
func MetaSMPTE(hour, minute, second, frame, fractionalFrame byte) Message {
	return _MetaMessage(byteSMPTEOffset, []byte{hour, minute, second, frame, fractionalFrame})
}

// MetaTempo returns a tempo meta message for the given beats per minute.
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

// MetaText returns a text meta message.
func MetaText(text string) Message {
	return _MetaMessage(byteText, []byte(text))
}

// MetaTrackSequenceName returns a track sequence name meta message.
func MetaTrackSequenceName(text string) Message {
	return _MetaMessage(byteTrackSequenceName, []byte(text))
}

// MetaUndefined returns an undefined meta message.
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

// MetaKey returns a key meta message.
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

// MetaMeter returns a time signature meta message.
func MetaMeter(num, denom uint8) Message {
	if denom == 0 {
		denom = 1
	}

	return MetaTimeSig(num, denom, 8, 8)
}

// MetaTimeSig returns a time signature meta message.
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
