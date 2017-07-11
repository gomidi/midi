package channel

import "fmt"

// AfterTouch represents a MIDI aftertouch message (aka "channel pressure")
type AfterTouch struct {
	channel  uint8
	pressure uint8
}

// Pressure returns the pressure of the aftertouch message.
func (a AfterTouch) Pressure() uint8 {
	return a.pressure
}

// Channel returns the channel of the aftertouch message.
func (a AfterTouch) Channel() uint8 {
	return a.channel
}

// Raw returns the raw bytes of the aftertouch message.
func (a AfterTouch) Raw() []byte {
	return channelMessage1(a.channel, 13, a.pressure)
}

// String returns human readable information about the aftertouch message.
func (a AfterTouch) String() string {
	return fmt.Sprintf("%T (\"ChannelPressure\") channel %v pressure %v", a, a.channel, a.pressure)
}

// set returns a new aftertouch message that is set to the parsed arguments
func (AfterTouch) set(channel uint8, firstArg uint8) setter1 {
	var m AfterTouch
	m.channel = channel
	m.pressure = parseUint7(firstArg)
	return m
}
