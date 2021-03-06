package portmididrv

import (
	"fmt"
	"sync"
	"time"

	"github.com/rakyll/portmidi"
	"gitlab.com/gomidi/midi/v2"
)

func newIn(driver *driver, id portmidi.DeviceID, name string) midi.In {
	return &in{driver: driver, id: id, name: name}
}

type in struct {
	id     portmidi.DeviceID
	stream *portmidi.Stream
	name   string

	driver *driver

	lastTimestamp portmidi.Timestamp
	mx            sync.Mutex
	stopped       bool
}

// IsOpen returns wether the MIDI in port is open.
func (i *in) IsOpen() bool {
	return i.stream != nil
}

// Underlying returns the underlying *portmidi.Stream. It will be nil, if the port is closed.
// Use it with type casting:
//   portIn := i.Underlying().(*portmidi.Stream)
func (i *in) Underlying() interface{} {
	return i.stream
}

// Number returns the number of the MIDI in port.
// Note that with portmidi, out and in ports are counted together.
// That means there should not be an out port with the same number as an in port.
func (i *in) Number() int {
	return int(i.id)
}

// String returns the name of the MIDI in port.
func (i *in) String() string {
	return i.name
}

// Close closes the MIDI in port
func (i *in) Close() error {
	if i.stream == nil {
		return nil
	}

	err := i.StopListening()
	if err != nil {
		panic("unreachable")
	}

	err = i.stream.Close()
	if err != nil {
		return fmt.Errorf("can't close MIDI in %v (%s): %v", i.Number(), i, err)
	}
	i.stream = nil
	return nil
}

// Open opens the MIDI in port
func (i *in) Open() (err error) {
	if i.stream != nil {
		return nil
	}
	i.stream, err = portmidi.NewInputStream(i.id, i.driver.buffersizeIn)
	if err != nil {
		i.stream = nil
		return fmt.Errorf("can't open MIDI in port %v (%s): %v", i.Number(), i, err)
	}
	i.driver.opened = append(i.driver.opened, i)
	return nil
}

// StopListening cancels the listening
func (i *in) StopListening() error {
	i.mx.Lock()
	i.stopped = true
	i.mx.Unlock()
	return nil
}

// read is an internal helper function
func (i *in) read(cb func([]byte, int64)) error {
	events, err := i.stream.Read(i.driver.buffersizeRead)

	if err != nil {
		return err
	}

	for _, ev := range events {
		var b = make([]byte, 3)
		b[0] = byte(ev.Status)
		b[1] = byte(ev.Data1)
		b[2] = byte(ev.Data2)
		// ev.Timestamp is in Milliseconds
		// we want deltaMicroseconds as int64
		cb(b, int64(ev.Timestamp-i.lastTimestamp)*1000)
	}

	return nil
}

// SendTo
func (i *in) SendTo(recv midi.Receiver) error {
	i.lastTimestamp = portmidi.Time()
	for i.stopped == false {
		has, _ := i.stream.Poll()
		if has {
			i.read(recv.Receive)
		}
		time.Sleep(i.driver.sleepingTime)
	}
	return nil
}


/*
// SetListener sets the listener
func (i *in) SetListener(listener func(data []byte, deltaMicroseconds int64)) error {
	i.lastTimestamp = portmidi.Time()
	for i.stopped == false {
		has, _ := i.stream.Poll()
		if has {
			i.read(listener)
		}
		time.Sleep(i.driver.sleepingTime)
	}
	return nil
}
*/
