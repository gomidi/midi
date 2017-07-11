package channel

// Message represents a channel message
type Message interface {
	String() string
	Raw() []byte
	Channel() uint8
}

var (
	_ Message = NoteOff{}
	_ Message = NoteOffVelocity{}
	_ Message = NoteOn{}
	_ Message = PolyphonicAfterTouch{}
	_ Message = ControlChange{}
	_ Message = ProgramChange{}
	_ Message = AfterTouch{}
	_ Message = PitchBend{}

	_ setter2 = NoteOff{}
	_ setter2 = NoteOffVelocity{}
	_ setter2 = NoteOn{}
	_ setter2 = PolyphonicAfterTouch{}
	_ setter2 = ControlChange{}
	_ setter2 = PitchBend{}

	_ setter1 = ProgramChange{}
	_ setter1 = AfterTouch{}
)
