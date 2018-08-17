// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package midi provides interfaces for reading and writing of MIDI messages.

Since they are handled slightly different, this packages introduces the terminology of
"live" MIDI reading/writing for dealing with MIDI messages as "over the wire" (in realtime)
as opposed to smf MIDI reading/writing to Standard MIDI Files (SMF).

However both variants can be used with io.Writer and io.Reader and can thus be "streamed".

This package provides a Reader and Writer interface that is common to live and SMF MIDI handling.
This should allow to easily develop transformations (e.g. quantization,
filtering) that may be used in both cases.

If you want a comfortable common package providing everything at a high level, use the
porcelain package

  github.com/gomidi/mid

The underlying core implementations can be found here:

  github.com/gomidi/midi/midireader (live reading)
  github.com/gomidi/midi/midiwriter (live writing)
  github.com/gomidi/midi/smf/smfreader   (SMF reading)
  github.com/gomidi/midi/smf/smfwriter   (SMF writing)
  github.com/gomidi/midi/smf/smftrack    (SMF modification)

The core of the MIDI messages that can be written or analyzed can be found here:

  github.com/gomidi/midi/midimessage/channel    (Channel Messages)
  github.com/gomidi/midi/midimessage/meta       (Meta Messages)
  github.com/gomidi/midi/midimessage/realtime   (System Realtime Messages)
  github.com/gomidi/midi/midimessage/syscommon  (System Common messages)
  github.com/gomidi/midi/midimessage/sysex      (System Exclusive messages)

Please keep in mind that that not all kinds of MIDI messages can be used in both scenarios.

System Realtime and System Common Messages are restricted to "over the wire",
while Meta Messages are restricted to SMF files. However System Realtime and System Common Messages
can be saved inside a SMF file which the help of SysEx escaping (F7).

*/
package midi
