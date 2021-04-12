package midi

import (
	"fmt"
)

var DRIVERS = map[string]Driver{}
var firstDriver string

// RegisterDriver register a driver
func RegisterDriver(d Driver) {
	DRIVERS[d.String()] = d
	if len(DRIVERS) == 0 {
		firstDriver = d.String()
	}
}

func GetDriver() Driver {
	if len(DRIVERS) == 0 {
		return nil
	}
	return DRIVERS[firstDriver]
}

func CloseDriver() {
	d := GetDriver()
	if d != nil {
		d.Close()
	}
}

func Ins() ([]In, error) {
	d := GetDriver()
	if d == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return d.Ins()
}

func Outs() ([]Out, error) {
	d := GetDriver()
	if d == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return d.Outs()
}

// Driver is a driver for MIDI connections.
type Driver interface {

	// Ins returns the available MIDI input ports.
	Ins() ([]In, error)

	// Outs returns the available MIDI output ports.
	Outs() ([]Out, error)

	// String returns the name of the driver.
	String() string

	// Close closes the driver. Must be called for cleanup at the end of a session.
	Close() error
}

// Port is an interface for a MIDI port.
type Port interface {

	// Open opens the MIDI port. An implementation should save the open state to make it
	// save to call open when the port is already open without getting an error.
	Open() error

	// Close closes the MIDI port. An implementation should save the open state to make it
	// save to call close when the port is already closed without getting an error.
	Close() error

	// IsOpen returns wether the MIDI port is open.
	IsOpen() bool

	// Number returns the number of the MIDI port. It is only guaranteed that the numbers are unique within
	// MIDI port groups i.e. within MIDI input ports and MIDI output ports. So there may be the same number
	// for a given MIDI input port and some MIDI output port. Or not - that depends on the underlying driver.
	Number() int

	// String represents the MIDI port by a string, aka name.
	String() string

	// Underlying returns the underlying driver to allow further adjustments.
	// When using the underlying driver, the user must take care of proper opening/closing etc.
	Underlying() interface{}
}

// println(big.NewRat(math.MaxInt64,1000 /* milliseonds */ *1000 /* seconds */ *60 /* minutes */ *60 /* hours */ *24 /* days */ *365 /* years */).FloatString(0))
// output: 292471
// => a ascending timestamp based on microseconds would wrap after 292471 years

// In is an interface for a MIDI input port
type In interface {
	Port
	SenderTo

	// StopListening stops the listening.
	// When closing a MIDI input port, StopListening must be called before (from the driver).
	StopListening() error
}

// Out is an interface for a MIDI output port.
type Out interface {
	Port
	Sender
}

// ErrPortClosed should be returned from a driver when trying to write to a closed port.
var ErrPortClosed = fmt.Errorf("ERROR: port is closed")

type listener struct {
	err              error
	in               In
	filter           Filter
	realtimeCallback func(msg Message, deltamicrosec int64)
}

func (l *listener) Error() error {
	return l.err
}

func Listen(portName string) *listener {
	l := &listener{}
	l.in, l.err = InByName(portName)
	return l
}

func (l *listener) Only(mtypes ...MsgType) *listener {
	if len(mtypes) > 0 {
		l.filter = Filter(mtypes)
	}
	return l
}

func (l *listener) RealTime(realtimeMsgCallback func(msg Message, deltamicrosec int64)) *listener {
	l.realtimeCallback = realtimeMsgCallback
	return l
}

func (l *listener) Do(fn func(msg Message, deltamicroSec int64)) (In, error) {
	if l.err != nil {
		return l.in, l.err
	}

	var rec Receiver

	if l.filter == nil {
		rec = NewReceiver(fn, l.realtimeCallback)
	} else {
		var fun = func(m Message, delta int64) {
			//m := NewMessage(msg)

			if m.MsgType.IsOneOf(l.filter...) {
				fn(m, delta)
			}
		}

		var funrt func(m Message, delta int64)
		if l.realtimeCallback != nil {
			funrt = func(m Message, delta int64) {
				//m := NewMessage(msg)

				if m.MsgType.IsOneOf(l.filter...) {
					l.realtimeCallback(m, delta)
				}
			}
		}

		rec = NewReceiver(fun, funrt)
	}

	l.in.SendTo(rec)
	return l.in, nil
}

// InByName opens the first midi in port that contains the given name
func InByName(portName string) (in In, err error) {
	drv := GetDriver()
	if drv == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return openIn(drv, -1, portName)
}

// InByNumber opens the midi in port with the given number
func InByNumber(portNumber int) (in In, err error) {
	drv := GetDriver()
	if drv == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return openIn(drv, portNumber, "")
}

// OutByName opens the first midi out port that contains the given name
func OutByName(portName string) (out Out, err error) {
	drv := GetDriver()
	if drv == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return openOut(drv, -1, portName)
}

// OutByNumber opens the midi out port with the given number
func OutByNumber(portNumber int) (out Out, err error) {
	drv := GetDriver()
	if drv == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return openOut(drv, portNumber, "")
}

// openIn opens a MIDI input port with the help of the given driver.
// To find the port by port number, pass a number >= 0.
// To find the port by port name, pass a number < 0 and a non empty string.
func openIn(d Driver, number int, name string) (in In, err error) {
	ins, err := d.Ins()
	if err != nil {
		return nil, fmt.Errorf("can't find MIDI input ports: %v", err)
	}

	if number >= 0 {
		for _, port := range ins {
			if number == port.Number() {
				in = port
				break
			}
		}
		if in == nil {
			return nil, fmt.Errorf("can't find MIDI input port %v", number)
		}
	} else {
		if name != "" {
			for _, port := range ins {
				if name == port.String() {
					in = port
					break
				}
			}
		}
		if in == nil {
			return nil, fmt.Errorf("can't find MIDI input port %v", name)
		}
	}

	// should not happen here, since we already returned above
	if in == nil {
		panic("unreachable")
	}

	err = in.Open()
	return
}

// openOut opens a MIDI output port with the help of the given driver.
// To find the port by port number, pass a number >= 0.
// To find the port by port name, pass a number < 0 and a non empty string.
func openOut(d Driver, number int, name string) (out Out, err error) {
	outs, err := d.Outs()
	if err != nil {
		return nil, fmt.Errorf("can't find MIDI output ports: %v", err)
	}

	if number >= 0 {
		for _, port := range outs {
			if number == port.Number() {
				out = port
				break
			}
		}
		if out == nil {
			return nil, fmt.Errorf("can't find MIDI output port %v", number)
		}
	} else {
		if name != "" {
			for _, port := range outs {
				if name == port.String() {
					out = port
					break
				}
			}
		}
		if out == nil {
			return nil, fmt.Errorf("can't find MIDI output port %v", name)
		}
	}

	// should not happen here, since we already returned above
	if out == nil {
		panic("unreachable")
	}

	err = out.Open()
	return
}
