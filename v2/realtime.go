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
	byteTimingClock: RealTimeMsg.Set(TimingClockMsg),
	byteTick:        RealTimeMsg.Set(TickMsg),
	byteStart:       RealTimeMsg.Set(StartMsg),
	byteContinue:    RealTimeMsg.Set(ContinueMsg),
	byteStop:        RealTimeMsg.Set(StopMsg),
	byteUndefined4:  RealTimeMsg.Set(UndefinedMsg),
	byteActivesense: RealTimeMsg.Set(ActiveSenseMsg),
	byteReset:       RealTimeMsg.Set(ResetMsg),
}

// TimingClock returns a MIDI timing clock message
func TimingClock() Message {
	return NewMessage([]byte{byteTimingClock})
}

// Tick returns a midi tick message
func Tick() Message {
	return NewMessage([]byte{byteTick})
}

// Start returns a MIDI start message
func Start() Message {
	return NewMessage([]byte{byteStart})
}

// Continue returns a MIDI continue message
func Continue() Message {
	return NewMessage([]byte{byteContinue})
}

// Stop returns a MIDI stop message
func Stop() Message {
	return NewMessage([]byte{byteStop})
}

// Undefined returns an undefined realtime message
func Undefined() Message {
	return NewMessage([]byte{byteUndefined4})
}

// Activesense returns a MIDI active sensing message
func Activesense() Message {
	return NewMessage([]byte{byteActivesense})
}

// Reset returns a MIDI reset message
func Reset() Message {
	return NewMessage([]byte{byteReset})
}
