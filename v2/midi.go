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
}

type receiver struct {
	realtimeMsgCallback func(msg Message, deltamicrosec int64)
	otherMsgCallback    func(msg Message, deltamicrosec int64)
}

func NewReceiver(otherMsgCallback func(msg Message, deltamicrosec int64), realtimeMsgCallback func(msg Message, deltamicrosec int64)) Receiver {
	return &receiver{
		realtimeMsgCallback: realtimeMsgCallback,
		otherMsgCallback:    otherMsgCallback,
	}
}

func (r *receiver) Receive(msg []byte, deltamicrosec int64) {
	// TODO find out if there is a realtime message

	m := NewMessage(msg)

	if m.IsOneOf(RealTimeMsg, SysCommonMsg) && r.realtimeMsgCallback != nil {
		r.realtimeMsgCallback(m, deltamicrosec)
		return
	}

	if r.otherMsgCallback != nil {
		r.otherMsgCallback(m, deltamicrosec)
	}
}
