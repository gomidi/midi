# midi

library for reading and writing of MIDI messages and SMF/MIDI files with Go.

Note: If you are reading this on Github, please note that the repo has moved to Gitlab (gitlab.com/gomidi/midi) and this is only a mirror.

- Go version: >= 1.16
- OS/architectures: everywhere Go runs (tested on Linux and Windows).

## Installation

```
go get gitlab.com/gomidi/midi/v2@latest
```

## Features

This package provides a unified way to read and write "over the wire" MIDI data and MIDI files (SMF).

- [x] implementation of complete MIDI standard ("cable" and SMF MIDI)
- [x] unified Driver interface (see below)
- [x] reading and optional writing with "running status"
- [x] seamless integration with io.Reader and io.Writer
- [x] no dependencies outside the standard library
- [x] typed Messages 
- [x] low overhead 
- [x] shortcuts for General MIDI, Sysex messages etc.

## Drivers

For "cable" communication you need a `Driver`to connect with the MIDI system of your OS.
Currently the following drivers available in the drivers subdirectory (all multi-platform):
- rtmididrv based on rtmidi (requires CGO)
- portmididrv based on portmidi (requires CGO)
- webmididrv based on the Web MIDI standard (produces webassembly)
- midicatdrv based on the midicat binaries via piping (stdin / stdout) (no CGO needed)
- testdrv for testing (no CGO needed)

### Example

```go
package main

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/gm"
	"gitlab.com/gomidi/midi/v2/smf"

	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	for _, o := range midi.OutPorts() {
		fmt.Printf("out: %s\n", o)
	}

	out := midi.FindOutPort("qsynth")
	if out < 0 {
		fmt.Printf("can't find qsynth")
		return
	}

	// create a SMF
	rd := bytes.NewReader(mkSMF())

	// read and play it
	smf.ReadTracksFrom(rd).Do(func(ev smf.TrackEvent) {
		fmt.Printf("track %v @%vms %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message)
		
		// use the ev.Message.Get* methods to get the data
	}).Play(out)
}

// makes a SMF and returns the bytes
func mkSMF() []byte {
	var (
		bf    bytes.Buffer
		clock = smf.MetricTicks(96) // resolution: 96 ticks per quarternote 960 is also common
		tr    smf.Track
	)

	// first track must have tempo and meter informations
	tr.Add(0, smf.MetaMeter(3, 4))
	tr.Add(0, smf.MetaTempo(140))
	tr.Add(0, smf.MetaInstrument("Brass"))
	tr.Add(0, midi.ProgramChange(0, gm.Instr_BrassSection.Value()))
	tr.Add(0, midi.NoteOn(0, 76, 120))
	tr.Add(clock.Ticks8th(), midi.NoteOn(0, 83, 120))
	// duration: a quarter note (96 ticks in our case)
	tr.Add(clock.Ticks4th()*2, midi.NoteOff(0, 76))
	tr.Add(0, midi.NoteOff(0, 83))
	tr.Close(0)

	// create the SMF and add the tracks
	s := smf.New()
	s.TimeFormat = clock
	s.Add(tr)
	s.WriteTo(&bf)
	return bf.Bytes()
}

```




## Documentation

[![Go Reference](https://pkg.go.dev/badge/gitlab.com/gomidi/midi/v2.svg)](https://pkg.go.dev/gitlab.com/gomidi/midi/v2)

## License

MIT (see LICENSE file) 

## Credits

Inspiration and low level code for MIDI reading (see internal midilib package) came from the http://github.com/afandian/go-midi package of Joe Wass which also helped as a starting point for the reading of SMF files.

## Alternatives

Matt Aimonetti is also working on MIDI inside https://github.com/mattetti/audio but I didn't try it.
