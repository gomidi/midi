package channel

import (
	"fmt"
	"github.com/gomidi/midi/internal/lib"
)

type AfterTouch struct {
	channel  uint8
	pressure uint8
}

func (a AfterTouch) Pressure() uint8 {
	return a.pressure
}

func (a AfterTouch) Channel() uint8 {
	return a.channel
}

func (a AfterTouch) Raw() []byte {
	return channelMessage1(a.channel, 13, a.pressure)
}

func (a AfterTouch) String() string {
	return fmt.Sprintf("%T channel %v pressure %v", a, a.channel, a.pressure)
}

func (AfterTouch) set(channel uint8, firstArg uint8) setter1 {
	var m AfterTouch
	m.channel = channel
	m.pressure = lib.ParseUint7(firstArg)
	return m
}
