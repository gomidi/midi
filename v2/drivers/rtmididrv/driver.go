package rtmididrv

import (
	"fmt"
	"strings"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi"
)

func init() {
	drv, err := New()
	if err != nil {
		panic(fmt.Sprintf("could not register rtmididrv: %s", err.Error()))
	}
	drivers.Register(drv)
}

type Driver struct {
	opened []drivers.Port
	//ignoreSysex       bool
	//ignoreTimeCode    bool
	//ignoreActiveSense bool
	//sync.RWMutex
}

func (d *Driver) String() string {
	return "rtmididrv"
}

// Close closes all open ports. It must be called at the end of a session.
func (d *Driver) Close() (err error) {
	//	d.Lock()
	var e CloseErrors

	for _, p := range d.opened {
		err = p.Close()
		if err != nil {
			e = append(e, err)
		}
	}

	//	d.Unlock()

	if len(e) == 0 {
		return nil
	}

	return e
}

/*
type Option func(*Driver)

func IgnoreSysex() Option {
	return func(d *Driver) {
		d.ignoreSysex = true
	}
}

func IgnoreTimeCode() Option {
	return func(d *Driver) {
		d.ignoreTimeCode = true
	}
}

func IgnoreActiveSense() Option {
	return func(d *Driver) {
		d.ignoreActiveSense = true
	}
}
*/

// New returns a driver based on the default rtmidi in and out
func New() (*Driver, error) {
	d := &Driver{}
	return d, nil
}

// OpenVirtualIn opens and returns a virtual MIDI in. We can't get the port number, so set it to -1.
func (d *Driver) OpenVirtualIn(name string) (drivers.In, error) {
	_in, err := rtmidi.NewMIDIInDefault()
	if err != nil {
		return nil, fmt.Errorf("can't open default MIDI in: %v", err)
	}

	err = _in.OpenVirtualPort(name)

	if err != nil {
		return nil, fmt.Errorf("can't open virtual in port: %s", err.Error())
	}

	//	d.Lock()
	//defer d.Unlock()
	//_in.IgnoreTypes(d.ignoreSysex, d.ignoreTimeCode, d.ignoreActiveSense)
	inPort := &in{driver: d, number: -1, name: name, midiIn: _in}
	d.opened = append(d.opened, inPort)
	return inPort, nil
}

// OpenVirtualOut opens and returns a virtual MIDI out. We can't get the port number, so set it to -1.
func (d *Driver) OpenVirtualOut(name string) (drivers.Out, error) {
	_out, err := rtmidi.NewMIDIOutDefault()
	if err != nil {
		return nil, fmt.Errorf("can't open default MIDI out: %v", err)
	}

	err = _out.OpenVirtualPort(name)

	if err != nil {
		return nil, fmt.Errorf("can't open virtual out port: %s", err.Error())
	}

	//d.Lock()
	//defer d.Unlock()
	outPort := &out{driver: d, number: -1, name: name, midiOut: _out}
	d.opened = append(d.opened, outPort)
	return outPort, nil
}

// Ins returns the available MIDI input ports
func (d *Driver) Ins() (ins []drivers.In, err error) {
	var in rtmidi.MIDIIn
	in, err = rtmidi.NewMIDIInDefault()
	if err != nil {
		return nil, fmt.Errorf("can't open default MIDI in: %v", err)
	}

	ports, err := in.PortCount()
	if err != nil {
		return nil, fmt.Errorf("can't get number of in ports: %s", err.Error())
	}

	for i := 0; i < ports; i++ {
		name, err := in.PortName(i)
		if err != nil {
			name = ""
		}
		ins = append(ins, newIn(d, i, name))
	}

	// don't destroy, destroy just panics
	// in.Destroy()
	err = in.Close()
	return
}

// Outs returns the available MIDI output ports
func (d *Driver) Outs() (outs []drivers.Out, err error) {
	var out rtmidi.MIDIOut
	out, err = rtmidi.NewMIDIOutDefault()
	if err != nil {
		return nil, fmt.Errorf("can't open default MIDI out: %v", err)
	}

	ports, err := out.PortCount()
	if err != nil {
		return nil, fmt.Errorf("can't get number of out ports: %s", err.Error())
	}

	for i := 0; i < ports; i++ {
		name, err := out.PortName(i)
		if err != nil {
			name = ""
		}
		outs = append(outs, newOut(d, i, name))
	}

	err = out.Close()
	return
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
