# midi

Modular library for reading and writing of MIDI messages and MIDI files with Go.

Note: If you are reading this on Github, please note that the repo has moved to Gitlab (gitlab.com/gomidi/midi) and this is only a mirror.

## Status

stable

- Go version: >= 1.12
- OS/architectures: everywhere Go runs (tested on Linux and Windows).

## Installation

```
go get gitlab.com/gomidi/midi@latest
```

## Features

This package provides a unified way to read and write "over the wire" MIDI data and MIDI files (SMF).

- [x] implementation of complete MIDI standard (live and SMF data)
- [x] reading and optional writing with "running status"
- [x] seemless integration with io.Reader and io.Writer
- [x] allows the reuse of same libraries for live writing and writing to SMF files
- [x] provide building blocks for other MIDI libraries and applications
- [x] stable API
- [x] no dependencies outside the standard library
- [x] small modular core packages
- [x] typed Messages 

## Non-Goals

- [ ] constructing of MIDI time code messages
- [ ] Multidimensional Polyphonic Expression (MPE)
- [ ] dealing with the inner structure of sysex messages
- [ ] CLI tools


## Drivers

For "over the wire" communication you need a `Driver`to connect with the MIDI system of your OS.
Currently there are two multi plattform drivers available:
- package `gitlab.com/gomidi/rtmididrv` based on rtmidi
- package `gitlab.com/gomidi/portmididrv` based on portmidi

## Porcelain package

For easy access, the porcelain package `gitlab.com/gomidi/midi/mid` is recommended.

The other packages are more low level and allow you to write your own implementations of the `midi.Reader` and `midi.Writer` interfaces
to wrap the given SMF and live readers/writers for your own application.

[Documentation porcelain package](https://pkg.go.dev/gitlab.com/gomidi/midi/mid)


### Example with external MIDI gear

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

### Example with io.Writer and io.Reader

A simple example with "live" MIDI and `io.Reader` and `io.Writer`.

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

## Low level packages

[Documentation main package](https://pkg.go.dev/gitlab.com/gomidi/midi)


### Example with low level packages

```go
package main

import (
    "bytes"
    "fmt"
    . "gitlab.com/gomidi/midi/midimessage/channel"
    "gitlab.com/gomidi/midi/midimessage/realtime"
    "gitlab.com/gomidi/midi/midireader"
    "gitlab.com/gomidi/midi/midiwriter"
    "io"
    "gitlab.com/gomidi/midi"
)

func main() {
    var bf bytes.Buffer

    wr := midiwriter.New(&bf)
    wr.Write(Channel2.Pitchbend(5000))
    wr.Write(Channel2.NoteOn(65, 90))
    wr.Write(realtime.Reset)
    wr.Write(Channel2.NoteOff(65))

    rthandler := func(m realtime.Message) {
        fmt.Printf("Realtime: %s\n", m)
    }

    rd := midireader.New(bytes.NewReader(bf.Bytes()), rthandler)

    var m midi.Message
    var err error

    for {
        m, err = rd.Read()

        // breaking at least with io.EOF
        if err != nil {
            break
        }

        // inspect
        fmt.Println(m)

        switch v := m.(type) {
        case NoteOn:
            fmt.Printf("NoteOn at channel %v: key: %v velocity: %v\n", v.Channel(), v.Key(), v.Velocity())
        case NoteOff:
            fmt.Printf("NoteOff at channel %v: key: %v\n", v.Channel(), v.Key())
        }

    }

    if err != io.EOF {
        panic("error: " + err.Error())
    }

    // Output:
    // channel.Pitchbend channel 2 value 5000 absValue 13192
    // channel.NoteOn channel 2 key 65 velocity 90
    // NoteOn at channel 2: key: 65 velocity: 90
    // Realtime: Reset
    // channel.NoteOff channel 2 key 65
    // NoteOff at channel 2: key: 65
}

```

### Modularity

Apart from the porcelain package there are small subpackages, so that you only need to import
what you really need.

This keeps packages and dependencies small, better testable and should result in a smaller memory footprint which should help smaller devices.

For reading and writing of live and SMF MIDI data io.Readers are accepted as input and io.Writers as output. Furthermore there are common interfaces for live and SMF MIDI data handling: midi.Reader and midi.Writer. The typed MIDI messages used in each case are the same.

To connect with MIDI libraries expecting and returning plain bytes (e.g. over the wire), use `midiio` subpackage.

## License

MIT (see LICENSE file) 

## Credits

Inspiration and low level code for MIDI reading (see internal midilib package) came from the http://github.com/afandian/go-midi package of Joe Wass which also helped as a starting point for the reading of SMF files.

## Alternatives

Matt Aimonetti is also working on MIDI inside https://github.com/mattetti/audio but I didn't try it.
