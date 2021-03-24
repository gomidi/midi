# midi

Modular library for reading and writing of MIDI messages and MIDI files with Go.

Note: If you are reading this on Github, please note that the repo has moved to Gitlab (gitlab.com/gomidi/midi) and this is only a mirror.

- Go version: >= 1.14
- OS/architectures: everywhere Go runs (tested on Linux and Windows).

## Installation

```
go get gitlab.com/gomidi/midi@latest
```

## Features

This package provides a unified way to read and write "over the wire" MIDI data and MIDI files (SMF).

- [x] implementation of complete MIDI standard ("cable" and SMF MIDI)
- [x] reading and optional writing with "running status"
- [x] seamless integration with io.Reader and io.Writer
- [x] allows the reuse of same libraries for live writing and writing to SMF files
- [x] provide building blocks for other MIDI libraries and applications
- [x] no dependencies outside the standard library
- [x] small modular core packages
- [x] typed Messages 

## Drivers

For "cable" communication you need a `Driver`to connect with the MIDI system of your OS.
Currently the following drivers available  (all multi-platform):
- package `gitlab.com/gomidi/rtmididrv` based on rtmidi (requires CGO)
- package `gitlab.com/gomidi/portmididrv` based on portmidi (requires CGO)
- package `gitlab.com/gomidi/webmididrv` based on the Web MIDI standard (produces webassembly)
- package `gitlab.com/gomidi/midicatdrv` based on the midicat binaries via piping (stdin / stdout) (no CGO needed)
- package `gitlab.com/gomidi/midi/testdrv` for testing (no CGO needed)

## Projects using this library

- [miti](https://pkg.go.dev/github.com/schollz/miti)
- [launchpad](https://pkg.go.dev/github.com/rainu/launchpad)
- [pamidicontrol](https://pkg.go.dev/github.com/solarnz/pamidicontrol)
- [bleep](https://pkg.go.dev/github.com/bspaans/bleep)
- [midi-macro](https://pkg.go.dev/github.com/vipul-sharma20/midi-macro)
- [hyperarp](https://pkg.go.dev/gitlab.com/gomidi/hyperarp)
- [golinnstrument](https://pkg.go.dev/gitlab.com/golinnstrument/linnstrument)
- [smfimage](https://pkg.go.dev/gitlab.com/gomidi/smfimage)
- [midispy](https://pkg.go.dev/gitlab.com/gomidi/midispy)

## Porcelain package

For easy access, the following packages are recommended:

- reading: `gitlab.com/gomidi/midi/reader` [![Go Reference](https://pkg.go.dev/badge/gitlab.com/gomidi/midi/reader.svg)](https://pkg.go.dev/gitlab.com/gomidi/midi/reader)
- writing: `gitlab.com/gomidi/midi/writer`  [![Go Reference](https://pkg.go.dev/badge/gitlab.com/gomidi/midi/writer.svg)](https://pkg.go.dev/gitlab.com/gomidi/midi/writer)

The other packages are more low level and allow you to write your own implementations of the `midi.Reader`, `midi.Writer`and `midi.Driver` interfaces to wrap the given SMF and live readers/writers/drivers for your own application.

### Example with MIDI cables

```go
package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"

	// replace with e.g. "gitlab.com/gomidi/rtmididrv" for real midi connections
	driver "gitlab.com/gomidi/midi/testdrv"
	"gitlab.com/gomidi/midi/writer"
)

// This example reads from the first input and and writes to the first output port
func main() {
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
```

### Example with MIDI file (SMF)


```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
)

type printer struct{}

func (pr printer) noteOn(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Track: %v Pos: %v NoteOn (ch %v: key %v vel: %v)\n", p.Track, p.AbsoluteTicks, channel, key, vel)
}

func (pr printer) noteOff(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Track: %v Pos: %v NoteOff (ch %v: key %v)\n", p.Track, p.AbsoluteTicks, channel, key)
}

func main() {
	dir := os.TempDir()
	f := filepath.Join(dir, "smf-test.mid")

	defer os.Remove(f)

	var p printer

	err := writer.WriteSMF(f, 2, func(wr *writer.SMF) error {
		
		wr.SetChannel(11) // sets the channel for the next messages
		writer.NoteOn(wr, 120, 50)
		wr.SetDelta(120)
		writer.NoteOff(wr, 120)
		
		wr.SetDelta(240)
		writer.NoteOn(wr, 125, 50)
		wr.SetDelta(20)
		writer.NoteOff(wr, 125)
		writer.EndOfTrack(wr)
		
		wr.SetChannel(2)
		writer.NoteOn(wr, 120, 50)
		wr.SetDelta(60)
		writer.NoteOff(wr, 120)
		writer.EndOfTrack(wr)
		return nil
	})

	if err != nil {
		fmt.Printf("could not write SMF file %v\n", f)
		return
	}

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(reader.NoLogger(),
		// set the functions for the messages you are interested in
		reader.NoteOn(p.noteOn),
		reader.NoteOff(p.noteOff),
	)

	err = reader.ReadSMFFile(rd, f)

	if err != nil {
		fmt.Printf("could not read SMF file %v\n", f)
	}

	// Output: Track: 0 Pos: 0 NoteOn (ch 11: key 120 vel: 50)
	// Track: 0 Pos: 120 NoteOff (ch 11: key 120)
	// Track: 0 Pos: 360 NoteOn (ch 11: key 125 vel: 50)
	// Track: 0 Pos: 380 NoteOff (ch 11: key 125)
	// Track: 1 Pos: 0 NoteOn (ch 2: key 120 vel: 50)
	// Track: 1 Pos: 60 NoteOff (ch 2: key 120)
}
```

## Low level packages

[![Go Reference](https://pkg.go.dev/badge/gitlab.com/gomidi/midi.svg)](https://pkg.go.dev/gitlab.com/gomidi/midi)


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

Apart from the porcelain packages there are small subpackages, so that you only need to import
what you really need.

This keeps packages and dependencies small, better testable and should result in a smaller memory footprint which should help smaller devices.

For reading and writing of cable and SMF MIDI data `io.Readers` are accepted as input and `io.Writers` as output. Furthermore there are common interfaces for live and SMF MIDI data handling: `midi.Reader` and `midi.Writer`. The typed MIDI messages used in each case are the same.

To connect with MIDI libraries expecting and returning plain bytes, use the `midiio` subpackage.

## Stability / API / Semantic Versioning

[This excellent blog post by Peter Bourgon](https://peter.bourgon.org/blog/2020/09/14/siv-is-unsound.html) describes
perfectly the problem, I have with `Go modules`' take on `semantic versioning`.

First of all, the API of this library is large (if you use all of it, which is unlikely) and generally
stable. There may be small incompatibilities from time to time that will only affect a small amount of users and
that will "blow up into your face"  when compiling. These are easy to fix with the help of the compiler
and the documentation. 

The `diamond dependency problem` should be unlikely with this library, since that would mean, 
you are using a library that is abstracting over MIDI and uses this library in order to do so. 
Well, then you should reconsider the use of that library, since there really is no point in further abstractions. 
For other uses however (e.g. integration),  the interfaces should be used and they won't change. 

As a last resort, you could fork that other library, add a replace statement to your `go.mod`
and fix the issue, followed by a pull request. If that libary requiring the old code is not maintained anymore, 
you would need to fork anyway in the future. If it is maintained, your pull request should get accepted.

**We are not going to bump the major version with every tiny incompatible change of this code base.
The costs are too high and the benefits too low. Instead, any `minor` version upgrade reflects incompatible changes and
the `patch`versions contain compatible changes (i.e. also compatible additions). A large rewrite will however end up in a major version bump.**

## License

MIT (see LICENSE file) 

## Credits

Inspiration and low level code for MIDI reading (see internal midilib package) came from the http://github.com/afandian/go-midi package of Joe Wass which also helped as a starting point for the reading of SMF files.

## Alternatives

Matt Aimonetti is also working on MIDI inside https://github.com/mattetti/audio but I didn't try it.
