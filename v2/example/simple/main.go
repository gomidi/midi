package main

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/gm"
	"gitlab.com/gomidi/midi/v2/smf"

	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	for _, o := range midi.OutPorts() {
		fmt.Printf("out: %s\n", o)
	}

	out := midi.FindOutPort("qsynth")
	if out < 0 {
		fmt.Printf("can't find qsynth")
		return
	}

	// create a SMF
	rd := bytes.NewReader(mkSMF())

	// read and play it
	smf.ReadTracksFrom(rd).Do(func(ev smf.TrackEvent) {
		fmt.Printf("track %v @%vms %s\n", ev.TrackNo, ev.AbsMicroSeconds/1000, ev.Message)
	}).Play(out)
}

// makes a SMF and returns the bytes
func mkSMF() []byte {
	var (
		bf    bytes.Buffer
		clock = smf.MetricTicks(96) // resolution: 96 ticks per quarternote 960 is also common
		tr    smf.Track
	)

	// first track must have tempo and meter informations
	tr.Add(0, smf.MetaMeter(3, 4))
	tr.Add(0, smf.MetaTempo(140))
	tr.Add(0, smf.MetaInstrument("Brass"))
	tr.Add(0, midi.ProgramChange(0, gm.Instr_BrassSection.Value()))
	tr.Add(0, midi.NoteOn(0, midi.Ab(3), 120))
	tr.Add(clock.Ticks8th(), midi.NoteOn(0, midi.C(4), 120))
	// duration: a quarter note (96 ticks in our case)
	tr.Add(clock.Ticks4th()*2, midi.NoteOff(0, midi.Ab(3)))
	tr.Add(0, midi.NoteOff(0, midi.C(4)))
	tr.Close(0)

	// create the SMF and add the tracks
	s := smf.New()
	s.TimeFormat = clock
	s.Add(tr)
	s.WriteTo(&bf)
	return bf.Bytes()
}
