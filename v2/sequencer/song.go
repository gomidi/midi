package sequencer

import (
	"fmt"
	"sort"

	"gitlab.com/gomidi/midi/v2/smf"
)

func New() *Song {
	return &Song{
		Ticks: smf.MetricTicks(960),
	}
}

type Song struct {
	Title      string
	Composer   string
	TrackNames []string
	Ticks      smf.MetricTicks
	lastTick   int64
	bars       Bars
}

func (s *Song) Bars() Bars {
	return s.bars
}

// TODO allow inserting a bar at a certain position (renumber all bars, move the keychanges and tempochanges accordingly
// TODO allow removing a bar (renumber all bars, move the keychanges and tempochanges accordingly
// TODO move a bar (renumber all bars, move the keychanges and tempochanges accordingly
// TODO allow copying a bar with all events from selected tracks (renumber all bars, move the keychanges and tempochanges accordingly
// TODO allow playing some or all tracks
// TODO allow to record to a track

func (s *Song) AddBar(b Bar) {
	// default time signature
	var timeSig = [2]uint8{4, 4}
	if len(s.bars) > 0 {
		// overwrite with last time signature
		timeSig = s.bars[len(s.bars)-1].TimeSig
	}

	// set time signature, if not set
	if b.TimeSig == [2]uint8{0, 0} {
		b.TimeSig = timeSig
	}

	s.bars = append(s.bars, &b)
	s.bars.Renumber()
}

func (s *Song) FindBar(absTicks int64) (found *Bar) {

	for _, b := range s.bars {
		if b.AbsTicks <= absTicks {
			found = b
		} else {
			return
		}
	}

	return
}

func (s *Song) SetBarAbsTicks() {
	var absticks int64
	for _, b := range s.bars {
		b.AbsTicks = absticks
		absticks += int64(b.Len()) * int64(s.Ticks.Ticks32th())
	}
	s.lastTick = absticks
}

func (s *Song) ToSMF0() smf.SMF {
	if s.Ticks == 0 {
		s.Ticks = smf.MetricTicks(960)
	}
	var sm smf.SMF
	sm.TimeFormat = s.Ticks

	var t smf.Track
	t.Add(0, smf.MetaText(s.Title))
	t.Add(0, smf.MetaCopyright(s.Composer))
	evts, _ := s.mkBarLine(s.Ticks)

	for _, b := range s.bars {
		evts = append(evts, b.trackEvents(s.Ticks)...)
	}

	sort.Sort(evts)

	var lasttick int64

	for i := 0; i < len(evts); i++ {
		t.Add(uint32(evts[i].AbsTicks-lasttick), evts[i].Message)
		lasttick = evts[i].AbsTicks
	}

	t.Close(uint32(s.lastTick - lasttick))
	sm.Add(t)

	return sm
}

func (s Song) ToSMF1() smf.SMF {
	if s.Ticks == 0 {
		s.Ticks = smf.MetricTicks(960)
	}
	var sm smf.SMF
	sm.TimeFormat = s.Ticks
	var barTrack smf.Track
	barTrack.Add(0, smf.MetaText(s.Title))
	barTrack.Add(0, smf.MetaCopyright(s.Composer))
	barTrack.Add(0, smf.MetaTrackSequenceName("bars"))
	barevts, _ := s.mkBarLine(s.Ticks)
	var lastMessageAbs int64
	for _, ev := range barevts {
		barTrack.Add(ev.Delta, ev.Message)
		lastMessageAbs = ev.AbsTicks
	}

	barTrack.Close(uint32(s.lastTick - lastMessageAbs))
	sm.Add(barTrack)

	var allevts smf.TrackEvents

	for _, b := range s.bars {
		allevts = append(allevts, b.trackEvents(s.Ticks)...)
	}

	var settracks = map[int]bool{}

	for _, ev := range allevts {
		settracks[ev.TrackNo] = true
	}

	var tracks []int

	for no, has := range settracks {
		if has {
			tracks = append(tracks, no)
		}
	}

	sort.Ints(tracks)

	sort.Sort(allevts)

	for _, trackno := range tracks {
		name := fmt.Sprintf("track-%v", trackno)

		if len(s.TrackNames) > trackno {
			name = s.TrackNames[trackno]
		}

		var t smf.Track
		t.Add(0, smf.MetaTrackSequenceName(name))

		var lasttick int64

		for _, ev := range allevts {
			if ev.TrackNo == trackno {
				delta := uint32(ev.AbsTicks - lasttick)
				t.Add(delta, ev.Message)
				lasttick = ev.AbsTicks
			}
		}
		sm.Add(t)

	}

	return sm
}

func (s *Song) mkBarLine(ticks smf.MetricTicks) (evts smf.TrackEvents, abslength int64) {
	//sm := si.SMF
	sort.Sort(s.bars)
	s.SetBarAbsTicks()

	var timesig = [2]uint8{4, 4}
	//var tempo float64 = 120

	for _, b := range s.bars {
		if b.TimeSig != [2]uint8{0, 0} && b.TimeSig != timesig {
			timesig = b.TimeSig
			evts = append(evts, &smf.TrackEvent{
				AbsTicks: b.AbsTicks,
				Event: smf.Event{
					Message: smf.MetaMeter(b.TimeSig[0], b.TimeSig[1]),
				},
			})
		}
	}

	//abslength = absTicks

	var lasttick int64

	for i := 0; i < len(evts); i++ {
		evts[i].Delta = uint32(evts[i].AbsTicks - lasttick)
		lasttick = evts[i].AbsTicks
	}
	return
}
