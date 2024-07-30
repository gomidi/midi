package smf

import (
	"fmt"
	"io"
	"runtime"
	"sort"
	"time"

	"reflect"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type Event struct {
	Delta   uint32
	Message Message
}

type Track []Event

func (t Track) IsClosed() bool {
	if len(t) == 0 {
		return false
	}

	last := t[len(t)-1]
	return reflect.DeepEqual(last.Message, EOT)
}

func (t Track) IsEmpty() bool {
	if t.IsClosed() {
		return len(t) == 1
	}
	return len(t) == 0
}

func (t *Track) Close(deltaticks uint32) {
	if t.IsClosed() {
		return
	}
	*t = append(*t, Event{Delta: deltaticks, Message: EOT})
}

func (t *Track) Add(deltaticks uint32, msgs ...[]byte) {
	if t.IsClosed() {
		return
	}
	for _, msg := range msgs {
		ev := Event{Delta: deltaticks, Message: msg}
		*t = append(*t, ev)
		deltaticks = 0
	}
}

func (t *Track) RecordFrom(inPort drivers.In, ticks MetricTicks, bpm float64) (stop func(), err error) {
	if !inPort.IsOpen() {
		err := inPort.Open()
		if err != nil {
			return nil, err
		}
	}
	t.Add(0, MetaTempo(bpm))
	var absmillisec int32
	return midi.ListenTo(inPort, func(msg midi.Message, absms int32) {
		deltams := absms - absmillisec
		absmillisec = absms
		delta := ticks.Ticks(bpm, time.Duration(deltams)*time.Millisecond)
		t.Add(delta, msg)
	})
}

func (t *Track) SendTo(resolution MetricTicks, tc TempoChanges, receiver func(m midi.Message, timestampms int32)) {
	var absDelta int64

	for _, ev := range *t {
		absDelta += int64(ev.Delta)
		if Message(ev.Message).IsPlayable() {
			ms := int32(resolution.Duration(tc.TempoAt(absDelta), ev.Delta).Microseconds() * 100)
			receiver(ev.Message.Bytes(), ms)
		}
	}
}

type TracksReader struct {
	smf    *SMF
	tracks map[int]bool
	filter []midi.Type
	err    error
}

func (t *TracksReader) Error() error {
	return t.err
}

func (t *TracksReader) SMF() *SMF {
	return t.smf
}

func (t *TracksReader) doTrack(tr int) bool {
	if len(t.tracks) == 0 {
		return true
	}

	return t.tracks[tr]
}

func ReadTracks(filepath string, tracks ...int) *TracksReader {
	t := &TracksReader{}
	t.tracks = map[int]bool{}
	for _, tr := range tracks {
		t.tracks[tr] = true
	}
	t.smf, t.err = ReadFile(filepath)
	if t.err != nil {
		return t
	}
	if _, ok := t.smf.TimeFormat.(MetricTicks); !ok {
		t.err = fmt.Errorf("SMF time format is not metric ticks, but %s (currently not supported)", t.smf.TimeFormat.String())
		return t
	}
	return t
}

func ReadTracksFrom(rd io.Reader, tracks ...int) *TracksReader {
	t := &TracksReader{}
	t.tracks = map[int]bool{}
	for _, tr := range tracks {
		t.tracks[tr] = true
	}

	t.smf, t.err = ReadFrom(rd)
	if t.err != nil {
		return t
	}
	if _, ok := t.smf.TimeFormat.(MetricTicks); !ok {
		t.err = fmt.Errorf("SMF time format is not metric ticks, but %s (currently not supported)", t.smf.TimeFormat.String())
		return t
	}
	return t
}

func (t *TracksReader) Only(mtypes ...midi.Type) *TracksReader {
	t.filter = mtypes
	return t
}

type TrackEvent struct {
	Event
	TrackNo         int
	AbsTicks        int64
	AbsMicroSeconds int64
}

type TrackEvents []*TrackEvent

func (b TrackEvents) Len() int {
	return len(b)
}

func (br TrackEvents) Swap(a, b int) {
	br[a], br[b] = br[b], br[a]
}

func (br TrackEvents) Less(a, b int) bool {
	return br[a].AbsTicks < br[b].AbsTicks
}

type playEvent struct {
	absTime int64
	sleep   time.Duration
	data    []byte
	out     drivers.Out
	trackNo int
	//str     string
}

type player []playEvent

func (p player) Swap(a, b int) {
	p[a], p[b] = p[b], p[a]
}

func (p player) Less(a, b int) bool {
	return p[a].absTime < p[b].absTime
}

func (p player) Len() int {
	return len(p)
}

// Play plays the tracks on the given out port
func (t *TracksReader) Play(out drivers.Out) error {
	if t.err != nil {
		return t.err
	}

	err := out.Open()
	if err != nil {
		return err
	}

	return t.MultiPlay(map[int]drivers.Out{-1: out})
}

// MultiPlay plays tracks to different out ports.
// If the map has an index of -1, it will be used to play all tracks that have no explicit out port.
func (t *TracksReader) MultiPlay(trackouts map[int]drivers.Out) error {
	if t.err != nil {
		return t.err
	}
	var pl player
	if len(trackouts) == 0 {
		t.err = fmt.Errorf("trackouts not set")
		return t.err
	}

	t.Do(
		func(te TrackEvent) {
			msg := te.Message
			if msg.IsPlayable() {
				var out drivers.Out

				if o, has := trackouts[te.TrackNo]; has {
					out = o
				} else {
					if def, hasDef := trackouts[-1]; hasDef {
						out = def
					} else {
						return
					}
				}

				pl = append(pl, playEvent{
					absTime: te.AbsMicroSeconds,
					data:    msg,
					out:     out,
					trackNo: te.TrackNo,
				})
			}
		},
	)

	sort.Sort(pl)

	var last time.Duration = 0

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	for i := range pl {
		last = t.play(last, pl[i])
	}

	return t.err
}

func (t *TracksReader) play(last time.Duration, p playEvent) time.Duration {
	current := (time.Microsecond * time.Duration(p.absTime))
	diff := current - last
	time.Sleep(diff)
	p.out.Send(p.data)
	return current
}

func (t *TracksReader) Do(fn func(TrackEvent)) *TracksReader {
	if t.err != nil {
		return t
	}
	tracks := t.smf.Tracks

	for no, tr := range tracks {
		if t.doTrack(no) {
			var absTicks int64
			for _, ev := range tr {
				te := TrackEvent{Event: ev, TrackNo: no}
				d := int64(ev.Delta)
				te.AbsTicks = absTicks + d
				te.AbsMicroSeconds = t.smf.TimeAt(te.AbsTicks)
				if t.filter == nil {
					fn(te)
				} else {
					msg := ev.Message
					ty := msg.Type()
					for _, f := range t.filter {
						if ty.Is(f) {
							fn(te)
						}
					}
				}
				absTicks = te.AbsTicks
			}
		}
	}

	return t
}
