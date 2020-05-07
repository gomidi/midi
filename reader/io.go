package reader

import (
	"bytes"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/smf"
)

type inReader struct {
	rd *Reader
	//in         midi.In
	midireader midi.Reader
	bf         bytes.Buffer
}

func (r *inReader) Read() (midi.Message, error) {
	return r.midireader.Read()
}

func (r *inReader) handleMessage(b []byte, deltaMicroseconds int64) {
	// use the fake position to get the ticks for the current tempo
	r.rd.pos = &Position{}
	r.rd.pos.DeltaTicks = Ticks(r.rd, time.Duration(deltaMicroseconds*1000)) // deltaticks
	r.rd.pos.AbsoluteTicks += uint64(r.rd.pos.DeltaTicks)
	r.bf.Write(b)
	r.rd.dispatchMessageFromReader()
}

// Duration returns the duration for the given delta ticks, respecting the current tempo
func Duration(r *Reader, deltaticks uint32) time.Duration {
	return r.resolution.FractionalDuration(r.TempoBPM(), deltaticks)
}

// Resolution returns the ticks of a quarternote
// If it can't be determined, 0 is returned
func Resolution(r *Reader) uint32 {
	if r.resolution == 0 {
		return 0
	}
	return r.resolution.Ticks4th()
}

// Ticks returns the ticks that correspond to a duration while respecting the current tempo
// If it can't be determined, 0 is returned
func Ticks(r *Reader, d time.Duration) uint32 {
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

// ListenTo connects the reader to the given MIDI in connection and listens for incoming messages.
func (r *Reader) ListenTo(in midi.In) error {
	r.resolution = LiveResolution
	r.reset()
	//rd := &inReader{rd: r, in: in}
	rd := &inReader{rd: r}
	rd.midireader = midireader.New(&rd.bf, r.dispatchRealTime, r.midiReaderOptions...)
	r.reader = rd
	//return rd.in.SetListener(rd.handleMessage)
	return in.SetListener(rd.handleMessage)
}

// LiveResolution is the resolution used for live over the wire reading with Reader.ReadFrom
const LiveResolution = smf.MetricTicks(1920)

//const LiveResolution = smf.MetricTicks(960)
