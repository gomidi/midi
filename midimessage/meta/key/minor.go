package key

import (
	"github.com/gomidi/midi/midimessage/meta"
)

// AMin returns the MIDI key signature meta message for A Minor
func AMin() meta.Key {
	return key(9, 0, false, false)
}

// BMin returns the MIDI key signature meta message for B Minor
func BMin() meta.Key {
	return key(11, 2, false, false)
}

// CSharpMin returns the MIDI key signature meta message for C# Minor
func CSharpMin() meta.Key {
	return key(1, 4, false, false)
}

// DSharpMin returns the MIDI key signature meta message for D# Minor
func DSharpMin() meta.Key {
	return key(3, 6, false, false)
}

// EMin returns the MIDI key signature meta message for E Minor
func EMin() meta.Key {
	return key(4, 1, false, false)
}

// FSharpMin returns the MIDI key signature meta message for F# Minor
func FSharpMin() meta.Key {
	return key(6, 3, false, false)
}

// GSharpMin returns the MIDI key signature meta message for G# Minor
func GSharpMin() meta.Key {
	return key(8, 5, false, false)
}

// DMin returns the MIDI key signature meta message for D Minor
func DMin() meta.Key {
	return key(2, 1, false, true)
}

// GMin returns the MIDI key signature meta message for G Minor
func GMin() meta.Key {
	return key(7, 2, false, true)
}

// CMin returns the MIDI key signature meta message for C Minor
func CMin() meta.Key {
	return key(0, 3, false, true)
}

// FMin returns the MIDI key signature meta message for F Minor
func FMin() meta.Key {
	return key(5, 4, false, true)
}

// BFlatMin returns the MIDI key signature meta message for Bb Minor
func BFlatMin() meta.Key {
	return key(10, 5, false, true)
}

// EFlatMin returns the MIDI key signature meta message for Eb Minor
func EFlatMin() meta.Key {
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
