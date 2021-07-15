package smf

type TempoChange struct {
	AbsDelta int64
	BPM      float64
}

type TempoChanges []TempoChange

func (t TempoChanges) Swap(a, b int) {
	t[a], t[b] = t[b], t[a]
}

func (t TempoChanges) Len() int {
	return len(t)
}

func (t TempoChanges) Less(a, b int) bool {
	return t[a].AbsDelta < t[b].AbsDelta
}

func (t TempoChanges) TempoAt(absDelta int64) (bpm float64) {
	bpm = 120.00
	for _, tc := range t {
		if tc.AbsDelta > absDelta {
			break
		}
		bpm = tc.BPM
	}
	return
}
