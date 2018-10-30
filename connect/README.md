# connect
Go interface for MIDI drivers


## Purpose

Unification of MIDI driver packages for Go. Currently two implementations exist: 
- for rtmidi: https://gitlab.com/gomidi/rtmididrv
- for portmidi: https://gitlab.com/gomidi/portmididrv

## Installation

It is recommended to use Go 1.11 with module support (`$GO111MODULE=on`).

For rtmidi (see https://github.com/thestk/rtmidi for more information)

```
// install the headers of alsa somehow, e.g. sudo apt-get install libasound2-dev
go get -d gitlab.com/gomidi/rtmididrv
```

For portaudio (see https://gitlab.com/rakyll/portmidi for more information)

```
// install the headers of portmidi somehow, e.g. sudo apt-get install libportmidi-dev
go get -d gitlab.com/gomidi/portmididrv
```

## Documentation

rtmididrv: [![rtmidi docs](http://godoc.org/gitlab.com/gomidi/rtmididrv?status.png)](http://godoc.org/gitlab.com/gomidi/rtmididrv)

portmididrv: [![portmidi docs](http://godoc.org/gitlab.com/gomidi/portmididrv?status.png)](http://godoc.org/gitlab.com/gomidi/portmididrv)

## Example

```go
package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/connect"
	"gitlab.com/gomidi/midi/mid"
	driver "gitlab.com/gomidi/rtmididrv"
	// for portmidi
	// driver "gitlab.com/gomidi/portmididrv" 
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// This example expects the first input and output port to be connected
// somehow (are either virtual MIDI through ports or physically connected).
// We write to the out port and listen to the in port.
func main() {
	drv, err := driver.New()
	must(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs()
	must(err)

	if len(os.Args) == 2 && os.Args[1] == "list" {
		printInPorts(ins)
		printOutPorts(outs)
		return
	}

	in, out := ins[0], outs[0]

	must(in.Open())
	must(out.Open())

	wr := mid.WriteTo(out)

	// listen for MIDI
	go mid.NewReader().ReadFrom(in)

	{ // write MIDI to out that passes it to in on which we listen.
		err := wr.NoteOn(60, 100)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Nanosecond)
		wr.NoteOff(60)
		time.Sleep(time.Nanosecond)

		wr.SetChannel(1)

		wr.NoteOn(70, 100)
		time.Sleep(time.Nanosecond)
		wr.NoteOff(70)
		time.Sleep(time.Second * 1)
	}
}

func printPort(port connect.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts(ports []connect.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func printOutPorts(ports []connect.Out) {
	fmt.Printf("MIDI OUT Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

```
