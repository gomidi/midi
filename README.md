# midi
Reading and writing of MIDI messages with Go.

[![Build Status Travis/Linux](https://travis-ci.org/gomidi/midi.svg?branch=master)](http://travis-ci.org/gomidi/midi) [![Build Status AppVeyor/Windows](https://ci.appveyor.com/api/projects/status/408nwdlo2b1lwdd1?svg=true)](https://ci.appveyor.com/project/metakeule/midi) [![Coverage Status](https://coveralls.io/repos/github/gomidi/midi/badge.svg)](https://coveralls.io/github/gomidi/midi) [![Go Report](https://goreportcard.com/badge/github.com/gomidi/midi)](https://goreportcard.com/report/github.com/gomidi/midi)

## Goals

- implementation of complete MIDI standard (live and SMF data)
- provide building blocks for MIDI applications
- stable API
- no dependencies outside the standard library
- small modular core packages and a fat comfortable porcelain package
- pure Go library (no C, no assembler) 

## Non-Goals

- constructing of MIDI time code messages
- dealing with the inner structure of sysex messages
- connection to MIDI devices (combine this lib with http://github.com/rakyll/portmidi or http://github.com/scgolang/midi )
- CLI tools

## Usage

Most users will want to use the porcelain package

```
go get -d -t github.com/gomidi/midi/mid
```

Docs: [![Documentation](http://godoc.org/github.com/gomidi/midi/mid?status.png)](http://godoc.org/github.com/gomidi/midi/mid)

You don't need to read here further, if you are not interested in the nitty gritty details. 

## Modularity

This package is divided into small subpackages, so that you only need to import
what you really need. This keeps packages and dependencies small, better testable and should result in a smaller memory footprint which should help smaller devices.

For reading and writing of live and SMF MIDI data io.Readers are accepted as input and io.Writers as output. Furthermore there are common interfaces for live and SMF MIDI data handling: midi.Reader and midi.Writer. The typed MIDI messages used in each case are the same.

## Perfomance

On my laptop, writing noteon and noteoff ("live")

    BenchmarkSameChannel            123 ns/op  12 B/op  3 allocs/op
    BenchmarkAlternatingChannel     123 ns/op  12 B/op  3 allocs/op
    BenchmarkRunningStatusDisabled  110 ns/op  12 B/op  3 allocs/op

On my laptop, reading noteon and noteoff ("live")
("Samechannel" makes use of running status byte).

    BenchmarkSameChannel            351 ns/op  12 B/op  7 allocs/op
    BenchmarkAlternatingChannel     425 ns/op  14 B/op  9 allocs/op

## Full Documentation

see http://godoc.org/github.com/gomidi/midi

## Status

beta - API mostly stable

- Supported Go versions: >= 1.2
- Supported OS/architecture: Should work on all OS/architectures that Go supports (is tested on Linux and Windows, but no OS specific code).


    package                   API stable          complete
    ----------------------------------------------------
    midiwriter                yes                 yes
    midireader                yes                 yes
    smf                       yes                 yes
    smf/smfwriter             yes                 yes
    smf/smfreader             yes                 yes
    midimessage/channel       almost              yes
    midimessage/cc            yes                 yes
    midimessage/meta          yes                 yes
    midimessage/realtime      yes                 yes
    midimessage/syscommon     yes                 yes
    midimessage/sysex         no                  yes
    midiio                    yes                 yes
    
	------- porcelain packages -------
    smf/smftrack              no                  no
    mid                       almost              yes


## License

MIT (see LICENSE file) 

## Credits

Inspiration and low level code for MIDI reading (see internal midilib package) came from the http://github.com/afandian/go-midi package of Joe Wass which also helped as a starting point for the reading of SMF files.

## Alternatives

Matt Aimonetti is also working on MIDI inside https://github.com/mattetti/audio but I didn't try it.
