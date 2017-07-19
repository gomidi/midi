package meta

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
			Copyright("(c) 2017"),
			"meta.Copyright: \"(c) 2017\"",
		},
		{
			CuePoint("verse"),
			"meta.CuePoint: \"verse\"",
		},
		{
			DevicePort("2"),
			"meta.DevicePort: \"2\"",
		},
		{
			EndOfTrack,
			"meta.endOfTrack",
		},
		{
			KeySignature{Key: 0, IsMajor: true, Num: 0, IsFlat: false},
			"meta.KeySignature: C maj.",
		},
		{
			Lyric("yeah"),
			"meta.Lyric: \"yeah\"",
		},
		{
			Marker("TODO"),
			"meta.Marker: \"TODO\"",
		},
		{
			MIDIChannel(3),
			"meta.MIDIChannel: 3",
		},
		{
			MIDIPort(10),
			"meta.MIDIPort: 10",
		},
		{
			ProgramName("violin"),
			"meta.ProgramName: \"violin\"",
		},
		{
			Sequence("A"),
			"meta.Sequence: \"A\"",
		},
		{
			SequenceNumber(18),
			"meta.SequenceNumber: 18",
		},
		{
			SequencerSpecific([]byte("hello world")),
			"meta.SequencerSpecific len 11",
		},
		{
			SMPTEOffset{
				Hour:            2,
				Minute:          3,
				Second:          4,
				Frame:           5,
				FractionalFrame: 6,
			},
			"meta.SMPTEOffset 2:3:4 5.6",
		},
		{
			Tempo(240),
			"meta.Tempo BPM: 240",
		},
		{
			Text("hi"),
			"meta.Text: \"hi\"",
		},
		{
			TimeSignature{3, 4},
			"meta.TimeSignature 3/4",
		},
		{
			TrackInstrument("1st violins"),
			"meta.TrackInstrument: \"1st violins\"",
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
			Copyright("(c) 2017"),
			"FF 02 08 28 63 29 20 32 30 31 37",
		},
		{
			CuePoint("verse"),
			"FF 07 05 76 65 72 73 65",
		},
		{
			DevicePort("2"),
			"FF 09 01 32",
		},
		{
			EndOfTrack,
			"FF 2F 00",
		},
		{
			KeySignature{Key: 0, IsMajor: true, Num: 0, IsFlat: false},
			"FF 59 02 00 00",
		},
		{
			Lyric("yeah"),
			"FF 05 04 79 65 61 68",
		},
		{
			Marker("TODO"),
			"FF 06 04 54 4F 44 4F",
		},
		{
			MIDIChannel(3),
			"FF 20 01 03",
		},
		{
			MIDIPort(10),
			"FF 21 01 0A",
		},
		{
			ProgramName("violin"),
			"FF 08 06 76 69 6F 6C 69 6E",
		},
		{
			Sequence("A"),
			"FF 03 01 41",
		},
		{
			SequenceNumber(18),
			"FF 00 02 00 12",
		},
		{
			SequencerSpecific([]byte("hello world")),
			"FF 7F 0B 68 65 6C 6C 6F 20 77 6F 72 6C 64",
		},
		{
			SMPTEOffset{
				Hour:            2,
				Minute:          3,
				Second:          4,
				Frame:           5,
				FractionalFrame: 6,
			},
			"FF 54 05 02 03 04 05 06",
		},
		{
			Tempo(240),
			"FF 51 03 03 D0 90",
		},
		{
			Text("hi"),
			"FF 01 02 68 69",
		},
		{
			TimeSignature{3, 4},
			"FF 58 04 03 02 08 08",
		},
		{
			TrackInstrument("1st violins"),
			"FF 04 0B 31 73 74 20 76 69 6F 6C 69 6E 73",
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
