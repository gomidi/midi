package midi_test

import (
	"testing"

	"gitlab.com/gomidi/midi/v2"
)

func TestRealTime(t *testing.T) {

	tests := []struct {
		msg      []byte
		expected string
	}{
		{
			midi.NewTimingClock(),
			"TimingClock",
		},
		{
			midi.NewTick(),
			"Tick",
		},
		{
			midi.NewStart(),
			"Start",
		},
		{
			midi.NewContinue(),
			"Continue",
		},
		{
			midi.NewStop(),
			"Stop",
		},
		/*
			{
				midi.NewUndefined(),
				"UnknownType",
			},
		*/
		{
			midi.NewActivesense(),
			"ActiveSense",
		},
		{
			midi.NewReset(),
			"Reset",
		},
	}

	for n, test := range tests {
		m := midi.Message(test.msg)

		if got, want := m.String(), test.expected; got != want {
			t.Errorf("[%v] (% X).String() = %#v; want %#v", n, test.msg, got, want)
		}

	}
}
