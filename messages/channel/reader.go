package channel

import (
	"io"
	"lib"
)

const (
	byteProgramChange         = 0xC
	byteChannelPressure       = 0xD
	byteNoteOff               = 0x8
	byteNoteOn                = 0x9
	bytePolyphonicKeyPressure = 0xA
	byteControlChange         = 0xB
	bytePitchWheel            = 0xE
	byteSystemCommon          = 0xF
	byteMetaEvents            = 0xF
)

// Reader read a channel message
type Reader interface {
	// Read reads a single channel message.
	// It may just be called once per Reader. A second call returns io.EOF
	Read() (Message, error)
}

// NewReader returns a reader that can read a single channel message
// Read may just be called once per Reader. A second call returns io.EOF
func NewReader(input io.Reader, status byte) Reader {
	return &reader{input, status, false}
}

type reader struct {
	input  io.Reader
	status byte
	done   bool
}

// Read may just be called once per Reader. A second call returns io.EOF
func (r *reader) Read() (msg Message, err error) {
	if r.done {
		return nil, io.EOF
	}
	var typ, channel, arg1 uint8

	typ, channel = lib.ParseStatus(r.status)

	arg1, err = lib.ReadByte(r.input)
	r.done = true

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
		arg2, err = lib.ReadByte(r.input)

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
		msg = AfterTouch{}
	default:
		// unsupported
		return nil
	}

	msg = msg.set(channel, arg)
	return
}

func (r *reader) getMsg2(typ uint8, channel uint8, arg1 uint8, arg2 uint8) (msg setter2) {

	switch typ {
	case byteNoteOff:
		msg = NoteOff{}
	case byteNoteOn:
		msg = NoteOn{}
	case bytePolyphonicKeyPressure:
		msg = PolyphonicAfterTouch{}
	case byteControlChange:
		msg = ControlChange{}
	case bytePitchWheel:
		msg = PitchWheel{}
	default:
		// unsupported
		return nil
	}

	msg = msg.set(channel, arg1, arg2)

	// handle noteOn messages with velocity of 0 as note offs
	if noteOn, is := msg.(NoteOn); is && noteOn.velocity == 0 {
		msg = NoteOff{}
		msg = msg.set(channel, arg1, 0)
	}
	return
}
