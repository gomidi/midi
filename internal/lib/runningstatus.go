package lib

type RunningStatus struct {
	status byte
}

/*
   his (http://midi.teragonaudio.com/tech/midispec.htm) take on running status buffer
   A recommended approach for a receiving device is to maintain its "running status buffer" as so:

       Buffer is cleared (ie, set to 0) at power up.
       Buffer stores the status when a Voice Category Status (ie, 0x80 to 0xEF) is received.
       Buffer is cleared when a System Common Category Status (ie, 0xF0 to 0xF7) is received.
       Nothing is done to the buffer when a RealTime Category message is received.
       Any data bytes are ignored when the buffer is 0. (I think that only holds for realtime midi)
*/

func (r *RunningStatus) HandleLive(canary byte) (status byte, changed bool) {
	if canary >= 0x80 && canary <= 0xEF {
		r.status = canary
		return r.status, true
	}

	if canary >= 0xF0 && canary <= 0xF7 {
		r.status = 0
		return 0, true
	}

	return 0, false
}

func (r *RunningStatus) HandleSMF(canary byte) (status byte, changed bool) {
	if canary >= 0x80 && canary <= 0xEF {
		r.status = canary
		return r.status, true
	}

	if canary >= 0xF0 && canary <= 0xF7 {
		r.status = 0
		return 0, true
	}

	// here we also clear for meta events
	if canary == 0xFF {
		r.status = 0
		return 0, true
	}

	return 0, false
}
