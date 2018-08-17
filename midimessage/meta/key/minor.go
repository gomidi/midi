package key

import (
	"github.com/gomidi/midi/midimessage/meta"
)

// AMin returns the MIDI key signature meta message for A Minor
func AMin() meta.KeySignature {
	return key(9, 0, false, false)
}

// BMin returns the MIDI key signature meta message for B Minor
func BMin() meta.KeySignature {
	return key(11, 2, false, false)
}

// CSharpMin returns the MIDI key signature meta message for C# Minor
func CSharpMin() meta.KeySignature {
	return key(1, 4, false, false)
}

// DSharpMin returns the MIDI key signature meta message for D# Minor
func DSharpMin() meta.KeySignature {
	return key(3, 6, false, false)
}

// EMin returns the MIDI key signature meta message for E Minor
func EMin() meta.KeySignature {
	return key(4, 1, false, false)
}

// FSharpMin returns the MIDI key signature meta message for F# Minor
func FSharpMin() meta.KeySignature {
	return key(6, 3, false, false)
}

// GSharpMin returns the MIDI key signature meta message for G# Minor
func GSharpMin() meta.KeySignature {
	return key(8, 5, false, false)
}

// DMin returns the MIDI key signature meta message for D Minor
func DMin() meta.KeySignature {
	return key(2, 1, false, true)
}

// GMin returns the MIDI key signature meta message for G Minor
func GMin() meta.KeySignature {
	return key(7, 2, false, true)
}

// CMin returns the MIDI key signature meta message for C Minor
func CMin() meta.KeySignature {
	return key(0, 3, false, true)
}

// FMin returns the MIDI key signature meta message for F Minor
func FMin() meta.KeySignature {
	return key(5, 4, false, true)
}

// BFlatMin returns the MIDI key signature meta message for Bb Minor
func BFlatMin() meta.KeySignature {
	return key(10, 5, false, true)
}

// EFlatMin returns the MIDI key signature meta message for Eb Minor
func EFlatMin() meta.KeySignature {
	return key(3, 6, false, true)
}

/*
func CFlatMin() meta.KeySignature {
	return key(11, 2, false, false)
}


func ESharpMin() meta.KeySignature {
	return key(5, 4, false, true)
}

func FFlatMin() meta.KeySignature {
	return key(4, 1, false, false)
}

func GFlatMin() meta.KeySignature {
	return key(5, 3, false, false)
}

func AFlatMin() meta.KeySignature {
	return key(8, 5, false, false)
}

func ASharpMin() meta.KeySignature {
	return key(10, 5, false, true)
}

func BSharpMin() meta.KeySignature {
	return key(0, 3, false, true)
}

func DFlatMin() meta.KeySignature {
	return key(1, 4, false, false)
}



*/
