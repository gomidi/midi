package smf

import (
	//"fmt"

	"gitlab.com/gomidi/midi/v2"
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

func (t *tracksReader) Do(fn func(trackNo int, msg midi.Message, delta int64, deltamicroSec int64)) (*SMF, error) {
	tracks := t.smf.Tracks()

	ticks := t.smf.TimeFormat.(MetricTicks)
	tc := t.smf.TempoChanges()

	for no, tr := range tracks {
		var absTicks int64
		if t.doTrack(no) {
			for _, ev := range tr.Events {
				bpm := tc.TempoAt(absTicks)
				dmsec := ticks.Duration(bpm, ev.Delta).Microseconds()
				d := int64(ev.Delta)
				if t.filter == nil {
					fn(no, ev.Message(), d, dmsec)
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
							fn(no, msg, d, dmsec)
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
				absTicks += d
			}
		}
	}

	return t.smf, t.err
}
