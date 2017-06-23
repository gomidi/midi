// Copyright (c) 2017 Marc René Arns. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

/*
	Package smf provides constants and interfaces for reading and writing of standard midi format (SMF) files.

	The readers/writers can be found here:

	  github.com/gomidi/midi/smf/smfreader (read midi messages from standard midi file)
	  github.com/gomidi/midi/smf/smfwriter (write midi messages to standard midi file)

	The midi messages that can be read/written from/to a SMF file can be found here:

	  github.com/gomidi/midi/messages/channel    (voice/channel messages)
	  github.com/gomidi/midi/messages/cc         (control change messages)
	  github.com/gomidi/midi/messages/meta       (meta messages)
	  github.com/gomidi/midi/messages/sysex      (system exclusive messages)

	For reading there is also a comfortable handler package:

	  github.com/gomidi/midi/handler    (reading midi messages from streams or SMF files)

*/
package smf
