package meter

import (
	"testing"

	"github.com/gomidi/midi/midimessage/meta"
)

func TestMeter(t *testing.T) {

	tests := []struct {
		input    meta.TimeSig
		expected string
	}{
		{M2_4(), "meta.TimeSig 2/4 clocksperclick 8 dsqpq 8"},
		{M3_4(), "meta.TimeSig 3/4 clocksperclick 8 dsqpq 8"},
		{M4_4(), "meta.TimeSig 4/4 clocksperclick 8 dsqpq 8"},
		{M5_8(), "meta.TimeSig 5/8 clocksperclick 8 dsqpq 8"},
		{M6_8(), "meta.TimeSig 6/8 clocksperclick 8 dsqpq 8"},
		{M7_8(), "meta.TimeSig 7/8 clocksperclick 8 dsqpq 8"},
		{M12_8(), "meta.TimeSig 12/8 clocksperclick 8 dsqpq 8"},
	}

	for _, test := range tests {

		if got, want := test.input.String(), test.expected; got != want {
			t.Errorf("(%v).String() = %v; want %v", test.input, got, want)
		}
	}

}
