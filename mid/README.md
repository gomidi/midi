# mid
Porcelain library for reading and writing MIDI and SMF (Standard MIDI File) 

## Description

Package mid provides an easy abstraction for reading and writing of "live" `MIDI` and `SMF` 
(Standard MIDI File) data.

`MIDI` data could be written the following ways:

- `NewWriter` is used to write "live" MIDI to an `io.Writer`.
- `NewSMF` is used to write SMF MIDI to an `io.Writer`.
- `NewSMFFile` is used to write a complete SMF file.
- `ConnectOut` creates a writer that writes "over the wire" to a MIDI out port

To read, create a `Reader` and attach callbacks to it.
Then MIDI data could be read the following ways:

- `Reader.ReadAllFrom` reads "live" MIDI from an `io.Reader`.
- `Reader.ReadAllSMF` reads SMF MIDI from an `io.Reader`.
- `Reader.ReadSMFFile` reads a complete SMF file.
- `ConnectIn` reads "over the wire" from a MIDI in port into the given reader

To connect to external MIDI devices, you need a driver implementing the `Driver` interface. There are currently two drivers available: 
- rtmidi based at `gitlab.com/gomidi/rtmididrv`
- portmidi based at `gitlab.com/gomidi/portmididrv`

## Example with external MIDI gear

```go
package main

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/rtmididrv"
)

// This example reads from the first input and and writes to the first output port
func main() {
	drv, err := rtmididrv.New()
	must(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs()
	must(err)

	// if the string `list` is passed as an argument,
	// print the available in and out ports
	if len(os.Args) == 2 && os.Args[1] == "list" {
		printInPorts(ins)
		printOutPorts(outs)
		return
	}

	in, out := ins[0], outs[0]

	must(in.Open())
	must(out.Open())

	defer in.Close()
	defer out.Close()

	// the writer we are writing to
	wr := mid.ConnectOut(out)

	// to disable logging, pass mid.NoLogger() as option
	rd := mid.NewReader()

	// write every message to the out port
	rd.Msg.Each = func(pos *mid.Position, msg midi.Message) {
		wr.Write(msg)
	}

	// listen for MIDI
	mid.ConnectIn(in, rd)
}

func printPort(port mid.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts(ports []mid.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func printOutPorts(ports []mid.Out) {
	fmt.Printf("MIDI OUT Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
```

For a simple example with "live" MIDI and `io.Reader` and `io.Writer` see the example below.

## Example with io.Writer and io.Reader

We use an `io.Writer` to write to and `io.Reader` to read from. They are connected by the same `io.Pipe`.

```go
package main

import (
	"fmt"
	"io"
	"time"

	"gitlab.com/gomidi/midi/mid"
)

// callback for note on messages
func noteOn(p *mid.Position, channel, key, vel uint8) {
	fmt.Printf("NoteOn (ch %v: key %v vel: %v)\n", channel, key, vel)
}

// callback for note off messages
func noteOff(p *mid.Position, channel, key, vel uint8) {
	fmt.Printf("NoteOff (ch %v: key %v)\n", channel, key)
}

func main() {
	fmt.Println()

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
		time.Sleep(time.Second) // let the note ring for 1 sec
		wr.NoteOff(120)
		pipewr.Close() // finishes the writing
	}()

	for {
		if rd.ReadAllFrom(piperd) == io.EOF {
			piperd.Close() // finishes the reading
			break
		}
	}

	// Output:
	// channel.NoteOn channel 11 key 120 velocity 50
	// NoteOn (ch 11: key 120 vel: 50)
	// channel.NoteOff channel 11 key 120
	// NoteOff (ch 11: key 120)
}
```
