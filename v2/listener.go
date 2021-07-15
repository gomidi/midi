package midi

/*
TODO

merge with NewWrapper, i.e. add In and filter to wrapperreceiver and rename it and delete this
*/

/*
// NewListener returns a new Listener that listens on the given MIDI in port by calling the given
// msgCallback, when the StartListening method is called.
// Before that, a message filter can be set via the Only method and a callback for realtime
// messages can be set via the RealTime method.
func NewListener(in In, rec Receiver) (l *Listener, err error) {
	if rec == nil {
		return nil, fmt.Errorf("rec must not be nil")
	}
	l = &Listener{}
	//l.msgCallback = msgCallback
	l.rec = NewWrapReceiver(rec)
	l.In = in
	err = l.In.Open()
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Listener is an utility struct to make listening on a MIDI port for (filtered) messages easy.
type Listener struct {
	In     In
	rec    *wrapreceiver
	filter []MsgType
	//msgCallback      func(msg Message, deltamicroSec int64)
	//realtimeCallback func(msg Message, deltamicrosec int64)
}

// Only sets the message types that should be listened on.
// I.e. any of the given MsgTypes will be passed to the given callback(s).
func (l *Listener) Only(mtypes ...MsgType) *Listener {
	if len(mtypes) > 0 {
		l.filter = mtypes
	}
	return l
}

// StartListening starts the listening.
func (l *Listener) StartListening() {
	if l.filter == nil {
		//rec = NewReceiver(l.msgCallback, l.realtimeCallback)
	} else {
		var inner = l.rec.channelCallback
		var fun = func(m Message, abs int64) {

			if m.MsgType.IsOneOf(l.filter...) {
				inner(m, abs)
			}
		}

		l.rec.channelCallback = fun

		var innerRT = l.rec.realtimeCallback
		var funrt func(mtype MsgType, abs int64)
		funrt = func(mtype MsgType, abs int64) {

			if mtype.IsOneOf(l.filter...) {
				innerRT(mtype, abs)
			}
		}

		l.rec.realtimeCallback = funrt
	}

	l.In.SendTo(l.rec)
}

// StopListening stops the listening
func (l *Listener) StopListening() {
	l.In.StopListening()
}

// Close closes the underlying MIDI In port.
func (l *Listener) Close() {
	l.In.StopListening()
}
*/
