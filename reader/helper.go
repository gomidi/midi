package reader

import (
	// "fmt"
	"time"
)

// tempoBasedOnMIDIClocks returns a tempo calculated by 4 timestamps of midi clock events in a row
// a is the first and d is the last timestamp
func tempoBasedOnMIDIClocks(a, b, c, d *time.Time) float64 {
	// the simplest way to do this is to build the difference in time to the last midi clock
	// 12 clocks = 8th, 6 = 16th, 3 = 32th
	// lets say we always takt the last 3 midi clocks their timespan must be a 32th note
	// bpm = qn/minute = 8 * 32th/60 sec
	// 3clocks = 32th <=> bpm = 8 * 3clocks/60 sec
	//

	// here we got three midi clocks, so calc the time spans
	last := d.Sub(*c)
	before := c.Sub(*b)
	bebefore := b.Sub(*a)
	one32th := (last + before + bebefore)
	tempo := (time.Second * 60 / (one32th * 8)).Nanoseconds()
	return float64(tempo)
}
