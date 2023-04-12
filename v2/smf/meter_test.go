package smf

import (
	"testing"
)

func TestMeter(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{MetaMeter(2, 4), "MetaTimeSig meter: 2/4"},
		{MetaMeter(3, 4), "MetaTimeSig meter: 3/4"},
		{MetaMeter(4, 4), "MetaTimeSig meter: 4/4"},
		{MetaMeter(5, 8), "MetaTimeSig meter: 5/8"},
		{MetaMeter(6, 8), "MetaTimeSig meter: 6/8"},
		{MetaMeter(7, 8), "MetaTimeSig meter: 7/8"},
		{MetaMeter(12, 8), "MetaTimeSig meter: 12/8"},
	}

	for _, test := range tests {

		if got, want := test.input.String(), test.expected; got != want {
			t.Errorf("(%v).String() = %v; want %v", test.input, got, want)
		}
	}

}
