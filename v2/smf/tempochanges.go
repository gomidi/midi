package smf

type TempoChange struct {
	AbsTicks        int64
	AbsTimeMicroSec int64
	BPM             float64
}

type TempoChanges []*TempoChange

func (t TempoChanges) Swap(a, b int) {
	t[a], t[b] = t[b], t[a]
}

func (t TempoChanges) Len() int {
	return len(t)
}

func (t TempoChanges) Less(a, b int) bool {
	return t[a].AbsTicks < t[b].AbsTicks
}

func (t TempoChanges) TempoAt(absTicks int64) (bpm float64) {
	tc := t.TempoChangeAt(absTicks)
	if tc == nil {
		return 120.00
	}
	return tc.BPM
}

func (t TempoChanges) TempoChangeAt(absTicks int64) (tch *TempoChange) {
	for _, tc := range t {
		if tc.AbsTicks > absTicks {
			break
		}
		tch = tc
	}
	return
}
