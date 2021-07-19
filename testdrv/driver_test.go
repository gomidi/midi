package testdrv

import (
	"bytes"
	"fmt"
	"testing"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/reader"
)

func TestDriver(t *testing.T) {
	var msgs = []midi.Message{
		channel.Channel0.NoteOn(45, 123),
		channel.Channel0.Aftertouch(45),
		channel.Channel0.NoteOff(45),
		channel.Channel0.PolyAftertouch(23, 56),
		channel.Channel0.Pitchbend(670),
	}

	d := New("testdrv")
	var bf bytes.Buffer
	rd := reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			fmt.Fprintf(&bf, "%s\n", msg.String())
		}),
	)

	in, _ := midi.OpenIn(d, 0, "")
	rd.ListenTo(in)
	out, _ := midi.OpenOut(d, 0, "")

	for _, msg := range msgs {
		out.Write(msg.Raw())
	}

	var expected = `channel.NoteOn channel 0 key 45 velocity 123
channel.Aftertouch channel 0 pressure 45
channel.NoteOff channel 0 key 45
channel.PolyAftertouch channel 0 key 23 pressure 56
channel.Pitchbend channel 0 value 670 absValue 8862
`

	var got = bf.String()

	if got != expected {
		t.Errorf("expected\n%s\n\ngot\n%s\n", expected, got)
	}

}
