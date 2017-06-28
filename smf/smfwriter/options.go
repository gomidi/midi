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

type Option func(*writer)

func NoRunningStatus() Option {
	return func(w *writer) {
		w.noRunningStatus = true
	}
}

// TimeFormat sets the timeformat to either smf.QuarterNoteTicks or smf.TimeCode
func TimeFormat(timeformat smf.TimeFormat) Option {
	if timeformat == nil {
		panic("timeformat must not be nil")
	}
	return func(w *writer) {
		w.TimeFormat = timeformat
	}
}

func NumTracks(ntracks uint16) Option {
	if ntracks == 0 {
		panic("ntracks must not be 0")
	}
	return func(w *writer) {
		w.NumTracks = ntracks
	}
}

func Format(f smf.Format) Option {
	return func(w *writer) {
		w.Format = f
	}
}
