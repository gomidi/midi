// Copyright (c) 2021 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package smf helps with reading and writing of Standard MIDI Files.

The most common time resolution for SMF files are metric ticks. They define, how many ticks a quarter note is
divided into.

A SMF has of one or more tracks. A track is a slice of events and each event has a delta in ticks to the previous event
and a message.

The smf package has its own Message type that forms the events and tracks. However it is fully transparent to the midi.Message type
and both types are used in tandem when adding messages to a track with the Add method.

A created track must be closed with its Close method. The track can then be added to the SMF, which then can be written
to an io.Writer or a file.

When reading, the tracks contain the resulting messages. The methods of the Message type can then be used to get the
different informations from the message.

There are also helper functions for playing and recording.

The TracksReader provides handy shortcuts for reading multiple tracks and also converts the time,
based on the tick resolution and the tempo changes.
*/
package smf
