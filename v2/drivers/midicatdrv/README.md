# midicatdrv

This driver is based on the slim midicat tool (see tools/midicat for more information).

## Installation

This is driver uses the `midicat` binary that you can get [here](https://github.com/gomidi/midicat/releases/download/v0.3.6/midicat-binaries.zip)
for Windows and Linux (it should be possible to compile it on your own, e.g. for the Mac).

The `midicat` binary is based on the rtmidi project and connects MIDI ports to Stdin and Stdout.
The idea is, to have just one binary that requires CGO (`midicat`) and for all the Go projects that need
to connect to MIDI ports just pipe the MIDI data from and to this binary.

This driver connects to the `midicat` binary via Stdin and Stdout while providing the same unified https://gitlab.com/gomidi/midi.Driver interface as `rtmididrv` and `portmididrv`. But projects importing this `midicatdrv` will not required CGO
(like that would be the case otherwise).

Download or compile the `midicat` binary and place it in your `PATH` before using this driver.
**midicat version >= 0.5.0 is required**.
