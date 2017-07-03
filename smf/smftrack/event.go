package smftrack

import (
	"fmt"
	"github.com/gomidi/midi"
	// "github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
	"time"
)

// Event is a MIDI message at absolute ticks in a track.
type Event struct {
	AbsTicks uint64
	midi.Message
	no uint64
}

// DurationTo returns the duration to a given target Event based on the tick resolution and the given tempo
func (e Event) DurationTo(resolution smf.MetricTicks, tempoBPM uint32, target Event) time.Duration {
	return resolution.TempoDuration(tempoBPM, uint32(target.AbsTicks-e.AbsTicks))
}

// TicksTo returns the absticks to the given target duration, based on the given tempo and the resolution
func (e Event) TicksTo(resolution smf.MetricTicks, tempoBPM uint32, timeDistance time.Duration) uint64 {
	return e.AbsTicks + uint64(resolution.TempoTicks(tempoBPM, timeDistance))
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
	if e[i].Message == meta.EndOfTrack {
		return true
	}

	if e[j].Message == meta.EndOfTrack {
		return false
	}

	t1 := fmt.Sprintf("%T", e[i].Message)
	t2 := fmt.Sprintf("%T", e[j].Message)

	// then tempo comes last
	if t1 == "meta.Tempo" {
		return true
	}

	if t2 == "meta.Tempo" {
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
