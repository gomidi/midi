package midi

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2/drivers"
	midilib "gitlab.com/gomidi/midi/v2/internal/utils"
)

func _channelMessage(typ, channel, data1, data2 byte) Msg {
	ch := Channel(channel)
	switch typ {
	case byteChannelPressure:
		return ch.Aftertouch(data1)
	case byteProgramChange:
		return ch.ProgramChange(data1)
	case byteControlChange:
		return ch.ControlChange(data1, data2)
	case byteNoteOn:
		return ch.NoteOn(data1, data2)
	case byteNoteOff:
		return ch.NoteOffVelocity(data1, data2)
	case bytePolyphonicKeyPressure:
		return ch.PolyAftertouch(data1, data2)
	case bytePitchWheel:
		rel, _ := midilib.ParsePitchWheelVals(data1, data2)
		return ch.Pitchbend(rel)
	default:
		panic("unknown typ")
	}
}

// ListenOptions are the options for the listening
type ListenOptions struct {

	// TimeCode lets the the timecode messages pass through, if set
	TimeCode bool

	// ActiveSense lets the the active sense messages pass through, if set
	ActiveSense bool

	// SysExBufferSize defines the size of the buffer for sysex messages (in bytes).
	// SysEx messages larger than this size will be ignored.
	// When SysExBufferSize is 0, the default buffersize (1024) is used.
	SysExBufferSize uint32
}

var ErrPortClosed = drivers.ErrPortClosed
var ErrListenStopped = drivers.ErrListenStopped

func ListenToPort(portnumber int, recv Receiver, opt ListenOptions) (stop func(), err error) {
	in, err := drivers.InByNumber(portnumber)
	if err != nil {
		return nil, err
	}

	var conf drivers.ListenConfig
	conf.SysExBufferSize = opt.SysExBufferSize
	conf.TimeCode = opt.TimeCode
	conf.ActiveSense = opt.ActiveSense

	if sysrc, has := recv.(SysExReceiver); has {
		conf.OnSysEx = sysrc.OnSysEx
	}

	if errrc, has := recv.(ErrorReceiver); has {
		conf.OnErr = errrc.OnError
	}

	var isStatusSet bool
	var typ, channel byte

	var onMsg = func(data [3]byte, millisec int32) {
		status := data[0]
		var msg Msg
		switch {
		// realtime message
		case status >= 0xF8:
			msg = NewMsg(status, 0, 0)
		// here we clear for System Common Category messages
		case status > 0xF0 && status < 0xF7:
			isStatusSet = false
			switch status {
			case byteSysTuneRequest:
				msg = Tune()
			case byteMIDITimingCodeMessage:
				msg = MTC(data[1])
			case byteSysSongPositionPointer:
				_, abs := midilib.ParsePitchWheelVals(data[1], data[2])
				msg = SPP(abs)
			case byteSysSongSelect:
				msg = SongSelect(data[1])
			default:
				// undefined syscommon message
				msg = Undefined()
			}

		// [MIDI] permits 0xF7 octets that are not part of a (0xF0, 0xF7) pair
		// to appear on a MIDI 1.0 DIN cable.  Unpaired 0xF7 octets have no
		// semantic meaning in MIDI apart from cancelling running status.
		case status == 0xF7:
			isStatusSet = false
		case status == 0xF0:
			errmsg := fmt.Sprintf("error in driver: %q receiving 0xF0 in non sysex callback", drivers.Get().String())
			if conf.OnErr != nil {
				conf.OnErr(fmt.Errorf(errmsg))
			} else {
				// TODO: maybe log
				panic(errmsg)
			}
		case status >= 0x80 && status <= 0xEF:
			isStatusSet = true
			typ, channel = midilib.ParseStatus(status)
			msg = _channelMessage(typ, channel, data[1], data[2])
		default:
			// running status
			if isStatusSet {
				msg = _channelMessage(typ, channel, data[1], data[2])
			}
		}

		recv.Receive(msg, millisec)
	}

	return in.Listen(onMsg, conf)
}
