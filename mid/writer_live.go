package mid

import (
	"io"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midimessage/syscommon"
	"gitlab.com/gomidi/midi/midiwriter"
)

// Writer writes live MIDI data. Its methods must not be called concurrently
type Writer struct {
	*midiWriter
}

var _ midi.Writer = &Writer{}

// NewWriter creates and new Writer for writing of "live" MIDI data ("over the wire")
// By default it makes no use of the running status.
func NewWriter(dest io.Writer, options ...midiwriter.Option) *Writer {
	options = append(
		[]midiwriter.Option{
			midiwriter.NoRunningStatus(),
		}, options...)

	wr := midiwriter.New(dest, options...)
	return &Writer{&midiWriter{wr: wr, Channel: channel.Channel0}}
}

// ActiveSensing writes the active sensing realtime message
func (w *Writer) Activesense() error {
	return w.midiWriter.wr.Write(realtime.Activesense)
}

// Continue writes the continue realtime message
func (w *Writer) Continue() error {
	return w.midiWriter.wr.Write(realtime.Continue)
}

// Reset writes the reset realtime message
func (w *Writer) Reset() error {
	return w.midiWriter.wr.Write(realtime.Reset)
}

// Start writes the start realtime message
func (w *Writer) Start() error {
	return w.midiWriter.wr.Write(realtime.Start)
}

// Stop writes the stop realtime message
func (w *Writer) Stop() error {
	return w.midiWriter.wr.Write(realtime.Stop)
}

// Tick writes the tick realtime message
func (w *Writer) Tick() error {
	return w.midiWriter.wr.Write(realtime.Tick)
}

// Clock writes the timing clock realtime message
func (w *Writer) Clock() error {
	return w.midiWriter.wr.Write(realtime.TimingClock)
}

// MTC writes the MIDI Timing Code system message
func (w *Writer) MTC(code uint8) error {
	return w.midiWriter.wr.Write(syscommon.MTC(code))
}

// SPP writes the song position pointer system message
func (w *Writer) SPP(ptr uint16) error {
	return w.midiWriter.wr.Write(syscommon.SPP(ptr))
}

// SongSelect writes the song select system message
func (w *Writer) SongSelect(song uint8) error {
	return w.midiWriter.wr.Write(syscommon.SongSelect(song))
}

// Tune writes the tune request system message
func (w *Writer) Tune() error {
	return w.midiWriter.wr.Write(syscommon.Tune)
}
