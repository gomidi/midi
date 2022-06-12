package smf

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
			MetaCopyright("(c) 2017"),
			"MetaCopyright text: \"(c) 2017\"",
		},
		{
			MetaCuepoint("verse"),
			"MetaCuepoint text: \"verse\"",
		},
		{
			MetaDevice("2"),
			"MetaDevice text: \"2\"",
		},
		{
			EOT,
			"MetaEndOfTrack",
		},
		{
			MetaKey(0, true, 0, false),
			"MetaKeySig key: CMaj",
		},
		{
			MetaLyric("yeah"),
			"MetaLyric text: \"yeah\"",
		},
		{
			MetaMarker("TODO"),
			"MetaMarker text: \"TODO\"",
		},
		{
			MetaChannel(3),
			"MetaChannel channel: 3",
		},
		{
			MetaPort(10),
			"MetaPort port: 10",
		},
		{
			MetaProgram("violin"),
			"MetaProgramName text: \"violin\"",
		},
		{
			MetaTrackSequenceName("A"),
			"MetaTrackName text: \"A\"",
		},
		{
			MetaSequenceNo(18),
			"MetaSeqNumber number: 18",
		},
		{
			MetaSequencerData([]byte("hello world")),
			"MetaSeqData bytes: 68 65 6C 6C 6F 20 77 6F 72 6C 64",
		},
		{
			MetaSMPTE(
				2, // hour
				3, // minute
				4, // second
				5, // frame
				6, // factional frame
			),
			"MetaSMPTEOffset hour: 2 minute: 3 second: 4 frame: 5 fractframe: 6",
		},
		{
			MetaTempo(240),
			"MetaTempo bpm: 240.00",
		},
		{
			MetaText("hi"),
			"MetaText text: \"hi\"",
		},
		{
			MetaTimeSig(3, 4, 8, 8),
			"MetaTimeSig meter: 3/4",
		},
		{
			MetaInstrument("1st violins"),
			"MetaInstrument text: \"1st violins\"",
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

		bf.Write(test.input.Bytes())

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("got: %#v; wanted %#v", got, want)
		}
	}

}
