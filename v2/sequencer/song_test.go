package sequencer

import (
	"testing"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func TestAddEvents(t *testing.T) {
	var s Song
	var sm smf.SMF
	ticks := smf.MetricTicks(960)

	sm.TimeFormat = ticks
	var tr0, tr1 smf.Track

	qn := ticks.Ticks4th()

	var m = smf.MetaMeter
	tr0.Add(0, smf.MetaTrackSequenceName("testmkbars"))

	// bar 0
	tr0.Add(0, m(3, 4))
	// bar 1
	tr0.Add(qn*3, midi.NoteOn(1, 60, 110))
	// bar 2
	tr0.Add(qn*2, m(4, 4))
	tr0.Add(qn*2, midi.NoteOff(1, 60))
	// bar 2 rest & bar 3
	tr0.Add(qn*6, m(3, 4))
	// bar 4 & 5
	tr0.Close(qn * 6)
	sm.Add(tr0)
	// 3/4 3/4 4/4 4/4 3/4 3/4

	// bar 0
	tr1.Add(qn, midi.NoteOn(2, 60, 120))
	tr1.Add(qn, midi.ControlChange(1, 22, 105))
	// bar 1
	// bar 2
	tr1.Add(qn*4, midi.NoteOff(2, 60))
	tr1.Close(0)
	sm.Add(tr1)

	si := smfimport{&s, sm}

	si.mkBars()
	si.addEvents()

	if len(s.TrackNames) != 2 {
		t.Errorf("len(s.Tracks) = %v // expected %v", len(s.TrackNames), 2)
	}

	bars := s.Bars()

	if len(bars) != 6 {
		t.Errorf("len(s.Bars()) = %v // expected %v", len(bars), 6)
	}

	if len(bars[0].Events) != 2 {
		t.Errorf("len(bars[0].Events) = %v // expected %v", len(bars[0].Events), 2)
	}

	got := bars[0].Events[0].Inspect()
	expected := `Event{TrackNo:1, Pos:8, Duration:40, Message: NoteOn channel: 2 key: 60 velocity: 120, absTicks: 960}`

	if got != expected {
		t.Errorf("bars[0].Events[0].Inspect() = %q // expected %q", got, expected)
	}

	got = bars[0].Events[1].Inspect()
	expected = `Event{TrackNo:1, Pos:16, Duration:0, Message: ControlChange channel: 1 controller: 22 value: 105, absTicks: 1920}`

	if got != expected {
		t.Errorf("bars[0].Events[1].Inspect() = %q // expected %q", got, expected)
	}

	if len(bars[1].Events) != 1 {
		t.Errorf("len(bars[1].Events) = %v // expected %v", len(bars[1].Events), 1)
	}

	got = bars[1].Events[0].Inspect()
	expected = `Event{TrackNo:0, Pos:0, Duration:32, Message: NoteOn channel: 1 key: 60 velocity: 110, absTicks: 2880}`

	if got != expected {
		t.Errorf("bars[1].Events[0].Inspect() = %q // expected %q", got, expected)
	}

	if len(bars[2].Events) != 0 {
		t.Errorf("len(bars[2].Events) = %v // expected %v", len(bars[2].Events), 0)
	}

}
