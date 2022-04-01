package smf

type Event struct {
	Delta   uint32
	Message Message
	//Data  []byte
}

/*
func (e *Event) Message() midi.Message {

	if len(e.Data) == 0 {
		var m midi.Message
		return m
	}

	if e.Data[0] == 0xFF {
		var msg MetaMessage
		msg.Data = e.Data
		msg.Type = GetMetaType(e.Data[1])
		return msg.Message
	}

	return midi.NewMsg(e.Data)
}
*/
