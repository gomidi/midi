package midiio

import (
	"bytes"
	"io"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/realtime"
	"github.com/gomidi/midi/midireader"
)

/*
use case:
we want to write bytes that are MIDI data and want to translate them to
midi.Messages that are passed to a custom midi.Writer
NewWriter allows us to pass our custom midi.Writer and returns an io.Writer
that we can write the bytes to.
*/
func NewWriter(to midi.Writer) io.Writer {
	p := &iowriter{}
	p.to = to
	p.from = midireader.New(&p.bf, p.writeRealtime)
	return p
}

type iowriter struct {
	bf   bytes.Buffer
	from midi.Reader
	to   midi.Writer
}

func (p *iowriter) writeRealtime(msg realtime.Message) {
	p.bf.Write(msg.Raw())
}

// each write does in fact write to the midi.Writer passed to new
func (p *iowriter) Write(data []byte) (n int, err error) {
	_, err = p.bf.Write(data)
	if err != nil {
		return
	}

	var msg midi.Message
	msg, err = p.from.Read()

	if err != nil {
		return
	}

	p.bf.Reset()
	return p.to.Write(msg)
}
