// Copyright (c) 2018 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package testdrv provides a Driver for testing.
*/
package testdrv

import (
	//"sync"

	"time"

	"gitlab.com/gomidi/midi/v2"
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
	now           time.Time
	stopListening bool
	rd            *drivers.Reader
	//wg            sync.WaitGroup
}

func New(name string) *Driver {
	d := &Driver{name: name}
	d.in = &in{name: name + "-in", Driver: d, number: 0}
	d.out = &out{name: name + "-out", Driver: d, number: 0}
	d.last = time.Now()
	d.now = d.last
	return d
}

func (f *Driver) Sleep(d time.Duration) {
	f.now = f.now.Add(d)
}

// wait until all messages are handled
/*
func (f *Driver) Wait() {
	f.wg.Wait()
}
*/

func (f *Driver) String() string               { return f.name }
func (f *Driver) Close() error                 { return nil }
func (f *Driver) Ins() ([]drivers.In, error)   { return []drivers.In{f.in}, nil }
func (f *Driver) Outs() ([]drivers.Out, error) { return []drivers.Out{f.out}, nil }

type in struct {
	number int
	name   string
	isOpen bool
	*Driver
}

func (f *in) String() string          { return f.name }
func (f *in) Number() int             { return f.number }
func (f *in) IsOpen() bool            { return f.isOpen }
func (f *in) Underlying() interface{} { return nil }

func (f *in) Listen(onMsg func(msg []byte, milliseconds int32), conf drivers.ListenConfig) (stopFn func(), err error) {
	//fmt.Printf("listeining from in port of %s\n", f.Driver.name)

	f.last = time.Now()

	stopFn = func() {
		f.stopListening = true
	}

	f.rd = drivers.NewReader(conf, func(m []byte, ms int32) {
		msg := midi.Message(m)

		if msg.Is(midi.ActiveSenseMsg) && !conf.ActiveSense {
			return
		}

		if msg.Is(midi.TimingClockMsg) && !conf.TimeCode {
			return
		}

		if msg.Is(midi.SysExMsg) && !conf.SysEx {
			return
		}

		//fmt.Printf("handle message % X at [%v] in driver %q\n", m, ms, f.Driver.name)
		onMsg(m, ms)
		//	f.wg.Done()
		//fmt.Println("msg handled")
	})
	f.rd.Reset()
	return stopFn, nil
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
	*Driver
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

func (f *out) Send(bt []byte) error {
	if !f.isOpen {
		return drivers.ErrPortClosed
	}

	if f.stopListening {
		return nil
	}

	dur := f.now.Sub(f.last)
	ts_ms := int32(dur.Milliseconds())
	f.last = f.now
	//f.wg.Add(1)
	//fmt.Printf("message added % X (len %v) at [%v] in driver %q\n", bt, len(bt), ts_ms, f.Driver.name)
	f.rd.EachMessage(bt, ts_ms)
	/*
		f.rd.SetDelta(ts_ms)
		for _, b := range bt {
			f.rd.EachByte(b)
		}
	*/
	return nil
}

func (f *out) Open() error {
	if f.isOpen {
		return nil
	}
	f.isOpen = true
	return nil
}
