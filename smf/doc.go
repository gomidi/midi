// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
	Package smf provides constants and interfaces for reading and writing of Standard MIDI Files (SMF).

	The readers/writers can be found here:

	  github.com/gomidi/midi/smf/smfreader (read MIDI messages from SMF)
	  github.com/gomidi/midi/smf/smfwriter (writes MIDI messages to SMF)

	The MIDI messages that can be read/written from/to a SMF file can be found here:

	  github.com/gomidi/midi/messages/channel    (Channel Messages)
	  github.com/gomidi/midi/messages/cc         (Control Change Messages)
	  github.com/gomidi/midi/messages/meta       (Meta Messages)
	  github.com/gomidi/midi/messages/sysex      (System Exclusive Messages)

	For reading there is also a comfortable handler package:

	  github.com/gomidi/midi/handler    (reading MIDI messages live or from SMF files)

*/
package smf
