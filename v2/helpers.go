package midi

import (
	"fmt"
	"io"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// bin2decDenom converts the binary denominator to the decimal
func bin2decDenom(bin uint8) uint8 {
	if bin == 0 {
		return 1
	}
	return 2 << (bin - 1)
}

// channelMessage1 returns the bytes for a single byte channel message
func channelMessage1(c uint8, status, msg byte) []byte {
	cm := &channelMessage{channel: c, status: status}
	cm.data[0] = msg
	return cm.bytes()
}

// channelMessage2 returns the bytes for a two bytes channel message
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
func ReadChannelMessage(status byte, arg1 byte, rd io.Reader) (m Message, err error) {
	typ, channel := utils.ParseStatus(status)

	if err != nil {
		return
	}

	switch typ {

	// one argument only
	case byteProgramChange, byteChannelPressure:
		m = getMsg1(typ, channel, arg1)

	// two Arguments needed
	default:
		var arg2 byte
		arg2, err = utils.ReadByte(rd)

		if err != nil {
			return
		}
		m = getMsg2(typ, channel, arg1, arg2)
	}
	return
}

// getMsg1 returns a 1-byte channel message (program change or aftertouch)
func getMsg1(typ uint8, channel uint8, arg uint8) (m Message) {
	m.MsgType = ChannelMsg.Set(channelType[channel])
	m.Data = channelMessage1(channel, typ, arg)

	switch typ {
	case byteProgramChange:
		m.MsgType = m.MsgType.Set(ProgramChangeMsg)
	case byteChannelPressure:
		m.MsgType = m.MsgType.Set(AfterTouchMsg)
	default:
		panic(fmt.Sprintf("must not happen (typ % X is not an channel message with one argument)", typ))
	}

	return
}

// getMsg1 returns a 2-byte channel message (noteon/noteoff, poly aftertouch, control change or pitchbend)
func getMsg2(typ uint8, channel uint8, arg1 uint8, arg2 uint8) (msg Message) {
	msg.MsgType = ChannelMsg.Set(channelType[channel])
	msg.Data = channelMessage2(channel, typ, arg1, arg2)

	switch typ {
	case byteNoteOff:
		msg.MsgType = msg.MsgType.Set(NoteOffMsg)
	case byteNoteOn:
		msg.MsgType = msg.MsgType.Set(NoteOnMsg)
	case bytePolyphonicKeyPressure:
		msg.MsgType = msg.MsgType.Set(PolyAfterTouchMsg)
	case byteControlChange:
		msg.MsgType = msg.MsgType.Set(ControlChangeMsg)
	case bytePitchWheel:
		msg.MsgType = msg.MsgType.Set(PitchBendMsg)
	default:
		panic(fmt.Sprintf("must not happen (typ % X is not an channel message with two arguments)", typ))
	}
	return
}
