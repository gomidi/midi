package rtmididrv

import (
	"fmt"
	"math"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi"
)

type in struct {
	number int
	sync.RWMutex
	listenerSet bool
	driver      *Driver
	name        string
	midiIn      rtmidi.MIDIIn
}

// IsOpen returns wether the MIDI in port is open
func (i *in) IsOpen() (open bool) {
	i.RLock()
	open = i.midiIn != nil
	i.RUnlock()
	return
}

// String returns the name of the MIDI in port.
func (i *in) String() string {
	return i.name
}

// Underlying returns the underlying rtmidi.MIDIIn. Use it with type casting:
//   rtIn := i.Underlying().(rtmidi.MIDIIn)
func (i *in) Underlying() interface{} {
	return i.midiIn
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

	i.StopListening()
	i.Lock()
	err = i.midiIn.Close()
	i.midiIn = nil
	i.Unlock()
	return
}

// Open opens the MIDI in port
func (i *in) Open() (err error) {
	if i.IsOpen() {
		return nil
	}

	i.Lock()

	i.midiIn, err = rtmidi.NewMIDIInDefault()
	if err != nil {
		i.midiIn = nil
		i.Unlock()
		return fmt.Errorf("can't open default MIDI in: %v", err)
	}

	err = i.midiIn.OpenPort(i.number, "")
	i.Unlock()

	if err != nil {
		i.Close()
		return fmt.Errorf("can't open MIDI in port %v (%s): %v", i.number, i, err)
	}

	i.driver.Lock()
	i.midiIn.IgnoreTypes(i.driver.ignoreSysex, i.driver.ignoreTimeCode, i.driver.ignoreActiveSense)
	i.driver.opened = append(i.driver.opened, i)
	i.driver.Unlock()

	return nil
}

func newIn(driver *Driver, number int, name string) drivers.In {
	return &in{driver: driver, number: number, name: name}
}

func (i *in) StartListening(callback func(data []byte, deltadecimilliseconds int32)) error {
	if !i.IsOpen() {
		//fmt.Printf("post closed\n")
		return drivers.ErrPortClosed
	}

	i.RLock()
	if i.listenerSet {
		i.RUnlock()
		return fmt.Errorf("listener already set")
	}
	i.RUnlock()
	//fmt.Println("pre lock")
	i.Lock()
	i.listenerSet = true
	i.driver.Lock()
	i.midiIn.IgnoreTypes(i.driver.ignoreSysex, i.driver.ignoreTimeCode, i.driver.ignoreActiveSense)
	i.driver.Unlock()
	i.Unlock()

	var tsdecimilliseconds int32

	// since i.midiIn.SetCallback is blocking on success, there is no meaningful way to get an error
	// and set the callback non blocking
	go i.midiIn.SetCallback(func(_ rtmidi.MIDIIn, bt []byte, deltaSeconds float64) {
		// convert to milliseconds (10^-5)
		tsdecimilliseconds += int32(math.Round(deltaSeconds * 1000))
		callback(bt, tsdecimilliseconds)
	})

	time.Sleep(time.Millisecond * 10)

	return nil
}

// StopListening cancels the listening
func (i *in) StopListening() (err error) {
	if !i.IsOpen() {
		return drivers.ErrPortClosed
	}
	i.Lock()
	if i.listenerSet {
		i.listenerSet = false
		err = i.midiIn.CancelCallback()
		if err != nil {
			err = fmt.Errorf("can't stop listening on MIDI in port %v (%s): %v", i.number, i, err)
		}
	}
	i.Unlock()
	return
}
