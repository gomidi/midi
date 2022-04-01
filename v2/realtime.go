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

// NewTimingClock returns a MIDI timing clock message
func TimingClock() Message {
	//return NewMessage([]byte{byteTimingClock})
	return []byte{byteTimingClock}
}

// NewTick returns a midi tick message
func Tick() Message {
	//return NewMessage([]byte{byteTick})
	return []byte{byteTick}
}

// NewStart returns a MIDI start message
func Start() Message {
	//return NewMessage([]byte{byteStart})
	return []byte{byteStart}
}

// NewContinue returns a MIDI continue message
func Continue() Message {
	//return NewMessage([]byte{byteContinue})
	return []byte{byteContinue}
}

// NewStop returns a MIDI stop message
func Stop() Message {
	//return NewMessage([]byte{byteStop})
	return []byte{byteStop}
}

/*
// NewUndefined returns an undefined realtime message
func NewUndefined() Message {
	return NewMsg([]byte{byteUndefined4})
}
*/

// NewActivesense returns a MIDI active sensing message
func Activesense() Message {
	//return NewMessage([]byte{byteActivesense})
	return []byte{byteActivesense}
}

// NewReset returns a MIDI reset message
func Reset() Message {
	//return NewMessage([]byte{byteReset})
	return []byte{byteReset}
}
