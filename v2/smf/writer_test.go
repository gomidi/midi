package smf

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"gitlab.com/gomidi/midi/v2"
)

// SpecSMF1Missing is like SpecSMF1 but with missing data for last track
var SpecSMF1Missing = []byte{
	// header chunk
	0x4D, 0x54, 0x68, 0x64, // MThd
	0x00, 0x00, 0x00, 0x06, // chunk length
	0x00, 0x01, // format 1
	0x00, 0x04, // four tracks
	0x00, 0x60, // 96 per quarter-note

	// first track (only time signature/tempo)

	// header of track 0
	0x4D, 0x54, 0x72, 0x6B, // MTrk
	0x00, 0x00, 0x00, 0x14, // chunk length (20)

	// data of track 0
	// delta          time event            comment
	0x00 /* delta */, 0xFF, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08, // time signature 4 bytes; 4/4 time; 24 MIDI clocks/click, 8 32nd notes/24 MIDI clocks
	0x00 /* delta */, 0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20, // tempo 120 BPM; 3 bytes: 500,000 usec/quarter note

	0x83, 0x00 /* delta */, 0xFF, 0x2F, 0x00, // end of track

	// second track

	// header of track 1
	0x4D, 0x54, 0x72, 0x6B, // MTrk
	0x00, 0x00, 0x00, 0x10, // chunk length (16)

	// data of track 1
	0x00 /* delta */, 0xC0, 0x05, // Ch.1 Program Change 5
	0x81, 0x40 /* delta */, 0x90, 0x4C, 0x20, // Ch.1 Note On E4, piano
	0x81, 0x40 /* delta */, 0x4C, 0x00, // Ch.1 Note On E4, velocity 0 (==noteoff) - running status
	0x00 /* delta */, 0xFF, 0x2F, 0x00, // end of track

	// third track

	// header of track 2
	0x4D, 0x54, 0x72, 0x6B, // MTrk
	0x00, 0x00, 0x00, 0x0F, // chunk length (15)

	// data of track 2
	0x00 /* delta */, 0xC1, 0x2E, // Ch.2 Program Change 46
	0x60 /* delta */, 0x91, 0x43, 0x40, // Ch.2 Note On G3, mezzo-forte
	0x82, 0x20 /* delta */, 0x43, 0x00, // Ch.2 Note On G3, velocity 0 (==noteoff) - running status
	0x00 /* delta */, 0xFF, 0x2F, 0x00, // end of track

}

// SpecSMF0 is an example from SMF spec for SMF type 0
var SpecSMF0 = []byte{
	// header chunk
	0x4D, 0x54, 0x68, 0x64, // MThd
	0x00, 0x00, 0x00, 0x06, // chunk length
	0x00, 0x00, // format 0
	0x00, 0x01, // one track
	0x00, 0x60, // 96 per quarter-note

	// first and only track

	// header of track
	0x4D, 0x54, 0x72, 0x6B, // MTrk
	0x00, 0x00, 0x00, 0x3B, // chunk length (59)

	// data of track
	// delta          time event            comment
	0x00 /* delta */, 0xFF, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08, // time signature 4 bytes; 4/4 time; 24 MIDI clocks/click, 8 32nd notes/24 MIDI clocks
	0x00 /* delta */, 0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20, // tempo 120 BPM; 3 bytes: 500,000 usec/quarter note
	0x00 /* delta */, 0xC0, 0x05, // Ch.1 Program Change 5
	0x00 /* delta */, 0xC1, 0x2E, // Ch.2 Program Change 46
	0x00 /* delta */, 0xC2, 0x46, // Ch.3 Program Change 70
	0x00 /* delta */, 0x92, 0x30, 0x60, // Ch.3 Note On C2, forte
	0x00 /* delta */, 0x3C, 0x60, // Ch.3 Note On C3, forte  - running status
	0x60 /* delta */, 0x91, 0x43, 0x40, // Ch.2 Note On G3, mezzo-forte
	0x60 /* delta */, 0x90, 0x4C, 0x20, // Ch.1 Note On E4, piano
	0x81, 0x40 /* delta */, 0x82, 0x30, 0x40, // two-byte delta-time; Ch.3 Note Off C2, standard
	0x00 /* delta */, 0x3C, 0x40, // Ch.3 Note Off C3, standard - running status
	0x00 /* delta */, 0x81, 0x43, 0x40, // Ch.2 Note Off G3, standard
	0x00 /* delta */, 0x80, 0x4C, 0x40, // Ch.1 Note Off E4, standard
	0x00 /* delta */, 0xFF, 0x2F, 0x00, // end of track

}

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

	var (
		bf  bytes.Buffer
		tr  Track
		res = MetricTicks(96)
	)

	tr.Add(0, MetaTimeSig(4, 4, 24, 8))
	tr.Add(0, MetaTempo(120))

	tr.Add(0, midi.ProgramChange(0, 5))
	tr.Add(0, midi.ProgramChange(1, 46))
	tr.Add(0, midi.ProgramChange(2, 70))

	tr.Add(0, midi.NoteOn(2, 48, 96))
	tr.Add(0, midi.NoteOn(2, 60, 96))

	tr.Add(res.Ticks4th(), midi.NoteOn(1, 67, 64))
	tr.Add(res.Ticks4th(), midi.NoteOn(0, 76, 32))

	tr.Add(res.Ticks4th()*2, midi.NoteOffVelocity(2, 48, 64))

	tr.Add(0, midi.NoteOffVelocity(2, 60, 64))
	tr.Add(0, midi.NoteOffVelocity(1, 67, 64))
	tr.Add(0, midi.NoteOffVelocity(0, 76, 64))

	tr.Close(0)

	smf := New()
	smf.TimeFormat = res
	smf.Add(tr)

	_, err := smf.WriteTo(&bf)

	if err != nil {
		t.Errorf("ERROR: %s", err.Error())
	}
	result := bf.Bytes()

	if got, want := result, SpecSMF0; !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n% X\n\nwanted:\n% X\n\n", got, want)
	}

}

// SpecSMF1 is an example from SMF spec for SMF type 1
var SpecSMF1 = []byte{
	// header chunk
	0x4D, 0x54, 0x68, 0x64, // MThd
	0x00, 0x00, 0x00, 0x06, // chunk length
	0x00, 0x01, // format 1
	0x00, 0x04, // four tracks
	0x00, 0x60, // 96 per quarter-note

	// first track (only time signature/tempo)

	// header of track 0
	0x4D, 0x54, 0x72, 0x6B, // MTrk
	0x00, 0x00, 0x00, 0x14, // chunk length (20)

	// data of track 0
	// delta          time event            comment
	0x00 /* delta */, 0xFF, 0x58, 0x04, 0x04, 0x02, 0x18, 0x08, // time signature 4 bytes; 4/4 time; 24 MIDI clocks/click, 8 32nd notes/24 MIDI clocks
	0x00 /* delta */, 0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20, // tempo 120 BPM; 3 bytes: 500,000 usec/quarter note

	0x83, 0x00 /* delta */, 0xFF, 0x2F, 0x00, // end of track

	// second track

	// header of track 1
	0x4D, 0x54, 0x72, 0x6B, // MTrk
	0x00, 0x00, 0x00, 0x10, // chunk length (16)

	// data of track 1
	0x00 /* delta */, 0xC0, 0x05, // Ch.1 Program Change 5
	0x81, 0x40 /* delta */, 0x90, 0x4C, 0x20, // Ch.1 Note On E4, piano
	0x81, 0x40 /* delta */, 0x4C, 0x00, // Ch.1 Note On E4, velocity 0 (==noteoff) - running status
	0x00 /* delta */, 0xFF, 0x2F, 0x00, // end of track

	// third track

	// header of track 2
	0x4D, 0x54, 0x72, 0x6B, // MTrk
	0x00, 0x00, 0x00, 0x0F, // chunk length (15)

	// data of track 2
	0x00 /* delta */, 0xC1, 0x2E, // Ch.2 Program Change 46
	0x60 /* delta */, 0x91, 0x43, 0x40, // Ch.2 Note On G3, mezzo-forte
	0x82, 0x20 /* delta */, 0x43, 0x00, // Ch.2 Note On G3, velocity 0 (==noteoff) - running status
	0x00 /* delta */, 0xFF, 0x2F, 0x00, // end of track

	// fourth track

	// header of track 3
	0x4D, 0x54, 0x72, 0x6B, // MTrk
	0x00, 0x00, 0x00, 0x15, // chunk length (21)

	// data of track 3
	0x00 /* delta */, 0xC2, 0x46, // Ch.3 Program Change 70
	0x00 /* delta */, 0x92, 0x30, 0x60, // Ch.3 Note On C2, forte
	0x00 /* delta */, 0x3C, 0x60, // Ch.3 Note On C3, forte  - running status
	0x83, 0x00 /* delta */, 0x30, 0x00, // two-byte delta-time; Ch.3 Note On C2, velocity 0 - running status
	0x00 /* delta */, 0x3C, 0x00, // Ch.3 Note Off C3, standard - running status
	0x00 /* delta */, 0xFF, 0x2F, 0x00, // end of track
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

	var (
		bf    bytes.Buffer
		track Track
		ticks = MetricTicks(96)
		smf   = New()
		beat  = ticks.Ticks4th()
	)

	smf.TimeFormat = ticks

	track = Track{}
	track.Add(0, MetaTimeSig(4, 4, 24, 8))
	track.Add(0, MetaTempo(120))
	track.Close(beat * 4)
	smf.Tracks = append(smf.Tracks, track)

	track = Track{}
	track.Add(0, midi.ProgramChange(0, 5))
	track.Add(beat*2, midi.NoteOn(0, 76, 32))
	track.Add(beat*2, midi.NoteOn(0, 76, 0))
	track.Close(0)
	smf.Tracks = append(smf.Tracks, track)

	track = Track{}
	track.Add(0, midi.ProgramChange(1, 46))
	track.Add(beat, midi.NoteOn(1, 67, 64))
	track.Add(beat*3, midi.NoteOn(1, 67, 0))
	track.Close(0)
	smf.Tracks = append(smf.Tracks, track)

	track = Track{}
	track.Add(0, midi.ProgramChange(2, 70))
	track.Add(0, midi.NoteOn(2, 48, 96))
	track.Add(0, midi.NoteOn(2, 60, 96))
	track.Add(beat*4, midi.NoteOn(2, 48, 0), midi.NoteOn(2, 60, 0))
	track.Close(0)
	smf.Tracks = append(smf.Tracks, track)

	_, err := smf.WriteTo(&bf)

	if err != nil {
		t.Errorf("ERROR: %s", err.Error())
	}

	if got, want := bf.Bytes(), SpecSMF1; !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n% X\n\nwanted:\n% X\n\n", got, want)
	}

}

type llog func(format string, vals ...interface{})

func (l llog) Printf(format string, vals ...interface{}) {
	l(format, vals...)
}

func TestWriteSysEx(t *testing.T) {
	var bf bytes.Buffer

	var smf = New()

	var tr Track
	tr.Add(0, midi.NoteOn(2, 65, 90))
	tr.Add(10, midi.SysEx([]byte{0x90, 0x51}))
	tr.Add(1, midi.NoteOff(2, 65))
	tr.Close(0)

	smf.Tracks = append(smf.Tracks, tr)

	_, err := smf.WriteTo(&bf)

	if err != nil {
		t.Errorf("Error while writing: %s\n", err.Error())
	}

	var lg = llog(func(format string, vals ...interface{}) {
		fmt.Printf(format, vals...)
	})

	_ = lg

	//rd, err := ReadFrom(&bf, Log(lg))
	rd, err := ReadFrom(&bf)

	if err != nil {
		t.Errorf("Error while reading: %s\n", err.Error())
	}

	trrd := rd.Tracks[0]

	var ch, key, velocity uint8

	var res bytes.Buffer
	res.WriteString("\n")

	for _, ev := range trrd {
		switch {
		case ev.Message.GetNoteOn(&ch, &key, &velocity):
			fmt.Fprintf(&res, "[%v] NoteOn at channel %v: key %v velocity: %v\n", ev.Delta, ch, key, velocity)
		case ev.Message.GetNoteOff(&ch, &key, &velocity):
			fmt.Fprintf(&res, "[%v] NoteOff at channel %v: key %v\n", ev.Delta, ch, key)
		default:
			if ev.Message.Is(midi.SysExMsg) {
				fmt.Fprintf(&res, "[%v] Sysex: % X\n", ev.Delta, ev.Message.Bytes())
			}
		}
	}

	expected := `
[0] NoteOn at channel 2: key 65 velocity: 90
[10] Sysex: F0 90 51 F7
[1] NoteOff at channel 2: key 65
`

	if got, want := res.String(), expected; got != want {
		t.Errorf("got\n%v\n\nwant\n%v\n\n", got, want)
	}

}

func TestRunningStatus(t *testing.T) {

	var bf bytes.Buffer

	var tr Track
	tr.Add(0, midi.NoteOn(0, 50, 33))
	tr.Add(2, midi.NoteOn(0, 50, 0)) // de facto a noteoff
	tr.Close(0)

	wr := New()
	wr.Tracks = []Track{tr}

	_, err := wr.WriteTo(&bf)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := "4D 54 68 64 00 00 00 06 00 00 00 01 03 C0 4D 54 72 6B 00 00 00 0B 00 90 32 21 02 32 00 00 FF 2F 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}

func TestNoRunningStatus(t *testing.T) {

	var bf bytes.Buffer

	var tr Track
	tr.Add(0, midi.NoteOn(0, 50, 33))
	tr.Add(2, midi.NoteOn(0, 50, 0)) // de facto a noteoff
	tr.Close(0)

	wr := New()
	wr.NoRunningStatus = true
	wr.Tracks = []Track{tr}

	_, err := wr.WriteTo(&bf)

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := "4D 54 68 64 00 00 00 06 00 00 00 01 03 C0 4D 54 72 6B 00 00 00 0C 00 90 32 21 02 90 32 00 00 FF 2F 00"

	if got, want := fmt.Sprintf("% X", bf.Bytes()), expected; got != want {
		t.Errorf("got:\n%#v\nwanted:\n%#v\n\n", got, want)
	}
}
