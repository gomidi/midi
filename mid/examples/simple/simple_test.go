package simple_test

import (
	"fmt"
	"io"
	"time"

	"gitlab.com/gomidi/midi/mid"
)

func noteOn(p *mid.Position, channel, key, vel uint8) {
	fmt.Printf("NoteOn (ch %v: key %v vel: %v)\n", channel, key, vel)
}

func noteOff(p *mid.Position, channel, key, vel uint8) {
	fmt.Printf("NoteOff (ch %v: key %v)\n", channel, key)
}

func Example() {
	// to disable logging, pass mid.NoLogger() as option
	rd := mid.NewReader()

	// set the functions for the messages you are interested in
	rd.Msg.Channel.NoteOn = noteOn
	rd.Msg.Channel.NoteOff = noteOff

	// to allow reading and writing concurrently in this example
	// we need a pipe
	piperd, pipewr := io.Pipe()

	go func() {
		wr := mid.NewWriter(pipewr)
		wr.SetChannel(11) // sets the channel for the next messages
		wr.NoteOn(120, 50)
		time.Sleep(time.Second)
		wr.NoteOff(120) // let the note ring for 1 sec
		pipewr.Close()  // finishes the writing
	}()

	for {
		if rd.ReadAllFrom(piperd) == io.EOF {
			piperd.Close() // finishes the reading
			break
		}
	}

	// Output: channel.NoteOn channel 11 key 120 velocity 50
	// NoteOn (ch 11: key 120 vel: 50)
	// channel.NoteOff channel 11 key 120
	// NoteOff (ch 11: key 120)
}
