package midiio

import (
	"bytes"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/live/midiwriter"
	"io"
)

/*
use case:
we want to read bytes that are MIDI data and come from a custom midi.Reader
with each read there should come the data of a single MIDI message.
*/
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

func (p *ioreader) Read(data []byte) (n int, err error) {
	msg, err := p.rd.Read()

	if err != nil {
		return
	}

	// midiwriter writes the running status
	_, err = p.to.Write(msg)

	if err != nil {
		return
	}

	return p.bf.Read(data)
}
