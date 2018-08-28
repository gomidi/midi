package smfreader_test

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gomidi/midi"
	. "github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfreader"
	"github.com/gomidi/midi/smf/smfwriter"
)

func mkMIDI() io.Reader {
	var bf bytes.Buffer

	wr := smfwriter.New(&bf)
	wr.Write(Channel2.Pitchbend(5000))
	wr.Write(Channel2.NoteOn(65, 90))
	wr.SetDelta(2)
	wr.Write(Channel2.NoteOff(65))
	wr.SetDelta(4)
	wr.Write(meta.EndOfTrack)
	return bytes.NewReader(bf.Bytes())
}

func Example() {
	fmt.Println()

	rd := smfreader.New(mkMIDI())

	var m midi.Message
	var err error

	for {
		m, err = rd.Read()

		// at the end, smf.ErrFinished will be returned
		if err != nil {
			break
		}

		// inspect
		fmt.Println(rd.Delta(), m)

		switch v := m.(type) {
		case NoteOn:
			fmt.Printf("[%v] NoteOn at channel %v: key %v velocity %v\n", rd.Delta(), v.Channel(), v.Key(), v.Velocity())
		case NoteOff:
			fmt.Printf("[%v] NoteOff at channel %v: key %v\n", rd.Delta(), v.Channel(), v.Key())
		}

	}

	if err != smf.ErrFinished {
		panic("error: " + err.Error())
	}

	// Output:
	// 0 channel.Pitchbend channel 2 value 5000 absValue 13192
	// 0 channel.NoteOn channel 2 key 65 velocity 90
	// [0] NoteOn at channel 2: key 65 velocity 90
	// 2 channel.NoteOff channel 2 key 65
	// [2] NoteOff at channel 2: key 65
	// 4 meta.EndOfTrack

}
