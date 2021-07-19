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

func rec(msg midi.Message, timestamp int32) {

	var channel, key, velocity, program, pressure uint8
	var pitch int16

	switch {
	case msg.NoteOn(&channel, &key, &velocity):
		fmt.Printf("Channel: %v key: %v %s\n", channel, key, msg)
	case msg.NoteOff(&channel, &key, &velocity):
		fmt.Printf("Channel: %v key: %v %s\n", channel, key, msg)
	case msg.AfterTouch(&channel, &pressure):
		fmt.Printf("Channel: %v Pressure: %v %s\n", channel, pressure, msg)
	case msg.ProgramChange(&channel, &program):
		fmt.Printf("Channel: %v Program: %v %s\n", channel, program, msg)
	case msg.PitchBend(&channel, &pitch, nil):
		fmt.Printf("Channel: %v Pitch: %v %s\n", channel, pitch, msg)
	case msg.Channel(&channel):
		fmt.Printf("Channel: %v %s\n", channel, msg)
	default:
		fmt.Printf("%s\n", msg)
	}
}

func main() {

	defer midi.CloseDriver()

	// allows you to get the ports when using "real" drivers like rtmididrv or portmididrv
	if len(os.Args) == 2 && os.Args[1] == "list" {
		printInPorts()
		printOutPorts()
		return
	}

	// here we take first out, for real drivers midi.OutByName should be more helpful
	s, err := midi.SenderToPort(0)
	must(err)

	// here we take first in, for real drivers midi.InByName should be more helpful
	err = midi.ListenToPort(0, midi.ReceiverFunc(rec))

	//listener, err := midi.NewListener(in, midi.ReceiverFunc(rec))

	must(err)

	//listener.Only(midi.ChannelMsg).StartListening()

	{ // write somehow MIDI
		ch := midi.Channel(0)
		err = s.Send(ch.NoteOn(60, 100))
		must(err)

		time.Sleep(time.Nanosecond)
		s.Send(ch.NoteOff(60))
		s.Send(ch.Pitchbend(-12))
		time.Sleep(time.Nanosecond)

		ch = midi.Channel(1)
		s.Send(ch.ProgramChange(12))

		s.Send(ch.NoteOn(70, 100))
		time.Sleep(time.Nanosecond)
		s.Send(ch.NoteOff(70))
		time.Sleep(time.Second * 1)
	}
}

func printInPorts() {
	fmt.Printf("MIDI IN Ports\n")
	ins := midi.InPorts()
	for i, port := range ins {
		fmt.Printf("[%v] %s\n", i, port)
	}
	fmt.Printf("\n\n")
}

func printOutPorts() {
	fmt.Printf("MIDI OUT Ports\n")
	outs := midi.OutPorts()
	for i, port := range outs {
		fmt.Printf("[%v] %s\n", i, port)
	}
	fmt.Printf("\n\n")
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
