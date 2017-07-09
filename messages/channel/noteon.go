package channel

import "fmt"

type NoteOn struct {
	channel  uint8
	key      uint8
	velocity uint8
}

func (n NoteOn) Key() uint8 {
	return n.key
}

func (n NoteOn) IsLiveMessage() {

}

func (n NoteOn) Velocity() uint8 {
	return n.velocity
}

func (n NoteOn) Channel() uint8 {
	return n.channel
}

func (n NoteOn) Raw() []byte {
	return channelMessage2(n.channel, 9, n.key, n.velocity)
}

func (n NoteOn) String() string {
	return fmt.Sprintf("%T channel %v key %v vel %v", n, n.channel, n.key, n.velocity)
}

func (NoteOn) set(channel, arg1, arg2 uint8) setter2 {
	var m NoteOn
	m.channel = channel
	m.key, m.velocity = parseTwoUint7(arg1, arg2)
	return m
}
