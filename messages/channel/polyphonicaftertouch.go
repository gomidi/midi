package channel

import "fmt"

// PolyphonicAfterTouch represents a MIDI polyphonic aftertouch message (aka "key pressure")
type PolyphonicAfterTouch struct {
	channel  uint8
	key      uint8
	pressure uint8
}

func (p PolyphonicAfterTouch) Key() uint8 {
	return p.key
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
	return fmt.Sprintf("%T (\"KeyPressure\") channel %v key %v pressure %v", p, p.channel, p.key, p.pressure)
}

func (p PolyphonicAfterTouch) Raw() []byte {
	return channelMessage2(p.channel, 10, p.key, p.pressure)
}

func (PolyphonicAfterTouch) set(channel uint8, arg1, arg2 uint8) setter2 {
	var m PolyphonicAfterTouch
	m.channel = channel
	m.key, m.pressure = parseTwoUint7(arg1, arg2)
	return m
}
