package smftrack

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func mkSMF0() []byte {
	var bf bytes.Buffer

	sm := smf.New()
	var tr smf.Track

	tr.Add(0, smf.Message(midi.NoteOn(2, 65, 90)))
	tr.Add(0, smf.Message(midi.NoteOn(1, 24, 100)))
	tr.Add(2, smf.Message(midi.NoteOff(1, 24)))
	tr.Add(8, smf.Message(midi.NoteOn(3, 80, 109)))
	tr.Add(10, smf.Message(midi.NoteOff(2, 65)))
	tr.Add(10, smf.Message(midi.NoteOff(3, 80)))
	tr.Close(0)

	sm.Add(tr)
	sm.WriteTo(&bf)
	return bf.Bytes()
}

func ExampleToSMF1() {

	var bf bytes.Buffer

	// get some SMF0
	srcBytes := mkSMF0()

	src, err := smf.ReadFrom(bytes.NewReader(srcBytes))

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("source: ", inspectSMF(src))

	(SMF0{}).ToSMF1(src, &bf)

	// read it back in
	rd, err := smf.ReadFrom(bytes.NewReader(bf.Bytes()))

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Println("result: ", inspectSMF(rd))

	for trno, tr := range rd.Tracks {
		var channel, val1, val2 uint8

		for _, ev := range tr {
			switch {
			case ev.Message.GetNoteStart(&channel, &val1, &val2):
				fmt.Printf("Track %v [%v] NoteOn at channel %v: key %v velocity %v\n", trno, ev.Delta, channel, val1, val2)
			case ev.Message.GetNoteEnd(&channel, &val1):
				fmt.Printf("Track %v [%v] NoteOff at channel %v: key %v\n", trno, ev.Delta, channel, val1)
			}
		}
	}

	// Output: source:  <Format: SMF0 (singletrack), NumTracks: 1, TimeFormat: 960 MetricTicks>
	// result:  <Format: SMF1 (multitrack), NumTracks: 4, TimeFormat: 960 MetricTicks>
	// Track 1 [0] NoteOn at channel 1: key 24 velocity 100
	// Track 1 [2] NoteOff at channel 1: key 24
	// Track 2 [0] NoteOn at channel 2: key 65 velocity 90
	// Track 2 [20] NoteOff at channel 2: key 65
	// Track 3 [10] NoteOn at channel 3: key 80 velocity 109
	// Track 3 [20] NoteOff at channel 3: key 80
}
