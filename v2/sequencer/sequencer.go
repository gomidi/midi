package sequencer

import (
	"fmt"
	"math"
	"sort"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

type Bars []*Bar

func (b Bars) Len() int {
	return len(b)
}

func (br Bars) Swap(a, b int) {
	br[a], br[b] = br[b], br[a]
}

func (br Bars) Less(a, b int) bool {
	return br[a].Number < br[b].Number
}

func (br Bars) Renumber() {
	for i := 0; i < len(br); i++ {
		br[i].Number = uint(i)
	}
}

type Events []*Event

func (e Events) Len() int {
	return len(e)
}

func (e Events) Swap(a, b int) {
	e[a], e[b] = e[b], e[a]
}

// warning: only works within a bar
func (e Events) Less(a, b int) bool {
	return e[a].Pos < e[b].Pos
}

type Bar struct {
	Number   uint
	TimeSig  [2]uint8
	Events   Events
	Key      *smf.Key // TODO a key change, if != nil
	absTicks int64
}

func (b *Bar) SortEvents() {
	evts := b.Events
	sort.Sort(evts)
	b.Events = evts
}

func (b Bar) Len() uint8 {
	return b.TimeSig[0] * 32 / b.TimeSig[1]
}

type Event struct {
	TrackNo  int
	Pos      uint8       // in 32th
	Duration uint8       // in 32th for noteOn messages, it is the length of the note, for all other messages, it is 0
	Message  smf.Message // may only be channel messages or sysex messages. no noteon velocity 0, or noteoff messages, this is expressed via Duration
	absTicks int64       // just for smf import
}

func (e Event) Inspect() string {
	return fmt.Sprintf("Event{TrackNo:%v, Pos:%v, Duration:%v, Message: %s, absTicks: %v}", e.TrackNo, int(e.Pos), int(e.Duration), e.Message.String(), e.absTicks)
}

func (e *Event) AbsTicks(b *Bar, ticks smf.MetricTicks) (start, end int64) {
	start = b.absTicks + int64(ticks.Ticks32th()*uint32(e.Pos))
	if e.Duration <= 0 {
		return start, 0
	}
	end = start + int64(ticks.Ticks32th()*uint32(e.Duration))
	return
}

func FromSMF(sm smf.SMF) *Song {
	switch sm.Format() {
	case 0:
		return fromSMF1(sm.ConvertToSMF1())
	case 1:
		return fromSMF1(sm)
	default:
		panic(fmt.Sprintf("SMF format %v is not supported", sm.Format()))
	}
}

func mkBars(s *Song, sm smf.SMF) {
	var timesigs smf.TrackEvents
	var keychanges smf.TrackEvents
	ticks := sm.TimeFormat.(smf.MetricTicks)

	var totalTicks int64
	var nkey smf.Key

	for _, tr := range sm.Tracks {
		var absTicks int64
		//var te []*smf.TrackEvent
		var instr string
		var track string

		for _, ev := range tr {
			absTicks += int64(ev.Delta)

			var num, denom uint8

			if absTicks > totalTicks {
				totalTicks = absTicks
			}

			switch {
			case ev.Message.GetMetaMeter(&num, &denom):
				timesigs = append(timesigs, &smf.TrackEvent{
					AbsTicks: absTicks,
					Event: smf.Event{
						Message: ev.Message,
					},
				})
			case ev.Message.GetMetaInstrument(&instr):
			case ev.Message.GetMetaTrackName(&track):
			case ev.Message.GetMetaKey(&nkey):
				keychanges = append(keychanges, &smf.TrackEvent{
					AbsTicks: absTicks,
					Event: smf.Event{
						Message: ev.Message,
					},
				})
			}
		}

		//	fmt.Printf("len(ts): %v\n", len(timesigs))

		var name string

		switch {
		case track != "":
			name = track
		case instr != "":
			name = instr
		}

		s.Tracks = append(s.Tracks, name)
	}

	sort.Sort(timesigs)
	if len(timesigs) == 0 || timesigs[0].AbsTicks != 0 {
		timesigs = append(timesigs, &smf.TrackEvent{
			AbsTicks: 0,
			Event: smf.Event{
				Message: smf.MetaMeter(4, 4),
			},
		})
		sort.Sort(timesigs)
	}

	var lastTick int64
	var lastNum, lastdenom uint8
	ticks32th := ticks.Ticks32th()
	var currAbsTick int64

	for _, ts := range timesigs {
		tickspassed := ts.AbsTicks - lastTick

		//fmt.Printf("tickspassed: %v\n", tickspassed)

		rounded := math.Round(float64(lastNum) * 32 / float64(lastdenom))

		n32thpassed := float64(tickspassed) / float64(ticks32th)
		//fmt.Printf("n32thpassed: %v\n", n32thpassed)

		var barspassed int

		if n32thpassed != 0 {
			barspassed = int(math.Round(n32thpassed / rounded))
		}
		//fmt.Printf("barspassed: %v\n", barspassed)

		for n := 1; n < barspassed; n++ {
			var b Bar
			b.TimeSig[0] = lastNum
			b.TimeSig[1] = lastdenom
			b.absTicks = currAbsTick
			currAbsTick += int64(b.Len()) * int64(ticks32th)
			s.AddBar(b)
			//	fmt.Printf("bar added\n")
		}

		if !ts.Message.GetMetaMeter(&lastNum, &lastdenom) {
			// something strange happend
			panic("whow!")
		}

		var b Bar
		b.TimeSig[0] = lastNum
		b.TimeSig[1] = lastdenom
		b.absTicks = currAbsTick
		currAbsTick += int64(b.Len()) * int64(ticks32th)
		s.AddBar(b)
		//fmt.Printf("bar added\n")

		lastTick = tickspassed
	}

	//fmt.Printf("totalTicks: %v currAbsTick: %v\n", totalTicks, currAbsTick)

	if totalTicks > currAbsTick {
		tickspassed := totalTicks - currAbsTick

		rounded := math.Round(float64(lastNum) * 32 / float64(lastdenom))

		n32thpassed := float64(tickspassed) / float64(ticks32th)
		barspassed := int(math.Round(n32thpassed / rounded))

		for n := 0; n < barspassed; n++ {
			var b Bar
			b.TimeSig[0] = lastNum
			b.TimeSig[1] = lastdenom
			s.AddBar(b)
		}
	}

	sort.Sort(keychanges)

	for _, keychange := range keychanges {
		b := s.FindBar(keychange.AbsTicks)
		if b != nil {
			var k smf.Key
			if keychange.Message.GetMetaKey(&k) {
				b.Key = &k
			}
		}
	}
	/*
		for _, s.bars {

		}
	*/
	return
}

func (b *Bar) barPos(absTicks int64, ticks smf.MetricTicks) uint8 {
	diff := absTicks - b.absTicks
	t32th := ticks.Ticks32th()
	return uint8(math.Round(float64(diff) / float64(t32th)))
}

func addEvents(s *Song, sm smf.SMF) {
	ticks := sm.TimeFormat.(smf.MetricTicks)
	t32th := int64(ticks.Ticks32th())
	//var keychanges []*smf.TrackEvent

	for trackno, tr := range sm.Tracks {
		//var t Track
		//var te []*smf.TrackEvent
		var noteOns = map[[2]uint8]*Event{} // [2]uint8{channel,key} to *Event

		var absTicks int64

		for _, ev := range tr {
			absTicks += int64(ev.Delta)
			b := s.FindBar(absTicks)
			//var text string
			//var num, denom uint8
			var nkey smf.Key
			var bpm float64
			var channel, key, velocity uint8
			var e Event
			e.absTicks = absTicks
			e.TrackNo = trackno
			e.Pos = b.barPos(absTicks, ticks)

			switch {
			case ev.Message.GetNoteStart(&channel, &key, &velocity):
				e.Message = ev.Message
				epoint := &e
				b.Events = append(b.Events, epoint)
				noteOns[[2]uint8{channel, key}] = epoint

			case ev.Message.GetNoteEnd(&channel, &key):
				noteonev, found := noteOns[[2]uint8{channel, key}]
				if found {
					noteonev.Duration = uint8((absTicks - noteonev.absTicks) / t32th)
				}
				delete(noteOns, [2]uint8{channel, key})

			case ev.Message.Is(midi.ChannelMsg):
				e.Message = ev.Message
				epoint := &e
				b.Events = append(b.Events, epoint)
			case ev.Message.Is(midi.SysExMsg):
				e.Message = ev.Message
				epoint := &e
				b.Events = append(b.Events, epoint)
			case ev.Message.GetMetaKey(&nkey):
				// ignore
			case ev.Message.GetMetaTempo(&bpm):
				// should not happen, since the tempochanges should all be sorted out to sm.TempoChanges()
				panic("should not happen")
			default:
				// ignore
			}
		}

		//s.AddTrack(&t)

	}

	//	var _ = tempochanges
}

func addTempoChanges(s *Song, sm smf.SMF) {
	//e.Pos = b.barPos(absTicks, ticks)
	tcs := sm.TempoChanges()
	ticks := sm.TimeFormat.(smf.MetricTicks)

	for _, tc := range tcs {
		var change Event

		b := s.FindBar(tc.AbsTicks)

		if b != nil {
			change.Pos = b.barPos(tc.AbsTicks, ticks)
			change.absTicks = tc.AbsTicks
			change.Message = smf.MetaTempo(tc.BPM)
			change.TrackNo = 0
			b.Events = append(b.Events, &change)
			b.SortEvents()
		}
	}
}

func fromSMF1(sm smf.SMF) *Song {
	s := &Song{}

	// create the bars and their abspos
	mkBars(s, sm)

	// add events to the bars
	addEvents(s, sm)

	// add tempochanges to the bars
	addTempoChanges(s, sm)

	return s
}

type Song struct {
	Title    string
	Composer string
	Tracks   []string
	bars     Bars
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
	//fmt.Printf("len bars is: %v\n", len(s.bars))
	var timeSig = [2]uint8{4, 4}
	if len(s.bars) > 0 {
		timeSig = s.bars[len(s.bars)-1].TimeSig
		//fmt.Printf("last TS was %v\n", timeSig)
	}
	if b.TimeSig == [2]uint8{0, 0} {
		//fmt.Printf("overwrite TS with: %v\n", timeSig)
		b.TimeSig = timeSig
	}

	//fmt.Printf("appending TS: %v\n", b.TimeSig)
	s.bars = append(s.bars, &b)
	s.bars.Renumber()
}

func (s *Song) FindBar(absTicks int64) (found *Bar) {

	for _, b := range s.bars {
		if b.absTicks <= absTicks {
			found = b
		} else {
			return
		}
	}

	return
}

func (s *Song) mkBarLine(ticks smf.MetricTicks) (evts smf.TrackEvents, abslength int64) {
	sort.Sort(s.bars)

	var timesig = [2]uint8{4, 4}
	//var tempo float64 = 120
	var absTicks int64

	for _, b := range s.bars {
		b.absTicks = absTicks
		if b.TimeSig != [2]uint8{0, 0} && b.TimeSig != timesig {
			timesig = b.TimeSig
			evts = append(evts, &smf.TrackEvent{
				AbsTicks: absTicks,
				Event: smf.Event{
					Message: smf.MetaMeter(b.TimeSig[0], b.TimeSig[1]),
				},
			})
		}
		absTicks += int64(ticks.Ticks32th() * uint32(b.Len()))
	}

	abslength = absTicks

	var lasttick int64

	for i := 0; i < len(evts); i++ {
		evts[i].Delta = uint32(evts[i].AbsTicks - lasttick)
		lasttick = evts[i].AbsTicks
	}
	return
}

func (b *Bar) trackEvents(ticks smf.MetricTicks) (evts smf.TrackEvents) {
	for _, ev := range b.Events {
		start, end := ev.AbsTicks(b, ticks)
		evts = append(evts, &smf.TrackEvent{
			AbsTicks: start,
			Event: smf.Event{
				Message: smf.Message(ev.Message.Bytes()),
			},
			TrackNo: ev.TrackNo,
		})

		var channel, key, velocity uint8
		if ev.Message.GetNoteStart(&channel, &key, &velocity) && end != 0 {
			evts = append(evts, &smf.TrackEvent{
				AbsTicks: end,
				Event: smf.Event{
					Message: smf.Message(midi.NoteOff(channel, key)),
				},
				TrackNo: ev.TrackNo,
			})
		}
	}

	sort.Sort(evts)

	var lasttick int64

	for i := 0; i < len(evts); i++ {
		evts[i].Delta = uint32(evts[i].AbsTicks - lasttick)
		lasttick = evts[i].AbsTicks
	}

	return evts
}

func (s *Song) ToSMF0() smf.SMF {
	var sm smf.SMF
	var ticks = smf.MetricTicks(960)
	sm.TimeFormat = ticks
	var t smf.Track
	t.Add(0, smf.MetaText(s.Title))
	t.Add(0, smf.MetaCopyright(s.Composer))
	evts, abslength := s.mkBarLine(ticks)

	for _, b := range s.bars {
		evts = append(evts, b.trackEvents(ticks)...)
	}

	sort.Sort(evts)

	var lasttick int64

	for i := 0; i < len(evts); i++ {
		t.Add(uint32(evts[i].AbsTicks-lasttick), evts[i].Message)
		lasttick = evts[i].AbsTicks
	}

	t.Close(uint32(abslength - lasttick))
	sm.Add(t)

	return sm
}

func (s Song) ToSMF1() smf.SMF {
	var sm smf.SMF
	var ticks = smf.MetricTicks(960)
	sm.TimeFormat = ticks
	var barTrack smf.Track
	barTrack.Add(0, smf.MetaText(s.Title))
	barTrack.Add(0, smf.MetaCopyright(s.Composer))
	barTrack.Add(0, smf.MetaTrackSequenceName("bars"))
	barevts, abslength := s.mkBarLine(ticks)
	var lastMessageAbs int64
	for _, ev := range barevts {
		barTrack.Add(ev.Delta, ev.Message)
		lastMessageAbs = ev.AbsTicks
	}

	barTrack.Close(uint32(abslength - lastMessageAbs))
	sm.Add(barTrack)

	var allevts smf.TrackEvents

	for _, b := range s.bars {
		allevts = append(allevts, b.trackEvents(ticks)...)
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

		if len(s.Tracks) > trackno {
			name = s.Tracks[trackno]
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
