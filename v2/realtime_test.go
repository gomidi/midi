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
			midi.TimingClock(),
			"TimingClock",
		},
		{
			midi.Tick(),
			"Tick",
		},
		{
			midi.Start(),
			"Start",
		},
		{
			midi.Continue(),
			"Continue",
		},
		{
			midi.Stop(),
			"Stop",
		},
		/*
			{
				midi.NewUndefined(),
				"UnknownType",
			},
		*/
		{
			midi.Activesense(),
			"ActiveSense",
		},
		{
			midi.Reset(),
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
