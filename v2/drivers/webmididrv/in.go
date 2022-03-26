package webmididrv

import (
	"math"
	"sync"
	"sync/atomic"
	"syscall/js"

	"gitlab.com/gomidi/midi/v2/drivers"
)

type in struct {
	number int
	sync.RWMutex
	driver   *Driver
	name     string
	isOpen   bool
	jsport   js.Value
	listener func(data []byte, timestamp int32)
}

// IsOpen returns wether the MIDI in port is open
func (o *in) IsOpen() (open bool) {
	o.RLock()
	open = o.isOpen
	o.RUnlock()
	return
}

// String returns the name of the MIDI in port.
func (i *in) String() string {
	return i.name
}

// Underlying returns the underlying driver. Here returns the js midi port.
func (i *in) Underlying() interface{} {
	return i.jsport
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

	i.Lock()
	i.isOpen = false
	i.jsport.Call("close")
	i.Unlock()
	return
}

// Open opens the MIDI in port
func (i *in) Open() (err error) {
	if i.IsOpen() {
		return nil
	}
	i.Lock()
	i.isOpen = true
	i.jsport.Call("open")
	i.Unlock()

	i.driver.Lock()
	i.driver.opened = append(i.driver.opened, i)
	i.driver.Unlock()

	return nil
}

func newIn(driver *Driver, number int, name string, jsport js.Value) drivers.In {
	return &in{driver: driver, number: number, name: name, jsport: jsport}
}

/*
    i.Lock()
	//i.listener = recv.Receive

	jsCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsdata := args[0].Get("data")
		jstime := args[0].Get("receivedTime")

		var data = make([]byte, 3)
		data[0] = byte(jsdata.Index(0).Int())
		data[1] = byte(jsdata.Index(1).Int())
		data[2] = byte(jsdata.Index(2).Int())
		var t = int32(-1)
		if jstime.Truthy() {
			// round to milliseconds
			t = int32(math.Round(jstime.Float()))
		}
		//i.listener(data, t)
		cb(data, t)
		return nil
	})

	i.jsport.Call("addEventListener", "midimessage", jsCallback)
	i.Unlock()

*/

func (i *in) Listen(onMsg func(msg []byte, milliseconds int32), config drivers.ListenConfig) (stopFn func(), err error) {

	var stop int32

	//stopWait := i.driver.sleepingTime * 2
	stopFn = func() {
		// lockless sync
		atomic.StoreInt32(&stop, 1)
		//time.Sleep(stopWait)
	}

	i.Lock()

	var stopped int32

	jsCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// lockless sync
		stopped = atomic.LoadInt32(&stop)

		if stopped == 1 {
			return nil
		}

		jsdata := args[0].Get("data")
		jstime := args[0].Get("receivedTime")

		var data = make([]byte, 3)
		data[0] = byte(jsdata.Index(0).Int())
		data[1] = byte(jsdata.Index(1).Int())
		data[2] = byte(jsdata.Index(2).Int())
		var t = int32(-1)
		if jstime.Truthy() {
			// round to milliseconds
			t = int32(math.Round(jstime.Float()))
		}
		onMsg(data, t)
		return nil
	})

	go i.jsport.Call("addEventListener", "midimessage", jsCallback)
	i.Unlock()

	return
}

/*
// SendTo
func (i *in) StartListening(cb func(data []byte, timestamp int32)) (err error) {
	if !i.IsOpen() {
		return drivers.ErrPortClosed
	}

	i.RLock()
	if i.listener != nil {
		i.RUnlock()
		return fmt.Errorf("listener already set")
	}
	i.RUnlock()
	i.Lock()
	//i.listener = recv.Receive

	jsCallback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsdata := args[0].Get("data")
		jstime := args[0].Get("receivedTime")

		var data = make([]byte, 3)
		data[0] = byte(jsdata.Index(0).Int())
		data[1] = byte(jsdata.Index(1).Int())
		data[2] = byte(jsdata.Index(2).Int())
		var t = int32(-1)
		if jstime.Truthy() {
			// round to milliseconds
			t = int32(math.Round(jstime.Float()))
		}
		//i.listener(data, t)
		cb(data, t)
		return nil
	})

	i.jsport.Call("addEventListener", "midimessage", jsCallback)
	i.Unlock()

	return nil
}

// StopListening cancels the listening
func (i *in) StopListening() (err error) {
	if !i.IsOpen() {
		return drivers.ErrPortClosed
	}

	// TODO
	return
}
*/
