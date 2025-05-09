package rpn_nrpn

import "gitlab.com/gomidi/midi/v2"

func cc(channel, ctl, val uint8) midi.Message {
	return midi.ControlChange(channel, ctl, val)
}

var (
	CC_RPN0 = uint8(101)
	CC_RPN1 = uint8(100)

	CC_NRPN0 = uint8(99)
	CC_NRPN1 = uint8(98)

	CC_INC = uint8(96)
	CC_DEC = uint8(97)

	CC_MSB = uint8(6)
	CC_LSB = uint8(38)

	VAL_SET   = uint8(127)
	VAL_UNSET = uint8(0)
)

func IsRPN_NRPN_CC(cc uint8) bool {
	switch cc {
	case CC_RPN0, CC_RPN1, CC_NRPN0, CC_NRPN1, CC_INC, CC_DEC, CC_MSB, CC_LSB:
		return true
	default:
		return false
	}
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
