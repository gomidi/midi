package realtime_test

import (
	"bytes"
	"fmt"
	"testing"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midiwriter"
)

func mkInput(events ...midi.Message) []byte {
	var bf bytes.Buffer
	wr := midiwriter.New(&bf, midiwriter.NoRunningStatus())

	for _, ev := range events {
		wr.Write(ev)
	}

	return bf.Bytes()

}

func TestRead(t *testing.T) {

	tests := []struct {
		input    []byte
		times    int
		length   int
		output   string
		rtoutput string
	}{
		{
			mkInput(channel.Channel1.NoteOn(65, 100)),
			1,
			3,
			"91 41 64",
			"",
		},
		{
			mkInput(realtime.Start, channel.Channel1.NoteOn(65, 100)),
			1,
			3,
			"91 41 64",
			"Start\n",
		},
		{
			mkInput(channel.Channel1.NoteOn(65, 100), channel.Channel1.NoteOff(65)),
			2,
			3,
			"91 41 64 91 41 00",
			"",
		},
		{
			mkInput(channel.Channel1.NoteOn(65, 100), realtime.Continue, channel.Channel1.NoteOff(65)),
			2,
			3,
			"91 41 64 91 41 00",
			"Continue\n",
		},
		{
			mkInput(realtime.Start, channel.Channel1.NoteOn(65, 100), realtime.Stop, channel.Channel1.NoteOff(65)),
			2,
			3,
			"91 41 64 91 41 00",
			"Start\nStop\n",
		},
		{
			mkInput(channel.Channel1.ProgramChange(3), channel.Channel1.ProgramChange(4), channel.Channel1.ProgramChange(5)),
			3,
			2,
			"C1 03 C1 04 C1 05",
			"",
		},
		{
			mkInput(realtime.Start, channel.Channel1.ProgramChange(3), realtime.Stop, channel.Channel1.ProgramChange(4), realtime.Continue, channel.Channel1.ProgramChange(5)),
			3,
			2,
			"C1 03 C1 04 C1 05",
			"Start\nStop\nContinue\n",
		},
	}

	for n, test := range tests {
		var rtout bytes.Buffer
		var out bytes.Buffer
		h := func(ev realtime.Message) {
			rtout.WriteString(ev.String() + "\n")
		}

		rd := realtime.NewReader(bytes.NewReader(test.input), h)
		handler := "handler"
		var err error

		for x := 0; x < 2; x++ {

			for i := 0; i < test.times; i++ {
				bf := make([]byte, test.length)
				_, err = rd.Read(bf)

				if err != nil {
					break
				}

				out.Write(bf)
			}

			if err != nil {
				t.Errorf("[%v] Read(% X, %s) returned error: %v", n, test.input, handler, err)
			}

			if got, want := fmt.Sprintf("% X", out.Bytes()), test.output; got != want {
				t.Errorf("[%v] Read(% X, %s) = %#v (output); want %#v", n, test.input, handler, got, want)
			}

			if x == 0 {
				if got, want := rtout.String(), test.rtoutput; got != want {
					t.Errorf("[%v] Read(% X, %s) = %#v (rtoutput); want %#v", n, test.input, handler, got, want)
				}
				out.Reset()
				rd = realtime.NewReader(bytes.NewReader(test.input), nil)
				handler = "nil"
			}

		}

	}

}
