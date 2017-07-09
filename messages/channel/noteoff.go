package channel

import "fmt"

// NoteOffPedantic is offered as an alternative to NoteOff for getting
// the "real" noteoff message (type 8) and preserving velocity.
type NoteOffPedantic struct {
	NoteOff
	velocity uint8
}

func (n NoteOffPedantic) IsLiveMessage() {

}

func (n NoteOffPedantic) Velocity() uint8 {
	return n.velocity
}

func (NoteOffPedantic) set(channel uint8, arg1, arg2 uint8) setter2 {
	var m NoteOffPedantic
	m.channel = channel
	m.key, m.velocity = parseTwoUint7(arg1, arg2)
	return m
}

// Raw returns the bytes for the noteoff message.
// Since NoteOff.Raw() returns in fact a noteon message (type 9) with velocity of 0 to allow running status,
// NoteOffPedantic.Raw() is offered as an alternative to make sure a "real" noteoff message (type 8) is returned.
func (n NoteOffPedantic) Raw() []byte {
	return channelMessage2(n.channel, 8, n.key, n.velocity)
}

func (m NoteOffPedantic) String() string {
	return fmt.Sprintf("%T channel %v key %v vel %v", m, m.channel, m.key, m.velocity)
}

// NoteOff creates a "fake" noteoff message by creating a NoteOn with velocity of 0 (helps for running status).
type NoteOff struct {
	channel uint8
	key     uint8
}

func (n NoteOff) IsLiveMessage() {

}

func (n NoteOff) Key() uint8 {
	return n.key
}

// Raw returns the bytes for the noteoff message.
// To allowing running status, here the bytes for a noteon message (type 9) with velocity = 0 are returned.
// If you need a "real" noteoff message, call NoteOffPedantic.Raw()
func (n NoteOff) Raw() []byte {
	return channelMessage2(n.channel, 9, n.key, 0)
}

func (n NoteOff) Channel() uint8 {
	return n.channel
}

func (m NoteOff) String() string {
	return fmt.Sprintf("%T channel %v key %v", m, m.channel, m.key)
}

func (NoteOff) set(channel uint8, arg1, arg2 uint8) setter2 {
	var m NoteOff
	m.channel = channel
	m.key, _ = parseTwoUint7(arg1, arg2)
	return m
}
