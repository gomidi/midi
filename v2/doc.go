// Copyright (c) 2021 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package midi helps with reading and writing of MIDI messages.

The heart of this library is the slim `Msg` type. It associates a `MsgType` with the raw MIDI data of the message.
When used with a `Driver` only the bytes are send and retrieved.

To process the received MIDI bytes, the `NewMsg` function creates a `Msg` from the given bytes.
That `Msg` has the `MsgType` set - a binary flag that allow further
examination by using its `Is` method for comparison.

The `Msg` provides methods to retrieve the associated data in a meaningful way.
Have a look at the different `MsgType`s where the meaning of type is explained and the corresponding methods of a concrete
Msg object are documented. For each `MsgType` there is also a corresponding function helps creating
the MIDI data for sending or writing. For channel messages these are methods of the `Channel` type.

To listen for MIDI messages coming from an `In` port, the `Listener` provides an easy and slim abstraction.
It allows to act on complete messages by taking care of `running status` bytes.

The `smf` subpackage helps with writing to and reading from `Simple MIDI Files` (SMF).
It seemlessly integrates with the `midi` package by using the same types.

Examples for the usage of both packages can be found in the `example` subdirectory.
Different cross plattform implementations of the `Driver` interface can be found in the `drivers` subdirectory.

The `tools` subdirectory provides command line tools to deal with MIDI data in files or one the wire.

*/
package midi
