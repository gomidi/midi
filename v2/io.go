package midi

// Sender sends MIDI messages.
type Sender interface {
	// Send sends the given MIDI message and returns any error.
	Send(msg []byte) error
}

// SenderTo sends MIDI messages.
type SenderTo interface {
	// SendTo sends MIDI messages to the given receiver.
	SendTo(Receiver) error
}

// Receiver receives MIDI messages.
type Receiver interface {
	// Receive receives a MIDI message. deltamicrosec is the delta to the previous note in microseconds (^-6)
	Receive(msg []byte, deltamicrosec int64)

	// println(big.NewRat(math.MaxInt64,1000 /* milliseonds */ *1000 /* seconds */ *60 /* minutes */ *60 /* hours */ *24 /* days */ *365 /* years */).FloatString(0))
	// output: 292471
	// => a ascending timestamp based on microseconds would wrap after 292471 years
}

// receiver implements the Receiver interface
type receiver struct {
	realtimeMsgCallback func(msg Message, deltamicrosec int64)
	otherMsgCallback    func(msg Message, deltamicrosec int64)
}

// NewReceiver returns a Receiver that calls msgCallback for every non-realtime message and if realtimeMsgCallback is not nil, calls it
// for every realtime message.
func NewReceiver(msgCallback func(msg Message, deltamicrosec int64), realtimeMsgCallback func(msg Message, deltamicrosec int64)) Receiver {
	return &receiver{
		realtimeMsgCallback: realtimeMsgCallback,
		otherMsgCallback:    msgCallback,
	}
}

func (r *receiver) Receive(msg []byte, deltamicrosec int64) {
	m := NewMessage(msg)

	if m.IsOneOf(RealTimeMsg) && r.realtimeMsgCallback != nil {
		r.realtimeMsgCallback(m, deltamicrosec)
		return
	}

	if r.otherMsgCallback != nil {
		r.otherMsgCallback(m, deltamicrosec)
	}
}
