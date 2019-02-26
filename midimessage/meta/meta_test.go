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
			Cuepoint("verse"),
			"meta.Cuepoint: \"verse\"",
		},
		{
			Device("2"),
			"meta.Device: \"2\"",
		},
		{
			EndOfTrack,
			"meta.EndOfTrack",
		},
		{
			Key{Key: 0, IsMajor: true, Num: 0, IsFlat: false},
			"meta.Key: C maj.",
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
			Channel(3),
			"meta.Channel: 3",
		},
		{
			Port(10),
			"meta.Port: 10",
		},
		{
			Program("violin"),
			"meta.Program: \"violin\"",
		},
		{
			TrackSequenceName("A"),
			"meta.TrackSequenceName: \"A\"",
		},
		{
			SequenceNo(18),
			"meta.SequenceNo: 18",
		},
		{
			SequencerData([]byte("hello world")),
			"meta.SequencerData len 11",
		},
		{
			SMPTE{
				Hour:            2,
				Minute:          3,
				Second:          4,
				Frame:           5,
				FractionalFrame: 6,
			},
			"meta.SMPTE 2:3:4 5.6",
		},
		{
			BPM(240),
			"meta.Tempo BPM: 240.00",
		},
		{
			Text("hi"),
			"meta.Text: \"hi\"",
		},
		{
			TimeSig{3, 4, 8, 8},
			"meta.TimeSig 3/4 clocksperclick 8 dsqpq 8",
		},
		{
			Instrument("1st violins"),
			"meta.Instrument: \"1st violins\"",
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
			Cuepoint("verse"),
			"FF 07 05 76 65 72 73 65",
		},
		{
			Device("2"),
			"FF 09 01 32",
		},
		{
			EndOfTrack,
			"FF 2F 00",
		},
		{
			Key{Key: 0, IsMajor: true, Num: 0, IsFlat: false},
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
			Channel(3),
			"FF 20 01 03",
		},
		{
			Port(10),
			"FF 21 01 0A",
		},
		{
			Program("violin"),
			"FF 08 06 76 69 6F 6C 69 6E",
		},
		{
			TrackSequenceName("A"),
			"FF 03 01 41",
		},
		{
			SequenceNo(18),
			"FF 00 02 00 12",
		},
		{
			SequencerData([]byte("hello world")),
			"FF 7F 0B 68 65 6C 6C 6F 20 77 6F 72 6C 64",
		},
		{
			SMPTE{
				Hour:            2,
				Minute:          3,
				Second:          4,
				Frame:           5,
				FractionalFrame: 6,
			},
			"FF 54 05 02 03 04 05 06",
		},
		{
			BPM(240),
			"FF 51 03 03 D0 90",
		},
		{
			Text("hi"),
			"FF 01 02 68 69",
		},
		{
			TimeSig{3, 4, 8, 8},
			"FF 58 04 03 02 08 08",
		},
		{
			Instrument("1st violins"),
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
