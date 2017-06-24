// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
	Package midiwriter provides a writer for live MIDI data.

	Usage


		import (
			"github.com/gomidi/midi/live/midiwriter"
			"github.com/gomidi/midi/messages/channel"     // (Channel Messages)
			"time"

			// you may also want to use these
			// github.com/gomidi/midi/messages/realtime   (System Realtime Messages)
			// github.com/gomidi/midi/messages/cc         (ControlChange Messages)
			// github.com/gomidi/midi/messages/syscommon  (System Common Messages)
			// github.com/gomidi/midi/messages/sysex      (system exclusive messages)
		)

		// given some output
		var output io.Writer

		wr := midiwriter.New(output)

		// simulates pressing down key 65 on MIDI channel 3 with velocity 90
		// MIDI channels 1-16 correspond to channel.Ch0 - channel.Ch15.
		wr.Write(channel.Ch2.NoteOn(65, 90))

		// simulates keep pressing for 1 sec
		time.Sleep(time.Second)

		// simulates releasing key 65 on MIDI channel 3
		wr.Write(channel.Ch2.NoteOff(65))

*/
package midiwriter
