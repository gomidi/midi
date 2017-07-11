package channel

import "fmt"

// PolyphonicAfterTouch represents a MIDI polyphonic aftertouch message (aka "key pressure")
type PolyphonicAfterTouch struct {
	channel  uint8
	key      uint8
	pressure uint8
}

// Key returns the key of the polyphonic aftertouch message
func (p PolyphonicAfterTouch) Key() uint8 {
	return p.key
}

// Pressure returns the pressure of the polyphonic aftertouch message
func (p PolyphonicAfterTouch) Pressure() uint8 {
	return p.pressure
}

// Channel returns the channel of the polyphonic aftertouch message
func (p PolyphonicAfterTouch) Channel() uint8 {
	return p.channel
}

// String returns human readable information about the polyphonic aftertouch message.
func (p PolyphonicAfterTouch) String() string {
	return fmt.Sprintf("%T (\"KeyPressure\") channel %v key %v pressure %v", p, p.channel, p.key, p.pressure)
}

// Raw returns the raw bytes of the polyphonic aftertouch message.
func (p PolyphonicAfterTouch) Raw() []byte {
	return channelMessage2(p.channel, 10, p.key, p.pressure)
}

// set returns a new polyphonic aftertouch message that is set to the parsed arguments
func (PolyphonicAfterTouch) set(channel uint8, arg1, arg2 uint8) setter2 {
	var m PolyphonicAfterTouch
	m.channel = channel
	m.key, m.pressure = parseTwoUint7(arg1, arg2)
	return m
}
