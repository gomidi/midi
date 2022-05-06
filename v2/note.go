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

func o(base uint8, oct uint8) uint8 {
	if oct > 10 {
		oct = 10
	}

	if oct == 0 {
		return base
	}

	res := base + uint8(12*oct)
	if res > 127 {
		res -= 12
	}

	return res
}

// C returns the key for the MIDI note C in the given octave
func C(oct uint8) uint8 {
	return o(0, oct)
}

// Db returns the key for the MIDI note Db in the given octave
func Db(oct uint8) uint8 {
	return o(1, oct)
}

// D returns the key for the MIDI note D in the given octave
func D(oct uint8) uint8 {
	return o(2, oct)
}

// Eb returns the key for the MIDI note Eb in the given octave
func Eb(oct uint8) uint8 {
	return o(3, oct)
}

// E returns the key for the MIDI note E in the given octave
func E(oct uint8) uint8 {
	return o(4, oct)
}

// F returns the key for the MIDI note F in the given octave
func F(oct uint8) uint8 {
	return o(5, oct)
}

// Gb returns the key for the MIDI note Gb in the given octave
func Gb(oct uint8) uint8 {
	return o(6, oct)
}

// G returns the key for the MIDI note G in the given octave
func G(oct uint8) uint8 {
	return o(7, oct)
}

// Ab returns the key for the MIDI note Ab in the given octave
func Ab(oct uint8) uint8 {
	return o(8, oct)
}

// A returns the key for the MIDI note A in the given octave
func A(oct uint8) uint8 {
	return o(9, oct)
}

// Bb returns the key for the MIDI note Bb in the given octave
func Bb(oct uint8) uint8 {
	return o(10, oct)
}

// B returns the key for the MIDI note B in the given octave
func B(oct uint8) uint8 {
	return o(11, oct)
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

func (n Note) Octave() uint8 {
	return uint8(n / 12)
}

// Equal returns true if noteX is the same as noteY
// they may be in different octaves.
func (n Note) Is(o Note) bool {
	return n%12 == o%12
}
