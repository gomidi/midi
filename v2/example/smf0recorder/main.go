package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"gitlab.com/gomidi/midi/v2/smf"
)

func main() {
	err := run()

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

func run() error {

	file := smf.New()
	ticks := file.TimeFormat.(smf.MetricTicks)
	bpm := float64(120)

	tr := smf.NewTrack()
	tr.Add(0, midi.MetaTempo(bpm))

	// single track recording, for multitrack we would have to collect the messages first (separated by port / midi channel)
	// and the write them after the recording on the different tracks
	in, err := midi.NewListener("port-description").
		Only(midi.Channel1Msg & midi.NoteMsg).
		Do(
			func(msg midi.Message, deltamicroSec int64) {
				delta := ticks.Ticks(bpm, time.Duration(deltamicroSec)*time.Microsecond)
				tr.Add(delta, msg.Data)
			},
		)

	in.Close()

	file.AddAndClose(0, tr)

	err2 := file.WriteFile("recorded.mid")
	if err != nil {
		return err
	}

	return err2
}
