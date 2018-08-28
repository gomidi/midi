package meta

import (
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

type endOfTrack bool

const (
	// EndOfTrack represents the end of track MIDI message. It must be written at the end of a track
	EndOfTrack = endOfTrack(true)
)

// String represents the end of track message as a string (for debugging)
func (m endOfTrack) String() string {
	return "meta.EndOfTrack"
}

// Raw returns the raw MIDI data
func (m endOfTrack) Raw() []byte {
	return (&metaMessage{
		Typ: byte(byteEndOfTrack),
	}).Bytes()
}

func (m endOfTrack) meta() {}

func (m endOfTrack) readFrom(rd io.Reader) (Message, error) {

	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 0 {
		err = unexpectedMessageLengthError("EndOfTrack expected length 0")
		return nil, err
	}

	return m, nil
}
