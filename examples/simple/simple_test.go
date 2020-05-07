package simple_test

import (
	"fmt"
	"io"
	"time"

	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
)

type printer struct{}

func (pr printer) noteOn(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("NoteOn (ch %v: key %v vel: %v)\n", channel, key, vel)
}

func (pr printer) noteOff(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("NoteOff (ch %v: key %v)\n", channel, key)
}

func Example() {

	var p printer

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(reader.NoLogger(),
		// set the callbacks for the messages you are interested in
		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
	)

	// to allow reading and writing concurrently in this example
	// we need a pipe
	piperd, pipewr := io.Pipe()

	go func() {
		wr := writer.New(pipewr)
		wr.SetChannel(11) // sets the channel for the next messages
		writer.NoteOn(wr, 120, 50)
		time.Sleep(time.Second)
		writer.NoteOff(wr, 120) // let the note ring for 1 sec
		pipewr.Close()          // finishes the writing
	}()

	for {
		if reader.ReadAllFrom(rd, piperd) == io.EOF {
			piperd.Close() // finishes the reading
			break
		}
	}

	// Output: NoteOn (ch 11: key 120 vel: 50)
	// NoteOff (ch 11: key 120)
}
