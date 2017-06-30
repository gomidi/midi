package smftrack

import (
	"fmt"
	"io"
	"sort"

	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf/smfreader"
	"github.com/gomidi/midi/smf/smfwriter"
	// "io"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/smf"
)

// Tracks helps sorting tracks
type Tracks []*Track

func (e Tracks) Len() int {
	return len(e)
}

func (e Tracks) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Tracks) Less(i, j int) bool {
	return e[i].Number < e[j].Number
}

type Track struct {
	Number uint16
	// map absolute ticks to messages
	events Events
	smf.Chunk
	lastEvendNo uint64
}

func New(number uint16) *Track {
	return &Track{Number: number}
}

// AddEvent adds an Event to the track
func (t *Track) AddEvent(ev *Event) {
	// don't add endoftrack messages
	if ev.Message != meta.EndOfTrack {
		ev.no = uint64(t.Number+1)*10000000 + t.lastEvendNo
		t.lastEvendNo++
		t.events = append(t.events, ev)
	}
}

// Remove removes Event from the track by the number
func (t *Track) RemoveEvents(numbers ...uint64) {
	skip := map[uint64]bool{}

	for _, no := range numbers {
		skip[no] = true
	}

	var evts Events

	for _, e := range t.events {
		if !skip[e.Number()] {
			evts = append(evts, e)
		}
	}

	t.events = evts
}

func NewEvent(absTicks uint64, msg midi.Message) *Event {
	return &Event{AbsTicks: absTicks, Message: msg}
}

func (t *Track) addMessage(absTicks uint64, msg midi.Message) {
	t.AddEvent(NewEvent(absTicks, msg))
}

// GetSMF0Track gets the track from a SMF0 file
func GetSMF0Track(rd smf.Reader) (tr *Track, err error) {

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

// GetSMF1Tracks gets the tracks with the given tracknos from rd.
// It returns an error for a SMF0 file.
// if no tracknos are given, all tracks are returned
func GetSMF1Tracks(rd smf.Reader, tracknos ...uint16) (tracks []*Track, err error) {

	err = rd.ReadHeader()

	if err != nil {
		return
	}

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
		if err != nil {
			return nil, err
		}

		if matchAll || match[uint16(rd.Track())] {

			currentTr.Number = uint16(rd.Track())
			found[uint16(rd.Track())] = true

			// don't write meta.EndOfTrack since track is handling it on its own
			if msg == meta.EndOfTrack {

				tracks = append(tracks, currentTr)
				absPos = 0
				currentTr = &Track{}
			} else {
				absPos += uint64(rd.Delta())
				currentTr.addMessage(absPos, msg)
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
			err = fmt.Errorf("could not find tracks %v (%v tracks in source)", notFound, rd.Header().NumTracks)
		}

	}
	return
}

// WriteSMF0 merges the given tracks to an SMF0 file and writes it to writer
// sysex data and meta messages other than copyright, cuepoint, marker, tempo, timesignature and keysignature
// get lost
func WriteSMF0(wr io.Writer, ticks smf.MetricTicks, tracks ...*Track) (nbytes int, err error) {
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

// WriteSMF1 writes a SMF1 file of the given tracks to the given io.Writer
// Tracks are ordered by Track.Number
func WriteSMF1(wr io.Writer, ticks smf.MetricTicks, tracks ...*Track) (nbytes int, err error) {
	w := smfwriter.New(wr,
		smfwriter.NumTracks(uint16(len(tracks))),
		smfwriter.TimeFormat(ticks),
		smfwriter.Format(smf.SMF1),
	)

	nbytes, err = w.WriteHeader()

	if err != nil {
		return
	}

	sortedTracks := Tracks(tracks)

	sort.Sort(sortedTracks)

	var n int
	for _, tr := range sortedTracks {
		n, err = tr.WriteTo(w)
		nbytes += n
		if err != nil {
			return
		}
	}

	return

}

// WriteTo writes the track to the given SMF writer
func (tr *Track) WriteTo(wr smf.Writer) (nbytes int, err error) {
	nbytes, err = wr.WriteHeader()
	if err != nil {
		return
	}
	sort.Sort(tr.events)

	var lastAbs uint64 = 0

	var n int
	for _, ev := range tr.events {
		delta := ev.AbsTicks - lastAbs
		wr.SetDelta(uint32(delta))
		lastAbs = ev.AbsTicks
		n, err = wr.Write(ev.Message)
		nbytes += n
		if err != nil {
			return
		}
	}

	n, err = wr.Write(meta.EndOfTrack)
	nbytes += n

	return
}

// Remove writes all tracks from rd to wr except the given track ids
// If rd is a SMF0 file it returns an error
func Remove(rd smf.Reader, wr io.Writer, tracknos ...uint16) (err error) {
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

	return nil
}

/*


// Track interface allows modification of midi tracks
// it relies on an absolute position; i.e. the max length is defined by uint64
// the track will grow as needed
// everything is handled by absolute time
type Track interface {
	Cursor() uint64
	SetCursor(abstime uint64) // sets cursor absolut time

	// GetMessages returns the events at the current position
	GetMessages() []midi.Message
	// adds the event at the current position
	AddMessage(midi.Message)

	RemoveMessages(num int)                // removes the given number of events at the current position
	MoveMessage(idx int, to uint64)        // moves the event with index idx at the current position to the given position
	MoveSlice(until uint64, target uint64) // moves all events between the current position and until to target (is the left/starting point)

	Len() uint64 // absolute length

	Cut(until uint64) // cuts from the current position to until

	Save() error // writes the track back to the file

	NextMessages() []midi.Message // returns the next events inside the track (from the current position), at the end, nil is returned

	PrevMessages() []midi.Message // returns the prev events inside the track (from the current position), at the start, nil is returned
}

func Get(f io.ReadWriteSeeker, trackno uint8) (Track, error) {
	return nil, nil
}
*/
