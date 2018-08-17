// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package smfreader provides a reader of Standard MIDI Files (SMF).

Usage

	import (
		"github.com/gomidi/midi/smf/smfreader"
		"github.com/gomidi/midi/midimessage/channel"    // (Channel Messages)

		// you may also want to use these
		// github.com/gomidi/midi/midimessage/meta       (Meta Messages)
		// github.com/gomidi/midi/midimessage/cc         (Control Change Messages)
		// github.com/gomidi/midi/midimessage/syscommon  (System Common Messages)
		// github.com/gomidi/midi/midimessage/sysex      (System Exclusive Messages)
	)

	var err1, err2 error

	readMIDI := func (rd smf.Reader) {

		var m midi.Message

		for {
			m, err1 = rd.Read()

			// to interrupt, the input.Read method must return io.EOF or any other error
			if err1 != nil {
				break
			}

			// deal with them based on a type switch
			switch msg := m.(type) {
			case channel.NoteOn:
				fmt.Printf(
				  "NoteOn at channel %v: key %v velocity: %v\n",
				  msg.Channel(), // MIDI channels 1-16 correspond to msg.Channel 0-15
				  msg.Key(),
				  msg.Velocity(),
				)
			case channel.NoteOff:
				...
			}
		}

	}

	err2 = smfreader.ReadFile("file.mid", readMIDI)

	// deal with err1 and err2

*/
package smfreader
