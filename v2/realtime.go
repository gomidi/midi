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

var rtMessages = map[byte]Type{
	byteTimingClock:/* RealTimeMsg.Set(TimingClockMsg), */ TimingClockMsg,
	byteTick:/* RealTimeMsg.Set(TickMsg), */ TickMsg,
	byteStart:/* RealTimeMsg.Set(StartMsg), */ StartMsg,
	byteContinue:/* RealTimeMsg.Set(ContinueMsg), */ ContinueMsg,
	byteStop:/* RealTimeMsg.Set(StopMsg), */ StopMsg,
	byteUndefined4:/* RealTimeMsg.Set(UndefinedMsg), */ UnknownMsg,
	byteActivesense:/* RealTimeMsg.Set(ActiveSenseMsg), */ ActiveSenseMsg,
	byteReset:/* RealTimeMsg.Set(ResetMsg), */ ResetMsg,
}

// TimingClock returns a timing clock message
func TimingClock() Message {
	return []byte{byteTimingClock}
}

// Tick returns a tick message
func Tick() Message {
	return []byte{byteTick}
}

// Start returns a start message
func Start() Message {
	return []byte{byteStart}
}

// Continue returns a continue message
func Continue() Message {
	return []byte{byteContinue}
}

// Stop returns a stop message
func Stop() Message {
	return []byte{byteStop}
}

/*
// NewUndefined returns an undefined realtime message
func NewUndefined() Message {
	return NewMsg([]byte{byteUndefined4})
}
*/

// Activesense returns an active sensing message
func Activesense() Message {
	return []byte{byteActivesense}
}

// Reset returns a reset message
func Reset() Message {
	return []byte{byteReset}
}
