package mid

import (
	"bytes"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfwriter"
	"testing"
	"time"
)

func TestTimeAt(t *testing.T) {
	mt := smf.MetricTicks(960)
	twobars := mt.Ticks4th() * 8

	tests := []struct {
		tempo1          uint32
		tempo2          uint32
		absPos          uint32
		durationSeconds int64
	}{

		{120, 120, twobars * 2, 8},
		{120, 120, twobars, 4},
		{120, 120, twobars / 2, 2},

		{60, 60, twobars * 2, 16},
		{60, 60, twobars, 8},
		{60, 60, twobars / 2, 4},

		{120, 60, twobars * 2, 12},
		{120, 30, twobars * 2, 20},
		{120, 30, twobars * 3, 36},

		{120, 30, twobars, 4},
	}

	for _, test := range tests {
		var bf bytes.Buffer

		wr := NewSMFWriter(&bf, 1, smfwriter.TimeFormat(mt))
		wr.Tempo(test.tempo1)
		wr.SetDelta(twobars)
		wr.Tempo(test.tempo2)
		wr.SetDelta(twobars)
		wr.SetDelta(twobars * 8)
		wr.NoteOn(64, 120)
		wr.SetDelta(twobars)
		wr.NoteOff(64)
		wr.EndOfTrack()

		h := NewHandler(NoLogger())
		h.ReadSMF(&bf)
		d := *h.TimeAt(uint64(test.absPos))
		// ms := int64(d / time.Millisecond)

		if got, want := int64(d/time.Second), test.durationSeconds; got != want {
			t.Errorf("tempo1, tempo2 = %v, %v; TimeAt(%v) = %v; want %v", test.tempo1, test.tempo2, test.absPos, got, want)
		}
	}

}
