package mid

import (
	"fmt"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/cc"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/sysex"
)

type midiWriter struct {
	wr              midi.Writer
	Channel         channel.Channel
	noteState       [16][128]bool
	noConsolidation bool
}

// GMReset resets the channel to some GM based standards (see Reset) and sets the given GM program.
func (w *midiWriter) GMReset(prog uint8) error {
	return w.Reset(0, prog)
}

// Reset "resets" channel to some established defaults
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
func (w *midiWriter) Reset(bank uint8, prog uint8) error {
	var msgs = []midi.Message{
		w.Channel.ControlChange(cc.BankSelectMSB, bank),
		w.Channel.ProgramChange(prog),
		w.Channel.ControlChange(cc.AllControllersOff, 0),
		w.Channel.ControlChange(cc.VolumeMSB, 100),
		w.Channel.ControlChange(cc.ExpressionMSB, 127),
		w.Channel.ControlChange(cc.HoldPedalSwitch, 0),
		w.Channel.ControlChange(cc.PanPositionMSB, 64),
	}

	for _, msg := range msgs {
		err := w.Write(msg)

		if err != nil {
			return fmt.Errorf("could not reset channel %v: %v", w.Channel.Channel(), err)
		}
	}

	err := w.PitchBendSensitivityRPN(2, 0)

	if err != nil {
		return fmt.Errorf("could not reset channel %v: %v", w.Channel.Channel(), err)
	}
	return nil
}

// SetChannel sets the channel for the following midi messages
// Channel numbers are counted from 0 to 15 (MIDI channel 1 to 16).
// The initial channel number is 0.
func (w *midiWriter) SetChannel(no uint8 /* 0-15 */) {
	w.Channel = channel.Channel(no)
}

// Aftertouch writes a channel pressure message for the current channel
func (w *midiWriter) Aftertouch(pressure uint8) error {
	return w.wr.Write(w.Channel.Aftertouch(pressure))
}

// PolyAftertouch writes a key pressure message for the current channel
func (w *midiWriter) PolyAftertouch(key, pressure uint8) error {
	return w.wr.Write(w.Channel.PolyAftertouch(key, pressure))
}

// NoteOff writes a note off message for the current channel
// By default, midi notes are consolidated (see ConsolidateNotes method)
func (w *midiWriter) NoteOff(key uint8) error {
	return w.Write(w.Channel.NoteOff(key))
}

// NoteOffVelocity writes a note off message for the current channel with a velocity.
// By default, midi notes are consolidated (see ConsolidateNotes method)
func (w *midiWriter) NoteOffVelocity(key, velocity uint8) error {
	return w.Write(w.Channel.NoteOffVelocity(key, velocity))
}

// NoteOn writes a note on message for the current channel
// By default, midi notes are consolidated (see ConsolidateNotes method)
func (w *midiWriter) NoteOn(key, veloctiy uint8) error {
	return w.Write(w.Channel.NoteOn(key, veloctiy))
}

// Pitchbend writes a pitch bend message for the current channel
// For reset value, use 0, for lowest -8191 and highest 8191
// Or use the pitch constants of midimessage/channel
func (w *midiWriter) Pitchbend(value int16) error {
	return w.wr.Write(w.Channel.Pitchbend(value))
}

// ProgramChange writes a program change message for the current channel
// Program numbers start with 0 for program 1.
func (w *midiWriter) ProgramChange(program uint8) error {
	return w.wr.Write(w.Channel.ProgramChange(program))
}

/*
CC101 00 Selects RPN function.
CC100 00 Selects pitch bend as the parameter you want to adjust.
CC06 XX Sensitivity in half steps. The range is 0-24.
*/

// PitchBendSensitivityRPN sets the pitch bend range via RPN
func (w *midiWriter) PitchBendSensitivityRPN(msbVal, lsbVal uint8) error {
	return w.RPN(0, 0, msbVal, lsbVal)
}

// FineTuningRPN
func (w *midiWriter) FineTuningRPN(msbVal, lsbVal uint8) error {
	return w.RPN(0, 1, msbVal, lsbVal)
}

// CoarseTuningRPN
func (w *midiWriter) CoarseTuningRPN(msbVal, lsbVal uint8) error {
	return w.RPN(0, 2, msbVal, lsbVal)
}

// TuningProgramSelectRPN
func (w *midiWriter) TuningProgramSelectRPN(msbVal, lsbVal uint8) error {
	return w.RPN(0, 3, msbVal, lsbVal)
}

// TuningBankSelectRPN
func (w *midiWriter) TuningBankSelectRPN(msbVal, lsbVal uint8) error {
	return w.RPN(0, 4, msbVal, lsbVal)
}

// ResetRPN aka Null
func (w *midiWriter) ResetRPN() error {
	msgs := append([]midi.Message{},
		w.Channel.ControlChange(101, 127),
		w.Channel.ControlChange(100, 127),
	)
	for _, msg := range msgs {
		err := w.wr.Write(msg)
		if err != nil {
			return fmt.Errorf("can't write ResetRPN: %v", msg)
		}
	}
	return nil
}

// RPN message consisting of a val101 and val100 to identify the RPN and a msb and lsb for the value
func (w *midiWriter) RPN(val101, val100, msbVal, lsbVal uint8) error {
	msgs := append([]midi.Message{},
		w.Channel.ControlChange(101, val101),
		w.Channel.ControlChange(100, val100),
		w.Channel.ControlChange(6, msbVal),
		w.Channel.ControlChange(38, lsbVal))

	for _, msg := range msgs {
		err := w.wr.Write(msg)
		if err != nil {
			return fmt.Errorf("can't write RPN(%v,%v): %v", val101, val100, msg)
		}
	}

	return w.ResetRPN()
}

func (w *midiWriter) RPNIncrement(val101, val100 uint8) error {
	msgs := append([]midi.Message{},
		w.Channel.ControlChange(101, val101),
		w.Channel.ControlChange(100, val100),
		w.Channel.ControlChange(96, 0))

	for _, msg := range msgs {
		err := w.wr.Write(msg)
		if err != nil {
			return fmt.Errorf("can't write RPNIncrement(%v,%v): %v", val101, val100, msg)
		}
	}

	return w.ResetRPN()
}

func (w *midiWriter) RPNDecrement(val101, val100 uint8) error {
	msgs := append([]midi.Message{},
		w.Channel.ControlChange(101, val101),
		w.Channel.ControlChange(100, val100),
		w.Channel.ControlChange(97, 0))

	for _, msg := range msgs {
		err := w.wr.Write(msg)
		if err != nil {
			return fmt.Errorf("can't write RPNDecrement(%v,%v): %v", val101, val100, msg)
		}
	}

	return w.ResetRPN()
}

func (w *midiWriter) NRPNIncrement(val99, val98 uint8) error {
	msgs := append([]midi.Message{},
		w.Channel.ControlChange(99, val99),
		w.Channel.ControlChange(98, val98),
		w.Channel.ControlChange(96, 0))

	for _, msg := range msgs {
		err := w.wr.Write(msg)
		if err != nil {
			return fmt.Errorf("can't write NRPNIncrement(%v,%v): %v", val99, val98, msg)
		}
	}
	return w.ResetNRPN()
}

func (w *midiWriter) NRPNDecrement(val99, val98 uint8) error {
	msgs := append([]midi.Message{},
		w.Channel.ControlChange(99, val99),
		w.Channel.ControlChange(98, val98),
		w.Channel.ControlChange(97, 0))

	for _, msg := range msgs {
		err := w.wr.Write(msg)
		if err != nil {
			return fmt.Errorf("can't write NRPNDecrement(%v,%v): %v", val99, val98, msg)
		}
	}
	return w.ResetNRPN()
}

// NRPN message consisting of a val99 and val98 to identify the RPN and a msb and lsb for the value
func (w *midiWriter) NRPN(val99, val98, msbVal, lsbVal uint8) error {
	msgs := append([]midi.Message{},
		w.Channel.ControlChange(99, val99),
		w.Channel.ControlChange(98, val98),
		w.Channel.ControlChange(6, msbVal),
		w.Channel.ControlChange(38, lsbVal))

	for _, msg := range msgs {
		err := w.wr.Write(msg)
		if err != nil {
			return fmt.Errorf("can't write NRPN(%v,%v): %v", val99, val98, msg)
		}
	}
	return w.ResetNRPN()
}

// ResetNRPN aka Null
func (w *midiWriter) ResetNRPN() error {
	msgs := append([]midi.Message{},
		w.Channel.ControlChange(99, 127),
		w.Channel.ControlChange(98, 127),
	)
	for _, msg := range msgs {
		err := w.wr.Write(msg)
		if err != nil {
			return fmt.Errorf("can't write ResetNRPN: %v", msg)
		}
	}
	return nil
}

// MsbLsb writes a Msb control change message, followed by a Lsb control change message
// for the current channel
// For more comfortable use, used it in conjunction with the gomidi/cc package
func (w *midiWriter) MsbLsb(msbControl, lsbControl uint8, value uint16) error {

	var b = make([]byte, 2)
	b[1] = byte(value & 0x7F)
	b[0] = byte((value >> 7) & 0x7F)

	/*
		r := midilib.MsbLsbSigned(value)

		var b = make([]byte, 2)

		binary.BigEndian.PutUint16(b, r)
	*/
	err := w.ControlChange(msbControl, b[0])
	if err != nil {
		return err
	}
	return w.ControlChange(lsbControl, b[1])
}

// ControlChange writes a control change message for the current channel
// For more comfortable use, used it in conjunction with the gomidi/cc package
func (w *midiWriter) ControlChange(controller, value uint8) error {
	return w.wr.Write(w.Channel.ControlChange(controller, value))
}

// CcOff writes a control change message with a value of 0 (=off) for the current channel
func (w *midiWriter) CcOff(controller uint8) error {
	return w.ControlChange(controller, 0)
}

// CcOn writes a control change message with a value of 127 (=on) for the current channel
func (w *midiWriter) CcOn(controller uint8) error {
	return w.ControlChange(controller, 127)
}

// SysEx writes system exclusive data
func (w *midiWriter) SysEx(data []byte) error {
	return w.wr.Write(sysex.SysEx(data))
}

// silentium for a single channel
func (w *midiWriter) silentium(ch uint8, force bool) (err error) {
	for key, val := range w.noteState[ch] {
		if force || val {
			err = w.wr.Write(channel.Channel(ch).NoteOff(uint8(key)))
		}
	}

	if force {
		err = w.wr.Write(channel.Channel(ch).ControlChange(cc.AllNotesOff, 0))
		err = w.wr.Write(channel.Channel(ch).ControlChange(cc.AllSoundOff, 0))
	}
	return
}

// Silence sends a  note off message for every running note
// on the given channel. If the channel is -1, every channel is affected.
// If note consolidation is switched off, notes aren't tracked and therefor
// every possible note will get a note off. This could also be enforced by setting
// force to true. If force is true, additionally the cc messages AllNotesOff (123) and AllSoundOff (120)
// are send.
// If channel is > 15, the method panics.
// The error or not error of the last write is returned.
// (i.e. it does not return on the first error, but tries everything instead to make it silent)
func (w *midiWriter) Silence(ch int8, force bool) (err error) {
	if ch > 15 {
		panic("invalid channel number")
	}
	if w.noConsolidation {
		force = true
	}

	// single channel
	if ch >= 0 {
		err = w.silentium(uint8(ch), force)
		// set note states for the channel
		w.noteState[ch] = [128]bool{}
		return
	}

	// all channels
	for c := 0; c < 16; c++ {
		err = w.silentium(uint8(c), force)
	}
	// reset all note states
	w.noteState = [16][128]bool{}
	return
}

// ConsolidateNotes enables/disables the midi note consolidation (default: enabled)
// When enabled, midi note on/off messages are consolidated, that means the on/off state of
// every possible note on every channel is tracked and note on messages are only
// written, if the corresponding note is off and vice versa. Note on messages
// with a velocity of 0 behave the same way as note offs.
// The tracking of the notes is not cross SMF tracks, i.e. a meta.EndOfTrack
// message will reset the tracking.
// The consolidation should prevent "hanging" notes in most cases.
// If on is true, the note will be started tracking again (fresh state), assuming no note is currently running.
func (w *midiWriter) ConsolidateNotes(on bool) {
	if on {
		w.noteState = [16][128]bool{}
	}
	w.noConsolidation = !on
}

// Write writes the given midi.Message. By default, midi notes are consolidated (see ConsolidateNotes method)
func (w *midiWriter) Write(msg midi.Message) error {
	if w.noConsolidation {
		return w.wr.Write(msg)
	}
	switch m := msg.(type) {
	case channel.NoteOn:
		if m.Velocity() > 0 && w.noteState[m.Channel()][m.Key()] {
			return fmt.Errorf("can't write %s. note already running.", msg)
		}
		if m.Velocity() == 0 && !w.noteState[m.Channel()][m.Key()] {
			return fmt.Errorf("can't write %s. note is not running.", msg)
		}
		w.noteState[m.Channel()][m.Key()] = m.Velocity() > 0
	case channel.NoteOff:
		if !w.noteState[m.Channel()][m.Key()] {
			return fmt.Errorf("can't write %s. note is not running.", msg)
		}
		w.noteState[m.Channel()][m.Key()] = false
	case channel.NoteOffVelocity:
		if !w.noteState[m.Channel()][m.Key()] {
			return fmt.Errorf("can't write %s. note is not running.", msg)
		}
		w.noteState[m.Channel()][m.Key()] = false
	}
	return w.wr.Write(msg)
}
