package rpn_nrpn

import "gitlab.com/gomidi/midi/v2"

func cc(channel, ctl, val uint8) midi.Message {
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
