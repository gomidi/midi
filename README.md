# midi
Reading and writing of MIDI messages with Go.

## Goals

- [ ] implementation of complete MIDI standard (including SMF)
- [x] common ground for live MIDI processing and processing of Standard MIDI Files (SMF)
- [ ] correctness
- [ ] high test coverage
- [ ] stable API
- [x] usage of small interfaces preferably of the standard library
- [x] provide building blocks for MIDI applications
- [ ] performance
- [x] no dependencies outside the standard library
- [x] as little dependencies to the standard library as possible
- [x] small modular packages (see below)
- [ ] beginner friendlyness
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

## Documentation

see http://godoc.org/github.com/gomidi/midi

## Status

[![Build Status](https://travis-ci.org/gomidi/midi.svg?branch=master)](http://travis-ci.org/gomidi/midi)

alpha (usable, but API may change)

## License

MIT (see LICENSE file) 
