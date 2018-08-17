// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package midiwriter provides a writer for live/streaming/"over the wire" MIDI data.

Usage


	import (
		"github.com/gomidi/midi/midiwriter"
		. "github.com/gomidi/midi/midimessage/channel"     // (Channel Messages)
		. "time"

		// you may also want to use these
		// github.com/gomidi/midi/midimessage/realtime   (System Realtime Messages)
		// github.com/gomidi/midi/midimessage/cc         (ControlChange Messages)
		// github.com/gomidi/midi/midimessage/syscommon  (System Common Messages)
		// github.com/gomidi/midi/midimessage/sysex      (system exclusive messages)
	)

	// given some output
	var output io.Writer

	wr := midiwriter.New(output)

	// simulates pressing down key 65 on MIDI channel 3 with velocity 90
	// MIDI channels 1-16 correspond to channel.Channel0 - channel.Channel15.
	wr.Write(Channel2.NoteOn(65, 90))

	// simulates keep pressing for 1 sec
	Sleep(Second)

	// simulates releasing key 65 on MIDI channel 3
	wr.Write(Channel2.NoteOff(65))

*/
package midiwriter
