package meta

import (
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
	"io"
)

/*
FF 7F len data Sequencer Specific Meta-Event
Special requirements for particular sequencers may use this event type: the first byte or bytes of data is a
manufacturer ID (these are one byte, or if the first byte is 00, three bytes). As with MIDI System Exclusive,
manufacturers who define something using this meta-event should publish it so that others may be used by a
sequencer which elects to use this as its only file format; sequencers with their established feature-specific
formats should probably stick to the standard features when using this format.
*/

type SequencerSpecific []byte

func (s SequencerSpecific) Data() []byte {
	return []byte(s)
}

func (s SequencerSpecific) Raw() []byte {
	return (&metaMessage{
		Typ:  byteSequencerSpecific,
		Data: s.Data(),
	}).Bytes()
}

func (s SequencerSpecific) Len() int {
	return len(s)
}

func (s SequencerSpecific) String() string {
	return fmt.Sprintf("%T len %v", s, s.Len())
}

func (s SequencerSpecific) meta() {

}

func (s SequencerSpecific) readFrom(rd io.Reader) (Message, error) {
	length, err := midilib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	bt, err := midilib.ReadNBytes(int(length), rd)

	if err != nil {
		return nil, err
	}

	return SequencerSpecific(bt), nil
}
