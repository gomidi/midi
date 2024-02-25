# midicatdrv

This driver is based on the slim midicat tool (see tools/midicat for more information).

## Installation

This is driver uses the `midicat` binary.

Download the binaries (for Windows) [here](https://gitlab.com/gomidi/midi/-/releases).

Or install them via 

    go install gitlab.com/gomidi/midi/tools/midicat@latest

(When using windows, run the commands inside `cmd.exe`.)

The `midicat` binary is based on the rtmidi project and connects MIDI ports to Stdin and Stdout.
The idea is, to have just one binary that requires CGO (`midicat`) and for all the Go projects that need
to connect to MIDI ports just pipe the MIDI data from and to this binary.

This driver connects to the `midicat` binary via Stdin and Stdout while providing the same unified https://gitlab.com/gomidi/v2/drivers.Driver interface as `rtmididrv` and `portmididrv`. But projects importing this `midicatdrv` will not required CGO
(like that would be the case otherwise).

Download or compile the `midicat` binary and place it in your `PATH` before using this driver.
**midicat version >= 0.6.8 is required**.
