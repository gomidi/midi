package rpn

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
)

/*
CC101 00 Selects RPN function.
CC100 00 Selects pitch bend as the parameter you want to adjust.
CC06 XX Sensitivity in half steps. The range is 0-24.
*/

type Channel uint8

func (c Channel) cc(ctl uint8, val uint8) midi.Message {
	return channel.Channel(c).ControlChange(ctl, val)
}

// PitchBendSensitivity sets the pitch bend range via RPN
func (c Channel) PitchBendSensitivity(msbVal, lsbVal uint8) []midi.Message {
	return c.RPN(0, 0, msbVal, lsbVal)
}

// FineTuning
func (c Channel) FineTuning(msbVal, lsbVal uint8) []midi.Message {
	return c.RPN(0, 1, msbVal, lsbVal)
}

// CoarseTuning
func (c Channel) CoarseTuning(msbVal, lsbVal uint8) []midi.Message {
	return c.RPN(0, 2, msbVal, lsbVal)
}

// TuningProgramSelect
func (c Channel) TuningProgramSelect(msbVal, lsbVal uint8) []midi.Message {
	return c.RPN(0, 3, msbVal, lsbVal)
}

// TuningBankSelect
func (c Channel) TuningBankSelect(msbVal, lsbVal uint8) []midi.Message {
	return c.RPN(0, 4, msbVal, lsbVal)
}

// Reset aka Null
func (c Channel) Reset() []midi.Message {
	return append([]midi.Message{},
		c.cc(101, 127),
		c.cc(100, 127),
	)
}

// RPN message consisting of a val101 and val100 to identify the RPN and a msb and lsb for the value
func (c Channel) RPN(val101, val100, msbVal, lsbVal uint8) []midi.Message {
	msgs := append([]midi.Message{},
		c.cc(101, val101),
		c.cc(100, val100),
		c.cc(6, msbVal),
		c.cc(38, lsbVal))

	return append(msgs, c.Reset()...)
}

func (c Channel) Increment(val101, val100 uint8) []midi.Message {
	msgs := append([]midi.Message{},
		c.cc(101, val101),
		c.cc(100, val100),
		c.cc(96, 0))

	return append(msgs, c.Reset()...)
}

func (c Channel) Decrement(val101, val100 uint8) []midi.Message {
	msgs := append([]midi.Message{},
		c.cc(101, val101),
		c.cc(100, val100),
		c.cc(97, 0))

	return append(msgs, c.Reset()...)
}
