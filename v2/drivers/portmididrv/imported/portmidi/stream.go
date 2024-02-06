// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package portmidi

// #cgo LDFLAGS: -lportmidi
//
// #include <stdlib.h>
// #include <portmidi.h>
import "C"

import (
	//"bytes"
	"encoding/hex"
	"errors"
	"strings"
	"time"
	"unsafe"
)

const (
	minEventBufferSize = 1
	maxEventBufferSize = 1024
)

var (
	ErrMaxBuffer         = errors.New("portmidi: max event buffer size is 1024")
	ErrMinBuffer         = errors.New("portmidi: min event buffer size is 1")
	ErrInputUnavailable  = errors.New("portmidi: input is unavailable")
	ErrOutputUnavailable = errors.New("portmidi: output is unavailable")
	ErrSysExOverflow     = errors.New("portmidi: SysEx message overflowed")
)

// Channel represent a MIDI channel. It should be between 1-16.
type Channel int

// Event represents a MIDI event.
type Event struct {
	Timestamp Timestamp
	Status    int64
	Data1     int64
	Data2     int64
	Rest      int64
	//SysEx     []byte
}

// Stream represents a portmidi stream.
type Stream struct {
	deviceID DeviceID
	pmStream *C.PmStream

	sysexBuffer [maxEventBufferSize]byte
}

// NewInputStream initializes a new input stream.
func NewInputStream(id DeviceID, bufferSize int64) (stream *Stream, err error) {
	var str *C.PmStream
	errCode := C.Pm_OpenInput(
		(*unsafe.Pointer)(unsafe.Pointer(&str)),
		C.PmDeviceID(id), nil, C.int32_t(bufferSize), nil, nil)
	if errCode != 0 {
		return nil, convertToError(errCode)
	}
	if info := Info(id); !info.IsInputAvailable {
		return nil, ErrInputUnavailable
	}
	return &Stream{deviceID: id, pmStream: str}, nil
}

// NewOutputStream initializes a new output stream.
func NewOutputStream(id DeviceID, bufferSize int64, latency int64) (stream *Stream, err error) {
	var str *C.PmStream
	errCode := C.Pm_OpenOutput(
		(*unsafe.Pointer)(unsafe.Pointer(&str)),
		C.PmDeviceID(id), nil, C.int32_t(bufferSize), nil, nil, C.int32_t(latency))
	if errCode != 0 {
		return nil, convertToError(errCode)
	}
	if info := Info(id); !info.IsOutputAvailable {
		return nil, ErrOutputUnavailable
	}
	return &Stream{deviceID: id, pmStream: str}, nil
}

// Close closes the MIDI stream.
func (s *Stream) Close() error {
	if s.pmStream == nil {
		return nil
	}
	return convertToError(C.Pm_Close(unsafe.Pointer(s.pmStream)))
}

// Abort aborts the MIDI stream.
func (s *Stream) Abort() error {
	if s.pmStream == nil {
		return nil
	}
	return convertToError(C.Pm_Abort(unsafe.Pointer(s.pmStream)))
}

// Write writes a buffer of MIDI events to the output stream.
func (s *Stream) Write(events []Event) error {
	size := len(events)
	if size > maxEventBufferSize {
		return ErrMaxBuffer
	}
	buffer := make([]C.PmEvent, size)
	for i, evt := range events {
		var event C.PmEvent
		//event.timestamp = C.PmTimestamp(evt.Timestamp)
		event.timestamp = C.int(pmTimestamp(evt.Timestamp))
		event.message = C.PmMessage((((evt.Data2 << 16) & 0xFF0000) | ((evt.Data1 << 8) & 0xFF00) | (evt.Status & 0xFF)))
		buffer[i] = event
	}
	return convertToError(C.Pm_Write(unsafe.Pointer(s.pmStream), &buffer[0], C.int32_t(size)))
}

// WriteShort writes a MIDI event of three bytes immediately to the output stream.
func (s *Stream) WriteShort(status int64, data1 int64, data2 int64) error {
	evt := Event{
		//Timestamp: Timestamp(C.Pt_Time()),
		Timestamp: Time(),
		Status:    status,
		Data1:     data1,
		Data2:     data2,
	}
	return s.Write([]Event{evt})
}

// WriteSysExBytes writes a system exclusive MIDI message given as a []byte to the output stream.
func (s *Stream) WriteSysExBytes(when Timestamp, msg []byte) error {
	return convertToError(C.Pm_WriteSysEx(unsafe.Pointer(s.pmStream), C.PmTimestamp(when), (*C.uchar)(unsafe.Pointer(&msg[0]))))
}

// WriteSysEx writes a system exclusive MIDI message given as a string of hexadecimal characters to
// the output stream. The string must only consist of hex digits (0-9A-F) and optional spaces. This
// function is case-insenstive.
func (s *Stream) WriteSysEx(when Timestamp, msg string) error {
	buf, err := hex.DecodeString(strings.Replace(msg, " ", "", -1))
	if err != nil {
		return err
	}

	return s.WriteSysExBytes(when, buf)
}

// SetChannelMask filters incoming stream based on channel.
// In order to filter from more than a single channel, or multiple channels.
// s.SetChannelMask(Channel(1) | Channel(10)) will both filter input
// from channel 1 and 10.
func (s *Stream) SetChannelMask(mask int) error {
	return convertToError(C.Pm_SetChannelMask(unsafe.Pointer(s.pmStream), C.int(mask)))
}

// Reads from the input stream, the max number events to be read are
// determined by max.
func (s *Stream) Read(max int) (events []Event, err error) {
	if max > maxEventBufferSize {
		return nil, ErrMaxBuffer
	}
	if max < minEventBufferSize {
		return nil, ErrMinBuffer
	}
	buffer := make([]C.PmEvent, max)
	numEvents := int(C.Pm_Read(unsafe.Pointer(s.pmStream), &buffer[0], C.int32_t(max)))
	if numEvents < 0 {
		return nil, convertToError(C.PmError(numEvents))
	}

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

	//fmt.Printf("PM got %v events\n", numEvents)

	events = make([]Event, 0, numEvents)
	for i := 0; i < numEvents; i++ {

		/*
			   typedef int32_t PmMessage;
				#define Pm_MessageStatus(msg) ((msg) & 0xFF)
				#define Pm_MessageData1(msg) (((msg) >> 8) & 0xFF)
				#define Pm_MessageData2(msg) (((msg) >> 16) & 0xFF)
		*/

		event := Event{
			Timestamp: Timestamp(buffer[i].timestamp),
			Status:    int64(buffer[i].message) & 0xFF,
			Data1:     (int64(buffer[i].message) >> 8) & 0xFF,
			Data2:     (int64(buffer[i].message) >> 16) & 0xFF,
			Rest:      (int64(buffer[i].message) >> 24) & 0xFF,
		}

		//fmt.Printf("event: %X %X %X %X\n", event.Status, event.Data1, event.Data2, event.Rest)

		//if event.Status&0xF0 == 0xF0 {
		// TODO we need to keep state here, since any sysex message has only the first event with the status
		// but might span across multiple buffers
		// also realtime messages might come in between
		// but since we are all handling this inside the reader,
		// it would be the best to simply put all the bytes together and pass them as is (i.e. without special sysex handling here)

		/*
			// rubbish
			if event.Status == 0xF0 {
				// Sysex message starts with 0xF0, ends with 0xF7
				read := 0
				for i+read < numEvents {
					copied := read * 4

					s.sysexBuffer[copied+0] = byte(buffer[i+read].message & 0xFF)
					s.sysexBuffer[copied+1] = byte((buffer[i+read].message >> 8) & 0xFF)
					s.sysexBuffer[copied+2] = byte((buffer[i+read].message >> 16) & 0xFF)
					s.sysexBuffer[copied+3] = byte((buffer[i+read].message >> 24) & 0xFF)

					if pos := bytes.IndexByte(s.sysexBuffer[copied:copied+4], 0xF7); pos >= 0 {
						size := copied + pos + 1
						event.SysEx = make([]byte, size)
						event.Data1 = 0
						event.Data2 = 0
						copy(event.SysEx, s.sysexBuffer[:size])
						break
					}

					read++
				}
				if event.SysEx == nil {
					// We didn't find a 0xF7, meaning the
					// event buffer was not large enough.
					// Comments on Pm_Read() indicate that
					// when a large SysEx message is not
					// fully received, the reader will
					// flush the buffer to avoid the next
					// read starting in the middle of the
					// unread SysEx message bytes.
					return nil, ErrSysExOverflow
				}
				i += read
			}
		*/

		events = append(events, event)
	}

	//fmt.Printf("sending events\n")
	return
}

// ReadSysExBytes reads 4*max sysex bytes from the input stream.
//
// Deprecated. Using this API may cause a loss of buffered data.  It
// is preferable to use Read() and inspect the Event.SysEx field to
// detect SysEx messages.
/*
func (s *Stream) ReadSysExBytes(max int) ([]byte, error) {
	evt, err := s.Read(max)
	if err != nil {
		return nil, err
	}
	return evt[0].SysEx, nil
}
*/

// Listen input stream for MIDI events.
func (s *Stream) Listen() <-chan Event {
	ch := make(chan Event)
	go func(s *Stream, ch chan Event) {
		for {
			// sleep for a while before the new polling tick,
			// otherwise operation is too intensive and blocking
			time.Sleep(10 * time.Millisecond)
			events, err := s.Read(maxEventBufferSize)
			// Note: It's not very reasonable to push sliced data into
			// a channel, several perf penalities there are.
			// This function is added as a handy utility.
			if err != nil {
				continue
			}
			for i := range events {
				ch <- events[i]
			}
		}
	}(s, ch)
	return ch
}

// HasHostError returns an async host error
func (s *Stream) HasHostError() error {
	hosterr := C.Pm_HasHostError(unsafe.Pointer(s.pmStream))
	if hosterr < 0 {
		return convertToError(C.PmError(hosterr))
	}
	return nil
}

// Poll reports whether there is input available in the stream.
func (s *Stream) Poll() (bool, error) {
	poll := C.Pm_Poll(unsafe.Pointer(s.pmStream))
	if poll < 0 {
		return false, convertToError(C.PmError(poll))
	}
	return poll > 0, nil
}

type Filter int

var FILTER_ACTIVE = C.PM_FILT_ACTIVE
var FILTER_SYSEX = C.PM_FILT_SYSEX
var FILTER_CLOCK = C.PM_FILT_CLOCK
var FILTER_PLAY = C.PM_FILT_PLAY
var FILTER_TICK = C.PM_FILT_TICK
var FILTER_FD = C.PM_FILT_FD
var FILTER_UNDEFINED = C.PM_FILT_UNDEFINED
var FILTER_RESET = C.PM_FILT_RESET
var FILTER_REALTIME = C.PM_FILT_REALTIME
var FILTER_NOTE = C.PM_FILT_NOTE
var FILTER_CHANNEL_AFTERTOUCH = C.PM_FILT_CHANNEL_AFTERTOUCH
var FILTER_POLY_AFTERTOUCH = C.PM_FILT_POLY_AFTERTOUCH
var FILTER_AFTERTOUCH = C.PM_FILT_AFTERTOUCH
var FILTER_PROGRAM = C.PM_FILT_PROGRAM
var FILTER_CONTROL = C.PM_FILT_CONTROL
var FILTER_PITCHBEND = C.PM_FILT_PITCHBEND
var FILTER_MTC = C.PM_FILT_MTC
var FILTER_SONG_POSITION = C.PM_FILT_SONG_POSITION
var FILTER_SONG_SELECT = C.PM_FILT_SONG_SELECT
var FILTER_TUNE = C.PM_FILT_TUNE
var FILTER_SYSTEMCOMMON = C.PM_FILT_SYSTEMCOMMON

func (s *Stream) SetFilter(filters int) error {
	filter := C.Pm_SetFilter(unsafe.Pointer(s.pmStream), C.int32_t(filters))
	if filter < 0 {
		return convertToError(C.PmError(filter))
	}
	return nil
}
