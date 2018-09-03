# midi
Modular library for reading and writing of MIDI messages with Go.

[![Build Status Travis/Linux](https://travis-ci.org/gomidi/midi.svg?branch=master)](http://travis-ci.org/gomidi/midi) [![Build Status AppVeyor/Windows](https://ci.appveyor.com/api/projects/status/408nwdlo2b1lwdd1?svg=true)](https://ci.appveyor.com/project/metakeule/midi) [![Coverage Status](https://coveralls.io/repos/github/gomidi/midi/badge.svg?branch=master)](https://coveralls.io/github/gomidi/midi?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gomidi/midi)](https://goreportcard.com/report/github.com/gomidi/midi) [![Documentation](http://godoc.org/github.com/gomidi/midi?status.png)](http://godoc.org/github.com/gomidi/midi)

**This package is meant for users that have some knowledge of the MIDI standards. Beginners, and people on the run might want to look at the porcelain package https://github.com/gomidi/mid.**

## Status

stable

- Go version: >= 1.10
- OS/architectures: everywhere Go runs (tested on Linux and Windows).

## Installation

```
go get github.com/gomidi/midi/...
```

## Documentation

see http://godoc.org/github.com/gomidi/midi

## Features

- [x] implementation of complete MIDI standard (live and SMF data)
- [x] reading and optional writing with "running status"
- [x] seemless integration with io.Reader and io.Writer
- [x] allows the reuse of same libraries for live writing and writing to SMF files
- [x] provide building blocks for other MIDI libraries and applications
- [x] stable API
- [x] no dependencies outside the standard library
- [x] small modular core packages
- [x] typed Messages 
- [x] pure Go library (no C, no assembler) 

## Non-Goals

- [ ] constructing of MIDI time code messages
- [ ] Multidimensional Polyphonic Expression (MPE)
- [ ] dealing with the inner structure of sysex messages
- [ ] connection to MIDI devices (for this combine it with https://github.com/gomidi/connect)
- [ ] CLI tools

## Usage / Example

```go
package main

import (
    "bytes"
    "fmt"
    . "github.com/gomidi/midi/midimessage/channel"
    "github.com/gomidi/midi/midimessage/realtime"
    "github.com/gomidi/midi/midireader"
    "github.com/gomidi/midi/midiwriter"
    "io"
    "github.com/gomidi/midi"
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



## Modularity

This package is divided into small subpackages, so that you only need to import
what you really need. This keeps packages and dependencies small, better testable and should result in a smaller memory footprint which should help smaller devices.

For reading and writing of live and SMF MIDI data io.Readers are accepted as input and io.Writers as output. Furthermore there are common interfaces for live and SMF MIDI data handling: midi.Reader and midi.Writer. The typed MIDI messages used in each case are the same.

To connect with MIDI libraries expecting and returning plain bytes (e.g. over the wire), use `midiio` subpackage.

## Perfomance

On my laptop, writing noteon and noteoff ("live")

    BenchmarkSameChannel            123 ns/op  12 B/op  3 allocs/op
    BenchmarkAlternatingChannel     123 ns/op  12 B/op  3 allocs/op
    BenchmarkRunningStatusDisabled  110 ns/op  12 B/op  3 allocs/op

On my laptop, reading noteon and noteoff ("live")
("Samechannel" makes use of running status byte).

    BenchmarkSameChannel            351 ns/op  12 B/op  7 allocs/op
    BenchmarkAlternatingChannel     425 ns/op  14 B/op  9 allocs/op


## License

MIT (see LICENSE file) 

## Credits

Inspiration and low level code for MIDI reading (see internal midilib package) came from the http://github.com/afandian/go-midi package of Joe Wass which also helped as a starting point for the reading of SMF files.

## Alternatives

Matt Aimonetti is also working on MIDI inside https://github.com/mattetti/audio but I didn't try it.
