package key

import (
	"gitlab.com/gomidi/midi/v2"
)

func key(key, num uint8, isMajor, isFlat bool) []byte {
	return midi.MetaKey(key, isMajor, num, isFlat)
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
