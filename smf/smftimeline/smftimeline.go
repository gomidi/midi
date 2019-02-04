package smftimeline

import (
	"math"

	"gitlab.com/gomidi/midi/smf"
)

func New(ticks smf.MetricTicks) *TimeLine {
	return &TimeLine{
		ticks: ticks,
	}
}

type TimeLine struct {
	ticks        smf.MetricTicks
	timeSigs     [][3]int64 // first: abs ticks, second: numerator, third: denominator
	tempoChanges [][2]int64 // first: abs ticks, second: bpm
	cursor       int64      // absticks
	lastDelta    int64      // absticks
}

// AddTimeSignature adds the given timesignature at the current cursor position
func (t *TimeLine) AddTimeSignature(num, denom uint8) {
	t.timeSigs = append(t.timeSigs, [3]int64{t.cursor, int64(num), int64(denom)})
}

// Reset resets the cursor and last delta
func (t *TimeLine) Reset() {
	t.cursor = 0
	t.lastDelta = 0
}

func (t *TimeLine) Ticks(num, denom uint32) int64 {
	return int64(math.Round((float64(t.ticks.Ticks4th()) * 4.0 * float64(num)) / float64(denom)))
}

// ForwardNBars checks the bar where the cursor currently is
// and goes n bars ahead and sets the cursor to the start of that bar.
func (t *TimeLine) ForwardNBars(nbars uint32) {
	var num, denom int64 = 4, 4
	var idx int
	var startOfBar int64

	for i, timeSig := range t.timeSigs {
		if timeSig[0] <= t.cursor {
			//			println("timeSig[0] <= t.cursor")
			idx = i
			startOfBar = timeSig[0]
			num = timeSig[1]
			denom = timeSig[2]
		} else {
			//			println("break")
			break
		}
	}

	//	println("start Of Bar with time signature", startOfBar)

	/*
		startOfBar is the start of the last bar that had the time signature.
		we want to find the start of the bar that contains the cursor.
	*/
	if t.cursor > startOfBar {
		diffTicks := t.cursor - startOfBar
		barLenTicks := t.Ticks(uint32(num), uint32(denom))
		if diffTicks >= barLenTicks {
			no := diffTicks / barLenTicks // removing the rest
			startOfBar += no * barLenTicks
			//			println("startOfBar corrected", startOfBar)
		}
	}

	t.cursor = startOfBar

	/*
		now check where the next time Signature change is
		then we advance bar by bar by the signature of the last bar
		until we either did the nbars or we got a different timesig
	*/
	for i := uint32(0); i < nbars; i++ {
		//println("i Forward", i, num, denom)
		t.Forward(uint32(num), uint32(denom))
		//println("idx  len(t.timeSigs)", idx, len(t.timeSigs))
		if idx < len(t.timeSigs) && t.timeSigs[idx][0] <= t.cursor {
			num = t.timeSigs[idx][1]
			denom = t.timeSigs[idx][2]
			idx++
		}
	}

}

// Forward sets the cursor forward for the given ratio of whole notes
func (t *TimeLine) Forward(num, denom uint32) {
	t.cursor += t.Ticks(num, denom)
}

// GetDelta returns the delta of the current cursor position
// to the last delta position and sets the last delta position to the current cursor position
// returns -1, if cursor is before last delta
func (t *TimeLine) GetDelta() int32 {
	if t.cursor < t.lastDelta {
		t.lastDelta = t.cursor
		return -1
	}

	res := int32(t.cursor - t.lastDelta)
	t.lastDelta = t.cursor
	return res
}
