package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"gitlab.com/gomidi/midi/v2/smf"
)

type playEvent struct {
	absTime uint64
	sleep   time.Duration
	data    [3]byte
	out     drivers.Out
	trackNo int
}

type player []playEvent

func run() error {
	out, err := drivers.OutByName("FLUID Synth")
	if err != nil {
		return err
	}

	defer out.Close()

	// single track playing
	// for multitrack we would have to collect the tracks events first and properly synchronize playback
	//_, err = smf.ReadTracks("Prelude4.mid", 2).
	_, err = smf.ReadTracks("Prelude4.mid", 1, 2, 3, 4).
		//_, err = smf.ReadTracks("VOYAGER.MID", 1).
		//Only(midi.NoteOnMsg, midi.NoteOffMsg).
		//Only(midi.NoteOnMsg, midi.NoteOffMsg, midi.MetaMsgType).
		Only(midi.NoteMsg).
		Do(
			func(trackNo int, msg smf.Message, delta int64, deltamicroSec int64) {
				fmt.Printf("%s (%s)\n", msg.String(), msg.Kind())
				if mm, ok := msg.(midi.Message); ok {
					time.Sleep(time.Microsecond * time.Duration(deltamicroSec))
					_ = mm
					out.Send(mm.Data)
				}
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
