package channel

import (
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
)

// ProgramChange represents a MIDI program change message
type ProgramChange struct {
	channel uint8
	program uint8
}

// Program returns the program of the program change message.
func (p ProgramChange) Program() uint8 {
	return p.program
}

// Channel returns the channel of the program change message.
func (m ProgramChange) Channel() uint8 {
	return m.channel
}

// Raw returns the raw bytes of the program change message.
func (c ProgramChange) Raw() []byte {
	return channelMessage1(c.channel, 12, c.program)
}

// String returns human readable information about the program change message.
func (c ProgramChange) String() string {
	return fmt.Sprintf("%T channel %v program %v", c, c.channel, c.program)
}

// set returns a new program change message that is set to the parsed arguments
func (ProgramChange) set(channel uint8, firstArg uint8) setter1 {
	var m ProgramChange
	m.channel = channel
	m.program = midilib.ParseUint7(firstArg)
	return m
}
