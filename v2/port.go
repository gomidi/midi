package midi

import (
	"fmt"
	"strings"
)

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

// Ins return the available MIDI in ports
func Ins() ([]In, error) {
	d := GetDriver()
	if d == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return d.Ins()
}

// Outs return the available MIDI out ports
func Outs() ([]Out, error) {
	d := GetDriver()
	if d == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return d.Outs()
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
				if strings.Contains(port.String(), name) {
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
				if strings.Contains(port.String(), name) {
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
