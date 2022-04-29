// +build js,wasm,!windows,!linux,!darwin

package webmididrv

import (
	"sync"
	"syscall/js"

	"gitlab.com/gomidi/midi/v2/drivers"
)

func newOut(driver *Driver, number int, name string, jsport js.Value) drivers.Out {
	o := &out{driver: driver, number: number, name: name, jsport: jsport}
	return o
}

type out struct {
	number int
	sync.RWMutex
	driver *Driver
	name   string
	jsport js.Value
	isOpen bool
}

// IsOpen returns wether the port is open
func (o *out) IsOpen() (open bool) {
	o.RLock()
	open = o.isOpen
	o.RUnlock()
	return
}

// Send writes a MIDI message to the MIDI output port
// If the output port is closed, it returns midi.ErrClosed
func (o *out) Send(b []byte) error {
	o.RLock()
	if !o.isOpen {
		o.RUnlock()
		return drivers.ErrPortClosed
	}
	o.RUnlock()

	var arr = make([]interface{}, len(b))
	for i, bt := range b {
		arr[i] = bt
	}

	o.jsport.Call("send", js.ValueOf(arr))
	return nil
}

// Underlying returns the underlying driver. Here it returns the js output port.
func (o *out) Underlying() interface{} {
	return o.jsport
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

	o.Lock()
	defer o.Unlock()
	o.isOpen = false
	o.jsport.Call("close")
	return err
}

// Open opens the MIDI out port
func (o *out) Open() (err error) {
	if o.IsOpen() {
		return nil
	}

	o.driver.Lock()
	o.isOpen = true
	o.jsport.Call("open")
	o.driver.opened = append(o.driver.opened, o)
	o.driver.Unlock()
	return nil
}
