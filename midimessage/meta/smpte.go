package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

/*
FF 54 05 hr mn se fr ff SMPTE Offset
This event, if present, designates the SMPTE time at which the track chunk is supposed to start. It should be
present at the beginning of the track, that is, before any nonzero delta-times, and before any transmittable
MIDI events. the hour must be encoded with the SMPTE format, just as it is in MIDI Time Code. In a format
1 file, the SMPTE Offset must be stored with the tempo map, and has no meaning in any of the other tracks.
The ff field contains fractional frames, in 100ths of a frame, even in SMPTE-based tracks which specify a
different frame subdivision for delta-times.

SMPTE timing is referenced from an
absolute "time of day". On the other hand, MIDI Clocks and Song Position
Pointer are based upon musical beats from the start of a song, played at a
specific Tempo. For many (non-musical) cues, it's easier for humans to
reference time in some absolute way rather than based upon musical beats at
a certain tempo.
*/

// SMPTE represents a smpte offset MIDI meta message
type SMPTE struct {
	Hour            byte
	Minute          byte
	Second          byte
	Frame           byte
	FractionalFrame byte
}

// Raw returns the raw bytes for the message
func (s SMPTE) Raw() []byte {
	return (&metaMessage{
		Typ:  byteSMPTEOffset,
		Data: []byte{s.Hour, s.Minute, s.Second, s.Frame, s.FractionalFrame},
	}).Bytes()
}

// String represents the smpte offset MIDI message as a string (for debugging)
func (s SMPTE) String() string {
	return fmt.Sprintf("%T %v:%v:%v %v.%0d", s, s.Hour, s.Minute, s.Second, s.Frame, s.FractionalFrame)
}

func (s SMPTE) readFrom(rd io.Reader) (Message, error) {
	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 5 {
		err = unexpectedMessageLengthError("SMPTEOffset expected length 5")
		return nil, err
	}

	bt, err := midilib.ReadNBytes(5, rd)

	if err != nil {
		return nil, err
	}

	s.Hour = bt[0]
	s.Minute = bt[1]
	s.Second = bt[2]
	s.Frame = bt[3]
	s.FractionalFrame = bt[4]

	return s, nil
}

func (s SMPTE) meta() {

}
