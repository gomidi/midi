package sequencer

import (
	"fmt"
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

type byAbsTicks []*smf.TrackEvent

func (b byAbsTicks) Len() int {
	return len(b)
}

func (br byAbsTicks) Swap(a, b int) {
	br[a], br[b] = br[b], br[a]
}

func (br byAbsTicks) Less(a, b int) bool {
	return br[a].AbsTicks < br[b].AbsTicks
}

type Track struct {
	Name    string
	Channel uint8
	Events  []Event
	te      []*smf.TrackEvent // only needed internally and temporary when importing from SMF
}

func (t Track) TrackEvents(ticks smf.MetricTicks) (evts []*smf.TrackEvent) {
	for _, ev := range t.Events {
		start, end := ev.AbsTicks(ticks)
		evts = append(evts, &smf.TrackEvent{
			AbsTicks: start,
			Event: smf.Event{
				Message: smf.Message(ev.Message.Bytes()),
			},
		})

		var channel, key, velocity uint8
		if ev.Message.GetNoteStart(&channel, &key, &velocity) && end != 0 {
			evts = append(evts, &smf.TrackEvent{
				AbsTicks: end,
				Event: smf.Event{
					Message: smf.Message(midi.NoteOff(channel, key)),
				},
			})
		}
	}

	sort.Sort(byAbsTicks(evts))

	var lasttick int64

	for i := 0; i < len(evts); i++ {
		evts[i].Delta = uint32(evts[i].AbsTicks - lasttick)
		lasttick = evts[i].AbsTicks
	}

	return evts
}

type Bar struct {
	Number   uint
	TimeSig  [2]uint8
	Key      *smf.Key // TODO a key change, if != nil TODO: export it to SMF, also move the key changes up to the song
	Tempo    float64  // TODO: move the tempochanges up to the song
	absTicks int64
}

func (b Bar) Len() uint8 {
	return b.TimeSig[0] * 32 / b.TimeSig[1]
}

type Event struct {
	Bar          *Bar
	Pos          uint8 // in 32th
	Duration     uint8 // in 32th for noteOn messages, it is the length of the note, for all other messages, it is 0
	midi.Message       // may only be channel messages or sysex messages. no noteon velocity 0, or noteoff messages, this is expressed via Duration
}

func (e *Event) AbsTicks(ticks smf.MetricTicks) (start, end int64) {
	start = e.Bar.absTicks + int64(ticks.Ticks32th()*uint32(e.Pos))
	if e.Duration <= 0 {
		return start, 0
	}
	end = start + int64(ticks.Ticks32th()*uint32(e.Duration))
	return
}

func FromSMF(sm smf.SMF) *Song {

	switch sm.Format() {
	case 0:
		return fromSMF0(sm)
	case 1:
		return fromSMF1(sm)
	default:
		panic(fmt.Sprintf("SMF format %v is not supported", sm.Format()))
	}
}

func fromSMF0(sm smf.SMF) *Song {
	// by channel
	//sm.TempoChanges()
	s := &Song{}
	return s
}

func fromSMF1(sm smf.SMF) *Song {
	// by track
	//sm.TempoChanges()
	s := &Song{}

	var timesigs []*smf.TrackEvent
	var tempochanges = sm.TempoChanges()
	var keychanges []*smf.TrackEvent

	for _, tr := range sm.Tracks {
		var t Track
		//var te []*smf.TrackEvent
		var absTicks int64
		for _, ev := range tr {
			absTicks += int64(ev.Delta)
			var text string
			var num, denom uint8
			var key smf.Key

			switch {
			case ev.Message.Is(midi.ChannelMsg):
				t.te = append(t.te, &smf.TrackEvent{
					AbsTicks: absTicks,
					Event: smf.Event{
						Message: ev.Message,
					},
				})
			case ev.Message.Is(midi.SysExMsg):
				t.te = append(t.te, &smf.TrackEvent{
					AbsTicks: absTicks,
					Event: smf.Event{
						Message: ev.Message,
					},
				})
			case ev.Message.GetMetaInstrument(&text):
				if t.Name == "" {
					t.Name = text
				}
			case ev.Message.GetMetaTrackName(&text):
				if t.Name == "" {
					t.Name = text
				}
			case ev.Message.GetMetaMeter(&num, &denom):
				timesigs = append(timesigs, &smf.TrackEvent{
					AbsTicks: absTicks,
					Event: smf.Event{
						Message: ev.Message,
					},
				})
			case ev.Message.GetMetaKey(&key):
				keychanges = append(keychanges, &smf.TrackEvent{
					AbsTicks: absTicks,
					Event: smf.Event{
						Message: ev.Message,
					},
				})
			default:
				// ignore
			}
		}

		s.AddTrack(&t)

	}

	var _ = tempochanges
	// TODO
	// 1. set the bars, based on timesig changes
	// 2. set the bar for each event of each track, based on the absticks, and convert the noteoffs to the duration of a note
	// 3. set keychanges and tempochanges

	return s
}

type Song struct {
	Title    string
	Composer string
	tracks   []*Track
	bars     Bars
}

// TODO allow inserting a bar at a certain position (renumber all bars, move the keychanges and tempochanges accordingly
// TODO allow removing a bar (renumber all bars, move the keychanges and tempochanges accordingly
// TODO move a bar (renumber all bars, move the keychanges and tempochanges accordingly
// TODO allow copying a bar with all events from selected tracks (renumber all bars, move the keychanges and tempochanges accordingly
// TODO allow playing some or all tracks
// TODO allow to record to a track

func (s *Song) AddTrack(t *Track) {
	s.tracks = append(s.tracks, t)
}

func (s *Song) AddBar(b *Bar) {
	s.bars = append(s.bars, b)
	s.bars.Renumber()
}

func (s *Song) mkBarLine(ticks smf.MetricTicks) (evts []*smf.TrackEvent, abslength int64) {
	sort.Sort(s.bars)

	var timesig = [2]uint8{4, 4}
	var tempo float64 = 120
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
		if b.Tempo != 0 && b.Tempo != tempo {
			tempo = b.Tempo
			evts = append(evts, &smf.TrackEvent{
				AbsTicks: absTicks,
				Event: smf.Event{
					Message: smf.MetaTempo(b.Tempo),
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

func (s Song) ToSMF1() smf.SMF {
	var sm smf.SMF
	var ticks = smf.MetricTicks(960)
	sm.TimeFormat = ticks
	var t smf.Track
	t.Add(0, smf.MetaText(s.Title))
	t.Add(0, smf.MetaCopyright(s.Composer))
	evts, abslength := s.mkBarLine(ticks)

	for _, tr := range s.tracks {
		evts = append(evts, tr.TrackEvents(ticks)...)
	}

	sort.Sort(byAbsTicks(evts))

	var lasttick int64

	for i := 0; i < len(evts); i++ {
		t.Add(uint32(evts[i].AbsTicks-lasttick), evts[i].Message)
		lasttick = evts[i].AbsTicks
	}

	t.Close(uint32(abslength - lasttick))
	sm.Add(t)

	return sm
}

func (s Song) ToSMF2() smf.SMF {
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

	for _, tr := range s.tracks {
		var t smf.Track
		t.Add(0, smf.MetaTrackSequenceName(tr.Name))
		t.Add(0, smf.MetaInstrument(tr.Name))
		evts := tr.TrackEvents(ticks)
		for _, ev := range evts {
			t.Add(ev.Delta, ev.Message)
		}
		sm.Add(t)
	}

	return sm
}
