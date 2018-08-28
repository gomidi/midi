package meta

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

/*
FF 7F len data Sequencer Specific Meta-Event
Special requirements for particular sequencers may use this event type: the first byte or bytes of data is a
manufacturer ID (these are one byte, or if the first byte is 00, three bytes). As with MIDI System Exclusive,
manufacturers who define something using this meta-event should publish it so that others may be used by a
sequencer which elects to use this as its only file format; sequencers with their established feature-specific
formats should probably stick to the standard features when using this format.
*/

// SequencerData is a sequencer specific meta message
type SequencerData []byte

// Data returns the sequencer specific data
func (s SequencerData) Data() []byte {
	return []byte(s)
}

// Raw returns the raw MIDI data
func (s SequencerData) Raw() []byte {
	return (&metaMessage{
		Typ:  byteSequencerSpecific,
		Data: s.Data(),
	}).Bytes()
}

// Len returns the length of the sequencer specific data
func (s SequencerData) Len() int {
	return len(s)
}

// String represents the sequencer spefici MIDI message as a string (for debugging)
func (s SequencerData) String() string {
	return fmt.Sprintf("%T len %v", s, s.Len())
}

func (s SequencerData) meta() {

}

func (s SequencerData) readFrom(rd io.Reader) (Message, error) {
	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	bt, err := midilib.ReadNBytes(int(length), rd)

	if err != nil {
		return nil, err
	}

	return SequencerData(bt), nil
}
