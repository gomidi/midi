package channel

// New returns a new channel. no must be <= 15, since there are 16 midi channels.
func New(no uint8) Channel {
	if no > 15 {
		panic("midi channel > 15 is not allowed")
	}

	return Channel{no}
}

var (
	Ch0  = New(0)
	Ch1  = New(1)
	Ch2  = New(2)
	Ch3  = New(3)
	Ch4  = New(4)
	Ch5  = New(5)
	Ch6  = New(6)
	Ch7  = New(7)
	Ch8  = New(8)
	Ch9  = New(9)
	Ch10 = New(10)
	Ch11 = New(11)
	Ch12 = New(12)
	Ch13 = New(13)
	Ch14 = New(14)
	Ch15 = New(15)
)

type Channel struct {
	number uint8
}

func (c Channel) Number() uint8 {
	return c.number
}

func (c Channel) NoteOff(pitch uint8) NoteOff {
	return NoteOff{channel: c.Number(), pitch: pitch}
}

func (c Channel) NoteOn(pitch uint8, veloctiy uint8) NoteOn {
	return NoteOn{channel: c.Number(), pitch: pitch, velocity: veloctiy}
}

func (c Channel) PolyphonicAfterTouch(pitch uint8, pressure uint8) PolyphonicAfterTouch {
	return PolyphonicAfterTouch{channel: c.Number(), pitch: pitch, pressure: pressure}
}

func (c Channel) ControlChange(controller uint8, value uint8) ControlChange {
	return ControlChange{channel: c.Number(), controller: controller, value: value}
}

func (c Channel) ProgramChange(program uint8) ProgramChange {
	return ProgramChange{channel: c.Number(), program: program}
}

func (c Channel) AfterTouch(value uint8) AfterTouch {
	return AfterTouch{channel: c.Number(), pressure: value}
}

func (c Channel) PitchWheel(value int16) PitchWheel {
	return PitchWheel{channel: c.Number(), value: value}
}
