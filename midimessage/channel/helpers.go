package channel

import (
	"github.com/gomidi/midi/internal/midilib"
)

type setter1 interface {
	Message
	set(channel uint8, firstArg uint8) setter1
}

type setter2 interface {
	Message
	set(channel uint8, firstArg, secondArg uint8) setter2
}

func channelMessage1(c uint8, status, msg byte) []byte {
	cm := &channelMessage{channel: c, status: status}
	cm.data[0] = msg
	return cm.bytes()
}

func channelMessage2(c uint8, status, msg1 byte, msg2 byte) []byte {
	cm := &channelMessage{channel: c, status: status}
	cm.data[0] = msg1
	cm.data[1] = msg2
	cm.twoBytes = true
	return cm.bytes()
}

type channelMessage struct {
	status   uint8
	channel  uint8
	twoBytes bool
	data     [2]byte
}

func (m *channelMessage) getCompleteStatus() uint8 {
	s := m.status << 4
	s = midilib.ClearBitU8(s, 0)
	s = midilib.ClearBitU8(s, 1)
	s = midilib.ClearBitU8(s, 2)
	s = midilib.ClearBitU8(s, 3)
	s = s | m.channel
	return s
}

func (m *channelMessage) bytes() []byte {
	if m.twoBytes {
		return []byte{m.getCompleteStatus(), m.data[0], m.data[1]}
	}
	return []byte{m.getCompleteStatus(), m.data[0]}
}
