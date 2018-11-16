package mid

import (
	"gitlab.com/gomidi/midi/midimessage/meta/meter"
	// "bytes"
	// "encoding/binary"
	"fmt"
	// "github.com/gomidi/midi/internal/midilib"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"

	// "github.com/gomidi/midi/midimessage/realtime"
	// "github.com/gomidi/midi/midimessage/syscommon"
	"io"
	"os"

	"gitlab.com/gomidi/midi/smf"
	"gitlab.com/gomidi/midi/smf/smfwriter"
	// "time"
)

// SMFWriter writes SMF MIDI data. Its methods must not be called concurrently
type SMFWriter struct {
	wr smf.Writer
	*midiWriter
	finishedTracks uint16
	dest           io.Writer
	smf.MetricTicks
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

	wr := smfwriter.New(dest, options...)

	smfwr := &SMFWriter{
		dest:       dest,
		wr:         wr,
		midiWriter: &midiWriter{wr: wr, ch: channel.Channel0},
	}

	if metr, isMetric := wr.Header().TimeFormat.(smf.MetricTicks); isMetric {
		smfwr.MetricTicks = metr
	}
	return smfwr
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
func (w *SMFWriter) SetDelta(deltatime uint32) {
	w.wr.SetDelta(deltatime)
}

// EndOfTrack signals the end of a track
func (w *SMFWriter) EndOfTrack() error {
	w.midiWriter.noteState = [16][128]bool{}
	if no := w.wr.Header().NumTracks; w.finishedTracks >= no {
		return fmt.Errorf("too many tracks: in header: %v, closed: %v", no, w.finishedTracks+1)
	}
	w.finishedTracks++
	return w.wr.Write(meta.EndOfTrack)
}

// Copyright writes the copyright meta message
func (w *SMFWriter) Copyright(text string) error {
	return w.wr.Write(meta.Copyright(text))
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

// Sequence writes the sequence (name) meta message
func (w *SMFWriter) Sequence(text string) error {
	return w.wr.Write(meta.Sequence(text))
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
// Numerator and Denominator are decimalw.
func (w *SMFWriter) Meter(numerator, denominator uint8) error {
	return w.wr.Write(meter.Meter(numerator, denominator))
}

// TimeSig writes the time signature meta message.
// Numerator and Denominator are decimalw.
// If you don't want to deal with clocks per click and demisemiquaverperquarter,
// user the Meter method instead.
func (w *SMFWriter) TimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter uint8) error {
	return w.wr.Write(meta.TimeSig{
		Numerator:                numerator,
		Denominator:              denominator,
		ClocksPerClick:           clocksPerClick,
		DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
	})
}

// Track writes the track name aka instrument name meta message
func (w *SMFWriter) Track(track string) error {
	return w.wr.Write(meta.Track(track))
}
