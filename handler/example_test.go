package handler_test

import (
	"bytes"
	"fmt"
	"github.com/gomidi/midi/handler"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf/smfwriter"
	"io"
)

func mkMIDI() io.Reader {
	var bf bytes.Buffer

	wr := smfwriter.New(&bf)
	wr.Write(channel.Ch2.NoteOn(65, 90))
	wr.SetDelta(2)
	wr.Write(channel.Ch2.NoteOff(65))
	wr.Write(meta.EndOfTrack)

	return bytes.NewReader(bf.Bytes())
}

func Example() {

	hd := handler.New(handler.NoLogger())

	// set the functions for the messages you are interested in
	hd.NoteOn = func(p *handler.Pos, channel, pitch, vel uint8) {
		fmt.Printf("[%v] NoteOn at channel %v: pitch %v velocity: %v\n", p.Delta, channel, pitch, vel)
	}

	hd.NoteOff = func(p *handler.Pos, channel, pitch uint8) {
		fmt.Printf("[%v] NoteOff at channel %v: pitch %v\n", p.Delta, channel, pitch)
	}

	hd.ReadSMF(mkMIDI())

	// Output: [0] NoteOn at channel 2: pitch 65 velocity: 90
	// [2] NoteOff at channel 2: pitch 65

}
