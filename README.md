# midi
Reading and writing of MIDI messages with Go.

## Goals

- implementation of complete MIDI standard (live and SMF data)
- provide building blocks for MIDI applications
- stable API
- no dependencies outside the standard library
- small modular packages (see below)
- pure Go library (no C, no assembler) 

## Non-Goals

- constructing of MIDI time code messages
- dealing with the inner structure of sysex messages
- connection to MIDI devices (combine this lib with http://github.com/rakyll/portmidi or http://github.com/scgolang/midi )
- CLI tools

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

## Documentation

see http://godoc.org/github.com/gomidi/midi

## Discussion

If you have questions that won't fit into tickets, larger ideas etc. please join the Google group at https://groups.google.com/forum/#!forum/gomidi

## Status

beta - API may change

The package and the subpackages are not officially stable yet.

However there are some ideas about how stable the subpackages are.
These ideas are reflected in the following table as a basis for communication
about how firm the current API is supposed to be. 

In other words: 

- If you have issues with an API that is declared as "stable"
in this table, please speak up soon, because there will probably be some important discussions that also effects other packages.  

- If you have issues with an API that is not declared as "stable" please have a look at the open issues that explain the current status and help to get it stable.

**Users of this package should ignore the following table, because as long as the whole package is in beta status, anything may change any time**.

However we want this package to come out of beta soon and the open issues are mainly about comfort and naming.

    package                   API stable          complete
    ----------------------------------------------------
    midiwriter                yes                 yes
    midireader                yes                 yes
    smf                       yes                 yes
    smf/smfwriter             yes                 yes
    smf/smfreader             yes                 yes
    midimessage/channel       almost              yes
    midimessage/cc            yes                 yes
    midimessage/meta          almost              yes
    midimessage/realtime      yes                 yes
    midimessage/syscommon     yes                 yes
    midimessage/sysex         no                  yes
    midiio                    yes                 yes
    
	------- porcelain packages -------
    smf/smftrack              no                  no
    midihandler               no                  yes


[![Build Status](https://travis-ci.org/gomidi/midi.svg?branch=master)](http://travis-ci.org/gomidi/midi)

[![Go Report](https://goreportcard.com/badge/github.com/gomidi/midi)](https://goreportcard.com/report/github.com/gomidi/midi)

- Supported Go versions: 1.4-1.8.
- Supported OS/architecture: Should work on all OS/architectures that Go supports (is tested on Linux, but no OS specific code).

## License

MIT (see LICENSE file) 

## Credits

Inspiration and low level code for MIDI reading (see internal midilib package) came from the http://github.com/afandian/go-midi package of Joe Wass which also helped as a starting point for the reading of SMF files.

## Alternatives

Matt Aimonetti is also working on MIDI inside https://github.com/mattetti/audio but I didn't try it.
