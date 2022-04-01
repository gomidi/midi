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
	byteTimingClock:/* RealTimeMsg.Set(TimingClockMsg), */ TimingClock,
	byteTick:/* RealTimeMsg.Set(TickMsg), */ Tick,
	byteStart:/* RealTimeMsg.Set(StartMsg), */ Start,
	byteContinue:/* RealTimeMsg.Set(ContinueMsg), */ Continue,
	byteStop:/* RealTimeMsg.Set(StopMsg), */ Stop,
	byteUndefined4:/* RealTimeMsg.Set(UndefinedMsg), */ UnknownType,
	byteActivesense:/* RealTimeMsg.Set(ActiveSenseMsg), */ ActiveSense,
	byteReset:/* RealTimeMsg.Set(ResetMsg), */ Reset,
}

// NewTimingClock returns a MIDI timing clock message
func NewTimingClock() []byte {
	//return NewMessage([]byte{byteTimingClock})
	return []byte{byteTimingClock}
}

// NewTick returns a midi tick message
func NewTick() []byte {
	//return NewMessage([]byte{byteTick})
	return []byte{byteTick}
}

// NewStart returns a MIDI start message
func NewStart() []byte {
	//return NewMessage([]byte{byteStart})
	return []byte{byteStart}
}

// NewContinue returns a MIDI continue message
func NewContinue() []byte {
	//return NewMessage([]byte{byteContinue})
	return []byte{byteContinue}
}

// NewStop returns a MIDI stop message
func NewStop() []byte {
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
func NewActivesense() []byte {
	//return NewMessage([]byte{byteActivesense})
	return []byte{byteActivesense}
}

// NewReset returns a MIDI reset message
func NewReset() []byte {
	//return NewMessage([]byte{byteReset})
	return []byte{byteReset}
}
