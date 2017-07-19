package meta

import (
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
	"io"
)

type MIDIChannel uint8

func (m MIDIChannel) Number() uint8 {
	return uint8(m)
}

func (m MIDIChannel) String() string {
	return fmt.Sprintf("%T: %v", m, m.Number())
}

func (m MIDIChannel) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteMIDIChannel),
		Data: []byte{byte(m)},
	}).Bytes()
}

func (m MIDIChannel) meta() {}

func (m MIDIChannel) readFrom(rd io.Reader) (Message, error) {

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

	return MIDIChannel(ch), nil

}
