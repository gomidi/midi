package midi

type Filter []MsgType

func (f Filter) SenderFunc(target Sender) func(Message) error {
	return func(m Message) error {
		if m.MsgType.IsOneOf(f...) {
			return target.Send(m.Data)
		}
		return nil
	}
}

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
