// Copyright (c) 2018 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package testdrv provides a gomidi/midi.Driver for testing.

*/
package testdrv

import (
	"io"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
)

func init() {
	drv := New("testdrv")
	drivers.Register(drv)
}

type Driver struct {
	in  *in
	out *out
	//reader   *drivers.DeviceReader
	callback        func([]byte, int32)
	name            string
	last            time.Time
	absdecimillisec int32
	mx              sync.Mutex
}

func New(name string) drivers.Driver {
	d := &Driver{name: name}
	d.in = &in{name: name + "-in", driver: d, number: 0}
	d.out = &out{name: name + "-out", driver: d, number: 0}
	d.last = time.Now()
	return d
}

func (f *Driver) String() string               { return f.name }
func (f *Driver) Close() error                 { return nil }
func (f *Driver) Ins() ([]drivers.In, error)   { return []drivers.In{f.in}, nil }
func (f *Driver) Outs() ([]drivers.Out, error) { return []drivers.Out{f.out}, nil }

type in struct {
	number int
	name   string
	isOpen bool
	driver *Driver
}

func (f *in) StopListening() error {
	f.driver.mx.Lock()
	f.driver.callback = nil
	f.driver.mx.Unlock()
	return nil
}
func (f *in) String() string          { return f.name }
func (f *in) Number() int             { return f.number }
func (f *in) IsOpen() bool            { return f.isOpen }
func (f *in) Underlying() interface{} { return nil }

func (f *in) StartListening(cb func([]byte, int32)) error {
	f.driver.mx.Lock()
	f.driver.absdecimillisec = 0
	f.driver.last = time.Now()
	//f.driver.reader = drivers.NewDeviceReader(recv)
	f.driver.callback = cb
	f.driver.mx.Unlock()
	return nil
}

/*
func (f *in) SetListener(listener func([]byte, int64)) error {
	f.driver.mx.Lock()
	f.driver.listener = listener
	f.driver.mx.Unlock()
	return nil
}
*/

func (f *in) Close() error {
	f.driver.mx.Lock()
	if !f.isOpen {
		f.driver.mx.Unlock()
		return nil
	}
	f.driver.mx.Unlock()
	f.StopListening()
	f.driver.mx.Lock()
	f.isOpen = false
	f.driver.mx.Unlock()
	return nil
}

func (f *in) Open() error {
	f.driver.mx.Lock()
	if f.isOpen {
		f.driver.mx.Unlock()
		return nil
	}
	f.isOpen = true
	f.driver.mx.Unlock()
	return nil
}

type out struct {
	number int
	name   string
	isOpen bool
	driver *Driver
}

func (f *out) Number() int             { return f.number }
func (f *out) IsOpen() bool            { return f.isOpen }
func (f *out) String() string          { return f.name }
func (f *out) Underlying() interface{} { return nil }

func (f *out) Close() error {
	f.driver.mx.Lock()
	if !f.isOpen {
		f.driver.mx.Unlock()
		return nil
	}
	f.isOpen = false
	f.driver.mx.Unlock()
	return nil
}
func (f *out) Send(data []byte) error {
	f.driver.mx.Lock()
	if !f.isOpen {
		f.driver.mx.Unlock()
		return drivers.ErrPortClosed
	}
	if f.driver.callback == nil {
		f.driver.mx.Unlock()
		return io.EOF
	}

	now := time.Now()
	dur := now.Sub(f.driver.last)
	//f.driver.last = now
	f.driver.absdecimillisec = int32(dur.Milliseconds())
	f.driver.callback(data, f.driver.absdecimillisec)
	f.driver.mx.Unlock()
	return nil
}

func (f *out) Open() error {
	f.driver.mx.Lock()
	if f.isOpen {
		f.driver.mx.Unlock()
		return nil
	}
	f.isOpen = true
	f.driver.mx.Unlock()
	return nil
}
