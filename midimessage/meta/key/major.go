package key

import (
	"github.com/gomidi/midi/midimessage/meta"
)

// CMaj returns the MIDI key signature meta message for C Major
func CMaj() meta.KeySignature {
	return key(0, 0, true, false)
}

// DMaj returns the MIDI key signature meta message for D Major
func DMaj() meta.KeySignature {
	return key(2, 2, true, false)
}

// EMaj returns the MIDI key signature meta message for E Major
func EMaj() meta.KeySignature {
	return key(4, 4, true, false)
}

// FSharpMaj returns the MIDI key signature meta message for F# Major
func FSharpMaj() meta.KeySignature {
	return key(6, 6, true, false)
}

// GMaj returns the MIDI key signature meta message for G Major
func GMaj() meta.KeySignature {
	return key(7, 1, true, false)
}

// AMaj returns the MIDI key signature meta message for A Major
func AMaj() meta.KeySignature {
	return key(9, 3, true, false)
}

// BMaj returns the MIDI key signature meta message for B Major
func BMaj() meta.KeySignature {
	return key(11, 5, true, false)
}

// FMaj returns the MIDI key signature meta message for F Major
func FMaj() meta.KeySignature {
	return key(5, 1, true, true)
}

// BFlatMaj returns the MIDI key signature meta message for Bb Major
func BFlatMaj() meta.KeySignature {
	return key(10, 2, true, true)
}

// EFlatMaj returns the MIDI key signature meta message for Eb Major
func EFlatMaj() meta.KeySignature {
	return key(3, 3, true, true)
}

// AFlatMaj returns the MIDI key signature meta message for Ab Major
func AFlatMaj() meta.KeySignature {
	return key(8, 4, true, true)
}

// DFlatMaj returns the MIDI key signature meta message for Db Major
func DFlatMaj() meta.KeySignature {
	return key(1, 5, true, true)
}

// GFlatMaj returns the MIDI key signature meta message for Gb Major
func GFlatMaj() meta.KeySignature {
	return key(6, 6, true, true)
}

/*
func CFlatMaj() meta.KeySignature {
	return key(11, 5, true, false)
}

func CSharpMaj() meta.KeySignature {
	return key(1, 5, true, true)
}

func DSharpMaj() meta.KeySignature {
	return key(3, 3, true, true)
}

func ESharpMaj() meta.KeySignature {
	return key(5, 1, true, true)
}

func FFlatMaj() meta.KeySignature {
	return key(4, 4, true, false)
}

func GSharpMaj() meta.KeySignature {
	return key(8, 4, true, true)
}

func ASharpMaj() meta.KeySignature {
	return key(10, 2, true, true)
}

func BSharpMaj() meta.KeySignature {
	return key(0, 0, true, false)
}
*/
