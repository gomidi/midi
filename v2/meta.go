package midi

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"math/big"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

func ReadMetaData(tp MessageType, rd io.Reader) (data []byte, err error) {
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

var metaMessages = map[byte]MessageType{
	byteEndOfTrack:        Set(MetaMsg, MetaEndOfTrackMsg),
	byteSequenceNumber:    Set(MetaMsg, MetaSeqNumberMsg),
	byteText:              Set(MetaMsg, MetaTextMsg),
	byteCopyright:         Set(MetaMsg, MetaCopyrightMsg),
	byteTrackSequenceName: Set(MetaMsg, MetaTrackNameMsg),
	byteInstrument:        Set(MetaMsg, MetaInstrumentMsg),
	byteLyric:             Set(MetaMsg, MetaLyricMsg),
	byteMarker:            Set(MetaMsg, MetaMarkerMsg),
	byteCuepoint:          Set(MetaMsg, MetaCuepointMsg),
	byteMIDIChannel:       Set(MetaMsg, MetaChannelMsg),
	byteDevicePort:        Set(MetaMsg, MetaDeviceMsg),
	byteMIDIPort:          Set(MetaMsg, MetaPortMsg),
	byteTempo:             Set(MetaMsg, MetaTempoMsg),
	byteTimeSignature:     Set(MetaMsg, MetaTimeSigMsg),
	byteKeySignature:      Set(MetaMsg, MetaKeySigMsg),
	byteSMPTEOffset:       Set(MetaMsg, MetaSMPTEOffsetMsg),
	byteSequencerSpecific: Set(MetaMsg, MetaSeqDataMsg),
	byteProgramName:       Set(MetaMsg, MetaProgramNameMsg),
}

const bpmFac = 60000000

func MetaMessage(typ byte, data []byte) []byte {
	b := []byte{byte(0xFF), typ}
	b = append(b, utils.VlqEncode(uint32(len(data)))...)
	if len(data) != 0 {
		b = append(b, data...)
	}
	return b
}

func MetaLyric(text string) []byte {
	return MetaMessage(byteLyric, []byte(text))
}

func MetaCopyright(text string) []byte {
	return MetaMessage(byteCopyright, []byte(text))
}

func MetaChannel(ch uint8) []byte {
	return MetaMessage(byteMIDIChannel, []byte{byte(ch)})
}

func MetaCuepoint(text string) []byte {
	return MetaMessage(byteCuepoint, []byte(text))
}

func MetaDevice(text string) []byte {
	return MetaMessage(byteDevicePort, []byte(text))
}

// EndOfTrack you should not use this. Use the smf package instead.
var EOT = MetaMessage(byteEndOfTrack, nil)

func MetaInstrument(text string) []byte {
	return MetaMessage(byteInstrument, []byte(text))
}

func MetaMarker(text string) []byte {
	return MetaMessage(byteMarker, []byte(text))
}

func MetaPort(p uint8) []byte {
	return MetaMessage(byteMIDIPort, []byte{byte(p)})
}

func MetaProgram(text string) []byte {
	return MetaMessage(byteProgramName, []byte(text))
}

func MetaSequenceNo(no uint16) []byte {
	var bf bytes.Buffer
	binary.Write(&bf, binary.BigEndian, no)
	return MetaMessage(byteSequenceNumber, bf.Bytes())
}

func MetaSequencerData(data []byte) []byte {
	return MetaMessage(byteSequencerSpecific, data)
}

func MetaSMPTE(hour, minute, second, frame, fractionalFrame byte) []byte {
	return MetaMessage(byteSMPTEOffset, []byte{hour, minute, second, frame, fractionalFrame})
}

func MetaTempo(bpm float64) []byte {
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

	return MetaMessage(byteTempo, b)
}

func MetaText(text string) []byte {
	return MetaMessage(byteText, []byte(text))
}

func MetaTrackSequenceName(text string) []byte {
	return MetaMessage(byteTrackSequenceName, []byte(text))
}

func MetaUndefined(typ byte, data []byte) []byte {
	return MetaMessage(typ, data)
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

func MetaKey(key uint8, isMajor bool, num uint8, isFlat bool) []byte {
	mi := int8(0)
	if !isMajor {
		mi = 1
	}
	sf := int8(num)

	if isFlat {
		sf = sf * (-1)
	}

	return MetaMessage(byteKeySignature, []byte{byte(sf), byte(mi)})
}

func MetaMeter(num, denom uint8) []byte {
	if denom == 0 {
		denom = 1
	}

	return MetaTimeSig(num, denom, 8, 8)
}

// Raw returns the raw MIDI data
func MetaTimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) []byte {
	cpcl := clocksPerClick
	if cpcl == 0 {
		cpcl = byte(8)
	}

	dsqpq := demiSemiQuaverPerQuarter
	if dsqpq == 0 {
		dsqpq = byte(8)
	}

	var denom = dec2binDenom(denominator)

	return MetaMessage(byteTimeSignature, []byte{numerator, denom, cpcl, dsqpq})

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

// CMaj returns the MIDI key signature meta message for C Major
func CMaj() []byte {
	return key(0, 0, true, false)
}

// DMaj returns the MIDI key signature meta message for D Major
func DMaj() []byte {
	return key(2, 2, true, false)
}

// EMaj returns the MIDI key signature meta message for E Major
func EMaj() []byte {
	return key(4, 4, true, false)
}

// FSharpMaj returns the MIDI key signature meta message for F# Major
func FSharpMaj() []byte {
	return key(6, 6, true, false)
}

// GMaj returns the MIDI key signature meta message for G Major
func GMaj() []byte {
	return key(7, 1, true, false)
}

// AMaj returns the MIDI key signature meta message for A Major
func AMaj() []byte {
	return key(9, 3, true, false)
}

// BMaj returns the MIDI key signature meta message for B Major
func BMaj() []byte {
	return key(11, 5, true, false)
}

// FMaj returns the MIDI key signature meta message for F Major
func FMaj() []byte {
	return key(5, 1, true, true)
}

// BFlatMaj returns the MIDI key signature meta message for Bb Major
func BFlatMaj() []byte {
	return key(10, 2, true, true)
}

// EFlatMaj returns the MIDI key signature meta message for Eb Major
func EFlatMaj() []byte {
	return key(3, 3, true, true)
}

// AFlatMaj returns the MIDI key signature meta message for Ab Major
func AFlatMaj() []byte {
	return key(8, 4, true, true)
}

// DFlatMaj returns the MIDI key signature meta message for Db Major
func DFlatMaj() []byte {
	return key(1, 5, true, true)
}

// GFlatMaj returns the MIDI key signature meta message for Gb Major
func GFlatMaj() []byte {
	return key(6, 6, true, true)
}

func key(key, num uint8, isMajor, isFlat bool) []byte {
	return MetaKey(key, isMajor, num, isFlat)
}

/*
func CFlatMaj() meta.Key {
	return key(11, 5, true, false)
}

func CSharpMaj() meta.Key {
	return key(1, 5, true, true)
}

func DSharpMaj() meta.Key {
	return key(3, 3, true, true)
}

func ESharpMaj() meta.Key {
	return key(5, 1, true, true)
}

func FFlatMaj() meta.Key {
	return key(4, 4, true, false)
}

func GSharpMaj() meta.Key {
	return key(8, 4, true, true)
}

func ASharpMaj() meta.Key {
	return key(10, 2, true, true)
}

func BSharpMaj() meta.Key {
	return key(0, 0, true, false)
}
*/

// AMin returns the MIDI key signature meta message for A Minor
func AMin() []byte {
	return key(9, 0, false, false)
}

// BMin returns the MIDI key signature meta message for B Minor
func BMin() []byte {
	return key(11, 2, false, false)
}

// CSharpMin returns the MIDI key signature meta message for C# Minor
func CSharpMin() []byte {
	return key(1, 4, false, false)
}

// DSharpMin returns the MIDI key signature meta message for D# Minor
func DSharpMin() []byte {
	return key(3, 6, false, false)
}

// EMin returns the MIDI key signature meta message for E Minor
func EMin() []byte {
	return key(4, 1, false, false)
}

// FSharpMin returns the MIDI key signature meta message for F# Minor
func FSharpMin() []byte {
	return key(6, 3, false, false)
}

// GSharpMin returns the MIDI key signature meta message for G# Minor
func GSharpMin() []byte {
	return key(8, 5, false, false)
}

// DMin returns the MIDI key signature meta message for D Minor
func DMin() []byte {
	return key(2, 1, false, true)
}

// GMin returns the MIDI key signature meta message for G Minor
func GMin() []byte {
	return key(7, 2, false, true)
}

// CMin returns the MIDI key signature meta message for C Minor
func CMin() []byte {
	return key(0, 3, false, true)
}

// FMin returns the MIDI key signature meta message for F Minor
func FMin() []byte {
	return key(5, 4, false, true)
}

// BFlatMin returns the MIDI key signature meta message for Bb Minor
func BFlatMin() []byte {
	return key(10, 5, false, true)
}

// EFlatMin returns the MIDI key signature meta message for Eb Minor
func EFlatMin() []byte {
	return key(3, 6, false, true)
}

/*
func CFlatMin() meta.Key {
	return key(11, 2, false, false)
}


func ESharpMin() meta.Key {
	return key(5, 4, false, true)
}

func FFlatMin() meta.Key {
	return key(4, 1, false, false)
}

func GFlatMin() meta.Key {
	return key(5, 3, false, false)
}

func AFlatMin() meta.Key {
	return key(8, 5, false, false)
}

func ASharpMin() meta.Key {
	return key(10, 5, false, true)
}

func BSharpMin() meta.Key {
	return key(0, 3, false, true)
}

func DFlatMin() meta.Key {
	return key(1, 4, false, false)
}



*/
