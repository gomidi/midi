package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
	// "github.com/gomidi/midi/smf"
)

/*
http://www.somascape.org/midi/tech/mfile.html

Time Signature

FF 58 04 nn dd cc bb

nn is a byte specifying the numerator of the time signature (as notated).
dd is a byte specifying the denominator of the time signature as a negative power of 2 (i.e. 2 represents a quarter-note, 3 represents an eighth-note, etc).
cc is a byte specifying the number of MIDI clocks between metronome clicks.
bb is a byte specifying the number of notated 32nd-notes in a MIDI quarter-note (24 MIDI Clocks). The usual value for this parameter is 8, though some sequencers allow the user to specify that what MIDI thinks of as a quarter note, should be notated as something else.
Examples

A time signature of 4/4, with a metronome click every 1/4 note, would be encoded :
FF 58 04 04 02 18 08
There are 24 MIDI Clocks per quarter-note, hence cc=24 (0x18).

A time signature of 6/8, with a metronome click every 3rd 1/8 note, would be encoded :
FF 58 04 06 03 24 08
Remember, a 1/4 note is 24 MIDI Clocks, therefore a bar of 6/8 is 72 MIDI Clocks.
Hence 3 1/8 notes is 36 (=0x24) MIDI Clocks.

There should generally be a Time Signature Meta event at the beginning of a track (at time = 0), otherwise a default 4/4 time signature will be assumed. Thereafter they can be used to effect an immediate time signature change at any point within a track.

For a format 1 MIDI file, Time Signature Meta events should only occur within the first MTrk chunk.

*/

// TimeSig sets the time signature according to the SMF spec.
// Denominator isn't a power of 2, but a readable decimal number (uint8).
// If you want an easy way without having to worry about ClocksPerClick
// and DemiSemiQuaverPerQuarter, use the meter subpackage.
type TimeSig struct {
	Numerator                uint8
	Denominator              uint8
	ClocksPerClick           uint8
	DemiSemiQuaverPerQuarter uint8
}

// Raw returns the raw MIDI data
func (m TimeSig) Raw() []byte {
	cpcl := m.ClocksPerClick
	if cpcl == 0 {
		cpcl = byte(8)
	}

	dsqpq := m.DemiSemiQuaverPerQuarter
	if dsqpq == 0 {
		dsqpq = byte(8)
	}

	var denom = dec2binDenom(m.Denominator)

	return (&metaMessage{
		Typ:  byteTimeSignature,
		Data: []byte{m.Numerator, denom, cpcl, dsqpq},
	}).Bytes()

}

// Signature returns the time signature in a readable way
func (m TimeSig) Signature() string {
	return fmt.Sprintf("%v/%v", m.Numerator, m.Denominator)
}

// String represents the time signature MIDI message as a string (for debugging)
func (m TimeSig) String() string {
	return fmt.Sprintf("%T %v/%v clocksperclick %v dsqpq %v", m, m.Numerator, m.Denominator, m.ClocksPerClick, m.DemiSemiQuaverPerQuarter)
	//return fmt.Sprintf("%T %v/%v", m, m.Numerator, m.Denominator)
}

// dec2binDenom converts the decimal denominator to the binary one
// it works, use it!
func dec2binDenom(dec uint8) (bin uint8) {
	if dec <= 1 {
		return 0
	}
	for dec > 2 {
		bin++
		dec = dec >> 1

	}
	return bin + 1
}

func (m TimeSig) readFrom(rd io.Reader) (Message, error) {
	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 4 {
		err = unexpectedMessageLengthError("TimeSignature expected length 4")
		return nil, err
	}

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

	m.DemiSemiQuaverPerQuarter = demiSemiQuaverPerQuarter
	m.ClocksPerClick = clocksPerClick
	m.Numerator = numerator
	m.Denominator = bin2decDenom(denominator)
	return m, nil
	/*
		// TODO: do something with clocksPerClick and demiSemiQuaverPerQuarter
		var _ = clocksPerClick
		var _ = demiSemiQuaverPerQuarter

		return TimeSignature{
			Numerator:   numerator,
			Denominator: 2 << (denominator - 1),
			// ClocksPerClick:           clocksPerClick,
			// DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
		}, nil
	*/
}

func (m TimeSig) meta() {}

// bin2decDenom converts the binary denominator to the decimal
func bin2decDenom(bin uint8) uint8 {
	if bin == 0 {
		return 1
	}
	return 2 << (bin - 1)
}
