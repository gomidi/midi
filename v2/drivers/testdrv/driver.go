// Copyright (c) 2018 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package testdrv provides a gomidi/midi.Driver for testing.

*/
package testdrv

import (
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
)

func init() {
	drv := New("testdrv")
	drivers.Register(drv)
}

type Driver struct {
	in            *in
	out           *out
	name          string
	last          time.Time
	stopListening bool
	rd            *drivers.Reader
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

func (f *in) String() string          { return f.name }
func (f *in) Number() int             { return f.number }
func (f *in) IsOpen() bool            { return f.isOpen }
func (f *in) Underlying() interface{} { return nil }

func (f *in) Listen(onMsg func([3]byte, int32), conf drivers.ListenConfig) (func(), error) {
	f.driver.last = time.Now()

	stopper := func() {
		f.driver.stopListening = true
	}

	f.driver.rd = drivers.NewReader(conf, onMsg)

	return stopper, nil
}

func (f *in) Close() error {
	if !f.isOpen {
		return nil
	}
	f.isOpen = false
	return nil
}

func (f *in) Open() error {
	if f.isOpen {
		return nil
	}
	f.isOpen = true
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
	if !f.isOpen {
		return nil
	}
	f.isOpen = false
	return nil
}

func (f *out) Send(b [3]byte) error {
	if !f.isOpen {
		return drivers.ErrPortClosed
	}

	if f.driver.stopListening {
		return nil
	}

	now := time.Now()
	dur := now.Sub(f.driver.last)
	ts_ms := int32(dur.Milliseconds())
	f.driver.last = now

	var bt []byte

	switch {
	case b[2] == 0 && b[1] == 0:
		bt = []byte{b[0]}
		//	case b[2] == 0:
	//	bt = []byte{b[0], b[1]}
	default:
		bt = []byte{b[0], b[1], b[2]}
	}

	f.driver.rd.EachMessage(bt, ts_ms)
	return nil
}

func (f *out) Open() error {
	if f.isOpen {
		return nil
	}
	f.isOpen = true
	return nil
}
