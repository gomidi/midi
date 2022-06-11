//go:build js && wasm && !windows && !linux && !darwin
// +build js,wasm,!windows,!linux,!darwin

package webmididrv

import (
	"fmt"
	"strings"
	"sync"
	"syscall/js"

	"gitlab.com/gomidi/midi/v2/drivers"
)

func init() {
	drv, err := New()
	if err != nil {
		panic(fmt.Sprintf("could not register webmididrv: %s", err.Error()))
	}
	drivers.Register(drv)
}

type Driver struct {
	opened []drivers.Port
	sync.RWMutex
	inputsJS  js.Value
	outputsJS js.Value
	wg        sync.WaitGroup
	Err       error
}

func (d *Driver) String() string {
	return "webmididrv"
}

// Close closes all open ports. It must be called at the end of a session.
func (d *Driver) Close() (err error) {
	d.Lock()
	var e CloseErrors

	for _, p := range d.opened {
		err = p.Close()
		if err != nil {
			e = append(e, err)
		}
	}

	d.Unlock()

	if len(e) == 0 {
		return nil
	}

	return e
}

// New returns a driver based on the js webmidi standard
func New() (*Driver, error) {
	jsDoc := js.Global().Get("navigator")
	if !jsDoc.Truthy() {
		return nil, fmt.Errorf("Unable to get navigator object")
	}

	// currently sysex messages are not allowed in the browser implementations
	var opts = map[string]interface{}{
		"sysex": "false",
	}

	jsOpts := js.ValueOf(opts)

	midiaccess := jsDoc.Call("requestMIDIAccess", jsOpts)
	if !midiaccess.Truthy() {
		return nil, fmt.Errorf("unable to get requestMIDIAccess")
	}

	drv := &Driver{}
	drv.wg.Add(1)
	midiaccess.Call("then", drv.onMIDISuccess(), drv.onMIDIFailure())
	drv.wg.Wait()
	return drv, nil
}

func (d *Driver) onMIDISuccess() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}

		d.inputsJS = args[0].Get("inputs")
		d.outputsJS = args[0].Get("outputs")
		d.wg.Done()
		return nil
	})
}

func (d *Driver) onMIDIFailure() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		d.Err = fmt.Errorf("Could not access the MIDI devices.")
		d.wg.Done()
		return nil
	})
}

// Ins returns the available MIDI input ports
func (d *Driver) Ins() (ins []drivers.In, err error) {
	if d.Err != nil {
		return nil, err
	}

	if !d.inputsJS.Truthy() {
		return nil, fmt.Errorf("no inputs")
	}

	var i = 0

	eachIn := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsport := args[0]
		var name = jsport.Get("name").String()
		ins = append(ins, newIn(d, i, name, jsport))
		i++
		return nil
	})

	d.inputsJS.Call("forEach", eachIn)
	return ins, nil
}

// Outs returns the available MIDI output ports
func (d *Driver) Outs() (outs []drivers.Out, err error) {
	if d.Err != nil {
		return nil, err
	}

	if !d.outputsJS.Truthy() {
		return nil, fmt.Errorf("no outputs")
	}

	var i = 0

	eachOut := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsport := args[0]
		var name = jsport.Get("name").String()
		outs = append(outs, newOut(d, i, name, jsport))
		i++
		return nil
	})

	d.outputsJS.Call("forEach", eachOut)

	return outs, nil
}

// CloseErrors collects error from closing multiple MIDI ports
type CloseErrors []error

func (c CloseErrors) Error() string {
	if len(c) == 0 {
		return "no errors"
	}

	var bd strings.Builder

	bd.WriteString("the following closing errors occured:\n")

	for _, e := range c {
		bd.WriteString(e.Error() + "\n")
	}

	return bd.String()
}
