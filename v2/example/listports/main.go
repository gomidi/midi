package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	defer midi.CloseDriver()

	err := run()

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

func run() error {
	if err := printInPorts(); err != nil {
		return err
	}

	if err := printOutPorts(); err != nil {
		return err
	}

	return nil
}

func printInPorts() error {
	ins, err := midi.Ins()
	if err != nil {
		return err
	}

	fmt.Println("MIDI input ports")

	for _, in := range ins {
		fmt.Printf("port no. %v %q\n", in.Number(), in.String())
	}

	fmt.Println("\n\n")

	return nil
}

func printOutPorts() error {
	outs, err := midi.Outs()
	if err != nil {
		return err
	}

	fmt.Println("MIDI output ports")

	for _, out := range outs {
		fmt.Printf("port no. %v %q\n", out.Number(), out.String())
	}

	fmt.Println("\n\n")

	return nil
}
