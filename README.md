# midi
Reading and writing of MIDI messages with Go.

## Goals

- [ ] implementation of complete MIDI standard (including SMF)
- [x] common ground for live MIDI processing and processing of Standard MIDI Files (SMF)
- [x] no known bugs 
- [ ] high test coverage
- [ ] stable API
- [x] usage of small interfaces preferably of the standard library
- [x] provide building blocks for MIDI applications
- [x] no dependencies outside the standard library
- [x] as little dependencies to the standard library as possible
- [x] small modular packages (see below)
- [x] shortcuts for known control changes
- [x] pure Go library (no C, no assembler) 
- [ ] quality documentation

## Non-Goals

- connection to MIDI devices (combine this lib with http://github.com/rakyll/portmidi or http://github.com/scgolang/midi )
- abstractions over the inner meat of sysex messages
- CLI tools (will be in separate package)
- shortcuts for certain devices (belong to separate packages)
- MIDI apps (belong to separate packages)

## Modularity

This package is divided into small subpackages, so that you only need to import
what you really need. This keeps packages and dependencies small, better testable and should result in a smaller memory footprint which should help smaller devices.

Also it allows for small interfaces in libraries reusing this building blocks.

For reading and writing of live and SMF MIDI data io.Readers are accepted as input and io.Writers as output. Furthermore there are common interfaces for live and SMF MIDI data handling: midi.Reader and midi.Writer. The typed MIDI messages used in each case are the same.

## Perfomance

On my laptop, sending 1000 messages (noteon and noteoff; live)

    BenchmarkSameChannel            124805 ns/op  12000 B/op  3000 allocs/op
    BenchmarkAlternatingChannel     123932 ns/op  12000 B/op  3000 allocs/op
    BenchmarkRunningStatusDisabled  113146 ns/op  12000 B/op  3000 allocs/op

On my laptop, reading 1000 messages (noteon and noteoff; live).
("Samechannel" makes use of running status byte).

    BenchmarkSameChannel            362482 ns/op  12000 B/op  7000 allocs/op
    BenchmarkAlternatingChannel     447461 ns/op  14000 B/op  8500 allocs/op

## Documentation

see http://godoc.org/github.com/gomidi/midi

## Status

alpha (usable, but API may change)

The implementation is almost complete (there are some missing meta messages).
Lots of test have to be written. Some open questions about naming of functions and types prevent the API from being stable. Maybe reading performance could be improved. 

[![Build Status](https://travis-ci.org/gomidi/midi.svg?branch=master)](http://travis-ci.org/gomidi/midi)

Supported Go versions: 1.2 to 1.8.
Supported OS/architecture: Should work on all OS/architectures that Go supports (is tested on Linux, but no OS specific code).

## License

MIT (see LICENSE file) 
