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

// SendTo returns a function that can be used to send messages to the given midi port.
func SendTo(portno int) (func(msg Message) error, error) {
	out, err := drivers.OutByNumber(portno)
	if err != nil {
		return nil, err
	}
	if !out.IsOpen() {
		err = out.Open()
		if err != nil {
			return nil, err
		}
	}
	return func(msg Message) error {
		return out.Send(msg.Bytes())
	}, nil
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
