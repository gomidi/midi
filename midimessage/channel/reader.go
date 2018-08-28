package channel

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

const (
	byteProgramChange         = 0xC
	byteChannelPressure       = 0xD
	byteNoteOff               = 0x8
	byteNoteOn                = 0x9
	bytePolyphonicKeyPressure = 0xA
	byteControlChange         = 0xB
	bytePitchWheel            = 0xE
)

// Reader read a channel message
type Reader interface {
	// Read reads a single channel message.
	// It may just be called once per Reader. A second call returns io.EOF
	Read(status, arg1 byte) (Message, error)
}

// ReaderOption is an option for the channel reader.
type ReaderOption func(*reader)

// ReadNoteOffVelocity lets the reader differentiate between "fake" noteoff messages
// (which are in fact noteon messages (typ 9) with velocity of 0) and "real" noteoff messages (typ 8)
// with own velocity.
// The former are returned as NoteOffVelocity messages and keep the given velocity, the later
// are returned as NoteOff messages without velocity. That means in order to get all noteoff messages,
// there must be checks for NoteOff and NoteOffVelocity (if this option is set).
// If this option is not set, both kinds are returned as NoteOff (default).
func ReadNoteOffVelocity() ReaderOption {
	return func(rd *reader) {
		rd.readNoteOffPedantic = true
	}
}

// NewReader returns a reader
func NewReader(input io.Reader, options ...ReaderOption) Reader {
	rd := &reader{input, false}

	for _, opt := range options {
		opt(rd)
	}

	return rd
}

type reader struct {
	input               io.Reader
	readNoteOffPedantic bool
}

// Read reads a channel message
func (r *reader) Read(status byte, arg1 byte) (msg Message, err error) {
	typ, channel := midilib.ParseStatus(status)

	// fmt.Printf("typ: %v channel: %v\n", typ, channel)

	// fmt.Printf("arg1: %v, err: %v\n", arg1, err)

	if err != nil {
		return
	}

	switch typ {

	// one argument only
	case byteProgramChange, byteChannelPressure:
		msg = r.getMsg1(typ, channel, arg1)

	// two Arguments needed
	default:
		var arg2 byte
		arg2, err = midilib.ReadByte(r.input)

		if err != nil {
			return
		}
		msg = r.getMsg2(typ, channel, arg1, arg2)
	}
	return
}

func (r *reader) getMsg1(typ uint8, channel uint8, arg uint8) (msg setter1) {
	switch typ {
	case byteProgramChange:
		msg = ProgramChange{}
	case byteChannelPressure:
		msg = Aftertouch{}
	default:
		panic(fmt.Sprintf("must not happen (typ % X is not an channel message with one argument)", typ))
	}

	msg = msg.set(channel, arg)
	return
}

func (r *reader) getMsg2(typ uint8, channel uint8, arg1 uint8, arg2 uint8) (msg setter2) {

	switch typ {
	case byteNoteOff:
		if r.readNoteOffPedantic {
			msg = NoteOffVelocity{}
		} else {
			msg = NoteOff{}
		}
	case byteNoteOn:
		msg = NoteOn{}
	case bytePolyphonicKeyPressure:
		msg = PolyAftertouch{}
	case byteControlChange:
		msg = ControlChange{}
	case bytePitchWheel:
		msg = Pitchbend{}
	default:
		panic(fmt.Sprintf("must not happen (typ % X is not an channel message with two arguments)", typ))
	}

	msg = msg.set(channel, arg1, arg2)

	// handle noteOn messages with velocity of 0 as note offs
	if noteOn, is := msg.(NoteOn); is && noteOn.velocity == 0 {
		msg = (NoteOff{}).set(channel, arg1, 0)
	}
	return
}
