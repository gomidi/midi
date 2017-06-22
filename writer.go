package midi

import (
	"io"

	"github.com/gomidi/midi/meta"
)

type Writer interface {
	// Write writes the given midi event
	// To finish writing, write a meta.EndOfTrack event
	Write(event Event) error
}

// NewWriter returns a new writer.
//
// Finish the writing by passing an meta.EndOfTrack or let dest return an error.
// The Writer does no buffering and makes no attempt to close dest.
func NewWriter(dest io.Writer) Writer {
	return writeLive(dest)
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
	wr io.Writer
}

// WriteTo writes a midi file to writer
// Pass NumTracks to write multiple tracks (SMF1), otherwise everything will be written
// into a single track (SMF0). However SMF1 can also be enforced with a single track by passing SMF1 as an option
func writeLive(wr io.Writer) *writer {
	return &writer{
		wr: wr,
	}

}

// WriteEvent writes the header on the first call, if e.writeHeader is true
// in realtime mode, no header and no track is written, instead each event is
// written as is to the output writer until an end of track event had come
// then io.EOF is returned
// WriteEvent returns any writing error or io.EOF if the last track has been written
func (e *writer) Write(ev Event) (err error) {
	if ev == meta.EndOfTrack {
		return io.EOF
	}

	_, err = e.wr.Write(ev.Raw())
	return err
}
