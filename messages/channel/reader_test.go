package channel_test

import (
	"bytes"
	// "fmt"
	// "io"
	"testing"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/midiwriter"
)

type readTest struct {
	input    *bytes.Buffer
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

	rd := bytes.NewBuffer(t.rawinput)

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
		mkTest(channel.Ch1.NoteOn(65, 100), "channel.NoteOn channel 1 pitch 65 vel 100"),
		mkTest(channel.Ch9.NoteOff(100), "channel.NoteOff channel 9 pitch 100"),
		mkTest(channel.Ch8.ProgramChange(3), "channel.ProgramChange channel 8 program 3"),
		mkTest(channel.Ch8.AfterTouch(30), "channel.AfterTouch channel 8 pressure 30"),
		mkTest(channel.Ch3.ControlChange(23, 25), "channel.ControlChange channel 3 controller 23 value 25"),
		mkTest(channel.Ch0.PitchWheel(123), "channel.PitchWheel channel 0 value 123 absValue 8315"),
		mkTest(channel.Ch15.PolyphonicAfterTouch(120, 106), "channel.PolyphonicAfterTouch channel 15 pitch 120 pressure 106"),
	}

	for n, test := range tests {
		var out bytes.Buffer

		ev, err := channel.NewReader(test.input, test.status).Read()

		if err != nil {
			t.Errorf("[%v] Read(% X) returned error: %v", n, test.rawinput, err)
			continue
		}
		out.WriteString(ev.String())

		if got, want := out.String(), test.expected; got != want {
			t.Errorf("[%v] Read(% X) = %#v; want %#v", n, test.rawinput, got, want)
		}

	}

}
