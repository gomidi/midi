package smftrack

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func TestMonoize(t *testing.T) {
	var bf bytes.Buffer

	var clock = smf.MetricTicks(960)
	sm := smf.New()
	sm.TimeFormat = clock

	var tr smf.Track
	tr.Add(0, smf.Message(midi.NoteOn(2, 45, 120)))
	tr.Add(clock.Ticks8th(), smf.Message(midi.NoteOn(2, 46, 100)))
	tr.Add(clock.Ticks16th(), smf.Message(midi.NoteOff(2, 45)))
	tr.Add(clock.Ticks8th(), smf.Message(midi.NoteOff(2, 46)))
	tr.Close(0)
	sm.Add(tr)
	sm.WriteTo(&bf)

	var outBf bytes.Buffer

	err := Monoize(&bf, &outBf, []int{0})

	if err != nil {
		t.Errorf("ERROR: %s", err.Error())
	}

	var outStr bytes.Buffer
	rd, err := smf.ReadFrom(&outBf)

	if err != nil && err != io.EOF {
		t.Errorf("ERROR: %s", err.Error())
	}

	var channel, val1, val2 uint8

	for _, ev := range rd.Tracks[0] {
		switch {
		case ev.Message.GetNoteStart(&channel, &val1, &val2):
			outStr.WriteString(fmt.Sprintf("@%v [%v] NoteOn(%v,%v)\n", ev.Delta, channel, val1, val2))
		case ev.Message.GetNoteEnd(&channel, &val1):
			outStr.WriteString(fmt.Sprintf("@%v [%v] NoteOff(%v)\n", ev.Delta, channel, val1))
		}
	}

	var expected = `@0 [2] NoteOn(45,120)
@480 [2] NoteOff(45)
@0 [2] NoteOn(46,100)
@240 [2] NoteOff(46)
`

	if outStr.String() != expected {
		t.Errorf("got:\n%s\nexpected:\n%s\n", outStr.String(), expected)
	}

}
