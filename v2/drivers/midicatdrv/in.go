package midicatdrv

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"gitlab.com/gomidi/midi/v2"
)

/*
TODO

opening means:
  - create the process of an midireader
  - and start reading from it (throw away the bytes)

setListener means:
  - don't throw away the bytes but pass them to the listener instead

stopListening means:
  - throw the bytes away, again

close mean:
  - kill the process

IMPORTANT:
we need to keep track of all open ports inside the driver, since for each port there could be
a process and we want to make sure that no port is opened twice.
All open ports (processes) must be closed when closing the driver.
Since each process must run in its own gorouting, we need channes to communicate:
- should stop reading
- has stopped reading
- should kill
- was killed
*/

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
	listener            func(data []byte, deltaMicroseconds int64)
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
	o.Unlock()
	go func(shouldStopListening <-chan bool, didStopListening chan<- bool, shouldKill <-chan bool, wasKilled chan<- bool) {
		cmd := midiCatInCmd(o.number)
		//cmd := midiCatCmd(fmt.Sprintf("in --index=%v --name='%s'", o.number, o.name))
		rd, wr := io.Pipe()
		cmd.Stdout = wr
		err := cmd.Start()
		if err != nil {
			o.Lock()
			o.hasProc = false
			o.Unlock()
			return
		}
		for {
			select {
			case <-shouldKill:
				if cmd.Process != nil {
					rd.Close()
					wr.Close()
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
				for {
					var data = make([]byte, 3)
					_, err := rd.Read(data)
					if err != nil {
						o.Lock()
						o.hasProc = false
						o.Unlock()
						return
					}
					o.Lock()
					if o.listener != nil {
						o.listener(data, -1)
					}
					o.Unlock()
					runtime.Gosched()
				}
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
	i.shouldStopListening <- true
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

/*
// SetListener makes the listener listen to the in port
func (i *in) SetListener(listener func(data []byte, deltaMicroseconds int64)) (err error) {
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
	i.listener = listener
	i.Unlock()

	return nil
}
*/

// StopListening cancels the listening
func (i *in) StopListening() (err error) {
	if !i.IsOpen() {
		return midi.ErrPortClosed
	}

	i.shouldStopListening <- true
	<-i.didStopListening
	return
}
