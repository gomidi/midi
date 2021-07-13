package midicatdrv

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midicat/lib"
)

type in struct {
	number int
	sync.RWMutex
	driver              *Driver
	name                string
	shouldStopListening chan bool
	didStopListening    chan bool
	shouldKill          chan bool
	wasKilled           chan bool
	hasProc             bool
	listener            func(data midi.Message, absMicrosecs int64)
}

func (o *in) fireCmd() error {
	o.Lock()
	if o.hasProc {
		o.Unlock()
		return fmt.Errorf("already running")
	}
	o.shouldStopListening = make(chan bool, 1)
	o.didStopListening = make(chan bool, 1)
	o.shouldKill = make(chan bool, 1)
	o.wasKilled = make(chan bool, 1)
	o.hasProc = true
	cmd := midiCatInCmd(o.number)
	rd, wr := io.Pipe()
	cmd.Stdout = wr
	err := cmd.Start()
	if err != nil {
		o.Lock()
		o.hasProc = false
		o.Unlock()
		return err
	}
	o.Unlock()
	go func() {
		for {
			data, err := lib.ReadAndConvert(rd)
			if err != nil {
				return
			}
			o.RLock()
			if !o.hasProc {
				o.RUnlock()
				return
			}

			if o.listener != nil {
				o.listener(midi.NewMessage(data), -1)
			}
			o.RUnlock()
			runtime.Gosched()
		}
	}()

	go func(shouldStopListening <-chan bool, didStopListening chan<- bool, shouldKill <-chan bool, wasKilled chan<- bool) {
		defer rd.Close()
		defer wr.Close()

		for {
			select {
			case <-shouldKill:
				if cmd.Process != nil {
					/*
						                                        rd.Close()
											wr.Close()
					*/
					cmd.Process.Kill()
				}
				o.Lock()
				o.hasProc = false
				o.Unlock()
				wasKilled <- true
				return
			case <-shouldStopListening:
				o.Lock()
				o.listener = nil
				o.Unlock()
				didStopListening <- true
			default:
				runtime.Gosched()
			}
		}
	}(o.shouldStopListening, o.didStopListening, o.shouldKill, o.wasKilled)

	return nil
}

// IsOpen returns wether the MIDI in port is open
func (o *in) IsOpen() (open bool) {
	o.RLock()
	open = o.hasProc
	o.RUnlock()
	return
}

// String returns the name of the MIDI in port.
func (i *in) String() string {
	return i.name
}

// Underlying returns the underlying driver. Here returns nil.
func (i *in) Underlying() interface{} {
	return nil
}

// Number returns the number of the MIDI in port.
// Note that with rtmidi, out and in ports are counted separately.
// That means there might exists out ports and an in ports that share the same number.
func (i *in) Number() int {
	return i.number
}

// Close closes the MIDI in port, after it has stopped listening.
func (i *in) Close() (err error) {
	if !i.IsOpen() {
		return nil
	}

	//i.shouldStopReading
	go func() {
		i.shouldStopListening <- true
	}()
	<-i.didStopListening

	i.shouldKill <- true
	<-i.wasKilled
	return
}

// Open opens the MIDI in port
func (i *in) Open() (err error) {
	if i.IsOpen() {
		return nil
	}

	err = i.fireCmd()
	if err != nil {
		i.Close()
		return fmt.Errorf("can't open MIDI in port %v (%s): %v", i.number, i, err)
	}

	i.driver.Lock()
	i.driver.opened = append(i.driver.opened, i)
	i.driver.Unlock()

	return nil
}

func newIn(driver *Driver, number int, name string) midi.In {
	return &in{driver: driver, number: number, name: name}
}

// SendTo makes the listener listen to the in port
func (i *in) SendTo(recv midi.Receiver) (err error) {
	if !i.IsOpen() {
		return midi.ErrPortClosed
	}

	i.RLock()
	if i.listener != nil {
		i.RUnlock()
		return fmt.Errorf("listener already set")
	}
	i.RUnlock()
	i.Lock()
	i.listener = recv.Receive
	i.Unlock()

	return nil
}

// StopListening cancels the listening
func (i *in) StopListening() (err error) {
	if !i.IsOpen() {
		return midi.ErrPortClosed
	}

	i.shouldStopListening <- true
	<-i.didStopListening
	return
}
