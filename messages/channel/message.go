package channel

type Message interface {
	String() string
	Raw() []byte
	Channel() uint8
	IsLiveMessage()
}

var (
	_ Message = NoteOff{}
	_ Message = NoteOffPedantic{}
	_ Message = NoteOn{}
	_ Message = PolyphonicAfterTouch{}
	_ Message = ControlChange{}
	_ Message = ProgramChange{}
	_ Message = AfterTouch{}
	_ Message = PitchWheel{}

	_ setter2 = NoteOff{}
	_ setter2 = NoteOffPedantic{}
	_ setter2 = NoteOn{}
	_ setter2 = PolyphonicAfterTouch{}
	_ setter2 = ControlChange{}
	_ setter2 = PitchWheel{}

	_ setter1 = ProgramChange{}
	_ setter1 = AfterTouch{}
)
