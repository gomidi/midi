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
	// and the write them after the recording to the different tracks
	listener, err := midi.NewListener("VMPK", func(msg midi.Message, deltamicroSec int64) {
		fmt.Printf("[%v] %s\n", deltamicroSec, msg.String())
		delta := ticks.Ticks(bpm, time.Duration(deltamicroSec)*time.Microsecond)
		tr.Add(delta, msg.Data)
	})

	if err != nil {
		return err
	}
	listener.Only(midi.ChannelMsg).StartListening()

	time.Sleep(5 * time.Second)

	listener.StopListening()
	listener.Close()

	file.AddAndClose(0, tr)

	err2 := file.WriteFile("recorded.mid")
	if err != nil {
		return err
	}

	return err2
}
