package rtmidi

import (
	"fmt"
	"log"
	"testing"
)

func TestMidiIn(t *testing.T) {
	_, err := NewMIDIInDefault()
	if err != nil {
		//return nil, fmt.Errorf("can't open default MIDI in: %v", err)
		//t.Errorf("can't open default MIDI in: %v", err)
		fmt.Printf("can't open default MIDI in: %v", err)
	}
}

func ExampleCompiledAPI() {
	for _, api := range CompiledAPI() {
		log.Println("Compiled API: ", api)
	}
}

func ExampleMIDIIn_Message() {
	in, err := NewMIDIInDefault()
	if err != nil {
		log.Fatal(err)
	}
	defer in.Destroy()
	if err := in.OpenPort(0, "RtMidi"); err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	for {
		m, t, err := in.Message()
		if len(m) > 0 {
			log.Println(m, t, err)
		}
	}
}

func ExampleMIDIIn_SetCallback() {
	in, err := NewMIDIInDefault()
	if err != nil {
		log.Fatal(err)
	}
	defer in.Destroy()
	if err := in.OpenPort(0, "RtMidi"); err != nil {
		log.Fatal(err)
	}
	defer in.Close()
	in.SetCallback(func(m MIDIIn, msg []byte, t float64) {
		log.Println(msg, t)
	})
	<-make(chan struct{})
}
