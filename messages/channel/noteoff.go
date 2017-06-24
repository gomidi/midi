package channel

import (
	"fmt"
	"github.com/gomidi/midi/internal/lib"
)

// NoteOffPedantic is offered as an alternative to NoteOff for getting
// the "real" noteoff message (type 8) and preserving velocity.
type NoteOffPedantic struct {
	NoteOff
	velocity uint8
}

func (n NoteOffPedantic) Velocity() uint8 {
	return n.velocity
}

func (NoteOffPedantic) set(channel uint8, arg1, arg2 uint8) setter2 {
	var m NoteOffPedantic
	m.channel = channel
	m.pitch, m.velocity = lib.ParseTwoUint7(arg1, arg2)
	return m
}

// Raw returns the bytes for the noteoff message.
// Since NoteOff.Raw() returns in fact a noteon message (type 9) with velocity of 0 to allow running status,
// NoteOffPedantic.Raw() is offered as an alternative to make sure a "real" noteoff message (type 8) is returned.
func (n NoteOffPedantic) Raw() []byte {
	return channelMessage2(n.channel, 8, n.pitch, n.velocity)
}

// NoteOff creates a "fake" noteoff message by creating a NoteOn with velocity of 0 (helps for running status).
type NoteOff struct {
	channel uint8
	pitch   uint8
}

func (n NoteOff) Pitch() uint8 {
	return n.pitch
}

// Raw returns the bytes for the noteoff message.
// To allowing running status, here the bytes for a noteon message (type 9) with velocity = 0 are returned.
// If you need a "real" noteoff message, call NoteOffPedantic.Raw()
func (n NoteOff) Raw() []byte {
	return channelMessage2(n.channel, 9, n.pitch, 0)
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
