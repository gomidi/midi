package midi

import (
	"fmt"
	"io"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

const (
	byteProgramChange         = 0xC
	byteChannelPressure       = 0xD
	byteNoteOff               = 0x8
	byteNoteOn                = 0x9
	bytePolyphonicKeyPressure = 0xA
	byteControlChange         = 0xB
	bytePitchWheel            = 0xE
)

// Read reads a channel message
func ReadChannelMessage(status byte, arg1 byte, rd io.Reader) (m Message, err error) {
	typ, channel := utils.ParseStatus(status)

	// fmt.Printf("typ: %v channel: %v\n", typ, channel)

	// fmt.Printf("arg1: %v, err: %v\n", arg1, err)

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

func getMsg1(typ uint8, channel uint8, arg uint8) (m Message) {
	m.MsgType = ChannelMsg.Set(channelType[channel])
	//m.Data = []byte{channel, arg}
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

func getMsg2(typ uint8, channel uint8, arg1 uint8, arg2 uint8) (msg Message) {
	msg.MsgType = ChannelMsg.Set(channelType[channel])
	//msg.Data = []byte{channel, arg1, arg2}
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

	/*
		// handle noteOn messages with velocity of 0 as note offs
		if noteOn, is := msg.(NoteOn); is && noteOn.velocity == 0 {
			msg = (NoteOff{}).set(channel, arg1, 0)
		}
	*/
	return
}

type Channel uint8

func (ch Channel) NoteOn(key, velocity uint8) []byte {
	return channelMessage2(uint8(ch), 9, key, velocity)
}

func (ch Channel) NoteOffVelocity(key, velocity uint8) []byte {
	return channelMessage2(uint8(ch), 8, key, velocity)
}

func (ch Channel) NoteOff(key uint8) []byte {
	return channelMessage2(uint8(ch), 9, key, 0)
}

func (ch Channel) ProgramChange(program uint8) []byte {
	return channelMessage1(uint8(ch), 12, program)
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
