package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
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
	tr.Add(0, smf.MetaTempo(bpm))

	defer midi.CloseDriver()
	in := midi.FindInPort("VMPK")

	if in < 0 {
		return fmt.Errorf("can't find MIDI in port %q", "VMPK")
	}

	var absmillisec int32

	// single track recording, for multitrack we would have to collect the messages first (separated by port / midi channel)
	// and the write them after the recording to the different tracks
	stop, err := midi.ListenTo(in, midi.ReceiverFunc(func(msg midi.Message, absms int32) {
		deltams := absms - absmillisec
		absmillisec = absms
		fmt.Printf("[%v] %s\n", deltams, msg.String())
		delta := ticks.Ticks(bpm, time.Duration(deltams)*time.Millisecond)
		tr.Add(delta, msg)
	}))

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	stop()

	file.AddAndClose(0, tr)

	err2 := file.WriteFile("recordedx.mid")
	if err != nil {
		return err
	}

	return err2
}
