// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
	Package smfwriter provides a writer for Standard MIDI Files (SMF).

	Usage

		import (
			"github.com/gomidi/midi/smf"
			"github.com/gomidi/midi/smf/smfwriter"
			"github.com/gomidi/midi/messages/meta"    // (Meta Messages)
			"github.com/gomidi/midi/messages/channel" // (Channel Messages)

			// you may also want to use these
			// github.com/gomidi/midi/messages/cc         (ControlChange Messages)
			// github.com/gomidi/midi/messages/sysex      (System Exclusive Messages)
		)

		var err1, err2 error

		writeMIDI := func (wr smf.Writer) {

			// always set the delta before writing
			// delta defaults to 960 ticks per quarter note
			wr.SetDelta(480)

			// starts MIDI note 65 on MIDI channel 3 with velocity 90 with delta of 480 to
			// the beginning of the track (note starts after a quaver pause)
			// MIDI channels 1-16 correspond to channel.Ch0 - channel.Ch15.
			_, err1 = wr.Write(channel.Ch2.NoteOn(65, 90))

			if err1 != nil {
				return
			}

			wr.SetDelta(960)

			// stops MIDI note 65 on MIDI channel 3 with delta of 960 to previous message
			// this results in a duration of 1 quarter note for midi note 65
			_, err1 = wr.Write(channel.Ch2.NoteOff(65))

			if err1 != nil {
				return
			}

			// finishes the first track and writes it to the file
			_, err1 = wr.Write(meta.EndOfTrack)

			if err1 != nil {
				return
			}

			// the next write writes to the second track
			// after writing delta is always 0, so we start here at the beginning of the second track
			_, err1 = wr.Write(meta.Text("hello second track!"))

			if err1 != nil {
				return
			}

			// finishes the second track and writes it to the file
			_, err1 = wr.Write(meta.EndOfTrack)
		}

		// the number passed to the NumTracks option must match the tracks written
		// if NumTracks is not passed, it defaults to single track (SMF0)
		// if numtracks > 1, SMF1 format is chosen.
		err2 = smfwriter.WriteFile("file.mid", writeMIDI, smfwriter.NumTracks(2))

		// deal with err1 and err2

*/
package smfwriter
