package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"gitlab.com/gomidi/midi/v2/smf"
)

func run() error {
	out, err := midi.OutByName("FLUID Synth")
	if err != nil {
		return err
	}

	defer out.Close()

	// single track playing
	// for multitrack we would have to collect the tracks events first and properly synchronize playback
	_, err = smf.ReadTracks("Prelude4.mid", 1).
		Only(midi.Channel0Msg).
		Do(
			func(trackNo int, msg midi.Message, delta int64, deltamicroSec int64) {
				fmt.Printf("%s\n", msg.String())
				time.Sleep(time.Microsecond * time.Duration(deltamicroSec))
				out.Send(msg.Data)
			},
		)
	return err
}

func main() {
	err := run()

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
