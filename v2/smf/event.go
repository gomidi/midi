package smf

import (
	"gitlab.com/gomidi/midi/v2"
)

type Event struct {
	Delta uint32
	Data  []byte
}

func (e *Event) Message() midi.Message {

	if len(e.Data) == 0 {
		var m midi.Msg
		return m
	}

	if e.Data[0] == 0xFF {
		var msg MetaMessage
		msg.Data = e.Data
		msg.MetaType = GetMetaType(e.Data[1])
		return msg
	}

	m := midi.NewMsg(e.Data)

	switch m.MsgType {
	case midi.SysEx:
		return midi.NewSysEx(e.Data[1 : len(e.Data)-1])
	default:
		return m
	}

}
