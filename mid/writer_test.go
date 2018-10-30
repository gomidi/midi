package mid

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMsbLsb(t *testing.T) {

	tests := []struct {
		msb    uint8
		lsb    uint8
		value  uint16
		valMSB uint8
		valLSB uint8
	}{
		{msb: 22, lsb: 54, value: 16350, valMSB: 127, valLSB: 94},
		{msb: 22, lsb: 54, value: 0, valMSB: 0, valLSB: 0},
		{msb: 22, lsb: 54, value: 8192, valMSB: 64, valLSB: 0},
		{msb: 22, lsb: 54, value: 11419, valMSB: 89, valLSB: 27},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		wr := NewWriter(&bf)

		wr.MsbLsb(test.msb, test.lsb, test.value)

		var result []uint8

		rd := NewReader(NoLogger())
		rd.Msg.Channel.ControlChange.Each = func(p *Position, channel, cc, val uint8) {
			result = append(result, cc, val)
		}
		rd.Read(&bf)

		if len(result) != 4 {
			t.Errorf("len(result) must be 4, but is: %v", len(result))
		}

		if got, want := result[0:2], []uint8{test.msb, test.valMSB}; !reflect.DeepEqual(got, want) {
			t.Errorf("MSB(%v) = %v; want %v", test.value, got, want)
		}

		if got, want := result[2:4], []uint8{test.lsb, test.valLSB}; !reflect.DeepEqual(got, want) {
			t.Errorf("LSB(%v) = %v; want %v", test.value, got, want)
		}
	}

}
