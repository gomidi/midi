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
func (p ProgramChange) Channel() uint8 {
	return p.channel
}

// Raw returns the raw bytes of the program change message.
func (p ProgramChange) Raw() []byte {
	return channelMessage1(p.channel, 12, p.program)
}

// String returns human readable information about the program change message.
func (p ProgramChange) String() string {
	return fmt.Sprintf("%T channel %v program %v", p, p.Channel(), p.Program())
}

// set returns a new program change message that is set to the parsed arguments
func (ProgramChange) set(channel uint8, firstArg uint8) setter1 {
	var p ProgramChange
	p.channel = channel
	p.program = midilib.ParseUint7(firstArg)
	return p
}
