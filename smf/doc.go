// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
Package smf provides constants and interfaces for reading and writing of Standard MIDI Files (SMF).

The readers/writers can be found here:

  gitlab.com/gomidi/midi/smf/smfreader (read MIDI messages from SMF)
  gitlab.com/gomidi/midi/smf/smfwriter (writes MIDI messages to SMF)

The MIDI messages that can be read/written from/to a SMF file can be found here:

  gitlab.com/gomidi/midi/midimessage/channel    (Channel Messages)
  gitlab.com/gomidi/midi/midimessage/cc         (Control Change Messages)
  gitlab.com/gomidi/midi/midimessage/meta       (Meta Messages)
  gitlab.com/gomidi/midi/midimessage/sysex      (System Exclusive Messages)
*/
package smf
