package sequencer

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2/smf"
)

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
	start = b.AbsTicks + int64(ticks.Ticks32th()*uint32(e.Pos))
	if e.Duration <= 0 {
		return start, 0
	}
	end = start + int64(ticks.Ticks32th()*uint32(e.Duration))
	return
}
