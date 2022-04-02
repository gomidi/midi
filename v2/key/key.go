package key

import (
	"gitlab.com/gomidi/midi/v2/smf"
)

func key(key, num uint8, isMajor, isFlat bool) smf.Message {
	return smf.MetaKey(key, isMajor, num, isFlat)
}

// CMaj returns the MIDI key signature meta message for C Major
func CMaj() smf.Message {
	return key(0, 0, true, false)
}

// DMaj returns the MIDI key signature meta message for D Major
func DMaj() smf.Message {
	return key(2, 2, true, false)
}

// EMaj returns the MIDI key signature meta message for E Major
func EMaj() smf.Message {
	return key(4, 4, true, false)
}

// FsharpMaj returns the MIDI key signature meta message for F# Major
func FsharpMaj() smf.Message {
	return key(6, 6, true, false)
}

// GMaj returns the MIDI key signature meta message for G Major
func GMaj() smf.Message {
	return key(7, 1, true, false)
}

// AMaj returns the MIDI key signature meta message for A Major
func AMaj() smf.Message {
	return key(9, 3, true, false)
}

// BMaj returns the MIDI key signature meta message for B Major
func BMaj() smf.Message {
	return key(11, 5, true, false)
}

// FMaj returns the MIDI key signature meta message for F Major
func FMaj() smf.Message {
	return key(5, 1, true, true)
}

// BbMaj returns the MIDI key signature meta message for Bb Major
func BbMaj() smf.Message {
	return key(10, 2, true, true)
}

// EbMaj returns the MIDI key signature meta message for Eb Major
func EbMaj() smf.Message {
	return key(3, 3, true, true)
}

// AbMaj returns the MIDI key signature meta message for Ab Major
func AbMaj() smf.Message {
	return key(8, 4, true, true)
}

// DbMaj returns the MIDI key signature meta message for Db Major
func DbMaj() smf.Message {
	return key(1, 5, true, true)
}

// GbMaj returns the MIDI key signature meta message for Gb Major
func GbMaj() smf.Message {
	return key(6, 6, true, true)
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
func AMin() smf.Message {
	return key(9, 0, false, false)
}

// BMin returns the MIDI key signature meta message for B Minor
func BMin() smf.Message {
	return key(11, 2, false, false)
}

// CsharpMin returns the MIDI key signature meta message for C# Minor
func CsharpMin() smf.Message {
	return key(1, 4, false, false)
}

// DsharpMin returns the MIDI key signature meta message for D# Minor
func DsharpMin() smf.Message {
	return key(3, 6, false, false)
}

// EMin returns the MIDI key signature meta message for E Minor
func EMin() smf.Message {
	return key(4, 1, false, false)
}

// FsharpMin returns the MIDI key signature meta message for F# Minor
func FsharpMin() smf.Message {
	return key(6, 3, false, false)
}

// GsharpMin returns the MIDI key signature meta message for G# Minor
func GsharpMin() smf.Message {
	return key(8, 5, false, false)
}

// DMin returns the MIDI key signature meta message for D Minor
func DMin() smf.Message {
	return key(2, 1, false, true)
}

// GMin returns the MIDI key signature meta message for G Minor
func GMin() smf.Message {
	return key(7, 2, false, true)
}

// CMin returns the MIDI key signature meta message for C Minor
func CMin() smf.Message {
	return key(0, 3, false, true)
}

// FMin returns the MIDI key signature meta message for F Minor
func FMin() smf.Message {
	return key(5, 4, false, true)
}

// BbMin returns the MIDI key signature meta message for Bb Minor
func BbMin() smf.Message {
	return key(10, 5, false, true)
}

// EbMin returns the MIDI key signature meta message for Eb Minor
func EbMin() smf.Message {
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
