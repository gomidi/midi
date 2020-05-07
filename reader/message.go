package reader

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/midi/midimessage/syscommon"
	"gitlab.com/gomidi/midi/midimessage/sysex"
	"gitlab.com/gomidi/midi/smf"
)

// Msg provides callbacks for MIDI messages
type message struct {

	// Each is called for every MIDI message in addition to the other callbacks.
	Each func(*Position, midi.Message)

	// Unknown is called for unknown messages
	Unknown func(p *Position, msg midi.Message)

	// Meta provides callbacks for meta messages (only in SMF files)
	Meta struct {

		// Copyright is called for the copyright message
		Copyright func(p Position, text string)

		// TempoBPM is called for the tempo (change) message, BPM is fractional
		TempoBPM func(p Position, bpm float64)

		// TimeSigis called for the time signature (change) message
		TimeSig func(p Position, num, denom uint8)

		// Key is called for the key signature (change) message
		Key func(p Position, key uint8, ismajor bool, num_accidentals uint8, accidentals_are_flat bool)

		// Instrument is called for the instrument (name) message
		Instrument func(p Position, name string)

		// TrackSequenceName is called for the sequence / track name message
		// If in a format 0 track, or the first track in a format 1 file, the name of the sequence. Otherwise, the name of the track.
		TrackSequenceName func(p Position, name string)

		// SequenceNo is called for the sequence number message
		SequenceNo func(p Position, number uint16)

		// Marker is called for the marker message
		Marker func(p Position, text string)

		// Cuepoint is called for the cuepoint message
		Cuepoint func(p Position, text string)

		// Text is called for the text message
		Text func(p Position, text string)

		// Lyric is called for the lyric message
		Lyric func(p Position, text string)

		// EndOfTrack is called for the end of a track message
		EndOfTrack func(p Position)

		// Device is called for the device port message
		Device func(p Position, name string)

		// Program is called for the program name message
		Program func(p Position, text string)

		// SMPTE is called for the smpte offset message
		SMPTE func(p Position, hour, minute, second, frame, fractionalFrame byte)

		// SequencerData is called for the sequencer specific message
		SequencerData func(p Position, data []byte)

		// Undefined is called for the undefined meta message
		Undefined func(p Position, typ byte, data []byte)

		Deprecated struct {
			// Channel is called for the deprecated MIDI channel message
			Channel func(p Position, channel uint8)

			// Port is called for the deprecated MIDI port message
			Port func(p Position, port uint8)
		}
	}

	// Channel provides callbacks for channel messages
	// They may occur in SMF files and in live MIDI.
	// For live MIDI *Position is nil.
	Channel struct {

		// NoteOn is just called for noteon messages with a velocity > 0.
		// Noteon messages with velocity == 0 will trigger NoteOff with a velocity of 0.
		NoteOn func(p *Position, channel, key, velocity uint8)

		// NoteOff is called for noteoff messages (then the given velocity is passed)
		// and for noteon messages of velocity 0 (then velocity is 0).
		NoteOff func(p *Position, channel, key, velocity uint8)

		// Pitchbend is called for pitch bend messages
		Pitchbend func(p *Position, channel uint8, value int16)

		// ProgramChange is called for program change messages. Program numbers start with 0.
		ProgramChange func(p *Position, channel, program uint8)

		// Aftertouch is called for aftertouch messages  (aka "channel pressure")
		Aftertouch func(p *Position, channel, pressure uint8)

		// PolyAftertouch is called for polyphonic aftertouch messages (aka "key pressure").
		PolyAftertouch func(p *Position, channel, key, pressure uint8)

		// ControlChange deals with control change messages
		ControlChange struct {

			// Each is called for every control change message
			// If RPN or NRPN callbacks are defined, the corresponding control change messages will not
			// be passed to each and the corrsponding RPN/NRPN callback are called.
			Each func(p *Position, channel, controller, value uint8)

			// RPN deals with Registered Program Numbers (RPN) and their values.
			// If the callbacks are set, the corresponding control change messages will not be passed of ControlChange.Each.
			RPN struct {

				// MSB is called, when the MSB of a RPN arrives
				MSB func(p *Position, channel, typ1, typ2, msbVal uint8)

				// LSB is called, when the MSB of a RPN arrives
				LSB func(p *Position, channel, typ1, typ2, lsbVal uint8)

				// Increment is called, when the increment of a RPN arrives
				Increment func(p *Position, channel, typ1, typ2 uint8)

				// Decrement is called, when the decrement of a RPN arrives
				Decrement func(p *Position, channel, typ1, typ2 uint8)

				// Reset is called, when the reset or null RPN arrives
				Reset func(p *Position, channel uint8)
			}

			// NRPN deals with Non-Registered Program Numbers (NRPN) and their values.
			// If the callbacks are set, the corresponding control change messages will not be passed of ControlChange.Each.
			NRPN struct {

				// MSB is called, when the MSB of a NRPN arrives
				MSB func(p *Position, channel uint8, typ1, typ2, msbVal uint8)

				// LSB is called, when the LSB of a NRPN arrives
				LSB func(p *Position, channel uint8, typ1, typ2, msbVal uint8)

				// Increment is called, when the increment of a NRPN arrives
				Increment func(p *Position, channel, typ1, typ2 uint8)

				// Decrement is called, when the decrement of a NRPN arrives
				Decrement func(p *Position, channel, typ1, typ2 uint8)

				// Reset is called, when the reset or null NRPN arrives
				Reset func(p *Position, channel uint8)
			}
		}
	}

	// Realtime provides callbacks for realtime messages.
	// They are only used with "live" MIDI
	Realtime struct {

		// Clock is called for a clock message
		Clock func()

		// Tick is called for a tick message
		Tick func()

		// Activesense is called for a active sense message
		Activesense func()

		// Start is called for a start message
		Start func()

		// Stop is called for a stop message
		Stop func()

		// Continue is called for a continue message
		Continue func()

		// Reset is called for a reset message
		Reset func()
	}

	// SysCommon provides callbacks for system common messages.
	// They are only used with "live" MIDI
	SysCommon struct {

		// Tune is called for a tune request message
		Tune func()

		// SongSelect is called for a song select message
		SongSelect func(num uint8)

		// SPP is called for a song position pointer message
		SPP func(pos uint16)

		// MTC is called for a MIDI timing code message
		MTC func(frame uint8)
	}

	// SysEx provides callbacks for system exclusive messages.
	// They may occur in SMF files and in live MIDI.
	// For live MIDI *Position is nil.
	SysEx struct {

		// Complete is called for a complete system exclusive message
		Complete func(p *Position, data []byte)

		// Start is called for a starting system exclusive message
		Start func(p *Position, data []byte)

		// Continue is called for a continuing system exclusive message
		Continue func(p *Position, data []byte)

		// End is called for an ending system exclusive message
		End func(p *Position, data []byte)

		// Escape is called for an escaping system exclusive message
		Escape func(p *Position, data []byte)
	}
}

func (r *Reader) dispatchMessage(m midi.Message) (err error) {
	if frd, ok := r.reader.(smf.Reader); ok && r.pos != nil {
		r.pos.DeltaTicks = frd.Delta()
		r.pos.AbsoluteTicks += uint64(r.pos.DeltaTicks)
		r.pos.Track = frd.Track()
	}

	if r.logger != nil {
		r.log(m)
	}

	if r.message.Each != nil {
		r.message.Each(r.pos, m)
	}

	switch msg := m.(type) {

	// most common event, should be exact
	case channel.NoteOn:
		if r.message.Channel.NoteOn != nil {
			r.message.Channel.NoteOn(r.pos, msg.Channel(), msg.Key(), msg.Velocity())
		}

	// proably second most common
	case channel.NoteOff:
		if r.message.Channel.NoteOff != nil {
			r.message.Channel.NoteOff(r.pos, msg.Channel(), msg.Key(), 0)
		}

	case channel.NoteOffVelocity:
		if r.message.Channel.NoteOff != nil {
			r.message.Channel.NoteOff(r.pos, msg.Channel(), msg.Key(), msg.Velocity())
		}

	// if send there often are a lot of them
	case channel.Pitchbend:
		if r.message.Channel.Pitchbend != nil {
			r.message.Channel.Pitchbend(r.pos, msg.Channel(), msg.Value())
		}

	case channel.PolyAftertouch:
		if r.message.Channel.PolyAftertouch != nil {
			r.message.Channel.PolyAftertouch(r.pos, msg.Channel(), msg.Key(), msg.Pressure())
		}

	case channel.Aftertouch:
		if r.message.Channel.Aftertouch != nil {
			r.message.Channel.Aftertouch(r.pos, msg.Channel(), msg.Pressure())
		}

	case channel.ControlChange:
		var (
			ch  = msg.Channel()
			val = msg.Value()
			cc  = msg.Controller()
		)

		switch cc {

		/*
			Ok, lets explain the reasoning behind this confusing RPN/NRPN handling a bit.
			There are the following observations:
				- a channel can either have a RPN message or a NRPN message at a point in time
				- the identifiers are sent via CC101 + CC100 for RPN and CC99 + CC98 for NRPN
			    - the order of the identifier CC messages may vary in reality
				- the identifiers are sent before the value
				- the MSB is sent via CC6
				- the LSB is sent via CC38

			RPN and NRPN are never mixed at the same time on the same channel.
			We want to always send complete valid RPN/NRPN messages to the callbacks.
			For this to happen, each identifier is cached and when the MSB arrives and both identifiers are there,
			the callback is called. If any of the conditions are not met, the callback is not called.
		*/

		// first identifier of a RPN/NRPN message
		case 101, 99:
			if (cc == 101 && !r.hasRPNCallback()) ||
				(cc == 99 && !r.hasNRPNCallback()) {
				return r.sendAsCC(ch, cc, val)
			}

			// RPN reset (127,127)
			if val+r.channelRPN_NRPN[ch][3] == 2*127 {
				r._RPN_NRPN_Reset(ch, cc == 101)
			} else {
				// register first ident cc
				r.channelRPN_NRPN[ch][0] = cc
				// track the first ident value
				r.channelRPN_NRPN[ch][2] = val
			}

		// second identifier of a RPN/NRPN message
		case 100, 98:
			if (cc == 100 && !r.hasRPNCallback()) ||
				(cc == 98 && !r.hasNRPNCallback()) {
				return r.sendAsCC(ch, cc, val)
			}

			// RPN reset (127,127)
			if val+r.channelRPN_NRPN[ch][2] == 2*127 {
				r._RPN_NRPN_Reset(ch, cc == 100)
			} else {
				// register second ident cc
				r.channelRPN_NRPN[ch][1] = cc
				// track the second ident value
				r.channelRPN_NRPN[ch][3] = val
			}

		// the data entry controller
		case 6:
			if r.hasNoRPNorNRPNCallback() {
				//println("early return on cc6")
				return r.sendAsCC(ch, cc, val)
			}
			switch {

			// is a valid RPN
			case r.channelRPN_NRPN[ch][0] == 101 && r.channelRPN_NRPN[ch][1] == 100:
				if r.message.Channel.ControlChange.RPN.MSB != nil {
					r.message.Channel.ControlChange.RPN.MSB(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3],
						val)
				}
				return

			// is a valid NRPN
			case r.channelRPN_NRPN[ch][0] == 99 && r.channelRPN_NRPN[ch][1] == 98:
				if r.message.Channel.ControlChange.NRPN.MSB != nil {
					r.message.Channel.ControlChange.NRPN.MSB(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3],
						val)
				}
				return

			// is no valid RPN/NRPN, send as controller change
			default:
				//				println("invalid RPN/NRPN on cc6")
				return r.sendAsCC(ch, cc, val)
			}

		// the lsb
		case 38:
			if r.hasNoRPNorNRPNCallback() {
				return r.sendAsCC(ch, cc, val)
			}

			switch {

			// is a valid RPN
			case r.channelRPN_NRPN[ch][0] == 101 && r.channelRPN_NRPN[ch][1] == 100:
				if r.message.Channel.ControlChange.RPN.LSB != nil {
					r.message.Channel.ControlChange.RPN.LSB(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3],
						val)
				}
				return

			// is a valid NRPN
			case r.channelRPN_NRPN[ch][0] == 99 && r.channelRPN_NRPN[ch][1] == 98:
				if r.message.Channel.ControlChange.NRPN.LSB != nil {
					r.message.Channel.ControlChange.NRPN.LSB(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3],
						val)
				}
				return

			// is no valid RPN/NRPN, send as controller change
			default:
				return r.sendAsCC(ch, cc, val)
			}

		// the increment
		case 96:
			if r.message.Channel.ControlChange.RPN.Increment == nil && r.message.Channel.ControlChange.NRPN.Increment == nil {
				return r.sendAsCC(ch, cc, val)
			}
			switch {

			// is a valid RPN
			case r.channelRPN_NRPN[ch][0] == 101 && r.channelRPN_NRPN[ch][1] == 100:
				if r.message.Channel.ControlChange.RPN.Increment != nil {
					r.message.Channel.ControlChange.RPN.Increment(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3])
				}
				return

			// is a valid NRPN
			case r.channelRPN_NRPN[ch][0] == 99 && r.channelRPN_NRPN[ch][1] == 98:
				if r.message.Channel.ControlChange.NRPN.Increment != nil {
					r.message.Channel.ControlChange.NRPN.Increment(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3])
				}
				return

			// is no valid RPN/NRPN, send as controller change
			default:
				return r.sendAsCC(ch, cc, val)
			}

		// the decrement
		case 97:
			if r.message.Channel.ControlChange.RPN.Decrement == nil && r.message.Channel.ControlChange.NRPN.Decrement == nil {
				return r.sendAsCC(ch, cc, val)
			}
			switch {

			// is a valid RPN
			case r.channelRPN_NRPN[ch][0] == 101 && r.channelRPN_NRPN[ch][1] == 100:
				if r.message.Channel.ControlChange.RPN.Decrement != nil {
					r.message.Channel.ControlChange.RPN.Decrement(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3])
				}
				return

			// is a valid NRPN
			case r.channelRPN_NRPN[ch][0] == 99 && r.channelRPN_NRPN[ch][1] == 98:
				if r.message.Channel.ControlChange.NRPN.Decrement != nil {
					r.message.Channel.ControlChange.NRPN.Decrement(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3])
				}
				return

			// is no valid RPN/NRPN, send as controller change
			default:
				return r.sendAsCC(ch, cc, val)
			}

		default:
			return r.sendAsCC(ch, cc, val)
		}

	case meta.SMPTE:
		if r.message.Meta.SMPTE != nil {
			r.message.Meta.SMPTE(*r.pos, msg.Hour, msg.Minute, msg.Second, msg.Frame, msg.FractionalFrame)
		}

	case meta.Tempo:
		r.saveTempoChange(*r.pos, msg.FractionalBPM())
		if r.message.Meta.TempoBPM != nil {
			r.message.Meta.TempoBPM(*r.pos, msg.FractionalBPM())
		}

	case meta.TimeSig:
		if r.message.Meta.TimeSig != nil {
			r.message.Meta.TimeSig(*r.pos, msg.Numerator, msg.Denominator)
		}

		// may be for karaoke we need to be fast
	case meta.Lyric:
		if r.message.Meta.Lyric != nil {
			r.message.Meta.Lyric(*r.pos, msg.Text())
		}

	// may be useful to synchronize by sequence number
	case meta.SequenceNo:
		if r.message.Meta.SequenceNo != nil {
			r.message.Meta.SequenceNo(*r.pos, msg.Number())
		}

	case meta.Marker:
		if r.message.Meta.Marker != nil {
			r.message.Meta.Marker(*r.pos, msg.Text())
		}

	case meta.Cuepoint:
		if r.message.Meta.Cuepoint != nil {
			r.message.Meta.Cuepoint(*r.pos, msg.Text())
		}

	case meta.Program:
		if r.message.Meta.Program != nil {
			r.message.Meta.Program(*r.pos, msg.Text())
		}

	case meta.SequencerData:
		if r.message.Meta.SequencerData != nil {
			r.message.Meta.SequencerData(*r.pos, msg.Data())
		}

	case sysex.SysEx:
		if r.message.SysEx.Complete != nil {
			r.message.SysEx.Complete(r.pos, msg.Data())
		}

	case sysex.Start:
		if r.message.SysEx.Start != nil {
			r.message.SysEx.Start(r.pos, msg.Data())
		}

	case sysex.End:
		if r.message.SysEx.End != nil {
			r.message.SysEx.End(r.pos, msg.Data())
		}

	case sysex.Continue:
		if r.message.SysEx.Continue != nil {
			r.message.SysEx.Continue(r.pos, msg.Data())
		}

	case sysex.Escape:
		if r.message.SysEx.Escape != nil {
			r.message.SysEx.Escape(r.pos, msg.Data())
		}

	// this usually takes some time
	case channel.ProgramChange:
		if r.message.Channel.ProgramChange != nil {
			r.message.Channel.ProgramChange(r.pos, msg.Channel(), msg.Program())
		}

	// the rest is not that interesting for performance
	case meta.Key:
		if r.message.Meta.Key != nil {
			r.message.Meta.Key(*r.pos, msg.Key, msg.IsMajor, msg.Num, msg.IsFlat)
		}

	case meta.TrackSequenceName:
		if r.message.Meta.TrackSequenceName != nil {
			r.message.Meta.TrackSequenceName(*r.pos, msg.Text())
		}

	case meta.Instrument:
		if r.message.Meta.Instrument != nil {
			r.message.Meta.Instrument(*r.pos, msg.Text())
		}

	case meta.Channel:
		if r.message.Meta.Deprecated.Channel != nil {
			r.message.Meta.Deprecated.Channel(*r.pos, msg.Number())
		}

	case meta.Port:
		if r.message.Meta.Deprecated.Port != nil {
			r.message.Meta.Deprecated.Port(*r.pos, msg.Number())
		}

	case meta.Text:
		if r.message.Meta.Text != nil {
			r.message.Meta.Text(*r.pos, msg.Text())
		}

	case syscommon.SongSelect:
		if r.message.SysCommon.SongSelect != nil {
			r.message.SysCommon.SongSelect(msg.Number())
		}

	case syscommon.SPP:
		if r.message.SysCommon.SPP != nil {
			r.message.SysCommon.SPP(msg.Number())
		}

	case syscommon.MTC:
		if r.message.SysCommon.MTC != nil {
			r.message.SysCommon.MTC(msg.QuarterFrame())
		}

	case meta.Copyright:
		if r.message.Meta.Copyright != nil {
			r.message.Meta.Copyright(*r.pos, msg.Text())
		}

	case meta.Device:
		if r.message.Meta.Device != nil {
			r.message.Meta.Device(*r.pos, msg.Text())
		}

	//case meta.Undefined, syscommon.Undefined4, syscommon.Undefined5:
	case meta.Undefined:
		if r.message.Meta.Undefined != nil {
			r.message.Meta.Undefined(*r.pos, msg.Typ, msg.Data)
		}

	default:
		switch m {
		case syscommon.Tune:
			if r.message.SysCommon.Tune != nil {
				r.message.SysCommon.Tune()
			}
		case meta.EndOfTrack:
			if _, ok := r.reader.(smf.Reader); ok && r.pos != nil {
				r.pos.DeltaTicks = 0
				r.pos.AbsoluteTicks = 0
			}
			if r.message.Meta.EndOfTrack != nil {
				r.message.Meta.EndOfTrack(*r.pos)
			}
		default:

			if r.message.Unknown != nil {
				r.message.Unknown(r.pos, m)
			}

		}

	}
	return nil
}
