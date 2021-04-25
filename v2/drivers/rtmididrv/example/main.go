package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"
	//"gitlab.com/gomidi/midi/reader"
	//"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	// when using portmidi, replace the line above with
	// driver gitlab.com/gomidi/portmididrv
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// This example expects the first input and output port to be connected
// somehow (are either virtual MIDI through ports or physically connected).
// We write to the out port and listen to the in port.
func main() {
	drv, err := driver.New()
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

	in, out := ins[0], outs[0]

	must(in.Open())
	must(out.Open())

	// listen for MIDI
	recv := midi.NewReceiver(func(msg midi.Message, deltamicrosec int64) {
		fmt.Printf("@%v: %s\n", deltamicrosec, msg)
	}, nil)

	go in.SendTo(recv)

	{ // write MIDI to out that passes it to in on which we listen.
		ch := midi.Channel(0)
		err := out.Send(ch.NoteOn(60, 100))

		//err := writer.NoteOn(wr, 60, 100)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Nanosecond)
		out.Send(ch.NoteOff(60))
		//writer.NoteOff(wr, 60)
		time.Sleep(time.Nanosecond)

		ch = midi.Channel(1)

		//writer.NoteOn(wr, 70, 100)
		out.Send(ch.NoteOn(70, 100))
		time.Sleep(time.Nanosecond)
		//writer.NoteOff(wr, 70)
		out.Send(ch.NoteOff(70))
		time.Sleep(time.Second * 1)
	}
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
