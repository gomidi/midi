package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"

	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
)

func main() {
	err := run()

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}

func run() error {

	defer midi.CloseDriver()
	in := midi.FindInPort("VMPK")

	if in < 0 {
		return fmt.Errorf("can't find MIDI in port %q", "VMPK")
	}

	stop, err := smf.RecordTo(in, 120, "recordedx.mid")

	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	return stop()
}
