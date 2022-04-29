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

// #cgo linux LDFLAGS: -lporttime
// #cgo windows LDFLAGS: -lporttime
// #cgo darwin LDFLAGS: -lportmidi
//
// #include <stdlib.h>
// #include <porttime.h>
import "C"

func ptStart() {
	C.Pt_Start(C.int(1), nil, nil)
}

func ptStop() {
	C.Pt_Stop()
}

// Time returns the portmidi timer's current time.
func Time() Timestamp {
	return Timestamp(C.Pt_Time())
}
