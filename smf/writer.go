package smf

import (
	"io"

	midi "github.com/gomidi/midi"
	"github.com/gomidi/midi/meta"
)

type WriteOption func(*writer)

// NewWriter returns a Writer
func NewWriter(dest io.Writer, opts ...WriteOption) Writer {
	return newWriter(dest, opts...)
}

type Writer interface {
	midi.Writer
	SetDelta(deltatime uint32)
}

func SetQuarterNoteTicks(ticks uint16) WriteOption {
	return func(e *writer) {
		e.header.TickHeader = quarterNote(ticks)
	}
}

func SMPTE24(ticksPerFrame int8) WriteOption {
	return func(e *writer) {
		e.header.TickHeader = smpteFrames{24, ticksPerFrame}
	}
}

func SMPTE25(ticksPerFrame int8) WriteOption {
	return func(e *writer) {
		e.header.TickHeader = smpteFrames{25, ticksPerFrame}
	}
}

func SMPTE30DropFrame(ticksPerFrame int8) WriteOption {
	return sMPTE29(ticksPerFrame)
}

func sMPTE29(ticksPerFrame int8) WriteOption {
	return func(e *writer) {
		e.header.TickHeader = smpteFrames{29, ticksPerFrame}
	}
}

func SMPTE30(ticksPerFrame int8) WriteOption {
	return func(e *writer) {
		e.header.TickHeader = smpteFrames{30, ticksPerFrame}
	}
}

func NumTracks(ntracks uint16) WriteOption {
	return func(e *writer) {
		e.header.NumTracks = ntracks
	}
}

func SMF2() WriteOption {
	return func(e *writer) {
		e.header.MidiFormat = SequentialTracks
	}
}

func SMF1() WriteOption {
	return func(e *writer) {
		e.header.MidiFormat = MultiTrack
	}
}

func SMF0() WriteOption {
	return func(e *writer) {
		e.header.MidiFormat = SingleTrack
	}
}

/*
from http://www.artandscienceofsound.com/article/standardmidifiles

Depending upon the application you are using to create the file in the first place, header information may automatically be saved from within parameters set in the application, or may need to be placed in a ‘set-up’ bar before the music data commences.

Either way, information that should be considered includes:

GM/GS Reset message

Per MIDI Channel
Bank Select (0=GM) / Program Change #
Reset All Controllers (not all devices may recognize this command so you may prefer to zero out or reset individual controllers)
Initial Volume (CC7) (standard level = 100)
Expression (CC11) (initial level set to 127)
Hold pedal (0 = off)
Pan (Center = 64)
Modulation (0)
Pitch bend range
Reverb (0 = off)
Chorus level (0 = off)

System Exclusive data

If RPNs or more detailed controller messages are being employed in the file these should also be reset or normalized in the header.

If you are inputting header data yourself it is advisable not to clump all such information together but rather space it out in intervals of 5-10 ticks. Certainly if a file is designed to be looped, having too much data play simultaneously will cause most playback devices to ‘choke, ’ and throw off your timing.
*/

type writer struct {
	header          *header
	wr              io.Writer
	qticks          uint16
	AutoEndTrack    bool
	writeHeader     bool
	currentTrack    *track
	tracksProcessed uint16
	deltatime       uint32
}

// WriteTo writes a midi file to writer
// Pass NumTracks to write multiple tracks (SMF1), otherwise everything will be written
// into a single track (SMF0). However SMF1 can also be enforced with a single track by passing SMF1 as an option
func newWriter(wr io.Writer, opts ...WriteOption) *writer {
	enc := &writer{
		header: &header{
			MidiFormat: format(10), // not existing, only for checking if it is undefined to be able to set the default
		},
		writeHeader:  true,
		wr:           wr,
		currentTrack: &track{},
	}

	for _, opt := range opts {
		opt(enc)
	}

	if enc.header.NumTracks == 0 {
		enc.header.NumTracks = 1
	}

	// if midiformat is undefined (see above), i.e. not set via options
	// set the default, which is format 0 for one track and format 1 for multitracks
	if enc.header.MidiFormat == format(10) {
		if enc.header.NumTracks == 1 {
			enc.header.MidiFormat = SingleTrack
		} else {
			enc.header.MidiFormat = MultiTrack
		}
	}

	if enc.header.TickHeader == nil {
		enc.header.TickHeader = quarterNote(960)
	}

	if qn, is := enc.header.TickHeader.(quarterNote); is {
		enc.qticks = qn.Ticks()
	}

	return enc
}

func (e *writer) SetDelta(deltatime uint32) {
	e.deltatime = deltatime
}

// WriteEvent writes the header on the first call, if e.writeHeader is true
// in realtime mode, no header and no track is written, instead each event is
// written as is to the output writer until an end of track event had come
// then io.EOF is returned
// WriteEvent returns any writing error or io.EOF if the last track has been written
func (e *writer) Write(ev midi.Event) (err error) {
	defer func() {
		e.deltatime = 0
	}()

	if e.writeHeader {
		err = e.WriteHeader()
		if err != nil {
			return err
		}
		e.writeHeader = false
	}
	// fmt.Printf("%T\n", ev)
	if ev == meta.EndOfTrack {
		e.currentTrack.Add(e.deltatime, ev.Raw())
		_, err = e.currentTrack.WriteTo(e.wr)
		e.tracksProcessed++
		if e.header.NumTracks == e.tracksProcessed {
			return io.EOF
		}
		e.currentTrack = &track{}
		return nil
	}
	e.currentTrack.Add(e.deltatime, ev.Raw())
	return nil
}

func (e *writer) WriteHeader() (err error) {
	_, err = e.header.WriteTo(e.wr)
	return
}
