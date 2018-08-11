package meter

import (
	"github.com/gomidi/midi/midimessage/meta"
	"testing"
)

func TestMeter(t *testing.T) {

	tests := []struct {
		input    meta.TimeSignature
		expected string
	}{
		{M2_4(), "meta.TimeSignature 2/4 clocksperclick 8 dsqpq 8"},
		{M3_4(), "meta.TimeSignature 3/4 clocksperclick 8 dsqpq 8"},
		{M4_4(), "meta.TimeSignature 4/4 clocksperclick 8 dsqpq 8"},
		{M5_8(), "meta.TimeSignature 5/8 clocksperclick 8 dsqpq 8"},
		{M6_8(), "meta.TimeSignature 6/8 clocksperclick 8 dsqpq 8"},
		{M7_8(), "meta.TimeSignature 7/8 clocksperclick 8 dsqpq 8"},
		{M12_8(), "meta.TimeSignature 12/8 clocksperclick 8 dsqpq 8"},
	}

	for _, test := range tests {

		if got, want := test.input.String(), test.expected; got != want {
			t.Errorf("(%v).String() = %v; want %v", test.input, got, want)
		}
	}

}
