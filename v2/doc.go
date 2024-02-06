// Copyright (c) 2021 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package midi helps with reading and writing of MIDI messages.

A `Message` is a slice of bytes with some methods to get the data.
Any complete midi message (i.e. without "running status") can be interpreted as a Message, by converting the type.
The data can be retrieved with the corresponding Get* method.

	// the raw bytes for a noteon message on channel 1 (2nd channel) for the key 60 with velocity of 120
	var b = []byte{0x91, 0x3C, 0x78}

	// convert to Message type
	msg := midi.Message(b)

	var channel, key, velocity uint8
	if msg.GetNoteOn(&channel, &key, &velocity) {
	  fmt.Printf("got %s: channel: %v key: %v, velocity: %v\n", msg.Type(), channel, key, velocity)
	}

Received messages can be categorized by their type, e.g.

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

Sending and retrieving of midi data is done via drivers. The drivers take care of converting a "running status" into a full status.
Most of the time, a single driver is sufficient.
Therefor it is handy, that the drivers autoregister themself. This also makes it easy to exchange a driver if needed (see the example).

Different cross plattform implementations of the `Driver` interface can be found in the `drivers` subdirectory.

The `smf` subpackage helps with writing to and reading from `Simple MIDI Files` (SMF) (see https://pkg.go.dev/gitlab.com/gomidi/midi/v2/smf).

The `tools` subdirectory provides command line tools and libraries based on this library.
*/
package midi
