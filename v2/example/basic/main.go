package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"

	// testdrv has one in port and one out port which is connected to the in port
	// which works fine for this example
	_ "gitlab.com/gomidi/midi/v2/drivers/testdrv"
	// when using rtmidi, replace the line above with
	//_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	// when using portmidi, replace the line above with
	//_ gitlab.com/gomidi/midi/v2/drivers/portmididrv
)

func main() {

	defer midi.CloseDriver()

	// allows you to get the ports when using "real" drivers like rtmididrv or portmididrv
	if len(os.Args) == 2 && os.Args[1] == "list" {
		printInPorts()
		printOutPorts()
		return
	}

	// here we take first out, for real drivers midi.OutByName should be more helpful
	out, err := midi.OutByNumber(0)
	must(err)
	out.Open()

	// here we take first in, for real drivers midi.InByName should be more helpful
	in, err := midi.InByNumber(0)
	must(err)

	listener, err := midi.NewListener(in, func(msg midi.Message, deltamicrosec int64) {
		switch {
		case msg.Is(midi.NoteMsg):
			fmt.Printf("Channel: %v key: %v %s\n", msg.Channel(), msg.Key(), msg)
		case msg.IsOneOf(midi.AfterTouchMsg, midi.PolyAfterTouchMsg):
			fmt.Printf("Channel: %v Pressure: %v %s\n", msg.Channel(), msg.Pressure(), msg)
		case msg.Is(midi.ProgramChangeMsg):
			fmt.Printf("Channel: %v Program: %v\n", msg.Channel(), msg.Program())
		case msg.Is(midi.PitchBendMsg):
			rel, _ := msg.Pitch()
			fmt.Printf("Channel: %v Pitch: %v\n", msg.Channel(), rel)
		default:
			fmt.Printf("Channel: %v %s\n", msg.Channel(), msg)
		}
	})

	must(err)
	listener.Only(midi.ChannelMsg).StartListening()

	{ // write somehow MIDI
		ch := midi.Channel(0)
		err = out.Send(ch.NoteOn(60, 100))
		must(err)

		time.Sleep(time.Nanosecond)
		out.Send(ch.NoteOff(60))
		out.Send(ch.Pitchbend(-12))
		time.Sleep(time.Nanosecond)

		ch = midi.Channel(1)
		out.Send(ch.ProgramChange(12))

		out.Send(ch.NoteOn(70, 100))
		time.Sleep(time.Nanosecond)
		out.Send(ch.NoteOff(70))
		time.Sleep(time.Second * 1)
	}
}

func printPort(port midi.Port) {
	fmt.Printf("[%v] %s\n", port.Number(), port.String())
}

func printInPorts() {
	fmt.Printf("MIDI IN Ports\n")
	ins, err := midi.Ins()
	must(err)
	for _, port := range ins {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func printOutPorts() {
	fmt.Printf("MIDI OUT Ports\n")
	outs, err := midi.Outs()
	must(err)
	for _, port := range outs {
		printPort(port)
	}
	fmt.Printf("\n\n")
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
