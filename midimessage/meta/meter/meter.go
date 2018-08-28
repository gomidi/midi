package meter

import (
	"github.com/gomidi/midi/midimessage/meta"
)

// Meter returns the MIDI time signature meta message for the given
// numerator and denominator.
func Meter(num, denom uint8) meta.TimeSig {
	if denom == 0 {
		denom = 1
	}

	return meta.TimeSig{
		Numerator:                num,
		Denominator:              denom,
		ClocksPerClick:           8,
		DemiSemiQuaverPerQuarter: 8,
	}
}

// M2_4 returns  the MIDI time signature meta message for a 2/4 meter
func M2_4() meta.TimeSig {
	return Meter(2, 4)
}

// M4_4 returns  the MIDI time signature meta message for a 4/4 meter
func M4_4() meta.TimeSig {
	return Meter(4, 4)
}

// M3_4 returns  the MIDI time signature meta message for a 3/4 meter
func M3_4() meta.TimeSig {
	return Meter(3, 4)
}

// M6_8 returns  the MIDI time signature meta message for a 6/8 meter
func M6_8() meta.TimeSig {
	return Meter(6, 8)
}

// M12_8 returns  the MIDI time signature meta message for a 12/8 meter
func M12_8() meta.TimeSig {
	return Meter(12, 8)
}

// M5_8 returns  the MIDI time signature meta message for a 5/8 meter
func M5_8() meta.TimeSig {
	return Meter(5, 8)
}

// M7_8 returns  the MIDI time signature meta message for a 7/8 meter
func M7_8() meta.TimeSig {
	return Meter(7, 8)
}

/*


type TimeSignature struct {
	Numerator   uint8
	Denominator uint8
	// ClocksPerClick           uint8
	// DemiSemiQuaverPerQuarter uint8
}

/*
func NewTimeSignature(num uint8, denom uint8) TimeSignature {
	return TimeSignature{Numerator: num, Denominator: denom}
}



func (m TimeSignature) Raw() []byte {
	// cpcl := m.ClocksPerClick
	// if cpcl == 0 {
	cpcl := byte(8)
	// }

	// dsqpq := m.DemiSemiQuaverPerQuarter
	// if dsqpq == 0 {
	dsqpq := byte(8)
	// }

	var denom = dec2binDenom(m.Denominator)

	return (&metaMessage{
		Typ:  byteTimeSignature,
		Data: []byte{m.Numerator, denom, cpcl, dsqpq},
	}).Bytes()

}

func (m TimeSignature) String() string {
	//return fmt.Sprintf("%T %v/%v clocksperclick %v dsqpq %v", m, m.Numerator, m.Denominator, m.ClocksPerClick, m.DemiSemiQuaverPerQuarter)
	return fmt.Sprintf("%T %v/%v", m, m.Numerator, m.Denominator)
}

func (m TimeSignature) readFrom(rd io.Reader) (Message, error) {
	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 4 {
		err = unexpectedMessageLengthError("TimeSignature expected length 4")
		return nil, err
	}

	// TODO TEST
	var numerator uint8
	numerator, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	var denominator uint8
	denominator, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	var clocksPerClick uint8
	clocksPerClick, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	var demiSemiQuaverPerQuarter uint8
	demiSemiQuaverPerQuarter, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	// TODO: do something with clocksPerClick and demiSemiQuaverPerQuarter
	var _ = clocksPerClick
	var _ = demiSemiQuaverPerQuarter

	return TimeSignature{
		Numerator:   numerator,
		Denominator: 2 << (denominator - 1),
		// ClocksPerClick:           clocksPerClick,
		// DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
	}, nil

}

func (m TimeSignature) meta() {}


*/
