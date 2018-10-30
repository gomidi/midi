package syscommon_test

import (
	"bytes"
	// "fmt"

	// "github.com/gomidi/midi/internal/midilib"
	// "fmt"
	"io"
	"testing"

	"gitlab.com/gomidi/midi"
	. "gitlab.com/gomidi/midi/midimessage/syscommon"
	"gitlab.com/gomidi/midi/midiwriter"
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

	var bt = make([]byte, 1)

	_, err := rd.Read(bt)
	if err != nil {
		panic(err.Error())
	}

	t.input = rd
	t.status = bt[0]
	t.expected = expected
	return t
}

func TestRead(t *testing.T) {

	tests := []*readTest{
		mkTest(
			MTC(3),
			"syscommon.MTC: 3",
		),

		mkTest(
			SPP(4),
			"syscommon.SPP: 4",
		),

		mkTest(
			SongSelect(2),
			"syscommon.SongSelect: 2",
		),
		mkTest(
			Tune,
			"syscommon.Tune",
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
