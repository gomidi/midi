package portmididrv

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/portmididrv/imported/portmidi"
)

func newIn(driver *Driver, deviceid portmidi.DeviceID, id int, name string) drivers.In {
	return &in{driver: driver, id: id, name: name, deviceid: deviceid}
}

type in struct {
	deviceid portmidi.DeviceID
	stream   *portmidi.Stream
	id       int
	name     string
	driver   *Driver
}

// IsOpen returns wether the MIDI in port is open.
func (i *in) IsOpen() bool {
	return i.stream != nil
}

// Number returns the number of the MIDI in port.
// Since portmidis ports counting is confusing (out and in ports are counted together),
// we do our own counting.
func (i *in) Number() int {
	return i.id
}

// String returns the name of the MIDI in port.
func (i *in) String() string {
	return i.name
}

// Close closes the MIDI in port
func (i *in) Close() error {
	if i.stream == nil {
		return nil
	}
	//fmt.Println("closing in port")
	err := i.stream.Close()
	if err != nil {
		return fmt.Errorf("can't close MIDI in %v (%s): %v", i.Number(), i, err)
	}
	i.stream = nil
	return nil
}

// Open opens the MIDI in port
func (i *in) Open() (err error) {
	if i.stream != nil {
		return nil
	}

	i.stream, err = portmidi.NewInputStream(i.deviceid, i.driver.buffersizeIn)

	if err != nil {
		i.stream = nil
		return fmt.Errorf("can't open MIDI in port %v (%s): %v", i.Number(), i, err)
	}

	err = i.stream.SetFilter(0)
	if err != nil {
		fmt.Printf("can't set filter: %v\n", err.Error())
	}

	i.driver.opened = append(i.driver.opened, i)
	return nil
}

// Listen
func (i *in) Listen(onMsg func(msg []byte, milliseconds int32), config drivers.ListenConfig) (stopFn func(), err error) {

	if onMsg == nil {
		return nil, fmt.Errorf("onMsg callback must not be nil")
	}

	var filters []int

	if !config.ActiveSense {
		filters = append(filters, portmidi.FILTER_ACTIVE)
	}

	if !config.TimeCode {
		filters = append(filters, portmidi.FILTER_CLOCK)
	}

	switch len(filters) {
	case 1:
		i.stream.SetFilter(filters[0])
	case 2:
		i.stream.SetFilter(filters[0] | filters[1])
	default:
		// let all pass through
		i.stream.SetFilter(0)
	}

	lastTimestamp := portmidi.Time()
	var inSysEx bool
	if config.SysExBufferSize == 0 {
		config.SysExBufferSize = 1024
	}
	maxlenSysex := int(config.SysExBufferSize)
	var sysexBf = make([]byte, maxlenSysex)
	var sysexlen int
	var lastSysExTimestamp int32

	var stop int32

	stopWait := i.driver.sleepingTime * 2
	stopFn = func() {
		// lockless sync
		atomic.StoreInt32(&stop, 1)
		time.Sleep(stopWait)
	}

	go func() {
		var stopped int32

		defer func() {
			if config.OnErr != nil {
				config.OnErr(drivers.ErrListenStopped)
			}
		}()

		for {

			// lockless sync
			stopped = atomic.LoadInt32(&stop)

			if stopped == 1 {
				return
			}

			events, err2 := i.stream.Read(i.driver.buffersizeRead)

			/*
									from portmidi docs:

				   Note that MIDI allows nested messages: the so-called "real-time" MIDI
				   messages can be inserted into the MIDI byte stream at any location,
				   including within a sysex message. MIDI real-time messages are one-byte
				   messages used mainly for timing (see the MIDI spec). PortMidi retains
				   the order of non-real-time MIDI messages on both input and output, but
				   it does not specify exactly how real-time messages are processed. This
				   is particulary problematic for MIDI input, because the input parser
				   must either prepare to buffer an unlimited number of sysex message
				   bytes or to buffer an unlimited number of real-time messages that
				   arrive embedded in a long sysex message. To simplify things, the input
				   parser is allowed to pass real-time MIDI messages embedded within a
				   sysex message, and it is up to the client to detect, process, and
				   remove these messages as they arrive.

				   When receiving sysex messages, the sysex message is terminated
				   by either an EOX status byte (anywhere in the 4 byte messages) or
				   by a non-real-time status byte in the low order byte of the message.
				   If you get a non-real-time status byte but there was no EOX byte, it
				   means the sysex message was somehow truncated. This is not
				   considered an error; e.g., a missing EOX can result from the user
				   disconnecting a MIDI cable during sysex transmission.

				   A real-time message can occur within a sysex message. A real-time
				   message will always occupy a full PmEvent with the status byte in
				   the low-order byte of the PmEvent message field. (This implies that
				   the byte-order of sysex bytes and real-time message bytes may not
				   be preserved -- for example, if a real-time message arrives after
				   3 bytes of a sysex message, the real-time message will be delivered
				   first. The first word of the sysex message will be delivered only
				   after the 4th byte arrives, filling the 4-byte PmEvent message field.

				   [...]

				   On input, the timestamp ideally denotes the arrival time of the
				   status byte of the message. The first timestamp on sysex message
				   data will be valid. Subsequent timestamps may denote
				   when message bytes were actually received, or they may be simply
				   copies of the first timestamp.

				   Timestamps for nested messages: If a real-time message arrives in
				   the middle of some other message, it is enqueued immediately with
				   the timestamp corresponding to its arrival time. The interrupted
				   non-real-time message or 4-byte packet of sysex data will be enqueued
				   later. The timestamp of interrupted data will be equal to that of
				   the interrupting real-time message to insure that timestamps are
				   non-decreasing.
			*/

			if err2 != nil && config.OnErr != nil {
				config.OnErr(err2)
			}

			if err2 == nil {

				for _, ev := range events {

					// ev.Timestamp is in Milliseconds
					ts := int32(ev.Timestamp - lastTimestamp)

					switch {

					// realtime message
					case ev.Status >= 0xF8:
						onMsg([]byte{byte(ev.Status)}, ts)

					// sysex
					case ev.Status == 0xF0 || inSysEx:
						if ev.Status == 0xF0 {
							inSysEx = true
							lastSysExTimestamp = ts
							// really? to cancel the running status, send the 0xF7 to the onMsg callback
							// onMsg([3]byte{0xF7, 0, 0}, ts)
						}

						bt := []byte{byte(ev.Status), byte(ev.Data1), byte(ev.Data2), byte(ev.Rest)}

						if config.SysEx {
							for _, b := range bt {
								if sysexlen >= maxlenSysex {
									if config.OnErr != nil {
										config.OnErr(fmt.Errorf("skip sysex; sysex buffer size is too small (%v byte)", config.SysExBufferSize))
									}
									break
								}
								sysexBf[sysexlen] = b
								sysexlen++
								if b == 0xF7 {
									if sysexlen >= maxlenSysex {
										sysexlen = 0
										inSysEx = false
										break
									}
									go func(bb []byte, l int) {
										var _bt = make([]byte, l)

										for i := 0; i < l; i++ {
											_bt[i] = bb[i]
										}
										onMsg(_bt, lastSysExTimestamp)
									}(sysexBf, sysexlen)
									sysexlen = 0
									inSysEx = false
									break
								}
							}
						} else {
							if pos := bytes.IndexByte(bt, 0xF7); pos >= 0 {
								inSysEx = false
							}
						}

					// [MIDI] permits 0xF7 octets that are not part of a (0xF0, 0xF7) pair
					// to appear on a MIDI 1.0 DIN cable.  Unpaired 0xF7 octets have no
					// semantic meaning in MIDI apart from cancelling running status.
					case ev.Status == 0xF7:
						if config.SysEx {
							sysexlen = 0
						}
						// to allow cancelling running status, send the 0xF7
						onMsg([]byte{0xF7}, ts)

					// all other messages (syscommon and channel messages)
					// TODO: only transmit the necessary bytes
					default:
						onMsg([]byte{byte(ev.Status), byte(ev.Data1), byte(ev.Data2)}, ts)
					}
				}

			}
			time.Sleep(i.driver.sleepingTime)
		}
	}()
	time.Sleep(time.Millisecond * 2)

	return
}

// Underlying returns the underlying *portmidi.Stream. It will be nil, if the port is closed.
// Use it with type casting:
//
//	portOut := o.Underlying().(*portmidi.Stream)
func (i *in) Underlying() interface{} {
	return i.stream
}
