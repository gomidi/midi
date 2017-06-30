package smftrack_test

import (
	"bytes"
	"fmt"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf/smfreader"
	"github.com/gomidi/midi/smf/smftrack"
	"github.com/gomidi/midi/smf/smfwriter"
)

func mkSMF0() []byte {
	var bf bytes.Buffer

	wr := smfwriter.New(&bf)

	wr.Write(channel.Ch2.NoteOn(65, 90))
	wr.Write(channel.Ch1.NoteOn(24, 100))
	wr.SetDelta(2)
	wr.Write(channel.Ch1.NoteOff(24))
	wr.SetDelta(8)
	wr.Write(channel.Ch3.NoteOn(80, 109))
	wr.SetDelta(10)
	wr.Write(channel.Ch2.NoteOff(65))
	wr.SetDelta(10)
	wr.Write(channel.Ch3.NoteOff(80))
	wr.Write(meta.EndOfTrack)

	//	return bytes.NewReader(bf.Bytes())
	return bf.Bytes()
}

func ExampleToSMF1() {

	var bf bytes.Buffer

	// get some SMF0
	srcBytes := mkSMF0()

	src := smfreader.New(bytes.NewReader(srcBytes))
	src.ReadHeader()
	fmt.Println("source: ", src.Header().String())

	(smftrack.SMF0{}).ToSMF1(src, &bf)

	// read it back in
	rd := smfreader.New(bytes.NewReader(bf.Bytes()))
	rd.ReadHeader()
	fmt.Println("result: ", rd.Header().String())

	var m midi.Message
	var err error

	for {
		m, err = rd.Read()

		if err == smfreader.ErrFinished {
			break
		}

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			break
			// panic(err.Error())
		}

		switch v := m.(type) {
		case channel.NoteOn:
			fmt.Printf("Track %v [%v] NoteOn at channel %v: pitch %v velocity: %v\n", rd.Track(), rd.Delta(), v.Channel(), v.Pitch(), v.Velocity())
		case channel.NoteOff:
			fmt.Printf("Track %v [%v] NoteOff at channel %v: pitch %v\n", rd.Track(), rd.Delta(), v.Channel(), v.Pitch())
		}

	}

	// Output: source:  <Format: SMF0 (singletrack), NumTracks: 1, TimeFormat: 960 MetricTicks>
	// result:  <Format: SMF1 (multitrack), NumTracks: 4, TimeFormat: 960 MetricTicks>
	// Track 1 [0] NoteOn at channel 1: pitch 24 velocity: 100
	// Track 1 [2] NoteOff at channel 1: pitch 24
	// Track 2 [0] NoteOn at channel 2: pitch 65 velocity: 90
	// Track 2 [20] NoteOff at channel 2: pitch 65
	// Track 3 [10] NoteOn at channel 3: pitch 80 velocity: 109
	// Track 3 [20] NoteOff at channel 3: pitch 80

}
