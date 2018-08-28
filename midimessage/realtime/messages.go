package realtime

const (
	// TimingClock is a MIDI timing clock message
	TimingClock = msg(0xF8)

	// Tick is a midi tick message
	Tick = msg(0xF9)

	// Start is a MIDI start message
	Start = msg(0xFA)

	// Continue is a MIDI continue message
	Continue = msg(0xFB)

	// Stop is a MIDI stop message
	Stop = msg(0xFC)

	// Undefined4 is an undefined realtime message
	Undefined4 = msg(0xFD)

	// Activesense is a MIDI active sensing message
	Activesense = msg(0xFE)

	// Reset is a MIDI reset message
	Reset = msg(0xFF)
)

// Message is a System Realtime Message
type Message interface {
	String() string
	Raw() []byte
	realTime()
	// IsLiveMessage()
}

type msg byte

// String represents the MIDI message as a string (for debugging)
func (m msg) String() string {
	return msg2String[m]
}

// Raw returns the raw bytes for the message
func (m msg) Raw() []byte {
	return []byte{byte(m)}
}

/*
func (m msg) IsLiveMessage() {

}
*/

func (m msg) realTime() {}

var msg2String = map[msg]string{
	TimingClock: "TimingClock",
	Tick:        "Tick",
	Start:       "Start",
	Continue:    "Continue",
	Stop:        "Stop",
	Undefined4:  "Undefined4",
	Activesense: "Activesense",
	Reset:       "Reset",
}
