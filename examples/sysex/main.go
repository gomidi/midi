package main

import (
	"fmt"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/testdrv" // autoregisters driver
	"gitlab.com/gomidi/midi/v2/sysex"
)

func main() {
	defer midi.CloseDriver()

	stop, err := midi.ListenTo(0, func(msg midi.Message, timestampms int32) {
		var bt []byte
		switch {
		case msg.GetSysEx(&bt):
			fmt.Printf("got sysex: % X\n", bt)
		default:
			// ignore
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	send, err := midi.SendTo(0)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	reset := sysex.GMReset.SysEx()
	fmt.Printf("sending reset:\n% X,\n", reset)

	send(reset)

	time.Sleep(time.Second * 1)

	stop()
}
