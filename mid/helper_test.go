package mid

import (
	// "github.com/metakeule/fmtdate"
	"testing"
	"time"
)

func quarternoteLength(tempo uint32) time.Duration {
	return (time.Second * 60) / time.Duration(tempo)
}

func thirtenthnoteLength(tempo uint32) time.Duration {
	return (time.Second * 60) / time.Duration(tempo*8)
}

func clock(bpm uint32) time.Duration {
	return thirtenthnoteLength(bpm) / 3 // three MIDIclocks make a 32th
}

func TestCalcTempoBasedOnMIDIClocks(t *testing.T) {
	format := "5.0000"
	now := time.Now().Round(time.Minute)

	tests := []struct {
		a        time.Time
		b        time.Time
		c        time.Time
		d        time.Time
		expected float64
	}{
		{now, now.Add(clock(120)), now.Add(clock(120)).Add(clock(120)), now.Add(clock(120)).Add(clock(120)).Add(clock(120)), 120},
		{now, now.Add(clock(130)), now.Add(clock(130)).Add(clock(130)), now.Add(clock(130)).Add(clock(130)).Add(clock(130)), 130},
		{now, now.Add(clock(120)), now.Add(clock(120)).Add(clock(118)), now.Add(clock(120)).Add(clock(118)).Add(clock(123)), 120},
	}

	for _, test := range tests {

		if got, want := tempoBasedOnMIDIClocks(&test.a, &test.b, &test.c, &test.d), test.expected; got != want {
			t.Errorf("tempoBasedOnMIDIClocks(%v,%v,%v,%v) = %v; want %v",
				test.a.Format(format),
				test.b.Format(format),
				test.c.Format(format),
				test.d.Format(format),
				got, want)
		}
	}

}
