package midi

import (
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

// listeningOptions are the options for the listening
type listeningOptions struct {

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

	// OnError handles occuring errors
	OnError func(error)
}

// Option is an option for listening
type Option func(*listeningOptions)

// UseTimeCode is an option to receive time code messages
func UseTimeCode() Option {
	return func(l *listeningOptions) {
		l.TimeCode = true
	}
}

// UseActiveSense is an option to receive active sense messages
func UseActiveSense() Option {
	return func(l *listeningOptions) {
		l.ActiveSense = true
	}
}

// UseSysEx is an option to receive system exclusive messages
func UseSysEx() Option {
	return func(l *listeningOptions) {
		l.SysEx = true
	}
}

// SysExBufferSize is an option to set the buffer size for sysex messages
func SysExBufferSize(size uint32) Option {
	return func(l *listeningOptions) {
		l.SysExBufferSize = size
	}
}

// HandleError sets an error handler when receiving messages
func HandleError(cb func(error)) Option {
	return func(l *listeningOptions) {
		l.OnError = cb
	}
}

var ErrPortClosed = drivers.ErrPortClosed
var ErrListenStopped = drivers.ErrListenStopped

// ListenTo listens on the given port and passes the received MIDI data to the given receiver.
// It returns a stop function that may be called to stop the listening.
func ListenTo(inPort drivers.In, recv func(msg Message, timestampms int32), opts ...Option) (stop func(), err error) {
	if !inPort.IsOpen() {
		err = inPort.Open()

		if err != nil {
			return nil, err
		}
	}

	var opt listeningOptions
	for _, o := range opts {
		o(&opt)
	}

	var conf drivers.ListenConfig
	conf.SysExBufferSize = opt.SysExBufferSize
	conf.TimeCode = opt.TimeCode
	conf.ActiveSense = opt.ActiveSense
	conf.SysEx = opt.SysEx
	conf.OnErr = opt.OnError

	var isStatusSet bool
	var typ, channel byte

	var onMsg = func(data []byte, millisec int32) {
		status := data[0]

		var msg Message
		switch {

		// realtime message
		case status >= 0xF8:
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

		// sysex message
		case status == 0xF0:
			isStatusSet = false
			msg = Message(data)

		// channel message
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

		recv(msg, millisec)
	}

	return inPort.Listen(onMsg, conf)
}
