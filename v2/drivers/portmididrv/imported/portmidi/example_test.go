package portmidi_test

import (
	"fmt"
	"log"
	"time"

	"github.com/rakyll/portmidi"
)

func ExampleStream_WriteShort() {
	out, err := portmidi.NewOutputStream(portmidi.DefaultOutputDeviceID(), 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Send "note on" events to play C major chord.
	out.WriteShort(0x90, 60, 100)
	out.WriteShort(0x90, 64, 100)
	out.WriteShort(0x90, 67, 100)

	// Notes will be sustained for 2 seconds.
	time.Sleep(2 * time.Second)

	// Note off events.
	out.WriteShort(0x80, 60, 100)
	out.WriteShort(0x80, 64, 100)
	out.WriteShort(0x80, 67, 100)

	out.Close()
}

func ExampleStream_WriteSysEx() {
	out, err := portmidi.NewOutputStream(portmidi.DefaultOutputDeviceID(), 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	if err = out.WriteSysEx(portmidi.Time(), "F0 0A 0A 1B 00 7F 30 F7"); err != nil {
		log.Fatal(err)
	}
}

func ExampleStream_WriteSysExBytes() {
	out, err := portmidi.NewOutputStream(portmidi.DefaultOutputDeviceID(), 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	if err = out.WriteSysExBytes(portmidi.Time(), []byte{0xF0, 0x0A, 0x0A, 0x1B, 0x00, 0x7F, 0x30, 0xF7}); err != nil {
		log.Fatal(err)
	}
}

func ExampleStream_ReadSysExBytes() {
	in, err := portmidi.NewInputStream(portmidi.DefaultInputDeviceID(), 1024)
	if err != nil {
		log.Fatal(err)
	}

	msg, err := in.Read(1024)
	if err != nil {
		log.Fatal(err)
	}

	for i, b := range msg {
		fmt.Printf("SysEx message byte %d = %02x\n", i, b)
	}
}

func ExampleStream_Poll() {
	in, err := portmidi.NewInputStream(portmidi.DefaultInputDeviceID(), 1024)
	if err != nil {
		log.Fatal(err)
	}

	result, err := in.Poll()
	if err != nil {
		log.Fatal(err)
	}

	if result {
		fmt.Println("New messages in the queue!")
	} else {
		fmt.Println("No new messages in the queue :(")
	}
}
