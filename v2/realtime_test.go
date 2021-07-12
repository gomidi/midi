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
			"TimingClockMsg",
		},
		{
			midi.Tick(),
			"TickMsg",
		},
		{
			midi.Start(),
			"StartMsg",
		},
		{
			midi.Continue(),
			"ContinueMsg",
		},
		{
			midi.Stop(),
			"StopMsg",
		},
		{
			midi.Undefined(),
			"UndefinedMsg",
		},
		{
			midi.Activesense(),
			"ActiveSenseMsg",
		},
		{
			midi.Reset(),
			"ResetMsg",
		},
	}

	for n, test := range tests {
		m := midi.NewMessage(test.msg)

		if got, want := m.String(), test.expected; got != want {
			t.Errorf("[%v] (% X).String() = %#v; want %#v", n, test.msg, got, want)
		}

	}
}
