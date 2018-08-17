package smfwriter

import (
	"github.com/gomidi/midi/smf"
)

/*
TODO
defaults:
	- store sysex if they are written
	- ignore realtime and System Common or System Real-Time messages if they are written
	- take deltas set via SetDelta
	- no quantization/rounding
options:
  - ignore incoming sysex
  - store incoming realtime and System Common or System Real-Time messages (escape them inside a sysex)
  - ignore SetDelta and measure time instead (autogenerating delta on the way)
  - allow delta quantization/rounding

*/

// Option is a Writer option
type Option func(*writer)

// NoRunningStatus forces the writer to always write the status byte.
// Without passing this option, running status will be used if possible (saving some bytes).
func NoRunningStatus() Option {
	return func(w *writer) {
		w.noRunningStatus = true
	}
}

// TimeFormat sets the timeformat. Allowed values are smf.MetricTicks and smf.TimeCode
// Without passing this option or when timeformat is nil, smf.MetricTicks(960) will be used.
func TimeFormat(timeformat smf.TimeFormat) Option {
	if timeformat == nil {
		timeformat = smf.MetricTicks(0)
	}
	return func(w *writer) {
		w.header.TimeFormat = timeformat
	}
}

// NumTracks sets the number of tracks in the file.
// Due to the SMF format and for performance reasons,
// the number of tracks must be given to Writer, before the MIDI events could be written.
// If the number of tracks is not given - or 0 - , it defaults to 1 track.
// If the given number of tracks has been written, any further writing returns an io.EOF error.
func NumTracks(ntracks uint16) Option {
	if ntracks == 0 {
		ntracks = 1
	}
	return func(w *writer) {
		w.header.NumTracks = ntracks
	}
}

// Format sets the SMF file format version.
// Valid values are: smf.SMF0 (single track), smf.SMF1 (multi track), smf.SMF2 (sequential track)
// If this option is not given, SMF0 will be used as default if the number of tracks is 1, otherwise SMF1.
func Format(f smf.Format) Option {
	return func(w *writer) {
		w.header.Format = f
	}
}
