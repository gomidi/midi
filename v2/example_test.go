package midi_test

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
)

func Example() {

	var eachMessage = func(msg midi.Message, timestampms int32) {
		if msg.Is(midi.RealTimeMsg) {
			// ignore realtime messages
			return
		}
		var channel, key, velocity uint8
		switch {
		// is better, than to use GetNoteOn, since note on messages with velocity of 0 also stop notes
		case msg.GetNoteStart(&channel, &key, &velocity):
			fmt.Printf("note started at %vms channel: %v key: %v velocity: %v\n", timestampms, channel, key, velocity)
		// is better, than to use GetNoteOff, since note on messages with velocity of 0 also stop notes
		case msg.GetNoteEnd(&channel, &key, &velocity):
			fmt.Printf("note ended at %vms channel: %v key: %v\n", timestampms, channel, key)
		default:
			fmt.Printf("received %s at %vms\n", msg, timestampms)
		}
	}

	// always good to close the driver at the end
	defer midi.CloseDriver()

	// allows you to get the ports when using "real" drivers like rtmididrv or portmididrv
	if len(os.Args) == 2 && os.Args[1] == "list" {
		fmt.Printf("MIDI IN Ports\n")
		for i, port := range midi.InPorts() {
			fmt.Printf("no: %v %q\n", i, port)
		}
		fmt.Printf("\n\nMIDI OUT Ports\n")
		for i, port := range midi.OutPorts() {
			fmt.Printf("no: %v %q\n", i, port)
		}
		fmt.Printf("\n\n")
		return
	}

	var out int = 0
	// here we take first out, for real drivers the following should be more helpful
	// var out = midi.OutByName("my synth")

	// creates a sender function to the midi out port
	send, _ := midi.SendTo(out)

	var in int = 0
	// here we take first in, for real drivers the following should be more helpful
	// var in = midi.InByName("my midi keyboard")

	// listens to the midi in port and calls the callback function eachMessage for each
	// message. Note, that any running status bytes are converted and only complete messages
	// are passed to the callback.
	stop, _ := midi.ListenTo(in, eachMessage)

	{ // send some MIDI via the sender
		send(midi.NoteOn(0, 60, 100))
		time.Sleep(time.Millisecond * 30)
		send(midi.NoteOff(0, 60))
		send(midi.Pitchbend(0, -12))
		time.Sleep(time.Millisecond * 20)
		send(midi.ProgramChange(1, 12))
	}

	// stops listening
	stop()

	// Output:
	// note started at 0ms channel: 0 key: 60 velocity: 100
	// note ended at 30ms channel: 0 key: 60
	// received PitchBend channel: 0 pitch: -12 (8180) at 30ms
	// received ProgramChange channel: 1 program: 12 at 50ms

}
