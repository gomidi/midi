package meta

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

// SequenceNumber represents the sequence number MIDI meta message
type SequenceNumber uint16

// Number returns the number of the sequence
func (s SequenceNumber) Number() uint16 {
	return uint16(s)
}

// String represents the MIDI sequence name message as a string (for debugging)
func (s SequenceNumber) String() string {
	return fmt.Sprintf("%T: %v", s, s.Number())
}

// Raw returns the raw bytes for the message
func (s SequenceNumber) Raw() []byte {
	var bf bytes.Buffer
	binary.Write(&bf, binary.BigEndian, s.Number())
	return (&metaMessage{
		Typ:  byteSequenceNumber,
		Data: bf.Bytes(),
	}).Bytes()
}

func (s SequenceNumber) readFrom(rd io.Reader) (Message, error) {
	length, err := midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	// Zero length sequences allowed according to http://home.roadrunner.com/~jgglatt/tech/midifile/seq.htm
	if length == 0 {
		return SequenceNumber(0), nil
	}

	// Otherwise length will be 2 to hold the uint16.
	var sequenceNumber uint16
	sequenceNumber, err = midilib.ReadUint16(rd)

	if err != nil {
		return nil, err
	}

	return SequenceNumber(sequenceNumber), nil
}

func (s SequenceNumber) meta() {}
