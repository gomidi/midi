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
	backup           struct {
		cursor           int64 // absticks
		lastDelta        int64 // absticks
		plannedCallbacks plannedCallbacks
	}
}

// Backup saves a backup of cursor, lastDelta and plannedCallbacks
func (t *TimeLine) Backup() {
	t.backup.cursor = t.cursor
	t.backup.lastDelta = t.lastDelta
	t.backup.plannedCallbacks = t.plannedCallbacks
}

// Restore restores cursor, lastDelta and plannedCallbacks from the backup
func (t *TimeLine) Restore() {
	t.cursor = t.backup.cursor
	t.lastDelta = t.backup.lastDelta
	t.plannedCallbacks = t.backup.plannedCallbacks
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
		oldCursor := t.cursor
		t.forwardNBars(nbars)
		t.lastDelta = t.runCallbacks(oldCursor, t.cursor)
	}

	if num > 0 && denom > 0 {
		t.forward(num, denom, true)
	}

}

func (t *TimeLine) forwardIgnoringCallbacks(nbars, num, denom uint32) {
	if nbars > 0 {
		t.forwardNBars(nbars)
	}

	if num > 0 && denom > 0 {
		t.forward(num, denom, false)
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
func (t *TimeLine) toNextBar() {
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

	t.cursor = startOfBar
	t.forward(uint32(num), uint32(denom), false)
}

// forwardNBars checks the bar where the cursor currently is
// and goes n bars ahead and sets the cursor to the start of that bar.
func (t *TimeLine) forwardNBars(nbars uint32) {
	for i := uint32(0); i < nbars; i++ {
		t.toNextBar()
	}
}

func (t *TimeLine) FinishPlanned() {
	t.lastDelta = t.runCallbacks(t.cursor, -1)
	t.cursor = t.lastDelta
}

// runCallbacks runs all planed callbacks until the given absolute position in ticks
// if until is < 0 all remaining callbacks are called
func (t *TimeLine) runCallbacks(from, until int64) (lastPos int64) {
	lastPos = t.lastDelta
	var rest plannedCallbacks

	for _, posCb := range t.plannedCallbacks {
		if posCb.position < from {
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
		//		t.lastDelta = lastPos
		//		t.cursor = lastPos
		posCb.callback(int32(delta))
	}

	//	t.cursor = lastPos
	sort.Sort(rest)
	t.plannedCallbacks = rest
	return
}

// forward sets the cursor forward for the given ratio of whole notes
func (t *TimeLine) forward(num, denom uint32, runCallbacks bool) {
	end := t.cursor + t.Ticks(num, denom)
	if runCallbacks {
		t.lastDelta = t.runCallbacks(t.cursor, end)
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
