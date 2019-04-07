package mid

import (
	"bytes"
	"testing"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/syscommon"
	"gitlab.com/gomidi/midi/midimessage/sysex"

	"reflect"

	mr "gitlab.com/gomidi/midi/midireader"
	mw "gitlab.com/gomidi/midi/midiwriter"
)

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
		var w = midiWriter{wr: mw.New(&bf)}
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
