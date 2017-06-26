# midi
reading and writing of MIDI data

## Goals

- implement complete MIDI standard (including SMF)
- common ground for live MIDI processing and processing of Standard MIDI Files (SMF)
- correctness and high test coverage
- stable API
- usage of small interfaces when possible of the standard library
- provide building blocks for MIDI applications
- performance
- no dependencies outside the standard library
- as little dependencies to the standard library as possible
- small modular packages (see below)
- beginner friendlyness
- quality documentation

## Modularity

This package is divided into small subpackages, so that you only need to import
what you really need. This keeps packages and dependencies small, better testable and should result in a smaller memory footprint which should help smaller devices.

Also it allows for small interfaces in libraries reusing this building blocks.

For reading and writing of live and SMF MIDI data io.Readers are accepted as input and io.Writers as output. Furthermore there is are common interfaces for live and SMF MIDI data handling: midi.Reader and midi.Writer. The typed MIDI messages used in each case are the same.

## Documentation

see http://godoc.org/github.com/gomidi/midi

## Status

[![Build Status](https://travis-ci.org/gomidi/midi.svg?branch=master)](http://travis-ci.org/gomidi/midi)

alpha (usable, but expect bugs and incompatible changes in the API)

## License

MIT (see LICENSE file) 
