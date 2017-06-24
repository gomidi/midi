package handler

import (
	"fmt"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/messages/syscommon"
	"github.com/gomidi/midi/messages/sysex"
	"github.com/gomidi/midi/smf"
)

// Pos is the position of the event inside a standard midi file (SMF).
type Pos struct {
	// the Track number
	Track uint16

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
//
// When reading an SMF file (via Handler.ReadSMF), the passed *Pos is set,
// when reading live data (via Handler.ReadLive) it is nil.
type Handler struct {
	pos *Pos
	//Event  Event
	logger Logger

	// SMF header informations
	Format           func(smf.Format) // the midi file format (0=single track,1=multitrack,2=sequential tracks)
	NumTracks        func(n uint16)   // number of tracks
	TimeCode         func(uint16)
	QuarterNoteTicks func(uint16)

	// SMF general settings
	Copyright     func(p *Pos, text string)
	Tempo         func(p *Pos, bpm uint32)
	TimeSignature func(p *Pos, num, denom uint8)
	KeySignature  func(p *Pos, key uint8, ismajor bool, num_accidentals uint8, accidentals_are_flat bool)

	// SMF tracks and sequence definitions
	TrackInstrument func(p *Pos, name string)
	Sequence        func(p *Pos, name string)
	SequenceNumber  func(p *Pos, number uint16)

	// SMF port description
	DevicePort func(p *Pos, name string)

	// SMF text entries
	Marker   func(p *Pos, text string)
	CuePoint func(p *Pos, text string)
	Text     func(p *Pos, text string)
	Lyric    func(p *Pos, text string)

	// SMF end of track
	EndOfTrack func(p *Pos)

	// channel messages
	// NoteOn is just called for noteon messages with a velocity > 0
	// noteon messages with velocity == 0 will trigger NoteOff with a velocity of 0
	NoteOn func(p *Pos, channel, pitch, velocity uint8)
	// NoteOff is triggered by noteoff messages (then the given velocity is passed)
	// and by noteon messages of velocity 0 (then velocity is 0)
	NoteOff              func(p *Pos, channel, pitch uint8, velocity uint8)
	PolyphonicAfterTouch func(p *Pos, channel, pitch, pressure uint8)
	ControlChange        func(p *Pos, channel, controller, value uint8)
	ProgramChange        func(p *Pos, channel, program uint8)
	AfterTouch           func(p *Pos, channel, pressure uint8)
	PitchWheel           func(p *Pos, channel uint8, value int16)

	// system messages
	SysExComplete func(p *Pos, data []byte)
	SysExStart    func(p *Pos, data []byte)
	SysExContinue func(p *Pos, data []byte)
	SysExEnd      func(p *Pos, data []byte)
	SysExEscape   func(p *Pos, data []byte)

	// system common
	TuneRequest         func()
	SongSelect          func(num uint8)
	SongPositionPointer func(pos uint16)
	MIDITimingCode      func(frame uint8)

	// realtime messages
	Reset       func()
	Clock       func()
	Tick        func()
	Start       func()
	Continue    func()
	Stop        func()
	ActiveSense func()

	// deprecated
	MIDIChannel func(p *Pos, channel uint8)
	MIDIPort    func(p *Pos, port uint8)

	// undefined
	UndefinedMeta       func(p *Pos, typ byte, data []byte)
	UndefinedSysCommon4 func(p *Pos)
	UndefinedSysCommon5 func(p *Pos)
	UndefinedRealtime4  func()
	Unknown             func(p *Pos, info string)

	// is called in addition to other functions, if set.
	Each func(*Pos, midi.Message)

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
// for realtime reading, the passed *Pos is nil
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

		if h.Each != nil {
			h.Each(h.pos, m)
		}

		switch msg := m.(type) {

		// most common event, should be exact
		case channel.NoteOn:
			if h.NoteOn != nil {
				h.NoteOn(h.pos, msg.Channel(), msg.Pitch(), msg.Velocity())
			}

		// proably second most common
		case channel.NoteOff:
			if h.NoteOff != nil {
				h.NoteOff(h.pos, msg.Channel(), msg.Pitch(), 0)
			}

		case channel.NoteOffPedantic:
			if h.NoteOff != nil {
				h.NoteOff(h.pos, msg.Channel(), msg.Pitch(), msg.Velocity())
			}

		// if send there often are a lot of them
		case channel.PitchWheel:
			if h.PitchWheel != nil {
				h.PitchWheel(h.pos, msg.Channel(), msg.Value())
			}

		case channel.PolyphonicAfterTouch:
			if h.PolyphonicAfterTouch != nil {
				h.PolyphonicAfterTouch(h.pos, msg.Channel(), msg.Pitch(), msg.Pressure())
			}

		case channel.AfterTouch:
			if h.AfterTouch != nil {
				h.AfterTouch(h.pos, msg.Channel(), msg.Pressure())
			}

		case channel.ControlChange:
			if h.ControlChange != nil {
				h.ControlChange(h.pos, msg.Channel(), msg.Controller(), msg.Value())
			}

		case meta.Tempo:
			if h.Tempo != nil {
				h.Tempo(h.pos, msg.BPM())
			}

		case meta.TimeSignature:
			if h.TimeSignature != nil {
				h.TimeSignature(h.pos, msg.Numerator, msg.Denominator)
			}

			// may be for karaoke we need to be fast
		case meta.Lyric:
			if h.Lyric != nil {
				h.Lyric(h.pos, msg.Text())
			}

		// may be useful to synchronize by sequence number
		case meta.SequenceNumber:
			if h.SequenceNumber != nil {
				h.SequenceNumber(h.pos, msg.Number())
			}

		// markers and cuepoints could also be useful when communication sections or sequences between devices
		case meta.Marker:
			if h.Marker != nil {
				h.Marker(h.pos, msg.Text())
			}

		case meta.CuePoint:
			if h.CuePoint != nil {
				h.CuePoint(h.pos, msg.Text())
			}

		case sysex.SysEx:
			if h.SysExComplete != nil {
				h.SysExComplete(h.pos, msg.Data())
			}

		case sysex.Start:
			if h.SysExStart != nil {
				h.SysExStart(h.pos, msg.Data())
			}

		case sysex.End:
			if h.SysExEnd != nil {
				h.SysExEnd(h.pos, msg.Data())
			}

		case sysex.Continue:
			if h.SysExContinue != nil {
				h.SysExContinue(h.pos, msg.Data())
			}

		case sysex.Escape:
			if h.SysExEscape != nil {
				h.SysExEscape(h.pos, msg.Data())
			}

		// this usually takes some time
		case channel.ProgramChange:
			if h.ProgramChange != nil {
				h.ProgramChange(h.pos, msg.Channel(), msg.Program())
			}

		// the rest is not that interesting for performance
		case meta.KeySignature:
			if h.KeySignature != nil {
				h.KeySignature(h.pos, msg.Key, msg.IsMajor, msg.Num, msg.IsFlat)
			}

		case meta.Sequence:
			if h.Sequence != nil {
				h.Sequence(h.pos, msg.Text())
			}

		case meta.TrackInstrument:
			if h.TrackInstrument != nil {
				h.TrackInstrument(h.pos, msg.Text())
			}

		case meta.MIDIChannel:
			if h.MIDIChannel != nil {
				h.MIDIChannel(h.pos, msg.Number())
			}

		case meta.MIDIPort:
			if h.MIDIPort != nil {
				h.MIDIPort(h.pos, msg.Number())
			}

		case meta.Text:
			if h.Text != nil {
				h.Text(h.pos, msg.Text())
			}

		case syscommon.SongSelect:
			if h.SongSelect != nil {
				h.SongSelect(msg.Number())
			}

		case syscommon.SongPositionPointer:
			if h.SongPositionPointer != nil {
				h.SongPositionPointer(msg.Number())
			}

		case syscommon.MIDITimingCode:
			if h.MIDITimingCode != nil {
				h.MIDITimingCode(msg.QuarterFrame())
			}

		case meta.Copyright:
			if h.Copyright != nil {
				h.Copyright(h.pos, msg.Text())
			}

		case meta.DevicePort:
			if h.DevicePort != nil {
				h.DevicePort(h.pos, msg.Text())
			}

		case meta.Undefined:
			if h.UndefinedMeta != nil {
				h.UndefinedMeta(h.pos, msg.Typ, msg.Data)
			}

		case syscommon.Undefined4:
			if h.UndefinedSysCommon4 != nil {
				h.UndefinedSysCommon4(h.pos)
			}

		case syscommon.Undefined5:
			if h.UndefinedSysCommon5 != nil {
				h.UndefinedSysCommon5(h.pos)
			}

		default:
			switch m {
			case syscommon.TuneRequest:
				if h.TuneRequest != nil {
					h.TuneRequest()
				}
			case meta.EndOfTrack:
				if h.EndOfTrack != nil {
					h.EndOfTrack(h.pos)
				}
			default:

				if h.Unknown != nil {
					h.Unknown(h.pos, fmt.Sprintf("%T %#v", m, m))
				}

			}

		}

	}

	return
}
