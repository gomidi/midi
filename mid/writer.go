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

type LiveWriter struct {
	*midiWriter
}

func (l *LiveWriter) ActiveSensing() error {
	_, err := l.midiWriter.wr.Write(realtime.ActiveSensing)
	return err
}

func (l *LiveWriter) Continue() error {
	_, err := l.midiWriter.wr.Write(realtime.Continue)
	return err
}

func (l *LiveWriter) Reset() error {
	_, err := l.midiWriter.wr.Write(realtime.Reset)
	return err
}

func (l *LiveWriter) Start() error {
	_, err := l.midiWriter.wr.Write(realtime.Start)
	return err
}

func (l *LiveWriter) Stop() error {
	_, err := l.midiWriter.wr.Write(realtime.Stop)
	return err
}

func (l *LiveWriter) Tick() error {
	_, err := l.midiWriter.wr.Write(realtime.Tick)
	return err
}

func (l *LiveWriter) TimingClock() error {
	_, err := l.midiWriter.wr.Write(realtime.TimingClock)
	return err
}

func (l *LiveWriter) MIDITimingCode(code uint8) error {
	_, err := l.midiWriter.wr.Write(syscommon.MIDITimingCode(code))
	return err
}

func (l *LiveWriter) SongPositionPointer(ptr uint16) error {
	_, err := l.midiWriter.wr.Write(syscommon.SongPositionPointer(ptr))
	return err
}

func (l *LiveWriter) SongSelect(song uint8) error {
	_, err := l.midiWriter.wr.Write(syscommon.SongSelect(song))
	return err
}

func (l *LiveWriter) TuneRequest() error {
	_, err := l.midiWriter.wr.Write(syscommon.TuneRequest)
	return err
}

type midiWriter struct {
	wr midi.Writer
	ch channel.Channel
}

func (m *midiWriter) SetChannel(no uint8 /* 0-15 */) {
	m.ch = channel.New(no)
}

func (m *midiWriter) ChannelPressure(pressure uint8) error {
	_, err := m.wr.Write(m.ch.ChannelPressure(pressure))
	return err
}

func (m *midiWriter) KeyPressure(key, pressure uint8) error {
	_, err := m.wr.Write(m.ch.KeyPressure(key, pressure))
	return err
}

func (m *midiWriter) NoteOff(key uint8) error {
	_, err := m.wr.Write(m.ch.NoteOff(key))
	return err
}

func (m *midiWriter) NoteOn(key, veloctiy uint8) error {
	_, err := m.wr.Write(m.ch.NoteOn(key, veloctiy))
	return err
}

func (m *midiWriter) PitchBend(value int16) error {
	_, err := m.wr.Write(m.ch.PitchBend(value))
	return err
}

func (m *midiWriter) ProgramChange(program uint8) error {
	_, err := m.wr.Write(m.ch.ProgramChange(program))
	return err
}

func (m *midiWriter) CC(cch channel.ControlChange) error {
	_, err := m.wr.Write(cch)
	return err
}

func (m *midiWriter) ControlChange(controller, value uint8) error {
	_, err := m.wr.Write(m.ch.ControlChange(controller, value))
	return err
}

func (m *midiWriter) SysEx(data []byte) error {
	_, err := m.wr.Write(sysex.SysEx(data))
	return err
}

type SMFWriter struct {
	wr smf.Writer
	*midiWriter
	finishedTracks uint16
	dest           io.Writer
}

func (s *SMFWriter) SetDelta(deltatime uint32) {
	s.wr.SetDelta(deltatime)
}

func (s *SMFWriter) EndOfTrack() error {

	if no := s.wr.Header().NumTracks; s.finishedTracks >= no {
		return fmt.Errorf("too many tracks: in header: %v, closed: %v", no, s.finishedTracks+1)
	}
	s.finishedTracks++
	_, err := s.wr.Write(meta.EndOfTrack)
	return err
}

func (s *SMFWriter) Copyright(text string) error {
	_, err := s.wr.Write(meta.Copyright(text))
	return err
}

func (s *SMFWriter) Cuepoint(text string) error {
	_, err := s.wr.Write(meta.Cuepoint(text))
	return err
}

func (s *SMFWriter) DevicePort(port string) error {
	_, err := s.wr.Write(meta.DevicePort(port))
	return err
}

func (s *SMFWriter) KeySignature(key uint8, ismajor bool, num uint8, isflat bool) error {
	_, err := s.wr.Write(meta.KeySignature{Key: key, IsMajor: ismajor, Num: num, IsFlat: isflat})
	return err
}

// Key is supposed to be used with the midi/midimessage/meta/key package
func (s *SMFWriter) Key(keysig meta.KeySignature) error {
	_, err := s.wr.Write(keysig)
	return err
}

func (s *SMFWriter) Lyric(text string) error {
	_, err := s.wr.Write(meta.Lyric(text))
	return err
}

func (s *SMFWriter) Marker(text string) error {
	_, err := s.wr.Write(meta.Marker(text))
	return err
}

func (s *SMFWriter) MIDIChannel(ch uint8) error {
	_, err := s.wr.Write(meta.MIDIChannel(ch))
	return err
}

func (s *SMFWriter) MIDIPort(port uint8) error {
	_, err := s.wr.Write(meta.MIDIPort(port))
	return err
}

func (s *SMFWriter) ProgramName(text string) error {
	_, err := s.wr.Write(meta.ProgramName(text))
	return err
}

func (s *SMFWriter) Sequence(text string) error {
	_, err := s.wr.Write(meta.Sequence(text))
	return err
}

func (s *SMFWriter) SequenceNumber(no uint16) error {
	_, err := s.wr.Write(meta.SequenceNumber(no))
	return err
}

func (s *SMFWriter) SequencerSpecific(data []byte) error {
	_, err := s.wr.Write(meta.SequencerSpecific(data))
	return err
}

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

func (s *SMFWriter) Tempo(bpm uint32) error {
	_, err := s.wr.Write(meta.Tempo(bpm))
	return err
}

func (s *SMFWriter) Text(text string) error {
	_, err := s.wr.Write(meta.Text(text))
	return err
}

func (s *SMFWriter) Meter(numerator, denominator uint8) error {
	_, err := s.wr.Write(meter.Meter(numerator, denominator))
	return err
}

// If you don't want to deal with clocks per click and demisemiquaverperquarter,
// user the Meter method instead
func (s *SMFWriter) TimeSignature(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) error {
	_, err := s.wr.Write(meta.TimeSignature{
		Numerator:                numerator,
		Denominator:              denominator,
		ClocksPerClick:           clocksPerClick,
		DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
	})
	return err
}

func (s *SMFWriter) Track(track string) error {
	_, err := s.wr.Write(meta.Track(track))
	return err
}

// panics if numtracks is == 0
func NewSMFWriter(dest io.Writer, numtracks uint16) *SMFWriter {
	if numtracks == 0 {
		panic("numtracks must be > 0")
	}
	wr := smfwriter.New(dest, smfwriter.NumTracks(numtracks), smfwriter.TimeFormat(smf.MetricTicks(0)))
	return &SMFWriter{
		dest:       dest,
		wr:         wr,
		midiWriter: &midiWriter{wr: wr, ch: channel.Ch0},
	}
}

// Creates a new SMF file and allows writer to write to it.
// The file is guaranteed to be closed when returning.
// The last track is closed automatically, if needed.
func NewSMFFile(file string, numtracks uint16, writer func(*SMFWriter) error) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	wr := NewSMFWriter(f, numtracks)
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

// NewLiveWriter creates and new LiveWriter
func NewLiveWriter(dest io.Writer) *LiveWriter {
	wr := midiwriter.New(dest, midiwriter.NoRunningStatus())
	return &LiveWriter{&midiWriter{wr: wr, ch: channel.Ch0}}
}

/*
func init() {
	var bf bytes.Buffer
	wr := NewLiveWriter(&bf)
	wr.SetChannel(2)
	wr.NoteOn(102, 80)
	time.Sleep(time.Second)
	wr.NoteOff(102)
}

func init() {
	err := NewSMFFile("test.mid", 2, func(wr *SMFWriter) (err error) {
		for {
			err = wr.Copyright("just a test")

			if err != nil {
				break
			}

			err = wr.Tempo(144)

			if err != nil {
				break
			}

			err = wr.TimeSignature(3, 4)

			if err != nil {
				break
			}

			wr.SetDelta(12)

			err = wr.NoteOn(60, 110)

			if err != nil {
				break
			}

			wr.SetDelta(32)

			err = wr.NoteOff(60)

			if err != nil {
				break
			}

			break // leave at the end
		}
		return
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

*/
