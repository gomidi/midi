package mid

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"gitlab.com/gomidi/midi/smf"
)

// makeSMF makes a SMF
func makeSMF() io.Reader {
	var bf bytes.Buffer
	wr := NewSMF(&bf, 1 /* number of tracks */)
	wr.TempoBPM(160 /* beats per minute */)
	wr.Meter(4, 4) /* set the meter of 4/4 */
	wr.SetChannel(2 /* valid: 0-15 */)
	wr.NoteOn(65 /* key */, 90 /* velocity */)
	wr.Forward(1, 0, 0) // forwards the cursor to the next bar (we make a whole note)
	wr.NoteOff(65 /* key */)
	wr.EndOfTrack() // required at the end of a track
	return bytes.NewReader(bf.Bytes())
}

// roundSec is a helper to round the duration to seconds
func roundSec(d time.Duration) time.Duration {
	return time.Second * time.Duration((d.Nanoseconds() / 1000000000))
}

type example struct {
	ticks smf.MetricTicks
	bpm   float64
	start time.Time
}

// SMFHeader tracks the ticks from the SMF file
func (e *example) SMFHeader(head smf.Header) {
	// here we ignore that the timeformat could also be SMPTE
	e.ticks = head.TimeFormat.(smf.MetricTicks)
}

// Tempo tracks a tempo change
func (e *example) TempoBPM(p Position, bpm float64) {
	e.bpm = bpm
}

// NoteOn responds to note on messages
func (e *example) NoteOn(p *Position, channel, key, vel uint8) {
	fmt.Printf("[%vs] NoteOn at channel %v: key %v velocity: %v\n",
		e.calcDuration(p).Seconds(), channel, key, vel)
}

func (e *example) NoteOff(p *Position, channel, key, vel uint8) {
	fmt.Printf("[%vs] NoteOff at channel %v: key %v velocity: %v\n",
		e.calcDuration(p).Seconds(),
		channel, key, vel)
}

// a helper to calculate the duration for both live and SMF messages
func (e *example) calcDuration(p *Position) (dur time.Duration) {
	// we are in a live setting
	if p == nil {
		dur = roundSec(time.Now().Sub(e.start))
		return
	}

	// here p is not nil - that means we are reading the SMF

	// calculate the time from the timeformat of the SMF file
	// to make it easy, we ignore the possibility that tempo information may be in another track
	// that is read after this track (the SMF spec recommends to write tempo changes to the first track)
	// however, since makeSMF just creates one track, we are safe
	return roundSec(e.ticks.FractionalDuration(e.bpm, uint32(p.AbsoluteTicks)))
}

func Example() {
	// This example illustrates how the same handler can be used for live and SMF MIDI messages
	// It also illustrates how live and SMF midi can be written

	rd := NewReader(NoLogger() /* disable default logging*/)

	var ex example
	ex.bpm = 120 // default according to SMF spec

	// setup the callbacks
	rd.SMFHeader = ex.SMFHeader
	rd.Msg.Meta.TempoBPM = ex.TempoBPM
	rd.Msg.Channel.NoteOn = ex.NoteOn
	rd.Msg.Channel.NoteOff = ex.NoteOff

	// handle the smf
	fmt.Println("-- SMF data --")
	rd.ReadAllSMF(makeSMF())

	// handle the live data
	fmt.Println("-- live data --")

	// we need a pipe to read and write concurrently
	piperd, pipewr := io.Pipe()

	go func() {
		wr := NewWriter(pipewr)

		// reset the time
		ex.start = time.Now()

		wr.SetChannel(11)

		// now write some live data
		wr.NoteOn(120, 50)
		time.Sleep(time.Second * 2)
		wr.NoteOff(120)
		pipewr.Close() // close the pipe we're done writing
	}()

	for {
		if rd.ReadAllFrom(piperd) == io.EOF {
			piperd.Close() // we're done reading
			break
		}
	}

	// Output: -- SMF data --
	// [0s] NoteOn at channel 2: key 65 velocity: 90
	// [1s] NoteOff at channel 2: key 65 velocity: 0
	// -- live data --
	// [0s] NoteOn at channel 11: key 120 velocity: 50
	// [2s] NoteOff at channel 11: key 120 velocity: 0
}
