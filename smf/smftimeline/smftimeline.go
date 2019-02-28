package smftimeline

import (
	"math"
	"sort"

	"gitlab.com/gomidi/midi/smf"
)

func New(ticks smf.MetricTicks) *TimeLine {
	return &TimeLine{
		ticks: ticks,
	}
}

type TimeLine struct {
	ticks            smf.MetricTicks
	timeSigs         [][3]int64 // first: abs ticks, second: numerator, third: denominator
	tempoChanges     [][2]int64 // first: abs ticks, second: bpm
	cursor           int64      // absticks
	lastDelta        int64      // absticks
	plannedCallbacks plannedCallbacks
}

type plannedCallbacks []plannedCallback

func (p plannedCallbacks) Len() int {
	return len(p)
}

func (p plannedCallbacks) Swap(a, b int) {
	p[a], p[b] = p[b], p[a]
}

func (p plannedCallbacks) Less(a, b int) bool {
	return p[a].position < p[b].position
}

type plannedCallback struct {
	callback func(delta int32)
	position int64 // absticks
}

// Forward checks the bar where the cursor currently is
// and goes nbars ahead and moves the cursor to the start of that bar.
// It then moves the cursor forward within the target bar for the given ratio of whole notes.
func (t *TimeLine) Forward(nbars, num, denom uint32) {
	if nbars > 0 {
		t.forwardNBars(nbars, true)
	}

	if num > 0 && denom > 0 {
		t.forward(num, denom, t.cursor, true)
	}
}

func (t *TimeLine) forwardIgnoringCallbacks(nbars, num, denom uint32) {
	if nbars > 0 {
		t.forwardNBars(nbars, false)
	}

	if num > 0 && denom > 0 {
		t.forward(num, denom, t.cursor, false)
	}
}

// Plan registers the given callback to be invoked, when the cursor moves to the
// position given by the delta resulting from nbars, num, denom from the current
// cursor position.
// I.e. when Forward is invoked, it checks, if it moves the cursor across planned callbacks
// and if so they are called. However these callbacks must not move the cursor or set a delta.
// They should simply write midi events to some smf.Writer.
func (t *TimeLine) Plan(nbars, num, denom uint32, callback func(delta int32)) {
	/*
	   1. calc the abs position for the callback by using forward
	   2. rewind cursor
	   3. register callback
	   4. sort planned callbacks
	*/

	savedCursor := t.cursor
	t.forwardIgnoringCallbacks(nbars, num, denom)
	pos := t.cursor
	t.cursor = savedCursor
	//	fmt.Printf("cursor: %v, pos: %v\n", t.cursor, pos)
	t.plannedCallbacks = append(t.plannedCallbacks, plannedCallback{callback: callback, position: pos})
	sort.Sort(t.plannedCallbacks)
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

// goes ahead and sets the cursor to the start of the next bar.
func (t *TimeLine) toNextBar(runCallbacks bool) {
	var num, denom int64 = 4, 4
	//	var idx int
	var startOfBar int64

	for _, timeSig := range t.timeSigs {
		if timeSig[0] <= t.cursor {
			//			println("timeSig[0] <= t.cursor")
			//			idx = i
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

	if runCallbacks {
		t.runCallbacks(startOfBar)
	}
	//t.cursor = startOfBar

	t.forward(uint32(num), uint32(denom), startOfBar, runCallbacks)
}

// forwardNBars checks the bar where the cursor currently is
// and goes n bars ahead and sets the cursor to the start of that bar.
func (t *TimeLine) forwardNBars(nbars uint32, runCallbacks bool) {
	for i := uint32(0); i < nbars; i++ {
		t.toNextBar(runCallbacks)
	}
}

func (t *TimeLine) FinishPlanned() {
	t.runCallbacks(-1)
}

// runCallbacks runs all planed callbacks until the given absolute position in ticks
// if until is < 0 all remaining callbacks are called
func (t *TimeLine) runCallbacks(until int64) {
	lastPos := t.cursor
	var rest plannedCallbacks

	for _, posCb := range t.plannedCallbacks {
		if posCb.position < t.cursor {
			continue
		}

		if until >= 0 {
			if posCb.position > until {
				rest = append(rest, posCb)
				continue
			}
		}

		delta := posCb.position - lastPos
		//		fmt.Printf("callback position %v, lastPos %v delta %v cursor %v\n", posCb.position, lastPos, delta, t.cursor)
		lastPos = posCb.position
		posCb.callback(int32(delta))
	}
	t.lastDelta = lastPos
	sort.Sort(rest)
	t.plannedCallbacks = rest
}

// forward sets the cursor forward for the given ratio of whole notes
func (t *TimeLine) forward(num, denom uint32, starter int64, runCallbacks bool) {
	end := starter + t.Ticks(num, denom)
	if runCallbacks {
		t.runCallbacks(end)
	}
	t.cursor = end
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
