package smftrack

import (
	"fmt"
	"io"
	"sort"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfreader"
	"github.com/gomidi/midi/smf/smfwriter"
)

// SMF1 is a namespace for methods reading from and writing to SMF1 (multitrack) files.
// However they should mostly also work with the sparingly used SMF2 (sequential tracks) files.
type SMF1 struct{}

// ToSMF0 converts a given SMF1 file to SMF0 and writes it to wr
// If src is no SMF1 file, an error is returned
// sysex data and meta messages other than copyright, cuepoint, marker, tempo, timesignature and keysignature
// get lost, since they can be bound to a certain track.
func (smf1 SMF1) ToSMF0(src smf.Reader, wr io.Writer) (trackNames []string, err error) {
	src.ReadHeader()

	if src.Header().Format == smf.SMF0 {
		return nil, fmt.Errorf("src is already an SMF0 file")
	}

	if src.Header().Format == smf.SMF2 {
		return nil, fmt.Errorf("can't write SMF2 file to SMF0")
	}

	var smf0 SMF0
	var tracks []*Track
	tracks, err = smf1.ReadFrom(src)

	if err != nil {
		return
	}

	for _, tr := range tracks {
		trackNames = append(trackNames, tr.Name())
	}

	_, err = smf0.WriteTo(wr, src.Header().TimeFormat, tracks...)
	return
}

// TracksOnDifferentChannels sets the track data on a different MIDI channel for each track
func (smf1 SMF1) TracksOnDifferentChannels(src smf.Reader, wr io.Writer) (err error) {
	src.ReadHeader()

	if src.Header().Format == smf.SMF0 {
		return fmt.Errorf("SMF0 file not supported")
	}

	if src.Header().Format == smf.SMF2 {
		return fmt.Errorf("SMF2 file not supported")
	}

	var tracks []*Track
	tracks, err = smf1.ReadFrom(src)

	if err != nil {
		return
	}

	for i, tr := range tracks {
		callback := func(ev Event) {
			if _, is := ev.Message.(channel.Message); is {
				ch := uint8(i % 12)

				switch v := ev.Message.(type) {
				case channel.NoteOn:
					ev.Message = channel.New(ch).NoteOn(v.Key(), v.Velocity())
				case channel.NoteOff:
					ev.Message = channel.New(ch).NoteOff(v.Key())
				case channel.AfterTouch:
					ev.Message = channel.New(ch).AfterTouch(v.Pressure())
				case channel.PolyphonicAfterTouch:
					ev.Message = channel.New(ch).PolyphonicAfterTouch(v.Key(), v.Pressure())
				case channel.ProgramChange:
					ev.Message = channel.New(ch).ProgramChange(v.Program())
				case channel.ControlChange:
					ev.Message = channel.New(ch).ControlChange(v.Controller(), v.Value())
				case channel.PitchBend:
					ev.Message = channel.New(ch).PitchBend(v.Value())
				}
				tr.UpdateEvents(ev)
			}
		}
		tr.EachEvent(callback)
	}

	_, err = smf1.WriteTo(wr, src.Header().TimeFormat, tracks...)

	return
}

// AddTracks adds the given tracks to the tracks in src and writes the resulting SMF1 to wr.
func (smf1 SMF1) AddTracks(src smf.Reader, wr io.Writer, tracks ...*Track) error {
	src.ReadHeader()

	if src.Header().Format != smf.SMF1 {
		return fmt.Errorf("can only add tracks to SMF1 file, got %s", src.Header().Format)
	}

	oldTracks, err := smf1.ReadFrom(src)

	if err != nil {
		return err
	}

	newTracksSorted := Tracks(tracks)
	sort.Sort(newTracksSorted)

	_, err = smf1.WriteTo(wr, src.Header().TimeFormat, append(Tracks(oldTracks), newTracksSorted...)...)
	return err
}

// ReadFrom reads the tracks with the given tracknos from rd.
// It returns an error for a SMF0 file.
// if no tracknos are given, all tracks are returned
func (SMF1) ReadFrom(rd smf.Reader, tracknos ...uint16) (tracks []*Track, err error) {

	rd.ReadHeader()

	if rd.Header().Format == smf.SMF0 {
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

	var msg midi.Message

	for {
		msg, err = rd.Read()

		if matchAll || match[uint16(rd.Track())] {

			currentTr.Number = uint16(rd.Track())
			found[uint16(rd.Track())] = true

			// don't write meta.EndOfTrack since track is handling it on its own
			if msg == meta.EndOfTrack {

				tracks = append(tracks, currentTr)
				absPos = 0
				currentTr = &Track{}
			} else {
				if err == nil {
					absPos += uint64(rd.Delta())
					currentTr.addMessage(absPos, msg)
				}
			}
		}

		if err != nil {
			//if err == smfreader.ErrFinished || err == io.EOF {
			if err == smfreader.ErrFinished {
				err = nil
				break
			}
			return nil, err
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
			err = fmt.Errorf("could not find tracks %v (%v tracks in source)", notFound, rd.Header().NumTracks)
		}

	}
	return
}

// WriteTo writes a SMF1 file of the given tracks to the given io.Writer
// Tracks are ordered by Track.Number
func (SMF1) WriteTo(wr io.Writer, timeFormat smf.TimeFormat, tracks ...*Track) (nbytes int64, err error) {
	w := smfwriter.New(wr,
		smfwriter.NumTracks(uint16(len(tracks))),
		smfwriter.TimeFormat(timeFormat),
		smfwriter.Format(smf.SMF1),
	)

	var n int
	n, err = w.WriteHeader()

	nbytes += int64(n)

	if err != nil {
		return
	}

	sortedTracks := Tracks(tracks)

	sort.Sort(sortedTracks)

	var nn int64

	for _, tr := range sortedTracks {
		nn, err = tr.WriteTo(w)
		nbytes += nn
		if err != nil {
			return
		}
	}

	return

}

// Remove writes all tracks from rd to wr except the given track ids
// If rd is a SMF0 file it returns an error
func (SMF1) Remove(rd smf.Reader, wr io.Writer, tracknos ...uint16) (err error) {
	err = rd.ReadHeader()
	if err != nil {
		return
	}

	if rd.Header().Format == smf.SMF0 {
		return fmt.Errorf("can't remove from SMF0 file")
	}

	var shouldSkip = map[uint16]bool{}

	for _, trackno := range tracknos {
		shouldSkip[trackno] = true
	}

	w := smfwriter.New(wr,
		smfwriter.Format(rd.Header().Format),
		smfwriter.TimeFormat(rd.Header().TimeFormat),
		smfwriter.NumTracks(rd.Header().NumTracks-uint16(len(tracknos))),
	)

	_, err = w.WriteHeader()
	if err != nil {
		return err
	}

	var found = map[uint16]bool{}

	var msg midi.Message

	for {
		msg, err = rd.Read()
		if err != nil {
			if err == smfreader.ErrFinished {
				break
			}
			return
		}

		if shouldSkip[uint16(rd.Track())] {
			found[uint16(rd.Track())] = true
			continue
		}

		w.SetDelta(rd.Delta())
		_, err = w.Write(msg)
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
		err = fmt.Errorf("could not find tracks %v (%v tracks in source)", notFound, rd.Header().NumTracks)
	}

	return
}
