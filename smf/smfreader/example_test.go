package smfreader_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gomidi/midi"
	. "github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/smf/smfreader"
	"github.com/gomidi/midi/smf/smfwriter"
)

func mkMIDI() io.Reader {
	var bf bytes.Buffer

	wr := smfwriter.New(&bf)
	wr.Write(Channel2.NoteOn(65, 90))
	wr.SetDelta(2)
	wr.Write(Channel2.NoteOff(65))
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
		case NoteOn:
			fmt.Printf("[%v] NoteOn at channel %v: key %v velocity %v\n", rd.Delta(), v.Channel(), v.Key(), v.Velocity())
		case NoteOff:
			fmt.Printf("[%v] NoteOff at channel %v: key %v\n", rd.Delta(), v.Channel(), v.Key())
		}

	}

	// Output: [0] NoteOn at channel 2: key 65 velocity 90
	// [2] NoteOff at channel 2: key 65

}
