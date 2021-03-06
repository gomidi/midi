// package midi
package midi

import (
	"fmt"
)

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

// TODO change signature to OpenIn(connstr string) (in In, err error)
// where connstr has syntax like:
// "rtmidi:my-device" (takes rtmidi driver)
// or
// "my-device" (takes first available driver)
// or
// "/device[0-9]/" (takes first device matching regular expression
// or
// "rtmidi:1"  (takes second device of rtmidi driver)
// or
// "rtmidi" (takes first device of rtmidi driver)
// or
// "1" (takes second device of first available driver)
// or
// "" (takes first device of first available driver)
// the drivers will register themselves into the midi library
// if a driver has to be initialized with options, it must be registered "by hand"
// drivers are only opened, when they are used the first time

// OpenIn opens a MIDI input port with the help of the given driver.
// To find the port by port number, pass a number >= 0.
// To find the port by port name, pass a number < 0 and a non empty string.
func OpenIn(d Driver, number int, name string) (in In, err error) {
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

// TODO change signature to OpenOut(connstr string) (out Out, err error)
// where connstr has syntax like:
// "rtmidi:my-device" (takes rtmidi driver)
// or
// "my-device" (takes first available driver)
// or
// "/device[0-9]/" (takes first device matching regular expression
// or
// "rtmidi:1"  (takes second device of rtmidi driver)
// or
// "rtmidi" (takes first device of rtmidi driver)
// or
// "1" (takes second device of first available driver)
// or
// "" (takes first device of first available driver)
// the drivers will register themselves into the midi library
// if a driver has to be initialized with options, it must be registered "by hand"
// drivers are only opened, when they are used the first time

// OpenOut opens a MIDI output port with the help of the given driver.
// To find the port by port number, pass a number >= 0.
// To find the port by port name, pass a number < 0 and a non empty string.
func OpenOut(d Driver, number int, name string) (out Out, err error) {
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
