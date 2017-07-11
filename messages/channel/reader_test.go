package channel_test

import (
	"bytes"

	"github.com/gomidi/midi/internal/midilib"
	// "fmt"
	"io"
	"testing"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/live/midiwriter"
	"github.com/gomidi/midi/messages/channel"
)

type readTest struct {
	input    io.Reader
	rawinput []byte
	status   byte
	expected string
}

func (r *readTest) func_name() {

}

func mkTest(event midi.Message, expected string) *readTest {
	var bf bytes.Buffer
	// we take no running status here, since the handling of running status
	// involves the runningstatus lib and midireader or smfreader
	wr := midiwriter.New(&bf, midiwriter.NoRunningStatus())
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
		mkTest(channel.Ch1.NoteOn(65, 100), "channel.NoteOn channel 1 key 65 velocity 100"),
		mkTest(channel.Ch9.NoteOff(100), "channel.NoteOff channel 9 key 100"),
		mkTest(channel.Ch9.NoteOffVelocity(120, 64), "channel.NoteOffVelocity channel 9 key 120 velocity 64"),
		mkTest(channel.Ch8.ProgramChange(3), "channel.ProgramChange channel 8 program 3"),
		mkTest(channel.Ch8.AfterTouch(30), "channel.AfterTouch (\"ChannelPressure\") channel 8 pressure 30"),
		mkTest(channel.Ch3.ControlChange(23, 25), "channel.ControlChange channel 3 controller 23 value 25"),
		mkTest(channel.Ch0.PitchBend(123), "channel.PitchBend (\"Portamento\") channel 0 value 123 absValue 8315"),
		mkTest(channel.Ch15.PolyphonicAfterTouch(120, 106), "channel.PolyphonicAfterTouch (\"KeyPressure\") channel 15 key 120 pressure 106"),
	}

	for n, test := range tests {
		var out bytes.Buffer

		// ignore running status (see above) and always read the first argument
		arg1, err := midilib.ReadByte(test.input)
		if err != nil {
			t.Errorf("[%v] ReadByte() returned error: %v", n, test.rawinput, err)
			continue
		}

		var m midi.Message

		m, err = channel.NewReader(test.input, channel.ReadNoteOffPedantic()).Read(test.status, arg1)

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
