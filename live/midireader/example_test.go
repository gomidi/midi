package midireader_test

import (
	"bytes"
	"fmt"
	"github.com/gomidi/midi/live/midireader"
	"github.com/gomidi/midi/live/midiwriter"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/realtime"

	"github.com/gomidi/midi"
)

func Example() {
	var bf bytes.Buffer

	wr := midiwriter.New(&bf)
	wr.Write(channel.Ch2.NoteOn(65, 90))
	wr.Write(realtime.Reset)
	wr.Write(channel.Ch2.NoteOff(65))

	rthandler := func(m realtime.Message) {
		fmt.Printf("Realtime: %s\n", m)
	}

	rd := midireader.New(bytes.NewReader(bf.Bytes()), rthandler)

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
			fmt.Printf("NoteOn at channel %v: pitch %v velocity: %v\n", v.Channel(), v.Pitch(), v.Velocity())
		case channel.NoteOff:
			fmt.Printf("NoteOff at channel %v: pitch %v\n", v.Channel(), v.Pitch())
		}

	}

	// Output: NoteOn at channel 2: pitch 65 velocity: 90
	// Realtime: Reset
	// NoteOff at channel 2: pitch 65
}
