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
	"io"
)

func mkSMF1() io.Reader {
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

	return bytes.NewReader(bf.Bytes())
}

func Example() {

	var bf bytes.Buffer

	// get some SMF1
	src := smfreader.New(mkSMF1())
	src.ReadHeader()
	fmt.Println("source: ", src.Header().String())

	// convert it to SMF0
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
			fmt.Printf("[%v] NoteOn at channel %v: pitch %v velocity: %v\n", rd.Delta(), v.Channel(), v.Pitch(), v.Velocity())
		case channel.NoteOff:
			fmt.Printf("[%v] NoteOff at channel %v: pitch %v\n", rd.Delta(), v.Channel(), v.Pitch())
		}

	}

	// Output: source:  <Format: SMF1 (multitrack), NumTracks: 2, TimeFormat: 960 MetricTicks>
	// result:  <Format: SMF0 (singletrack), NumTracks: 1, TimeFormat: 960 MetricTicks>
	// [0] NoteOn at channel 2: pitch 65 velocity: 90
	// [0] NoteOn at channel 1: pitch 24 velocity: 100
	// [2] NoteOff at channel 1: pitch 24
	// [18] NoteOff at channel 2: pitch 65

}
