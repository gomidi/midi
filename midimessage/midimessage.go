package midimessage

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/gomidi/midi/midimessage/syscommon"
	"github.com/gomidi/midi/midimessage/sysex"
)

// IsLive checks if msg can be send to a MIDI device (as "live" MIDI/ over the wire)
func IsLive(msg midi.Message) bool {
	if _, ok := msg.(channel.Message); ok {
		return true
	}

	if _, ok := msg.(realtime.Message); ok {
		return true
	}

	if _, ok := msg.(syscommon.Message); ok {
		return true
	}

	if _, ok := msg.(sysex.SysEx); ok {
		return true
	}

	return false

}
