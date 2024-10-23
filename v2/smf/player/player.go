package player

import (
	"context"
	"errors"
	"io"
	"runtime"
	"sort"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/smf"
)

/*
stolen from https://github.com/Minoxs/gomidi-player
*/

type (
	// message is an smf.Message to be played, followed by the sleeping time to the next message
	message struct {
		msg   smf.Message
		sleep time.Duration
	}

	// Player plays SMF data
	Player struct {
		mutex      sync.RWMutex
		isPlaying  bool
		ctx        context.Context
		cancelFn   context.CancelCauseFunc
		currentDur time.Duration
		totalDur   time.Duration
		currentMsg int
		messages   []message
		outPort    drivers.Out
	}
)

// New returns a Player that plays on the given output port
func New(outPort drivers.Out) *Player {
	return &Player{
		ctx:      UnavailableContext(),
		cancelFn: func(cause error) {},
		outPort:  outPort,
	}
}

// stop will signal the player to stop playing
func (p *Player) stop(cause error) error {
	if !p.isPlaying {
		return ErrIsStopped
	}
	p.cancelFn(cause)
	return nil
}

// SetSMF takes in a SMF and creates units that can be played.
func (p *Player) SetSMF(smfdata io.Reader, tracks ...int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.isPlaying {
		return ErrIsPlaying
	}

	p.currentMsg = 0
	p.currentDur = 0
	p.messages = nil

	var events smfReader
	err := events.read(smfdata, tracks...)
	if err != nil {
		return err
	}

	p.messages, p.totalDur = events.getMessages()
	return nil
}

// Start starts playing. It is non-blocking. Call wait to wait until it is finished.
func (p *Player) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.messages == nil {
		return ErrNoSMFData
	}

	if p.isPlaying {
		return ErrIsPlaying
	}

	p.isPlaying = true
	p.ctx, p.cancelFn = context.WithCancelCause(context.Background())

	go p.playOn(p.outPort)
	return nil
}

// Stop stops the playing.
func (p *Player) Stop() (err error) {
	return p.stop(errStopped)
}

// Pause pauses the playing.
func (p *Player) Pause() (err error) {
	return p.stop(errPaused)
}

// Wait will block until the player finishes playing.
func (p *Player) Wait() {
	for p.isPlaying {
		time.Sleep(50 * time.Millisecond)
	}
}

// IsPlaying returns wether the player is playing
func (p *Player) IsPlaying() bool {
	return p.isPlaying
}

// Duration is the total duration of the song.
func (p *Player) Duration() time.Duration {
	return p.totalDur
}

// Current is the current time on the song.
func (p *Player) Current() time.Duration {
	return p.currentDur
}

// Remaining is the remaining duration of the smf data.
func (p *Player) Remaining() time.Duration {
	return p.totalDur - p.currentDur
}

type smfReader struct {
	trackEvents []smf.TrackEvent
}

// read reads SMF data and parses all the track events.
func (e *smfReader) read(smfdata io.Reader, tracks ...int) (err error) {
	e.trackEvents = make([]smf.TrackEvent, 0, 100)
	return WrapOnError(
		smf.ReadTracksFrom(smfdata, tracks...).Do(e.readEvent).Error(),
		ErrInvalidSMF,
	)
}

// readEvent reads a single event from the track.
func (e *smfReader) readEvent(event smf.TrackEvent) {
	if event.Message.IsPlayable() {
		e.trackEvents = append(e.trackEvents, event)
	}
}

// getMessages parses the track events and returns the playable units.
func (e *smfReader) getMessages() (messages []message, totalDur time.Duration) {
	e.sortTrackEvents()
	messages = make([]message, len(e.trackEvents))
	totalDur = 0
	for i := 0; i < len(messages); i++ {
		var event = e.trackEvents[i]
		messages[i] = message{
			msg:   event.Message,
			sleep: time.Microsecond * time.Duration(event.AbsMicroSeconds-totalDur.Microseconds()),
		}
		totalDur += messages[i].sleep
	}
	return
}

// sortTrackEvents makes sure the song is ordered by time.
func (e *smfReader) sortTrackEvents() {
	sort.SliceStable(
		e.trackEvents, func(i, j int) bool {
			return e.trackEvents[i].AbsMicroSeconds < e.trackEvents[j].AbsMicroSeconds
		},
	)
}

// playOn will play the current song in the given out port
func (p *Player) playOn(out drivers.Out) {
	defer p.cleanupAfterPlaying()

	// Drivers may invoke CGO
	// Makes sure thread is locked to avoid weird errors
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	//	defer IgnoreError(out.Close)

	// Creates timer to properly time sound writes
	var sleep = time.NewTimer(0)
	defer sleep.Stop()

	// Makes sure channel is drained
	<-sleep.C

	// play all messages
	for i, m := range p.messages[p.currentMsg:] {
		p.currentDur += m.sleep
		p.currentMsg = i
		if m.sleep > 0 {
			sleep.Reset(m.sleep)
			select {
			case <-sleep.C:
				break
			case <-p.ctx.Done():
				return
			}
		}
		_ = out.Send(m.msg)
	}
}

// cleanupAfterPlaying does the cleaning up after playOn ends
func (p *Player) cleanupAfterPlaying() {
	_ = p.stop(errDone)
	var err = context.Cause(p.ctx)

	switch {
	case errors.Is(err, errDone):
		fallthrough
	case errors.Is(err, errStopped):
		p.currentMsg = 0
		p.currentDur = 0
	}

	p.isPlaying = false
}

func Play(out drivers.Out, smfdata io.Reader, tracks ...int) (*Player, error) {
	player := New(out)
	err := player.SetSMF(smfdata, tracks...)
	if err != nil {
		return nil, err
	}

	err = player.Start()
	if err != nil {
		return nil, err
	}

	player.Wait()
	return player, nil
}
