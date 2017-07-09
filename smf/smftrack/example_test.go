package smftrack_test

import (
	"bytes"
	"fmt"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfreader"
	"github.com/gomidi/midi/smf/smftrack"
	"github.com/gomidi/midi/smf/smfwriter"
)

func mkSMF1() []byte {
	var bf bytes.Buffer

	wr := smfwriter.New(&bf,
		smfwriter.NumTracks(2),
		smfwriter.Format(smf.SMF1), // not neccessary, since it is automatically set for numtracks > 1
	)

	wr.Write(channel.Ch2.NoteOn(65, 90))
	wr.SetDelta(20)
	wr.Write(channel.Ch2.NoteOff(65))
	wr.Write(meta.EndOfTrack)

	wr.Write(channel.Ch1.NoteOn(24, 100))
	wr.SetDelta(2)
	wr.Write(channel.Ch1.NoteOff(24))
	wr.Write(meta.EndOfTrack)

	//	return bytes.NewReader(bf.Bytes())
	return bf.Bytes()
}

func Example() {

	var bf bytes.Buffer

	// get some SMF1
	srcBytes := mkSMF1()

	src := smfreader.New(bytes.NewReader(srcBytes))
	src.ReadHeader()
	fmt.Println("source: ", src.Header().String())

	// create a new track
	tr := smftrack.New(2)
	tr.AddEvents(
		smftrack.NewEvent(10, channel.Ch3.NoteOn(80, 109)),
		smftrack.NewEvent(30, channel.Ch3.NoteOff(80)),
	)

	// add the track to the SMF1
	(smftrack.SMF1{}).AddTracks(src, &bf, tr)

	// read it back in
	src = smfreader.New(bytes.NewReader(bf.Bytes()))
	src.ReadHeader()
	fmt.Println("with added track: ", src.Header().String())

	// convert it to SMF0
	bf.Reset()
	(smftrack.SMF1{}).ToSMF0(src, &bf)

	// read back in to check
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
			panic(err.Error())
		}

		switch v := m.(type) {
		case channel.NoteOn:
			fmt.Printf("[%v] NoteOn at channel %v: key %v velocity %v\n", rd.Delta(), v.Channel(), v.Key(), v.Velocity())
		case channel.NoteOff:
			fmt.Printf("[%v] NoteOff at channel %v: key %v\n", rd.Delta(), v.Channel(), v.Key())
		}

	}

	// Output: source:  <Format: SMF1 (multitrack), NumTracks: 2, TimeFormat: 960 MetricTicks>
	// with added track:  <Format: SMF1 (multitrack), NumTracks: 3, TimeFormat: 960 MetricTicks>
	// result:  <Format: SMF0 (singletrack), NumTracks: 1, TimeFormat: 960 MetricTicks>
	// [0] NoteOn at channel 2: key 65 velocity 90
	// [0] NoteOn at channel 1: key 24 velocity 100
	// [2] NoteOff at channel 1: key 24
	// [8] NoteOn at channel 3: key 80 velocity 109
	// [10] NoteOff at channel 2: key 65
	// [10] NoteOff at channel 3: key 80

}
