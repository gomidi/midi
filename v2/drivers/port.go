package drivers

import (
	"fmt"
	"strings"
)

var ErrPortClosed = fmt.Errorf("ERROR: port is closed")
var ErrListenStopped = fmt.Errorf("ERROR: stopped listening")

// Port is an interface for a MIDI port.
// In order to be lockless (for realtime), a port is not threadsafe, so none of its method may be called
// from different goroutines.
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

	// Underlying returns the underlying implementation of the driver
	Underlying() interface{}
}

// ListenConfig defines the configuration for in port listening
type ListenConfig struct {

	// TimeCode lets the timecode messages pass through, if set
	TimeCode bool

	// ActiveSense lets the active sense messages pass through, if set
	ActiveSense bool

	// SysEx lets the sysex messaes pass through, if set
	SysEx bool

	// SysExBufferSize defines the size of the buffer for sysex messages (in bytes).
	// SysEx messages larger than this size will be ignored.
	// When SysExBufferSize is 0, the default buffersize (1024) is used.
	SysExBufferSize uint32

	// OnErr is the callback that is called for any error happening during the listening.
	OnErr func(error)
}

// In is an interface for a MIDI input port
type In interface {
	Port

	// Listen listens for incoming messages. It returns a function that must be used to stop listening.
	// The onMsg callback is called for every non-sysex message. The onMsg callback must not be nil.
	// The config defines further listening options (see ListenConfig)
	// The listening must be stopped before the port may be closed.
	Listen(
		onMsg func(msg []byte, milliseconds int32),
		config ListenConfig,
	) (
		stopFn func(),
		err error,
	)
}

// Out is an interface for a MIDI output port.
type Out interface {
	Port

	//Send(data [3]byte) error
	Send(data []byte) error
}

// Ins return the available MIDI in ports
func Ins() ([]In, error) {
	d := Get()
	if d == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return d.Ins()
}

// Outs return the available MIDI out ports
func Outs() ([]Out, error) {
	d := Get()
	if d == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return d.Outs()
}

// InByName opens the first midi in port that contains the given name
func InByName(portName string) (in In, err error) {
	drv := Get()
	if drv == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return openIn(drv, -1, portName)
}

// InByNumber opens the midi in port with the given number
func InByNumber(portNumber int) (in In, err error) {
	drv := Get()
	if drv == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return openIn(drv, portNumber, "")
}

// OutByName opens the first midi out port that contains the given name
func OutByName(portName string) (out Out, err error) {
	drv := Get()
	if drv == nil {
		return nil, fmt.Errorf("no driver registered")
	}
	return openOut(drv, -1, portName)
}

// OutByNumber opens the midi out port with the given number
func OutByNumber(portNumber int) (out Out, err error) {
	drv := Get()
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
