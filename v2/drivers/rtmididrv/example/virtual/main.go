package main

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi/v2"
	//"gitlab.com/gomidi/midi/reader"
	//"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// run this in two terminals. first terminal without args to create the virtual ports and
// second terminal with argument "list" to see the ports.
func main() {
	drv, err := rtmididrv.New()
	must(err)

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs()
	must(err)

	if len(os.Args) == 2 && os.Args[1] == "list" {
		printInPorts(ins)
		printOutPorts(outs)
		return
	}

	var in midi.In
	in, err = drv.OpenVirtualIn("test-virtual-in")

	must(err)

	var out midi.Out
	out, err = drv.OpenVirtualOut("test-virtual-out")

	must(err)

	// listen for MIDI
	recv := midi.NewReceiver(func(msg midi.Message, deltamicrosec int64) {
		out.Send(msg.Data)
	}, nil)

	// example to write received messages from the virtual in port to the virtual out port
	c := make(chan int, 10)
	go in.SendTo(recv)
	//	go rd.ListenTo(in)
	<-c
}

func printPort(port midi.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts(ports []midi.In) {
	fmt.Printf("MIDI IN Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func printOutPorts(ports []midi.Out) {
	fmt.Printf("MIDI OUT Ports\n")
	for _, port := range ports {
		printPort(port)
	}
	fmt.Printf("\n\n")
}
