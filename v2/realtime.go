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
	byteTimingClock:/* RealTimeMsg.Set(TimingClockMsg), */ TimingClockMsg,
	byteTick:/* RealTimeMsg.Set(TickMsg), */ TickMsg,
	byteStart:/* RealTimeMsg.Set(StartMsg), */ StartMsg,
	byteContinue:/* RealTimeMsg.Set(ContinueMsg), */ ContinueMsg,
	byteStop:/* RealTimeMsg.Set(StopMsg), */ StopMsg,
	byteUndefined4:/* RealTimeMsg.Set(UndefinedMsg), */ UndefinedMsgType,
	byteActivesense:/* RealTimeMsg.Set(ActiveSenseMsg), */ ActiveSenseMsg,
	byteReset:/* RealTimeMsg.Set(ResetMsg), */ ResetMsg,
}

// TimingClock returns a MIDI timing clock message
func TimingClock() Msg {
	return NewMsg([]byte{byteTimingClock})
}

// Tick returns a midi tick message
func Tick() Msg {
	return NewMsg([]byte{byteTick})
}

// Start returns a MIDI start message
func Start() Msg {
	return NewMsg([]byte{byteStart})
}

// Continue returns a MIDI continue message
func Continue() Msg {
	return NewMsg([]byte{byteContinue})
}

// Stop returns a MIDI stop message
func Stop() Msg {
	return NewMsg([]byte{byteStop})
}

// Undefined returns an undefined realtime message
func Undefined() Msg {
	return NewMsg([]byte{byteUndefined4})
}

// Activesense returns a MIDI active sensing message
func Activesense() Msg {
	return NewMsg([]byte{byteActivesense})
}

// Reset returns a MIDI reset message
func Reset() Msg {
	return NewMsg([]byte{byteReset})
}
