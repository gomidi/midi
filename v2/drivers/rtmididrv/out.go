//go:build !js
// +build !js

package rtmididrv

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi"
)

func newOut(driver *Driver, number int, name string) drivers.Out {
	o := &out{driver: driver, number: number, name: name}
	return o
}

type out struct {
	number int
	//sync.RWMutex
	driver  *Driver
	name    string
	midiOut rtmidi.MIDIOut
}

// IsOpen returns wether the port is open
func (o *out) IsOpen() (open bool) {
	//	o.RLock()
	open = o.midiOut != nil
	//	o.RUnlock()
	return
}

// Send writes a MIDI sysex message to the outut port
func (o *out) SendSysEx(data []byte) error {
	//fmt.Printf("try to send sysex\n")

	if o.midiOut == nil {
		//o.RUnlock()
		return drivers.ErrPortClosed
	}
	//o.mx.RUnlock()

	// since we always open the outputstream with a latency of 0
	// the timestamp is ignored
	//var ts portmidi.Timestamp // or portmidi.Time()

	//o.mx.Lock()
	//	defer o.mx.Unlock()
	//fmt.Printf("sending sysex % X\n", data)
	//err := o.stream.WriteSysExBytes(ts, data)
	err := o.midiOut.SendMessage(data)
	if err != nil {
		return fmt.Errorf("could not send sysex message to MIDI out %v (%s): %v", o.Number(), o, err)
	}
	return nil
}

func (o *out) Send(b []byte) error {
	if o.midiOut == nil {
		//o.RUnlock()
		return drivers.ErrPortClosed
	}
	//	o.RUnlock()

	//fmt.Printf("send % X\n", m.Data)
	/*
		var bt []byte

		switch {
		case b[2] == 0 && b[1] == 0:
			bt = []byte{b[0]}
			//	case b[2] == 0:
		//	bt = []byte{b[0], b[1]}
		default:
			bt = []byte{b[0], b[1], b[2]}
		}

		//bt := []byte{b[0], b[1], b[2]}
		err := o.midiOut.SendMessage(bt)
	*/
	err := o.midiOut.SendMessage(b)
	if err != nil {
		return fmt.Errorf("could not send message to MIDI out %v (%s): %v", o.number, o, err)
	}
	return nil
}

/*
// Send writes a MIDI message to the MIDI output port
// If the output port is closed, it returns midi.ErrClosed
func (o *out) send(bt []byte) error {
	//o.RLock()
	o.Lock()
	defer o.Unlock()
	if o.midiOut == nil {
		//o.RUnlock()
		return drivers.ErrPortClosed
	}
	//	o.RUnlock()

	//fmt.Printf("send % X\n", m.Data)
	err := o.midiOut.SendMessage(bt)
	if err != nil {
		return fmt.Errorf("could not send message to MIDI out %v (%s): %v", o.number, o, err)
	}
	return nil
}
*/

// Underlying returns the underlying rtmidi.MIDIOut. Use it with type casting:
//
//	rtOut := o.Underlying().(rtmidi.MIDIOut)
func (o *out) Underlying() interface{} {
	return o.midiOut
}

// Number returns the number of the MIDI out port.
// Note that with rtmidi, out and in ports are counted separately.
// That means there might exists out ports and an in ports that share the same number
func (o *out) Number() int {
	return o.number
}

// String returns the name of the MIDI out port.
func (o *out) String() string {
	return o.name
}

// Close closes the MIDI out port
func (o *out) Close() (err error) {
	if !o.IsOpen() {
		return nil
	}
	//o.Lock()
	//defer o.Unlock()

	err = o.midiOut.Close()
	o.midiOut = nil

	if err != nil {
		err = fmt.Errorf("can't close MIDI out %v (%s): %v", o.number, o, err)
	}

	return
}

// Open opens the MIDI out port
func (o *out) Open() (err error) {
	if o.IsOpen() {
		return nil
	}
	//	o.Lock()
	//defer o.Unlock()
	o.midiOut, err = rtmidi.NewMIDIOutDefault()
	if err != nil {
		o.midiOut = nil
		return fmt.Errorf("can't open default MIDI out: %v", err)
	}

	err = o.midiOut.OpenPort(o.number, "")
	if err != nil {
		o.midiOut = nil
		return fmt.Errorf("can't open MIDI out port %v (%s): %v", o.number, o, err)
	}

	//	o.driver.Lock()
	o.driver.opened = append(o.driver.opened, o)
	//	o.driver.Unlock()

	return nil
}
