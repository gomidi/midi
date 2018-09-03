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

		// too high values
		{
			Channel1.Aftertouch(130),
			"channel.Aftertouch channel 1 pressure 127",
		},
		{
			Channel8.ControlChange(137, 130),
			"channel.ControlChange channel 8 controller 127 (\"Poly Operation\") value 127",
		},
		{
			Channel2.NoteOn(130, 130),
			"channel.NoteOn channel 2 key 127 velocity 127",
		},
		{
			Channel3.NoteOff(180),
			"channel.NoteOff channel 3 key 127",
		},
		{
			Channel4.NoteOffVelocity(180, 220),
			"channel.NoteOffVelocity channel 4 key 127 velocity 127",
		},
		{
			Channel4.Pitchbend(12300),
			"channel.Pitchbend channel 4 value 8191 absValue 0",
		},
		{
			Channel4.PolyAftertouch(186, 190),
			"channel.PolyAftertouch channel 4 key 127 pressure 127",
		},
		{
			Channel4.ProgramChange(183),
			"channel.ProgramChange channel 4 program 127",
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

func TestSetChannel(t *testing.T) {

	tests := []struct {
		input     Message
		toChannel uint8
		expected  string
	}{
		{
			Channel1.Aftertouch(120),
			5,
			"channel.Aftertouch channel 5 pressure 120",
		},
		{
			Channel8.ControlChange(7, 110),
			9,
			"channel.ControlChange channel 9 controller 7 (\"Volume (MSB)\") value 110",
		},
		{
			Channel2.NoteOn(100, 80),
			0,
			"channel.NoteOn channel 0 key 100 velocity 80",
		},
		{
			Channel3.NoteOff(80),
			2,
			"channel.NoteOff channel 2 key 80",
		},
		{
			Channel4.NoteOffVelocity(80, 20),
			11,
			"channel.NoteOffVelocity channel 11 key 80 velocity 20",
		},
		{
			Channel4.Pitchbend(300),
			14,
			"channel.Pitchbend channel 14 value 300 absValue 0",
		},
		{
			Channel4.PolyAftertouch(86, 109),
			2,
			"channel.PolyAftertouch channel 2 key 86 pressure 109",
		},
		{
			Channel4.ProgramChange(83),
			0,
			"channel.ProgramChange channel 0 program 83",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		msg := SetChannel(test.input, test.toChannel)

		bf.WriteString(msg.String())

		if got, want := bf.String(), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}
