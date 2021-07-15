package smf

import (
	"gitlab.com/gomidi/midi/v2"
)

type Event struct {
	Delta uint32
	Data  []byte
}

func (e Event) Message() midi.Message {
	return midi.NewMessage(e.Data)
}

func (e *Event) MsgType() midi.MsgType {
	return midi.GetMsgType(e.Data)
}
