package midiwriter

import (
	"github.com/gomidi/midi"
	"io"
)

// New returns a new midi.Writer.
//
// The Writer does no buffering and makes no attempt to close dest.
func New(dest io.Writer) midi.Writer {
	return &writer{dest}
}

type writer struct {
	output io.Writer
}

// WriteEvent writes the header on the first call, if e.writeHeader is true
// in realtime mode, no header and no track is written, instead each event is
// written as is to the output writer until an end of track event had come
// then io.EOF is returned
// WriteEvent returns any writing error or io.EOF if the last track has been written
func (w *writer) Write(msg midi.Message) (err error) {
	_, err = w.output.Write(msg.Raw())
	return err
}
