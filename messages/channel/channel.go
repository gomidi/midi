package channel

// New returns a new channel. no must be <= 15, since there are 16 midi channels.
func New(no uint8) Channel {
	if no > 15 {
		panic("midi channel > 15 is not allowed")
	}

	return Channel{no}
}

var (
	// MIDI channel 1
	Ch0 = New(0)

	// MIDI channel 2
	Ch1 = New(1)

	// MIDI channel 3
	Ch2 = New(2)

	// MIDI channel 4
	Ch3 = New(3)

	// MIDI channel 5
	Ch4 = New(4)

	// MIDI channel 6
	Ch5 = New(5)

	// MIDI channel 7
	Ch6 = New(6)

	// MIDI channel 8
	Ch7 = New(7)

	// MIDI channel 9
	Ch8 = New(8)

	// MIDI channel 10
	Ch9 = New(9)

	// MIDI channel 11
	Ch10 = New(10)

	// MIDI channel 12
	Ch11 = New(11)

	// MIDI channel 13
	Ch12 = New(12)

	// MIDI channel 14
	Ch13 = New(13)

	// MIDI channel 15
	Ch14 = New(14)

	// MIDI channel 16
	Ch15 = New(15)
)

// Channel represents a MIDI channel
type Channel struct {
	number uint8
}

// Number returns the number of the MIDI channel (0-15)
func (c Channel) Number() uint8 {
	return c.number
}

// NoteOff creates a note-off message on the channel for the given key
// The note-off message is "faked" by a note-on message of velocity 0.
// This allows saving bandwidth by using running status.
// If you need a "real" note-off message with velocity, use NoteOffVelocity.
func (c Channel) NoteOff(key uint8) NoteOff {
	return NoteOff{channel: c.Number(), key: key}
}

// NoteOffVelocity creates a note-off message with velocity on the channel.
func (c Channel) NoteOffVelocity(key uint8, velocity uint8) NoteOffVelocity {
	return NoteOffVelocity{NoteOff{channel: c.Number(), key: key}, velocity}
}

// NoteOn creates a note-on message on the channel
func (c Channel) NoteOn(key uint8, veloctiy uint8) NoteOn {
	return NoteOn{channel: c.Number(), key: key, velocity: veloctiy}
}

// KeyPressure creates a polyphonic aftertouch message on the channel
func (c Channel) KeyPressure(key uint8, pressure uint8) PolyphonicAfterTouch {
	return c.PolyphonicAfterTouch(key, pressure)
}

// PolyphonicAfterTouch creates a polyphonic aftertouch message on the channel
func (c Channel) PolyphonicAfterTouch(key uint8, pressure uint8) PolyphonicAfterTouch {
	return PolyphonicAfterTouch{channel: c.Number(), key: key, pressure: pressure}
}

// ControlChange creates a control change message on the channel
func (c Channel) ControlChange(controller uint8, value uint8) ControlChange {
	return ControlChange{channel: c.Number(), controller: controller, value: value}
}

// ProgramChange creates a program change message on the channel
func (c Channel) ProgramChange(program uint8) ProgramChange {
	return ProgramChange{channel: c.Number(), program: program}
}

// ChannelPressure creates an aftertouch message on the channel
func (c Channel) ChannelPressure(pressure uint8) AfterTouch {
	return c.AfterTouch(pressure)
}

// AfterTouch creates an aftertouch message on the channel
func (c Channel) AfterTouch(pressure uint8) AfterTouch {
	return AfterTouch{channel: c.Number(), pressure: pressure}
}

// PitchBend creates a pitch bend message on the channel
func (c Channel) PitchBend(value int16) PitchBend {
	return PitchBend{channel: c.Number(), value: value}
}

// Portamento creates a pitch bend message on the channel
func (c Channel) Portamento(value int16) PitchBend {
	return c.PitchBend(value)
}
