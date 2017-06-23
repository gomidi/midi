package smfreader_test

import (
	"bytes"
	"fmt"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf/smfreader"
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

	rd := smfreader.New(mkMIDI())

	var m midi.Message
	var err error

	for {
		m, err = rd.Read()

		// breaking at least with io.EOF
		if err != nil {
			break
		}

		switch v := m.(type) {
		case channel.NoteOn:
			fmt.Printf("[%v] NoteOn at channel %v: pitch %v velocity: %v\n", rd.Delta(), v.Channel(), v.Pitch(), v.Velocity())
		case channel.NoteOff:
			fmt.Printf("[%v] NoteOff at channel %v: pitch %v\n", rd.Delta(), v.Channel(), v.Pitch())
		}

	}

	// Output: [0] NoteOn at channel 2: pitch 65 velocity: 90
	// [2] NoteOff at channel 2: pitch 65

}
