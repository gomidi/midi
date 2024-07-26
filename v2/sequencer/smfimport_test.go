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

	s.TrackNames = []string{
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
#### SMF Format: 1 TimeFormat: 960 MetricTicks NumTracks: 3 ####
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
#1 [960] MetaEndOfTrack
## TRACK 2 ##
#2 [0] MetaTrackName text: "track-1"
#2 [960] NoteOn channel: 0 key: 60 velocity: 100
#2 [480] NoteOff channel: 0 key: 60
#2 [2400] ControlChange channel: 0 controller: 22 value: 127
#2 [1920] MetaEndOfTrack
`)

	if got != expected {
		t.Errorf("got:\n%s\nexpected:\n%s\n", got, expected)
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
		var s = New()
		var sm smf.SMF

		sm.TimeFormat = ticks
		var tr smf.Track

		tr.Add(0, smf.MetaTrackSequenceName(fmt.Sprintf("testmkbars%v", i)))

		for _, ev := range test.evts {
			tr.Add(ev.delta, ev.msg)
		}

		tr.Close(test.closeTicks)
		sm.Add(tr)

		si := smfimport{s, sm}
		si.mkBars()

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

func TestFromSMF1(t *testing.T) {

	ticks := smf.MetricTicks(960)

	var sm smf.SMF

	sm.TimeFormat = ticks

	var t0, t1, t2 smf.Track

	t0.Add(0, smf.AbMaj())
	t0.Add(0, smf.MetaMeter(3, 4))
	t0.Add(6*ticks.Ticks4th(), smf.MetaMeter(4, 4))
	t0.Add(0, smf.MetaTempo(140.00))
	t0.Close(4 * ticks.Ticks4th())
	sm.Add(t0)

	t1.Add(ticks.Ticks4th(), midi.NoteOn(0, 50, 100))
	t1.Add(ticks.Ticks4th(), midi.NoteOff(0, 50))
	t1.Close(0)
	sm.Add(t1)

	t2.Add(2*ticks.Ticks4th(), midi.NoteOn(1, 30, 100))
	t2.Add(2*ticks.Ticks4th(), midi.NoteOff(1, 30))
	t2.Close(0)
	sm.Add(t2)

	song := FromSMF(sm)

	if len(song.Bars()) != 3 {
		t.Errorf("wrong number of bars: %v // expected: %v", len(song.Bars()), 3)
	}

	if len(song.TrackNames) != 3 {
		t.Errorf("wrong number of tracks: %v // expected: %v", len(song.TrackNames), 3)
	}

	// TODO: further tests

}
