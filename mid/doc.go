// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package mid provides an easy abstraction for reading and writing of "live" MIDI and SMF (Standard MIDI File) data.

MIDI data could be written the following ways:

	- WriteTo writes "live" MIDI to an OutConnection, aka MIDI out port
	- NewWriter is used to write "live" MIDI to an io.Writer.
	- NewSMF is used to write SMF MIDI to an io.Writer.
	- NewSMFFile is used to write a complete SMF file.

To read, create a Reader and attach callbacks to it.
Then MIDI data could be read the following ways:

	- Reader.ReadFrom reads "live" MIDI from an InConnection, aka MIDI in port
	- Reader.Read reads "live" MIDI from an io.Reader.
	- Reader.ReadSMF reads SMF MIDI from an io.Reader.
	- Reader.ReadSMFFile reads a complete SMF file.

For a simple example with "live" MIDI and io.Reader and io.Writer see examples/simple/simple_test.go.

To connect with the MIDI ports of your computer (via mid.In and mid.Out), use it with
the adapter package for rtmidi (gitlab.com/gomidi/rtmididrv) or portmidi (gitlab.com/gomidi/portmididrv).

In the README.md you can find a simple example how to do it.
*/
package mid
