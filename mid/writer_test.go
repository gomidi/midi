package mid

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"gitlab.com/gomidi/midi"
)

type captureLogger struct {
	bf bytes.Buffer
}

func (c *captureLogger) String() string {
	return c.bf.String()
}

func (c *captureLogger) Printf(format string, vals ...interface{}) {
	c.bf.WriteString(fmt.Sprintf(format, vals...))
}

func TestPlan(t *testing.T) {

	var bf bytes.Buffer

	wr := NewSMF(&bf, 1)

	wr.NoteOn(1, 120)
	wr.Plan(0, 4, 32, wr.Channel.NoteOff(1))
	wr.Forward(0, 8, 32)
	wr.NoteOn(2, 120)
	wr.Plan(0, 4, 32, wr.Channel.NoteOff(2))
	wr.Forward(1, 0, 0)
	wr.NoteOn(3, 120)
	wr.Plan(0, 4, 32, wr.Channel.NoteOff(3))
	wr.FinishPlanned()
	wr.EndOfTrack()

	var res captureLogger
	var expected = `
#0 [0 d:0] channel.NoteOn channel 0 key 1 velocity 120
#0 [480 d:480] channel.NoteOff channel 0 key 1
#0 [960 d:480] channel.NoteOn channel 0 key 2 velocity 120
#0 [1440 d:480] channel.NoteOff channel 0 key 2
#0 [3840 d:2400] channel.NoteOn channel 0 key 3 velocity 120
#0 [4320 d:480] channel.NoteOff channel 0 key 3
#0 [4320 d:0] meta.EndOfTrack	
`

	expected = strings.TrimSpace(expected)

	rd := NewReader(SetLogger(&res))
	//rd := NewReader()
	rd.Msg.Each = func(p *Position, msg midi.Message) {
		//		result = append(result, cc, val)
	}
	rd.ReadAllSMF(&bf)

	if got := strings.TrimSpace(res.String()); got != expected {
		t.Errorf("got\n%s\nexpected: \n\n%s\n", got, expected)
	}
}

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
		rd.ReadAllFrom(&bf)

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
