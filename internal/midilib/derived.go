package midilib

import (
	"io"

	"github.com/gomidi/midi"
)

/*
This file contains functions that are modifications of the functions found
in the github.com/afandian/go-midi package of Joe Wass.

See the file midi_functions.go for the original functions.
*/

// ReadUint16 reads a 2-byte 16 bit integer from a Reader.
// It returns the 16-bit value and an error.
// This is a slightly modified variant of the parseUint16 function
// from Joe Wass. See the file midi_functions.go for the original.
func ReadUint16(rd io.Reader) (uint16, error) {
	b, err := ReadNBytes(2, rd)

	if err != nil {
		return 0, err
	}

	var val uint16 = 0x00
	val |= uint16(b[1]) << 0
	val |= uint16(b[0]) << 8

	return val, nil
}

// ReadUint24 parse a 3-byte 24 bit integer from a Reader.
// It returns the 32-bit value and an error.
// This is a slightly modified variant of the parseUint24 function
// from Joe Wass. See the file midi_functions.go for the original.
func ReadUint24(rd io.Reader) (uint32, error) {
	b, err := ReadNBytes(3, rd)

	if err != nil {
		return 0, err
	}

	var val uint32 = 0x00
	val |= uint32(b[2]) << 0
	val |= uint32(b[1]) << 8
	val |= uint32(b[0]) << 16

	return val, nil
}

// ReadUint32 parse a 4-byte 32 bit integer from a Reader.
// It returns the 32-bit value and an error.
// This is a slightly modified variant of the parseUint32 function
// from Joe Wass. See the file midi_functions.go for the original.
func ReadUint32(rd io.Reader) (uint32, error) {
	b, err := ReadNBytes(4, rd)

	if err != nil {
		return 0, err
	}

	var val uint32 = 0x00
	val |= uint32(b[3]) << 0
	val |= uint32(b[2]) << 8
	val |= uint32(b[1]) << 16
	val |= uint32(b[0]) << 24

	return val, nil
}

// ReadVarLength reads a variable length value from a Reader.
// It returns the [up to] 32-bit value and an error.
// This is a slightly modified variant of the parseVarLength function
// from Joe Wass. See the file midi_functions.go for the original.
func ReadVarLength(reader io.Reader) (uint32, error) {

	// Single byte buffer to read byte by byte.
	var buffer []byte = make([]uint8, 1)

	// The number of bytes returned.
	// Should always be 1 unless we reach the EOF
	var num int = 1

	// Result value
	var result uint32 = 0x00

	// RTFM.
	var first = true
	for (first || (buffer[0]&0x80 == 0x80)) && (num > 0) {
		result = result << 7

		num, _ = reader.Read(buffer)
		result |= (uint32(buffer[0]) & 0x7f)
		first = false
	}

	if num == 0 && !first {
		return result, midi.ErrUnexpectedEOF
	}

	return result, nil
}

// ReadVarLengthData reads data that is prefixed by a varLength that tells the length of the data
//
// This is a slightly modified variant of the parseText function
// from Joe Wass. See the file midi_functions.go for the original.
func ReadVarLengthData(reader io.Reader) ([]byte, error) {
	length, err := ReadVarLength(reader)

	if err != nil {
		return []byte{}, err
	}

	var buffer []byte = make([]byte, length)

	num, err := reader.Read(buffer)

	// If we couldn't read the entire expected-length buffer, that's a problem.
	if num != int(length) {
		return []byte{}, midi.ErrUnexpectedEOF
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return []byte{}, err
	}

	return buffer, nil
}

// ParseStatus parses the status byte and returns type and channel
//
// This is a slightly modified variant of the readStatusByte function
// from Joe Wass. See the file midi_functions.go for the original.
func ParseStatus(b byte) (messageType uint8, messageChannel uint8) {
	messageType = (b & 0xF0) >> 4
	messageChannel = b & 0x0F
	return
}

// -----------------------------------------

// ParseUint7 parses a 7-bit bit integer from a byte, ignoring the high bit.
//
// This is a slightly modified variant of the parseUint7 function
// from Joe Wass. See the file midi_functions.go for the original.
func ParseUint7(b byte) uint8 {
	return b & 0x7f
}

// ParseTwoUint7 parses two 7-bit bit integer stored in two bytes, ignoring the high bit in each.
//
// This is a slightly modified variant of the parseTwoUint7 function
// from Joe Wass. See the file midi_functions.go for the original.
func ParseTwoUint7(b1, b2 byte) (uint8, uint8) {
	return (b1 & 0x7f), (b2 & 0x7f)
}

// ParsePitchWheelVals parses a 14-bit signed value, which becomes a signed int16.
//
// This is a slightly modified variant of the parsePitchWheelValue function
// from Joe Wass. See the file midi_functions.go for the original.
func ParsePitchWheelVals(b1 byte, b2 byte) (relative int16, absolute uint16) {
	var val uint16

	val = uint16((b2)&0x7f) << 7
	val |= uint16(b1) & 0x7f

	// Turn into a signed value relative to the centre.
	relative = int16(val) - 0x2000

	return relative, val
}

// Supplied to KeySignature
const (
	majorMode = 0
	minorMode = 1
)

// KeyFromSharpsOrFlats Taking a signed number of sharps or flats (positive for sharps, negative for flats) and a mode (0 for major, 1 for minor)
// decide the key signature.
//
// This is a slightly modified variant of the keySignatureFromSharpsOrFlats function
// from Joe Wass. See the file music.go for the original.
func KeyFromSharpsOrFlats(sharpsOrFlats int8, mode uint8) uint8 {
	tmp := int(sharpsOrFlats * 7)

	// Relative Minor.
	if mode == minorMode {
		tmp -= 3
	}

	// Clamp to Octave 0-11.
	for tmp < 0 {
		tmp += 12
	}

	return uint8(tmp % 12)
}
