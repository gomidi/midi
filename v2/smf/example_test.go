package smf_test

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/gm"
	smf "gitlab.com/gomidi/midi/v2/smf"
)

func Example() {
	// we write a SMF file into a buffer and read it back

	var (
		bf                   bytes.Buffer
		clock                = smf.MetricTicks(96) // resolution: 96 ticks per quarternote 960 is also a common choice
		general, piano, bass smf.Track             // our tracks
	)

	// first track must have tempo and meter informations
	general.Add(0, smf.MetaTrackSequenceName("general"))
	general.Add(0, smf.MetaMeter(3, 4))
	general.Add(0, smf.MetaTempo(140))
	general.Add(clock.Ticks4th()*6, smf.MetaTempo(130))
	general.Add(clock.Ticks4th(), smf.MetaTempo(135))
	general.Close(0) // don't forget to close a track

	piano.Add(0, smf.MetaInstrument("Piano"))
	piano.Add(0, midi.ProgramChange(0, gm.Instr_HonkytonkPiano.Value()))
	piano.Add(0, midi.NoteOn(0, 76, 120))
	// duration: a quarter note (96 ticks in our case)
	piano.Add(clock.Ticks4th(), midi.NoteOff(0, 76))
	piano.Close(0)

	bass.Add(0, smf.MetaInstrument("Bass"))
	bass.Add(0, midi.ProgramChange(1, gm.Instr_AcousticBass.Value()))
	bass.Add(clock.Ticks4th(), midi.NoteOn(1, 47, 64))
	bass.Add(clock.Ticks4th()*3, midi.NoteOff(1, 47))
	bass.Close(0)

	// create the SMF and add the tracks
	s := smf.New()
	s.TimeFormat = clock
	s.Add(general)
	s.Add(piano)
	s.Add(bass)

	// write the bytes to the buffer
	err := s.WriteTo(&bf)

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}

	// read the bytes
	s, err = smf.ReadFrom(bytes.NewReader(bf.Bytes()))

	if err != nil {
		fmt.Printf("ERROR: %s", err.Error())
		return
	}

	fmt.Printf("got %v tracks\n", len(s.Tracks))

	for no, track := range s.Tracks {

		// it might be a good idea to go from delta ticks to absolute ticks.
		var absTicks uint64

		var trackname string
		var channel, program uint8
		var gm_name string

		for _, ev := range track {
			absTicks += uint64(ev.Delta)
			msg := ev.Message

			if msg.Type() == smf.MetaEndOfTrackMsg {
				// ignore
				continue
			}

			switch {
			case msg.GetMetaTrackName(&trackname): // set the trackname
			case msg.GetMetaInstrument(&trackname): // set the trackname based on instrument name
			case msg.GetProgramChange(&channel, &program):
				gm_name = "(" + gm.Instr(program).String() + ")"
			default:
				fmt.Printf("track %v %s %s @%v %s\n", no, trackname, gm_name, absTicks, ev.Message)
			}
		}
	}

	// Output:
	// got 3 tracks
	// track 0 general  @0 MetaTimeSig meter: 3/4
	// track 0 general  @0 MetaTempo bpm: 140.00
	// track 0 general  @576 MetaTempo bpm: 130.00
	// track 0 general  @672 MetaTempo bpm: 135.00
	// track 1 Piano (HonkytonkPiano) @0 NoteOn channel: 0 key: 76 velocity: 120
	// track 1 Piano (HonkytonkPiano) @96 NoteOff channel: 0 key: 76
	// track 2 Bass (AcousticBass) @96 NoteOn channel: 1 key: 47 velocity: 64
	// track 2 Bass (AcousticBass) @384 NoteOff channel: 1 key: 47
}
