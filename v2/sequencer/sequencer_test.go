package sequencer

import (
	"fmt"
	"strings"
	"testing"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

type testevent struct {
	delta uint32
	msg   smf.Message
}

type testevents []testevent

func TestToSMF0(t *testing.T) {
	var s Song
	s.Composer = "me"
	s.Title = "testpiece"

	var b Bar
	var k smf.Key
	smf.AbMaj().GetMetaKey(&k)
	b.Key = &k
	b.TimeSig = [2]uint8{3, 4}

	b.Events = append(b.Events, &Event{
		Message:  smf.Message(midi.NoteOn(0, 50, 100)),
		Duration: 16,
	})
	b.Events = append(b.Events, &Event{
		Message:  smf.Message(midi.NoteOn(1, 60, 100)),
		Pos:      8,
		Duration: 4,
	})

	//	fmt.Printf("len bars is: %v\n", len(s.bars))

	//fmt.Printf("adding first bar\n")
	s.AddBar(b)

	b = Bar{}
	b.Events = append(b.Events, &Event{
		Message:  smf.Message(midi.NoteOn(0, 55, 90)),
		Duration: 16,
	})
	b.Events = append(b.Events, &Event{
		Message: smf.Message(midi.ControlChange(1, 22, 127)),
		Pos:     8,
	})

	//	fmt.Printf("adding second bar\n")
	s.AddBar(b)

	sm := s.ToSMF0()

	_ = sm

	got := strings.TrimSpace(sm.String())
	expected := strings.TrimSpace(`
#### SMF Format: 0 TimeFormat: 960 MetricTicks NumTracks: 1 ####
## TRACK 0 ##
#0 [0] MetaText text: "testpiece"
#0 [0] MetaCopyright text: "me"
#0 [0] MetaTimeSig meter: 3/4
#0 [0] NoteOn channel: 0 key: 50 velocity: 100
#0 [960] NoteOn channel: 1 key: 60 velocity: 100
#0 [480] NoteOff channel: 1 key: 60
#0 [480] NoteOff channel: 0 key: 50
#0 [960] NoteOn channel: 0 key: 55 velocity: 90
#0 [960] ControlChange channel: 1 controller: 22 value: 127
#0 [960] NoteOff channel: 0 key: 55
#0 [960] MetaEndOfTrack
`)

	if got != expected {
		t.Errorf("got:\n%s\nexpected:\n%s\n", got, expected)
	}
}

func TestToSMF1(t *testing.T) {
	var s Song
	s.Composer = "me"
	s.Title = "testpiece"
	/*
		s.Tracks = []string{
			"first",
			"second",
		}
	*/

	s.Tracks = []string{
		"first track",
	}

	var b Bar
	var k smf.Key
	smf.AbMaj().GetMetaKey(&k)
	b.Key = &k
	b.TimeSig = [2]uint8{3, 4}

	b.Events = append(b.Events, &Event{
		Message:  smf.Message(midi.NoteOn(0, 50, 100)),
		Duration: 16,
		TrackNo:  0,
	})
	b.Events = append(b.Events, &Event{
		Message:  smf.Message(midi.NoteOn(0, 60, 100)),
		Pos:      8,
		Duration: 4,
		TrackNo:  1,
	})

	//	fmt.Printf("len bars is: %v\n", len(s.bars))

	//fmt.Printf("adding first bar\n")
	s.AddBar(b)

	b = Bar{}
	b.Events = append(b.Events, &Event{
		Message:  smf.Message(midi.NoteOn(0, 55, 90)),
		Duration: 16,
		TrackNo:  0,
	})
	b.Events = append(b.Events, &Event{
		Message: smf.Message(midi.ControlChange(0, 22, 127)),
		Pos:     8,
		TrackNo: 1,
	})

	//	fmt.Printf("adding second bar\n")
	s.AddBar(b)

	sm := s.ToSMF1()

	_ = sm

	got := strings.TrimSpace(sm.String())
	expected := strings.TrimSpace(`
#### SMF Format: 0 TimeFormat: 960 MetricTicks NumTracks: 3 ####
## TRACK 0 ##
#0 [0] MetaText text: "testpiece"
#0 [0] MetaCopyright text: "me"
#0 [0] MetaTrackName text: "bars"
#0 [0] MetaTimeSig meter: 3/4
#0 [5760] MetaEndOfTrack
## TRACK 1 ##
#1 [0] MetaTrackName text: "first track"
#1 [0] NoteOn channel: 0 key: 50 velocity: 100
#1 [1920] NoteOff channel: 0 key: 50
#1 [960] NoteOn channel: 0 key: 55 velocity: 90
#1 [1920] NoteOff channel: 0 key: 55
## TRACK 2 ##
#2 [0] MetaTrackName text: "track-1"
#2 [960] NoteOn channel: 0 key: 60 velocity: 100
#2 [480] NoteOff channel: 0 key: 60
#2 [2400] ControlChange channel: 0 controller: 22 value: 127
`)

	if got != expected {
		t.Errorf("got:\n%s\nexpected:\n%s\n", got, expected)
	}
}

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

	mkBars(&s, sm)
	addEvents(&s, sm)

	if len(s.Tracks) != 2 {
		t.Errorf("len(s.Tracks) = %v // expected %v", len(s.Tracks), 2)
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

func TestMkBars(t *testing.T) {
	// mkBars(s *Song, sm smf.SMF) (keychanges []*smf.TrackEvent) {
	ticks := smf.MetricTicks(960)
	qn := ticks.Ticks4th()

	var m = smf.MetaMeter

	tests := []struct {
		descr      string
		evts       testevents
		closeTicks uint32
		numBars    int
		bars       string
	}{
		{ // 0
			"mixture of meters. each bar has a meter change",
			testevents{
				{0, m(4, 4)},
				{qn * 4, m(3, 4)},
				{qn * 3, m(4, 4)},
			},
			qn * 4,
			3,
			"4/4 3/4 4/4",
		},
		{ // 1
			"mixture of meters. second bar has no meter change",
			testevents{
				{0, m(4, 4)},
				{qn * 8, m(3, 4)},
				{qn * 3, m(4, 4)},
			},
			qn * 4,
			4,
			"4/4 4/4 3/4 4/4",
		},
		{ // 2
			"mixture of meters. third bar has no meter change",
			testevents{
				{0, m(4, 4)},
				{qn * 4, m(3, 4)},
				{qn * 6, m(4, 4)},
			},
			qn * 4,
			4,
			"4/4 3/4 3/4 4/4",
		},
		{ // 3
			"mixture of meters. last bar has no meter change",
			testevents{
				{0, m(4, 4)},
				{qn * 4, m(3, 4)},
				{qn * 3, m(4, 4)},
			},
			qn * 8,
			4,
			"4/4 3/4 4/4 4/4",
		},
		{ // 4
			"single meter for 8 bars",
			testevents{
				{0, m(3, 4)},
			},
			qn * 3 * 8,
			8,
			"3/4 3/4 3/4 3/4 3/4 3/4 3/4 3/4",
		},
		{ // 5
			"default meter for 8 bars",
			testevents{},
			qn * 4 * 8,
			8,
			"4/4 4/4 4/4 4/4 4/4 4/4 4/4 4/4",
		},
		{ // 6
			"default meter for 8 bars followed by 2 bars with other meter",
			testevents{
				{qn * 4 * 8, m(3, 4)},
			},
			qn * 3 * 2,
			10,
			"4/4 4/4 4/4 4/4 4/4 4/4 4/4 4/4 3/4 3/4",
		},
	}

	for i, test := range tests {
		var s Song
		var sm smf.SMF

		sm.TimeFormat = ticks
		var tr smf.Track

		tr.Add(0, smf.MetaTrackSequenceName(fmt.Sprintf("testmkbars%v", i)))

		for _, ev := range test.evts {
			tr.Add(ev.delta, ev.msg)
		}

		tr.Close(test.closeTicks)
		sm.Add(tr)

		mkBars(&s, sm)

		bars := s.Bars()
		got := len(bars)
		expected := test.numBars

		if got != expected {
			t.Errorf("[%v] %s\n\tlen(bars) = %v // expected: %v", i, test.descr, got, expected)
		}

		var bf strings.Builder

		for _, b := range bars {
			bf.WriteString(fmt.Sprintf("%v/%v ", b.TimeSig[0], b.TimeSig[1]))
		}

		gotBars := strings.TrimSpace(bf.String())

		if gotBars != test.bars {
			t.Errorf("[%v] %s:\n\t%s // expected: %s", i, test.descr, gotBars, test.bars)
		}

	}
}
