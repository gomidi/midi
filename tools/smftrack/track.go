package smftrack

import (
	"sort"

	"gitlab.com/gomidi/midi/v2/smf"
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
	Number      uint16
	events      Events
	lastEvendNo uint64
	instrument  string
	program     string
	trackName   string
	channel     int // main MIDI channel for the track -1 if no channel events are on the track
}

func (t Track) Name() string {
	if t.trackName != "" {
		return t.trackName
	}

	if t.instrument != "" {
		return t.instrument
	}

	if t.program != "" {
		return t.program
	}

	return ""
}

func New(number uint16) *Track {
	return &Track{Number: number}
}

func (t *Track) Len() int {
	return len(t.events)
}

// addEvent add events to the track
func (t *Track) addEvent(ev Event) {
	// don't add endoftrack messages
	if !ev.Message.Is(smf.MetaEndOfTrackMsg) {
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

func (t *Track) GetEventsAt(absPos uint64) (evts []Event) {
	if t.events == nil {
		return
	}
	sort.Sort(t.events)

	for _, e := range t.events {
		if e.AbsTicks == absPos {
			evts = append(evts, e)
		}

		if e.AbsTicks > absPos {
			return
		}

	}
	return
}

func (t *Track) EachEvent(callback func(Event)) {
	if t.events == nil {
		return
	}
	sort.Sort(t.events)

	for _, e := range t.events {
		callback(e)
	}
}

// NextPosition returns the next position after absPos where an event happens
// absPos is returned if there is no next event
func (t *Track) NextPosition(absPos uint64) uint64 {
	if t.events == nil {
		return 0
	}
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
	if t.events == nil {
		return
	}
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
	if t.events == nil {
		return
	}
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

func NewEvent(absTicks uint64, msg smf.Message) Event {
	return Event{AbsTicks: absTicks, Message: msg}
}

func (t *Track) addMessage(absTicks uint64, msg smf.Message) {

	var channel uint8
	var text string
	switch {
	case msg.GetMetaInstrument(&text):
		t.instrument = text
	case msg.GetMetaProgramName(&text):
		t.program = text
	case msg.GetMetaTrackName(&text):
		t.trackName = text
	case msg.GetChannel(&channel):
		if t.channel == 0 {
			t.channel = int(channel) + 1
		} else {
			if t.channel != int(channel)+1 {
				t.channel = -1
			}
		}
	}
	t.addEvent(NewEvent(absTicks, msg))
}

func (tr *Track) Channel() int {
	return tr.channel - 1
}

// WriteTo writes the track to the given SMF writer
// TODO: get the correct nbytes by writing to a buffer somehow
// and getting the length of the buffer
// func (tr *Track) WriteTo(wr smf.Writer) (nbytes int64, err error) {
func (tr *Track) WriteTo(wr *smf.SMF) (nbytes int64, err error) {
	var track smf.Track
	sort.Sort(tr.events)

	var lastAbs uint64 = 0

	for _, ev := range tr.events {
		delta := ev.AbsTicks - lastAbs
		track.Add(uint32(delta), ev.Message.Bytes())
		lastAbs = ev.AbsTicks
		nbytes += int64(len(ev.Message.Bytes()))
	}

	track.Close(0)
	wr.Add(track)
	nbytes += int64(len(smf.EOT.Bytes()))
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
