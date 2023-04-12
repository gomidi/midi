package midi

import (
	"bytes"
	"fmt"
	"testing"
)

func TestChannelString(t *testing.T) {

	tests := []struct {
		input    []byte
		expected string
	}{
		{
			AfterTouch(1, 120),
			"AfterTouch channel: 1 pressure: 120",
		},
		{
			ControlChange(8, 7, 110),
			//"ControlChange channel: 8 controller: 7 (\"Volume (MSB)\") value 110",
			"ControlChange channel: 8 controller: 7 value: 110",
		},
		{
			NoteOn(2, 100, 80),
			"NoteOn channel: 2 key: 100 velocity: 80",
		},
		{
			NoteOff(3, 80),
			"NoteOff channel: 3 key: 80",
		},
		{
			NoteOffVelocity(4, 80, 20),
			"NoteOff channel: 4 key: 80 velocity: 20",
		},
		{
			Pitchbend(4, 300),
			"PitchBend channel: 4 pitch: 300 (8492)",
		},
		{
			PolyAfterTouch(4, 86, 109),
			"PolyAfterTouch channel: 4 key: 86 pressure: 109",
		},
		{
			ProgramChange(4, 83),
			"ProgramChange channel: 4 program: 83",
		},

		// too high values
		{
			AfterTouch(1, 130),
			"AfterTouch channel: 1 pressure: 127",
		},
		{
			ControlChange(8, 137, 130),
			//"ControlChange channel: 8 controller: 127 (\"Poly Operation\") value 127",
			"ControlChange channel: 8 controller: 127 value: 127",
		},
		{
			NoteOn(2, 130, 130),
			"NoteOn channel: 2 key: 127 velocity: 127",
		},
		{
			NoteOff(3, 180),
			"NoteOff channel: 3 key: 127",
		},
		{
			NoteOffVelocity(4, 180, 220),
			"NoteOff channel: 4 key: 127 velocity: 127",
		},
		{
			Pitchbend(4, 12300),
			"PitchBend channel: 4 pitch: 8191 (16383)",
		},
		{
			PolyAfterTouch(4, 186, 190),
			"PolyAfterTouch channel: 4 key: 127 pressure: 127",
		},
		{
			ProgramChange(4, 183),
			"ProgramChange channel: 4 program: 127",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		bf.WriteString(Message(test.input).String())

		if got, want := bf.String(), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}

func TestChannelRaw(t *testing.T) {

	tests := []struct {
		input    []byte
		expected string
	}{
		{ // 0
			AfterTouch(1, 120),
			"D1 78",
		},
		{ // 1
			ControlChange(8, 7, 110),
			"B8 07 6E",
		},
		{ // 2
			NoteOn(2, 100, 80),
			"92 64 50",
		},
		{ // 3
			NoteOff(3, 80),
			"83 50 00",
		},
		{
			NoteOffVelocity(4, 80, 20),
			"84 50 14",
		},
		{
			Pitchbend(4, 300),
			"E4 2C 42",
		},
		{
			PolyAfterTouch(4, 86, 109),
			"A4 56 6D",
		},
		{
			ProgramChange(4, 83),
			"C4 53",
		},
	}

	for i, test := range tests {

		var bf bytes.Buffer

		bf.Write(test.input)

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("[%v] got: %#v; wanted %#v", i, got, want)
		}
	}

}

func TestGetChannel(t *testing.T) {

	tests := []struct {
		input    []byte
		expected int
	}{
		{ // 0
			AfterTouch(1, 120),
			1,
		},
		{ // 1
			ControlChange(8, 7, 110),
			8,
		},
		{ // 2
			NoteOn(2, 100, 80),
			2,
		},
		{ // 3
			NoteOff(3, 80),
			3,
		},
		{
			NoteOffVelocity(4, 80, 20),
			4,
		},
		{
			Pitchbend(4, 300),
			4,
		},
		{
			PolyAfterTouch(4, 86, 109),
			4,
		},
		{
			ProgramChange(4, 83),
			4,
		},
	}

	for i, test := range tests {

		var ch uint8
		Message(test.input).GetChannel(&ch)

		if got, want := int(ch), test.expected; got != want {
			t.Errorf("[%v] got: %v; wanted %v", i, got, want)
		}
	}

}

/*
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
*/
