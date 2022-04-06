package midi

import (
	"fmt"
)

/*
in order to be able to easy deal with tones and notes,
all notes are defined on the basic octave which happens to be
c1 (=60), so that it is easy to go below and above
*/

/*
const (
	C  Note = 60
	Db Note = 61
	D  Note = 62
	Eb Note = 63
	E  Note = 64
	F  Note = 65
	Gb Note = 66
	G  Note = 67
	Ab Note = 68
	A  Note = 69
	Bb Note = 70
	B  Note = 71
)
*/

func o(base uint8, oct int8) uint8 {
	if oct > 9 {
		oct = 9
	}

	if oct == 0 {
		return base
	}

	if oct < -1 {
		oct = -1
	}

	res := base + uint8(12*oct)
	if res > 127 {
		res -= 12
	}

	return res
}

// C returns the key for the MIDI note C in the given octave
func C(oct int8) uint8 {
	return o(12, oct)
}

// Db returns the key for the MIDI note Db in the given octave
func Db(oct int8) uint8 {
	return o(13, oct)
}

// D returns the key for the MIDI note D in the given octave
func D(oct int8) uint8 {
	return o(14, oct)
}

// Eb returns the key for the MIDI note Eb in the given octave
func Eb(oct int8) uint8 {
	return o(15, oct)
}

// E returns the key for the MIDI note E in the given octave
func E(oct int8) uint8 {
	return o(16, oct)
}

// F returns the key for the MIDI note F in the given octave
func F(oct int8) uint8 {
	return o(17, oct)
}

// Gb returns the key for the MIDI note Gb in the given octave
func Gb(oct int8) uint8 {
	return o(18, oct)
}

// G returns the key for the MIDI note G in the given octave
func G(oct int8) uint8 {
	return o(19, oct)
}

// Ab returns the key for the MIDI note Ab in the given octave
func Ab(oct int8) uint8 {
	return o(20, oct)
}

// A returns the key for the MIDI note A in the given octave
func A(oct int8) uint8 {
	return o(21, oct)
}

// Bb returns the key for the MIDI note Bb in the given octave
func Bb(oct int8) uint8 {
	return o(22, oct)
}

// B returns the key for the MIDI note B in the given octave
func B(oct int8) uint8 {
	return o(23, oct)
}

type Note uint8

func (n Note) Value() uint8 {
	return uint8(n)
}

func (n Note) Name() (name string) {
	switch n % 12 {
	case 0:
		name = "C"
	case 1:
		name = "Db"
	case 2:
		name = "D"
	case 3:
		name = "Eb"
	case 4:
		name = "E"
	case 5:
		name = "F"
	case 6:
		name = "Gb"
	case 7:
		name = "G"
	case 8:
		name = "Ab"
	case 9:
		name = "A"
	case 10:
		name = "Bb"
	case 11:
		name = "B"
	default:
		panic("unreachable")
	}

	return name
}

func (n Note) String() string {
	name := n.Name()
	if name != "" {
		name += fmt.Sprintf("%v", n.Octave())
	}
	return name
}

func (n Note) Octave() int8 {
	return int8(n/12) - 1
}

// Equal returns true if noteX is the same as noteY
// they may be in different octaves.
func (n Note) Is(o Note) bool {
	return n%12 == o%12
}
