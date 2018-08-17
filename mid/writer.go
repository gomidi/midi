package mid

import (
	"github.com/gomidi/midi/midimessage/meta/meter"
	// "bytes"
	// "encoding/binary"
	"fmt"
	"github.com/gomidi/midi"
	// "github.com/gomidi/midi/internal/midilib"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/gomidi/midi/midimessage/syscommon"
	"github.com/gomidi/midi/midimessage/sysex"
	"github.com/gomidi/midi/midiwriter"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfwriter"
	"io"
	"os"
	// "time"
)

// LiveWriter writes live MIDI data. Its methods must not be called concurrently
// The LiveWrite does not make use of the running status.
type LiveWriter struct {
	*midiWriter
}

// ActiveSensing writes the active sensing realtime message
func (l *LiveWriter) ActiveSensing() error {
	return l.midiWriter.wr.Write(realtime.ActiveSensing)
}

// Continue writes the continue realtime message
func (l *LiveWriter) Continue() error {
	return l.midiWriter.wr.Write(realtime.Continue)
}

// Reset writes the reset realtime message
func (l *LiveWriter) Reset() error {
	return l.midiWriter.wr.Write(realtime.Reset)
}

// Start writes the start realtime message
func (l *LiveWriter) Start() error {
	return l.midiWriter.wr.Write(realtime.Start)
}

// Stop writes the stop realtime message
func (l *LiveWriter) Stop() error {
	return l.midiWriter.wr.Write(realtime.Stop)
}

// Tick writes the tick realtime message
func (l *LiveWriter) Tick() error {
	return l.midiWriter.wr.Write(realtime.Tick)
}

// TimingClock writes the timing clock realtime message
func (l *LiveWriter) TimingClock() error {
	return l.midiWriter.wr.Write(realtime.TimingClock)
}

// MIDITimingCode writes the MIDI Timing Code system message
func (l *LiveWriter) MIDITimingCode(code uint8) error {
	return l.midiWriter.wr.Write(syscommon.MIDITimingCode(code))
}

// SongPositionPointer writes the song position pointer system message
func (l *LiveWriter) SongPositionPointer(ptr uint16) error {
	return l.midiWriter.wr.Write(syscommon.SongPositionPointer(ptr))
}

// SongSelect writes the song select system message
func (l *LiveWriter) SongSelect(song uint8) error {
	return l.midiWriter.wr.Write(syscommon.SongSelect(song))
}

// TuneRequest writes the tune request system message
func (l *LiveWriter) TuneRequest() error {
	return l.midiWriter.wr.Write(syscommon.TuneRequest)
}

type midiWriter struct {
	wr midi.Writer
	ch channel.Channel
}

// SetChannel sets the channel for the following midi messages
// Channel numbers are counted from 0 to 15 (MIDI channel 1 to 16).
// The initial channel number is 0.
func (m *midiWriter) SetChannel(no uint8 /* 0-15 */) {
	m.ch = channel.Channel(no)
}

// ChannelPressure writes a channel pressure message for the current channel
func (m *midiWriter) ChannelPressure(pressure uint8) error {
	return m.wr.Write(m.ch.ChannelPressure(pressure))
}

// KeyPressure writes a key pressure message for the current channel
func (m *midiWriter) KeyPressure(key, pressure uint8) error {
	return m.wr.Write(m.ch.KeyPressure(key, pressure))
}

// NoteOff writes a note off message for the current channel
func (m *midiWriter) NoteOff(key uint8) error {
	return m.wr.Write(m.ch.NoteOff(key))
}

// NoteOn writes a note on message for the current channel
func (m *midiWriter) NoteOn(key, veloctiy uint8) error {
	return m.wr.Write(m.ch.NoteOn(key, veloctiy))
}

// PitchBend writes a pitch bend message for the current channel
// For reset value, use 0, for lowest -8191 and highest 8191
// Or use the pitch constants of midimessage/channel
func (m *midiWriter) PitchBend(value int16) error {
	return m.wr.Write(m.ch.PitchBend(value))
}

// ProgramChange writes a program change message for the current channel
// Program numbers start with 0 for program 1.
func (m *midiWriter) ProgramChange(program uint8) error {
	return m.wr.Write(m.ch.ProgramChange(program))
}

// MsbLsb writes a Msb control change message, followed by a Lsb control change message
// for the current channel
// For more comfortable use, used it in conjunction with the gomidi/cc package
func (m *midiWriter) MsbLsb(msb, lsb uint8, value uint16) error {

	var b = make([]byte, 2)
	b[1] = byte(value & 0x7F)
	b[0] = byte((value >> 7) & 0x7F)

	/*
		r := midilib.MsbLsbSigned(value)

		var b = make([]byte, 2)

		binary.BigEndian.PutUint16(b, r)
	*/
	err := m.ControlChange(msb, b[0])
	if err != nil {
		return err
	}
	return m.ControlChange(lsb, b[1])
}

// ControlChange writes a control change message for the current channel
// For more comfortable use, used it in conjunction with the gomidi/cc package
func (m *midiWriter) ControlChange(controller, value uint8) error {
	return m.wr.Write(m.ch.ControlChange(controller, value))
}

func (m *midiWriter) ControlChangeOff(controller uint8) error {
	return m.ControlChange(controller, 0)
}

func (m *midiWriter) ControlChangeOn(controller uint8) error {
	return m.ControlChange(controller, 127)
}

// SysEx writes sysex data
func (m *midiWriter) SysEx(data []byte) error {
	return m.wr.Write(sysex.SysEx(data))
}

// SMFWriter writes SMF MIDI data. Its methods must not be called concurrently
type SMFWriter struct {
	wr smf.Writer
	*midiWriter
	finishedTracks uint16
	dest           io.Writer
}

// SetDelta sets the delta ticks to the next message
func (s *SMFWriter) SetDelta(deltatime uint32) {
	s.wr.SetDelta(deltatime)
}

// EndOfTrack signals the end of a track
func (s *SMFWriter) EndOfTrack() error {

	if no := s.wr.Header().NumTracks; s.finishedTracks >= no {
		return fmt.Errorf("too many tracks: in header: %v, closed: %v", no, s.finishedTracks+1)
	}
	s.finishedTracks++
	return s.wr.Write(meta.EndOfTrack)
}

// Copyright writes the copyright meta message
func (s *SMFWriter) Copyright(text string) error {
	return s.wr.Write(meta.Copyright(text))
}

// Cuepoint writes the cuepoint meta message
func (s *SMFWriter) Cuepoint(text string) error {
	return s.wr.Write(meta.Cuepoint(text))
}

// DevicePort writes the device port meta message
func (s *SMFWriter) DevicePort(port string) error {
	return s.wr.Write(meta.DevicePort(port))
}

// KeySignature writes the key signature meta message.
// A more comfortable way is to use the Key method in conjunction
// with the midimessage/meta/key package
func (s *SMFWriter) KeySignature(key uint8, ismajor bool, num uint8, isflat bool) error {
	return s.wr.Write(meta.KeySignature{Key: key, IsMajor: ismajor, Num: num, IsFlat: isflat})
}

// Key writes the given key signature meta message.
// It is supposed to be used with the midimessage/meta/key package
func (s *SMFWriter) Key(keysig meta.KeySignature) error {
	return s.wr.Write(keysig)
}

// Lyric writes the lyric meta message
func (s *SMFWriter) Lyric(text string) error {
	return s.wr.Write(meta.Lyric(text))
}

// Marker writes the marker meta message
func (s *SMFWriter) Marker(text string) error {
	return s.wr.Write(meta.Marker(text))
}

// MIDIChannel writes the deprecated MIDI channel meta message
func (s *SMFWriter) MIDIChannel(ch uint8) error {
	return s.wr.Write(meta.MIDIChannel(ch))
}

// MIDIPort writes the deprecated MIDI port meta message
func (s *SMFWriter) MIDIPort(port uint8) error {
	return s.wr.Write(meta.MIDIPort(port))
}

// ProgramName writes the program name meta message
func (s *SMFWriter) ProgramName(text string) error {
	return s.wr.Write(meta.ProgramName(text))
}

// Sequence writes the sequence (name) meta message
func (s *SMFWriter) Sequence(text string) error {
	return s.wr.Write(meta.Sequence(text))
}

// SequenceNumber writes the sequence number meta message
func (s *SMFWriter) SequenceNumber(no uint16) error {
	return s.wr.Write(meta.SequenceNumber(no))
}

// SequencerSpecific writes a custom sequences specific meta message
func (s *SMFWriter) SequencerSpecific(data []byte) error {
	return s.wr.Write(meta.SequencerSpecific(data))
}

// SMPTEOffset writes the SMPTE offset meta message
func (s *SMFWriter) SMPTEOffset(hour, minute, second, frame, fractionalFrame byte) error {
	return s.wr.Write(meta.SMPTEOffset{
		Hour:            hour,
		Minute:          minute,
		Second:          second,
		Frame:           frame,
		FractionalFrame: fractionalFrame,
	})
}

// Tempo writes the tempo meta message
func (s *SMFWriter) Tempo(bpm uint32) error {
	return s.wr.Write(meta.Tempo(bpm))
}

// Text writes the text meta message
func (s *SMFWriter) Text(text string) error {
	return s.wr.Write(meta.Text(text))
}

// Meter writes the time signature meta message in a more comfortable way.
// Numerator and Denominator are decimals.
func (s *SMFWriter) Meter(numerator, denominator uint8) error {
	return s.wr.Write(meter.Meter(numerator, denominator))
}

// TimeSignature writes the time signature meta message.
// Numerator and Denominator are decimals.
// If you don't want to deal with clocks per click and demisemiquaverperquarter,
// user the Meter method instead.
func (s *SMFWriter) TimeSignature(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) error {
	return s.wr.Write(meta.TimeSignature{
		Numerator:                numerator,
		Denominator:              denominator,
		ClocksPerClick:           clocksPerClick,
		DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
	})
}

// Track writes the track name aka instrument name meta message
func (s *SMFWriter) Track(track string) error {
	return s.wr.Write(meta.Track(track))
}

// NewSMFWriter returns a new SMFWriter that writes to dest.
// It panics if numtracks is == 0.
func NewSMFWriter(dest io.Writer, numtracks uint16, options ...smfwriter.Option) *SMFWriter {
	if numtracks == 0 {
		panic("numtracks must be > 0")
	}

	options = append(
		[]smfwriter.Option{
			smfwriter.NumTracks(numtracks),
			smfwriter.TimeFormat(smf.MetricTicks(0)),
		}, options...)

	wr := smfwriter.New(dest, options...)
	return &SMFWriter{
		dest:       dest,
		wr:         wr,
		midiWriter: &midiWriter{wr: wr, ch: channel.Channel0},
	}
}

// NewSMFFile creates a new SMF file and allows writer to write to it.
// The file is guaranteed to be closed when returning.
// The last track is closed automatically, if needed.
func NewSMFFile(file string, numtracks uint16, writer func(*SMFWriter) error, options ...smfwriter.Option) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	wr := NewSMFWriter(f, numtracks, options...)
	if writer != nil {
		err = writer(wr)
		if err != nil {
			return err
		}
	}

	if no := wr.wr.Header().NumTracks; wr.finishedTracks < no {
		err := wr.EndOfTrack()
		if err != nil {
			return err
		}
	}

	return err
}

// NewLiveWriter creates and new LiveWriter.
func NewLiveWriter(dest io.Writer) *LiveWriter {
	wr := midiwriter.New(dest, midiwriter.NoRunningStatus())
	return &LiveWriter{&midiWriter{wr: wr, ch: channel.Channel0}}
}
