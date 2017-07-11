package sysex

import (
	"fmt"
)

const (
	byteSysExStart = byte(0xF0)
	byteSysExEnd   = byte(0xF7)
)

// Escape is a sysex escape sequence with a prefixed 0xF7
// it may only used within SMF files (not for live MIDI)
type Escape []byte

// Data returns the escaped data
func (m Escape) Data() []byte {
	return []byte(m)
}

func (m Escape) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

// Raw() returns the data with the escape prefix 0xF7
func (m Escape) Raw() []byte {
	var b = []byte{0xF7}
	b = append(b, []byte(m)...)
	return b
}

func (m Escape) Len() int {
	return len(m)
}

func (c Escape) sysex() {}

// SysExStart is an incomplete sysex that is the start of several sysexes (casio style)
// i.e. beginning with 0xF0 but no 0xF7 at the end
// when used within a SMF file, the first byte (0xF0) must be followed by a length
// when used live, no messages apart from realtime messages may be send before the
// rest of the sysex was send
type Start []byte

// Data returns the inner sysex data
func (m Start) Data() []byte {
	return []byte(m)
}

// Raw() returns the data with the prefix 0xF0
func (m Start) Raw() []byte {
	var b = []byte{0xF0}
	b = append(b, []byte(m)...)
	return b
}

func (m Start) Len() int {
	return len(m)
}

func (c Start) sysex() {}

func (m Start) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

// Continue is an incomplete sysex that is following Start or SysExContinue but not ending it.
// It starts with an 0xF7 but does not end with a 0xF7
// when used within a SMF file, the first byte (0xF7) must be followed by a length
// when used live, no messages apart from realtime messages may be send before the
// rest of the sysex was send
type Continue []byte

// Data returns the inner sysex data
func (m Continue) Data() []byte {
	return []byte(m)
}

func (c Continue) sysex() {}

func (m Continue) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

// Raw() returns the data with the prefix 0xF7
func (m Continue) Raw() []byte {
	var b = []byte{0xF7}
	b = append(b, []byte(m)...)
	return b
}

func (m Continue) Len() int {
	return len(m)
}

// End is an incomplete sysex that is continuing Start or Continue and ending it.
// It starts with an 0xF7 and ends with a 0xF7
// when used within a SMF file, the first byte (0xF7) must be followed by a length (including the last F7 but excluding the first)
// when used live, no messages apart from realtime messages may be send in between the preceeding Start or Continue and
// this one
type End []byte

// Data returns the inner sysex data
func (m End) Data() []byte {
	return []byte(m)
}

func (m End) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

func (m End) Len() int {
	return len(m)
}

func (c End) sysex() {}

// Raw() returns the data with the prefix 0xF7 and the postfix 0xF7
func (m End) Raw() []byte {
	var b = []byte{0xF7}
	b = append(b, []byte(m)...)
	b = append(b, 0xF7)
	return b
}

// Message is a System Exclusive Message
type Message interface {
	String() string
	Raw() []byte
	Len() int
	Data() []byte
	sysex()
}

var _ Message = SysEx([]byte{})
var _ Message = Escape([]byte{})
var _ Message = Start([]byte{})
var _ Message = End([]byte{})
var _ Message = Continue([]byte{})

// SysEx is a sysex that is complete (i.e. starting with 0xF0 and ending with 0xF7
// it may be used within SMF files and with live MIDI.
// However when used within a SMF file, the first byte (0xF0) must be followed by a length
// before the rest comes (including the 0xF7)
type SysEx []byte

// Data returns the inner sysex data
func (m SysEx) Data() []byte {
	return []byte(m)
}

func (m SysEx) sysex() {}

func (m SysEx) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

func (m SysEx) Len() int {
	return len(m)
}

// Raw returns the sysex data with the prefixed 0xF0 and
// the postfix 0xF7
func (m SysEx) Raw() []byte {
	var b = []byte{0xF0}
	b = append(b, []byte(m)...)
	b = append(b, 0xF7)
	return b
}
