package realtime

const (
	// MIDI Clock
	TimingClock = _event(0xF8)
	// Tick
	Tick = _event(0xF9)
	// MIDI Start
	Start = _event(0xFA)
	// MIDI Continue
	Continue = _event(0xFB)
	// MIDI Stop
	Stop = _event(0xFC)
	// Undefined4
	Undefined4 = _event(0xFD)
	// Active Sense
	ActiveSensing = _event(0xFE)
	// Reset
	Reset = _event(0xFF)
)

type Event interface {
	String() string
	Raw() []byte
	realTime()
}

type _event byte

func (ev _event) String() string {
	return event2String[ev]
}

func (ev _event) Raw() []byte {
	return []byte{byte(ev)}
}

func (ev _event) realTime() {}

var event2String = map[_event]string{
	TimingClock:   "TimingClock",
	Tick:          "Tick",
	Start:         "Start",
	Continue:      "Continue",
	Stop:          "Stop",
	Undefined4:    "Undefined4",
	ActiveSensing: "ActiveSensing",
	Reset:         "Reset",
}

func Dispatch(b byte) Event {
	ev := _event(b)
	if _, has := event2String[ev]; !has {
		return nil
	}
	return ev
}
