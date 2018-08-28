package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

/* http://www.somascape.org/midi/tech/mfile.html
Key Signature

FF 59 02 sf mi

sf is a byte specifying the number of flats (-ve) or sharps (+ve) that identifies the key signature (-7 = 7 flats, -1 = 1 flat, 0 = key of C, 1 = 1 sharp, etc).
mi is a byte specifying a major (0) or minor (1) key.

For a format 1 MIDI file, Key Signature Meta events should only occur within the first MTrk chunk.

*/

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

// Key sets the key/scale of the SMF file.
// If you want a more comfortable way to set the key, use the key subpackage.
type Key struct {
	Key     uint8
	IsMajor bool
	Num     uint8
	//	SharpsOrFlats int8
	IsFlat bool
}

/*
// NewKeySignature returns a key signature event.
// key is the key of the scale (C=0 add the corresponding number of semitones). ismajor indicates whether it is a major or minor scale
// num is the number of accidentals. isflat indicates whether the accidentals are flats or sharps
func NewKeySignature(key uint8, ismajor bool, num uint8, isflat bool) KeySignature {
	return KeySignature{Key: key, IsMajor: ismajor, Num: num, IsFlat: isflat}
}
*/

// Raw returns the raw MIDI data
func (m Key) Raw() []byte {
	mi := int8(0)
	if !m.IsMajor {
		mi = 1
	}
	sf := int8(m.Num)

	if m.IsFlat {
		sf = sf * (-1)
	}

	return (&metaMessage{
		Typ:  byteKeySignature,
		Data: []byte{byte(sf), byte(mi)},
	}).Bytes()
}

// String represents the key signature message as a string (for debugging)
func (m Key) String() string {
	return fmt.Sprintf("%T: %s", m, m.Text())
}

var keyNotes = map[uint8]string{
	degreeC:  "C",
	degreeD:  "D",
	degreeE:  "E",
	degreeF:  "F",
	degreeG:  "G",
	degreeA:  "A",
	degreeB:  "B",
	degreeCs: "C♯",
	degreeDs: "D♯",
	degreeFs: "F♯",
	degreeGs: "G♯",
	degreeAs: "A♯",
}

var keyNotesFlat = map[uint8]string{
	degreeCs: "D♭",
	degreeDs: "E♭",
	degreeFs: "G♭",
	degreeGs: "A♭",
	degreeAs: "B♭",
}

// Note returns the note of the key signature as a string, e.g. C♯ or E♭
func (m Key) Note() (note string) {
	if m.IsFlat {
		if nt, has := keyNotesFlat[m.Key]; has {
			return nt
		}
	}

	return keyNotes[m.Key]
}

// Text returns a the text of the key signature
func (m Key) Text() string {
	if m.IsMajor {
		return m.Note() + " maj."
	}

	return m.Note() + " min."
}

func (m Key) readFrom(rd io.Reader) (Message, error) {

	// fmt.Println("Key signature")
	// TODO TEST
	var sharpsOrFlats int8
	var mode uint8

	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 2 {
		err = unexpectedMessageLengthError("KeySignature expected length 2")
		return nil, err
	}

	// Signed int, positive is sharps, negative is flats.
	var b byte
	b, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	sharpsOrFlats = int8(b)

	// Mode is Major or Minor.
	mode, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	num := sharpsOrFlats
	if num < 0 {
		num = num * (-1)
	}

	key := midilib.KeyFromSharpsOrFlats(sharpsOrFlats, mode)

	return Key{
		Key:     key,
		Num:     uint8(num),
		IsMajor: mode == majorMode,
		IsFlat:  sharpsOrFlats < 0,
	}, nil

}

func (m Key) meta() {}
