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

In sum the only MIDI message handling that could be shared in practice is the handling of channel
and sysex messages.

That is reflected in the type signature of the callbacks and results in three kinds of callbacks:
  - The ones that don't receive a SMFPosition are called for messages that may only appear live
  - The ones with that receive a SMFPosition are called for messages that may only appear
    within a SMF file
  - The onse with that receive a *SMFPosition are called for messages that may appear live
    and within a SMF file. In a SMF file the pointer is never nil while in a live situation it always is.

Due to the nature of a SMF file, the tracks have no "assigned numbers", and can in the worst case just be
distinguished by the order in which they appear inside the file.

In addition the delta times are just relative to the previous message inside the track.

This has some severe consequences:

  - Eeven if the user is just interested in some kind of messages (e.g. note on and note off),
    he still has to deal with the delta times of every message.
  - At the beginning of each track he has to reset the tracking time.

This procedure is error prone and therefor SMFPosition provides a helper that contains not just
the original delta (which shouldn't be needed most of the time), but also the absolute position (in ticks)
and order number of the track in which the message appeared.

This way the user can ignore the messages and tracks he is not interested in.

See the example for a handler that handles both live and SMF messages.
*/
package midihandler
