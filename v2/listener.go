package midi

import "fmt"

// NewListener returns a new Listener that listens on the given MIDI in port by calling the given
// msgCallback, when the StartListening method is called.
// Before that, a message filter can be set via the Only method and a callback for realtime
// messages can be set via the RealTime method.
func NewListener(in In, msgCallback func(msg Message, deltamicroSec int64)) (l *Listener, err error) {
	if msgCallback == nil {
		return nil, fmt.Errorf("msgCallback must not be nil")
	}
	l = &Listener{}
	l.msgCallback = msgCallback
	l.In = in
	err = l.In.Open()
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Listener is an utility struct to make listening on a MIDI port for (filtered) messages easy.
type Listener struct {
	In               In
	filter           []MsgType
	msgCallback      func(msg Message, deltamicroSec int64)
	realtimeCallback func(msg Message, deltamicrosec int64)
}

// Only sets the message types that should be listened on.
// I.e. any of the given MsgTypes will be passed to the given callback(s).
func (l *Listener) Only(mtypes ...MsgType) *Listener {
	if len(mtypes) > 0 {
		l.filter = mtypes
	}
	return l
}

// RealTime sets a callback for realtime messages.
func (l *Listener) RealTime(realtimeMsgCallback func(msg Message, deltamicrosec int64)) *Listener {
	l.realtimeCallback = realtimeMsgCallback
	return l
}

// StartListening starts the listening.
func (l *Listener) StartListening() {
	var rec Receiver

	if l.filter == nil {
		rec = NewReceiver(l.msgCallback, l.realtimeCallback)
	} else {
		var fun = func(m Message, delta int64) {
			//m := NewMessage(msg)

			if m.MsgType.IsOneOf(l.filter...) {
				l.msgCallback(m, delta)
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

	l.In.SendTo(rec)
}

// StopListening stops the listening
func (l *Listener) StopListening() {
	l.In.StopListening()
}

// Close closes the underlying MIDI In port.
func (l *Listener) Close() {
	l.In.StopListening()
}
