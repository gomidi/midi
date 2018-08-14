package midiwriter

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/internal/runningstatus"
	"io"
)

// New returns a new midi.Writer.
//
// The Writer does no buffering and makes no attempt to close dest.
//
// By default the writer uses running status for efficiency.
// You can disable that behaviour by passing the NoRunningStatus() option.
// If you don't know what running status is, keep the default.
func New(dest io.Writer, opts ...Option) (wr midi.Writer) {
	var c = &config{}

	for _, opt := range opts {
		opt(c)
	}

	if c.noRunningStatus {
		wr = &notRunningWriter{output: dest}
	} else {
		wr = &runningWriter{
			runningstatus: runningstatus.NewLiveWriter(dest),
		}
	}

	return wr
}

type notRunningWriter struct {
	output io.Writer
}

// Write writes a midi.Message to a midi (live) stream.
// It does no caching and makes no use of running status.
func (w *notRunningWriter) Write(msg midi.Message) (err error) {
	_, err = w.output.Write(msg.Raw())
	return
}

type runningWriter struct {
	runningstatus runningstatus.Writer
}

// Write writes a midi.Message to a midi (live) stream.
// It does no caching but makes use of running status.
func (w *runningWriter) Write(msg midi.Message) (err error) {
	_, err = w.runningstatus.Write(msg.Raw())
	return
}
