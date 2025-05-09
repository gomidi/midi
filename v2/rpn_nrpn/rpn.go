package rpn_nrpn

import "gitlab.com/gomidi/midi/v2"

/*
CC101 00 Selects RPN function.
CC100 00 Selects pitch bend as the parameter you want to adjust.
CC06 XX Sensitivity in half steps. The range is 0-24.
*/

// Reset aka Null
func RPNReset(channel uint8) []midi.Message {
	return append([]midi.Message{},
		cc(channel, CC_RPN0, VAL_SET),
		cc(channel, CC_RPN1, VAL_SET),
	)
}

// RPN message consisting of a val101 and val100 to identify the RPN and a msb and lsb for the value
func RPN(channel, val101, val100, msbVal, lsbVal uint8) []midi.Message {
	return []midi.Message{
		cc(channel, CC_RPN0, val101),
		cc(channel, CC_RPN1, val100),
		cc(channel, CC_MSB, msbVal),
		cc(channel, CC_LSB, lsbVal),
	}
}

func RPNIncrement(channel, val101, val100 uint8) []midi.Message {
	return []midi.Message{
		cc(channel, CC_RPN0, val101),
		cc(channel, CC_RPN1, val100),
		cc(channel, CC_INC, VAL_UNSET),
	}
}

func RPNDecrement(channel, val101, val100 uint8) []midi.Message {
	return []midi.Message{
		cc(channel, CC_RPN0, val101),
		cc(channel, CC_RPN1, val100),
		cc(channel, CC_DEC, VAL_UNSET),
	}
}
