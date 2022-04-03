// Copyright (c) 2021 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package midi helps with reading and writing of MIDI messages.

The heart of this library is the `Message` type. It is simply a slice of bytes with some methods to scan the data.
Any complete midi message (i.e. without "running status") can be interpreted as a Message, by simply converting the type.
Then its data can be retrieved by simply calling the corresponding Scan* method.

Example

   // this are the raw bytes for a noteon message on channel 1 (2nd channel) for the key 60 with velocity of 120
   var b = []byte{0x91, 0x3C, 0x78}

   // convert to Message type
   msg := midi.Message(b)

   var channel, key, velocity uint8
   if msg.ScanNoteOn(&channel, &key, &velocity) {
     fmt.Printf("got %s: channel: %v key: %v, velocity: %v\n", msg.Type(), channel, key, velocity)
   }

Received messages can be categorized via their types, e.g.

  switch msg.Type() {
  case midi.NoteOnMsg, midi.NoteOffMsg:
    // do something
  case midi.ControlChangeMsg:
    // do some other thing
  default:
    if msg.Is(midi.RealTimeMsg) || msg.Is(midi.SysCommonMsg) || msg.Is(midi.SysExMsg) {
      // ignore
    }
  }

A new message is created by the corresponding function, e.g.

   msg := midi.NoteOn(1, 60, 120) // channel, key, velocity
   fmt.Printf("% X\n", []byte(msg)) // prints 91 3C 78

Sending and retrieving of midi data is done via drivers. Most of the time, the use of a single driver is sufficient, since
it can provide access to all available midi ports on a system.

Therefor it is handy, that the drivers autoregister themself when including the corresponding package.
This also makes it easy to exchange a driver when needed.

Complete Example (for readability errors are not handled here):

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
	)

	func main() {

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
	}

	func eachMessage(msg midi.Message, timestampms int32) {
		if msg.Is(midi.RealTimeMsg) {
			// ignore realtime messages
			return
		}
		var channel, key, velocity uint8
		switch {
		// is better, than to use ScanNoteOn, since note on messages with velocity of 0 also stop notes
		case msg.ScanNoteStart(&channel, &key, &velocity):
			fmt.Printf("note started at %vms channel: %v key: %v velocity: %v\n", timestampms, channel, key, velocity)
		// is better, than to use ScanNoteOff, since note on messages with velocity of 0 also stop notes
		case msg.ScanNoteEnd(&channel, &key, &velocity):
			fmt.Printf("note ended at %vms channel: %v key: %v\n", timestampms, channel, key)
		default:
			fmt.Printf("received %s at %vms\n", msg, timestampms)
		}
	}



The `smf` subpackage helps with writing to and reading from `Simple MIDI Files` (SMF).

Examples for the usage of both packages can be found in the `example` subdirectory.

Different cross plattform implementations of the `Driver` interface can be found in the `drivers` subdirectory.

The `tools` subdirectory provides command line tools to deal with MIDI data in files or one the wire.

*/
package midi
