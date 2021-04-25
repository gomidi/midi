package midi

/*
// Filter is a filter that may send and receive MIDI messages that have one of the given message types.
type Filter []MsgType

// SenderFunc returns a function that can be used to send the filtered messages to the given Sender.
func (f Filter) SenderFunc(target Sender) func(Message) error {
	return func(m Message) error {
		if m.MsgType.IsOneOf(f...) {
			return target.Send(m.Data)
		}
		return nil
	}
}

// Receiver returns a Receiver that calls the given function for each message that is received.
func (f Filter) Receiver(fn func(msg Message, deltamicrosec int64)) Receiver {
	return &filteringReceiver{f, fn}
}

type filteringReceiver struct {
	Filter
	fn func(msg Message, deltamicrosec int64)
}

// Read reads the midi message and silently filters the unwanted out
func (f *filteringReceiver) Receive(msg []byte, deltamicrosec int64) {
	m := NewMessage(msg)
	if m.MsgType.IsOneOf(f.Filter...) {
		f.fn(m, deltamicrosec)
	}
	return
}
*/
