package midiio

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midiwriter"
	// "github.com/gomidi/midi"
	"github.com/gomidi/midi/midireader"
)

func mkMIDI() []byte {
	var bf bytes.Buffer

	wr := midiwriter.New(&bf)
	wr.Write(channel.Ch2.NoteOn(65, 90))
	wr.Write(channel.Ch2.NoteOff(65))

	return bf.Bytes()
}

func TestReader(t *testing.T) {

	bt := mkMIDI()
	// fmt.Printf("% X\n", bt)

	mr := midireader.New(bytes.NewReader(bt), nil)

	rd := NewReader(mr)

	tests := []struct {
		expected string
	}{
		{"92 41 5A"},
		{"41 00 00"}, // running status
	}

	for i, test := range tests {
		var b = make([]byte, 3)
		_, err := rd.Read(b)
		if err != nil {
			t.Fatalf("Error: %s", err.Error())
		}

		if got, want := fmt.Sprintf("% X", b), test.expected; got != want {
			t.Errorf("Read()[%v] = %#v; want %#v", i, got, want)
		}
	}

}
