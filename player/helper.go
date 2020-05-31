package player

import (
	"time"

	"gitlab.com/gomidi/midi"
)

type message struct {
	pos   uint64
	msg   midi.Message
	track int16
}

type tempo struct {
	pos      uint64
	tempoBPM float64
}

type tempi []tempo
type messages []message

func (t tempi) Len() int {
	return len(t)
}

func (t tempi) Swap(a, b int) {
	t[a], t[b] = t[b], t[a]
}

func (t tempi) Less(a, b int) bool {
	return t[a].pos < t[b].pos
}

func (t messages) Len() int {
	return len(t)
}

func (t messages) Swap(a, b int) {
	t[a], t[b] = t[b], t[a]
}

func (t messages) Less(a, b int) bool {
	if t[a].pos == t[b].pos {
		// TODO for same track, send controll messages before note messages and note off before note ons
		return t[a].track < t[b].track
	}
	return t[a].pos < t[b].pos
}

type durCalc struct {
	player       *Player
	lastTempoIdx int
	lastMsgIdx   int
}

func (c *durCalc) findTempoChanges(from, to uint64) (tc []tempo) {
	for i, t := range c.player.tempi {
		if t.pos >= from && t.pos < to {
			//fmt.Printf("tempochange at: %v to: %v\n", t.pos, t.tempoBPM)
			tc = append(tc, t)
			c.lastTempoIdx = i
		}
	}
	return
}

func (c *durCalc) getDur(lastTempo tempo, from, to uint64) (d time.Duration) {
	tc := c.findTempoChanges(from, to)
	last := from

	for _, t := range tc {
		diff := t.pos - last
		d += c.player.duration(lastTempo.tempoBPM, uint32(diff))
		lastTempo = t
		last = t.pos
	}

	diff := to - last
	if diff > 0 {
		d += c.player.duration(lastTempo.tempoBPM, uint32(diff))
	}

	return
}

// if there is no message left, msg is nil
func (c *durCalc) nextMessage() (d time.Duration, msg midi.Message, track int16) {
	if c.lastMsgIdx+1 > len(c.player.messages)-1 {
		//fmt.Printf("no messages left\n")
		return 0, nil, -1
	}

	nextMsg := c.player.messages[c.lastMsgIdx+1]
	var lastMsgPos uint64 = 0
	if c.lastMsgIdx >= 0 {
		lastMsg := c.player.messages[c.lastMsgIdx]
		lastMsgPos = lastMsg.pos
	}

	c.lastMsgIdx++

	var lastTempo = tempo{pos: 0, tempoBPM: 120}
	if c.lastTempoIdx >= 0 {
		lastTempo = c.player.tempi[c.lastTempoIdx]
	}

	d = c.getDur(lastTempo, lastMsgPos, nextMsg.pos)
	return d, nextMsg.msg, nextMsg.track
}
