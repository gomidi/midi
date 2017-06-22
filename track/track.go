package track

import (
	"io"

	midi "github.com/gomidi/midi"
	"github.com/gomidi/midi/smf"
)

// records live input from srcLive and writes it to destFile until error or EOF
// any tempo changes, time signatures etc that come in will make it to the file
// recording starts with the first incoming end of track message
// recording stops with the next incoming end of track message (and so on)
// Record won't return until everything is written to destfile.
// to what tracks what is recorded depends on TrackReader and TrackWriter
func Record(from midi.Reader, to smf.Writer) error {
	return nil
}

// adds a track to the given SMF file
func Add(wr io.ReadWriteSeeker, track smf.Writer) error {
	return nil
}

// removes a track to the given SMF file
func Remove(wr io.ReadWriteSeeker, trackno uint8) error {
	return nil
}

// Track interface allows modification of midi tracks
// it relies on an absolute position; i.e. the max length is defined by uint64
// the track will grow as needed
// everything is handled by absolute time
type Track interface {
	Cursor() uint64
	SetCursor(abstime uint64) // sets cursor absolut time

	// GetEvents returns the events at the current position
	GetEvents() []midi.Event
	// adds the event at the current position
	AddEvent(midi.Event)

	RemoveEvents(num int)                  // removes the given number of events at the current position
	MoveEvent(idx int, to uint64)          // moves the event with index idx at the current position to the given position
	MoveSlice(until uint64, target uint64) // moves all events between the current position and until to target (is the left/starting point)

	Len() uint64 // absolute length

	Cut(until uint64) // cuts from the current position to until

	Save() error // writes the track back to the file

	NextEvents() []midi.Event // returns the next events inside the track (from the current position), at the end, nil is returned

	PrevEvents() []midi.Event // returns the prev events inside the track (from the current position), at the start, nil is returned
}

func Get(f io.ReadWriteSeeker, trackno uint8) (Track, error) {
	return nil, nil
}
