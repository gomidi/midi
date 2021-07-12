package midi_test

import (
	"testing"

	"gitlab.com/gomidi/midi/v2"
)

func TestSysCommon(t *testing.T) {

	tests := []struct {
		msg      []byte
		expected string
	}{
		{
			midi.MTC(3),
			"MTCMsg mtc: 3",
		},
		{
			midi.Tune(),
			"TuneMsg",
		},
		{
			midi.SongSelect(5),
			"SongSelectMsg song: 5",
		},
		{
			midi.SPP(4),
			"SPPMsg spp: 4",
		},
		{
			midi.SPP(4000),
			"SPPMsg spp: 4000",
		},
	}

	for n, test := range tests {
		m := midi.NewMessage(test.msg)

		if got, want := m.String(), test.expected; got != want {
			t.Errorf("[%v] (% X).String() = %#v; want %#v", n, test.msg, got, want)
		}

	}
}
