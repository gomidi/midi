package midi

// ResetChannel "resets" channel to some established defaults
/*
  bank select -> bank
  program change -> prog
  cc all controllers off
  cc volume -> 100
  cc expression -> 127
  cc hold pedal -> off
  cc pan position -> 64
  cc RPN pitch bend sensitivity -> 2 (semitones)
*/
func ResetChannel(ch uint8, bank uint8, prog uint8) []Message {
	var msgs = []Message{
		ControlChange(ch, BankSelectMSB, bank),
		ProgramChange(ch, prog),
		ControlChange(ch, AllControllersOff, 0),
		ControlChange(ch, VolumeMSB, 100),
		ControlChange(ch, ExpressionMSB, 127),
		ControlChange(ch, HoldPedalSwitch, 0),
		ControlChange(ch, PanPositionMSB, 64),
	}

	//err := PitchBendSensitivityRPN(w, 2, 0)

	return msgs
}

// SilenceChannel sends a  note off message for every running note
// on the given channel. If the channel is -1, every channel is affected.
// If note consolidation is switched off, notes aren't tracked and therefor
// every possible note will get a note off. This could also be enforced by setting
// force to true. If force is true, additionally the cc messages AllNotesOff (123) and AllSoundOff (120)
// are send.
// If channel is > 15, the method panics.
// The error or not error of the last write is returned.
// (i.e. it does not return on the first error, but tries everything instead to make it silent)
func SilenceChannel(ch int8) (out []Message) {
	if ch > 15 {
		panic("invalid channel number")
	}

	// single channel
	if ch >= 0 {
		out = silentium(uint8(ch))
		return
	}

	// all channels
	for c := 0; c < 16; c++ {
		out = append(out, silentium(uint8(c))...)
	}
	return
}

// silentium for a single channel
func silentium(ch uint8) (out []Message) {
	out = append(out, ControlChange(ch, AllNotesOff, 0))
	out = append(out, ControlChange(ch, AllSoundOff, 0))
	return
}
