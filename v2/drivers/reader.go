package drivers

import (
	"fmt"

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

const (
	byteMIDITimingCodeMessage  = byte(0xF1)
	byteSysSongPositionPointer = byte(0xF2)
	byteSysSongSelect          = byte(0xF3)
	byteSysTuneRequest         = byte(0xF6)
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

type Reader struct {
	//	maxlenSysex     int
	sysexBf  []byte
	sysexlen int

	ts_ms      int32
	sysexTS    int32
	state      readerState
	statusByte uint8
	issetBf    bool
	bf         byte
	typ        uint8

	SysExBufferSize uint32
	OnMsg           func([]byte, int32)
	HandleSysex     bool
	OnErr           func(error)
}

func (r *Reader) withinChannelMessage(b byte) {
	//fmt.Println("withinChannelMessage")
	switch r.typ {
	case byteChannelPressure:
		r.issetBf = false
		r.state = readerStateClean
		//p.receiver.Receive(Channel(p.channel).Aftertouch(b), p.timestamp)
		r.OnMsg([]byte{r.statusByte, b, 0}, r.ts_ms)
	case byteProgramChange:
		r.issetBf = false // first: is set, second: the byte
		r.state = readerStateClean
		//p.receiver.Receive(Channel(p.channel).ProgramChange(b), p.timestamp)
		r.OnMsg([]byte{r.statusByte, b, 0}, r.ts_ms)
	case byteControlChange:
		if r.issetBf {
			r.issetBf = false // first: is set, second: the byte
			r.state = readerStateClean
			//p.receiver.Receive(Channel(p.channel).ControlChange(p.getBf(), b), p.timestamp)
			r.OnMsg([]byte{r.statusByte, r.bf, b}, r.ts_ms)
		} else {
			r.issetBf = true
			r.bf = b
		}
	case byteNoteOn:
		if r.issetBf {
			r.issetBf = false // first: is set, second: the byte
			r.state = readerStateClean
			//p.receiver.Receive(Channel(p.channel).NoteOn(p.getBf(), b), p.timestamp)
			r.OnMsg([]byte{r.statusByte, r.bf, b}, r.ts_ms)
		} else {
			r.issetBf = true
			r.bf = b
		}
	case byteNoteOff:
		if r.issetBf {
			r.issetBf = false // first: is set, second: the byte
			r.state = readerStateClean
			//p.receiver.Receive(Channel(p.channel).NoteOffVelocity(p.getBf(), b), p.timestamp)
			r.OnMsg([]byte{r.statusByte, r.bf, b}, r.ts_ms)
		} else {
			r.issetBf = true
			r.bf = b
		}
	case bytePolyphonicKeyPressure:
		if r.issetBf {
			r.issetBf = false // first: is set, second: the byte
			r.state = readerStateClean
			//p.receiver.Receive(Channel(p.channel).PolyAftertouch(p.getBf(), b), p.timestamp)
			r.OnMsg([]byte{r.statusByte, r.bf, b}, r.ts_ms)
		} else {
			r.issetBf = true
			r.bf = b
		}
	case bytePitchWheel:
		if r.issetBf {
			//rel, abs := midilib.ParsePitchWheelVals(bf, b)
			//_ = abs
			r.issetBf = false // first: is set, second: the byte
			r.state = readerStateClean
			//p.receiver.Receive(Channel(p.channel).Pitchbend(rel), p.timestamp)
			r.OnMsg([]byte{r.statusByte, r.bf, b}, r.ts_ms)
		} else {
			r.issetBf = true
			r.bf = b
		}
	default:
		panic("unknown typ")
	}
}

func (r *Reader) cleanState(b byte) {
	//fmt.Println("clean state")
	switch {

	/* start sysex */
	case b == 0xF0:
		r.statusByte = 0
		r.sysexBf = make([]byte, r.SysExBufferSize)
		//sysexBf.Reset()
		//sysexBf.WriteByte(b)
		r.sysexBf[0] = b
		r.sysexlen = 1
		r.sysexTS = r.ts_ms
		r.state = readerStateInSysEx
	// end sysex
	// [MIDI] permits 0xF7 octets that are not part of a (0xF0, 0xF7) pair
	// to appear on a MIDI 1.0 DIN cable.  Unpaired 0xF7 octets have no
	// semantic meaning in MIDI apart from cancelling running status.
	case b == 0xF7:
		r.sysexBf = nil
		r.sysexlen = 0
		r.statusByte = 0
		r.OnMsg([]byte{b, 0, 0}, r.ts_ms)

	// here we clear for System Common Category messages
	case b > 0xF0 && b < 0xF7:
		r.statusByte = 0
		r.issetBf = false // reset buffer
		//fmt.Printf("sys common msg started\n")
		switch b {
		case byteMIDITimingCodeMessage, byteSysSongPositionPointer, byteSysSongSelect:
			r.state = readerStateWithinSysCommon
			r.typ = b
		case byteSysTuneRequest:
			r.OnMsg([]byte{b, 0, 0}, r.ts_ms)
			/*
				if p.syscommonHander != nil {
					p.syscommonHander(Tune(), p.timestamp)
				}
			*/
			return
		default:
			// 0xF4, 0xF5, or 0xFD
			r.state = readerStateWithinUnknown
			return
		}

	// channel message with status byte
	case b >= 0x80 && b <= 0xEF:
		//fmt.Println("channel message")
		r.statusByte = b
		r.issetBf = false // reset buffer
		//typ, channel = midilib.ParseStatus(statusByte)
		r.typ, _ = midilib.ParseStatus(r.statusByte)
		r.state = readerStateWithinChannelMessage
	default:
		if r.statusByte != 0 {
			r.state = readerStateWithinChannelMessage
			r.withinChannelMessage(b)
		}
	}
}

func (r *Reader) eachByte(b byte) {
	if b >= 0xF8 {
		//r.OnMsg([]byte{b, 0, 0}, r.ts_ms)
		r.OnMsg([]byte{b}, r.ts_ms)
		return
	}

	//fmt.Printf("state: %v\n", p.state)

	switch r.state {
	case readerStateInSysEx:
		//fmt.Println("readerStateInSysEx")
		/* interrupted sysex, discard old data */
		if b == 0xF0 {
			r.statusByte = 0
			r.sysexBf = make([]byte, r.SysExBufferSize)
			r.sysexTS = r.ts_ms
			//sysexBf.Reset()
			//sysexBf.WriteByte(b)
			r.sysexBf[0] = b
			r.sysexlen = 1
			r.state = readerStateInSysEx
			return
		}

		if b == 0xF7 {
			/*
				if p.sysexHandler != nil {
					p.sysexBf.WriteByte(b)
					bt := p.sysexBf.Bytes()
					p.sysexBf.Reset()
					p.sysexHandler(bt, p.timestamp)
				}
			*/
			r.state = readerStateClean
			if r.HandleSysex {
				r.sysexBf[r.sysexlen] = b
				r.sysexlen++
				//go
				func(bb []byte, l int) {
					var _bt = make([]byte, l)

					for i := 0; i < l; i++ {
						_bt[i] = bb[i]
					}
					r.OnMsg(_bt, r.sysexTS)
				}(r.sysexBf, r.sysexlen)
			}
			r.sysexBf = nil
			r.sysexlen = 0
			return
		}
		if midilib.IsStatusByte(b) {
			//p.sysexBf.Reset()
			r.sysexBf = nil
			r.sysexlen = 0
			r.state = readerStateClean
			r.cleanState(b)
			return
		}

		if r.HandleSysex {
			r.sysexBf[r.sysexlen] = b
			r.sysexlen++
		}

		/*
			if p.sysexHandler != nil {
				p.sysexBf.WriteByte(b)
			}
		*/
	case readerStateClean:
		//fmt.Println("readerStateClean")
		r.cleanState(b)
	case readerStateWithinUnknown:
		//fmt.Println("readerStateWithinUnknown")
		//p.withinUnknown(b)
		if midilib.IsStatusByte(b) {
			r.state = readerStateClean
			r.cleanState(b)
		}
	case readerStateWithinSysCommon:
		//fmt.Println("readerStateWithinSysCommon")
		switch r.typ {
		case byteMIDITimingCodeMessage:
			/*
				if p.syscommonHander != nil {
					p.syscommonHander(MTC(b), p.timestamp)
				}
			*/
			r.issetBf = false
			r.state = readerStateClean
			r.OnMsg([]byte{r.typ, b, 0}, r.ts_ms)
		case byteSysSongPositionPointer:
			if r.issetBf {
				/*
					if p.syscommonHander != nil {
						_, abs := midilib.ParsePitchWheelVals(p.getBf(), b)
						p.syscommonHander(SPP(abs), p.timestamp)
					}
				*/
				r.issetBf = false
				r.state = readerStateClean
				r.OnMsg([]byte{r.typ, r.bf, b}, r.ts_ms)
			} else {
				r.issetBf = true
				r.bf = b
			}
		case byteSysSongSelect:
			/*
				if p.syscommonHander != nil {
					p.syscommonHander(SongSelect(b), p.timestamp)
				}
			*/
			r.issetBf = false
			r.state = readerStateClean
			r.OnMsg([]byte{r.typ, b, 0}, r.ts_ms)
		case byteSysTuneRequest:
			//panic("must not be handled here, but within clean state")
		default:
			if r.OnErr != nil {
				r.OnErr(fmt.Errorf("unknown syscommon message: % X", b))
			}
			//panic("unknown syscommon")
		}
	case readerStateWithinChannelMessage:
		//fmt.Println("readerStateWithinChannelMessage")
		r.withinChannelMessage(b)
	default:
		panic(fmt.Sprintf("unknown state %v, must not happen", r.state))
	}
}

func (r *Reader) Reset() {
	//fmt.Println("reset")

	if r.SysExBufferSize == 0 {
		r.SysExBufferSize = 1024
	}

	r.sysexBf = make([]byte, r.SysExBufferSize)
	r.sysexlen = 0
	r.ts_ms = 0
	r.statusByte = 0
	r.issetBf = false
	r.state = readerStateClean
}

func NewReader(config ListenConfig, onMsg func([]byte, int32)) *Reader {
	var r Reader
	r.OnMsg = onMsg
	r.OnErr = config.OnErr
	//r.OnSysEx = config.OnSysEx
	r.HandleSysex = config.SysEx
	r.SysExBufferSize = config.SysExBufferSize
	r.Reset()
	return &r
}

func (r *Reader) setDelta(deltaMilliSeconds int32) {
	r.ts_ms += deltaMilliSeconds
}

func (r *Reader) resetStatus() {
	r.statusByte = 0
	r.issetBf = false // first: is set, second: the byte
}

// func (r *Reader) EachMessage(bt []byte, deltaSeconds float64) {
func (r *Reader) EachMessage(bt []byte, deltaMilliSeconds int32) {

	// TODO: verify
	// assume that each call is without running state
	//r.ResetStatus()

	r.setDelta(deltaMilliSeconds) // int32(math.Round(deltaSeconds * 1000))

	//fmt.Printf("got % X\n", bt)

	for _, b := range bt {
		// => realtime message
		r.eachByte(b)

	}

}
