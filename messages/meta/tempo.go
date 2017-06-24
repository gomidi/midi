package meta

import (
	"fmt"
	"github.com/gomidi/midi/internal/lib"
	"io"
	"math/big"
)

// value is equal to BPM
type Tempo uint32

func (t Tempo) BPM() uint32 {
	return uint32(t)
}

func (m Tempo) String() string {
	return fmt.Sprintf("%T BPM: %v", m, m.BPM())
}

func (m Tempo) Raw() []byte {

	f := float64(60000000) / float64(m.BPM())

	muSecPerQuarterNote := uint32(f)

	if muSecPerQuarterNote > 0xFFFFFF {
		muSecPerQuarterNote = 0xFFFFFF
	}
	b4 := big.NewInt(int64(muSecPerQuarterNote)).Bytes()
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
	// TODO TEST
	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 3 {
		err = lib.UnexpectedMessageLengthError("Tempo expected length 3")
		return nil, err
	}

	var microsecondsPerCrotchet uint32
	microsecondsPerCrotchet, err = lib.ReadUint24(rd)

	if err != nil {
		return nil, err
	}

	// Also beats per minute
	var bpm uint32
	bpm = 60000000 / microsecondsPerCrotchet

	return Tempo(bpm), nil
}
