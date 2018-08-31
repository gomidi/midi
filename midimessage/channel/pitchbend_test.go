package channel_test

import (
	"bytes"
	"testing"

	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midireader"
	"github.com/gomidi/midi/midiwriter"
)

func TestPitchbend(t *testing.T) {

	tests := []struct {
		in       int16
		expected uint16
	}{
		{
			in:       0,
			expected: 8192,
		},
		{
			in:       channel.PitchHighest,
			expected: 16383,
		},
		{
			in:       channel.PitchLowest,
			expected: 0,
		},
	}

	for _, test := range tests {
		var bf bytes.Buffer

		wr := midiwriter.New(&bf)
		rd := midireader.New(&bf, nil)

		wr.Write(channel.Channel0.Pitchbend(test.in))
		msg, _ := rd.Read()

		got := msg.(channel.Pitchbend).AbsValue()

		if got != test.expected {
			t.Errorf("Pitchbend(%v).absValue = %v; wanted %v", test.in, got, test.expected)
		}
	}
}
