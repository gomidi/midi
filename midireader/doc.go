// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package midireader provides a reader for live/streaming/"over the wire" MIDI data.

Usage

	import (
		"github.com/gomidi/midi/midireader"
		. "github.com/gomidi/midi/midimessage/channel"    // (Channel Messages)
		"github.com/gomidi/midi/midimessage/realtime"   // (System Realtime Messages)

		// you may also want to use these
		// github.com/gomidi/midi/midimessage/cc         (Control Change Messages)
		// github.com/gomidi/midi/midimessage/syscommon  (System Common Messages)
		// github.com/gomidi/midi/midimessage/sysex      (System Exclusive Messages)
	)

	// given some MIDI input
	var input io.Reader

	// create a callback for realtime messages
	rthandler := func(m realtime.Message) {
		// deal with it
		if m == realtime.Start {
			...
		}
	}

	rd := midireader.New(input), rthandler)

	// everything but realtime messages, since they are covered by rthandler
	var m midi.Message
	var err error

	for {
		m, err = rd.Read()

		// to interrupt, the input.Read method must return io.EOF or any other error
		if err != nil {
			break
		}

		// deal with them based on a type switch
		switch msg := m.(type) {
		case NoteOn:
			fmt.Printf(
			  "NoteOn at channel %v: key %v velocity: %v\n",
			  msg.Channel(), // MIDI channels 1-16 correspond to msg.Channel 0-15
			  msg.Key(),
			  msg.Velocity(),
			)
		case NoteOff:
			...
		}
	}

	if err != io.EOF {
	  // real error here
	}


*/
package midireader
