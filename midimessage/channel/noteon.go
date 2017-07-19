package channel

import (
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
)

// NoteOn represents a note-on message
type NoteOn struct {
	channel  uint8
	key      uint8
	velocity uint8
}

// Key returns the key of the note-on message
func (n NoteOn) Key() uint8 {
	return n.key
}

// Velocity returns the velocity of the note-on message
func (n NoteOn) Velocity() uint8 {
	return n.velocity
}

// Channel returns the channel of the note-on message
func (n NoteOn) Channel() uint8 {
	return n.channel
}

// Raw returns the bytes for the noteon message.
func (n NoteOn) Raw() []byte {
	return channelMessage2(n.channel, 9, n.key, n.velocity)
}

// String returns human readable information about the note-on message.
func (n NoteOn) String() string {
	return fmt.Sprintf("%T channel %v key %v velocity %v", n, n.Channel(), n.Key(), n.Velocity())
}

// set returns a new note-on message that is set to the parsed arguments
func (NoteOn) set(channel, arg1, arg2 uint8) setter2 {
	var m NoteOn
	m.channel = channel
	m.key, m.velocity = midilib.ParseTwoUint7(arg1, arg2)
	return m
}
