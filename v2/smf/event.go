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
		return NewMetaMessage(e.Data[1:2][0], e.Data[2:])
	}

	/*
		b1 := e.Data[0]

		var b2, b3 byte
		switch len(e.Data) {
		case 0, 1:
		case 2:
			b2 = e.Data[1]
		default:
			b1 = e.Data[0]
			b2 = e.Data[1]
			b3 = e.Data[2]
		}
	*/
	m := midi.NewMsg(e.Data)

	switch m.MsgType {
	case midi.SysExMsgType:
		return midi.SysEx(e.Data[1 : len(e.Data)-1])
	default:
		return m
	}

}

/*
func (e *Event) MsgType() midi.MessageType {
	return e.Message().Type()

		var b1, b2 byte = 0, 0
		switch len(e.Data) {
		case 0:
		case 1:
			b1 = e.Data[0]
		default:
			b1 = e.Data[0]
			b2 = e.Data[1]
		}
		return midi.GetMsgType(b1, b2)
}
*/
