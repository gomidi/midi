package channel

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMessagesString(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			Channel1.Aftertouch(120),
			"channel.Aftertouch channel 1 pressure 120",
		},
		{
			Channel8.ControlChange(7, 110),
			"channel.ControlChange channel 8 controller 7 (\"Volume (MSB)\") value 110",
		},
		{
			Channel2.NoteOn(100, 80),
			"channel.NoteOn channel 2 key 100 velocity 80",
		},
		{
			Channel3.NoteOff(80),
			"channel.NoteOff channel 3 key 80",
		},
		{
			Channel4.NoteOffVelocity(80, 20),
			"channel.NoteOffVelocity channel 4 key 80 velocity 20",
		},
		{
			Channel4.Pitchbend(300),
			"channel.Pitchbend channel 4 value 300 absValue 0",
		},
		{
			Channel4.PolyAftertouch(86, 109),
			"channel.PolyAftertouch channel 4 key 86 pressure 109",
		},
		{
			Channel4.ProgramChange(83),
			"channel.ProgramChange channel 4 program 83",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		bf.WriteString(test.input.String())

		if got, want := bf.String(), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}

func TestMessagesRaw(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			Channel1.Aftertouch(120),
			"D1 78",
		},
		{
			Channel8.ControlChange(7, 110),
			"B8 07 6E",
		},
		{
			Channel2.NoteOn(100, 80),
			"92 64 50",
		},
		{
			Channel3.NoteOff(80),
			"93 50 00",
		},
		{
			Channel4.NoteOffVelocity(80, 20),
			"84 50 14",
		},
		{
			Channel4.Pitchbend(300),
			"E4 2C 42",
		},
		{
			Channel4.PolyAftertouch(86, 109),
			"A4 56 6D",
		},
		{
			Channel4.ProgramChange(83),
			"C4 53",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		bf.Write(test.input.Raw())

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}
