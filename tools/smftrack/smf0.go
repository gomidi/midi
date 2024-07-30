package smftrack

import (
	"fmt"
	"io"

	"gitlab.com/gomidi/midi/v2/smf"
)

// SMF0 is a namespace for methods reading from and writing to SMF0 (singletrack) files.
type SMF0 struct{}

// ReadFrom the track from a SMF0 file
func (SMF0) ReadFrom(rd *smf.SMF) (tr *Track, err error) {
	if rd.Format() != 0 {
		return nil, fmt.Errorf("wrong file SMF format: %v", rd.Format())
	}

	var absPos uint64
	tr = &Track{
		Number: 0,
	}

	for _, ev := range rd.Tracks[0] {
		// don't write meta.EndOfTrack since track is handling it on its own
		if ev.Message.Is(smf.MetaEndOfTrackMsg) {
			break
		}

		absPos += uint64(ev.Delta)
		tr.addMessage(absPos, ev.Message)
	}
	return
}

const smfHeaderLen int64 = 14

// WriteTo merges the given tracks to an SMF0 file and writes it to writer
// sysex data and meta messages other than copyright, cuepoint, marker, tempo, timesignature and keysignature
// get lost
// func (SMF0) WriteTo(wr io.Writer, timeformat smf.TimeFormat, tracks ...*Track) (nbytes int64, err error) {
func (SMF0) WriteTo(wr io.Writer, timeformat smf.TimeFormat, tracks ...*Track) (nbytes int64, err error) {
	var sm = smf.New()
	sm.TimeFormat = timeformat
	nbytes = smfHeaderLen // SMF file header size

	mergeTrack := &Track{}

	for _, tr := range tracks {
		for _, ev := range tr.events {

			shouldAdd := false
			var channel uint8

			// skipping sysex and most meta messages
			switch {
			case ev.Message.IsOneOf(smf.MetaCopyrightMsg, smf.MetaCuepointMsg, smf.MetaMarkerMsg, smf.MetaTempoMsg, smf.MetaTimeSigMsg, smf.MetaKeySigMsg):
				shouldAdd = true
			case ev.Message.GetChannel(&channel):
				shouldAdd = true
			default:
				shouldAdd = false
			}

			if shouldAdd {
				mergeTrack.addMessage(ev.AbsTicks, ev.Message)
			}

		}

	}

	var n int64
	n, err = mergeTrack.WriteTo(sm)
	nbytes += n
	if err != nil {
		return
	}

	_, err = sm.WriteTo(wr)
	return
}

// ToSMF1 converts a given SMF0 file to SMF1 and writes it to wr
// If src is no SMF0 file, an error is returned
// channel messages are distributed over the tracks by their channels
// e.g. channel 0 -> track 1, channel 1 -> track 2 etc.
// and everything else stays in track 0
func (smf0 SMF0) ToSMF1(src *smf.SMF, wr io.Writer) (err error) {

	tr, err := smf0.ReadFrom(src)

	if err != nil {
		return err
	}

	var removedFromTrack0 []Event
	var channelTracks [16]*Track
	var channel uint8

	tr.EachEvent(func(ev Event) {
		if ev.Message.GetChannel(&channel) {
			removedFromTrack0 = append(removedFromTrack0, ev)
			chr := channelTracks[int(channel)]
			if chr == nil {
				chr = New(uint16(channel) + 1)
				channelTracks[int(channel)] = chr
			}
			chr.AddEvents(ev)
		}
	})

	tr.RemoveEvents(removedFromTrack0...)
	tr.Number = 0

	tracks := []*Track{tr}

	for chNum, chtr := range channelTracks {
		if chtr == nil {
			continue
		}

		if chtr.Len() > 0 {
			chtr.Number = uint16(chNum) + 1
			tracks = append(tracks, chtr)
		}
	}

	_, err = (SMF1{}).WriteTo(wr, src.TimeFormat, tracks...)
	return err
}
