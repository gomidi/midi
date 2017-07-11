package messages

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/realtime"
	"github.com/gomidi/midi/messages/syscommon"
	"github.com/gomidi/midi/messages/sysex"
)

// IsLiveMessage checks if msg can be send to a MIDI device
func IsLiveMessage(msg midi.Message) bool {
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
