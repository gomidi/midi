package midi

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2/drivers"
	midilib "gitlab.com/gomidi/midi/v2/internal/utils"
)

func _channelMessage(typ, channel, data1, data2 byte) Message {
	switch typ {
	case byteChannelPressure:
		return AfterTouch(channel, data1)
	case byteProgramChange:
		return ProgramChange(channel, data1)
	case byteControlChange:
		return ControlChange(channel, data1, data2)
	case byteNoteOn:
		return NoteOn(channel, data1, data2)
	case byteNoteOff:
		return NoteOffVelocity(channel, data1, data2)
	case bytePolyphonicKeyPressure:
		return PolyAfterTouch(channel, data1, data2)
	case bytePitchWheel:
		rel, _ := midilib.ParsePitchWheelVals(data1, data2)
		return Pitchbend(channel, rel)
	default:
		panic("unknown typ")
	}
}

// ListenOptions are the options for the listening
type ListenOptions struct {

	// TimeCode lets the timecode messages pass through, if set
	TimeCode bool

	// ActiveSense lets the active sense messages pass through, if set
	ActiveSense bool

	// SysEx lets the system exclusive messages pass through, if set
	SysEx bool

	// SysExBufferSize defines the size of the buffer for sysex messages (in bytes).
	// SysEx messages larger than this size will be ignored.
	// When SysExBufferSize is 0, the default buffersize (1024) is used.
	SysExBufferSize uint32
}

type ListenOption func(*ListenOptions)

func GetTimeCode() ListenOption {
	return func(l *ListenOptions) {
		l.TimeCode = true
	}
}

func GetActiveSense() ListenOption {
	return func(l *ListenOptions) {
		l.ActiveSense = true
	}
}

func GetSysEx() ListenOption {
	return func(l *ListenOptions) {
		l.SysEx = true
	}
}

func SysExBufferSize(size uint32) ListenOption {
	return func(l *ListenOptions) {
		l.SysExBufferSize = size
	}
}

var ErrPortClosed = drivers.ErrPortClosed
var ErrListenStopped = drivers.ErrListenStopped

// ListenTo listens on the given port number and passes the received MIDI data to the given receiver.
// It returns a stop function that may be called to stop the listening.
func ListenTo(portnumber int, recv Receiver, opts ...ListenOption) (stop func(), err error) {
	in, err := drivers.InByNumber(portnumber)
	if err != nil {
		return nil, err
	}

	var opt ListenOptions
	for _, o := range opts {
		o(&opt)
	}

	var conf drivers.ListenConfig
	conf.SysExBufferSize = opt.SysExBufferSize
	conf.TimeCode = opt.TimeCode
	conf.ActiveSense = opt.ActiveSense
	conf.SysEx = opt.SysEx

	if errrc, has := recv.(ErrorReceiver); has {
		conf.OnErr = errrc.OnError
	}

	var isStatusSet bool
	var typ, channel byte

	var onMsg = func(data []byte, millisec int32) {
		status := data[0]
		//var msg []byte
		var msg Message
		switch {
		// realtime message
		case status >= 0xF8:
			//msg = NewMessage([]byte{status})
			msg = []byte{status}
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
				//				msg = NewUndefined()
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
