package midi

// NewListener returns a new Listener
func NewListener(portName string) *Listener {
	l := &Listener{}
	l.in, l.err = InByName(portName)
	return l
}

// Listener is an utility struct to make listening on a MIDI port for (filtered) messages easy.
type Listener struct {
	err              error
	in               In
	filter           Filter
	realtimeCallback func(msg Message, deltamicrosec int64)
}

func (l *Listener) Error() error {
	return l.err
}

func (l *Listener) Only(mtypes ...MsgType) *Listener {
	if len(mtypes) > 0 {
		l.filter = Filter(mtypes)
	}
	return l
}

func (l *Listener) RealTime(realtimeMsgCallback func(msg Message, deltamicrosec int64)) *Listener {
	l.realtimeCallback = realtimeMsgCallback
	return l
}

func (l *Listener) Do(fn func(msg Message, deltamicroSec int64)) (In, error) {
	if l.err != nil {
		return l.in, l.err
	}

	var rec Receiver

	if l.filter == nil {
		rec = NewReceiver(fn, l.realtimeCallback)
	} else {
		var fun = func(m Message, delta int64) {
			//m := NewMessage(msg)

			if m.MsgType.IsOneOf(l.filter...) {
				fn(m, delta)
			}
		}

		var funrt func(m Message, delta int64)
		if l.realtimeCallback != nil {
			funrt = func(m Message, delta int64) {
				//m := NewMessage(msg)

				if m.MsgType.IsOneOf(l.filter...) {
					l.realtimeCallback(m, delta)
				}
			}
		}

		rec = NewReceiver(fun, funrt)
	}

	l.in.SendTo(rec)
	return l.in, nil
}
