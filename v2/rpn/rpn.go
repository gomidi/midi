package rpn

import (
	"gitlab.com/gomidi/midi/v2"
)

/*
CC101 00 Selects RPN function.
CC100 00 Selects pitch bend as the parameter you want to adjust.
CC06 XX Sensitivity in half steps. The range is 0-24.
*/

func cc(channel, ctl uint8, val uint8) midi.Message {
	return midi.ControlChange(channel, ctl, val)
}

// PitchBendSensitivity sets the pitch bend range via RPN
func PitchBendSensitivity(channel, msbVal, lsbVal uint8) []midi.Message {
	return RPN(channel, 0, 0, msbVal, lsbVal)
}

// FineTuning
func FineTuning(channel, msbVal, lsbVal uint8) []midi.Message {
	return RPN(channel, 0, 1, msbVal, lsbVal)
}

// CoarseTuning
func CoarseTuning(channel, msbVal, lsbVal uint8) []midi.Message {
	return RPN(channel, 0, 2, msbVal, lsbVal)
}

// TuningProgramSelect
func TuningProgramSelect(channel, msbVal, lsbVal uint8) []midi.Message {
	return RPN(channel, 0, 3, msbVal, lsbVal)
}

// TuningBankSelect
func TuningBankSelect(channel, msbVal, lsbVal uint8) []midi.Message {
	return RPN(channel, 0, 4, msbVal, lsbVal)
}

// Reset aka Null
func Reset(channel uint8) []midi.Message {
	return append([]midi.Message{},
		cc(channel, 101, 127),
		cc(channel, 100, 127),
	)
}

// RPN message consisting of a val101 and val100 to identify the RPN and a msb and lsb for the value
func RPN(channel, val101, val100, msbVal, lsbVal uint8) []midi.Message {
	msgs := append([]midi.Message{},
		cc(channel, 101, val101),
		cc(channel, 100, val100),
		cc(channel, 6, msbVal),
		cc(channel, 38, lsbVal))

	return append(msgs, Reset(channel)...)
}

func Increment(channel, val101, val100 uint8) []midi.Message {
	msgs := append([]midi.Message{},
		cc(channel, 101, val101),
		cc(channel, 100, val100),
		cc(channel, 96, 0))

	return append(msgs, Reset(channel)...)
}

func Decrement(channel, val101, val100 uint8) []midi.Message {
	msgs := append([]midi.Message{},
		cc(channel, 101, val101),
		cc(channel, 100, val100),
		cc(channel, 97, 0))

	return append(msgs, Reset(channel)...)
}
