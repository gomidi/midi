package rpn_nrpn

import (
	"gitlab.com/gomidi/midi/v2"
)

// Reset aka Null
func NRPNReset(channel uint8) []midi.Message {
	return append([]midi.Message{},
		cc(channel, CC_NRPN0, VAL_SET),
		cc(channel, CC_NRPN1, VAL_SET),
	)
}

func NRPNIncrement(channel uint8, val99, val98 uint8) []midi.Message {
	return []midi.Message{
		cc(channel, CC_NRPN0, val99),
		cc(channel, CC_NRPN1, val98),
		cc(channel, CC_INC, CC_INC),
	}
}

func NRPNDecrement(channel uint8, val99, val98 uint8) []midi.Message {
	return []midi.Message{
		cc(channel, CC_NRPN0, val99),
		cc(channel, CC_NRPN1, val98),
		cc(channel, CC_DEC, VAL_UNSET),
	}
}

// NRPN message consisting of a val99 and val98 to identify the NRPN and a msb and lsb for the value
func NRPN(channel uint8, val99, val98, msbVal, lsbVal uint8) []midi.Message {
	return []midi.Message{
		cc(channel, CC_NRPN0, val99),
		cc(channel, CC_NRPN1, val98),
		cc(channel, CC_MSB, msbVal),
		cc(channel, CC_LSB, lsbVal),
	}
}
