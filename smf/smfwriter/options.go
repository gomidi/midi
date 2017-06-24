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

func QuarterNoteTicks(ticks uint16) Option {
	return func(e *writer) {
		e.header.TickHeader = resQuarterNote(ticks)
	}
}

func SMPTE24(ticksPerFrame int8) Option {
	return func(e *writer) {
		e.header.TickHeader = resSmpteFrames{24, ticksPerFrame}
	}
}

func SMPTE25(ticksPerFrame int8) Option {
	return func(e *writer) {
		e.header.TickHeader = resSmpteFrames{25, ticksPerFrame}
	}
}

func SMPTE30DropFrame(ticksPerFrame int8) Option {
	return sMPTE29(ticksPerFrame)
}

func sMPTE29(ticksPerFrame int8) Option {
	return func(e *writer) {
		e.header.TickHeader = resSmpteFrames{29, ticksPerFrame}
	}
}

func SMPTE30(ticksPerFrame int8) Option {
	return func(e *writer) {
		e.header.TickHeader = resSmpteFrames{30, ticksPerFrame}
	}
}

func NumTracks(ntracks uint16) Option {
	return func(e *writer) {
		e.header.NumTracks = ntracks
	}
}

func SMF2() Option {
	return func(e *writer) {
		e.header.MidiFormat = smf.SequentialTracks.Number()
	}
}

func SMF1() Option {
	return func(e *writer) {
		e.header.MidiFormat = smf.MultiTrack.Number()
	}
}

func SMF0() Option {
	return func(e *writer) {
		e.header.MidiFormat = smf.SingleTrack.Number()
	}
}
