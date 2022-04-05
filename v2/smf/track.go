package smf

import (
	"fmt"
	"io"
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

/*
func NewTrack() (t Track) {
	return
}
*/

/*
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
*/

func (t *Track) Close(deltaticks uint32) {
	if t.IsClosed() {
		return
	}
	*t = append(*t, Event{Delta: deltaticks, Message: EOT})
	//fmt.Printf("appending bytes: % X\n", EOT.Data)
	//t.Closed = true
}

func (t *Track) Add(deltaticks uint32, msgs ...[]byte) {
	if t.IsClosed() {
		return
	}
	for _, msg := range msgs {
		ev := Event{Delta: deltaticks, Message: msg}
		//fmt.Printf("appending bytes: % X, evtype: %s\n", ev.Data, ev.MsgType())
		*t = append(*t, ev)
		deltaticks = 0
	}
}

func (t *Track) RecordFrom(portno int, ticks MetricTicks, bpm float64) (stop func(), err error) {
	t.Add(0, MetaTempo(bpm))
	var absmillisec int32
	//ticks := file.TimeFormat.(smf.MetricTicks)
	return midi.ListenTo(portno, func(msg midi.Message, absms int32) {
		deltams := absms - absmillisec
		absmillisec = absms
		//fmt.Printf("[%v] %s\n", deltams, msg.String())
		delta := ticks.Ticks(bpm, time.Duration(deltams)*time.Millisecond)
		t.Add(delta, msg)
	})
}

func (t *Track) SendTo(resolution MetricTicks, tc TempoChanges, receiver func(m midi.Message, timestampms int32)) {
	var absDelta int64

	for _, ev := range *t {
		absDelta += int64(ev.Delta)
		if Message(ev.Message).IsPlayable() {
			//		if m, ok := ev.Message().Type() <  .(midi.Msg); ok {
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
	if _, ok := t.smf.TimeFormat.(MetricTicks); !ok {
		t.err = fmt.Errorf("SMF time format is not metric ticks, but %s (currently not supported)", t.smf.TimeFormat.String())
		return nil
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
	if _, ok := t.smf.TimeFormat.(MetricTicks); !ok {
		t.err = fmt.Errorf("SMF time format is not metric ticks, but %s (currently not supported)", t.smf.TimeFormat.String())
		return nil
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
func (t *TracksReader) Play(out int) error {
	o, err := drivers.OutByNumber(out)
	if err != nil {
		return err
	}
	err = o.Open()
	if err != nil {
		return err
	}

	return t.MultiPlay(map[int]drivers.Out{-1: o})
}

// MultiPlay plays tracks to different out ports.
// If the map has an index of -1, it will be used to play all tracks that have no explicit out port.
func (t *TracksReader) MultiPlay(trackouts map[int]drivers.Out) error {
	var pl player
	if len(trackouts) == 0 {
		t.err = fmt.Errorf("trackouts not set")
		return t.err
	}

	t.Do(
		func(te TrackEvent) {
			msg := te.Message
			//ty := msg.Type
			if msg.IsPlayable() {
				//if mm, ok := msg.(midi.Msg); ok {
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
					//str:     msg.String(),
				})
			}
		},
	)

	sort.Sort(pl)

	var last time.Duration = 0

	for i, _ := range pl {
		last = t.play(last, pl[i])
	}

	return t.err
}

func (t *TracksReader) play(last time.Duration, p playEvent) time.Duration {
	current := (time.Microsecond * time.Duration(p.absTime))
	diff := current - last
	//fmt.Printf("sleeping %s\n", diff)
	time.Sleep(diff)
	//fmt.Printf("[%v] %q % X\n", p.trackNo, p.str, p.data)
	//fmt.Printf("[%v] %q\n", p.trackNo, p.str)
	p.out.Send(p.data)
	return current
}

func (t *TracksReader) Do(fn func(TrackEvent)) *TracksReader {
	tracks := t.smf.Tracks

	//	ticks := t.smf.TimeFormat.(MetricTicks)
	//tc := t.smf.TempoChanges()

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
					/*
						if ev.MsgType().IsOneOf(t.filter...) {
							fn(no, ev.Message(), d, dmsec)
						}
					*/
					msg := ev.Message
					ty := msg.Type()
					for _, f := range t.filter {
						//fmt.Printf("%s [%s] %s [%s]\n", f, f.Kind().String(), ty, ty.Kind().String())
						if ty.Is(f) {
							//if Is(f, ty) {
							//fn(no, msg, d, dmsec)
							fn(te)
						}
					}
					/*
						kind := ev.Message().Kind()
						for _, f := range t.filter {
							if kind == f {
								fn(no, ev.Message(), d, dmsec)
							}
						}
					*/
				}
				absTicks = te.AbsTicks
			}
		}
	}

	return t
}
