package smfwriter

import (
	"github.com/gomidi/midi/internal/runningstatus"
	"io"
	"os"

	"github.com/gomidi/midi"

	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
)

// WriteFile creates file, calls callback with a writer and closes file
// WriteFile makes sure that the data of the last track is written by sending
// an meta.EndOfTrack message after callback has been run.
// So callback may skip the sending of the last meta.EndOfTrack message although
// it does no harm if send twice. Especially for single track (SMF0) files this
// is interesting, since no meta.EndOfTrack message must then be send from callback.
func WriteFile(file string, callback func(smf.Writer), options ...Option) error {
	f, err := os.Create(file)

	if err != nil {
		return err
	}

	defer func() {
		f.Close()
	}()

	wr := New(f, options...)
	callback(wr)

	// make sure the data of the last track is written
	wr.Write(meta.EndOfTrack)

	return nil
}

// New returns a Writer
// If no options are passed, a single track midi file is written (SMT0).
// Each track must be finished with a meta.EndOfTrack message.
func New(dest io.Writer, opts ...Option) smf.Writer {
	return newWriter(dest, opts...)
}

type writer struct {
	header *header
	wr     io.Writer
	qticks uint16
	// AutoEndTrack    bool
	writeHeader     bool
	currentTrack    *track
	tracksProcessed uint16
	deltatime       uint32
	noRunningStatus bool
}

// WriteTo writes a midi file to writer
// Pass NumTracks to write multiple tracks (SMF1), otherwise everything will be written
// into a single track (SMF0). However SMF1 can also be enforced with a single track by passing SMF1 as an option
func newWriter(output io.Writer, opts ...Option) *writer {
	enc := &writer{
		header: &header{
		// MidiFormat: format(10), // not existing, only for checking if it is undefined to be able to set the default
		},
		writeHeader:  true,
		currentTrack: &track{},
	}

	for _, opt := range opts {
		opt(enc)
	}

	enc.wr = output

	if !enc.noRunningStatus {
		enc.currentTrack.runningWriter = runningstatus.NewSMFWriter()
	}

	if enc.header.NumTracks == 0 {
		enc.header.NumTracks = 1
	}

	// if midiformat is undefined (see above), i.e. not set via options
	// set the default, which is format 0 for one track and format 1 for multitracks
	// if enc.header.MidiFormat == format(10) {
	if enc.header.MidiFormat != smf.SMF2 && enc.header.NumTracks > 1 {
		enc.header.MidiFormat = smf.SMF1
	}
	// }

	if enc.header.TimeFormat == nil {
		enc.header.TimeFormat = smf.QuarterNoteTicks(960)
	}

	if qn, is := enc.header.TimeFormat.(smf.QuarterNoteTicks); is {
		enc.qticks = qn.Ticks()
	}

	return enc
}

func (e *writer) SetDelta(deltatime uint32) {
	e.deltatime = deltatime
}

// Write writes a midi message to the SMF file.
// Due to the nature of SMF files there is some maybe surprising behavior.
// - If the header has not been written yet, it will be written before writing the first message.
// - The first message will be written to track 0 which will be implicetly created.
// - All messages of a track will be buffered inside the track and only be written if an EndOfTrack
//   message is written.
// - The number of tracks that are written will never execeed the NumTracks that have been defined as
//   an option. If the last track has been written, io.EOF will be returned. (Also for any further attempt to write).
// - It is the responsability of the caller to make sure the provided NumTracks (which defaults to 1) is not
//   larger as the number of tracks in the file. smfreader is tolerant when reading such a file; so may be other
//   SMF readers.
// - It is the responsability of the caller to open and close any file where appropriate. The writer just uses an io.Writer.
// Keep the above in mind when examinating the written nbytes that are returned. They reflect the number of bytes
// that have been physically written.
func (e *writer) Write(m midi.Message) (nbytes int, err error) {
	defer func() {
		e.deltatime = 0
	}()

	if e.header.NumTracks == e.tracksProcessed {
		err = io.EOF
		return
	}

	if e.writeHeader {
		nbytes, err = e.WriteHeader()
		if err != nil {
			return
		}
		e.writeHeader = false
	}
	// fmt.Printf("%T\n", ev)
	if m == meta.EndOfTrack {
		e.currentTrack.Add(e.deltatime, m)
		var tnum int
		tnum, err = e.currentTrack.WriteTo(e.wr)
		nbytes += tnum
		e.tracksProcessed++
		if e.header.NumTracks == e.tracksProcessed {
			err = io.EOF
			return
		}
		e.currentTrack = &track{}

		if !e.noRunningStatus {
			e.currentTrack.runningWriter = runningstatus.NewSMFWriter()
		}
		return
	}
	e.currentTrack.Add(e.deltatime, m)
	return
}

func (e *writer) WriteHeader() (nbytes int, err error) {
	return e.header.WriteTo(e.wr)
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
