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
		return ch.NewAfterTouch(data1)
	case byteProgramChange:
		return ch.NewProgramChange(data1)
	case byteControlChange:
		return ch.NewControlChange(data1, data2)
	case byteNoteOn:
		return ch.NewNoteOn(data1, data2)
	case byteNoteOff:
		return ch.NewNoteOffVelocity(data1, data2)
	case bytePolyphonicKeyPressure:
		return ch.NewPolyAfterTouch(data1, data2)
	case bytePitchWheel:
		rel, _ := midilib.ParsePitchWheelVals(data1, data2)
		return ch.NewPitchbend(rel)
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

	SysEx bool
}

var ErrPortClosed = drivers.ErrPortClosed
var ErrListenStopped = drivers.ErrListenStopped

// ListenToPort listens on the given port number and passes the received MIDI data to the given receiver.
// It returns a stop function that may be called to stop the listening.
func ListenToPort(portnumber int, recv Receiver, opt ListenOptions) (stop func(), err error) {
	in, err := drivers.InByNumber(portnumber)
	if err != nil {
		return nil, err
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
		var msg Msg
		switch {
		// realtime message
		case status >= 0xF8:
			msg = NewMsg([]byte{status})
		// here we clear for System Common Category messages
		case status > 0xF0 && status < 0xF7:
			isStatusSet = false
			switch status {
			case byteSysTuneRequest:
				msg = NewTune()
			case byteMIDITimingCodeMessage:
				msg = NewMTC(data[1])
			case byteSysSongPositionPointer:
				_, abs := midilib.ParsePitchWheelVals(data[1], data[2])
				msg = NewSPP(abs)
			case byteSysSongSelect:
				msg = NewSongSelect(data[1])
			default:
				// undefined syscommon message
				msg = NewUndefined()
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
