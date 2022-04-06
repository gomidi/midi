package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	//_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv" // autoregisters driver
)

func printPorts() {
	outs := midi.OutPorts()
	for _, o := range outs {
		fmt.Printf("out: %s\n", o)
	}
}

func run() error {

	out := midi.FindOutPort("qsynth")
	if out < 0 {
		return fmt.Errorf("can't find qsynth")
	}

	//result := smf.ReadTracks("Prelude4.mid", 2).
	//result := smf.ReadTracks("Prelude4.mid", 1, 2, 3, 4, 5, 6, 7).
	return smf.ReadTracks("Prelude4.mid").
		//result := smf.ReadTracks("VOYAGER.MID").
		//result := smf.ReadTracks("VOYAGER.MID", 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20).
		//Only(midi.NoteOnMsg, midi.NoteOffMsg).
		//Only(midi.NoteOnMsg, midi.NoteOffMsg, midi.MetaMsgType).
		//Only(midi.NoteMsg, midi.ControlChangeMsg, midi.ProgramChangeMsg).
		//Only(midi.NoteOnMsg, midi.NoteOffMsg, midi.ControlChangeMsg, midi.ProgramChangeMsg, smf.MetaTrackNameMsg).
		//Only(midi.ProgramChangeMsg, smf.MetaTrackNameMsg, smf.MetaTempoMsg, smf.MetaTimeSigMsg).
		//Only(midi.MetaMsg).
		Do(
			func(te smf.TrackEvent) {
				if te.Message.IsMeta() {
					fmt.Printf("[%v] %s\n", te.TrackNo, te.Message.String())
					/*
						var t string
						if mm.Text(&t) {
							//fmt.Printf("[%v] %s %s (%s): %q\n", te.TrackNo, msg.Type().Kind(), msg.String(), msg.Type(), t)
							fmt.Printf("[%v] %s: %q\n", te.TrackNo, te.Type, t)
							//fmt.Printf("[%v] %s %s (%s): %q\n", te.TrackNo, mm.Type().Kind(), mm.String(), mm.Type(), t)
						}
						var bpm float64
						if mm.Tempo(&bpm) {
							fmt.Printf("[%v] %s: %v\n", te.TrackNo, te.Type, math.Round(bpm))
						}
					*/
				} else {
					fmt.Printf("[%v] %s\n", te.TrackNo, te.Message)
				}
			},
		).Play(out)
}

func main() {
	defer midi.CloseDriver()
	err := run()

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
