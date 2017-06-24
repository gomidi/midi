package realtime

const (
	// MIDI Clock
	TimingClock = msg(0xF8)

	// Tick
	Tick = msg(0xF9)

	// MIDI Start
	Start = msg(0xFA)

	// MIDI Continue
	Continue = msg(0xFB)

	// MIDI Stop
	Stop = msg(0xFC)

	// Undefined4
	Undefined4 = msg(0xFD)

	// Active Sense
	ActiveSensing = msg(0xFE)

	// Reset
	Reset = msg(0xFF)
)

// Message is a System Realtime Message
type Message interface {
	String() string
	Raw() []byte
	realTime()
}

type msg byte

func (m msg) String() string {
	return msg2String[m]
}

func (m msg) Raw() []byte {
	return []byte{byte(m)}
}

func (m msg) realTime() {}

var msg2String = map[msg]string{
	TimingClock:   "TimingClock",
	Tick:          "Tick",
	Start:         "Start",
	Continue:      "Continue",
	Stop:          "Stop",
	Undefined4:    "Undefined4",
	ActiveSensing: "ActiveSensing",
	Reset:         "Reset",
}
