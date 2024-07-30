package smftrack

import (
	"fmt"
	"io"
	"sort"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

// SMF1 is a namespace for methods reading from and writing to SMF1 (multitrack) files.
// However they should mostly also work with the sparingly used SMF2 (sequential tracks) files.
type SMF1 struct{}

type Option func(cfg *config)

func SkipTracks(tracknos ...int) Option {
	return func(cfg *config) {
		cfg.skipTracks = map[int]bool{}

		for _, trackno := range tracknos {
			cfg.skipTracks[trackno] = true
		}
	}
}

type config struct {
	skipTracks map[int]bool
}

// ToSMF0 converts a given SMF1 file to SMF0 and writes it to wr
// If src is no SMF1 file, an error is returned
// sysex data and meta messages other than copyright, cuepoint, marker, tempo, timesignature and keysignature
// get lost, since they can be bound to a certain track.
func (smf1 SMF1) ToSMF0(src *smf.SMF, wr io.Writer, opts ...Option) (trackNames []string, err error) {

	cfg := &config{}
	cfg.skipTracks = map[int]bool{}

	for _, opt := range opts {
		opt(cfg)
	}

	if src.Format() == 0 {
		return nil, fmt.Errorf("src is already an SMF0 file")
	}

	if src.Format() == 2 {
		return nil, fmt.Errorf("can't write SMF2 file to SMF0")
	}

	var smf0 SMF0
	var tracks []*Track
	tracks, err = smf1.ReadFrom(src)

	if err != nil {
		return
	}

	var processingTracks []*Track
	var tracknameconsolidated = map[int]string{}
	var maxTrackCh int

	for trackno, tr := range tracks {
		ch := tr.Channel()

		if ch == -2 {
			fmt.Printf("track %v %q has mixed channels, skipping\n", trackno, tr.Name())
		}

		if !cfg.skipTracks[ch] {
			if old, has := tracknameconsolidated[ch]; has {
				tracknameconsolidated[ch] = old + "/" + tr.Name()
			} else {
				tracknameconsolidated[ch] = tr.Name()
			}

			processingTracks = append(processingTracks, tr)
			if maxTrackCh < ch {
				maxTrackCh = ch
			}
		} else {
			fmt.Printf("skipping track %q with MIDI channel %v\n", tr.Name(), ch)
		}
	}

	trackNames = make([]string, maxTrackCh+1)

	for trno, trname := range tracknameconsolidated {
		if trno >= 0 {
			trackNames[trno] = trname
		}
	}

	_, err = smf0.WriteTo(wr, src.TimeFormat, processingTracks...)
	return
}

// TracksOnDifferentChannels sets the track data on a different MIDI channel for each track
func (smf1 SMF1) TracksOnDifferentChannels(src *smf.SMF, wr io.Writer) (err error) {
	if src.Format() == 0 {
		return fmt.Errorf("SMF0 file not supported")
	}

	if src.Format() == 2 {
		return fmt.Errorf("SMF2 file not supported")
	}

	var tracks []*Track
	tracks, err = smf1.ReadFrom(src)

	if err != nil {
		return
	}

	var channel, val1, val2 uint8
	var rel int16
	var abs uint16

	for i, tr := range tracks {
		callback := func(ev Event) {
			if ev.Message.GetChannel(&channel) {
				ch := uint8(i % 12)

				switch {
				case ev.Message.GetNoteStart(&channel, &val1, &val2):
					ev.Message = smf.Message(midi.NoteOn(ch, val1, val2))
				case ev.Message.GetNoteEnd(&channel, &val1):
					ev.Message = smf.Message(midi.NoteOff(ch, val1))
				case ev.Message.GetAfterTouch(&channel, &val1):
					ev.Message = smf.Message(midi.AfterTouch(ch, val1))
				case ev.Message.GetPolyAfterTouch(&channel, &val1, &val2):
					ev.Message = smf.Message(midi.PolyAfterTouch(ch, val1, val2))
				case ev.Message.GetProgramChange(&channel, &val1):
					ev.Message = smf.Message(midi.ProgramChange(ch, val1))
				case ev.Message.GetControlChange(&channel, &val1, &val2):
					ev.Message = smf.Message(midi.ControlChange(ch, val1, val2))
				case ev.Message.GetPitchBend(&channel, &rel, &abs):
					ev.Message = smf.Message(midi.Pitchbend(ch, rel))
				}

				tr.UpdateEvents(ev)
			}
		}
		tr.EachEvent(callback)
	}

	_, err = smf1.WriteTo(wr, src.TimeFormat, tracks...)

	return
}

// AddTracks adds the given tracks to the tracks in src and writes the resulting SMF1 to wr.
func (smf1 SMF1) AddTracks(src *smf.SMF, wr io.Writer, tracks ...*Track) error {
	if src.Format() != 1 {
		return fmt.Errorf("can only add tracks to SMF1 file, got %v", src.Format())
	}

	oldTracks, err := smf1.ReadFrom(src)

	if err != nil {
		return err
	}

	newTracksSorted := Tracks(tracks)
	sort.Sort(newTracksSorted)

	_, err = smf1.WriteTo(wr, src.TimeFormat, append(Tracks(oldTracks), newTracksSorted...)...)
	return err
}

// ReadFrom reads the tracks with the given tracknos from rd.
// It returns an error for a SMF0 file.
// if no tracknos are given, all tracks are returned
func (SMF1) ReadFrom(rd *smf.SMF, tracknos ...uint16) (tracks []*Track, err error) {
	if rd.Format() == 0 {
		return nil, fmt.Errorf("can't get tracks from SMF0 file")
	}

	var match = map[uint16]bool{}

	for _, trackno := range tracknos {
		match[trackno] = true
	}

	var matchAll bool
	if len(tracknos) == 0 {
		matchAll = true
	}

	var found = map[uint16]bool{}

	var absPos uint64
	currentTr := &Track{}

	for i, tr := range rd.Tracks {
		if matchAll || match[uint16(i)] {

			currentTr.Number = uint16(i)
			found[uint16(i)] = true

			for _, ev := range tr {
				switch {
				case ev.Message.Is(smf.MetaEndOfTrackMsg):
					tracks = append(tracks, currentTr)
					absPos = 0
					currentTr = &Track{}
				default:
					absPos += uint64(ev.Delta)
					currentTr.addMessage(absPos, ev.Message)
				}
			}
		}
	}

	if !matchAll {
		var notFound []uint16

		for tn := range match {
			if !found[tn] {
				notFound = append(notFound, tn)
			}
		}

		if len(notFound) > 0 {
			err = fmt.Errorf("could not find tracks %v (%v tracks in source)", notFound, len(rd.Tracks))
		}

	}
	return
}

// WriteTo writes a SMF1 file of the given tracks to the given io.Writer
// Tracks are ordered by Track.Number
func (SMF1) WriteTo(wr io.Writer, timeFormat smf.TimeFormat, tracks ...*Track) (nbytes int64, err error) {
	sm := smf.NewSMF1()
	sm.TimeFormat = timeFormat
	nbytes = smfHeaderLen
	sortedTracks := Tracks(tracks)
	sort.Sort(sortedTracks)
	var nn int64

	for _, tr := range sortedTracks {
		nn, err = tr.WriteTo(sm)
		nbytes += nn
		if err != nil {
			return
		}
	}

	_, err = sm.WriteTo(wr)

	return

}

// Remove writes all tracks from rd to wr except the given track ids
// If rd is a SMF0 file it returns an error
func (SMF1) Remove(rd *smf.SMF, wr io.Writer, tracknos ...uint16) (err error) {
	if rd.Format() == 0 {
		return fmt.Errorf("can't remove from SMF0 file")
	}

	var shouldSkip = map[uint16]bool{}

	for _, trackno := range tracknos {
		shouldSkip[trackno] = true
	}

	sm := smf.NewSMF1()
	sm.TimeFormat = rd.TimeFormat
	var found = map[uint16]bool{}

	for i, tr := range rd.Tracks {
		var newTr smf.Track

		if shouldSkip[uint16(i)] {
			found[uint16(i)] = true
			continue
		}

		for _, ev := range tr {
			newTr.Add(ev.Delta, ev.Message)
		}

		newTr.Close(0)
		err = sm.Add(newTr)
		if err != nil {
			return
		}
	}

	var notFound []uint16

	for tn := range shouldSkip {
		if !found[tn] {
			notFound = append(notFound, tn)
		}
	}

	if len(notFound) > 0 {
		err = fmt.Errorf("could not find tracks %v (%v tracks in source)", notFound, len(rd.Tracks))
	}

	_, err = sm.WriteTo(wr)
	return
}
