package channel

// TODO do with iota
const (
	// MIDI channel 1
	Channel0 = Channel(0)

	// MIDI channel 2
	Channel1 = Channel(1)

	// MIDI channel 3
	Channel2 = Channel(2)

	// MIDI channel 4
	Channel3 = Channel(3)

	// MIDI channel 5
	Channel4 = Channel(4)

	// MIDI channel 6
	Channel5 = Channel(5)

	// MIDI channel 7
	Channel6 = Channel(6)

	// MIDI channel 8
	Channel7 = Channel(7)

	// MIDI channel 9
	Channel8 = Channel(8)

	// MIDI channel 10
	Channel9 = Channel(9)

	// MIDI channel 11
	Channel10 = Channel(10)

	// MIDI channel 12
	Channel11 = Channel(11)

	// MIDI channel 13
	Channel12 = Channel(12)

	// MIDI channel 14
	Channel13 = Channel(13)

	// MIDI channel 15
	Channel14 = Channel(14)

	// MIDI channel 16
	Channel15 = Channel(15)
)

func SetChannel(msg Message, ch uint8) Message {
	if ch > 15 {
		panic("invalid channel number")
	}

	c := Channel(ch)

	switch v := msg.(type) {
	case Aftertouch:
		return c.Aftertouch(v.Pressure())
	case ControlChange:
		return c.ControlChange(v.Controller(), v.Value())
	case NoteOn:
		return c.NoteOn(v.Key(), v.Velocity())
	case NoteOff:
		return c.NoteOff(v.Key())
	case NoteOffVelocity:
		return c.NoteOffVelocity(v.Key(), v.Velocity())
	case Pitchbend:
		return c.Pitchbend(v.Value())
	case PolyAftertouch:
		return c.PolyAftertouch(v.Key(), v.Pressure())
	case ProgramChange:
		return c.ProgramChange(v.Program())
	}

	panic("unreachable")
}

// Channel represents a MIDI channel
// there must not be more than 16 MIDI channels (0-15)
type Channel uint8

// Channel returns the number of the MIDI channel (0-15)
func (c Channel) Channel() uint8 {
	return uint8(c)
}

// NoteOff creates a note-off message on the channel for the given key
// The note-off message is "faked" by a note-on message of velocity 0.
// This allows saving bandwidth by using running status.
// If you need a "real" note-off message with velocity, use NoteOffVelocity.
func (c Channel) NoteOff(key uint8) NoteOff {
	if key > 127 {
		key = 127
	}
	return NoteOff{channel: c.Channel(), key: key}
}

// NoteOffVelocity creates a note-off message with velocity on the channel.
func (c Channel) NoteOffVelocity(key uint8, velocity uint8) NoteOffVelocity {
	if key > 127 {
		key = 127
	}
	if velocity > 127 {
		velocity = 127
	}
	return NoteOffVelocity{NoteOff{channel: c.Channel(), key: key}, velocity}
}

// NoteOn creates a note-on message on the channel
func (c Channel) NoteOn(key uint8, velocity uint8) NoteOn {
	if key > 127 {
		key = 127
	}
	if velocity > 127 {
		velocity = 127
	}
	return NoteOn{channel: c.Channel(), key: key, velocity: velocity}
}

// PolyAftertouch creates a polyphonic aftertouch message on the channel
func (c Channel) PolyAftertouch(key uint8, pressure uint8) PolyAftertouch {
	if key > 127 {
		key = 127
	}
	if pressure > 127 {
		pressure = 127
	}
	return PolyAftertouch{channel: c.Channel(), key: key, pressure: pressure}
}

// ControlChange creates a control change message on the channel
func (c Channel) ControlChange(controller uint8, value uint8) ControlChange {
	if controller > 127 {
		controller = 127
	}
	if value > 127 {
		value = 127
	}
	return ControlChange{channel: c.Channel(), controller: controller, value: value}
}

// ProgramChange creates a program change message on the channel
func (c Channel) ProgramChange(program uint8) ProgramChange {
	if program > 127 {
		program = 127
	}
	return ProgramChange{channel: c.Channel(), program: program}
}

// Aftertouch creates an aftertouch message on the channel
func (c Channel) Aftertouch(pressure uint8) Aftertouch {
	if pressure > 127 {
		pressure = 127
	}
	return Aftertouch{channel: c.Channel(), pressure: pressure}
}

// Pitchbend creates a pitch bend message on the channel
func (c Channel) Pitchbend(value int16) Pitchbend {
	if value < PitchLowest {
		value = PitchLowest
	}

	if value > PitchHighest {
		value = PitchHighest
	}
	return Pitchbend{channel: c.Channel(), value: value}
}
