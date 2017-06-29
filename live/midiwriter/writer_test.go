package midiwriter

import (
	"bytes"
	"fmt"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"testing"
)

func TestRunningStatus(t *testing.T) {

	var bf bytes.Buffer

	wr := New(&bf)

	wr.Write(channel.Ch0.NoteOn(50, 33))
	wr.Write(channel.Ch0.NoteOff(50))

	expected := "90 32 21 32 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}

func TestNoRunningStatus(t *testing.T) {

	var bf bytes.Buffer

	wr := New(&bf, NoRunningStatus())

	wr.Write(channel.Ch0.NoteOn(50, 33))
	wr.Write(channel.Ch0.NoteOff(50))

	expected := "90 32 21 90 32 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}

func TestSkipNonLiveMessages(t *testing.T) {

	var bf bytes.Buffer

	wr := New(&bf, SkipNonLiveMessages())

	wr.Write(channel.Ch0.NoteOn(50, 33))
	wr.Write(meta.Text("hi"))
	wr.Write(channel.Ch0.NoteOff(50))

	expected := "90 32 21 32 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}

func TestCheckMessageType(t *testing.T) {

	var bf bytes.Buffer

	wr := New(&bf, CheckMessageType())

	_, err1 := wr.Write(channel.Ch0.NoteOn(20, 20))
	if err1 != nil {
		t.Errorf("expected no error when writing channel message, got: %s", err1)
	}

	_, err2 := wr.Write(meta.Text("hi"))

	if err2 == nil {
		t.Errorf("expected error when writing meta message, got nil")
	}

}
