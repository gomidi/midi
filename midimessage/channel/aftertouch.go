package channel

import (
	"fmt"

	"github.com/gomidi/midi/internal/midilib"
)

// Aftertouch represents a MIDI aftertouch message (aka "channel pressure")
type Aftertouch struct {
	channel  uint8
	pressure uint8
}

// Pressure returns the pressure of the aftertouch message.
func (a Aftertouch) Pressure() uint8 {
	return a.pressure
}

// Channel returns the channel of the aftertouch message.
func (a Aftertouch) Channel() uint8 {
	return a.channel
}

// Raw returns the raw bytes of the aftertouch message.
func (a Aftertouch) Raw() []byte {
	return channelMessage1(a.channel, 13, a.pressure)
}

// String returns human readable information about the aftertouch message.
func (a Aftertouch) String() string {
	return fmt.Sprintf("%T channel %v pressure %v", a, a.Channel(), a.Pressure())
}

// set returns a new aftertouch message that is set to the parsed arguments
func (Aftertouch) set(channel uint8, firstArg uint8) setter1 {
	var m Aftertouch
	m.channel = channel
	m.pressure = midilib.ParseUint7(firstArg)
	return m
}
