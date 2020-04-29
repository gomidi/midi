package mid

import (
	"bytes"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/smf"
)

type outWriter struct {
	out midi.Out
}

func (w *outWriter) Write(b []byte) (int, error) {
	return len(b), w.out.Send(b)
}

// ConnectOut returns a Writer that writes to the given MIDI out connection.
func ConnectOut(out midi.Out) *Writer {
	return NewWriter(&outWriter{out})
}

type inReader struct {
	rd         *Reader
	in         midi.In
	midireader midi.Reader
	bf         bytes.Buffer
}

func (r *inReader) Read() (midi.Message, error) {
	return r.midireader.Read()
}

func (r *inReader) handleMessage(b []byte, deltaMicroseconds int64) {
	// use the fake position to get the ticks for the current tempo
	r.rd.pos = &Position{}
	r.rd.pos.DeltaTicks = r.rd.Ticks(time.Duration(deltaMicroseconds * 1000)) // deltaticks
	r.rd.pos.AbsoluteTicks += uint64(r.rd.pos.DeltaTicks)
	r.bf.Write(b)
	r.rd.dispatchMessageFromReader()
}

// Duration returns the duration for the given delta ticks, respecting the current tempo
func (r *Reader) Duration(deltaticks uint32) time.Duration {
	return r.resolution.FractionalDuration(r.TempoBPM(), deltaticks)
}

// Resolution returns the ticks of a quarternote
// If it can't be determined, 0 is returned
func (r *Reader) Resolution() uint32 {
	if r.resolution == 0 {
		return 0
	}
	return r.resolution.Ticks4th()
}

// Ticks returns the ticks that correspond to a duration while respecting the current tempo
// If it can't be determined, 0 is returned
func (r *Reader) Ticks(d time.Duration) uint32 {
	if r.resolution == 0 {
		return 0
	}
	return r.resolution.FractionalTicks(r.TempoBPM(), d)
}

/*
// Tempo returns the current tempo in BPM (beats per minute)
func (r *Reader) Tempo() uint32 {
	tempochange := r.tempoChanges[len(r.tempoChanges)-1]
	return tempochange.bpm
}
*/

// BPM returns the current tempo in BPM (beats per minute)
func (r *Reader) TempoBPM() float64 {
	tempochange := r.tempoChanges[len(r.tempoChanges)-1]
	return tempochange.bpm
}

// ConnectIn connects the reader to the given MIDI in connection and starts reading from it.
func ConnectIn(in midi.In, r *Reader) error {
	r.resolution = LiveResolution
	r.reset()
	rd := &inReader{rd: r, in: in}
	rd.midireader = midireader.New(&rd.bf, r.dispatchRealTime, r.midiReaderOptions...)
	r.reader = rd
	return rd.in.SetListener(rd.handleMessage)
}

// LiveResolution is the resolution used for live over the wire reading with Reader.ReadFrom
const LiveResolution = smf.MetricTicks(1920)

//const LiveResolution = smf.MetricTicks(960)
