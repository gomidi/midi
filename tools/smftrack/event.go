package smftrack

import (
	"time"

	"gitlab.com/gomidi/midi/v2/smf"
)

// Event is a MIDI message at absolute ticks in a track.
type Event struct {
	AbsTicks uint64
	smf.Message
	no uint64
}

// DurationTo returns the duration to a given target Event based on the tick resolution and the given tempo
func (e Event) DurationTo(resolution smf.MetricTicks, tempoBPM float64, target Event) time.Duration {
	return resolution.Duration(tempoBPM, uint32(target.AbsTicks-e.AbsTicks))
}

// TicksTo returns the absticks to the given target duration, based on the given tempo and the resolution
func (e Event) TicksTo(resolution smf.MetricTicks, tempoBPM float64, timeDistance time.Duration) uint64 {
	return e.AbsTicks + uint64(resolution.Ticks(tempoBPM, timeDistance))
}

// Number returns the number of the event as part of a track. (If it is 0, the event has not been part of a track).
func (e Event) Number() uint64 {
	return e.no
}

// Events helps sorting events
type Events []Event

func (e Events) Len() int {
	return len(e)
}

func (e Events) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Events) Less(i, j int) bool {
	if e[i].AbsTicks > e[j].AbsTicks {
		return false
	}
	if e[i].AbsTicks < e[j].AbsTicks {
		return true
	}

	// end of track always comes last
	if e[i].Message.Is(smf.MetaEndOfTrackMsg) {
		return true
	}

	if e[j].Message.Is(smf.MetaEndOfTrackMsg) {
		return false
	}

	// then tempo comes last
	if e[i].Message.Is(smf.MetaTempoMsg) {
		return true
	}

	if e[j].Message.Is(smf.MetaTempoMsg) {
		return false
	}

	/*
		if t1 == "channel.NoteOn" && t2 == "channel.NoteOn" {
			n1 := e[i].Message.(channel.NoteOn)
			n2 := e[j].Message.(channel.NoteOn)
			return (n1.Pitch() / 12) < (n2.Pitch() / 12)
		}
	*/

	// rest doesn't matter
	return false

}
