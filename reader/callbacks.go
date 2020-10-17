package reader

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/smf"
)

func SMFHeader(cb func(smf.Header)) func(r *Reader) {
	return func(r *Reader) {
		r.smfheader = cb
	}
}

// Each is called for every MIDI message in addition to the other callbacks.
func Each(cb func(*Position, midi.Message)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Each = cb
	}
}

// Unknown is called for unknown messages
func Unknown(cb func(p *Position, msg midi.Message)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Unknown = cb
	}
}

// Copyright is called for the copyright message
func Copyright(cb func(p Position, text string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Copyright = cb
	}
}

// TempoBPM is called for the tempo (change) message, BPM is fractional
func TempoBPM(cb func(p Position, bpm float64)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.TempoBPM = cb
	}
}

// TimeSigis called for the time signature (change) message
func TimeSig(cb func(p Position, num, denom uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.TimeSig = cb
	}
}

// Key is called for the key signature (change) message
func Key(cb func(p Position, key uint8, ismajor bool, num_accidentals uint8, accidentals_are_flat bool)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Key = cb
	}
}

// Instrument is called for the instrument (name) message
func Instrument(cb func(p Position, name string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Instrument = cb
	}
}

// TrackSequenceName is called for the sequence / track name message
// If in a format 0 track, or the first track in a format 1 file, the name of the sequence. Otherwise, the name of the track.
func TrackSequenceName(cb func(p Position, name string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.TrackSequenceName = cb
	}
}

// SequenceNo is called for the sequence number message
func SequenceNo(cb func(p Position, number uint16)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.SequenceNo = cb
	}
}

// Marker is called for the marker message
func Marker(cb func(p Position, text string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Marker = cb
	}
}

// Cuepoint is called for the cuepoint message
func Cuepoint(cb func(p Position, text string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Cuepoint = cb
	}
}

// Text is called for the text message
func Text(cb func(p Position, text string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Text = cb
	}
}

// Lyric is called for the lyric message
func Lyric(cb func(p Position, text string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Lyric = cb
	}
}

// EndOfTrack is called for the end of a track message
func EndOfTrack(cb func(p Position)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.EndOfTrack = cb
	}
}

// Device is called for the device port message
func Device(cb func(p Position, name string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Device = cb
	}
}

// Program is called for the program name message
func Program(cb func(p Position, text string)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Program = cb
	}
}

// SMPTE is called for the smpte offset message
func SMPTE(cb func(p Position, hour, minute, second, frame, fractionalFrame byte)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.SMPTE = cb
	}
}

// SequencerData is called for the sequencer specific message
func SequencerData(cb func(p Position, data []byte)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.SequencerData = cb
	}
}

// Undefined is called for the undefined meta message
func Undefined(cb func(p Position, typ byte, data []byte)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Undefined = cb
	}
}

// Channel is called for the deprecated MIDI channel message
func DeprecatedChannel(cb func(p Position, channel uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Deprecated.Channel = cb
	}
}

// Port is called for the deprecated MIDI port message
func DeprecatedPort(cb func(p Position, port uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Meta.Deprecated.Port = cb
	}
}

// Channel provides callbacks for channel messages
// They may occur in SMF files and in live MIDI.
// For live MIDI *Position is nil.

// NoteOn is just called for noteon messages with a velocity > 0.
// Noteon messages with velocity == 0 will trigger NoteOff with a velocity of 0.
func NoteOn(cb func(p *Position, channel, key, velocity uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.NoteOn = cb
	}
}

// NoteOff is called for noteoff messages (then the given velocity is passed)
// and for noteon messages of velocity 0 (then velocity is 0).
func NoteOff(cb func(p *Position, channel, key, velocity uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.NoteOff = cb
	}
}

// Pitchbend is called for pitch bend messages
func Pitchbend(cb func(p *Position, channel uint8, value int16)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.Pitchbend = cb
	}
}

// ProgramChange is called for program change messages. Program numbers start with 0.
func ProgramChange(cb func(p *Position, channel, program uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ProgramChange = cb
	}
}

// Aftertouch is called for aftertouch messages  (aka "channel pressure")
func Aftertouch(cb func(p *Position, channel, pressure uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.Aftertouch = cb
	}
}

// PolyAftertouch is called for polyphonic aftertouch messages (aka "key pressure").
func PolyAftertouch(cb func(p *Position, channel, key, pressure uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.PolyAftertouch = cb
	}
}

// ControlChange deals with control change messages

// Each is called for every control change message
// If RPN or NRPN callbacks are defined, the corresponding control change messages will not
// be passed to each and the corrsponding RPN/NRPN callback are called.
func ControlChange(cb func(p *Position, channel, controller, value uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.Each = cb
	}
}

// RPN deals with Registered Program Numbers (RPN) and their values.
// If the callbacks are set, the corresponding control change messages will not be passed of ControlChange.Each.

// MSB is called, when the MSB of a RPN arrives
func RpnMSB(cb func(p *Position, channel, typ1, typ2, msbVal uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.RPN.MSB = cb
	}
}

// LSB is called, when the MSB of a RPN arrives
func RpnLSB(cb func(p *Position, channel, typ1, typ2, lsbVal uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.RPN.LSB = cb
	}
}

// Increment is called, when the increment of a RPN arrives
func RpnIncrement(cb func(p *Position, channel, typ1, typ2 uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.RPN.Increment = cb
	}
}

// Decrement is called, when the decrement of a RPN arrives
func RpnDecrement(cb func(p *Position, channel, typ1, typ2 uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.RPN.Decrement = cb
	}
}

// Reset is called, when the reset or null RPN arrives
func RpnReset(cb func(p *Position, channel uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.RPN.Reset = cb
	}
}

// NRPN deals with Non-Registered Program Numbers (NRPN) and their values.
// If the callbacks are set, the corresponding control change messages will not be passed of ControlChange.Each.

// MSB is called, when the MSB of a NRPN arrives
func NrpnMSB(cb func(p *Position, channel uint8, typ1, typ2, msbVal uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.NRPN.MSB = cb
	}
}

// LSB is called, when the LSB of a NRPN arrives
func NrpnLSB(cb func(p *Position, channel uint8, typ1, typ2, msbVal uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.NRPN.LSB = cb
	}
}

// Increment is called, when the increment of a NRPN arrives
func NrpnIncrement(cb func(p *Position, channel, typ1, typ2 uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.NRPN.Increment = cb
	}
}

// Decrement is called, when the decrement of a NRPN arrives
func NrpnDecrement(cb func(p *Position, channel, typ1, typ2 uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.NRPN.Decrement = cb
	}
}

// Reset is called, when the reset or null NRPN arrives
func NrpnReset(cb func(p *Position, channel uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.Channel.ControlChange.NRPN.Reset = cb
	}
}

// Realtime provides callbacks for realtime messages.
// They are only used with "live" MIDI

// Clock is called for a clock message
func RTClock(cb func()) func(r *Reader) {
	return func(r *Reader) {
		r.message.Realtime.Clock = cb
	}
}

// Tick is called for a tick message
func RTTick(cb func()) func(r *Reader) {
	return func(r *Reader) {
		r.message.Realtime.Tick = cb
	}
}

// Activesense is called for a active sense message
func RTActivesense(cb func()) func(r *Reader) {
	return func(r *Reader) {
		r.message.Realtime.Activesense = cb
	}
}

// Start is called for a start message
func RTStart(cb func()) func(r *Reader) {
	return func(r *Reader) {
		r.message.Realtime.Start = cb
	}
}

// Stop is called for a stop message
func RTStop(cb func()) func(r *Reader) {
	return func(r *Reader) {
		r.message.Realtime.Stop = cb
	}
}

// Continue is called for a continue message
func RTContinue(cb func()) func(r *Reader) {
	return func(r *Reader) {
		r.message.Realtime.Continue = cb
	}
}

// Reset is called for a reset message
func RTReset(cb func()) func(r *Reader) {
	return func(r *Reader) {
		r.message.Realtime.Reset = cb
	}
}

// SysCommon provides callbacks for system common messages.
// They are only used with "live" MIDI

// Tune is called for a tune request message
func Tune(cb func()) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysCommon.Tune = cb
	}
}

// SongSelect is called for a song select message
func SongSelect(cb func(num uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysCommon.SongSelect = cb
	}
}

// SPP is called for a song position pointer message
func SPP(cb func(pos uint16)) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysCommon.SPP = cb
	}
}

// MTC is called for a MIDI timing code message
func MTC(cb func(frame uint8)) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysCommon.MTC = cb
	}
}

// SysEx provides callbacks for system exclusive messages.
// They may occur in SMF files and in live MIDI.
// For live MIDI *Position is nil.

// Complete is called for a complete system exclusive message
func SysEx(cb func(p *Position, data []byte)) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysEx.Complete = cb
	}
}

// Start is called for a starting system exclusive message
func SysExStart(cb func(p *Position, data []byte)) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysEx.Start = cb
	}
}

// Continue is called for a continuing system exclusive message
func SysExContinue(cb func(p *Position, data []byte)) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysEx.Continue = cb
	}
}

// End is called for an ending system exclusive message
func SysExEnd(cb func(p *Position, data []byte)) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysEx.End = cb
	}
}

// Escape is called for an escaping system exclusive message
func SysExEscape(cb func(p *Position, data []byte)) func(r *Reader) {
	return func(r *Reader) {
		r.message.SysEx.Escape = cb
	}
}
