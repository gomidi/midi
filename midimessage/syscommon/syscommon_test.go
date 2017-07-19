package syscommon

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
			MIDITimingCode(3),
			"syscommon.MIDITimingCode: 3",
		},
		{
			SongPositionPointer(4),
			"syscommon.SongPositionPointer: 4",
		},
		{
			SongSelect(2),
			"syscommon.SongSelect: 2",
		},
		{
			TuneRequest,
			"syscommon.tuneRequest",
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
			MIDITimingCode(3),
			"F1 03",
		},
		{
			SongPositionPointer(4),
			"F2 00",
		},
		{
			SongSelect(2),
			"F3 02",
		},
		{
			TuneRequest,
			"F6",
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
