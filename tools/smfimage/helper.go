package smfimage

import (
	"image/color"
	"sort"
	"strconv"

	"gitlab.com/gomidi/midi/tools/smftrack"
)

func floatToInt(x float64) int {
	return int(roundFloat(x, 0))
}

// roundFloat rounds the given float by the given decimals after the dot
func roundFloat(x float64, decimals int) float64 {
	frep := strconv.FormatFloat(x, 'f', decimals, 64)
	f, _ := strconv.ParseFloat(frep, 64)
	return f
}

type trackRange struct {
	*smftrack.Track
	highestKey byte
	lowestKey  byte
	num        int
}

func newTrackRange(tr *smftrack.Track, i int) (rg trackRange) {
	rg.Track = tr
	rg.num = i

	tr.EachEvent(func(ev smftrack.Event) {
		var channel, key, velocity uint8
		switch {
		case ev.Message.GetNoteStart(&channel, &key, &velocity):
			if rg.highestKey == 0 || rg.highestKey < key {
				rg.highestKey = key
			}

			if rg.lowestKey == 0 || rg.lowestKey > key {
				rg.lowestKey = key
			}
		}
	})

	return
}

type trackRanges []trackRange

func (tr trackRanges) Len() int {
	return len(tr)
}

func (tr trackRanges) Less(i, j int) bool {
	if tr[i].lowestKey < tr[j].lowestKey {
		return true
	}

	if tr[i].highestKey < tr[j].highestKey {
		return true
	}

	return false
}

func (tr trackRanges) Swap(i, j int) {
	tr[i], tr[j] = tr[j], tr[i]
}

func sortByRange(tr ...*smftrack.Track) (trackOrder []int) {
	trackOrder = make([]int, len(tr))

	rg := make(trackRanges, len(tr))

	for i, t := range tr {
		rg[i] = newTrackRange(t, i)
	}

	sort.Sort(rg)

	for i, r := range rg {
		trackOrder[i] = r.num
	}

	return

}

type coloredNote struct {
	pitch byte
	//rgbColor rgbColor
	color    color.RGBA
	velocity byte
}
