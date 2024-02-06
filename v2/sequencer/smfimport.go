package sequencer

import (
	"fmt"
	"math"
	"sort"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

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

func fromSMF1(sm smf.SMF) *Song {
	si := smfimport{
		song: &Song{
			Ticks: sm.TimeFormat.(smf.MetricTicks),
		},
		SMF: sm,
	}

	// create the bars and their abspos
	si.mkBars()

	// add events to the bars
	si.addEvents()

	// add tempochanges to the bars
	si.addTempoChanges()

	return si.song
}

type smfimport struct {
	song *Song
	SMF  smf.SMF
}

func (si *smfimport) addTempoChanges() {
	s := si.song
	sm := si.SMF

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

func (si *smfimport) mkBars() {
	s := si.song
	sm := si.SMF
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

		s.TrackNames = append(s.TrackNames, name)
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
			b.AbsTicks = currAbsTick
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
		b.AbsTicks = currAbsTick
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

func (si *smfimport) addEvents() {
	s := si.song
	sm := si.SMF
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
				//panic("should not happen")
			default:
				// ignore
			}
		}

		//s.AddTrack(&t)

	}

	//	var _ = tempochanges
}
