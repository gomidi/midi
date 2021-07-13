package midi

// Sender sends MIDI messages.
type Sender interface {
	// Send sends the given MIDI message and returns any error.
	Send(msg Message) error
}

// SenderTo sends MIDI messages.
type SenderTo interface {
	// SendTo sends MIDI messages to the given receiver.
	SendTo(Receiver) error
}

type ReceiverFunc func(msg Message, absmicrosec int64)

func (r ReceiverFunc) Receive(msg Message, absmicrosec int64) {
	r(msg, absmicrosec)
}

// Receiver receives MIDI messages.
type Receiver interface {
	// Receive receives a MIDI message. deltamicrosec is the delta to the previous note in microseconds (^-6)
	Receive(msg Message, absmicrosec int64)

	// println(big.NewRat(math.MaxInt64,1000 /* milliseonds */ *1000 /* seconds */ *60 /* minutes */ *60 /* hours */ *24 /* days */ *365 /* years */).FloatString(0))
	// output: 292471
	// => a ascending timestamp based on microseconds would wrap after 292471 years
	// so absolute timestamp should be preferred
}

type SysExReceiver interface {
	Receiver
	ReceiveSysEx(data []byte)
}

type RealtimeReceiver interface {
	Receiver
	ReceiveRealTime(typ MsgType, absmicrosec int64)
}

type SysCommonReceiver interface {
	Receiver
	ReceiveSysCommon(msg Message, absmicrosec int64)
}

// wrapreceiver implements the Receiver interface
type wrapreceiver struct {
	realtimeCallback  func(mtype MsgType, absmicrosec int64)
	channelCallback   func(msg Message, absmicrosec int64)
	syscommonCallback func(m Message, absmicrosec int64)
	sysExCallback     func(data []byte)
}

// NewWrapReceiver returns a Receiver that calls msgCallback for every non-realtime message and if realtimeMsgCallback is not nil, calls it
// for every realtime message.
func NewWrapReceiver(inner Receiver) *wrapreceiver {
	r := &wrapreceiver{
		channelCallback: inner.Receive,
	}

	if rt, is := inner.(RealtimeReceiver); is {
		r.realtimeCallback = rt.ReceiveRealTime
	}

	if sc, is := inner.(SysCommonReceiver); is {
		r.syscommonCallback = sc.ReceiveSysCommon
	}

	if sx, is := inner.(SysExReceiver); is {
		r.sysExCallback = sx.ReceiveSysEx
	}

	return r
}

func (r *wrapreceiver) Receive(m Message, absmicrosec int64) {
	switch {
	case m.Is(RealTimeMsg):
		if r.realtimeCallback != nil {
			r.realtimeCallback(m.MsgType, absmicrosec)
		}
	case m.Is(SysCommonMsg):
		if r.syscommonCallback != nil {
			r.syscommonCallback(m, absmicrosec)
		}
	case m.Is(SysExMsg):
		if r.sysExCallback != nil {
			r.sysExCallback(m.Data[1 : len(m.Data)-1])
		}
	default:
		r.channelCallback(m, absmicrosec)
	}
}
