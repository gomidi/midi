package smfimage

type timeSignatureChange struct {
	Numerator   uint8
	Denominator uint8
	_32th       int
}

type timeSignatureChanges []timeSignatureChange

func (ts timeSignatureChanges) Less(i, j int) bool {
	return ts[i]._32th < ts[j]._32th
}

func (ts timeSignatureChanges) Len() int {
	return len(ts)
}

func (ts timeSignatureChanges) Swap(i, j int) {
	ts[i], ts[j] = ts[j], ts[i]
}
