package handler

import (
	"fmt"
	"io"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/live/midireader"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/messages/realtime"
	"github.com/gomidi/midi/messages/syscommon"
	"github.com/gomidi/midi/messages/sysex"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfreader"
)

// Logger is the inferface used by Handler for logging incoming messages.
type Logger interface {
	Printf(format string, vals ...interface{})
}

type logfunc func(format string, vals ...interface{})

func (l logfunc) Printf(format string, vals ...interface{}) {
	l(format, vals...)
}

func printf(format string, vals ...interface{}) {
	fmt.Printf(format, vals...)
}

// Pos is the position of the event inside a standard midi file (SMF).
type Pos struct {
	// the Track number
	Track uint16

	// the delta time to the previous message in the same track
	Delta uint32

	// the absolute time from the beginning of the track
	AbsTime uint64
}

// Option configures the handler
type Option func(*Handler)

// SetLogger allows to set a custom logger for the handler
func SetLogger(l Logger) Option {
	return func(h *Handler) {
		h.logger = l
	}
}

// NoLogger is an option to disable the defaut logging of a handler
func NoLogger() Option {
	return func(h *Handler) {
		h.logger = nil
	}
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
	SysEx func(p *Pos, data []byte)

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

// ReadLive reads midi messages from src until an error or io.EOF happens.
//
// If io.EOF happend the returned error is nil.
//
// ReadLive does not close the src.
//
// The messages are dispatched to the corresponding attached functions of the handler.
//
// They must be attached before Handler.ReadLive is called
// and they must not be unset or replaced until ReadLive returns.
//
// The *Pos parameter that is passed to the functions is nil, because we are in a live setting.
func (h *Handler) ReadLive(src io.Reader, options ...midireader.Option) (err error) {
	rthandler := func(m realtime.Message) {
		switch m {
		// ticks (most important, must be sent every 10 milliseconds) comes first
		case realtime.Tick:
			if h.Tick != nil {
				h.Tick()
			}

		// clock a bit slower synchronization method (24 MIDI Clocks in every quarter note) comes next
		case realtime.TimingClock:
			if h.Clock != nil {
				h.Clock()
			}

		// ok starting and continuing should not take too lpng
		case realtime.Start:
			if h.Start != nil {
				h.Start()
			}

		case realtime.Continue:
			if h.Continue != nil {
				h.Continue()
			}

		// Active Sense must come every 300 milliseconds (but is seldom implemented)
		case realtime.ActiveSensing:
			if h.ActiveSense != nil {
				h.ActiveSense()
			}

		// put any user defined realtime message here
		case realtime.Undefined4:
			if h.UndefinedRealtime4 != nil {
				h.UndefinedRealtime4()
			}

		// ok, stopping is not so urgent
		case realtime.Stop:
			if h.Stop != nil {
				h.Stop()
			}

		// reset may take some time
		case realtime.Reset:
			if h.Reset != nil {
				h.Reset()
			}

		}
	}

	rd := midireader.New(src, rthandler, options...)
	err = h.read(rd)

	if err == io.EOF {
		return nil
	}

	return
}

// read reads the messages from the midi.Reader (which might be an smf reader
// for realtime reading, the passed *Pos is nil
func (h *Handler) read(rd midi.Reader) (err error) {
	var evt midi.Message

	for {
		evt, err = rd.Read()
		if err != nil {
			break
		}

		if frd, ok := rd.(smf.Reader); ok && h.pos != nil {
			h.pos.Delta = frd.Delta()
			h.pos.AbsTime += uint64(h.pos.Delta)
			h.pos.Track = frd.Track()
		}

		if h.logger != nil {
			h.log(evt)
		}

		if h.Each != nil {
			h.Each(h.pos, evt)
		}

		switch ev := evt.(type) {

		// most common event, should be exact
		case channel.NoteOn:
			if h.NoteOn != nil {
				h.NoteOn(h.pos, ev.Channel(), ev.Pitch(), ev.Velocity())
			}

		// proably second most common
		case channel.NoteOff:
			if h.NoteOff != nil {
				h.NoteOff(h.pos, ev.Channel(), ev.Pitch(), 0)
			}

		case channel.NoteOffPedantic:
			if h.NoteOff != nil {
				h.NoteOff(h.pos, ev.Channel(), ev.Pitch(), ev.Velocity())
			}

		// if send there often are a lot of them
		case channel.PitchWheel:
			if h.PitchWheel != nil {
				h.PitchWheel(h.pos, ev.Channel(), ev.Value())
			}

		case channel.PolyphonicAfterTouch:
			if h.PolyphonicAfterTouch != nil {
				h.PolyphonicAfterTouch(h.pos, ev.Channel(), ev.Pitch(), ev.Pressure())
			}

		case channel.AfterTouch:
			if h.AfterTouch != nil {
				h.AfterTouch(h.pos, ev.Channel(), ev.Pressure())
			}

		case channel.ControlChange:
			if h.ControlChange != nil {
				h.ControlChange(h.pos, ev.Channel(), ev.Controller(), ev.Value())
			}

		case meta.Tempo:
			if h.Tempo != nil {
				h.Tempo(h.pos, ev.BPM())
			}

		case meta.TimeSignature:
			if h.TimeSignature != nil {
				h.TimeSignature(h.pos, ev.Numerator, ev.Denominator)
			}

			// may be for karaoke we need to be fast
		case meta.Lyric:
			if h.Lyric != nil {
				h.Lyric(h.pos, ev.Text())
			}

		// may be useful to synchronize by sequence number
		case meta.SequenceNumber:
			if h.SequenceNumber != nil {
				h.SequenceNumber(h.pos, ev.Number())
			}

		// markers and cuepoints could also be useful when communication sections or sequences between devices
		case meta.Marker:
			if h.Marker != nil {
				h.Marker(h.pos, ev.Text())
			}

		case meta.CuePoint:
			if h.CuePoint != nil {
				h.CuePoint(h.pos, ev.Text())
			}

		case sysex.SysEx:
			if h.SysEx != nil {
				h.SysEx(h.pos, ev.Bytes())
			}

		// this usually takes some time
		case channel.ProgramChange:
			if h.ProgramChange != nil {
				h.ProgramChange(h.pos, ev.Channel(), ev.Program())
			}

		// the rest is not that interesting for performance
		case meta.KeySignature:
			if h.KeySignature != nil {
				h.KeySignature(h.pos, ev.Key, ev.IsMajor, ev.Num, ev.IsFlat)
			}

		case meta.Sequence:
			if h.Sequence != nil {
				h.Sequence(h.pos, ev.Text())
			}

		case meta.TrackInstrument:
			if h.TrackInstrument != nil {
				h.TrackInstrument(h.pos, ev.Text())
			}

		case meta.MIDIChannel:
			if h.MIDIChannel != nil {
				h.MIDIChannel(h.pos, ev.Number())
			}

		case meta.MIDIPort:
			if h.MIDIPort != nil {
				h.MIDIPort(h.pos, ev.Number())
			}

		case meta.Text:
			if h.Text != nil {
				h.Text(h.pos, ev.Text())
			}

		case syscommon.SongSelect:
			if h.SongSelect != nil {
				h.SongSelect(ev.Number())
			}

		case syscommon.SongPositionPointer:
			if h.SongPositionPointer != nil {
				h.SongPositionPointer(ev.Number())
			}

		case syscommon.MIDITimingCode:
			if h.MIDITimingCode != nil {
				h.MIDITimingCode(ev.QuarterFrame())
			}

		case meta.Copyright:
			if h.Copyright != nil {
				h.Copyright(h.pos, ev.Text())
			}

		case meta.DevicePort:
			if h.DevicePort != nil {
				h.DevicePort(h.pos, ev.Text())
			}

		case meta.Undefined:
			if h.UndefinedMeta != nil {
				h.UndefinedMeta(h.pos, ev.Typ, ev.Data)
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
			switch evt {
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
					h.Unknown(h.pos, fmt.Sprintf("%T %#v", evt, evt))
				}

			}

		}

	}

	return
}

// ReadSMFFile open, reads and closes a complete SMF file.
// If the read content was a valid midi file, nil is returned.
//
// The messages are dispatched to the corresponding attached functions of the handler.
//
// They must be attached before Handler.ReadSMF is called
// and they must not be unset or replaced until ReadSMF returns.
//
// The *Pos parameter that is passed to the functions is always, because we are reading a file.
func (h *Handler) ReadSMFFile(file string, options ...smfreader.Option) error {
	h.errSMF = nil
	err := smfreader.ReadFile(file, h.readSMF, options...)
	if err != nil {
		return err
	}
	return h.errSMF
}

// ReadSMF reads midi messages from src (which is supposed to be the content of a standard midi file (SMF))
// until an error or io.EOF happens.
//
// ReadSMF does not close the src.
//
// If the read content was a valid midi file, nil is returned.
//
// The messages are dispatched to the corresponding attached functions of the handler.
//
// They must be attached before Handler.ReadSMF is called
// and they must not be unset or replaced until ReadSMF returns.
//
// The *Pos parameter that is passed to the functions is always, because we are reading a file.
func (h *Handler) ReadSMF(src io.Reader, options ...smfreader.Option) error {
	h.errSMF = nil
	h.pos = &Pos{}
	rd := smfreader.New(src, options...)
	h.readSMF(rd)
	return h.errSMF
}

func (h *Handler) readSMF(rd smf.Reader) {
	hd, err := rd.ReadHeader()

	if err != nil {
		h.errSMF = err
		return
	}

	if h.Format != nil {
		h.Format(hd.Format())
	}

	if h.NumTracks != nil {
		h.NumTracks(hd.NumTracks())
	}

	tf, tval := hd.TimeFormat()

	if tf == smf.TimeCode && h.TimeCode != nil {
		h.TimeCode(tval)
	}

	if tf == smf.QuarterNoteTicks && h.QuarterNoteTicks != nil {
		h.QuarterNoteTicks(tval)
	}

	// use err here
	err = h.read(rd)
	if err != io.EOF {
		h.errSMF = err
	}

	return
}
