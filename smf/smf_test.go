package smf

import (
	"testing"
	"time"
)

func TestTicksDuration(t *testing.T) {
	tests := []struct {
		resolution MetricTicks
		tempo      uint32
		deltaTicks uint32
		duration   time.Duration
	}{
		{96, 120, 96, 500 * time.Millisecond},
		{96, 120, 48, 250 * time.Millisecond},
		{96, 120, 192, 1000 * time.Millisecond},
		{90, 240, 90, 250 * time.Millisecond},
	}

	for _, test := range tests {

		if got, want := test.resolution.Duration(test.tempo, test.deltaTicks), test.duration; got != want {
			t.Errorf(
				"MetricTicks(%v).Duration(%v, %v) = %s; want %s",
				uint16(test.resolution),
				test.tempo,
				test.deltaTicks,
				got,
				want,
			)
		}

		if got, want := test.resolution.Ticks(test.tempo, test.duration), test.deltaTicks; got != want {
			t.Errorf(
				"MetricTicks(%v).Ticks(%v, %v) = %v; want %v",
				uint16(test.resolution),
				test.tempo,
				test.duration,
				got,
				want,
			)
		}
	}

}
