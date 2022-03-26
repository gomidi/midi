package rtmididrv

import (
	"fmt"
	"math"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv/imported/rtmidi"
)

type in struct {
	number int
	//sync.RWMutex
	//listenerSet bool
	driver *Driver
	name   string
	midiIn rtmidi.MIDIIn
}

// IsOpen returns wether the MIDI in port is open
func (i *in) IsOpen() (open bool) {
	//	i.RLock()
	open = i.midiIn != nil
	//i.RUnlock()
	return
}

// String returns the name of the MIDI in port.
func (i *in) String() string {
	return i.name
}

// Underlying returns the underlying rtmidi.MIDIIn. Use it with type casting:
//   rtIn := i.Underlying().(rtmidi.MIDIIn)
func (i *in) Underlying() interface{} {
	return i.midiIn
}

// Number returns the number of the MIDI in port.
// Note that with rtmidi, out and in ports are counted separately.
// That means there might exists out ports and an in ports that share the same number.
func (i *in) Number() int {
	return i.number
}

// Close closes the MIDI in port, after it has stopped listening.
func (i *in) Close() (err error) {
	if !i.IsOpen() {
		return nil
	}

	//	i.StopListening()
	//	i.Lock()
	err = i.midiIn.Close()
	i.midiIn = nil
	//	i.Unlock()
	return
}

// Open opens the MIDI in port
func (i *in) Open() (err error) {
	if i.IsOpen() {
		return nil
	}

	//i.Lock()

	i.midiIn, err = rtmidi.NewMIDIInDefault()
	if err != nil {
		i.midiIn = nil
		//i.Unlock()
		return fmt.Errorf("can't open default MIDI in: %v", err)
	}

	err = i.midiIn.OpenPort(i.number, "")
	//i.Unlock()

	if err != nil {
		i.Close()
		return fmt.Errorf("can't open MIDI in port %v (%s): %v", i.number, i, err)
	}

	//i.driver.Lock()
	//i.midiIn.IgnoreTypes(i.driver.ignoreSysex, i.driver.ignoreTimeCode, i.driver.ignoreActiveSense)
	i.driver.opened = append(i.driver.opened, i)
	//i.driver.Unlock()

	return nil
}

/*
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
*/

func newIn(driver *Driver, number int, name string) drivers.In {
	return &in{driver: driver, number: number, name: name}
}

func (i *in) Listen(onMsg func(msg []byte, milliseconds int32), config drivers.ListenConfig) (stopFn func(), err error) {

	if onMsg == nil {
		return nil, fmt.Errorf("onMsg callback must not be nil")
	}

	i.midiIn.IgnoreTypes(!config.SysEx, !config.TimeCode, !config.ActiveSense)

	//var inSysEx bool
	if config.SysExBufferSize == 0 {
		config.SysExBufferSize = 1024
	}

	/*
		maxlenSysex := int(config.SysExBufferSize)
		var sysexBf = make([]byte, maxlenSysex)
		var sysexlen int

		var ts_ms int32
		//var stop int32
		var state = readerStateClean
		var statusByte uint8
		//var channel uint8
		//var bf [2]byte // first: is set, second: the byte
		var issetBf bool
		var bf byte
		var typ uint8
	*/

	var rd = drivers.NewReader(config, onMsg)
	/*
		rd.OnErr = config.OnErr
		rd.OnSysEx = config.OnSysEx
		rd.OnMsg = onMsg
		rd.SysExBufferSize = config.SysExBufferSize
	*/

	//stopWait := i.driver.sleepingTime * 2
	stopFn = func() {
		// lockless sync
		//	atomic.StoreInt32(&stop, 1)
		//	fmt.Println("stopping")
		i.midiIn.CancelCallback()
		//time.Sleep(stopWait)
	}

	/*
		withinChannelMessage := func(b byte) {
			switch typ {
			case byteChannelPressure:
				issetBf = false
				state = readerStateClean
				//p.receiver.Receive(Channel(p.channel).Aftertouch(b), p.timestamp)
				onMsg([3]byte{statusByte, b, 0}, ts_ms)
			case byteProgramChange:
				issetBf = false // first: is set, second: the byte
				state = readerStateClean
				//p.receiver.Receive(Channel(p.channel).ProgramChange(b), p.timestamp)
				onMsg([3]byte{statusByte, b, 0}, ts_ms)
			case byteControlChange:
				if issetBf {
					issetBf = false // first: is set, second: the byte
					state = readerStateClean
					//p.receiver.Receive(Channel(p.channel).ControlChange(p.getBf(), b), p.timestamp)
					onMsg([3]byte{statusByte, b, bf}, ts_ms)
				} else {
					issetBf = true
					bf = b
				}
			case byteNoteOn:
				if issetBf {
					issetBf = false // first: is set, second: the byte
					state = readerStateClean
					//p.receiver.Receive(Channel(p.channel).NoteOn(p.getBf(), b), p.timestamp)
					onMsg([3]byte{statusByte, b, bf}, ts_ms)
				} else {
					issetBf = true
					bf = b
				}
			case byteNoteOff:
				if issetBf {
					issetBf = false // first: is set, second: the byte
					state = readerStateClean
					//p.receiver.Receive(Channel(p.channel).NoteOffVelocity(p.getBf(), b), p.timestamp)
					onMsg([3]byte{statusByte, b, bf}, ts_ms)
				} else {
					issetBf = true
					bf = b
				}
			case bytePolyphonicKeyPressure:
				if issetBf {
					issetBf = false // first: is set, second: the byte
					state = readerStateClean
					//p.receiver.Receive(Channel(p.channel).PolyAftertouch(p.getBf(), b), p.timestamp)
					onMsg([3]byte{statusByte, b, bf}, ts_ms)
				} else {
					issetBf = true
					bf = b
				}
			case bytePitchWheel:
				if issetBf {
					//rel, abs := midilib.ParsePitchWheelVals(bf, b)
					//_ = abs
					issetBf = false // first: is set, second: the byte
					state = readerStateClean
					//p.receiver.Receive(Channel(p.channel).Pitchbend(rel), p.timestamp)
					onMsg([3]byte{statusByte, b, bf}, ts_ms)
				} else {
					issetBf = true
					bf = b
				}
			default:
				panic("unknown typ")
			}
		}
	*/

	/*
		cleanState := func(b byte) {
			switch {

			// start sysex
			case b == 0xF0:
				statusByte = 0
				sysexBf = make([]byte, maxlenSysex)
				//sysexBf.Reset()
				//sysexBf.WriteByte(b)
				sysexBf[0] = b
				sysexlen = 1
				state = readerStateInSysEx
			// end sysex
			// [MIDI] permits 0xF7 octets that are not part of a (0xF0, 0xF7) pair
			// to appear on a MIDI 1.0 DIN cable.  Unpaired 0xF7 octets have no
			// semantic meaning in MIDI apart from cancelling running status.
			case b == 0xF7:
				sysexBf = nil
				sysexlen = 0
				statusByte = 0
				onMsg([3]byte{b, 0, 0}, ts_ms)

			// here we clear for System Common Category messages
			case b > 0xF0 && b < 0xF7:
				statusByte = 0
				issetBf = false // reset buffer
				//fmt.Printf("sys common msg started\n")
				switch b {
				case byteMIDITimingCodeMessage, byteSysSongPositionPointer, byteSysSongSelect:
					state = readerStateWithinSysCommon
					typ = b
				case byteSysTuneRequest:
					onMsg([3]byte{b, 0, 0}, ts_ms)
					//
					//	if p.syscommonHander != nil {
					//		p.syscommonHander(Tune(), p.timestamp)
					//	}

					return
				default:
					// 0xF4, 0xF5, or 0xFD
					state = readerStateWithinUnknown
					return
				}

			// channel message with status byte
			case b >= 0x80 && b <= 0xEF:
				statusByte = b
				issetBf = false // reset buffer
				//typ, channel = midilib.ParseStatus(statusByte)
				typ, _ = midilib.ParseStatus(statusByte)
				state = readerStateWithinChannelMessage
			default:
				if statusByte != 0 {
					state = readerStateWithinChannelMessage
					withinChannelMessage(b)
				}
			}
		}
	*/

	go i.midiIn.SetCallback(func(in rtmidi.MIDIIn, bt []byte, deltaSeconds float64) {
		/*
			var stopped int32

			// lockless sync
			stopped = atomic.LoadInt32(&stop)

			if stopped == 1 {
				if config.OnErr != nil {
					config.OnErr(drivers.ErrListenStopped)
				}
				fmt.Printf("stopped")
				in.CancelCallback()
				return
			}
		*/

		rd.EachMessage(bt, int32(math.Round(deltaSeconds*1000)))

		/*
			// TODO: verify
			// assume that each call is without running state
			statusByte = 0
			issetBf = false // first: is set, second: the byte

			ts_ms += int32(math.Round(deltaSeconds * 1000))

			//	fmt.Printf("got % X\n", bt)

			for _, b := range bt {
				// => realtime message
				if b >= 0xF8 {
					onMsg([3]byte{b, 0, 0}, ts_ms)
					continue
				}

				//fmt.Printf("state: %v\n", p.state)

				switch state {
				case readerStateInSysEx:
					// interrupted sysex, discard old data
					if b == 0xF0 {
						statusByte = 0
						sysexBf = make([]byte, maxlenSysex)
						//sysexBf.Reset()
						//sysexBf.WriteByte(b)
						sysexBf[0] = b
						sysexlen = 1
						state = readerStateInSysEx
						continue
					}

					if b == 0xF7 {
						state = readerStateClean
						if config.OnSysEx != nil {
							sysexBf[sysexlen] = b
							sysexlen++
							go func(bb []byte, l int) {
								var _bt = make([]byte, l)

								for i := 0; i < l; i++ {
									_bt[i] = bb[i]
								}
								config.OnSysEx(_bt)
							}(sysexBf, sysexlen)
						}
						sysexBf = nil
						sysexlen = 0
						continue
					}
					if midilib.IsStatusByte(b) {
						//p.sysexBf.Reset()
						sysexBf = nil
						sysexlen = 0
						state = readerStateClean
						cleanState(b)
						continue
					}

					if config.OnSysEx != nil {
						sysexBf[sysexlen] = b
						sysexlen++
					}

				case readerStateClean:
					cleanState(b)
				case readerStateWithinUnknown:
					//p.withinUnknown(b)
					if midilib.IsStatusByte(b) {
						state = readerStateClean
						cleanState(b)
					}
				case readerStateWithinSysCommon:
					switch typ {
					case byteMIDITimingCodeMessage:
						issetBf = false
						state = readerStateClean
						onMsg([3]byte{typ, b, 0}, ts_ms)
					case byteSysSongPositionPointer:
						if issetBf {
							issetBf = false
							state = readerStateClean
							onMsg([3]byte{typ, b, bf}, ts_ms)
						} else {
							issetBf = true
							bf = b
						}
					case byteSysSongSelect:
						issetBf = false
						state = readerStateClean
						onMsg([3]byte{typ, b, 0}, ts_ms)
					case byteSysTuneRequest:
						//panic("must not be handled here, but within clean state")
					default:
						if config.OnErr != nil {
							config.OnErr(fmt.Errorf("unknown syscommon message: % X", b))
						}
						//panic("unknown syscommon")
					}
				case readerStateWithinChannelMessage:
					withinChannelMessage(b)
				default:
					panic(fmt.Sprintf("unknown state %v, must not happen", state))
				}
			}
		*/

	})

	return stopFn, nil
}

/*
func (i *in) StartListening(callback func(data []byte, deltadecimilliseconds int32)) error {
	if !i.IsOpen() {
		//fmt.Printf("post closed\n")
		return drivers.ErrPortClosed
	}

	i.RLock()
	if i.listenerSet {
		i.RUnlock()
		return fmt.Errorf("listener already set")
	}
	i.RUnlock()
	//fmt.Println("pre lock")
	i.Lock()
	i.listenerSet = true
	i.driver.Lock()
	//i.midiIn.IgnoreTypes(i.driver.ignoreSysex, i.driver.ignoreTimeCode, i.driver.ignoreActiveSense)
	i.driver.Unlock()
	i.Unlock()

	var tsdecimilliseconds int32

	// since i.midiIn.SetCallback is blocking on success, there is no meaningful way to get an error
	// and set the callback non blocking
	go i.midiIn.SetCallback(func(_ rtmidi.MIDIIn, bt []byte, deltaSeconds float64) {
		// convert to milliseconds (10^-5)
		tsdecimilliseconds += int32(math.Round(deltaSeconds * 1000))
		callback(bt, tsdecimilliseconds)
	})

	time.Sleep(time.Millisecond * 10)

	return nil
}
*/

/*
// StopListening cancels the listening
func (i *in) StopListening() (err error) {
	if !i.IsOpen() {
		return drivers.ErrPortClosed
	}
	i.Lock()
	if i.listenerSet {
		i.listenerSet = false
		err = i.midiIn.CancelCallback()
		if err != nil {
			err = fmt.Errorf("can't stop listening on MIDI in port %v (%s): %v", i.number, i, err)
		}
	}
	i.Unlock()
	return
}
*/
