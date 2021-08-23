package portmididrv

import (
	"fmt"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/portmididrv/imported/portmidi"
)

func init() {
	drv, err := New()
	if err != nil {
		panic(fmt.Sprintf("could not register portmididrv: %s", err.Error()))
	}
	drivers.Register(drv)
}

var _ drivers.Driver = &Driver{}

//var _ drivers.SysExListener = &in{}
//var _ drivers.SysCommonListener = &in{}
//var _ drivers.RealtimeListener = &in{}
//var _ drivers.SysExSender = &out{}

type Driver struct {
	buffersizeRead int
	buffersizeIn   int64
	buffersizeOut  int64
	sleepingTime   time.Duration
	sync.Mutex
	opened []drivers.Port
}

func (d *Driver) String() string {
	return "portmididrv"
}

// Close closes all open ports. It must be called at the end of a session.
func (d *Driver) Close() (err error) {
	d.Lock()
	defer d.Unlock()

	//fmt.Println("close out devices")
	for _, p := range d.opened {
		if _, isOut := p.(*out); isOut {
			err = p.Close()
		}
	}

	//fmt.Println("close in devices")
	for _, p := range d.opened {
		if _, isIn := p.(*in); isIn {
			err = p.Close()
		}
	}

	//fmt.Println("all devices closed")
	// return just the last error to allow closing the other ports.
	// to ensure that all ports have been closed, this function must
	// return nil anyways
	return
}

// New returns a new driver
func New(options ...Option) (*Driver, error) {
	err := portmidi.Initialize()
	if err != nil {
		return nil, fmt.Errorf("can't initialize portmidi: %v", err)
	}
	dr := &Driver{}

	dr.buffersizeRead = 1024
	//dr.buffersizeRead = 1
	//dr.buffersizeIn = 1024
	//dr.buffersizeIn = 100
	dr.buffersizeIn = 1024
	//dr.buffersizeIn = 10
	dr.buffersizeOut = 1024
	//dr.buffersizeOut = 100
	//dr.buffersizeOut = 0

	// sleepingTime of 0.5ms should be fine to prevent busy waiting
	// and still fast enough for performances
	//dr.sleepingTime = time.Nanosecond * 1000 * 100
	//dr.sleepingTime = time.Millisecond
	//dr.sleepingTime = time.Millisecond
	dr.sleepingTime = time.Microsecond * 400
	//dr.sleepingTime = time.Millisecond * 10
	//dr.sleepingTime = time.Nanosecond * 1000 * 500

	for _, opt := range options {
		opt(dr)
	}

	return dr, nil
}

// Ins returns the available MIDI in ports
func (d *Driver) Ins() (ins []drivers.In, err error) {
	var num int
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if info != nil && info.IsInputAvailable {
			ins = append(ins, newIn(d, portmidi.DeviceID(i), num, info.Name))
			num++
		}
	}
	return
}

// Outs returns the available MIDI out ports
func (d *Driver) Outs() (outs []drivers.Out, err error) {
	var num int
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		//		fmt.Printf("%q devideID %v\n", info.Name, portmidi.DeviceID(i))
		if info != nil && info.IsOutputAvailable {
			//		fmt.Printf("registering out port %q, number [%v], devideID %v\n", info.Name, num, portmidi.DeviceID(i))
			outs = append(outs, newOut(d, portmidi.DeviceID(i), num, info.Name))
			num++
		}
	}
	return
}
