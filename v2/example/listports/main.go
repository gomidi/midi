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
	fmt.Println("MIDI input ports")

	for n, in := range midi.InPorts() {
		fmt.Printf("port no. %v %q\n", n, in)
	}

	fmt.Println("\n\n")

	return nil
}

func printOutPorts() error {
	fmt.Println("MIDI output ports")

	for n, out := range midi.OutPorts() {
		fmt.Printf("port no. %v %q\n", n, out)
	}

	fmt.Println("\n\n")

	return nil
}
