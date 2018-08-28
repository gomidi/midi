package meta

import (
	"io"
)

const (
	// End of track
	// the handler is supposed to keep track of the current track

	byteEndOfTrack        = byte(0x2F)
	byteSequenceNumber    = byte(0x00)
	byteText              = byte(0x01)
	byteCopyright         = byte(0x02)
	byteSequence          = byte(0x03)
	byteTrack             = byte(0x04)
	byteLyric             = byte(0x05)
	byteMarker            = byte(0x06)
	byteCuepoint          = byte(0x07)
	byteMIDIChannel       = byte(0x20)
	byteDevicePort        = byte(0x9)
	byteMIDIPort          = byte(0x21)
	byteTempo             = byte(0x51)
	byteTimeSignature     = byte(0x58)
	byteKeySignature      = byte(0x59)
	byteSequencerSpecific = byte(0x7F)
	byteSMPTEOffset       = byte(0x54)
	byteProgramName       = byte(0x8)
)

var metaMessages = map[byte]Message{
	byteEndOfTrack:        EndOfTrack,
	byteSequenceNumber:    SequenceNo(0),
	byteText:              Text(""),
	byteCopyright:         Copyright(""),
	byteSequence:          Sequence(""),
	byteTrack:             Track(""),
	byteLyric:             Lyric(""),
	byteMarker:            Marker(""),
	byteCuepoint:          Cuepoint(""),
	byteMIDIChannel:       Channel(0),
	byteDevicePort:        Device(""),
	byteMIDIPort:          Port(0),
	byteTempo:             Tempo(0),
	byteTimeSignature:     TimeSig{},
	byteKeySignature:      Key{},
	byteSMPTEOffset:       SMPTE{},
	byteSequencerSpecific: SequencerData(nil),
	byteProgramName:       Program(""),
}

// Reader reads a Meta Message
type Reader interface {
	// Read reads a single Meta Message.
	// It may just be called once per Reader. A second call returns io.EOF
	Read() (Message, error)
}

// NewReader returns a reader that can read a single Meta Message
// Read may just be called once per Reader. A second call returns io.EOF
func NewReader(input io.Reader, typ byte) Reader {
	return &reader{input, typ, false}
}

type reader struct {
	input io.Reader
	typ   byte
	done  bool
}

// Read may just be called once per Reader. A second call returns io.EOF
func (r *reader) Read() (Message, error) {
	if r.done {
		return nil, io.EOF
	}

	m := metaMessages[r.typ]
	if m == nil {
		m = Undefined{Typ: r.typ}
	}

	return m.readFrom(r.input)
}
