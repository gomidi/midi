package nrpn

import (
	"gitlab.com/gomidi/midi/v2"
)

func cc(channel, ctl, val uint8) midi.Message {
	return midi.ControlChange(channel, ctl, val)
}

// Reset aka Null
func Reset(channel uint8) []midi.Message {
	return append([]midi.Message{},
		cc(channel, 99, 127),
		cc(channel, 98, 127),
	)
}

func Increment(channel uint8, val99, val98 uint8) []midi.Message {
	msgs := append([]midi.Message{},
		cc(channel, 99, val99),
		cc(channel, 98, val98),
		cc(channel, 96, 0))

	return append(msgs, Reset(channel)...)
}

func Decrement(channel uint8, val99, val98 uint8) []midi.Message {
	msgs := append([]midi.Message{},
		cc(channel, 99, val99),
		cc(channel, 98, val98),
		cc(channel, 97, 0))

	return append(msgs, Reset(channel)...)
}

// NRPN message consisting of a val99 and val98 to identify the NRPN and a msb and lsb for the value
func NRPN(channel uint8, val99, val98, msbVal, lsbVal uint8) []midi.Message {
	msgs := append([]midi.Message{},
		cc(channel, 99, val99),
		cc(channel, 98, val98),
		cc(channel, 6, msbVal),
		cc(channel, 38, lsbVal))

	return append(msgs, Reset(channel)...)
}
