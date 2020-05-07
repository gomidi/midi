package writer

import (
	"io"
	"os"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/smf"
	"gitlab.com/gomidi/midi/smf/smftimeline"
	"gitlab.com/gomidi/midi/smf/smfwriter"
)

var _ smf.Writer = &SMF{}

// SMFWriter writes SMF MIDI data. Its methods must not be called concurrently
type SMF struct {
	wr smf.Writer
	*Writer
	finishedTracks uint16
	//dest           io.Writer
	smf.MetricTicks
	timeline *smftimeline.TimeLine
	delta    uint32
}

func (wr *SMF) Delta() uint32 {
	return wr.delta
}

// SetDelta sets the delta ticks to the next message
// It should mostly not be needed, use Forward instead to advance in musical time.
func (w *SMF) SetDelta(deltatime uint32) {
	w.delta = deltatime
	w.wr.SetDelta(deltatime)
}

func (wr *SMF) Header() smf.Header {
	return wr.wr.Header()
}

func (wr *SMF) WriteHeader() error {
	return wr.wr.WriteHeader()
}

// Position returns absolute position of the last written message in ticks
func (w *SMF) Position() uint64 {
	return w.wr.Position()
}

// WrapSMF returns a new SMF for a given smf.Writer
// The TimeFormat of the smf.Writer must be metric or this function will panic.
func WrapSMF(wr smf.Writer) *SMF {
	smfwr := &SMF{
		wr:     wr,
		Writer: &Writer{wr: wr, channel: 0},
	}

	metr, isMetric := wr.Header().TimeFormat.(smf.MetricTicks)

	if !isMetric {
		panic("timeformat must be metric")
	}
	smfwr.MetricTicks = metr
	smfwr.timeline = smftimeline.New(metr)

	return smfwr
}

// NewSMF returns a new SMF that writes to dest.
// It panics if numtracks is == 0.
func NewSMF(dest io.Writer, numtracks uint16, options ...smfwriter.Option) *SMF {
	if numtracks == 0 {
		panic("numtracks must be > 0")
	}

	options = append(
		[]smfwriter.Option{
			smfwriter.NumTracks(numtracks),
			smfwriter.TimeFormat(smf.MetricTicks(960)),
		}, options...)

	return WrapSMF(smfwriter.New(dest, options...))
}

func BackupTimeline(wr *SMF) {
	wr.timeline.Backup()
}

func RestoreTimeline(wr *SMF) {
	wr.timeline.Restore()
}

// WriteSMF creates a new SMF file and allows writer to write to it.
// The file is guaranteed to be closed when returning.
// The last track is closed automatically, if needed.
// It panics if numtracks is == 0.
func WriteSMF(file string, numtracks uint16, writer func(*SMF) error, options ...smfwriter.Option) error {
	if numtracks == 0 {
		panic("numtracks must be > 0")
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	wr := NewSMF(f, numtracks, options...)
	if writer != nil {
		err = writer(wr)
		if err != nil {
			return err
		}
	}

	if no := wr.wr.Header().NumTracks; wr.finishedTracks < no {
		err := EndOfTrack(wr)
		if err != nil && err != smf.ErrFinished {
			return err
		}
	}

	return nil
}

// Forward sets the cursor based on the given number of bars and ratio of whole notes.
// The cursor is the current position where the next event will be inserted. In the background
// it sets the delta to the next event. The cursor can only move forward.
//
// Examples:
//
// To move the cursor to the 2nd next bar (respecting time signature changes), use
//   Forward(2,0,0)
// To move the cursor by 23 8ths (independent from time signatures), use
//   Forward(0,23,8)
// To move the cursor to the 3rd 4th of the next bar (respecting time signature changes), use
//   Forward(1,3,4)
//
// Important notes:
//   1. Always put time signature changes at the beginning of a bar.
//   2. Never forward more than once without setting a event in between.
func Forward(w *SMF, nbars, num, denom uint32) {
	w.timeline.Forward(nbars, num, denom)
	delta := w.timeline.GetDelta()
	if delta < 0 {
		panic("cursor before last delta, must not happen")
	}
	w.SetDelta(uint32(delta))
}

// Plan plans the given midi.Message at the given position. That leads to the message being written
// when the Forward method is crossing the corresponding position
func Plan(w *SMF, nbars, num, denom uint32, msg midi.Message) {
	w.timeline.Plan(nbars, num, denom, func(delta int32) {
		w.SetDelta(uint32(delta))
		w.Write(msg)
	})
}

// FinishPlanned finishes the planned midi.Messages
func FinishPlanned(w *SMF) {
	w.timeline.FinishPlanned()
}
