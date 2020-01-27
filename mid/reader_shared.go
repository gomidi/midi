package mid

import (
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/midi/midimessage/syscommon"
	"gitlab.com/gomidi/midi/midimessage/sysex"
	"gitlab.com/gomidi/midi/smf"
)

func (r *Reader) reset() {
	r.tempoChanges = []tempoChange{tempoChange{0, 120}}

	for c := 0; c < 16; c++ {
		r.channelRPN_NRPN[c] = [4]uint8{0, 0, 0, 0}
	}
}

func (r *Reader) saveTempoChange(pos Position, bpm float64) {
	r.tempoChanges = append(r.tempoChanges, tempoChange{pos.AbsoluteTicks, bpm})
}

// TimeAt returns the time.Duration at the given absolute position counted
// from the beginning of the file, respecting all the tempo changes in between.
// If the time format is not of type smf.MetricTicks, nil is returned.
func (r *Reader) TimeAt(absTicks uint64) *time.Duration {
	if r.resolution == 0 {
		return nil
	}

	var tc = tempoChange{0, 120}
	var lastTick uint64
	var lastDur time.Duration
	for _, t := range r.tempoChanges {
		if t.absTicks >= absTicks {
			// println("stopping")
			break
		}
		// println("pre", "lastDur", lastDur, "lastTick", lastTick, "bpm", tc.bpm)
		lastDur += calcDeltaTime(r.resolution, uint32(t.absTicks-lastTick), tc.bpm)
		tc = t
		lastTick = t.absTicks
	}
	result := lastDur + calcDeltaTime(r.resolution, uint32(absTicks-lastTick), tc.bpm)
	return &result
}

// log does the logging
func (r *Reader) log(m midi.Message) {
	if r.pos != nil {
		r.logger.Printf("#%v [%v d:%v] %s\n", r.pos.Track, r.pos.AbsoluteTicks, r.pos.DeltaTicks, m)
	} else {
		r.logger.Printf("%s\n", m)
	}
}

// dispatch dispatches the messages from the midi.Reader (which might be an smf reader)
// for realtime reading, the passed *SMFPosition is nil
func (r *Reader) dispatch() (err error) {
	for {
		err = r.dispatchMessageFromReader()
		if err != nil {
			return
		}
	}
}

func (r *Reader) _RPN_NRPN_Reset(ch uint8, isRPN bool) {
	// reset tracking on this channel
	r.channelRPN_NRPN[ch] = [4]uint8{0, 0, 0, 0}

	if isRPN {
		if r.Msg.Channel.ControlChange.RPN.Reset != nil {
			r.Msg.Channel.ControlChange.RPN.Reset(r.pos, ch)
			return
		}
		if r.Msg.Channel.ControlChange.RPN.MSB != nil {
			r.Msg.Channel.ControlChange.RPN.MSB(r.pos, ch, 127, 127, 0)
		}

		return
	}

	if r.Msg.Channel.ControlChange.NRPN.Reset != nil {
		r.Msg.Channel.ControlChange.NRPN.Reset(r.pos, ch)
		return
	}
	if r.Msg.Channel.ControlChange.NRPN.MSB != nil {
		r.Msg.Channel.ControlChange.NRPN.MSB(r.pos, ch, 127, 127, 0)
	}

}

func (r *Reader) sendAsCC(ch, cc, val uint8) error {
	if r.Msg.Channel.ControlChange.Each != nil {
		r.Msg.Channel.ControlChange.Each(r.pos, ch, cc, val)
	}
	return nil
}

func (r *Reader) hasRPNCallback() bool {
	return !(r.Msg.Channel.ControlChange.RPN.MSB == nil && r.Msg.Channel.ControlChange.RPN.LSB == nil)
}

func (r *Reader) hasNRPNCallback() bool {
	return !(r.Msg.Channel.ControlChange.NRPN.MSB == nil && r.Msg.Channel.ControlChange.NRPN.LSB == nil)
}

func (r *Reader) hasNoRPNorNRPNCallback() bool {
	return !r.hasRPNCallback() && !r.hasNRPNCallback()
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

	if r.Msg.Each != nil {
		r.Msg.Each(r.pos, m)
	}

	switch msg := m.(type) {

	// most common event, should be exact
	case channel.NoteOn:
		if r.Msg.Channel.NoteOn != nil {
			r.Msg.Channel.NoteOn(r.pos, msg.Channel(), msg.Key(), msg.Velocity())
		}

	// proably second most common
	case channel.NoteOff:
		if r.Msg.Channel.NoteOff != nil {
			r.Msg.Channel.NoteOff(r.pos, msg.Channel(), msg.Key(), 0)
		}

	case channel.NoteOffVelocity:
		if r.Msg.Channel.NoteOff != nil {
			r.Msg.Channel.NoteOff(r.pos, msg.Channel(), msg.Key(), msg.Velocity())
		}

	// if send there often are a lot of them
	case channel.Pitchbend:
		if r.Msg.Channel.Pitchbend != nil {
			r.Msg.Channel.Pitchbend(r.pos, msg.Channel(), msg.Value())
		}

	case channel.PolyAftertouch:
		if r.Msg.Channel.PolyAftertouch != nil {
			r.Msg.Channel.PolyAftertouch(r.pos, msg.Channel(), msg.Key(), msg.Pressure())
		}

	case channel.Aftertouch:
		if r.Msg.Channel.Aftertouch != nil {
			r.Msg.Channel.Aftertouch(r.pos, msg.Channel(), msg.Pressure())
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
				if r.Msg.Channel.ControlChange.RPN.MSB != nil {
					r.Msg.Channel.ControlChange.RPN.MSB(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3],
						val)
				}
				return

			// is a valid NRPN
			case r.channelRPN_NRPN[ch][0] == 99 && r.channelRPN_NRPN[ch][1] == 98:
				if r.Msg.Channel.ControlChange.NRPN.MSB != nil {
					r.Msg.Channel.ControlChange.NRPN.MSB(
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
				if r.Msg.Channel.ControlChange.RPN.LSB != nil {
					r.Msg.Channel.ControlChange.RPN.LSB(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3],
						val)
				}
				return

			// is a valid NRPN
			case r.channelRPN_NRPN[ch][0] == 99 && r.channelRPN_NRPN[ch][1] == 98:
				if r.Msg.Channel.ControlChange.NRPN.LSB != nil {
					r.Msg.Channel.ControlChange.NRPN.LSB(
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
			if r.Msg.Channel.ControlChange.RPN.Increment == nil && r.Msg.Channel.ControlChange.NRPN.Increment == nil {
				return r.sendAsCC(ch, cc, val)
			}
			switch {

			// is a valid RPN
			case r.channelRPN_NRPN[ch][0] == 101 && r.channelRPN_NRPN[ch][1] == 100:
				if r.Msg.Channel.ControlChange.RPN.Increment != nil {
					r.Msg.Channel.ControlChange.RPN.Increment(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3])
				}
				return

			// is a valid NRPN
			case r.channelRPN_NRPN[ch][0] == 99 && r.channelRPN_NRPN[ch][1] == 98:
				if r.Msg.Channel.ControlChange.NRPN.Increment != nil {
					r.Msg.Channel.ControlChange.NRPN.Increment(
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
			if r.Msg.Channel.ControlChange.RPN.Decrement == nil && r.Msg.Channel.ControlChange.NRPN.Decrement == nil {
				return r.sendAsCC(ch, cc, val)
			}
			switch {

			// is a valid RPN
			case r.channelRPN_NRPN[ch][0] == 101 && r.channelRPN_NRPN[ch][1] == 100:
				if r.Msg.Channel.ControlChange.RPN.Decrement != nil {
					r.Msg.Channel.ControlChange.RPN.Decrement(
						r.pos,
						ch,
						r.channelRPN_NRPN[ch][2],
						r.channelRPN_NRPN[ch][3])
				}
				return

			// is a valid NRPN
			case r.channelRPN_NRPN[ch][0] == 99 && r.channelRPN_NRPN[ch][1] == 98:
				if r.Msg.Channel.ControlChange.NRPN.Decrement != nil {
					r.Msg.Channel.ControlChange.NRPN.Decrement(
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
		if r.Msg.Meta.SMPTE != nil {
			r.Msg.Meta.SMPTE(*r.pos, msg.Hour, msg.Minute, msg.Second, msg.Frame, msg.FractionalFrame)
		}

	case meta.Tempo:
		r.saveTempoChange(*r.pos, msg.FractionalBPM())
		if r.Msg.Meta.TempoBPM != nil {
			r.Msg.Meta.TempoBPM(*r.pos, msg.FractionalBPM())
		}

	case meta.TimeSig:
		if r.Msg.Meta.TimeSig != nil {
			r.Msg.Meta.TimeSig(*r.pos, msg.Numerator, msg.Denominator)
		}

		// may be for karaoke we need to be fast
	case meta.Lyric:
		if r.Msg.Meta.Lyric != nil {
			r.Msg.Meta.Lyric(*r.pos, msg.Text())
		}

	// may be useful to synchronize by sequence number
	case meta.SequenceNo:
		if r.Msg.Meta.SequenceNo != nil {
			r.Msg.Meta.SequenceNo(*r.pos, msg.Number())
		}

	case meta.Marker:
		if r.Msg.Meta.Marker != nil {
			r.Msg.Meta.Marker(*r.pos, msg.Text())
		}

	case meta.Cuepoint:
		if r.Msg.Meta.Cuepoint != nil {
			r.Msg.Meta.Cuepoint(*r.pos, msg.Text())
		}

	case meta.Program:
		if r.Msg.Meta.Program != nil {
			r.Msg.Meta.Program(*r.pos, msg.Text())
		}

	case meta.SequencerData:
		if r.Msg.Meta.SequencerData != nil {
			r.Msg.Meta.SequencerData(*r.pos, msg.Data())
		}

	case sysex.SysEx:
		if r.Msg.SysEx.Complete != nil {
			r.Msg.SysEx.Complete(r.pos, msg.Data())
		}

	case sysex.Start:
		if r.Msg.SysEx.Start != nil {
			r.Msg.SysEx.Start(r.pos, msg.Data())
		}

	case sysex.End:
		if r.Msg.SysEx.End != nil {
			r.Msg.SysEx.End(r.pos, msg.Data())
		}

	case sysex.Continue:
		if r.Msg.SysEx.Continue != nil {
			r.Msg.SysEx.Continue(r.pos, msg.Data())
		}

	case sysex.Escape:
		if r.Msg.SysEx.Escape != nil {
			r.Msg.SysEx.Escape(r.pos, msg.Data())
		}

	// this usually takes some time
	case channel.ProgramChange:
		if r.Msg.Channel.ProgramChange != nil {
			r.Msg.Channel.ProgramChange(r.pos, msg.Channel(), msg.Program())
		}

	// the rest is not that interesting for performance
	case meta.Key:
		if r.Msg.Meta.Key != nil {
			r.Msg.Meta.Key(*r.pos, msg.Key, msg.IsMajor, msg.Num, msg.IsFlat)
		}

	case meta.TrackSequenceName:
		if r.Msg.Meta.TrackSequenceName != nil {
			r.Msg.Meta.TrackSequenceName(*r.pos, msg.Text())
		}

	case meta.Instrument:
		if r.Msg.Meta.Instrument != nil {
			r.Msg.Meta.Instrument(*r.pos, msg.Text())
		}

	case meta.Channel:
		if r.Msg.Meta.Deprecated.Channel != nil {
			r.Msg.Meta.Deprecated.Channel(*r.pos, msg.Number())
		}

	case meta.Port:
		if r.Msg.Meta.Deprecated.Port != nil {
			r.Msg.Meta.Deprecated.Port(*r.pos, msg.Number())
		}

	case meta.Text:
		if r.Msg.Meta.Text != nil {
			r.Msg.Meta.Text(*r.pos, msg.Text())
		}

	case syscommon.SongSelect:
		if r.Msg.SysCommon.SongSelect != nil {
			r.Msg.SysCommon.SongSelect(msg.Number())
		}

	case syscommon.SPP:
		if r.Msg.SysCommon.SPP != nil {
			r.Msg.SysCommon.SPP(msg.Number())
		}

	case syscommon.MTC:
		if r.Msg.SysCommon.MTC != nil {
			r.Msg.SysCommon.MTC(msg.QuarterFrame())
		}

	case meta.Copyright:
		if r.Msg.Meta.Copyright != nil {
			r.Msg.Meta.Copyright(*r.pos, msg.Text())
		}

	case meta.Device:
		if r.Msg.Meta.Device != nil {
			r.Msg.Meta.Device(*r.pos, msg.Text())
		}

	//case meta.Undefined, syscommon.Undefined4, syscommon.Undefined5:
	case meta.Undefined:
		if r.Msg.Meta.Undefined != nil {
			r.Msg.Meta.Undefined(*r.pos, msg.Typ, msg.Data)
		}

	default:
		switch m {
		case syscommon.Tune:
			if r.Msg.SysCommon.Tune != nil {
				r.Msg.SysCommon.Tune()
			}
		case meta.EndOfTrack:
			if _, ok := r.reader.(smf.Reader); ok && r.pos != nil {
				r.pos.DeltaTicks = 0
				r.pos.AbsoluteTicks = 0
			}
			if r.Msg.Meta.EndOfTrack != nil {
				r.Msg.Meta.EndOfTrack(*r.pos)
			}
		default:

			if r.Msg.Unknown != nil {
				r.Msg.Unknown(r.pos, m)
			}

		}

	}
	return nil
}

// dispatchMessageFromReader dispatches a single message from the midi.Reader (which might be an smf reader)
// for realtime reading, the passed *SMFPosition is nil
func (r *Reader) dispatchMessageFromReader() (err error) {
	var m midi.Message
	m, err = r.reader.Read()
	if err != nil {
		return
	}

	return r.dispatchMessage(m)
}

type tempoChange struct {
	absTicks uint64
	bpm      float64
}

func calcDeltaTime(mt smf.MetricTicks, deltaTicks uint32, bpm float64) time.Duration {
	return mt.FractionalDuration(bpm, deltaTicks)
}
