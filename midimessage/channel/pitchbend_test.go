package channel_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"gitlab.com/gomidi/midi/internal/midilib"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/midiwriter"
)

/*
{msb: 22, lsb: 54, value: 16350, valMSB: 127, valLSB: 94},
		{msb: 22, lsb: 54, value: 0, valMSB: 0, valLSB: 0},
		{msb: 22, lsb: 54, value: 8192, valMSB: 64, valLSB: 0},
		{msb: 22, lsb: 54, value: 11419, valMSB: 89, valLSB: 27},
*/

func TestMsbLsb(t *testing.T) {
	tests := []struct {
		unsigned uint16
		signed   int16
		msb      uint8
		lsb      uint8
	}{
		{16350, 16350, 127, 94},
	}

	for i, test := range tests {

		_, absValue := midilib.ParsePitchWheelVals(test.lsb, test.msb)

		if absValue != test.unsigned {
			t.Errorf("[%v] unsigned = %v; wanted %v", i, absValue, test.unsigned)
		}

		r := midilib.MsbLsbUnsigned(test.unsigned)
		var b = make([]byte, 2)
		binary.BigEndian.PutUint16(b, r)
		//		_ = b[0], b[1]

		if b[1] != test.msb {
			t.Errorf("[%v] msb = %v; wanted %v", i, b[1], test.msb)
		}

		if b[0] != test.lsb {
			t.Errorf("[%v] lsb = %v; wanted %v", i, b[0], test.lsb)
		}

	}

}

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
