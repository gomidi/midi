package midi

import (
	"bytes"

	"gitlab.com/gomidi/midi/v2/drivers"
	midilib "gitlab.com/gomidi/midi/v2/internal/utils"
)

type readerState int

const (
	readerStateClean                readerState = 0
	readerStateWithinChannelMessage readerState = 1
	readerStateWithinSysCommon      readerState = 2
	readerStateInSysEx              readerState = 3
	readerStateWithinUnknown        readerState = 4
)

func ListenToPort(portnumber int, recv Receiver) error {
	rd := newReader(recv)
	in, err := drivers.InByNumber(portnumber)
	if err != nil {
		return err
	}

	return in.StartListening(func(data []byte, timestamp int64) {
		rd.Write(data, timestamp)
	})
}

func newReader(receiver Receiver) *reader {
	p := &reader{
		receiver: receiver,
		state:    readerStateClean,
	}

	if rt, is := receiver.(RealtimeReceiver); is {
		p.rtHandler = rt.ReceiveRealtime
	}

	if rs, is := receiver.(SysExReceiver); is {
		p.sysexHandler = rs.ReceiveSysEx
	}

	if rsc, is := receiver.(SysCommonReceiver); is {
		p.syscommonHander = rsc.ReceiveSysCommon
	}

	return p
}

type reader struct {
	receiver Receiver
	sysexBf  bytes.Buffer
	//bf              bytes.Buffer
	bf              [2]byte // first: is set, second: the byte
	state           readerState
	rtHandler       func(MsgType, int64)
	sysexHandler    func([]byte)
	syscommonHander func(Message, int64)
	statusByte      uint8
	channel         uint8
	typ             uint8
	timestamp       int64
}

func (p *reader) setBf(b byte) {
	p.bf[0] = 1
	p.bf[1] = b
}

func (p *reader) resetBf() {
	p.bf[0] = 0
}

func (p *reader) issetBf() bool {
	return p.bf[0] == 1
}

func (p *reader) getBf() (b byte) {
	return p.bf[1]
}

func (p *reader) Write(bt []byte, timestamp int64) {
	//fmt.Printf("% X\n", bt)
	for _, b := range bt {
		p.put(b, timestamp)
	}
}

func (p *reader) put(b byte, timestamp int64) {

	// => realtime message
	if b >= 0xF8 {
		if p.rtHandler != nil {
			p.rtHandler(NewMessage([]byte{b}).MsgType, timestamp)
		}
		return
	}

	p.timestamp = timestamp

	//fmt.Printf("state: %v\n", p.state)

	switch p.state {
	case readerStateInSysEx:
		p.inSysEx(b)
	case readerStateClean:
		p.cleanState(b)
	case readerStateWithinUnknown:
		p.withinUnknown(b)
	case readerStateWithinSysCommon:
		p.withinSysCommon(b)
	case readerStateWithinChannelMessage:
		p.withinChannelMessage(b)
	default:
		panic("unknown state")
	}
}

func (p *reader) inSysEx(b byte) {
	/*

	 */
	if b == 0xF7 {
		if p.sysexHandler != nil {
			p.sysexHandler(p.sysexBf.Bytes())
			p.sysexBf.Reset()
		}
		p.state = readerStateClean
		return
	}
	if midilib.IsStatusByte(b) {
		p.sysexBf.Reset()
		p.state = readerStateClean
		p.cleanState(b)
		return
	}

	if p.sysexHandler != nil {
		p.sysexBf.WriteByte(b)
	}

}

func (p *reader) cleanState(b byte) {
	switch {

	/* start sysex */
	case b == 0xF0:
		p.statusByte = 0
		p.sysexBf.Reset()
		p.state = readerStateInSysEx
	// end sysex
	case b == 0xF7:
		p.sysexBf.Reset()
		p.statusByte = 0

	// here we clear for System Common Category messages
	case b > 0xF0 && b < 0xF7:
		p.statusByte = 0
		p.resetBf()
		switch b {
		case byteMIDITimingCodeMessage, byteSysSongPositionPointer, byteSysSongSelect:
			p.state = readerStateWithinSysCommon
			p.typ = b
		case byteSysTuneRequest:
			if p.syscommonHander != nil {
				p.syscommonHander(Tune(), p.timestamp)
			}
			return
		default:
			// 0xF4, 0xF5, or 0xFD
			p.state = readerStateWithinUnknown
			return
		}
	case b >= 0x80 && b <= 0xEF:
		p.statusByte = b
		p.resetBf()
		p.typ, p.channel = midilib.ParseStatus(p.statusByte)
		p.state = readerStateWithinChannelMessage
	default:
		if p.statusByte != 0 {
			p.state = readerStateWithinChannelMessage
			p.withinChannelMessage(b)
		}
	}
}

func (p *reader) withinChannelMessage(b byte) {
	//p.bf.WriteByte(b)
	//l := p.bf.Len()

	switch p.typ {
	case byteChannelPressure:
		p.receiver.Receive(Channel(p.channel).Aftertouch(b), p.timestamp)
		p.resetBf()
		p.state = readerStateClean
	case byteProgramChange:
		p.receiver.Receive(Channel(p.channel).ProgramChange(b), p.timestamp)
		p.resetBf()
		p.state = readerStateClean
	case byteControlChange:
		if p.issetBf() {
			p.receiver.Receive(Channel(p.channel).ControlChange(p.getBf(), b), p.timestamp)
			p.resetBf()
			p.state = readerStateClean
		} else {
			p.setBf(b)
		}
	case byteNoteOn:
		if p.issetBf() {
			p.receiver.Receive(Channel(p.channel).NoteOn(p.getBf(), b), p.timestamp)
			p.resetBf()
			p.state = readerStateClean
		} else {
			p.setBf(b)
		}
	case byteNoteOff:
		if p.issetBf() {
			p.receiver.Receive(Channel(p.channel).NoteOffVelocity(p.getBf(), b), p.timestamp)
			p.resetBf()
			p.state = readerStateClean
		} else {
			p.setBf(b)
		}
	case bytePolyphonicKeyPressure:
		if p.issetBf() {
			p.receiver.Receive(Channel(p.channel).PolyAftertouch(p.getBf(), b), p.timestamp)
			p.resetBf()
			p.state = readerStateClean
		} else {
			p.setBf(b)
		}
	case bytePitchWheel:
		if p.issetBf() {
			rel, abs := midilib.ParsePitchWheelVals(p.getBf(), b)
			_ = abs
			p.receiver.Receive(Channel(p.channel).Pitchbend(rel), p.timestamp)
			p.resetBf()
			p.state = readerStateClean
		} else {
			p.setBf(b)
		}
	default:
		panic("unknown typ")
	}
}

func (p *reader) withinSysCommon(b byte) {

	switch p.typ {
	case byteMIDITimingCodeMessage:
		if p.syscommonHander != nil {
			p.syscommonHander(MTC(b), p.timestamp)
		}
		p.resetBf()
		p.state = readerStateClean
	case byteSysSongPositionPointer:
		if p.issetBf() {
			if p.syscommonHander != nil {
				_, abs := midilib.ParsePitchWheelVals(p.getBf(), b)
				p.syscommonHander(SPP(abs), p.timestamp)
			}
			p.resetBf()
			p.state = readerStateClean
		} else {
			p.setBf(b)
		}
	case byteSysSongSelect:
		if p.syscommonHander != nil {
			p.syscommonHander(SongSelect(b), p.timestamp)
		}
		p.resetBf()
		p.state = readerStateClean
	case byteSysTuneRequest:
		panic("must not be handled here, but within clean state")
	default:
		panic("unknown syscommon")
	}
}

func (p *reader) withinUnknown(b byte) {
	if midilib.IsStatusByte(b) {
		p.state = readerStateClean
		p.cleanState(b)
	}
	return
}

func (p *reader) HasSysExHandler() bool {
	return p.sysexHandler != nil
}

func (p *reader) HasSysCommonHandler() bool {
	return p.syscommonHander != nil
}

func (p *reader) HasRealTimeHandler() bool {
	return p.rtHandler != nil
}
