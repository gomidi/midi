// Copyright (c) 2021 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package midi helps with reading and writing of MIDI messages.

The heart of this library is the `Message` type. It associates a `Type` with the raw MIDI data of the message.
When used with a `Driver` only the bytes are send and retrieved.

To process the received MIDI bytes, the `NewMessage` function creates a `Message` from the given bytes.
That `Message` has the `Type` set; that allows further examination by using its `Is` method for comparison.

The `Message` provides a method for each `Type` to retrieve the associated data.

For each `Type` there is also a corresponding function that helps creating
the MIDI data for sending or writing. For channel messages these are methods of the `Channel` type.

To listen for MIDI messages coming from an `In` port, the `Listener` provides an easy abstraction.
It allows to act on complete messages by taking care of `running status` bytes.

The `smf` subpackage helps with writing to and reading from `Simple MIDI Files` (SMF).
It seemlessly integrates with the `midi` package by using the same types.
However there are additional meta messages that could be used only in SMF files.
A meta messages is just a `Message` with certain types that belong to the meta types universe.
The `smf`package has helper functions to create such messages. In order to retrieve the associated data of a meta message,
one has to wrap the `Message` with the type alias `smf.MetaMessage` and use its functions in a similar way than for normal messages.

A message can be checked, if it is a meta message, by calling the `Is` method and passing `smf.MetaType` as parameter. 

Examples for the usage of both packages can be found in the `example` subdirectory.
Different cross plattform implementations of the `Driver` interface can be found in the `drivers` subdirectory.

The `tools` subdirectory provides command line tools to deal with MIDI data in files or one the wire.

*/
package midi
