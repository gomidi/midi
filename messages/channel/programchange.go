package channel

import "fmt"

type ProgramChange struct {
	channel uint8
	program uint8
}

func (p ProgramChange) Program() uint8 {
	return p.program
}

func (m ProgramChange) Channel() uint8 {
	return m.channel
}

func (c ProgramChange) Raw() []byte {
	return channelMessage1(c.channel, 12, c.program)
}

func (c ProgramChange) String() string {
	return fmt.Sprintf("%T channel %v program %v", c, c.channel, c.program)
}

func (ProgramChange) set(channel uint8, firstArg uint8) setter1 {
	var m ProgramChange
	m.channel = channel
	m.program = parseUint7(firstArg)
	return m
}
