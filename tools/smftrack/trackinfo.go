package smftrack

import (
	"bytes"
	"fmt"
	"strings"

	"gitlab.com/gomidi/midi/v2/gm"
	"gitlab.com/gomidi/midi/v2/smf"
)

/*
TODO

make tempolane (TempoChange: Bar/Tempo) and TimeSignature lane (TimeSigChange: Bar/TimeSig)
or simply bars info (number bars, and any tempo and time signature changes (only once per bar)
*/

type TrackInfo struct {
	Number         int
	Name           string
	Instrument     string
	Channels       map[uint8]bool
	ProgramChanges map[uint8]bool
	ControlChanges map[uint8]bool
	Program        string
	Text           string
	Lyrics         string
	CopyRight      string
}

func (t *TrackInfo) trimSpace() {
	t.Name = strings.TrimSpace(t.Name)
	t.Instrument = strings.TrimSpace(t.Instrument)
	t.Program = strings.TrimSpace(t.Program)
	t.Text = strings.TrimSpace(t.Text)
	t.Lyrics = strings.TrimSpace(t.Lyrics)
	t.CopyRight = strings.TrimSpace(t.CopyRight)
}

func (t TrackInfo) String() string {
	var bf bytes.Buffer

	bf.WriteString(fmt.Sprintf("Track #%v", t.Number))

	if t.Name != "" || t.Instrument != "" {
		if t.Name != "" && t.Instrument != "" {
			bf.WriteString(fmt.Sprintf(" %q / %q", t.Name, t.Instrument))
		} else {
			if t.Name != "" {
				bf.WriteString(fmt.Sprintf(" %q", t.Name))
			} else {
				bf.WriteString(fmt.Sprintf(" %q", t.Instrument))
			}
		}
	}

	if len(t.Channels) == 1 {
		bf.WriteString(" channel ")

		for ch := range t.Channels {
			bf.WriteString(fmt.Sprintf("%v", ch))
		}
	}

	if len(t.Channels) > 1 {
		bf.WriteString(" channels ")

		for ch := range t.Channels {
			bf.WriteString(fmt.Sprintf("%v,", ch))
		}
	}

	if len(t.ProgramChanges) == 1 {
		bf.WriteString(" programchange ")

		for pc := range t.ProgramChanges {
			bf.WriteString(fmt.Sprintf("%v %q", pc, gm.Instr(pc).String()))
		}
	}

	if len(t.ProgramChanges) > 1 {
		bf.WriteString(" programchanges ")

		for pc := range t.ProgramChanges {
			bf.WriteString(fmt.Sprintf("%v %q,", pc, gm.Instr(pc).String()))
		}
	}

	if t.Program != "" {
		bf.WriteString(fmt.Sprintf(" program %q", t.Program))
	}

	/*
		if len(t.ControlChanges) == 1 {
			bf.WriteString(" CC ")

			for c, _ := range t.ControlChanges {
				bf.WriteString(fmt.Sprintf("%v %q", c, cc.Name[c]))
			}
		}

		if len(t.ControlChanges) > 1 {
			bf.WriteString(" CCs ")

			for c, _ := range t.ControlChanges {
				bf.WriteString(fmt.Sprintf("%v %q,", c, cc.Name[c]))
			}
		}
	*/

	if t.CopyRight != "" {
		bf.WriteString(fmt.Sprintf(" (c) %q", t.CopyRight))
	}

	if t.Text != "" {
		var s = t.Text
		if len(t.Text) > 10 {
			s = t.Text[:9] + "..."
		}
		bf.WriteString(fmt.Sprintf(" text '%s'", s))
	}

	if t.Lyrics != "" {
		var s = t.Lyrics
		if len(t.Lyrics) > 10 {
			s = t.Lyrics[:9] + "..."
		}
		bf.WriteString(fmt.Sprintf(" lyrics '%s'", s))
	}

	return bf.String()
}

func TracksInfos(inFile string) (tracks []TrackInfo, err error) {
	var currentTrack int

	sm, err := smf.ReadFile(inFile)
	if err != nil {
		return nil, fmt.Errorf("could not read SMF file %v\n", inFile)
	}

	for i, tr := range sm.Tracks {
		currentTrack = i

		for _, ev := range tr {
			tracks = append(tracks, TrackInfo{Number: 0, Channels: map[uint8]bool{}, ProgramChanges: map[uint8]bool{}, ControlChanges: map[uint8]bool{}})

			var text string
			var channel, val1, val2 uint8
			var rel int16
			var abs uint16

			switch {
			case ev.Message.GetMetaProgramName(&text):
				tracks[currentTrack].Program += text + " "
			case ev.Message.GetProgramChange(&channel, &val1):
				tracks[currentTrack].ProgramChanges[val1] = true
				tracks[currentTrack].Channels[channel] = true
			case ev.Message.GetMetaInstrument(&text):
				tracks[currentTrack].Instrument += text + " "
			case ev.Message.GetMetaTrackName(&text):
				tracks[currentTrack].Name += text + " "
			case ev.Message.GetMetaText(&text):
				tracks[currentTrack].Text += text + " "
			case ev.Message.GetMetaLyric(&text):
				tracks[currentTrack].Lyrics += text + " "
			case ev.Message.GetMetaCopyright(&text):
				tracks[currentTrack].CopyRight += text + " "
			case ev.Message.GetNoteOn(&channel, &val1, &val2):
				tracks[currentTrack].Channels[channel] = true
			case ev.Message.GetNoteOff(&channel, &val1, &val2):
				tracks[currentTrack].Channels[channel] = true
			case ev.Message.GetAfterTouch(&channel, &val1):
				tracks[currentTrack].Channels[channel] = true
			case ev.Message.GetPolyAfterTouch(&channel, &val1, &val2):
				tracks[currentTrack].Channels[channel] = true
			case ev.Message.GetPitchBend(&channel, &rel, &abs):
				tracks[currentTrack].Channels[channel] = true
			case ev.Message.GetControlChange(&channel, &val1, &val2):
				tracks[currentTrack].ControlChanges[val1] = true
				tracks[currentTrack].Channels[channel] = true
			}
		}
	}

	if len(tracks) < 2 {
		return nil, nil
	}

	for i, tr := range tracks {
		tr.trimSpace()
		tracks[i] = tr
	}

	return tracks[0 : len(tracks)-2], nil
}
