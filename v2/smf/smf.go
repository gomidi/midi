package smf

import (
	"gitlab.com/gomidi/midi/v2"
)

//var ErrFinished = errors.New("SMF action finished successfully")

type Track struct {
	Events []Event
	Closed bool
}

func NewTrack() *Track {
	return &Track{}
}

func (t *Track) IsEmpty() bool {
	if t.Closed {
		return len(t.Events) == 0 || len(t.Events) == 1
	}
	return len(t.Events) == 0
}

func (t *Track) Close(deltaticks uint32) {
	if t.Closed {
		return
	}
	t.Events = append(t.Events, Event{Delta: deltaticks, Data: midi.EOT})
	t.Closed = true
}

func (t *Track) Add(deltaticks uint32, msgs ...[]byte) {
	if t.Closed {
		return
	}
	for _, msg := range msgs {
		t.Events = append(t.Events, Event{Delta: deltaticks, Data: msg})
		deltaticks = 0
	}
}

func (t *Track) SendTo(resolution MetricTicks, tc TempoChanges, receiver midi.Receiver) {
	var absDelta int64

	for _, ev := range t.Events {
		absDelta += int64(ev.Delta)
		ms := resolution.Duration(tc.TempoAt(absDelta), ev.Delta).Microseconds()
		receiver.Receive(midi.NewMessage(ev.Data), ms)
	}
}

type Event struct {
	Delta uint32
	Data  []byte
}

func (e Event) Message() midi.Message {
	return midi.NewMessage(e.Data)
}

func (e *Event) MsgType() midi.MsgType {
	return midi.GetMsgType(e.Data)
}

type TempoChange struct {
	AbsDelta int64
	BPM      float64
}

type TempoChanges []TempoChange

func (t TempoChanges) Swap(a, b int) {
	t[a], t[b] = t[b], t[a]
}

func (t TempoChanges) Len() int {
	return len(t)
}

func (t TempoChanges) Less(a, b int) bool {
	return t[a].AbsDelta < t[b].AbsDelta
}

func (s SMF) Format() uint16 {
	return s.format
}

/*
type Config struct {
	NoRunningStatus bool
	Logger          Logger
	TimeFormat      TimeFormat
	//Format          uint16 // only valid: 0,1 and 2
}
*/

type SMF struct {
	//Header       SMFHeader
	// Format is the SMF file format: SMF0, SMF1 or SMF2.
	format uint16
	//Format

	// NumTracks is the number of tracks (0 indicates that the number is not set yet).
	numTracks uint16

	tracks []*Track

	// TimeFormat is the time format (either MetricTicks or TimeCode).
	//	timeFormat TimeFormat

	tempoChanges TempoChanges

	finished bool

	//opts []Option
	//Config Config

	NoRunningStatus bool
	Logger          Logger
	TimeFormat      TimeFormat
}

func (s *SMF) TempoChanges() TempoChanges {
	return s.tempoChanges
}

func (s *SMF) Tracks() []*Track {
	return s.tracks
}

func (s *SMF) NumTracks() uint16 {
	return uint16(len(s.tracks))
}

/*
func (s *SMF) TimeFormat() TimeFormat {
	return s.timeFormat
}

type Option func(*writer)

func OptionTimeFormat(tf TimeFormat) Option {
	return func(s *writer) {
		s.SMF.timeFormat = tf
	}
}
*/

/*
func (s smf) NumTracks() uint16 {
	return s.numTracks
}
*/

/*
func (s *smf) TempoAt(absDelta int64) (bpm float64) {
	bpm = 120.00
	for _, tc := range s.TempoChanges {
		if tc.AbsDelta > absDelta {
			break
		}
		bpm = tc.BPM
	}
	return
}
*/

func (t TempoChanges) TempoAt(absDelta int64) (bpm float64) {
	bpm = 120.00
	for _, tc := range t {
		if tc.AbsDelta > absDelta {
			break
		}
		bpm = tc.BPM
	}
	return
}

/*
func (s *SMF) WriteToTrack(trackNo int16, data []byte, deltaticks uint32) {
	s.Tracks[int(trackNo)].Write(deltaticks) = append(s.Tracks[int(trackNo)], event{
		Delta: deltaticks,
		Data:  data,
	})
}
*/
