package smftrack

import (
	"sort"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/meta"
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

// addEvent add events to the track
func (t *Track) addEvent(ev Event) {
	// don't add endoftrack messages
	if ev.Message != meta.EndOfTrack {
		ev.no = uint64(t.Number+1)*10000000 + t.lastEvendNo
		t.lastEvendNo++
		t.events = append(t.events, ev)
	}
}

// AddEvents add events to the track
func (t *Track) AddEvents(events ...Event) {
	for _, ev := range events {
		t.addEvent(ev)
	}
}

// Next returns the next position after absPos where an event happens
// absPos is returned if there is no next event
func (t *Track) Next(absPos uint64) uint64 {

	sort.Sort(t.events)

	for _, e := range t.events {
		if e.AbsTicks > absPos {
			return e.AbsTicks
		}

	}

	return absPos
}

// UpdateEvents update the events with the same Number
func (t *Track) UpdateEvents(events ...Event) {
	// TODO check if the tracknumber matches based on track event number calculation
	updaters := map[uint64]*Event{}

	for _, ev := range events {
		updaters[ev.Number()] = &ev
	}

	var evts Events

	for _, e := range t.events {
		if up := updaters[e.Number()]; up != nil {
			evts = append(evts, *up)
		} else {
			evts = append(evts, e)
		}
	}

	t.events = evts
}

// Remove removes Event from the track by matching the number
func (t *Track) RemoveEvents(events ...Event) {
	// TODO check if the tracknumber matches based on track event number calculation
	skip := map[uint64]bool{}

	for _, ev := range events {
		skip[ev.Number()] = true
	}

	var evts Events

	for _, e := range t.events {
		if !skip[e.Number()] {
			evts = append(evts, e)
		}
	}

	t.events = evts
}

func NewEvent(absTicks uint64, msg midi.Message) Event {
	return Event{AbsTicks: absTicks, Message: msg}
}

func (t *Track) addMessage(absTicks uint64, msg midi.Message) {
	t.addEvent(NewEvent(absTicks, msg))
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
