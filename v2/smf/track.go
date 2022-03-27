package smf

import (
	"fmt"
	"sort"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

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
	t.Events = append(t.Events, Event{Delta: deltaticks, Data: EOT.Data})
	//fmt.Printf("appending bytes: % X\n", EOT.Data)
	t.Closed = true
}

func (t *Track) Add(deltaticks uint32, msgs ...midi.Message) {
	if t.Closed {
		return
	}
	for _, msg := range msgs {
		ev := Event{Delta: deltaticks, Data: msg.Bytes()}
		//fmt.Printf("appending bytes: % X, evtype: %s\n", ev.Data, ev.MsgType())
		t.Events = append(t.Events, ev)
		deltaticks = 0
	}
}

func (t *Track) SendTo(resolution MetricTicks, tc TempoChanges, receiver midi.Receiver) {
	var absDelta int64

	for _, ev := range t.Events {
		absDelta += int64(ev.Delta)
		if m, ok := ev.Message().(midi.Msg); ok {
			ms := int32(resolution.Duration(tc.TempoAt(absDelta), ev.Delta).Microseconds() * 100)
			receiver.Receive(m, ms)
		}
	}
}

type tracksReader struct {
	smf    *SMF
	tracks map[int]bool
	filter []midi.MsgType
	err    error
}

func (t *tracksReader) Error() error {
	return t.err
}

func (t *tracksReader) SMF() *SMF {
	return t.smf
}

func (t *tracksReader) doTrack(tr int) bool {
	if len(t.tracks) == 0 {
		return true
	}

	return t.tracks[tr]
}

func ReadTracks(filepath string, tracks ...int) *tracksReader {
	t := &tracksReader{}
	t.tracks = map[int]bool{}
	for _, tr := range tracks {
		t.tracks[tr] = true
	}
	t.smf, t.err = ReadFile(filepath)
	return t
}

func (t *tracksReader) Only(mtypes ...midi.MsgType) *tracksReader {
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
	//bytes   []byte
	out     drivers.Out
	trackNo int
	str     string
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
func (t *tracksReader) Play(out drivers.Out) *tracksReader {
	return t.MultiPlay(map[int]drivers.Out{-1: out})
}

// MultiPlay plays tracks to different out ports.
// If the map has an index of -1, it will be used to play all tracks that have no explicit out port.
func (t *tracksReader) MultiPlay(trackouts map[int]drivers.Out) *tracksReader {
	var pl player
	if len(trackouts) == 0 {
		t.err = fmt.Errorf("trackouts not set")
		return t
	}

	t.Do(
		func(te TrackEvent) {
			msg := te.Message()
			if mm, ok := msg.(midi.Msg); ok {
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
					data:    mm.Data,
					out:     out,
					trackNo: te.TrackNo,
					str:     msg.String(),
				})
			}
		},
	)

	sort.Sort(pl)

	var last time.Duration = 0

	for i, _ := range pl {
		last = t.play(last, pl[i])
	}

	return t
}

func (t *tracksReader) play(last time.Duration, p playEvent) time.Duration {
	current := (time.Microsecond * time.Duration(p.absTime))
	diff := current - last
	//fmt.Printf("sleeping %s\n", diff)
	time.Sleep(diff)
	//fmt.Printf("[%v] %q % X\n", p.trackNo, p.str, p.data)
	//fmt.Printf("[%v] %q\n", p.trackNo, p.str)
	p.out.Send(p.data)
	return current
}

func (t *tracksReader) Do(fn func(TrackEvent)) *tracksReader {
	tracks := t.smf.Tracks()

	//	ticks := t.smf.TimeFormat.(MetricTicks)
	//tc := t.smf.TempoChanges()

	for no, tr := range tracks {
		if t.doTrack(no) {
			var absTicks int64
			for _, ev := range tr.Events {
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
					msg := ev.Message()
					ty := msg.Type()
					for _, f := range t.filter {
						if midi.Is(f, ty) {
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
