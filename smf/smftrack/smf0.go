package smftrack

import (
	"fmt"
	"io"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfreader"
	"github.com/gomidi/midi/smf/smfwriter"
)

// SMF0 is a namespace for methods reading from and writing to SMF0 (singletrack) files.
type SMF0 struct{}

// ReadFrom the track from a SMF0 file
func (SMF0) ReadFrom(rd smf.Reader) (tr *Track, err error) {

	rd.ReadHeader()

	if rd.Header().Format != smf.SMF0 {
		return nil, fmt.Errorf("wrong file format: %s", rd.Header().Format)
	}

	var absPos uint64
	tr = &Track{
		Number: 0,
	}

	var msg midi.Message

	for {
		msg, err = rd.Read()

		if err != nil {
			if err == smfreader.ErrFinished {
				err = nil
				break
			}
			return nil, err
		}

		// don't write meta.EndOfTrack since track is handling it on its own
		if msg != meta.EndOfTrack {
			absPos += uint64(rd.Delta())

			tr.addMessage(absPos, msg)
		}

	}

	return
}

// WriteTo merges the given tracks to an SMF0 file and writes it to writer
// sysex data and meta messages other than copyright, cuepoint, marker, tempo, timesignature and keysignature
// get lost
func (SMF0) WriteTo(wr io.Writer, timeformat smf.TimeFormat, tracks ...*Track) (nbytes int64, err error) {
	w := smfwriter.New(wr,
		smfwriter.NumTracks(1),
		smfwriter.TimeFormat(timeformat),
		smfwriter.Format(smf.SMF0),
	)

	var n int
	n, err = w.WriteHeader()

	nbytes += int64(n)

	if err != nil {
		return
	}

	mergeTrack := &Track{}

	for _, tr := range tracks {
		for _, ev := range tr.events {

			shouldAdd := false

			// skipping sysex and most meta messages
			switch ev.Message.(type) {

			// allow meta messages that are inpedendant of a track
			case meta.Copyright, meta.CuePoint, meta.Marker, meta.Tempo, meta.TimeSignature, meta.KeySignature:
				shouldAdd = true
			default:
				_, shouldAdd = ev.Message.(channel.Message)

			}

			if shouldAdd {
				mergeTrack.addMessage(ev.AbsTicks, ev.Message)
			}

		}

	}

	return mergeTrack.WriteTo(w)
}

// ToSMF1 converts a given SMF0 file to SMF1 and writes it to wr
// If src is no SMF0 file, an error is returned
// channel messages are distributed over the tracks by their channels
// e.g. channel 0 -> track 1, channel 1 -> track 2 etc.
// and everything else stays in track 0
func (smf0 SMF0) ToSMF1(src smf.Reader, wr io.Writer) (err error) {

	tr, err := smf0.ReadFrom(src)

	if err != nil {
		return err
	}

	var removedFromTrack0 []Event
	var channelTracks [16]*Track

	tr.EachEvent(func(ev Event) {
		if chMsg, ok := ev.Message.(channel.Message); ok {
			removedFromTrack0 = append(removedFromTrack0, ev)
			chr := channelTracks[int(chMsg.Channel())]
			if chr == nil {
				chr = New(uint16(chMsg.Channel()) + 1)
				channelTracks[int(chMsg.Channel())] = chr
			}
			// fmt.Printf("adding %s\n", chMsg)
			chr.AddEvents(ev)
		}
	})

	tr.RemoveEvents(removedFromTrack0...)

	// var num = uint16(0)

	tr.Number = 0

	tracks := []*Track{tr}

	for chNum, chtr := range channelTracks {
		if chtr == nil {
			continue
		}
		// fmt.Printf("got track: %v\n", chtr)
		if chtr.Len() > 0 {
			chtr.Number = uint16(chNum) + 1
			tracks = append(tracks, chtr)
		}
	}

	_, err = (SMF1{}).WriteTo(wr, src.Header().TimeFormat, tracks...)
	return err
}
