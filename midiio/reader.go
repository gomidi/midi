package midiio

import (
	"bytes"
	"io"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midiwriter"
)

/*
use case:
we want to read bytes that are MIDI data and come from a custom midi.Reader
with each read there should come the data of a single MIDI message.
*/

// NewReader returns an io.Reader that can be read of by third party
// libraries that expect midi as bytes. When the data is read, it
// gets the midi as typed midi.Messages from the given midi.Reader
func NewReader(from midi.Reader) io.Reader {
	p := &ioreader{}
	p.rd = from
	// we need the writer for the running status
	p.to = midiwriter.New(&p.bf)
	return p
}

type ioreader struct {
	bf bytes.Buffer
	to midi.Writer
	rd midi.Reader
}

// Read reads the typed midi.Messages from the midi.Reader that have
// been passed to NewReader and returns them as bytes.
func (p *ioreader) Read(data []byte) (n int, err error) {
	msg, err := p.rd.Read()

	if err != nil {
		return
	}

	// midiwriter writes the running status
	err = p.to.Write(msg)

	if err != nil {
		return
	}

	return p.bf.Read(data)
}
