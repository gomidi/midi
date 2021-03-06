package portmididrv

import (
	"fmt"

	"github.com/rakyll/portmidi"
	"gitlab.com/gomidi/midi/v2"
)

func newOut(driver *driver, id portmidi.DeviceID, name string) midi.Out {
	return &out{driver: driver, id: id, name: name}
}

type out struct {
	id     portmidi.DeviceID
	stream *portmidi.Stream
	name   string
	driver *driver
}

// IsOpen returns, wether the port is open
func (o *out) IsOpen() bool {
	return o.stream != nil
}

// Send writes a MIDI message to the outut port
// If the output port is closed, it returns midi.ErrPortClosed
func (o *out) Send(b []byte) error {
	if o.stream == nil {
		return midi.ErrPortClosed
	}

	if len(b) < 2 {
		return fmt.Errorf("cannot send less than two message bytes")
	}

	var last int64
	// ProgramChange messages only have 2 bytes
	if len(b) > 2 {
		last = int64(b[2])
	}

	err := o.stream.WriteShort(int64(b[0]), int64(b[1]), last)
	if err != nil {
		return fmt.Errorf("could not send message to MIDI out %v (%s): %v", o.Number(), o, err)
	}
	return nil
}

// Underlying returns the underlying *portmidi.Stream. It will be nil, if the port is closed.
// Use it with type casting:
//   portOut := o.Underlying().(*portmidi.Stream)
func (o *out) Underlying() interface{} {
	return o.stream
}

// Number returns the number of the MIDI out port.
// Note that with portmidi, out and in ports are counted together.
// That means there should not be an out port with the same number as an in port.
func (o *out) Number() int {
	return int(o.id)
}

// String returns the name of the MIDI out port.
func (o *out) String() string {
	return o.name
}

// Close closes the MIDI out port
func (o *out) Close() error {
	if o.stream == nil {
		return nil
	}

	err := o.stream.Close()
	if err != nil {
		return fmt.Errorf("can't close MIDI out %v (%s): %v", o.Number(), o, err)
	}
	o.stream = nil
	return nil
}

// Open opens the MIDI output port
func (o *out) Open() (err error) {
	if o.stream != nil {
		return nil
	}
	o.stream, err = portmidi.NewOutputStream(o.id, o.driver.buffersizeOut, 0)
	if err != nil {
		o.stream = nil
		return fmt.Errorf("can't open MIDI out port %v (%s): %v", o.Number(), o, err)
	}
	o.driver.opened = append(o.driver.opened, o)
	return nil
}
