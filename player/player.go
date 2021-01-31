package player

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/midi/midimessage/sysex"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/smf"
	"gitlab.com/gomidi/midi/smf/smfreader"
	"gitlab.com/gomidi/midi/writer"
)

// Player is a SMF (MIDI file) player.
type Player struct {
	messages messages
	tempi    tempi
	metric   smf.MetricTicks
	tracks   map[int16]string
	durCalc  *durCalc
}

func (p *Player) duration(fractionalBPM float64, ticks uint32) time.Duration {
	return p.metric.FractionalDuration(fractionalBPM, ticks)
}

// SMF creates a new Player for the given file. If you need to pass an opened file, use New.
func SMF(file string, options ...smfreader.Option) (*Player, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer func() {
		f.Close()
	}()

	return New(f, options...)
}

// New creates a new Player, based on the given io.Reader
func New(rd io.Reader, options ...smfreader.Option) (*Player, error) {
	p := &Player{
		//tracks: map[int16]messageTrack{},
		//out: out,
		tracks: map[int16]string{},
	}

	r := reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			switch v := msg.(type) {
			case meta.Tempo:
				//fmt.Printf("tempochange at: %v to: %v\n", pos.AbsoluteTicks, v.FractionalBPM())
				p.tempi = append(p.tempi, tempo{pos.AbsoluteTicks, v.FractionalBPM()})
			case channel.Message:
				p.messages = append(p.messages, message{pos: pos.AbsoluteTicks, msg: msg, track: pos.Track})
			case sysex.Message:
				p.messages = append(p.messages, message{pos: pos.AbsoluteTicks, msg: msg, track: pos.Track})
			case meta.Instrument:
				p.tracks[pos.Track] = p.tracks[pos.Track] + " " + v.Text()
			case meta.Program:
				p.tracks[pos.Track] = p.tracks[pos.Track] + " " + v.Text()
			case meta.TrackSequenceName:
				p.tracks[pos.Track] = p.tracks[pos.Track] + " " + v.Text()
			default:
				// ignore
			}
		}))

	for tr, nm := range p.tracks {
		p.tracks[tr] = strings.TrimSpace(nm)
	}

	err := reader.ReadSMF(r, rd, options...)

	if err != nil {
		return nil, err
	}

	header := r.Header()
	if ticks, ok := header.TimeFormat.(smf.MetricTicks); ok {
		p.metric = ticks
	} else {
		return nil, fmt.Errorf("can only play metric tricks timeformat, but got: %s", header.TimeFormat.String())
	}

	sort.Sort(p.messages)
	sort.Sort(p.tempi)

	return p, nil
}

func (p *Player) newDurCalc() {
	p.durCalc = &durCalc{
		player:       p,
		lastTempoIdx: -1,
		lastMsgIdx:   -1,
	}
}

// GetMessages calls the given callback for each midi message and passes the sleeping time
// until that message is to be played as well as the track number and the message itself.
// The callback is called until there are no midi messages left.
func (p *Player) GetMessages(callback func(wait time.Duration, m midi.Message, track int16)) {
	p.newDurCalc()

	for {
		d, msg, tr := p.durCalc.nextMessage()
		callback(d, msg, tr)

		if msg == nil {
			return
		}
	}
}

// PlayAll is a shortcut for PlayAllTo when using a MIDI port as output.
// It stops hanging notes before returning and when stopping.
// If you need to pass a specific midi.Writer, use PlayAllTo.
func (p *Player) PlayAll(out midi.Out, stop <-chan bool, finished chan<- bool) {
	wr := writer.New(out)
	p.playAllTo(wr, stop, finished)
	return
}

// PlayAllTo plays all tracks to the given midi.Writer until there are no messages left or
// a boolean is inserted into the given stop channel.
// When the function returns, the playing has finished.
// It stops hanging notes, when finished / stopped.
// If you need to play to different writers per track or channel, use GetMessages and define your own playing style.
func (p *Player) PlayAllTo(wr midi.Writer, stop <-chan bool, finished chan<- bool) {
	if w, ok := wr.(*writer.Writer); ok {
		p.playAllTo(w, stop, finished)
		return
	}
	p.playAllTo(writer.Wrap(wr), stop, finished)
}

// playAllTo plays all tracks to the given midi.Writer until there are no messages left or
// a boolean is inserted into the given stop channel.
// It stops hanging notes, when finished / stopped.
func (p *Player) playAllTo(wr *writer.Writer, stop <-chan bool, finished chan<- bool) {
	// stop hanging notes
	wr.Silence(-1, true)

	// give it some time
	time.Sleep(10 * time.Millisecond)

	p.newDurCalc()

	d, msg, _ := p.durCalc.nextMessage()

	if msg == nil {
		finished <- true
		return
	}

	go func() {
		for {
			select {
			case <-stop:
				//fmt.Printf("received stop\n")
				wr.Silence(-1, true)
				finished <- true
				return
			default:
				//fmt.Printf("waiting %v\n", d)
				time.Sleep(d)
				wr.Write(msg)
				d, msg, _ = p.durCalc.nextMessage()
				//fmt.Printf("next at %v\n", d)
				if msg == nil {
					wr.Silence(-1, true)
					finished <- true
					return
				}
				runtime.Gosched()
			}
		}

	}()

	return
}

/*
// is too slow!!
func (p *Player) PlayAll2(wr midi.Writer, stop chan bool) (finished chan bool) {
	c := &duractionCalculator{
		player:       p,
		lastTempoIdx: -1,
		lastMsgIdx:   -1,
	}

	finished = make(chan bool)

	d, msg, _ := c.NextMessage()

	if msg == nil {
		finished <- true
		return
	}

	after := time.After(d)

	go func() {
		for {
			select {
			case <-stop:
				//fmt.Printf("received stop\n")
				after = nil
				finished <- true
				return
			case <-after:
				wr.Write(msg)
				d, msg, _ = c.NextMessage()
				if msg == nil {
					finished <- true
					return
				}
				after = time.After(d)
			default:
			}
		}

	}()
	return
}
*/
