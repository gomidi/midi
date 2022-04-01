# midispy
Spy on MIDI data transmitted between a sending and a receiving device

[![Documentation](http://godoc.org/gitlab.com/gomidi/midispy?status.png)](http://godoc.org/gitlab.com/gomidi/midispy)

## Usage (CLI)

    go get -d gitlab.com/gomidi/midispy/cmd/midispy
    go install gitlab.com/gomidi/midispy/cmd/midispy

To get a list of available MIDI devices, run

    midispy list

Example of output:

    ---MIDI input ports---
    0 "Midi Through:Midi Through Port-0 14:0"
    1 "Virtual Raw MIDI 2-0:VirMIDI 2-0 24:0"
    2 "MPKmini2:MPKmini2 MIDI 1 28:0"

    ---MIDI output ports---
    0 "Midi Through:Midi Through Port-0 14:0"
    1 "Virtual Raw MIDI 2-0:VirMIDI 2-0 24:0"
    2 "MPKmini2:MPKmini2 MIDI 1 28:0"
    3 "FLUID Synth (qsynth):Synth input port (qsynth:0) 128:0"

Then you can use the given ids to tell midispy to listen:

    midispy in=10 out=11
    [10] "VMPK Output:VMPK Output 128:0"
    ->
    [11] "FLUID Synth (qsynth):Synth input port (qsynth:0) 130:0"
    -----------------------
    [10->11] 13:38:00.152725 channel.NoteOn channel 0 key 58 velocity 100
    [10->11] 13:38:00.265286 channel.NoteOff channel 0 key 58
    [10->11] 13:38:01.180286 channel.NoteOn channel 0 key 58 velocity 100
    [10->11] 13:38:01.276850 channel.NoteOff channel 0 key 58
    [10->11] 13:38:01.701086 channel.NoteOn channel 0 key 71 velocity 100
    [10->11] 13:38:01.786206 channel.NoteOff channel 0 key 71
    [10->11] 13:38:01.990798 channel.NoteOn channel 0 key 69 velocity 100
    [10->11] 13:38:02.081213 channel.NoteOff channel 0 key 69
    [10->11] 13:38:02.255353 channel.NoteOn channel 0 key 58 velocity 100
    [10->11] 13:38:02.345848 channel.NoteOff channel 0 key 58
    [10->11] 13:38:32.193058 channel.NoteOn channel 0 key 58 velocity 100
    [10->11] 13:38:33.215804 channel.NoteOn channel 0 key 74 velocity 100
    [10->11] 13:38:34.087821 channel.NoteOff channel 0 key 74
    [10->11] 13:38:34.580737 channel.NoteOff channel 0 key 58


To get help:

    midispy help   

## Usage (library)

    go get gitlab.com/gomidi/midispy@latest

[![Documentation](http://godoc.org/gitlab.com/gomidi/midispy?status.png)](http://godoc.org/gitlab.com/gomidi/midispy)

```go
package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/mid"
	"gitlab.com/gomidi/midispy"
	driver "gitlab.com/gomidi/rtmididrv"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	drv, err := driver.New()
	must(err)

	defer drv.Close()

	in, err := mid.OpenIn(drv, 6 /* change to fit your needs */, "")
	must(err)

	out, err := mid.OpenOut(drv, 4 /* change to fit your needs */, "")
	must(err)

	rd := mid.NewReader(mid.NoLogger())

	// see gitlab.com/gomidi/midi/mid#Reader
	// to learn how to listen to other messages
	rd.Msg.Channel.NoteOn = func(_ *mid.Position, channel, key, velocity uint8) {
		fmt.Printf("note on (channel: %d key: %d velocity: %d)\n", channel, key, velocity)
	}

	err = midispy.Run(in, out, rd)
	must(err)

	fmt.Println("start listening...")

	time.Sleep(time.Minute)

	fmt.Println("stop listening")

}
```
