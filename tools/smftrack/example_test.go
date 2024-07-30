package smftrack

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func inspectSMF(sm *smf.SMF) string {
	var bf bytes.Buffer

	switch sm.Format() {
	case 0:
		bf.WriteString("<Format: SMF0 (singletrack), ")
	case 1:
		bf.WriteString("<Format: SMF1 (multitrack), ")
	default:
		bf.WriteString("<Format: SMF2 (sequences), ")
	}

	bf.WriteString(fmt.Sprintf("NumTracks: %v, TimeFormat: %s>", len(sm.Tracks), sm.TimeFormat.String()))

	return bf.String()
}

func mkSMF1() []byte {
	var bf bytes.Buffer

	var sm = smf.New()
	var tr1, tr2 smf.Track

	tr1.Add(0, smf.Message(midi.NoteOn(2, 65, 90)))
	tr1.Add(20, smf.Message(midi.NoteOff(2, 65)))
	tr1.Close(0)
	sm.Add(tr1)

	tr2.Add(0, smf.Message(midi.NoteOn(1, 24, 100)))
	tr2.Add(2, smf.Message(midi.NoteOff(1, 24)))
	tr2.Close(0)
	sm.Add(tr2)

	sm.WriteTo(&bf)
	return bf.Bytes()
}

func Example() {

	var bf bytes.Buffer

	// get some SMF1
	srcBytes := mkSMF1()

	sm, err := smf.ReadFrom(bytes.NewReader(srcBytes))

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("source: ", inspectSMF(sm))

	// create a new track
	tr := New(2)
	tr.AddEvents(
		NewEvent(10, smf.Message(midi.NoteOn(3, 80, 109))),
		NewEvent(30, smf.Message(midi.NoteOff(3, 80))),
	)

	// add the track to the SMF1
	(SMF1{}).AddTracks(sm, &bf, tr)

	// read it back in
	sm, err = smf.ReadFrom(bytes.NewReader(bf.Bytes()))

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("with added track: ", inspectSMF(sm))

	// convert it to SMF0
	bf.Reset()
	(SMF1{}).ToSMF0(sm, &bf)

	// read back in to check
	sm, err = smf.ReadFrom(bytes.NewReader(bf.Bytes()))

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("result: ", inspectSMF(sm))

	for _, ev := range sm.Tracks[0] {
		var channel, val1, val2 uint8

		switch {
		case ev.Message.GetNoteStart(&channel, &val1, &val2):
			fmt.Printf("[%v] NoteOn at channel %v: key %v velocity %v\n", ev.Delta, channel, val1, val2)
		case ev.Message.GetNoteEnd(&channel, &val1):
			fmt.Printf("[%v] NoteOff at channel %v: key %v\n", ev.Delta, channel, val1)
		}
	}

	// Output: source:  <Format: SMF1 (multitrack), NumTracks: 2, TimeFormat: 960 MetricTicks>
	// with added track:  <Format: SMF1 (multitrack), NumTracks: 3, TimeFormat: 960 MetricTicks>
	// result:  <Format: SMF0 (singletrack), NumTracks: 1, TimeFormat: 960 MetricTicks>
	// [0] NoteOn at channel 2: key 65 velocity 90
	// [0] NoteOn at channel 1: key 24 velocity 100
	// [2] NoteOff at channel 1: key 24
	// [8] NoteOn at channel 3: key 80 velocity 109
	// [10] NoteOff at channel 2: key 65
	// [10] NoteOff at channel 3: key 80

}
