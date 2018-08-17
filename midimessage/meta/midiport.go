package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

// MIDIPort represents the deprecated MIDI port message
type MIDIPort uint8

// Number returns the number of the port
func (m MIDIPort) Number() uint8 {
	return uint8(m)
}

// String represents the MIDI port message as a string (for debugging)
func (m MIDIPort) String() string {
	return fmt.Sprintf("%T: %v", m, m.Number())
}

// Raw returns the raw MIDI data
func (m MIDIPort) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteMIDIPort),
		Data: []byte{byte(m)},
	}).Bytes()
}

func (m MIDIPort) meta() {}

func (m MIDIPort) readFrom(rd io.Reader) (Message, error) {

	// Obsolete 'MIDI Port'
	//	we can't ignore it, since it advanced in deltatime

	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 1 {
		return nil, unexpectedMessageLengthError("MIDI Port Message expected length 1")
	}

	var port uint8
	port, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return MIDIPort(port), nil

}
