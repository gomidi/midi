package sequencer

import (
	"math"
	"sort"

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

type Events []Event

func (e Events) Len() int {
	return len(e)
}

func (e Events) Swap(a, b int) {
	e[a], e[b] = e[b], e[a]
}

// warning: only works within a bar
func (e Events) Less(a, b int) bool {
	return e[a].Position() < e[b].Position()
}

type Event interface {
	Position() uint8
	TrackEvents(b *Bar, ticks smf.MetricTicks) (evts smf.TrackEvents)
	Inspect() string
}

type Bar struct {
	Number  uint
	TimeSig [2]uint8
	Events   Events
	Key      *smf.Key // TODO a key change, if != nil
	AbsTicks int64
}

func (b *Bar) SortEvents() {
	evts := b.Events
	sort.Sort(evts)
	b.Events = evts
}

func (b *Bar) AddEvent(ev Event) {
	b.Events = append(b.Events, ev)
}

func (b Bar) Len() uint8 {
	return b.TimeSig[0] * 32 / b.TimeSig[1]
}

func (b *Bar) barPos(absTicks int64, ticks smf.MetricTicks) uint8 {
	diff := absTicks - b.AbsTicks
	t32th := ticks.Ticks32th()
	return uint8(math.Round(float64(diff) / float64(t32th)))
}

func (b *Bar) TrackEvents(ticks smf.MetricTicks) (evts smf.TrackEvents) {
	for _, tevts := range b.Events {
		evts = append(evts, tevts.TrackEvents(b, ticks)...)
	}

	sort.Sort(evts)

	var lasttick int64

	for i := 0; i < len(evts); i++ {
		evts[i].Delta = uint32(evts[i].AbsTicks - lasttick)
		lasttick = evts[i].AbsTicks
	}

	return evts
}

/*
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
*/
