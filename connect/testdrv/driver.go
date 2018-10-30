package testdrv

import (
	"io"

	"gitlab.com/gomidi/midi/connect"
)

type Driver struct {
	in       *in
	out      *out
	listener func([]byte, int64)
	name     string
}

func New(name string) connect.Driver {
	d := &Driver{name: name}
	d.in = &in{name: name + "-in", driver: d, number: 0}
	d.out = &out{name: name + "-out", driver: d, number: 0}
	return d
}

func (f *Driver) String() string               { return f.name }
func (f *Driver) Close() error                 { return nil }
func (f *Driver) Ins() ([]connect.In, error)   { return []connect.In{f.in}, nil }
func (f *Driver) Outs() ([]connect.Out, error) { return []connect.Out{f.out}, nil }

type in struct {
	number int
	name   string
	isOpen bool
	driver *Driver
}

func (f *in) StopListening() error    { f.driver.listener = nil; return nil }
func (f *in) String() string          { return f.name }
func (f *in) Number() int             { return f.number }
func (f *in) IsOpen() bool            { return f.isOpen }
func (f *in) Underlying() interface{} { return nil }

func (f *in) SetListener(listener func([]byte, int64)) error {
	f.driver.listener = listener
	return nil
}
func (f *in) Close() error {
	if !f.isOpen {
		return nil
	}
	f.StopListening()
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
func (f *out) Send(b []byte) error {
	if !f.isOpen {
		return connect.ErrClosed
	}
	if f.driver.listener == nil {
		return io.EOF
	}
	f.driver.listener(b, 0)
	return nil
}

func (f *out) Open() error {
	if f.isOpen {
		return nil
	}
	f.isOpen = true
	return nil
}
