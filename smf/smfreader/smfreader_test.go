package smfreader

import (
	"bytes"
	"fmt"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/internal/examples"
	"testing"
)

func testRead(t *testing.T, input []byte, options ...Option) string {
	var out bytes.Buffer
	out.WriteString("\n")
	rd := New(bytes.NewReader(input), options...)
	hd, err := rd.ReadHeader()

	if err != nil {
		t.Fatalf("can't read header: %v", err)
	}

	out.WriteString(fmt.Sprintf("SMF%v\n", hd.Format.Type()))
	out.WriteString(fmt.Sprintf("%v Track(s)\n", hd.NumTracks))
	out.WriteString(fmt.Sprintf("TimeFormat: %s\n", hd.TimeFormat))

	var _ = hd
	var msg midi.Message
	for {
		msg, err = rd.Read()

		if err != nil {
			break
		}

		out.WriteString(fmt.Sprintf("Track %v@%v %s\n", rd.Track(), rd.Delta(), msg))
	}

	return out.String()

}

func TestReadSMF0(t *testing.T) {
	var expected = `
SMF0
1 Track(s)
TimeFormat: 96 MetricTicks
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
`

	if got, want := testRead(t, examples.SpecSMF0), expected; got != want {
		t.Errorf("got:\n%v\n\nwanted\n%v\n\n", got, want)
	}

}

func TestReadSMF1(t *testing.T) {
	var expected = `
SMF1
4 Track(s)
TimeFormat: 96 MetricTicks
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
`

	if got, want := testRead(t, examples.SpecSMF1), expected; got != want {
		t.Errorf("got:\n%v\n\nwanted\n%v\n\n", got, want)
	}

}

func TestReadSMF1NoteOffPedantic(t *testing.T) {
	var expected = `
SMF1
4 Track(s)
TimeFormat: 96 MetricTicks
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
`

	if got, want := testRead(t, examples.SpecSMF1, ReadNoteOffPedantic()), expected; got != want {
		t.Errorf("got:\n%v\n\nwanted\n%v\n\n", got, want)
	}

}

func TestReadSMF0NoteOffPedantic(t *testing.T) {
	var expected = `
SMF0
1 Track(s)
TimeFormat: 96 MetricTicks
Track 0@0 meta.TimeSignature 4/4
Track 0@0 meta.Tempo BPM: 120
Track 0@0 channel.ProgramChange channel 0 program 5
Track 0@0 channel.ProgramChange channel 1 program 46
Track 0@0 channel.ProgramChange channel 2 program 70
Track 0@0 channel.NoteOn channel 2 pitch 48 vel 96
Track 0@0 channel.NoteOn channel 2 pitch 60 vel 96
Track 0@96 channel.NoteOn channel 1 pitch 67 vel 64
Track 0@96 channel.NoteOn channel 0 pitch 76 vel 32
Track 0@192 channel.NoteOffPedantic channel 2 pitch 48 velocity: 64
Track 0@0 channel.NoteOffPedantic channel 2 pitch 60 velocity: 64
Track 0@0 channel.NoteOffPedantic channel 1 pitch 67 velocity: 64
Track 0@0 channel.NoteOffPedantic channel 0 pitch 76 velocity: 64
Track 0@0 meta.endOfTrack
`

	if got, want := testRead(t, examples.SpecSMF0, ReadNoteOffPedantic()), expected; got != want {
		t.Errorf("got:\n%v\n\nwanted\n%v\n\n", got, want)
	}

}
