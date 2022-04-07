package midi

import (
	"io"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// channelMessage1 returns the bytes for a single byte channel message
func channelMessage1(c uint8, status, msg byte) Message {
	cm := &channelMessage{channel: c, status: status}
	cm.data[0] = msg
	return cm.bytes()
}

// channelMessage2 returns the bytes for a two bytes channel message
func channelMessage2(c uint8, status, msg1 byte, msg2 byte) Message {
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
	s = utils.ClearBitU8(s, 0)
	s = utils.ClearBitU8(s, 1)
	s = utils.ClearBitU8(s, 2)
	s = utils.ClearBitU8(s, 3)
	s = s | m.channel
	return s
}

func (m *channelMessage) bytes() []byte {
	if m.twoBytes {
		return []byte{m.getCompleteStatus(), m.data[0], m.data[1]}
	}
	return []byte{m.getCompleteStatus(), m.data[0]}
}

const (
	byteProgramChange         = 0xC
	byteChannelPressure       = 0xD
	byteNoteOff               = 0x8
	byteNoteOn                = 0x9
	bytePolyphonicKeyPressure = 0xA
	byteControlChange         = 0xB
	bytePitchWheel            = 0xE
)

// ReadChannelMessage reads a channel message for the given status byte from the given reader.
// Don't use this function as a user, it is only internal to the library.
func ReadChannelMessage(status byte, arg1 byte, rd io.Reader) (m Message, err error) {
	typ, channel := utils.ParseStatus(status)

	if err != nil {
		return
	}

	switch typ {

	// one argument only
	case byteProgramChange, byteChannelPressure:
		m = channelMessage1(channel, typ, arg1)

	// two Arguments needed
	default:
		var arg2 byte
		arg2, err = utils.ReadByte(rd)

		if err != nil {
			return
		}
		m = channelMessage2(channel, typ, arg1, arg2)
	}
	return
}
