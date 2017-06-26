package channel

import "fmt"

type PolyphonicAfterTouch struct {
	channel  uint8
	pitch    uint8
	pressure uint8
}

func (p PolyphonicAfterTouch) Pitch() uint8 {
	return p.pitch
}

func (p PolyphonicAfterTouch) Pressure() uint8 {
	return p.pressure
}

func (p PolyphonicAfterTouch) IsLiveMessage() {

}

func (p PolyphonicAfterTouch) Channel() uint8 {
	return p.channel
}

func (p PolyphonicAfterTouch) String() string {
	return fmt.Sprintf("%T channel %v pitch %v pressure %v", p, p.channel, p.pitch, p.pressure)
}

func (p PolyphonicAfterTouch) Raw() []byte {
	return channelMessage2(p.channel, 10, p.pitch, p.pressure)
}

func (PolyphonicAfterTouch) set(channel uint8, arg1, arg2 uint8) setter2 {
	var m PolyphonicAfterTouch
	m.channel = channel
	m.pitch, m.pressure = parseTwoUint7(arg1, arg2)
	return m
}
