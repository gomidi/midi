package meta

import (
	"fmt"
	"io"
	"math"
	"math/big"

	"github.com/gomidi/midi/internal/midilib"
)

const bpmFac = 60000000

// BPM returns the meta tempo message that corresponds to the given bpm (beats per minute) value
func BPM(bpm uint32) Tempo {
	return FractionalBPM(float64(bpm))
}

// FractionalBPM returns the meta tempo message that corresponds to the given fractional bpm (beats per minute) value
func FractionalBPM(fbpm float64) Tempo {
	return Tempo(uint32(math.Round(bpmFac / fbpm)))
}

// Tempo represents a MIDI tempo (change) message in microseconds per crotchet
type Tempo uint32

// BPM returns the tempo in beats per minute
func (m Tempo) BPM() uint32 {
	return uint32(math.Round(m.FractionalBPM()))
}

// MuSecPerQN returns the tempo in microseconds per quarternote
func (m Tempo) MuSecPerQN() uint32 {
	return uint32(m)
}

// FractionalBPM returns the tempo in fractional beats per minute
func (m Tempo) FractionalBPM() float64 {
	return float64(bpmFac) / float64(m)
}

// String represents the tempo message as a string (for debugging)
func (m Tempo) String() string {
	return fmt.Sprintf("%T BPM: %0.2f", m, m.FractionalBPM())
}

// Raw returns the raw MIDI data
func (m Tempo) Raw() []byte {
	r := uint32(m)
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

	return (&metaMessage{
		Typ:  byteTempo,
		Data: b,
	}).Bytes()
}

func (m Tempo) meta() {}

func (m Tempo) readFrom(rd io.Reader) (Message, error) {
	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 3 {
		err = unexpectedMessageLengthError("Tempo expected length 3")
		return nil, err
	}

	var microsecondsPerCrotchet uint32
	microsecondsPerCrotchet, err = midilib.ReadUint24(rd)
	if err != nil {
		return nil, err
	}

	return Tempo(microsecondsPerCrotchet), nil
}
