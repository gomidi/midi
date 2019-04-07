package mid

import (
	"fmt"
	"io"
	"os"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/midi/midimessage/meta/meter"
	"gitlab.com/gomidi/midi/smf"
	"gitlab.com/gomidi/midi/smf/smftimeline"
	"gitlab.com/gomidi/midi/smf/smfwriter"
)

var _ smf.Writer = &SMFWriter{}

// SMFWriter writes SMF MIDI data. Its methods must not be called concurrently
type SMFWriter struct {
	wr smf.Writer
	*midiWriter
	finishedTracks uint16
	//dest           io.Writer
	smf.MetricTicks
	timeline *smftimeline.TimeLine
	delta    uint32
}

func (wr *SMFWriter) Delta() uint32 {
	return wr.delta
}

func (wr *SMFWriter) BackupTimeline() {
	wr.timeline.Backup()
}

func (wr *SMFWriter) RestoreTimeline() {
	wr.timeline.Restore()
}

func (wr *SMFWriter) Header() smf.Header {
	return wr.wr.Header()
}

func (wr *SMFWriter) WriteHeader() error {
	return wr.wr.WriteHeader()
}

// NewSMFWriter returns a new SMFWriter for a given smf.Writer
// The TimeFormat of the smf.Writer must be metric or this function will panic.
func NewSMFWriter(wr smf.Writer) *SMFWriter {
	smfwr := &SMFWriter{
		wr:         wr,
		midiWriter: &midiWriter{wr: wr, Channel: channel.Channel0},
	}

	metr, isMetric := wr.Header().TimeFormat.(smf.MetricTicks)

	if !isMetric {
		panic("timeformat must be metric")
	}
	smfwr.MetricTicks = metr
	smfwr.timeline = smftimeline.New(metr)

	return smfwr
}

// NewSMF returns a new SMFWriter that writes to dest.
// It panics if numtracks is == 0.
func NewSMF(dest io.Writer, numtracks uint16, options ...smfwriter.Option) *SMFWriter {
	if numtracks == 0 {
		panic("numtracks must be > 0")
	}

	options = append(
		[]smfwriter.Option{
			smfwriter.NumTracks(numtracks),
			smfwriter.TimeFormat(smf.MetricTicks(960)),
		}, options...)

	return NewSMFWriter(smfwriter.New(dest, options...))
}

// NewSMFFile creates a new SMF file and allows writer to write to it.
// The file is guaranteed to be closed when returning.
// The last track is closed automatically, if needed.
// It panics if numtracks is == 0.
func NewSMFFile(file string, numtracks uint16, writer func(*SMFWriter) error, options ...smfwriter.Option) error {
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
		err := wr.EndOfTrack()
		if err != nil && err != smf.ErrFinished {
			return err
		}
	}

	return nil
}

// SetDelta sets the delta ticks to the next message
// It should mostly not be needed, use Forward instead to advance in musical time.
func (w *SMFWriter) SetDelta(deltatime uint32) {
	w.delta = deltatime
	w.wr.SetDelta(deltatime)
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
func (w *SMFWriter) Forward(nbars, num, denom uint32) {
	w.timeline.Forward(nbars, num, denom)
	delta := w.timeline.GetDelta()
	if delta < 0 {
		panic("cursor before last delta, must not happen")
	}
	w.SetDelta(uint32(delta))
}

// Position returns absolute position of the last written message in ticks
func (w *SMFWriter) Position() uint64 {
	return w.wr.Position()
}

// Plan plans the given midi.Message at the given position. That leads to the message being written
// when the Forward method is crossing the corresponding position
func (w *SMFWriter) Plan(nbars, num, denom uint32, msg midi.Message) {
	w.timeline.Plan(nbars, num, denom, func(delta int32) {
		w.SetDelta(uint32(delta))
		w.Write(msg)
	})
}

// FinishPlanned finishes the planned midi.Messages
func (w *SMFWriter) FinishPlanned() {
	w.timeline.FinishPlanned()
}

// EndOfTrack signals the end of a track
func (w *SMFWriter) EndOfTrack() error {
	w.midiWriter.noteState = [16][128]bool{}
	if no := w.wr.Header().NumTracks; w.finishedTracks >= no {
		return fmt.Errorf("too many tracks: in header: %v, closed: %v", no, w.finishedTracks+1)
	}
	w.finishedTracks++
	if w.timeline != nil {
		w.timeline.Reset()
	}
	return w.wr.Write(meta.EndOfTrack)
}

// Copyright writes the copyright meta message
func (w *SMFWriter) Copyright(text string) error {
	return w.wr.Write(meta.Copyright(text))
}

// Writes an undefined meta message
func (w *SMFWriter) Undefined(typ byte, bt []byte) error {
	return w.wr.Write(meta.Undefined{typ, bt})
}

// Cuepoint writes the cuepoint meta message
func (w *SMFWriter) Cuepoint(text string) error {
	return w.wr.Write(meta.Cuepoint(text))
}

// Device writes the device port meta message
func (w *SMFWriter) Device(port string) error {
	return w.wr.Write(meta.Device(port))
}

// KeySig writes the key signature meta message.
// A more comfortable way is to use the Key method in conjunction
// with the gomidi/midi/midimessage/meta/key package
func (w *SMFWriter) KeySig(key uint8, ismajor bool, num uint8, isflat bool) error {
	return w.wr.Write(meta.Key{Key: key, IsMajor: ismajor, Num: num, IsFlat: isflat})
}

// Key writes the given key signature meta message.
// It is supposed to be used with the gomidi/midi/midimessage/meta/key package
func (w *SMFWriter) Key(keysig meta.Key) error {
	return w.wr.Write(keysig)
}

// Lyric writes the lyric meta message
func (w *SMFWriter) Lyric(text string) error {
	return w.wr.Write(meta.Lyric(text))
}

// Marker writes the marker meta message
func (w *SMFWriter) Marker(text string) error {
	return w.wr.Write(meta.Marker(text))
}

// DeprecatedChannel writes the deprecated MIDI channel meta message
func (w *SMFWriter) DeprecatedChannel(ch uint8) error {
	return w.wr.Write(meta.Channel(ch))
}

// DeprecatedPort writes the deprecated MIDI port meta message
func (w *SMFWriter) DeprecatedPort(port uint8) error {
	return w.wr.Write(meta.Port(port))
}

// Program writes the program name meta message
func (w *SMFWriter) Program(text string) error {
	return w.wr.Write(meta.Program(text))
}

// TrackSequenceName writes the track / sequence name meta message
// If in a format 0 track, or the first track in a format 1 file, the name of the sequence. Otherwise, the name of the track.
func (w *SMFWriter) TrackSequenceName(name string) error {
	return w.wr.Write(meta.TrackSequenceName(name))
}

// SequenceNo writes the sequence number meta message
func (w *SMFWriter) SequenceNo(no uint16) error {
	return w.wr.Write(meta.SequenceNo(no))
}

// SequencerData writes a custom sequences specific meta message
func (w *SMFWriter) SequencerData(data []byte) error {
	return w.wr.Write(meta.SequencerData(data))
}

// SMPTE writes the SMPTE offset meta message
func (w *SMFWriter) SMPTE(hour, minute, second, frame, fractionalFrame byte) error {
	return w.wr.Write(meta.SMPTE{
		Hour:            hour,
		Minute:          minute,
		Second:          second,
		Frame:           frame,
		FractionalFrame: fractionalFrame,
	})
}

// Tempo writes the tempo meta message
func (w *SMFWriter) TempoBPM(bpm float64) error {
	return w.wr.Write(meta.FractionalBPM(bpm))
}

// Text writes the text meta message
func (w *SMFWriter) Text(text string) error {
	return w.wr.Write(meta.Text(text))
}

// Meter writes the time signature meta message in a more comfortable way.
// Numerator and Denominator are decimals.
func (w *SMFWriter) Meter(numerator, denominator uint8) error {
	w.timeline.AddTimeSignature(numerator, denominator)
	return w.wr.Write(meter.Meter(numerator, denominator))
}

// TimeSig writes the time signature meta message.
// Numerator and Denominator are decimal.
// If you don't want to deal with clocks per click and demisemiquaverperquarter,
// user the Meter method instead.
func (w *SMFWriter) TimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) error {
	w.timeline.AddTimeSignature(numerator, denominator)
	return w.wr.Write(meta.TimeSig{
		Numerator:                numerator,
		Denominator:              denominator,
		ClocksPerClick:           clocksPerClick,
		DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
	})
}

// Instrument writes the instrument name meta message
func (w *SMFWriter) Instrument(name string) error {
	return w.wr.Write(meta.Instrument(name))
}
