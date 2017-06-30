package smftrack

import (
	"fmt"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/meta"
)

type Event struct {
	AbsTicks uint64
	midi.Message
	no uint64
}

// Only events that are inside a track have a number
func (e *Event) Number() uint64 {
	return e.no
}

// Events helps sorting events
type Events []*Event

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

	// rest doesn't matter
	return false

}
