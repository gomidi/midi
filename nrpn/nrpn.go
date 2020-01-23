package nrpn

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
)

type Channel uint8

func (c Channel) cc(ctl uint8, val uint8) midi.Message {
	return channel.Channel(c).ControlChange(ctl, val)
}

// Reset aka Null
func (c Channel) Reset() []midi.Message {
	return append([]midi.Message{},
		c.cc(99, 127),
		c.cc(98, 127),
	)
}

func (c Channel) Increment(val99, val98 uint8) []midi.Message {
	msgs := append([]midi.Message{},
		c.cc(99, val99),
		c.cc(98, val98),
		c.cc(96, 0))

	return append(msgs, c.Reset()...)
}

func (c Channel) Decrement(val99, val98 uint8) []midi.Message {
	msgs := append([]midi.Message{},
		c.cc(99, val99),
		c.cc(98, val98),
		c.cc(97, 0))

	return append(msgs, c.Reset()...)
}

// NRPN message consisting of a val99 and val98 to identify the NRPN and a msb and lsb for the value
func (c Channel) NRPN(val99, val98, msbVal, lsbVal uint8) []midi.Message {
	msgs := append([]midi.Message{},
		c.cc(99, val99),
		c.cc(98, val98),
		c.cc(6, msbVal),
		c.cc(38, lsbVal))

	return append(msgs, c.Reset()...)
}
