package syscommon

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

/*
MTC Quarter Frame

These are the MTC (i.e. SMPTE based) equivalent of the F8 Timing Clock messages,
though offer much higher resolution, as they are sent at a rate of 96 to 120 times
a second (depending on the SMPTE frame rate). Each Quarter Frame message provides
partial timecode information, 8 sequential messages being required to fully
describe a timecode instant. The reconstituted timecode refers to when the first
partial was received. The most significant nibble of the data byte indicates the
partial (aka Message Type).

Partial	Data byte	Usage
1	0000 bcde	Frame number LSBs 	abcde = Frame number (0 to frameRate-1)
2	0001 000a	Frame number MSB
3	0010 cdef	Seconds LSBs 	abcdef = Seconds (0-59)
4	0011 00ab	Seconds MSBs
5	0100 cdef	Minutes LSBs 	abcdef = Minutes (0-59)
6	0101 00ab	Minutes MSBs
7	0110 defg	Hours LSBs 	ab = Frame Rate (00 = 24, 01 = 25, 10 = 30drop, 11 = 30nondrop)
cdefg = Hours (0-23)
8	0111 0abc	Frame Rate, and Hours MSB
*/

// MTC represents a MIDI timing code message (quarter frame)
type MTC uint8

// String represents the MIDI timing code message as a string (for debugging)
func (m MTC) String() string {
	return fmt.Sprintf("%T: %v", m, m.QuarterFrame())
}

// Raw returns the raw bytes for the message
func (m MTC) Raw() []byte {
	// TODO check - it is a guess
	return []byte{byte(0xF1), byte(m)}
}

// QuarterFrame returns the quarter frame
func (m MTC) QuarterFrame() uint8 {
	return uint8(m)
}

func (m MTC) readFrom(rd io.Reader) (Message, error) {
	b, err := midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return MTC(b), nil
}

func (m MTC) sysCommon() {}
