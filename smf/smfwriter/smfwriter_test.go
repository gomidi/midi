package smfwriter

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/internal/examples"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/midimessage/sysex"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfreader"
	// "log"
	// "os"
)

/*
SMF0
1 Track(s)
TimeFormat: 96 QuarterNoteTicks
Track 0@0 meta.TimeSignature 4/4
Track 0@0 meta.Tempo BPM: 120.00
Track 0@0 channel.ProgramChange channel 0 program 5
Track 0@0 channel.ProgramChange channel 1 program 46
Track 0@0 channel.ProgramChange channel 2 program 70
Track 0@0 channel.NoteOn channel 2 pitch 48 vel 96
Track 0@0 channel.NoteOn channel 2 pitch 60 vel 96
Track 0@96 channel.NoteOn channel 1 pitch 67 vel 64
Track 0@96 channel.NoteOn channel 0 pitch 76 vel 32
Track 0@192 channel.NoteOff channel 2 pitch 48
Track 0@0 channel.NoteOff channel 2 pitch 60
Track 0@0 channel.NoteOff channel 1 pitch 67
Track 0@0 channel.NoteOff channel 0 pitch 76
Track 0@0 meta.endOfTrack
*/

func TestWriteSMF0(t *testing.T) {
	var bf bytes.Buffer

	resolution := smf.MetricTicks(96)

	wr := New(&bf, TimeFormat(resolution), Format(smf.SMF0))
	wr.Write(meta.TimeSig{
		Numerator:                4,
		Denominator:              4,
		ClocksPerClick:           24,
		DemiSemiQuaverPerQuarter: 8,
	})
	wr.Write(meta.BPM(120))
	wr.Write(channel.Channel0.ProgramChange(5))
	wr.Write(channel.Channel1.ProgramChange(46))
	wr.Write(channel.Channel2.ProgramChange(70))

	wr.Write(channel.Channel2.NoteOn(48, 96))
	wr.Write(channel.Channel2.NoteOn(60, 96))

	wr.SetDelta(resolution.Ticks4th())
	wr.Write(channel.Channel1.NoteOn(67, 64))

	wr.SetDelta(resolution.Ticks4th())
	wr.Write(channel.Channel0.NoteOn(76, 32))

	wr.SetDelta(resolution.Ticks4th() * 2)
	wr.Write(channel.Channel2.NoteOffVelocity(48, 64))

	wr.Write(channel.Channel2.NoteOffVelocity(60, 64))
	wr.Write(channel.Channel1.NoteOffVelocity(67, 64))
	wr.Write(channel.Channel0.NoteOffVelocity(76, 64))

	wr.Write(meta.EndOfTrack)

	res := bf.Bytes()

	if got, want := res, examples.SpecSMF0; !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n% X\n\nwanted:\n% X\n\n", got, want)
	}

}

/*
SMF1
4 Track(s)
TimeFormat: 96 QuarterNoteTicks
Track 0@0 meta.TimeSignature 4/4
Track 0@0 meta.Tempo BPM: 120.00
Track 0@384 meta.endOfTrack
Track 1@0 channel.ProgramChange channel 0 program 5
Track 1@192 channel.NoteOn channel 0 pitch 76 vel 32
Track 1@192 channel.NoteOff channel 0 pitch 76
Track 1@0 meta.endOfTrack
Track 2@0 channel.ProgramChange channel 1 program 46
Track 2@96 channel.NoteOn channel 1 pitch 67 vel 64
Track 2@288 channel.NoteOff channel 1 pitch 67
Track 2@0 meta.endOfTrack
Track 3@0 channel.ProgramChange channel 2 program 70
Track 3@0 channel.NoteOn channel 2 pitch 48 vel 96
Track 3@0 channel.NoteOn channel 2 pitch 60 vel 96
Track 3@384 channel.NoteOff channel 2 pitch 48
Track 3@0 channel.NoteOff channel 2 pitch 60
Track 3@0 meta.endOfTrack
*/

func TestWriteSMF1(t *testing.T) {
	var bf bytes.Buffer

	resolution := smf.MetricTicks(96)

	wr := New(&bf, NumTracks(4), TimeFormat(resolution))
	wr.Write(meta.TimeSig{
		Numerator:                4,
		Denominator:              4,
		ClocksPerClick:           24,
		DemiSemiQuaverPerQuarter: 8,
	})
	wr.Write(meta.BPM(120))

	wr.SetDelta(resolution.Ticks4th() * 4)
	wr.Write(meta.EndOfTrack)

	wr.Write(channel.Channel0.ProgramChange(5))

	wr.SetDelta(resolution.Ticks4th() * 2)
	wr.Write(channel.Channel0.NoteOn(76, 32))

	wr.SetDelta(resolution.Ticks4th() * 2)
	wr.Write(channel.Channel0.NoteOff(76))

	wr.Write(meta.EndOfTrack)

	wr.Write(channel.Channel1.ProgramChange(46))

	wr.SetDelta(resolution.Ticks4th())
	wr.Write(channel.Channel1.NoteOn(67, 64))

	wr.SetDelta(resolution.Ticks4th() * 3)
	wr.Write(channel.Channel1.NoteOff(67))

	wr.Write(meta.EndOfTrack)

	wr.Write(channel.Channel2.ProgramChange(70))

	wr.Write(channel.Channel2.NoteOn(48, 96))
	wr.Write(channel.Channel2.NoteOn(60, 96))

	wr.SetDelta(resolution.Ticks4th() * 4)
	wr.Write(channel.Channel2.NoteOff(48))
	wr.Write(channel.Channel2.NoteOff(60))

	wr.Write(meta.EndOfTrack)

	res := bf.Bytes()

	if got, want := res, examples.SpecSMF1; !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n% X\n\nwanted:\n% X\n\n", got, want)
	}

}

func TestWriteSysEx(t *testing.T) {
	var bf bytes.Buffer

	wr := New(&bf)
	wr.SetDelta(0)
	wr.Write(channel.Channel2.NoteOn(65, 90))
	wr.SetDelta(10)
	wr.Write(sysex.SysEx([]byte{0x90, 0x51}))
	wr.SetDelta(1)
	wr.Write(channel.Channel2.NoteOff(65))
	wr.Write(meta.EndOfTrack)

	rd := smfreader.New(bytes.NewReader(bf.Bytes()))

	var m midi.Message
	var err error

	var res bytes.Buffer
	res.WriteString("\n")
	for {
		m, err = rd.Read()

		// breaking at least with io.EOF
		if err != nil {
			break
		}

		switch v := m.(type) {
		case sysex.SysEx:
			fmt.Fprintf(&res, "[%v] Sysex: % X\n", rd.Delta(), v.Data())
		case channel.NoteOn:
			fmt.Fprintf(&res, "[%v] NoteOn at channel %v: key %v velocity: %v\n", rd.Delta(), v.Channel(), v.Key(), v.Velocity())
		case channel.NoteOff:
			fmt.Fprintf(&res, "[%v] NoteOff at channel %v: key %v\n", rd.Delta(), v.Channel(), v.Key())
		}

	}

	expected := `
[0] NoteOn at channel 2: key 65 velocity: 90
[10] Sysex: 90 51
[1] NoteOff at channel 2: key 65
`

	if got, want := res.String(), expected; got != want {
		t.Errorf("got\n%v\n\nwant\n%v\n\n", got, want)
	}

}

func TestRunningStatus(t *testing.T) {

	var bf bytes.Buffer

	wr := New(&bf)

	err := wr.WriteHeader()

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	wr.Write(channel.Channel0.NoteOn(50, 33))
	wr.SetDelta(2)
	wr.Write(channel.Channel0.NoteOff(50))
	wr.Write(meta.EndOfTrack)

	expected := "4D 54 68 64 00 00 00 06 00 00 00 01 03 C0 4D 54 72 6B 00 00 00 0B 00 90 32 21 02 32 00 00 FF 2F 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}

func TestNoRunningStatus(t *testing.T) {

	var bf bytes.Buffer

	wr := New(&bf, NoRunningStatus())

	err := wr.WriteHeader()

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	wr.Write(channel.Channel0.NoteOn(50, 33))
	wr.SetDelta(2)
	wr.Write(channel.Channel0.NoteOff(50))
	wr.Write(meta.EndOfTrack)

	expected := "4D 54 68 64 00 00 00 06 00 00 00 01 03 C0 4D 54 72 6B 00 00 00 0C 00 90 32 21 02 90 32 00 00 FF 2F 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}
