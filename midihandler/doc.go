// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
  Package midihandler provides an easy abstraction for reading of MIDI data live or from SMF files.

	It provides a common handler for live reading and SMF reading, so that code can be shared in both use cases.
	The user attaches callback functions to the handler and they get invoked as the MIDI data is read.

	However there are some differences between the two:
	  - SMF events have a message and a delta time while live MIDI just has messages
	  - a SMF file has a header
	  - a SMF file may have meta messages (e.g. track name, copyright etc) while live MIDI must not have them
	  - live MIDI may have realtime and syscommon messages while a SMF must not have them

	In sum the only MIDI message handling that could be shared in practise is the handling of channel messages
	and sysex messages.

	That is reflected in the type signature of the callbacks and results in three kinds of callbacks:
	  - The ones without a `midihandler.SMFPosition` are called for messages that may only appear live
	  - The ones with a `midihandler.SMFPosition` are called for messages that may only appear
	    within a SMF file
	  - The onse with a `*midihandler.SMFPosition` are called for messages that may appear live
	    and within a SMF file. In a SMF file pointer is never `nil` while in a live situation it is always `nil`.

  Due to the nature of a SMF file, the tracks have no "assigned numbers", and can in the worst case just be
  distinguished by the order in which they appear inside the file. In addition the delta times are just differences
  to the previous message inside the track. That means, if the user is just interested in some kind of MIDI messages
  (e.g. note on and note off), he still has to deal with the delta times of all other messages, because they affect
  the absolute timing of the following messages. Also he has to reset the delta time at the beginning of each track.
  This is error prone and therefor `midihandler.SMFPosition` provides a helper that not just contains he original delta
  (which shouldn't be needed most of the time), but also the number (order) of the track in which the message appeared
  and the absolute position (in ticks). So the user can focus just on the messages and tracks he is interested in.

	See the example for a handler that handles both live and SMF messages.
*/
package midihandler
