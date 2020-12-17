package filter

import (
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
)

// Or returns a filter that is true if any of the given filters is true
func Or(filters ...Filter) Filter {
	return func(msg midi.Message) bool {
		for _, f := range filters {
			if f(msg) {
				return true
			}
		}
		return false
	}
}

// And returns a filter that is true if all of the given filters are true
func And(filters ...Filter) Filter {
	return func(msg midi.Message) bool {
		for _, f := range filters {
			if !f(msg) {
				return false
			}
		}
		return true
	}
}

// Channel returns a Filter that triggers only for messages on the given MIDI channel.
// If ch is < 0, any channel will trigger.
func Channel(ch int8) Filter {
	return func(msg midi.Message) bool {
		chmsg, is := msg.(channel.Message)

		if !is {
			return false
		}
		return (ch < 0) || uint8(ch) == chmsg.Channel()
	}
}

// Filter detects  metronome beats on the metronome track.
// It is a function that returns true when a metronome beat was detected.
type Filter func(msg midi.Message) bool

// NoteOn returns a Filter that triggers only for the given key.
// If key is < 0, any key will trigger.
func NoteOn(key int8) Filter {
	return func(msg midi.Message) bool {
		switch nt := msg.(type) {
		case channel.NoteOn:
			if nt.Velocity() == 0 {
				return false
			}
			return key < 0 || uint8(key) == nt.Key()
		default:
			return false
		}
	}
}

// NoteOff returns a Filter that triggers only for the given key.
// If key is < 0, any key will trigger.
func NoteOff(key int8) Filter {
	return func(msg midi.Message) bool {
		switch nt := msg.(type) {
		case channel.NoteOn:
			if nt.Velocity() > 0 {
				return false
			}
			return key < 0 || uint8(key) == nt.Key()
		case channel.NoteOff:
			return key < 0 || uint8(key) == nt.Key()
		case channel.NoteOffVelocity:
			return key < 0 || uint8(key) == nt.Key()
		default:
			return false
		}
	}
}

// NoteOffVelocity returns a Filter that triggers only for the given key.
// If key is < 0, any key will trigger.
func NoteOffVelocity(key int8) Filter {
	return func(msg midi.Message) bool {
		switch nt := msg.(type) {
		case channel.NoteOffVelocity:
			if nt.Velocity() == 0 {
				return false
			}
			return key < 0 || uint8(key) == nt.Key()
		default:
			return false
		}
	}
}

// CC returns a Filter that triggers only for the given controller.
// If controller is < 0, any controller will trigger.
func CC(controller int8) Filter {
	return func(msg midi.Message) bool {
		switch c := msg.(type) {
		case channel.ControlChange:
			if c.Value() == 0 {
				return false
			}
			return controller < 0 || uint8(controller) == c.Controller()
		default:
			return false
		}
	}
}

// Aftertouch returns a Filter that triggers only for the aftertouch messages.
func Aftertouch() Filter {
	return func(msg midi.Message) bool {
		switch at := msg.(type) {
		case channel.Aftertouch:
			return at.Pressure() != 0
		default:
			return false
		}
	}
}

// PolyAftertouch returns a Filter that triggers only for the polyaftertouch messages of the given key.
// If key is < 0, any key will trigger.
func PolyAftertouch(key int8) Filter {
	return func(msg midi.Message) bool {
		switch pa := msg.(type) {
		case channel.PolyAftertouch:
			if pa.Pressure() == 0 {
				return false
			}
			return key < 0 || uint8(key) == pa.Key()
		default:
			return false
		}
	}
}

// Pitchbend returns a Filter that triggers only for the pitchbend messages > 0.
func Pitchbend() Filter {
	return func(msg midi.Message) bool {
		switch pb := msg.(type) {
		case channel.Pitchbend:
			return pb.Value() > 0
		default:
			return false
		}
	}
}
