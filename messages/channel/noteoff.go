package channel

import (
	"fmt"
	"github.com/gomidi/midi/internal/lib"
)

type NoteOff struct {
	channel uint8
	pitch   uint8
}

func (n NoteOff) Pitch() uint8 {
	return n.pitch
}

func (n NoteOff) Raw() []byte {
	return channelMessage2(n.channel, 8, n.pitch, 0)
}

func (n NoteOff) Channel() uint8 {
	return n.channel
}

func (m NoteOff) String() string {
	return fmt.Sprintf("%T channel %v pitch %v", m, m.channel, m.pitch)
}

func (NoteOff) set(channel uint8, arg1, arg2 uint8) setter2 {
	var m NoteOff
	m.channel = channel
	m.pitch, _ = lib.ParseTwoUint7(arg1, arg2)
	return m
}
