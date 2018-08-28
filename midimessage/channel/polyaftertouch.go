package channel

import (
	"fmt"

	"github.com/gomidi/midi/internal/midilib"
)

// PolyAftertouch represents a MIDI polyphonic aftertouch message (aka "key pressure")
type PolyAftertouch struct {
	channel  uint8
	key      uint8
	pressure uint8
}

// Key returns the key of the polyphonic aftertouch message
func (p PolyAftertouch) Key() uint8 {
	return p.key
}

// Pressure returns the pressure of the polyphonic aftertouch message
func (p PolyAftertouch) Pressure() uint8 {
	return p.pressure
}

// Channel returns the channel of the polyphonic aftertouch message
func (p PolyAftertouch) Channel() uint8 {
	return p.channel
}

// String returns human readable information about the polyphonic aftertouch message.
func (p PolyAftertouch) String() string {
	return fmt.Sprintf("%T channel %v key %v pressure %v", p, p.Channel(), p.Key(), p.Pressure())
}

// Raw returns the raw bytes of the polyphonic aftertouch message.
func (p PolyAftertouch) Raw() []byte {
	return channelMessage2(p.channel, 10, p.key, p.pressure)
}

// set returns a new polyphonic aftertouch message that is set to the parsed arguments
func (PolyAftertouch) set(channel uint8, arg1, arg2 uint8) setter2 {
	var m PolyAftertouch
	m.channel = channel
	m.key, m.pressure = midilib.ParseTwoUint7(arg1, arg2)
	return m
}
