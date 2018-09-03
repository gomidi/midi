package midilib

/*
functions in this file are _not_ derived from the work of Joe Wass.
*/

import (
	"io"
)

func clearBitU16(n uint16, pos uint16) uint16 {
	mask := ^(uint16(1) << pos)
	n &= mask
	return n
}

// ClearBitU8 clears the bit at position pos within n
func ClearBitU8(n uint8, pos uint8) uint8 {
	mask := ^(uint8(1) << pos)
	n &= mask
	return n
}

// MsbLsbSigned returns the uint16 for a signed MSB LSB message combination
func MsbLsbSigned(n int16) uint16 {

	//		if n > 8191 {
	//			panic("n must not overflow 14bits (max 8191)")
	//		}
	//		if n < -8191 {
	//			panic("n must not overflow 14bits (min -8191)")
	//		}

	return msbLsbUnsigned(uint16(n + 8192))
}

// takes a 14bit uint and pads it to 16 bit like in the specs for e.g. pitchbend
func msbLsbUnsigned(n uint16) uint16 {
	if n > 16383 {
		panic("n must not overflow 14bits (max 16383)")
	}

	lsb := n << 8
	lsb = clearBitU16(lsb, 15)
	lsb = clearBitU16(lsb, 7)

	// 0x7f = 127 = 0000000001111111
	msb := 0x7f & (n >> 7)
	return lsb | msb
}

func hasBitU8(n uint8, pos uint8) bool {
	val := n & (1 << pos)
	return (val > 0)
}

// IsChannelMessage returns if the given byte is a channel message
func IsChannelMessage(b uint8) bool {
	return !hasBitU8(b, 6)
}

// IsStatusByte returns if the given byte is a status byte
func IsStatusByte(b uint8) bool {
	return hasBitU8(b, 7)
}

// ReadNBytes reads n bytes from the reader
func ReadNBytes(n int, rd io.Reader) ([]byte, error) {
	var b []byte = make([]byte, n)
	num, err := rd.Read(b)

	// if num is correct, we are not interested in io.EOF errors
	if num == n {
		err = nil
	}

	return b, err
}

// ReadByte reads a byte from the reader
func ReadByte(rd io.Reader) (byte, error) {
	b, err := ReadNBytes(1, rd)

	if err != nil {
		return 0, err
	}

	return b[0], nil
}
