// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Copyright (c) 2020 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package reader provides an easy abstraction for reading of cable MIDI and SMF (Standard MIDI File) data.

First create new Reader with callbacks. To read a SMF file use ReadSMFFile. To read cable MIDI,
use the ListenTo method with an input port.

See the examples in the top level examples folder.

To connect with the MIDI ports of your computer use it with
the adapter package for rtmidi (gitlab.com/gomidi/rtmididrv) or portmidi (gitlab.com/gomidi/portmididrv).
*/
package reader
