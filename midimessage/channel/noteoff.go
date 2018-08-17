package channel

import (
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
)

// NoteOffVelocity is offered as an alternative to NoteOff for
// a "real" noteoff message (type 8) that has velocity.
type NoteOffVelocity struct {
	NoteOff
	velocity uint8
}

// Velocity returns the velocity of the note-off message
func (n NoteOffVelocity) Velocity() uint8 {
	return n.velocity
}

// set returns a new note-off message with velocity that is set to the parsed arguments
func (NoteOffVelocity) set(channel uint8, arg1, arg2 uint8) setter2 {
	var n NoteOffVelocity
	n.channel = channel
	n.key, n.velocity = midilib.ParseTwoUint7(arg1, arg2)
	return n
}

// Raw returns the bytes for the noteoff message.
// Since NoteOff.Raw() returns in fact a noteon message (type 9) with velocity of 0 to allow running status,
// NoteOffPedantic.Raw() is offered as an alternative to make sure a "real" noteoff message (type 8) is returned.
func (n NoteOffVelocity) Raw() []byte {
	return channelMessage2(n.channel, 8, n.key, n.velocity)
}

// String returns human readable information about the note-off message that includes velocity.
func (n NoteOffVelocity) String() string {
	return fmt.Sprintf("%T channel %v key %v velocity %v", n, n.Channel(), n.Key(), n.Velocity())
}

// NoteOff represents a note-off message by a note-on message with velocity of 0 (helps for running status).
// This is the normal way to go. If you need the velocity of a note-off message, use NoteOffVelocity.
type NoteOff struct {
	channel uint8
	key     uint8
}

// Key returns the key of the note off message
func (n NoteOff) Key() uint8 {
	return n.key
}

// Raw returns the bytes for the noteoff message.
// To allowing running status, here the bytes for a noteon message (type 9) with velocity = 0 are returned.
// If you need a "real" noteoff message, call NoteOffPedantic.Raw()
func (n NoteOff) Raw() []byte {
	return channelMessage2(n.channel, 9, n.key, 0)
}

// Channel returns the channel of the note-off message
func (n NoteOff) Channel() uint8 {
	return n.channel
}

// String returns human readable information about the note-off message.
func (n NoteOff) String() string {
	return fmt.Sprintf("%T channel %v key %v", n, n.Channel(), n.Key())
}

// set returns a new note-off message that is set to the parsed arguments
func (NoteOff) set(channel uint8, arg1, arg2 uint8) setter2 {
	var n NoteOff
	n.channel = channel
	n.key, _ = midilib.ParseTwoUint7(arg1, arg2)
	return n
}
