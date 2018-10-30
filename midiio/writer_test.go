package midiio

import (
	"bytes"
	"fmt"
	"testing"

	"gitlab.com/gomidi/midi/midimessage/realtime"

	"gitlab.com/gomidi/midi/midiwriter"
)

func TestWriter(t *testing.T) {
	var bf bytes.Buffer

	wr := NewWriter(midiwriter.New(&bf))

	wr.Write([]byte{0x92, 0x41, 0x5A})
	wr.Write(realtime.Tick.Raw())
	wr.Write([]byte{0x41, 0x00})

	expected := "92 41 5A F9 41 00"

	if got, want := fmt.Sprintf("% X", bf.String()), expected; got != want {
		t.Errorf("Write() = %#v; want %#v", got, want)
	}

}
