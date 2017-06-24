package midireader

import (
	"bytes"
	"github.com/gomidi/midi/live/midiwriter"
	"github.com/gomidi/midi/messages/channel"
	"io"
	"testing"
)

func mkMIDI() io.Reader {
	var bf bytes.Buffer

	wr := midiwriter.New(&bf)

	wr.Write(channel.New(1).NoteOn(65, 100))
	wr.Write(channel.New(1).NoteOff(65))
	return bytes.NewReader(bf.Bytes())
}

func TestRead(t *testing.T) {

	rd := New(mkMIDI(), nil)

	ev, err := rd.Read()

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := "channel.NoteOn channel 1 pitch 65 vel 100"
	if ev.String() != expected {
		t.Errorf("expected: %#v, got: %#v", expected, ev.String())
	}

}
