package midispy

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

// Run will read the messages from the given in port,
// pass them to the given midi.Receiver and given an out port
// that is >= 0, write them to the out port.
// All given port must be opened. Run will not close any ports.
// Stop the spying by closing the ports.
func Run(in drivers.In, out drivers.Out, recv midi.Receiver) error {
	if in == nil {
		panic("MIDI in port is nil")
	}
	s := &spy{
		out: out,
		In:  in,
	}

	return s.SetListener(recv.Receive)
}

// spy connects a MIDI in port and a MIDI out port and allows
// inspecting the MIDI messages flowing from in to out with the help of
// a mid.Reader.
type spy struct {
	out drivers.Out
	drivers.In
}

// SetListener sets the given outer function as listener
func (s *spy) SetListener(outer func(midi.Message, int32)) (err error) {

	var listener func([]byte, int32)

	switch {
	case s.out == nil && outer == nil:
		listener = func(bt []byte, millisecs int32) {}
	case s.out == nil && outer != nil:
		listener = func(bt []byte, millisecs int32) { outer(bt, millisecs) }
	case s.out != nil && outer == nil:
		listener = func(bt []byte, millisecs int32) {
			s.out.Send(bt)
		}
	case s.out != nil && outer != nil:
		listener = func(bt []byte, millisecs int32) {
			outer(bt, millisecs)
			s.out.Send(bt)
		}
	}

	var stop func()

	if listener != nil {
		fmt.Println("setting listener")
		var o drivers.ListenConfig
		stop, err = s.In.Listen(listener, o)
		//err = s.In.SetListener(listener)
		if err != nil {
			return
		}
	}
	_ = stop
	return nil
}
