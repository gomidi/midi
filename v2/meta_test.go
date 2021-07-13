package midi

import (
	"bytes"
	"fmt"
	"testing"
)

/*
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
*/

func TestMessagesText(t *testing.T) {

	tests := []struct {
		input    Message
		expected string
	}{
		{
			MetaCopyright("(c) 2017"),
			"(c) 2017",
		},
		{
			MetaCuepoint("verse"),
			"verse",
		},
		{
			MetaDevice("2"),
			"2",
		},
		{
			MetaLyric("yeah"),
			"yeah",
		},
		{
			MetaMarker("TODO"),
			"TODO",
		},
		{
			MetaProgram("violin"),
			"violin",
		},
		{
			MetaTrackSequenceName("A"),
			"A",
		},
		{
			MetaText("hi"),
			"hi",
		},
		{
			MetaInstrument("1st violins"),
			"1st violins",
		},
	}

	for _, test := range tests {
		var got string

		test.input.text(&got)

		if want := test.expected; got != want {
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
			MetaCopyright("(c) 2017"),
			"FF 02 08 28 63 29 20 32 30 31 37",
		},
		{
			MetaCuepoint("verse"),
			"FF 07 05 76 65 72 73 65",
		},
		{
			MetaDevice("2"),
			"FF 09 01 32",
		},
		{
			EOT,
			"FF 2F 00",
		},
		{
			MetaKey(0, true, 0, false),
			"FF 59 02 00 00",
		},
		{
			MetaLyric("yeah"),
			"FF 05 04 79 65 61 68",
		},
		{
			MetaMarker("TODO"),
			"FF 06 04 54 4F 44 4F",
		},
		{
			MetaChannel(3),
			"FF 20 01 03",
		},
		{
			MetaPort(10),
			"FF 21 01 0A",
		},
		{
			MetaProgram("violin"),
			"FF 08 06 76 69 6F 6C 69 6E",
		},
		{
			MetaTrackSequenceName("A"),
			"FF 03 01 41",
		},
		{
			MetaSequenceNo(18),
			"FF 00 02 00 12",
		},
		{
			MetaSequencerData([]byte("hello world")),
			"FF 7F 0B 68 65 6C 6C 6F 20 77 6F 72 6C 64",
		},
		{
			MetaSMPTE(
				2,
				3,
				4,
				5,
				6,
			),
			"FF 54 05 02 03 04 05 06",
		},
		{
			MetaTempo(240),
			"FF 51 03 03 D0 90",
		},
		{
			MetaText("hi"),
			"FF 01 02 68 69",
		},
		{
			MetaTimeSig(3, 4, 8, 8),
			"FF 58 04 03 02 08 08",
		},
		{
			MetaInstrument("1st violins"),
			"FF 04 0B 31 73 74 20 76 69 6F 6C 69 6E 73",
		},
	}

	for _, test := range tests {

		var bf bytes.Buffer

		bf.Write(test.input.Data)

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}
