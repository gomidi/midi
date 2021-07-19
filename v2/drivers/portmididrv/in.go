package portmididrv

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/portmididrv/imported/portmidi"
)

func newIn(driver *Driver, deviceid portmidi.DeviceID, id int, name string) drivers.In {
	return &in{driver: driver, id: id, name: name, deviceid: deviceid}
}

type in struct {
	deviceid portmidi.DeviceID
	id       int
	stream   *portmidi.Stream
	name     string

	driver  *Driver
	inSysEx bool

	//lastTimestamp portmidi.Timestamp
	mx            sync.RWMutex
	sysexcallback func(data []byte)
	sysexBf       bytes.Buffer
	stopped       chan bool
	shouldstop    bool
}

// IsOpen returns wether the MIDI in port is open.
func (i *in) IsOpen() bool {
	i.mx.RLock()
	defer i.mx.RUnlock()
	return i.stream != nil
}

// Underlying returns the underlying *portmidi.Stream. It will be nil, if the port is closed.
// Use it with type casting:
//   portIn := i.Underlying().(*portmidi.Stream)
func (i *in) Underlying() interface{} {
	return i.stream
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
	i.mx.RLock()
	if i.stream == nil {
		i.mx.RUnlock()
		return nil
	}
	i.mx.RUnlock()

	err := i.StopListening()
	if err != nil {
		return err
	}

	//fmt.Println("did stop, now closing stream")

	// to prevent the race between polling and closing
	// give the listener some time to stop listening
	//time.Sleep(i.driver.sleepingTime * 2)

	i.mx.Lock()
	defer i.mx.Unlock()
	err = i.stream.Close()
	if err != nil {
		return fmt.Errorf("can't close MIDI in %v (%s): %v", i.Number(), i, err)
	}
	i.stream = nil
	return nil
}

// Open opens the MIDI in port
func (i *in) Open() (err error) {
	i.mx.RLock()
	if i.stream != nil {
		i.mx.RUnlock()
		return nil
	}
	i.mx.RUnlock()

	i.mx.Lock()
	defer i.mx.Unlock()
	i.stream, err = portmidi.NewInputStream(i.deviceid, i.driver.buffersizeIn)
	if err != nil {
		i.stream = nil
		return fmt.Errorf("can't open MIDI in port %v (%s): %v", i.Number(), i, err)
	}
	//i.stream.SetFilter(portmidi.FILTER_ACTIVE)
	//err = i.stream.SetFilter(portmidi.FILTER_NOTE | portmidi.FILTER_ACTIVE)
	//err = i.stream.SetFilter(portmidi.FILTER_PITCHBEND | portmidi.FILTER_NOTE)
	err = i.stream.SetFilter(0)
	if err != nil {
		fmt.Printf("can't set filter: %v\n", err.Error())
	}
	i.driver.Lock()
	defer i.driver.Unlock()
	i.driver.opened = append(i.driver.opened, i)
	return nil
}

// StopListening cancels the listening
func (i *in) StopListening() error {
	i.mx.RLock()
	shouldstop := i.shouldstop
	i.mx.RUnlock()
	if shouldstop {
		return nil
	}
	i.mx.Lock()
	i.stopped = make(chan bool)
	i.shouldstop = true
	i.mx.Unlock()
	<-i.stopped
	return nil
}

// read is an internal helper function
func (i *in) read(cb func([]byte, int32)) error {

	//fmt.Printf("reading\n")

	//events, err := i.stream.Read(int(i.driver.buffersizeIn))
	i.mx.Lock()
	events, err := i.stream.Read(i.driver.buffersizeRead)
	i.mx.Unlock()

	// TODO
	// if an error occurs (buffer overflow), discard any sysexmessage
	// BTW that should be done automatically with the next status byte
	if err != nil {
		fmt.Printf("error when reading: %s\n", err)
		return err
	}

	//fmt.Printf("reader: got %v events\n", len(events))
	hassysexcallback := i.sysexcallback != nil
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

	for _, ev := range events {

		// => realtime message
		if ev.Status >= 0xF8 {
			// ev.Timestamp is in Milliseconds
			// we want deciMilliseconds as int32
			ts := int32(ev.Timestamp * 10)

			cb([]byte{byte(ev.Status)}, ts)
			continue
		}

		var inSysEx bool

		// start a new sysex
		if ev.Status == 0xF0 {
			i.mx.Lock()
			i.inSysEx = true
			inSysEx = true
			if hassysexcallback {
				i.sysexBf.Reset()
			}
			i.mx.Unlock()
		} else {
			i.mx.RLock()
			inSysEx = i.inSysEx
			i.mx.RUnlock()
		}

		//fmt.Printf("status: %X in sysex: %v\n", ev.Status, inSysEx)

		/*
			if pos := bytes.IndexByte(s.sysexBuffer[copied:copied+4], 0xF7); pos >= 0 {
							size := copied + pos + 1
							event.SysEx = make([]byte, size)
							event.Data1 = 0
							event.Data2 = 0
							copy(event.SysEx, s.sysexBuffer[:size])
							break
						}
		*/

		if inSysEx {
			bt := []byte{byte(ev.Status), byte(ev.Data1), byte(ev.Data2), byte(ev.Rest)}

			if pos := bytes.IndexByte(bt, 0xF7); pos >= 0 {
				i.mx.Lock()
				i.inSysEx = false
				if hassysexcallback {
					i.sysexBf.Write(bt[:pos+1])
					i.sysexcallback(i.sysexBf.Bytes())
					i.sysexBf.Reset()
				}
				i.mx.Unlock()
			} else {
				if hassysexcallback {
					i.mx.Lock()
					i.sysexBf.Write(bt)
					i.mx.Unlock()
				}
			}

			continue
		}
		// not in sysex

		// ignore errorneous sysex finish
		if ev.Status == 0xF7 {
			if hassysexcallback {
				i.mx.Lock()
				i.sysexBf.Reset()
				i.mx.Unlock()
			}
			continue
		}

		// in the start or middle of a sysex
		if inSysEx {

			continue
		}

		// ev.Timestamp is in Milliseconds
		// we want deciMilliseconds as int32
		ts := int32(ev.Timestamp * 10)

		// all other messages (syscommon and channel messages)
		var b = make([]byte, 3)
		b[0] = byte(ev.Status)
		b[1] = byte(ev.Data1)
		b[2] = byte(ev.Data2)
		// ev.Timestamp is in Milliseconds
		cb(b, ts)
	}

	return nil
}

// StartListening
func (i *in) StartListening(callback func(data []byte, deltadecimilliseconds int32)) error {
	time.Sleep(time.Millisecond)
	//i.lastTimestamp = portmidi.Time()
	go func() {
		//fmt.Println("start reading")
		for {

			i.mx.RLock()
			shouldstop := i.shouldstop
			i.mx.RUnlock()

			if shouldstop {
				i.stopped <- true
				//	fmt.Println("stopped")
				return
			}

			//i.driver.Lock()
			//has, _ := i.stream.Poll()
			//i.driver.Unlock()

			/*
				has, err := i.stream.Poll()
				if err != nil {
					fmt.Printf("error while polling: %s\n", err.Error())
				}
			*/

			//if has {
			//i.mx.Lock()
			i.read(callback)
			//i.mx.Unlock()
			//}

			/*
				if err := i.stream.HasHostError(); err != nil {
					fmt.Printf("hosterror: %v\n", err.Error())
				}
			*/
			time.Sleep(i.driver.sleepingTime)
		}
	}()
	return nil
}

func (i *in) StartListeningForSysEx(callback func(data []byte)) error {
	i.mx.Lock()
	i.sysexcallback = callback
	i.mx.Unlock()
	return nil
}
