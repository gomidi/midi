package mid

import (
	"github.com/gomidi/midi/midimessage/meta/meter"
	// "bytes"
	"fmt"
	"github.com/gomidi/midi"
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
	_, err := l.midiWriter.wr.Write(realtime.ActiveSensing)
	return err
}

// Continue writes the continue realtime message
func (l *LiveWriter) Continue() error {
	_, err := l.midiWriter.wr.Write(realtime.Continue)
	return err
}

// Reset writes the reset realtime message
func (l *LiveWriter) Reset() error {
	_, err := l.midiWriter.wr.Write(realtime.Reset)
	return err
}

// Start writes the start realtime message
func (l *LiveWriter) Start() error {
	_, err := l.midiWriter.wr.Write(realtime.Start)
	return err
}

// Stop writes the stop realtime message
func (l *LiveWriter) Stop() error {
	_, err := l.midiWriter.wr.Write(realtime.Stop)
	return err
}

// Tick writes the tick realtime message
func (l *LiveWriter) Tick() error {
	_, err := l.midiWriter.wr.Write(realtime.Tick)
	return err
}

// TimingClock writes the timing clock realtime message
func (l *LiveWriter) TimingClock() error {
	_, err := l.midiWriter.wr.Write(realtime.TimingClock)
	return err
}

// MIDITimingCode writes the MIDI Timing Code system message
func (l *LiveWriter) MIDITimingCode(code uint8) error {
	_, err := l.midiWriter.wr.Write(syscommon.MIDITimingCode(code))
	return err
}

// SongPositionPointer writes the song position pointer system message
func (l *LiveWriter) SongPositionPointer(ptr uint16) error {
	_, err := l.midiWriter.wr.Write(syscommon.SongPositionPointer(ptr))
	return err
}

// SongSelect writes the song select system message
func (l *LiveWriter) SongSelect(song uint8) error {
	_, err := l.midiWriter.wr.Write(syscommon.SongSelect(song))
	return err
}

// TuneRequest writes the tune request system message
func (l *LiveWriter) TuneRequest() error {
	_, err := l.midiWriter.wr.Write(syscommon.TuneRequest)
	return err
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
	_, err := m.wr.Write(m.ch.ChannelPressure(pressure))
	return err
}

// KeyPressure writes a key pressure message for the current channel
func (m *midiWriter) KeyPressure(key, pressure uint8) error {
	_, err := m.wr.Write(m.ch.KeyPressure(key, pressure))
	return err
}

// NoteOff writes a note off message for the current channel
func (m *midiWriter) NoteOff(key uint8) error {
	_, err := m.wr.Write(m.ch.NoteOff(key))
	return err
}

// NoteOn writes a note on message for the current channel
func (m *midiWriter) NoteOn(key, veloctiy uint8) error {
	_, err := m.wr.Write(m.ch.NoteOn(key, veloctiy))
	return err
}

// PitchBend writes a pitch bend message for the current channel
// For reset value, use 0, for lowest -8191 and highest 8191
// Or use the pitch constants of midimessage/channel
func (m *midiWriter) PitchBend(value int16) error {
	_, err := m.wr.Write(m.ch.PitchBend(value))
	return err
}

// ProgramChange writes a program change message for the current channel
// Program numbers start with 0 for program 1.
func (m *midiWriter) ProgramChange(program uint8) error {
	_, err := m.wr.Write(m.ch.ProgramChange(program))
	return err
}

// CC writes a control change message. It is meant to be used in conjunction
// with the midimessages/cc package.
func (m *midiWriter) CC(cch channel.ControlChange) error {
	_, err := m.wr.Write(cch)
	return err
}

// ControlChange writes a control change message for the current channel
func (m *midiWriter) ControlChange(controller, value uint8) error {
	_, err := m.wr.Write(m.ch.ControlChange(controller, value))
	return err
}

// SysEx writes sysex data
func (m *midiWriter) SysEx(data []byte) error {
	_, err := m.wr.Write(sysex.SysEx(data))
	return err
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
	_, err := s.wr.Write(meta.EndOfTrack)
	return err
}

// Copyright writes the copyright meta message
func (s *SMFWriter) Copyright(text string) error {
	_, err := s.wr.Write(meta.Copyright(text))
	return err
}

// Cuepoint writes the cuepoint meta message
func (s *SMFWriter) Cuepoint(text string) error {
	_, err := s.wr.Write(meta.Cuepoint(text))
	return err
}

// DevicePort writes the device port meta message
func (s *SMFWriter) DevicePort(port string) error {
	_, err := s.wr.Write(meta.DevicePort(port))
	return err
}

// KeySignature writes the key signature meta message.
// A more comfortable way is to use the Key method in conjunction
// with the midimessage/meta/key package
func (s *SMFWriter) KeySignature(key uint8, ismajor bool, num uint8, isflat bool) error {
	_, err := s.wr.Write(meta.KeySignature{Key: key, IsMajor: ismajor, Num: num, IsFlat: isflat})
	return err
}

// Key writes the given key signature meta message.
// It is supposed to be used with the midimessage/meta/key package
func (s *SMFWriter) Key(keysig meta.KeySignature) error {
	_, err := s.wr.Write(keysig)
	return err
}

// Lyric writes the lyric meta message
func (s *SMFWriter) Lyric(text string) error {
	_, err := s.wr.Write(meta.Lyric(text))
	return err
}

// Marker writes the marker meta message
func (s *SMFWriter) Marker(text string) error {
	_, err := s.wr.Write(meta.Marker(text))
	return err
}

// MIDIChannel writes the deprecated MIDI channel meta message
func (s *SMFWriter) MIDIChannel(ch uint8) error {
	_, err := s.wr.Write(meta.MIDIChannel(ch))
	return err
}

// MIDIPort writes the deprecated MIDI port meta message
func (s *SMFWriter) MIDIPort(port uint8) error {
	_, err := s.wr.Write(meta.MIDIPort(port))
	return err
}

// ProgramName writes the program name meta message
func (s *SMFWriter) ProgramName(text string) error {
	_, err := s.wr.Write(meta.ProgramName(text))
	return err
}

// Sequence writes the sequence (name) meta message
func (s *SMFWriter) Sequence(text string) error {
	_, err := s.wr.Write(meta.Sequence(text))
	return err
}

// SequenceNumber writes the sequence number meta message
func (s *SMFWriter) SequenceNumber(no uint16) error {
	_, err := s.wr.Write(meta.SequenceNumber(no))
	return err
}

// SequencerSpecific writes a custom sequences specific meta message
func (s *SMFWriter) SequencerSpecific(data []byte) error {
	_, err := s.wr.Write(meta.SequencerSpecific(data))
	return err
}

// SMPTEOffset writes the SMPTE offset meta message
func (s *SMFWriter) SMPTEOffset(hour, minute, second, frame, fractionalFrame byte) error {
	_, err := s.wr.Write(meta.SMPTEOffset{
		Hour:            hour,
		Minute:          minute,
		Second:          second,
		Frame:           frame,
		FractionalFrame: fractionalFrame,
	})
	return err
}

// Tempo writes the tempo meta message
func (s *SMFWriter) Tempo(bpm uint32) error {
	_, err := s.wr.Write(meta.Tempo(bpm))
	return err
}

// Text writes the text meta message
func (s *SMFWriter) Text(text string) error {
	_, err := s.wr.Write(meta.Text(text))
	return err
}

// Meter writes the time signature meta message in a more comfortable way.
// Numerator and Denominator are decimals.
func (s *SMFWriter) Meter(numerator, denominator uint8) error {
	_, err := s.wr.Write(meter.Meter(numerator, denominator))
	return err
}

// TimeSignature writes the time signature meta message.
// Numerator and Denominator are decimals.
// If you don't want to deal with clocks per click and demisemiquaverperquarter,
// user the Meter method instead.
func (s *SMFWriter) TimeSignature(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) error {
	_, err := s.wr.Write(meta.TimeSignature{
		Numerator:                numerator,
		Denominator:              denominator,
		ClocksPerClick:           clocksPerClick,
		DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
	})
	return err
}

// Track writes the track name aka instrument name meta message
func (s *SMFWriter) Track(track string) error {
	_, err := s.wr.Write(meta.Track(track))
	return err
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
