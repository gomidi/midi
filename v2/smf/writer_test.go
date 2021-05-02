package smf

import (
	"bytes"
	"reflect"
	"testing"

	. "gitlab.com/gomidi/midi/v2"
)

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
		bf    bytes.Buffer
		track = NewTrack()
		res   = MetricTicks(96)
		ch0   = Channel(0)
		ch1   = Channel(1)
		ch2   = Channel(2)
	)

	track.Add(0, MetaTimeSig(4, 4, 24, 8))
	track.Add(0, MetaTempo(120))

	track.Add(0, ch0.ProgramChange(5))
	track.Add(0, ch1.ProgramChange(46))
	track.Add(0, ch2.ProgramChange(70))

	track.Add(0, ch2.NoteOn(48, 96))
	track.Add(0, ch2.NoteOn(60, 96))

	track.Add(res.Ticks4th(), ch1.NoteOn(67, 64))
	track.Add(res.Ticks4th(), ch0.NoteOn(76, 32))

	track.Add(res.Ticks4th()*2, ch2.NoteOffVelocity(48, 64))

	track.Add(0, ch2.NoteOffVelocity(60, 64))
	track.Add(0, ch1.NoteOffVelocity(67, 64))
	track.Add(0, ch0.NoteOffVelocity(76, 64))

	smf := New()
	smf.TimeFormat = res
	smf.AddAndClose(0, track)
	err := smf.WriteTo(&bf)

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
		track *Track
		ticks = MetricTicks(96)
		smf   = New()
		beat  = ticks.Ticks4th()
		ch0   = Channel(0)
		ch1   = Channel(1)
		ch2   = Channel(2)
	)

	smf.TimeFormat = ticks

	track = NewTrack()
	track.Add(0, MetaTimeSig(4, 4, 24, 8))
	track.Add(0, MetaTempo(120))
	smf.AddAndClose(beat*4, track)

	track = NewTrack()
	track.Add(0, ch0.ProgramChange(5))
	track.Add(beat*2, ch0.NoteOn(76, 32))
	track.Add(beat*2, ch0.NoteOn(76, 0))
	smf.AddAndClose(0, track)

	track = NewTrack()
	track.Add(0, ch1.ProgramChange(46))
	track.Add(beat, ch1.NoteOn(67, 64))
	track.Add(beat*3, ch1.NoteOn(67, 0))
	smf.AddAndClose(0, track)

	track = NewTrack()
	track.Add(0, ch2.ProgramChange(70))
	track.Add(0, ch2.NoteOn(48, 96))
	track.Add(0, ch2.NoteOn(60, 96))
	track.Add(beat*4, ch2.NoteOn(48, 0), ch2.NoteOn(60, 0))
	smf.AddAndClose(0, track)

	err := smf.WriteTo(&bf)

	if err != nil {
		t.Errorf("ERROR: %s", err.Error())
	}

	if got, want := bf.Bytes(), SpecSMF1; !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n% X\n\nwanted:\n% X\n\n", got, want)
	}

}

/*
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
*/
