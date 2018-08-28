package syscommon

import (
	"io"
)

// Reader read a syscommon
type Reader interface {
	// Read reads a single syscommon message.
	// It may just be called once per Reader. A second call returns io.EOF
	Read() (Message, error)
}

// NewReader returns a reader that can read a single syscommon message
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
func (r *reader) Read() (msg Message, err error) {
	if r.done {
		return nil, io.EOF
	}
	msg = dispatch(r.typ)

	// unknown or undefined message
	if msg == nil {
		return
	}
	return msg.readFrom(r.input)
}

func dispatch(b byte) Message {
	return systemMessages[b]
}

var systemMessages = map[byte]Message{
	byteMIDITimingCodeMessage:  MTC(0),
	byteSysSongPositionPointer: SPP(0),
	byteSysSongSelect:          SongSelect(0),
	byteSysTuneRequest:         Tune,
	//	0xF4:                       Undefined4(0), // unused (ignore them)
	//	0xF5:                       Undefined5(0), // unused (ignore them)
}

const (
	byteMIDITimingCodeMessage  = byte(0xF1)
	byteSysSongPositionPointer = byte(0xF2)
	byteSysSongSelect          = byte(0xF3)
	byteSysTuneRequest         = byte(0xF6)
)
