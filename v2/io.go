package midi

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/gomidi/midi/v2/drivers"
)

// CloseDriver closes the default driver.
func CloseDriver() {
	drivers.Close()
}

// SendTo returns a function that can be used to send messages to the given midi port.
func SendTo(outPort drivers.Out) (func(msg Message) error, error) {
	if !outPort.IsOpen() {
		err := outPort.Open()
		if err != nil {
			return nil, err
		}
	}
	return func(msg Message) error {
		return outPort.Send(msg.Bytes())
	}, nil
}

type InPorts []drivers.In

func (ip InPorts) String() string {
	var bf strings.Builder

	for i, p := range ip {
		bf.WriteString(fmt.Sprintf("[%v] %s\n", i, p))
	}

	return bf.String()
}

type OutPorts []drivers.Out

func (op OutPorts) String() string {
	var bf strings.Builder

	for i, p := range op {
		bf.WriteString(fmt.Sprintf("[%v] %s\n", i, p))
	}

	return bf.String()
}

// GetInPorts returns the MIDI input ports
func GetInPorts() InPorts {
	ins, err := drivers.Ins()

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get midi in ports: %s\n", err.Error())
		return nil
	}

	return ins
}

// GetOutPorts returns the MIDI output ports
func GetOutPorts() OutPorts {
	outs, err := drivers.Outs()

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get midi out ports: %s\n", err.Error())
		return nil
	}

	return outs
}
