package midi

type Filter []MsgType

func (f Filter) Contains(msg Message) bool {
	return msg.MsgType.IsOneOf(f...)
}

func (f Filter) Receiver(rec Receiver) Receiver {
	return ReceiverFunc(func(m Message, ts int32) {
		if f.Contains(m) {
			rec.Receive(m, ts)
		}
	})
}
