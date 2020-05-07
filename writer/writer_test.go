package writer

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/syscommon"
	"gitlab.com/gomidi/midi/midimessage/sysex"
	"gitlab.com/gomidi/midi/reader"

	mr "gitlab.com/gomidi/midi/midireader"
	mw "gitlab.com/gomidi/midi/midiwriter"
	//"gitlab.com/gomidi/midi/writer"
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

	wr := NewSMF(&bf, 2)

	Meter(wr, 4, 4)
	Forward(wr, 0, 8, 4)
	Meter(wr, 3, 4)
	EndOfTrack(wr)

	// 1
	NoteOn(wr, 1, 120)
	// 1&
	Plan(wr, 0, 4, 32, Channel(wr).NoteOff(1))
	// 2
	Forward(wr, 0, 8, 32)
	NoteOn(wr, 2, 120)
	// 2&
	Plan(wr, 0, 4, 32, Channel(wr).NoteOff(2))

	// 1
	Forward(wr, 1, 0, 0)
	NoteOn(wr, 3, 120)

	// 1&
	Plan(wr, 0, 4, 32, Channel(wr).NoteOff(3))

	// 2
	Forward(wr, 1, 8, 32)
	NoteOn(wr, 4, 120)
	// 2&
	Plan(wr, 0, 4, 32, Channel(wr).NoteOff(4))

	FinishPlanned(wr)
	EndOfTrack(wr)

	var res captureLogger
	var expected = `
#0 [0 d:0] meta.TimeSig 4/4 clocksperclick 8 dsqpq 8
#0 [7680 d:7680] meta.TimeSig 3/4 clocksperclick 8 dsqpq 8
#0 [7680 d:0] meta.EndOfTrack
#1 [0 d:0] channel.NoteOn channel 0 key 1 velocity 120
#1 [480 d:480] channel.NoteOff channel 0 key 1
#1 [960 d:480] channel.NoteOn channel 0 key 2 velocity 120
#1 [1440 d:480] channel.NoteOff channel 0 key 2
#1 [3840 d:2400] channel.NoteOn channel 0 key 3 velocity 120
#1 [4320 d:480] channel.NoteOff channel 0 key 3
#1 [8640 d:4320] channel.NoteOn channel 0 key 4 velocity 120
#1 [9120 d:480] channel.NoteOff channel 0 key 4
#1 [9120 d:0] meta.EndOfTrack	
`

	expected = strings.TrimSpace(expected)

	rd := reader.New(reader.SetLogger(&res), reader.Each(func(p *reader.Position, msg midi.Message) {
		//		result = append(result, cc, val)
	}))
	//rd := NewReader()
	reader.ReadSMF(rd, &bf)

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

		wr := New(&bf)

		MsbLsb(wr, test.msb, test.lsb, test.value)

		var result []uint8

		rd := reader.New(reader.NoLogger(), reader.ControlChange(func(p *reader.Position, channel, cc, val uint8) {
			result = append(result, cc, val)
		}))
		reader.ReadAllFrom(rd, &bf)

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

func m(msgs ...midi.Message) []midi.Message {
	return msgs
}

func TestConsolidation(t *testing.T) {

	tests := []struct {
		input         []midi.Message
		expected      []midi.Message
		consolidation bool
	}{
		{nil, nil, true},
		{nil, nil, false},
		{
			m(channel.Channel1.NoteOn(127, 34)),
			m(channel.Channel1.NoteOn(127, 34)),
			true,
		},
		{
			m(channel.Channel1.NoteOff(127)),
			nil,
			true,
		},
		{
			m(channel.Channel1.NoteOff(127)),
			m(channel.Channel1.NoteOff(127)),
			false,
		},
		{
			m(channel.Channel1.NoteOn(127, 34), channel.Channel1.NoteOn(127, 36)),
			m(channel.Channel1.NoteOn(127, 34)),
			true,
		},
		{
			m(channel.Channel1.NoteOn(127, 34), channel.Channel1.NoteOff(127), channel.Channel1.NoteOn(127, 36)),
			m(channel.Channel1.NoteOn(127, 34), channel.Channel1.NoteOff(127), channel.Channel1.NoteOn(127, 36)),
			true,
		},
		{
			m(channel.Channel1.NoteOn(127, 34), channel.Channel1.NoteOn(120, 36)),
			m(channel.Channel1.NoteOn(127, 34), channel.Channel1.NoteOn(120, 36)),
			true,
		},
		{
			m(channel.Channel1.NoteOn(127, 34), syscommon.Tune, channel.Channel1.NoteOn(127, 36), channel.Channel1.NoteOff(127)),
			m(channel.Channel1.NoteOn(127, 34), syscommon.Tune, channel.Channel1.NoteOff(127)),
			true,
		},
		{
			m(channel.Channel1.NoteOn(127, 34), sysex.SysEx([]byte{45}), channel.Channel1.NoteOn(127, 36), channel.Channel1.NoteOff(127)),
			m(channel.Channel1.NoteOn(127, 34), sysex.SysEx([]byte{45}), channel.Channel1.NoteOff(127)),
			true,
		},
	}

	for _, test := range tests {
		var bf bytes.Buffer
		var w = Writer{wr: mw.New(&bf)}
		w.ConsolidateNotes(test.consolidation)

		for _, msg := range test.input {
			w.Write(msg)
		}

		r := mr.New(bytes.NewReader(bf.Bytes()), nil)

		var result []midi.Message
		for {
			m, err := r.Read()
			if m != nil {
				result = append(result, m)
			}
			if err != nil {
				break
			}
		}

		if got, want := result, test.expected; !reflect.DeepEqual(got, want) {
			t.Errorf("ConsolidateNotes(%v).Write(%v) = %v; want %v", test.consolidation, test.input, got, want)
		}
	}

}
