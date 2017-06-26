package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
	"github.com/gomidi/midi/smf"
)

type MIDIPort uint8

func (m MIDIPort) Number() uint8 {
	return uint8(m)
}

func (m MIDIPort) String() string {
	return fmt.Sprintf("%T: %v", m, uint8(m))
}

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
		return nil, smf.UnexpectedMessageLengthError("MIDI Port Message expected length 1")
	}

	var port uint8
	port, err = midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return MIDIPort(port), nil

}
