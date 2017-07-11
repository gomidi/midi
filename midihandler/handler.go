package midihandler

import (
	"github.com/gomidi/midi/smf"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/midimessage/syscommon"
	"github.com/gomidi/midi/midimessage/sysex"
)

// SMFPosition is the position of the event inside a standard midi file (SMF).
type SMFPosition struct {
	// the Track number
	Track int16

	// the delta time to the previous message in the same track
	Delta uint32

	// the absolute time from the beginning of the track
	AbsTime uint64
}

// New returns a new handler
func New(opts ...Option) *Handler {
	h := &Handler{logger: logfunc(printf)}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// Handler handles the midi messages coming from an SMF file or a live stream.
//
// The messages are dispatched to the corresponding functions that are not nil.
//
// The desired functions must be attached before Handler.ReadLive or Handler.ReadSMF is called
// and they must not be changed while these methods are running.
type Handler struct {

	// callback functions for SMF (Standard MIDI File) header data
	SMFHeader func(smf.Header)

	// callback functions for MIDI messages
	Message struct {
		// is called in addition to other functions, if set.
		Each func(*SMFPosition, midi.Message)

		// undefined or unknown messages
		Unknown func(p *SMFPosition, msg midi.Message)

		// meta messages (only in SMF files)
		Meta struct {
			// SMF general settings
			Copyright     func(p SMFPosition, text string)
			Tempo         func(p SMFPosition, bpm uint32)
			TimeSignature func(p SMFPosition, num, denom uint8)
			KeySignature  func(p SMFPosition, key uint8, ismajor bool, num_accidentals uint8, accidentals_are_flat bool)

			// SMF tracks and sequence definitions
			TrackInstrument func(p SMFPosition, name string)
			Sequence        func(p SMFPosition, name string)
			SequenceNumber  func(p SMFPosition, number uint16)

			// SMF text entries
			Marker   func(p SMFPosition, text string)
			CuePoint func(p SMFPosition, text string)
			Text     func(p SMFPosition, text string)
			Lyric    func(p SMFPosition, text string)

			// SMF diverse
			EndOfTrack        func(p SMFPosition)
			DevicePort        func(p SMFPosition, name string)
			ProgramName       func(p SMFPosition, text string)
			SMPTEOffset       func(p SMFPosition, hour, minute, second, frame, fractionalFrame byte)
			SequencerSpecific func(p SMFPosition, data []byte)

			// deprecated
			MIDIChannel func(p SMFPosition, channel uint8)
			MIDIPort    func(p SMFPosition, port uint8)
		}

		// channel messages, may be in SMF files and in live data
		// for live data *SMFPosition is nil
		Channel struct {
			// NoteOn is just called for noteon messages with a velocity > 0
			// noteon messages with velocity == 0 will trigger NoteOff with a velocity of 0
			NoteOn func(p *SMFPosition, channel, key, velocity uint8)

			// NoteOff is triggered by noteoff messages (then the given velocity is passed)
			// and by noteon messages of velocity 0 (then velocity is 0)
			NoteOff              func(p *SMFPosition, channel, key, velocity uint8)
			PolyphonicAfterTouch func(p *SMFPosition, channel, key, pressure uint8)
			ControlChange        func(p *SMFPosition, channel, controller, value uint8)
			ProgramChange        func(p *SMFPosition, channel, program uint8)
			AfterTouch           func(p *SMFPosition, channel, pressure uint8)
			PitchBend            func(p *SMFPosition, channel uint8, value int16)
		}

		// realtime messages: just in live data
		Realtime struct {
			// realtime messages
			Reset       func()
			Clock       func()
			Tick        func()
			Start       func()
			Continue    func()
			Stop        func()
			ActiveSense func()
		}

		// system common messages: just in live data
		SysCommon struct {
			// system common
			TuneRequest         func()
			SongSelect          func(num uint8)
			SongPositionPointer func(pos uint16)
			MIDITimingCode      func(frame uint8)
		}

		// system exclusive, may be in SMF files and in live data
		// for live data *SMFPosition is nil
		SysEx struct {
			// system exclusive
			Complete func(p *SMFPosition, data []byte)
			Start    func(p *SMFPosition, data []byte)
			Continue func(p *SMFPosition, data []byte)
			End      func(p *SMFPosition, data []byte)
			Escape   func(p *SMFPosition, data []byte)
		}
	}

	// optional logger
	logger Logger

	pos    *SMFPosition
	errSMF error
}

// log does the logging
func (h *Handler) log(m midi.Message) {
	if h.pos != nil {
		h.logger.Printf("#%v [%v d:%v] %#v\n", h.pos.Track, h.pos.AbsTime, h.pos.Delta, m)
	} else {
		h.logger.Printf("%#v\n", m)
	}
}

// read reads the messages from the midi.Reader (which might be an smf reader
// for realtime reading, the passed *SMFPosition is nil
func (h *Handler) read(rd midi.Reader) (err error) {
	var m midi.Message

	for {
		m, err = rd.Read()
		if err != nil {
			break
		}

		if frd, ok := rd.(smf.Reader); ok && h.pos != nil {
			h.pos.Delta = frd.Delta()
			h.pos.AbsTime += uint64(h.pos.Delta)
			h.pos.Track = frd.Track()
		}

		if h.logger != nil {
			h.log(m)
		}

		if h.Message.Each != nil {
			h.Message.Each(h.pos, m)
		}

		switch msg := m.(type) {

		// most common event, should be exact
		case channel.NoteOn:
			if h.Message.Channel.NoteOn != nil {
				h.Message.Channel.NoteOn(h.pos, msg.Channel(), msg.Key(), msg.Velocity())
			}

		// proably second most common
		case channel.NoteOff:
			if h.Message.Channel.NoteOff != nil {
				h.Message.Channel.NoteOff(h.pos, msg.Channel(), msg.Key(), 0)
			}

		case channel.NoteOffVelocity:
			if h.Message.Channel.NoteOff != nil {
				h.Message.Channel.NoteOff(h.pos, msg.Channel(), msg.Key(), msg.Velocity())
			}

		// if send there often are a lot of them
		case channel.PitchBend:
			if h.Message.Channel.PitchBend != nil {
				h.Message.Channel.PitchBend(h.pos, msg.Channel(), msg.Value())
			}

		case channel.PolyphonicAfterTouch:
			if h.Message.Channel.PolyphonicAfterTouch != nil {
				h.Message.Channel.PolyphonicAfterTouch(h.pos, msg.Channel(), msg.Key(), msg.Pressure())
			}

		case channel.AfterTouch:
			if h.Message.Channel.AfterTouch != nil {
				h.Message.Channel.AfterTouch(h.pos, msg.Channel(), msg.Pressure())
			}

		case channel.ControlChange:
			if h.Message.Channel.ControlChange != nil {
				h.Message.Channel.ControlChange(h.pos, msg.Channel(), msg.Controller(), msg.Value())
			}

		case meta.SMPTEOffset:
			if h.Message.Meta.SMPTEOffset != nil {
				h.Message.Meta.SMPTEOffset(*h.pos, msg.Hour, msg.Minute, msg.Second, msg.Frame, msg.FractionalFrame)
			}

		case meta.Tempo:
			if h.Message.Meta.Tempo != nil {
				h.Message.Meta.Tempo(*h.pos, msg.BPM())
			}

		case meta.TimeSignature:
			if h.Message.Meta.TimeSignature != nil {
				h.Message.Meta.TimeSignature(*h.pos, msg.Numerator, msg.Denominator)
			}

			// may be for karaoke we need to be fast
		case meta.Lyric:
			if h.Message.Meta.Lyric != nil {
				h.Message.Meta.Lyric(*h.pos, msg.Text())
			}

		// may be useful to synchronize by sequence number
		case meta.SequenceNumber:
			if h.Message.Meta.SequenceNumber != nil {
				h.Message.Meta.SequenceNumber(*h.pos, msg.Number())
			}

		case meta.Marker:
			if h.Message.Meta.Marker != nil {
				h.Message.Meta.Marker(*h.pos, msg.Text())
			}

		case meta.CuePoint:
			if h.Message.Meta.CuePoint != nil {
				h.Message.Meta.CuePoint(*h.pos, msg.Text())
			}

		case meta.ProgramName:
			if h.Message.Meta.ProgramName != nil {
				h.Message.Meta.ProgramName(*h.pos, msg.Text())
			}

		case meta.SequencerSpecific:
			if h.Message.Meta.SequencerSpecific != nil {
				h.Message.Meta.SequencerSpecific(*h.pos, msg.Data())
			}

		case sysex.SysEx:
			if h.Message.SysEx.Complete != nil {
				h.Message.SysEx.Complete(h.pos, msg.Data())
			}

		case sysex.Start:
			if h.Message.SysEx.Start != nil {
				h.Message.SysEx.Start(h.pos, msg.Data())
			}

		case sysex.End:
			if h.Message.SysEx.End != nil {
				h.Message.SysEx.End(h.pos, msg.Data())
			}

		case sysex.Continue:
			if h.Message.SysEx.Continue != nil {
				h.Message.SysEx.Continue(h.pos, msg.Data())
			}

		case sysex.Escape:
			if h.Message.SysEx.Escape != nil {
				h.Message.SysEx.Escape(h.pos, msg.Data())
			}

		// this usually takes some time
		case channel.ProgramChange:
			if h.Message.Channel.ProgramChange != nil {
				h.Message.Channel.ProgramChange(h.pos, msg.Channel(), msg.Program())
			}

		// the rest is not that interesting for performance
		case meta.KeySignature:
			if h.Message.Meta.KeySignature != nil {
				h.Message.Meta.KeySignature(*h.pos, msg.Key, msg.IsMajor, msg.Num, msg.IsFlat)
			}

		case meta.Sequence:
			if h.Message.Meta.Sequence != nil {
				h.Message.Meta.Sequence(*h.pos, msg.Text())
			}

		case meta.TrackInstrument:
			if h.Message.Meta.TrackInstrument != nil {
				h.Message.Meta.TrackInstrument(*h.pos, msg.Text())
			}

		case meta.MIDIChannel:
			if h.Message.Meta.MIDIChannel != nil {
				h.Message.Meta.MIDIChannel(*h.pos, msg.Number())
			}

		case meta.MIDIPort:
			if h.Message.Meta.MIDIPort != nil {
				h.Message.Meta.MIDIPort(*h.pos, msg.Number())
			}

		case meta.Text:
			if h.Message.Meta.Text != nil {
				h.Message.Meta.Text(*h.pos, msg.Text())
			}

		case syscommon.SongSelect:
			if h.Message.SysCommon.SongSelect != nil {
				h.Message.SysCommon.SongSelect(msg.Number())
			}

		case syscommon.SongPositionPointer:
			if h.Message.SysCommon.SongPositionPointer != nil {
				h.Message.SysCommon.SongPositionPointer(msg.Number())
			}

		case syscommon.MIDITimingCode:
			if h.Message.SysCommon.MIDITimingCode != nil {
				h.Message.SysCommon.MIDITimingCode(msg.QuarterFrame())
			}

		case meta.Copyright:
			if h.Message.Meta.Copyright != nil {
				h.Message.Meta.Copyright(*h.pos, msg.Text())
			}

		case meta.DevicePort:
			if h.Message.Meta.DevicePort != nil {
				h.Message.Meta.DevicePort(*h.pos, msg.Text())
			}

		//case meta.Undefined, syscommon.Undefined4, syscommon.Undefined5:
		case meta.Undefined:
			if h.Message.Unknown != nil {
				h.Message.Unknown(h.pos, m)
			}

		default:
			switch m {
			case syscommon.TuneRequest:
				if h.Message.SysCommon.TuneRequest != nil {
					h.Message.SysCommon.TuneRequest()
				}
			case meta.EndOfTrack:
				if _, ok := rd.(smf.Reader); ok && h.pos != nil {
					h.pos.Delta = 0
					h.pos.AbsTime = 0
				}
				if h.Message.Meta.EndOfTrack != nil {
					h.Message.Meta.EndOfTrack(*h.pos)
				}
			default:

				if h.Message.Unknown != nil {
					h.Message.Unknown(h.pos, m)
				}

			}

		}

	}

	return
}
