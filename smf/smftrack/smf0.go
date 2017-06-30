package smftrack

import (
	"fmt"
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfwriter"
	"io"
)

// SMF0 is a namespace for methods reading from and writing to SMF0 (singletrack) files.
type SMF0 struct{}

// ReadFrom the track from a SMF0 file
func (SMF0) ReadFrom(rd smf.Reader) (tr *Track, err error) {

	err = rd.ReadHeader()

	if err != nil {
		return
	}

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
func (SMF0) WriteTo(wr io.Writer, ticks smf.MetricTicks, tracks ...*Track) (nbytes int, err error) {
	w := smfwriter.New(wr,
		smfwriter.NumTracks(1),
		smfwriter.TimeFormat(ticks),
		smfwriter.Format(smf.SMF0),
	)

	nbytes, err = w.WriteHeader()

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
