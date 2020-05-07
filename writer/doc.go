// Copyright (c) 2020 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package writer provides an easy abstraction for writing of cable MIDI and SMF (Standard MIDI File) data.

Create a new Writer to write to a MIDI cable (an io.Writer). Use WriteSMF to write a complete SMF file.

See the examples in the top level examples folder.

To connect with the MIDI ports of your computer use it with
the adapter package for rtmidi (gitlab.com/gomidi/rtmididrv) or portmidi (gitlab.com/gomidi/portmididrv).
*/
package writer
