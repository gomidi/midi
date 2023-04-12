package midi

import (
	"bytes"
	"testing"
	//"gitlab.com/gomidi/midi/v2/sysex"
)

func TestMessageString(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			//midi.SysEx(sysex.GMSystem(2, true)),
			SysEx([]byte{0x7E, 0x02, 0x09, 0x01}),
			"SysExType data: 7E 02 09 01",
		},
		{
			SPP(4),
			"SPP spp: 4",
		},
		{
			SongSelect(2),
			"SongSelect song: 2",
		},
		{
			Tune(),
			"Tune",
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

/*
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
*/

/*
func TestMessagesRaw(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			SysEx([]byte("12,3")),
			"F0 31 32 2C 33 F7",
		},
		{
			Start([]byte("12,3")),
			"F0 31 32 2C 33",
		},
		{
			Continue([]byte("12,3")),
			"F7 31 32 2C 33",
		},
		{
			End([]byte("12,3")),
			"F7 31 32 2C 33 F7",
		},
		{
			Escape([]byte{0xFF, 0xF2}),
			"F7 FF F2",
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

func TestMessagesData(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			SysEx([]byte("12,3")),
			"31 32 2C 33",
		},
		{
			Start([]byte("12,3")),
			"31 32 2C 33",
		},
		{
			Continue([]byte("12,3")),
			"31 32 2C 33",
		},
		{
			End([]byte("12,3")),
			"31 32 2C 33",
		},
		{
			Escape([]byte{0xFF, 0xF2}),
			"FF F2",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		bf.Write(test.input.Data())

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}

func TestMessagesString(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			SysEx([]byte("12,3")),
			"sysex.SysEx len: 4",
		},
		{
			Start([]byte("12,3")),
			"sysex.Start len: 4",
		},
		{
			Continue([]byte("12,3")),
			"sysex.Continue len: 4",
		},
		{
			End([]byte("12,3")),
			"sysex.End len: 4",
		},
		{
			Escape([]byte{0xFF, 0xF2}),
			"sysex.Escape len: 2",
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
*/
