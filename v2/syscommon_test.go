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
			midi.NewMTC(3),
			"MTC mtc: 3",
		},
		{
			midi.NewTune(),
			"Tune",
		},
		{
			midi.NewSongSelect(5),
			"SongSelect song: 5",
		},
		{
			midi.NewSPP(4),
			"SPP spp: 4",
		},
		{
			midi.NewSPP(4000),
			"SPP spp: 4000",
		},
	}

	for n, test := range tests {
		m := midi.NewMessage(test.msg)

		if got, want := m.String(), test.expected; got != want {
			t.Errorf("[%v] (% X).String() = %#v; want %#v", n, test.msg, got, want)
		}

	}
}
