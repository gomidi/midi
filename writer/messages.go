package writer

import (
	"fmt"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/cc"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/midi/midimessage/meta/meter"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midimessage/syscommon"
	"gitlab.com/gomidi/midi/midimessage/sysex"
	"gitlab.com/gomidi/midi/nrpn"
	"gitlab.com/gomidi/midi/rpn"
)

// ActiveSensing writes the active sensing realtime message
func RTActivesense(w *Writer) error {
	return w.wr.Write(realtime.Activesense)
}

// Continue writes the continue realtime message
func RTContinue(w *Writer) error {
	return w.wr.Write(realtime.Continue)
}

// Reset writes the reset realtime message
func RTReset(w *Writer) error {
	return w.wr.Write(realtime.Reset)
}

// Start writes the start realtime message
func RTStart(w *Writer) error {
	return w.wr.Write(realtime.Start)
}

// Stop writes the stop realtime message
func RTStop(w *Writer) error {
	return w.wr.Write(realtime.Stop)
}

// Tick writes the tick realtime message
func RTTick(w *Writer) error {
	return w.wr.Write(realtime.Tick)
}

// Clock writes the timing clock realtime message
func RTClock(w *Writer) error {
	return w.wr.Write(realtime.TimingClock)
}

// MTC writes the MIDI Timing Code system message
func MTC(w *Writer, code uint8) error {
	return w.wr.Write(syscommon.MTC(code))
}

// SPP writes the song position pointer system message
func SPP(w *Writer, ptr uint16) error {
	return w.wr.Write(syscommon.SPP(ptr))
}

// SongSelect writes the song select system message
func SongSelect(w *Writer, song uint8) error {
	return w.wr.Write(syscommon.SongSelect(song))
}

// Tune writes the tune request system message
func Tune(w *Writer) error {
	return w.wr.Write(syscommon.Tune)
}

// EndOfTrack signals the end of a track
func EndOfTrack(w *SMF) error {
	w.Writer.noteState = [16][128]bool{}
	if no := w.wr.Header().NumTracks; w.finishedTracks >= no {
		return fmt.Errorf("too many tracks: in header: %v, closed: %v", no, w.finishedTracks+1)
	}
	w.finishedTracks++
	if w.timeline != nil {
		w.timeline.Reset()
	}
	return w.wr.Write(meta.EndOfTrack)
}

// Copyright writes the copyright meta message
func Copyright(w *SMF, text string) error {
	return w.wr.Write(meta.Copyright(text))
}

// Writes an undefined meta message
func Undefined(w *SMF, typ byte, bt []byte) error {
	return w.wr.Write(meta.Undefined{typ, bt})
}

// Cuepoint writes the cuepoint meta message
func Cuepoint(w *SMF, text string) error {
	return w.wr.Write(meta.Cuepoint(text))
}

// Device writes the device port meta message
func Device(w *SMF, port string) error {
	return w.wr.Write(meta.Device(port))
}

// KeySig writes the key signature meta message.
// A more comfortable way is to use the Key method in conjunction
// with the gomidi/midi/midimessage/meta/key package
func KeySig(w *SMF, key uint8, ismajor bool, num uint8, isflat bool) error {
	return w.wr.Write(meta.Key{Key: key, IsMajor: ismajor, Num: num, IsFlat: isflat})
}

// Key writes the given key signature meta message.
// It is supposed to be used with the gomidi/midi/midimessage/meta/key package
func Key(w *SMF, keysig meta.Key) error {
	return w.wr.Write(keysig)
}

// Lyric writes the lyric meta message
func Lyric(w *SMF, text string) error {
	return w.wr.Write(meta.Lyric(text))
}

// Marker writes the marker meta message
func Marker(w *SMF, text string) error {
	return w.wr.Write(meta.Marker(text))
}

// DeprecatedChannel writes the deprecated MIDI channel meta message
func DeprecatedChannel(w *SMF, ch uint8) error {
	return w.wr.Write(meta.Channel(ch))
}

// DeprecatedPort writes the deprecated MIDI port meta message
func DeprecatedPort(w *SMF, port uint8) error {
	return w.wr.Write(meta.Port(port))
}

// Program writes the program name meta message
func Program(w *SMF, text string) error {
	return w.wr.Write(meta.Program(text))
}

// TrackSequenceName writes the track / sequence name meta message
// If in a format 0 track, or the first track in a format 1 file, the name of the sequence. Otherwise, the name of the track.
func TrackSequenceName(w *SMF, name string) error {
	return w.wr.Write(meta.TrackSequenceName(name))
}

// SequenceNo writes the sequence number meta message
func SequenceNo(w *SMF, no uint16) error {
	return w.wr.Write(meta.SequenceNo(no))
}

// SequencerData writes a custom sequences specific meta message
func SequencerData(w *SMF, data []byte) error {
	return w.wr.Write(meta.SequencerData(data))
}

// SMPTE writes the SMPTE offset meta message
func SMPTE(w *SMF, hour, minute, second, frame, fractionalFrame byte) error {
	return w.wr.Write(meta.SMPTE{
		Hour:            hour,
		Minute:          minute,
		Second:          second,
		Frame:           frame,
		FractionalFrame: fractionalFrame,
	})
}

// Tempo writes the tempo meta message
func TempoBPM(w *SMF, bpm float64) error {
	return w.wr.Write(meta.FractionalBPM(bpm))
}

// Text writes the text meta message
func Text(w *SMF, text string) error {
	return w.wr.Write(meta.Text(text))
}

// Meter writes the time signature meta message in a more comfortable way.
// Numerator and Denominator are decimals.
func Meter(w *SMF, numerator, denominator uint8) error {
	w.timeline.AddTimeSignature(numerator, denominator)
	return w.wr.Write(meter.Meter(numerator, denominator))
}

// TimeSig writes the time signature meta message.
// Numerator and Denominator are decimal.
// If you don't want to deal with clocks per click and demisemiquaverperquarter,
// user the Meter method instead.
func TimeSig(w *SMF, numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) error {
	w.timeline.AddTimeSignature(numerator, denominator)
	return w.wr.Write(meta.TimeSig{
		Numerator:                numerator,
		Denominator:              denominator,
		ClocksPerClick:           clocksPerClick,
		DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
	})
}

// Instrument writes the instrument name meta message
func Instrument(w *SMF, name string) error {
	return w.wr.Write(meta.Instrument(name))
}

// GMReset resets the channel to some GM based standards (see Reset) and sets the given GM program.
func GMReset(w ChannelWriter, prog uint8) error {
	return Reset(w, 0, prog)
}

func Channel(w ChannelWriter) channel.Channel {
	return channel.Channel(w.Channel())
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
func Reset(w ChannelWriter, bank uint8, prog uint8) error {
	c := Channel(w)
	var msgs = []midi.Message{
		c.ControlChange(cc.BankSelectMSB, bank),
		c.ProgramChange(prog),
		c.ControlChange(cc.AllControllersOff, 0),
		c.ControlChange(cc.VolumeMSB, 100),
		c.ControlChange(cc.ExpressionMSB, 127),
		c.ControlChange(cc.HoldPedalSwitch, 0),
		c.ControlChange(cc.PanPositionMSB, 64),
	}

	for _, msg := range msgs {
		err := w.Write(msg)

		if err != nil {
			return fmt.Errorf("could not reset channel %v: %v", c.Channel(), err)
		}
	}

	err := PitchBendSensitivityRPN(w, 2, 0)

	if err != nil {
		return fmt.Errorf("could not reset channel %v: %v", c.Channel(), err)
	}
	return nil
}

// Aftertouch writes a channel pressure message for the current channel
func Aftertouch(w ChannelWriter, pressure uint8) error {
	return w.Write(Channel(w).Aftertouch(pressure))
}

// PolyAftertouch writes a key pressure message for the current channel
func PolyAftertouch(w ChannelWriter, key, pressure uint8) error {
	return w.Write(Channel(w).PolyAftertouch(key, pressure))
}

// NoteOff writes a note off message for the current channel
// By default, midi notes are consolidated (see ConsolidateNotes method)
func NoteOff(w ChannelWriter, key uint8) error {
	return w.Write(Channel(w).NoteOff(key))
}

// NoteOffVelocity writes a note off message for the current channel with a velocity.
// By default, midi notes are consolidated (see ConsolidateNotes method)
func NoteOffVelocity(w ChannelWriter, key, velocity uint8) error {
	return w.Write(Channel(w).NoteOffVelocity(key, velocity))
}

// NoteOn writes a note on message for the current channel
// By default, midi notes are consolidated (see ConsolidateNotes method)
func NoteOn(w ChannelWriter, key, velocity uint8) error {
	return w.Write(Channel(w).NoteOn(key, velocity))
}

// Pitchbend writes a pitch bend message for the current channel
// For reset value, use 0, for lowest -8191 and highest 8191
// Or use the pitch constants of midimessage/channel
func Pitchbend(w ChannelWriter, value int16) error {
	return w.Write(Channel(w).Pitchbend(value))
}

// ProgramChange writes a program change message for the current channel
// Program numbers start with 0 for program 1.
func ProgramChange(w ChannelWriter, program uint8) error {
	return w.Write(Channel(w).ProgramChange(program))
}

/*
CC101 00 Selects RPN function.
CC100 00 Selects pitch bend as the parameter you want to adjust.
CC06 XX Sensitivity in half steps. The range is 0-24.
*/

// PitchBendSensitivityRPN sets the pitch bend range via RPN
func PitchBendSensitivityRPN(w ChannelWriter, msbVal, lsbVal uint8) error {
	return RPN(w, 0, 0, msbVal, lsbVal)
}

// FineTuningRPN
func FineTuningRPN(w ChannelWriter, msbVal, lsbVal uint8) error {
	return RPN(w, 0, 1, msbVal, lsbVal)
}

// CoarseTuningRPN
func CoarseTuningRPN(w ChannelWriter, msbVal, lsbVal uint8) error {
	return RPN(w, 0, 2, msbVal, lsbVal)
}

// TuningProgramSelectRPN
func TuningProgramSelectRPN(w ChannelWriter, msbVal, lsbVal uint8) error {
	return RPN(w, 0, 3, msbVal, lsbVal)
}

// TuningBankSelectRPN
func TuningBankSelectRPN(w ChannelWriter, msbVal, lsbVal uint8) error {
	return RPN(w, 0, 4, msbVal, lsbVal)
}

// ResetRPN aka Null
func ResetRPN(w ChannelWriter) error {
	err := WriteMessages(w, rpn.Channel(w.Channel()).Reset())
	if err != nil {
		err = fmt.Errorf("can't write ResetRPN: %v", err)
	}
	return err
}

func WriteMessages(w ChannelWriter, msgs []midi.Message) error {
	for _, msg := range msgs {
		err := w.Write(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// RPN message consisting of a val101 and val100 to identify the RPN and a msb and lsb for the value
func RPN(w ChannelWriter, val101, val100, msbVal, lsbVal uint8) error {
	err := WriteMessages(w, rpn.Channel(w.Channel()).RPN(val101, val100, msbVal, lsbVal))
	if err != nil {
		err = fmt.Errorf("can't write RPN(%v,%v): %v", val101, val100, err)
	}
	return err
}

func RPNIncrement(w ChannelWriter, val101, val100 uint8) error {
	err := WriteMessages(w, rpn.Channel(w.Channel()).Increment(val101, val100))
	if err != nil {
		err = fmt.Errorf("can't write RPNIncrement(%v,%v): %v", val101, val100, err)
	}
	return err
}

func RPNDecrement(w ChannelWriter, val101, val100 uint8) error {
	err := WriteMessages(w, rpn.Channel(w.Channel()).Decrement(val101, val100))
	if err != nil {
		err = fmt.Errorf("can't write RPNDecrement(%v,%v): %v", val101, val100, err)
	}
	return err
}

func NRPNIncrement(w ChannelWriter, val99, val98 uint8) error {
	err := WriteMessages(w, nrpn.Channel(w.Channel()).Increment(val99, val98))
	if err != nil {
		err = fmt.Errorf("can't write NRPNIncrement(%v,%v): %v", val99, val98, err)
	}
	return err
}

func NRPNDecrement(w ChannelWriter, val99, val98 uint8) error {
	err := WriteMessages(w, nrpn.Channel(w.Channel()).Decrement(val99, val98))
	if err != nil {
		err = fmt.Errorf("can't write NRPNDecrement(%v,%v): %v", val99, val98, err)
	}
	return err
}

// NRPN message consisting of a val99 and val98 to identify the RPN and a msb and lsb for the value
func NRPN(w ChannelWriter, val99, val98, msbVal, lsbVal uint8) error {
	err := WriteMessages(w, nrpn.Channel(w.Channel()).NRPN(val99, val98, msbVal, lsbVal))
	if err != nil {
		err = fmt.Errorf("can't write NRPN(%v,%v): %v", val99, val98, err)
	}
	return err
}

// ResetNRPN aka Null
func ResetNRPN(w ChannelWriter) error {
	err := WriteMessages(w, nrpn.Channel(w.Channel()).Reset())
	if err != nil {
		err = fmt.Errorf("can't write ResetNRPN: %v", err)
	}
	return err
}

// MsbLsb writes a Msb control change message, followed by a Lsb control change message
// for the current channel
// For more comfortable use, used it in conjunction with the gomidi/cc package
func MsbLsb(w ChannelWriter, msbControl, lsbControl uint8, value uint16) error {

	var b = make([]byte, 2)
	b[1] = byte(value & 0x7F)
	b[0] = byte((value >> 7) & 0x7F)

	/*
		r := midilib.MsbLsbSigned(value)

		var b = make([]byte, 2)

		binary.BigEndian.PutUint16(b, r)
	*/
	err := ControlChange(w, msbControl, b[0])
	if err != nil {
		return err
	}
	return ControlChange(w, lsbControl, b[1])
}

// ControlChange writes a control change message for the current channel
// For more comfortable use, used it in conjunction with the gomidi/cc package
func ControlChange(w ChannelWriter, controller, value uint8) error {
	return w.Write(Channel(w).ControlChange(controller, value))
}

// CcOff writes a control change message with a value of 0 (=off) for the current channel
func CcOff(w ChannelWriter, controller uint8) error {
	return ControlChange(w, controller, 0)
}

// CcOn writes a control change message with a value of 127 (=on) for the current channel
func CcOn(w ChannelWriter, controller uint8) error {
	return ControlChange(w, controller, 127)
}

// SysEx writes system exclusive data
func SysEx(w ChannelWriter, data []byte) error {
	return w.Write(sysex.SysEx(data))
}
