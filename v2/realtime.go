package midi

const (
	byteTimingClock = 0xF8
	byteTick        = 0xF9
	byteStart       = 0xFA
	byteContinue    = 0xFB
	byteStop        = 0xFC
	byteUndefined4  = 0xFD
	byteActivesense = 0xFE
	byteReset       = 0xFF
)

var rtMessages = map[byte]MsgType{
	byteTimingClock:/* RealTimeMsg.Set(TimingClockMsg), */ TimingClock,
	byteTick:/* RealTimeMsg.Set(TickMsg), */ Tick,
	byteStart:/* RealTimeMsg.Set(StartMsg), */ Start,
	byteContinue:/* RealTimeMsg.Set(ContinueMsg), */ Continue,
	byteStop:/* RealTimeMsg.Set(StopMsg), */ Stop,
	byteUndefined4:/* RealTimeMsg.Set(UndefinedMsg), */ Undefined,
	byteActivesense:/* RealTimeMsg.Set(ActiveSenseMsg), */ ActiveSense,
	byteReset:/* RealTimeMsg.Set(ResetMsg), */ Reset,
}

// TimingClock returns a MIDI timing clock message
func NewTimingClock() Msg {
	return NewMsg([]byte{byteTimingClock})
}

// Tick returns a midi tick message
func NewTick() Msg {
	return NewMsg([]byte{byteTick})
}

// Start returns a MIDI start message
func NewStart() Msg {
	return NewMsg([]byte{byteStart})
}

// Continue returns a MIDI continue message
func NewContinue() Msg {
	return NewMsg([]byte{byteContinue})
}

// Stop returns a MIDI stop message
func NewStop() Msg {
	return NewMsg([]byte{byteStop})
}

// Undefined returns an undefined realtime message
func NewUndefined() Msg {
	return NewMsg([]byte{byteUndefined4})
}

// Activesense returns a MIDI active sensing message
func NewActivesense() Msg {
	return NewMsg([]byte{byteActivesense})
}

// Reset returns a MIDI reset message
func NewReset() Msg {
	return NewMsg([]byte{byteReset})
}
