package meta_test

import (
	"bytes"

	// "github.com/gomidi/midi/internal/midilib"
	// "fmt"
	"io"
	"testing"

	"github.com/gomidi/midi"
	. "github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/midiwriter"
)

type readTest struct {
	input    io.Reader
	rawinput []byte
	status   byte
	expected string
}

func mkTest(event midi.Message, expected string) *readTest {
	var bf bytes.Buffer
	wr := midiwriter.New(&bf)
	wr.Write(event)

	t := &readTest{}
	t.rawinput = bf.Bytes()

	rd := bytes.NewReader(t.rawinput)

	var bt = make([]byte, 2)

	_, err := rd.Read(bt)
	if err != nil {
		panic(err.Error())
	}

	t.input = rd
	t.status = bt[1]
	t.expected = expected
	return t
}

func TestRead(t *testing.T) {

	tests := []*readTest{
		mkTest(
			Copyright("(c) 2017"),
			"meta.Copyright: \"(c) 2017\"",
		),
		mkTest(
			Cuepoint("verse"),
			"meta.Cuepoint: \"verse\"",
		),
		mkTest(
			Device("2"),
			"meta.Device: \"2\"",
		),
		mkTest(
			EndOfTrack,
			"meta.EndOfTrack",
		),
		mkTest(
			Key{Key: 0, IsMajor: true, Num: 0, IsFlat: false},
			"meta.Key: C maj.",
		),
		mkTest(
			Lyric("yeah"),
			"meta.Lyric: \"yeah\"",
		),
		mkTest(
			Marker("TODO"),
			"meta.Marker: \"TODO\"",
		),
		mkTest(
			Channel(3),
			"meta.Channel: 3",
		),
		mkTest(
			Port(10),
			"meta.Port: 10",
		),
		mkTest(
			Program("violin"),
			"meta.Program: \"violin\"",
		),
		mkTest(
			Sequence("A"),
			"meta.Sequence: \"A\"",
		),
		mkTest(
			SequenceNo(18),
			"meta.SequenceNo: 18",
		),
		mkTest(
			SequencerData([]byte("hello world")),
			"meta.SequencerData len 11",
		),
		mkTest(
			SMPTE{
				Hour:            2,
				Minute:          3,
				Second:          4,
				Frame:           5,
				FractionalFrame: 6,
			},
			"meta.SMPTE 2:3:4 5.6",
		),
		mkTest(
			Tempo(240),
			"meta.Tempo BPM: 240",
		),
		mkTest(
			Text("hi"),
			"meta.Text: \"hi\"",
		),
		mkTest(
			TimeSig{3, 4, 8, 8},
			"meta.TimeSig 3/4 clocksperclick 8 dsqpq 8",
		),
		mkTest(
			Track("1st violins"),
			"meta.Track: \"1st violins\"",
		),
	}

	for n, test := range tests {
		var out bytes.Buffer

		m, err := NewReader(test.input, test.status).Read()

		if err != nil {
			t.Errorf("[%v] Read(% X) returned error: %v", n, test.rawinput, err)
			continue
		}
		out.WriteString(m.String())

		if got, want := out.String(), test.expected; got != want {
			t.Errorf("[%v] Read(% X) = %#v; want %#v", n, test.rawinput, got, want)
		}

	}

}
