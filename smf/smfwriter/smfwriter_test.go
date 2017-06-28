package smfwriter

import (
	"bytes"
	"github.com/gomidi/midi/internal/examples"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
	"reflect"
	"testing"
)

/*
SMF0
1 Track(s)
TimeFormat: 96 QuarterNoteTicks
Track 0@0 meta.TimeSignature 4/4
Track 0@0 meta.Tempo BPM: 120
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

	resolution := smf.MetricResolution(96)

	wr := New(&bf, TimeFormat(resolution))
	wr.Write(meta.TimeSignatureDetailed{
		Numerator:                4,
		Denominator:              4,
		ClocksPerClick:           24,
		DemiSemiQuaverPerQuarter: 8,
	})
	wr.Write(meta.Tempo(120))
	wr.Write(channel.Ch0.ProgramChange(5))
	wr.Write(channel.Ch1.ProgramChange(46))
	wr.Write(channel.Ch2.ProgramChange(70))

	wr.Write(channel.Ch2.NoteOn(48, 96))
	wr.Write(channel.Ch2.NoteOn(60, 96))

	wr.SetDelta(resolution.N4())
	wr.Write(channel.Ch1.NoteOn(67, 64))

	wr.SetDelta(resolution.N4())
	wr.Write(channel.Ch0.NoteOn(76, 32))

	wr.SetDelta(resolution.N2())
	wr.Write(channel.Ch2.NoteOffPedantic(48, 64))

	wr.Write(channel.Ch2.NoteOffPedantic(60, 64))
	wr.Write(channel.Ch1.NoteOffPedantic(67, 64))
	wr.Write(channel.Ch0.NoteOffPedantic(76, 64))

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
Track 0@0 meta.Tempo BPM: 120
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

	resolution := smf.MetricResolution(96)

	wr := New(&bf, NumTracks(4), TimeFormat(resolution))
	wr.Write(meta.TimeSignatureDetailed{
		Numerator:                4,
		Denominator:              4,
		ClocksPerClick:           24,
		DemiSemiQuaverPerQuarter: 8,
	})
	wr.Write(meta.Tempo(120))
	wr.SetDelta(resolution.N4() * 4)
	wr.Write(meta.EndOfTrack)

	wr.Write(channel.Ch0.ProgramChange(5))
	wr.SetDelta(resolution.N2())
	wr.Write(channel.Ch0.NoteOn(76, 32))
	wr.SetDelta(resolution.N2())
	wr.Write(channel.Ch0.NoteOff(76))
	wr.Write(meta.EndOfTrack)

	wr.Write(channel.Ch1.ProgramChange(46))
	wr.SetDelta(resolution.N4())
	wr.Write(channel.Ch1.NoteOn(67, 64))
	wr.SetDelta(resolution.N4() * 3)
	wr.Write(channel.Ch1.NoteOff(67))
	wr.Write(meta.EndOfTrack)

	wr.Write(channel.Ch2.ProgramChange(70))
	wr.Write(channel.Ch2.NoteOn(48, 96))
	wr.Write(channel.Ch2.NoteOn(60, 96))
	wr.SetDelta(resolution.N4() * 4)
	wr.Write(channel.Ch2.NoteOff(48))
	wr.Write(channel.Ch2.NoteOff(60))
	wr.Write(meta.EndOfTrack)

	res := bf.Bytes()

	if got, want := res, examples.SpecSMF1; !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n% X\n\nwanted:\n% X\n\n", got, want)
	}

}
