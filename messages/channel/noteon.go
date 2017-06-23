package channel

import (
	"fmt"
	"lib"
)

type NoteOn struct {
	channel  uint8
	pitch    uint8
	velocity uint8
}

func (n NoteOn) Pitch() uint8 {
	return n.pitch
}

func (n NoteOn) Velocity() uint8 {
	return n.velocity
}

func (n NoteOn) Channel() uint8 {
	return n.channel
}

func (n NoteOn) Raw() []byte {
	return channelMessage2(n.channel, 9, n.pitch, n.velocity)
}

func (n NoteOn) String() string {
	return fmt.Sprintf("%T channel %v pitch %v vel %v", n, n.channel, n.pitch, n.velocity)
}

func (NoteOn) set(channel, arg1, arg2 uint8) setter2 {
	var m NoteOn
	m.channel = channel
	m.pitch, m.velocity = lib.ParseTwoUint7(arg1, arg2)
	return m
}
