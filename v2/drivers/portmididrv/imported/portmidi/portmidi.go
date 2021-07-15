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

// Package portmidi provides PortMidi bindings.
package portmidi

// #cgo CFLAGS:  -I/usr/local/include
// #cgo LDFLAGS: -lportmidi -L/usr/local/lib
//
// #include <stdlib.h>
// #include <portmidi.h>
// #include <porttime.h>
import "C"

import (
	"errors"
)

// DeviceID is a MIDI device ID.
type DeviceID int

// DeviceInfo provides info about a MIDI device.
type DeviceInfo struct {
	Interface         string
	Name              string
	IsInputAvailable  bool
	IsOutputAvailable bool
	IsOpened          bool
}

type Timestamp int64

// Initialize initializes the portmidi. Needs to be called before
// making any other call from the portmidi package.
// Once portmidi package is no longer required, Terminate should be
// called to free the underlying resources.
func Initialize() error {
	if code := C.Pm_Initialize(); code != 0 {
		return convertToError(code)
	}
	C.Pt_Start(C.int(1), nil, nil)
	return nil
}

// Terminate terminates and cleans up the midi streams.
func Terminate() error {
	C.Pt_Stop()
	return convertToError(C.Pm_Terminate())
}

// DefaultInputDeviceID returns the default input device's ID.
func DefaultInputDeviceID() DeviceID {
	return DeviceID(C.Pm_GetDefaultInputDeviceID())
}

// DefaultOutputDeviceID returns the default output device's ID.
func DefaultOutputDeviceID() DeviceID {
	return DeviceID(C.Pm_GetDefaultOutputDeviceID())
}

// CountDevices returns the number of MIDI devices.
func CountDevices() int {
	return int(C.Pm_CountDevices())
}

// Info returns the device info for the device indentified with deviceID.
// If deviceID is out of range, Info returns nil.
func Info(deviceID DeviceID) *DeviceInfo {
	info := C.Pm_GetDeviceInfo(C.PmDeviceID(deviceID))
	if info == nil {
		return nil
	}
	return &DeviceInfo{
		Interface:         C.GoString(info.interf),
		Name:              C.GoString(info.name),
		IsInputAvailable:  info.input > 0,
		IsOutputAvailable: info.output > 0,
		IsOpened:          info.opened > 0,
	}
}

// Time returns the portmidi timer's current time.
func Time() Timestamp {
	return Timestamp(C.Pt_Time())
}

// convertToError converts a portmidi error code to a Go error.
func convertToError(code C.PmError) error {
	if code >= 0 {
		return nil
	}
	return errors.New(C.GoString(C.Pm_GetErrorText(code)))
}
