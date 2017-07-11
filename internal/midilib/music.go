// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

/*
 * Functions that deal with musical concepts.
 */

package midilib

/*
// import "fmt"

// Taking a signed number of sharps or flats (positive for sharps, negative for flats) and a mode (0 for major, 1 for minor)
// decide the key signature.
func keySignatureFromSharpsOrFlats(sharpsOrFlats int8, mode uint8) (key ScaleDegree, resultMode KeySignatureMode) {
	// 0 is C.
	var tmp int = int(DegreeC + sharpsOrFlats*7)

	// Relative Minor.
	if mode == MinorMode {
		tmp -= 3
	}

	// Clamp to Octave 0-11.
	for tmp < 0 {
		tmp += 12
	}

	tmp = tmp % 12

	resultMode = KeySignatureMode(mode)
	key = ScaleDegree(tmp)

	return
}
*/
