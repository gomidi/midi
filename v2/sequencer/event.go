package sequencer

import (
	"bytes"
	"fmt"
	"math"
	"sort"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

// NTupleEvent are AtomEvents spread evenly over the time
type NTupleEvent struct {
	Duration uint8
	evts     []*AtomEvent
}

func (e *NTupleEvent) Len() int {
	return len(e.evts)
}

func (e *NTupleEvent) durSingle() float64 {
	return float64(e.Duration) / float64(len(e.evts))
}

func (e *NTupleEvent) Swap(a, b int) {
	e.evts[a], e.evts[b] = e.evts[b], e.evts[a]
}

// warning: only works within a bar
func (e *NTupleEvent) Less(a, b int) bool {
	return e.evts[a].Position() < e.evts[b].Position()
}

// TODO: test
func (ntev *NTupleEvent) TrackEvents(b *Bar, ticks smf.MetricTicks) (evts smf.TrackEvents) {
	if len(ntev.evts) == 0 {
		return nil
	}

	start, _ := ntev.evts[0].AbsTicks(b, ticks)
	durSingle := ntev.durSingle()

	for _, ev := range ntev.evts {
		_, end := ev.AbsTicks(b, ticks)
		evts = append(evts, &smf.TrackEvent{
			AbsTicks: start,
			Event: smf.Event{
				Message: smf.Message(ev.Message.Bytes()),
			},
			TrackNo: ev.TrackNo,
		})

		var channel, key, velocity uint8
		var nextStart = int64(math.Round(float64(start) + durSingle))
		if ev.Message.GetNoteStart(&channel, &key, &velocity) && end != 0 {
			evts = append(evts, &smf.TrackEvent{
				AbsTicks: nextStart,
				Event: smf.Event{
					Message: smf.Message(midi.NoteOff(channel, key)),
				},
				TrackNo: ev.TrackNo,
			})
		}
		start = nextStart
	}
	sort.Sort(evts)
	return evts
}

func (ntev *NTupleEvent) Position() uint8 {
	if len(ntev.evts) == 0 {
		return 0
	}
	return ntev.evts[0].Position()
}

func (ntev *NTupleEvent) Inspect() string {
	var bf bytes.Buffer
	bf.WriteString("{")

	for _, ev := range ntev.evts {
		bf.WriteString(ev.Inspect() + ", ")
	}

	bf.WriteString("}")

	return bf.String()
}

var _ Event = &NTupleEvent{}

// MultiEvent are multiple AtomEvents at the same time
type MultiEvent []*AtomEvent

var _ Event = &MultiEvent{}

// TODO: test
func (mev MultiEvent) TrackEvents(b *Bar, ticks smf.MetricTicks) (evts smf.TrackEvents) {
	if len(mev) == 0 {
		return nil
	}

	start, _ := mev[0].AbsTicks(b, ticks)

	for _, ev := range mev {
		_, end := ev.AbsTicks(b, ticks)
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
	return evts
}

func (mev MultiEvent) Position() uint8 {
	if len(mev) == 0 {
		return 0
	}
	return mev[0].Position()
}

func (mev MultiEvent) Inspect() string {
	var bf bytes.Buffer
	bf.WriteString("[")

	for _, ev := range mev {
		bf.WriteString(ev.Inspect() + ", ")
	}

	bf.WriteString("]")

	return bf.String()
}

func (e MultiEvent) Len() int {
	return len(e)
}

func (e MultiEvent) Swap(a, b int) {
	e[a], e[b] = e[b], e[a]
}

// warning: only works within a bar
func (e MultiEvent) Less(a, b int) bool {
	return e[a].Pos < e[b].Pos
}

type EmptyEvent uint8

func (ev EmptyEvent) Position() uint8 {
	return uint8(ev)
}

func (ev EmptyEvent) TrackEvents(b *Bar, ticks smf.MetricTicks) (evts smf.TrackEvents) {
	return nil
}

func (e EmptyEvent) Inspect() string {
	return ""
}

var _ Event = EmptyEvent(0)

type AtomEvent struct {
	TrackNo  int
	Pos      uint8       // in 32th
	Duration uint8       // in 32th for noteOn messages, it is the length of the note, for all other messages, it is 0
	Message  smf.Message // may only be channel messages or sysex messages. no noteon velocity 0, or noteoff messages, this is expressed via Duration
	absTicks int64       // just for smf import
}

func (ev *AtomEvent) Position() uint8 {
	return ev.Pos
}

func (ev *AtomEvent) TrackEvents(b *Bar, ticks smf.MetricTicks) (evts smf.TrackEvents) {
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

	return evts
}

var _ Event = &AtomEvent{}

func (e AtomEvent) Inspect() string {
	return fmt.Sprintf("Event{TrackNo:%v, Pos:%v, Duration:%v, Message: %s, absTicks: %v}", e.TrackNo, int(e.Pos), int(e.Duration), e.Message.String(), e.absTicks)
}

func (e *AtomEvent) AbsTicks(b *Bar, ticks smf.MetricTicks) (start, end int64) {
	start = b.AbsTicks + int64(ticks.Ticks32th()*uint32(e.Pos))
	if e.Duration <= 0 {
		return start, 0
	}
	end = start + int64(ticks.Ticks32th()*uint32(e.Duration))
	return
}
