// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
// #include <porttime.h>
import "C"

import (
	"bytes"
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
	SysEx     []byte
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
		event.timestamp = C.PmTimestamp(evt.Timestamp)
		event.message = C.PmMessage((((evt.Data2 << 16) & 0xFF0000) | ((evt.Data1 << 8) & 0xFF00) | (evt.Status & 0xFF)))
		buffer[i] = event
	}
	return convertToError(C.Pm_Write(unsafe.Pointer(s.pmStream), &buffer[0], C.int32_t(size)))
}

// WriteShort writes a MIDI event of three bytes immediately to the output stream.
func (s *Stream) WriteShort(status int64, data1 int64, data2 int64) error {
	evt := Event{
		Timestamp: Timestamp(C.Pt_Time()),
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
	events = make([]Event, 0, numEvents)
	for i := 0; i < numEvents; i++ {
		event := Event{
			Timestamp: Timestamp(buffer[i].timestamp),
			Status:    int64(buffer[i].message) & 0xFF,
			Data1:     (int64(buffer[i].message) >> 8) & 0xFF,
			Data2:     (int64(buffer[i].message) >> 16) & 0xFF,
		}

		if event.Status&0xF0 == 0xF0 {
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

		events = append(events, event)
	}
	return
}

// ReadSysExBytes reads 4*max sysex bytes from the input stream.
//
// Deprecated. Using this API may cause a loss of buffered data.  It
// is preferable to use Read() and inspect the Event.SysEx field to
// detect SysEx messages.
func (s *Stream) ReadSysExBytes(max int) ([]byte, error) {
	evt, err := s.Read(max)
	if err != nil {
		return nil, err
	}
	return evt[0].SysEx, nil
}

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

// Poll reports whether there is input available in the stream.
func (s *Stream) Poll() (bool, error) {
	poll := C.Pm_Poll(unsafe.Pointer(s.pmStream))
	if poll < 0 {
		return false, convertToError(C.PmError(poll))
	}
	return poll > 0, nil
}

// TODO: add bindings for Pm_SetFilter
