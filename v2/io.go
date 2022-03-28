package midi

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi/v2/drivers"
)

// CloseDriver closes the default driver.
func CloseDriver() {
	drivers.Close()
}

// Sender sends MIDI messages.
type Sender interface {
	// Send sends the given MIDI message and returns any error.
	Send(msg Message) error
}

// ReceiverFunc is a function that receives MIDI messages
type ReceiverFunc func(msg Message, absdecimillisec int32)

func (r ReceiverFunc) Receive(msg Message, absdecimillisec int32) {
	r(msg, absdecimillisec)
}

// Receiver receives MIDI messages.
type Receiver interface {
	// Receive receives a MIDI message. absmillisec is the absolute timestamp in milliseconds
	Receive(msg Message, absmillisec int32)
}

// ErrorReceiver is a receiver that can receive errors.
type ErrorReceiver interface {
	Receiver
	OnError(error)
}

// InPorts returns the MIDI input ports
func InPorts() []string {
	ins, err := drivers.Ins()

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get midi in ports: %s\n", err.Error())
		return nil
	}

	res := make([]string, len(ins))

	for i, in := range ins {
		res[i] = in.String()
	}

	return res
}

// OutPorts returns the MIDI output ports
func OutPorts() []string {
	outs, err := drivers.Outs()

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get midi out ports: %s\n", err.Error())
		return nil
	}

	res := make([]string, len(outs))

	for i, out := range outs {
		res[i] = out.String()
	}

	return res
}
