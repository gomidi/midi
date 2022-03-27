package midi

/*
type Filter []MsgType

func (f Filter) Contains(msg Msg) bool {
	for _, t := range f {
		if Is(t, msg.MsgType) {
			return true
		}
	}
	return false
}

func (f Filter) Receiver(rec Receiver) Receiver {
	return ReceiverFunc(func(m Msg, ts int32) {
		if f.Contains(m) {
			rec.Receive(m, ts)
		}
	})
}
*/
