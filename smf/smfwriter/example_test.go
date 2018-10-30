package smfwriter_test

import (
	"bytes"
	"fmt"

	// "io/ioutil"
	"gitlab.com/gomidi/midi"
	. "gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/midi/smf/smfreader"
	"gitlab.com/gomidi/midi/smf/smfwriter"
)

func Example() {
	fmt.Println()
	var bf bytes.Buffer

	wr := smfwriter.New(&bf)
	wr.Write(Channel2.Pitchbend(5000))
	wr.SetDelta(2)
	wr.Write(Channel2.NoteOn(65, 90))
	wr.SetDelta(4)
	wr.Write(Channel2.NoteOff(65))
	wr.Write(meta.EndOfTrack)

	rd := smfreader.New(bytes.NewReader(bf.Bytes()))

	var m midi.Message
	var err error

	for {
		m, err = rd.Read()

		// breaking at least with io.EOF
		if err != nil {
			break
		}

		// inspect
		fmt.Println(rd.Delta(), m)

		switch v := m.(type) {
		case NoteOn:
			fmt.Printf("NoteOn at channel %v: key %v velocity: %v\n", v.Channel(), v.Key(), v.Velocity())
		case NoteOff:
			fmt.Printf("NoteOff at channel %v: key %v\n", v.Channel(), v.Key())
		}

	}

	// Output:
	// 0 channel.Pitchbend channel 2 value 5000 absValue 13192
	// 2 channel.NoteOn channel 2 key 65 velocity 90
	// NoteOn at channel 2: key 65 velocity: 90
	// 4 channel.NoteOff channel 2 key 65
	// NoteOff at channel 2: key 65
	// 0 meta.EndOfTrack

}
