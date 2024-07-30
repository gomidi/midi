package smfimage

import "strings"

type Note byte
type Interval byte

const (
	Prime      Interval = 0
	MinSecond  Interval = 1
	MajSecond  Interval = 2
	MinThird   Interval = 3
	MajThird   Interval = 4
	Fourth     Interval = 5
	Tritone    Interval = 6
	Fifth      Interval = 7
	MinSixth   Interval = 8
	MajSixth   Interval = 9
	MinSeventh Interval = 10
	MajSeventh Interval = 11

	Unison = Prime
	Octave = Prime

	C      Note = 0
	CSharp Note = 1
	D      Note = 2
	DSharp Note = 3
	E      Note = 4
	F      Note = 5
	FSharp Note = 6
	G      Note = 7
	GSharp Note = 8
	A      Note = 9
	ASharp Note = 10
	B      Note = 11

	DFlat  = CSharp
	EFlat  = DSharp
	ESharp = F
	FFlat  = E
	GFlat  = FSharp
	AFlat  = GSharp
	BFlat  = ASharp
	BSharp = C
	CFlat  = B
)

var noteToNumber = map[string]int{
	"c":  0,
	"c#": 1, "cis": 1, "des": 1, "db": 1,
	"d":  2,
	"d#": 3, "dis": 3, "es": 3, "eb": 3,
	"e": 4, "fes": 4, "fb": 4,
	"f": 5, "eis": 5, "e#": 5,
	"f#": 6, "fis": 6, "ges": 6, "gb": 6,
	"g":  7,
	"g#": 8, "gis": 8, "as": 8, "ab": 8,
	"a":  9,
	"a#": 10, "ais": 10, "bb": 10,
	"b": 11,
}

func NoteToNumber(name string) int {
	key := strings.ToLower(strings.TrimSpace(name))

	if k, has := noteToNumber[key]; has {
		return k
	}

	return -1
}
