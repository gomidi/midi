package midiwriter

import (
	"fmt"
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

	if c.checkMessageType {
		return &checkWriter{wr}
	}

	if c.ignoreWrongMessageType {
		return &skipNonLiveWriter{wr}
	}
	return wr
}

type liveMessage interface {
	IsLiveMessage()
}

type skipNonLiveWriter struct {
	midi.Writer
}

// Write checks if msg is valid for live usage before writing
func (w *skipNonLiveWriter) Write(msg midi.Message) (int, error) {
	if _, ok := msg.(liveMessage); !ok {
		return 0, nil
	}

	return w.Writer.Write(msg)
}

type checkWriter struct {
	midi.Writer
}

// Write checks if msg is valid for live usage before writing
func (w *checkWriter) Write(msg midi.Message) (int, error) {
	if _, ok := msg.(liveMessage); !ok {
		return 0, fmt.Errorf("not a MIDI live message: %s", msg)
	}

	return w.Writer.Write(msg)
}

type notRunningWriter struct {
	output io.Writer
}

// Write writes a midi.Message to a midi (live) stream.
// It does no caching and makes no use of running status.
func (w *notRunningWriter) Write(msg midi.Message) (int, error) {
	return w.output.Write(msg.Raw())
}

type runningWriter struct {
	runningstatus runningstatus.Writer
}

// Write writes a midi.Message to a midi (live) stream.
// It does no caching but makes use of running status.
func (w *runningWriter) Write(msg midi.Message) (int, error) {
	return w.runningstatus.Write(msg.Raw())
}
