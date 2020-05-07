package cables_test

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"

	// replace with e.g. "gitlab.com/gomidi/midi/rtmididrv" for real midi connections
	driver "gitlab.com/gomidi/midi/testdrv"
	"gitlab.com/gomidi/midi/writer"
)

// This example reads from the first input and and writes to the first output port
func Example() {
	// you would take a real driver here e.g. rtmididrv.New()
	drv := driver.New("fake cables: messages written to output port 0 are received on input port 0")

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs()
	must(err)

	in, out := ins[0], outs[0]

	must(in.Open())
	must(out.Open())

	defer in.Close()
	defer out.Close()

	// the writer we are writing to
	wr := writer.New(out)

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		reader.NoLogger(),
		// write every message to the out port
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Printf("got %s\n", msg)
		}),
	)

	// listen for MIDI
	err = rd.ListenTo(in)
	must(err)

	err = writer.NoteOn(wr, 60, 100)
	must(err)

	time.Sleep(1)
	err = writer.NoteOff(wr, 60)

	must(err)
	// Output: got channel.NoteOn channel 0 key 60 velocity 100
	// got channel.NoteOff channel 0 key 60
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
