package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

// Channel represents the deprecated MIDI channel meta message
type Channel uint8

// Number returns the number of the MIDI channel (starting with 0)
func (m Channel) Number() uint8 {
	return uint8(m)
}

// String represents the MIDIChannel message as a string (for debugging)
func (m Channel) String() string {
	return fmt.Sprintf("%T: %v", m, m.Number())
}

// Raw returns the raw bytes for the message
func (m Channel) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteMIDIChannel),
		Data: []byte{byte(m)},
	}).Bytes()
}

func (m Channel) meta() {}

func (m Channel) readFrom(rd io.Reader) (Message, error) {

	// Obsolete 'MIDI Channel'
	//	we can't ignore it, since it advanced in deltatime

	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 1 {
		return nil, unexpectedMessageLengthError("Midi Channel Message expected length 1")
	}

	var ch uint8
	ch, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return Channel(ch), nil

}
