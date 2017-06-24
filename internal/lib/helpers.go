/*
functions in this  file are taken from github.com/afandian/go-midi (Copyright Joe Wass)
some of them are slightly modified
*/

// Copyright 2012 Joe Wass. All rights reserved.
// Use of this source code is governed by the MIT license
// which can be found in the LICENSE file.

// MIDI package
// A package for reading Standard Midi Files, written in Go.
// Joe Wass 2012
// joe@afandian.com

/*
 * Functions for reading actual MIDI data in the various formats that crop up.
 */

package lib

import (
	// "fmt"
	"fmt"
	"io"
	"math/big"
	"strconv"
)

// Variable-Length Quantity (VLQ) is an way of representing arbitrarly
// see https://blogs.infosupport.com/a-primer-on-vlq/
// we use the variant of the midi-spec
// stolen and converted to go from https://github.com/dvberkel/VLQKata/blob/master/src/main/java/nl/dvberkel/kata/Kata.java#L12

const (
	continueBit = 128
	mask        = 127
)

var (
	baseBig = big.NewInt(128)
)

func encodeVLQ(n *big.Int) (out []byte) {
	var resNo, resRemainder big.Int

	resNo.QuoRem(n, baseBig, &resRemainder)
	out = append(out, leastSignificantBits(&resRemainder))

	for resNo.Cmp(big.NewInt(0)) == 1 {
		out = append(out, leastSignificantBits(&resNo)|continueBit)
		resNo.QuoRem(&resNo, baseBig, &resRemainder)
	}

	reverse(out)
	return
}

func decodeVLQ(source []byte) *big.Int {
	var res big.Int

	for i := 0; i < len(source); i++ {
		var n = valueOf(source[i] & mask)
		for (source[i] & continueBit) != 0 {
			i++
			n = n.Mul(n, baseBig)
			n.Add(n, valueOf(source[i]&mask))
		}
		res.Add(&res, n)
	}

	return &res
}

// leastSignificantBits returns the least significant bits (a byte)
// from the bit representation of n
func leastSignificantBits(n *big.Int) byte {
	b := toBytes(n)
	if len(b) == 0 {
		return 0
	}
	return b[0]
}

// stolen from http://stackoverflow.com/questions/19239449/how-do-i-reverse-an-array-in-go
func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

// toBytes converts n to the byte slice representing n in binary form
func toBytes(n *big.Int) []byte {
	var resRemainder, x big.Int
	x.QuoRem(n, baseBig, &resRemainder)
	return resRemainder.Bytes()
}

// valueOf converts a byte in binary form representing an int string into an int
func valueOf(b byte) *big.Int {
	var i big.Int
	i.SetBytes([]byte{b})
	return &i
}

func ReadN(n int, rd io.Reader) ([]byte, error) {
	var b []byte = make([]byte, n)
	num, err := rd.Read(b)

	if err != nil {
		return nil, err
	}
	if num != n {
		return nil, ErrUnexpectedEOF
	}

	return b, nil
}

// readUint32 parse a 4-byte 32 bit integer from a ReadSeeker.
// It returns the 32-bit value and an error.
func ReadUint32(rd io.Reader) (uint32, error) {
	b, err := ReadN(4, rd)

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

// readUint24 parse a 3-byte 24 bit integer from a ReadSeeker.
// It returns the 32-bit value and an error.
// TODO TEST
func ReadUint24(rd io.Reader) (uint32, error) {
	b, err := ReadN(3, rd)

	if err != nil {
		return 0, err
	}

	var val uint32 = 0x00
	val |= uint32(b[2]) << 0
	val |= uint32(b[1]) << 8
	val |= uint32(b[0]) << 16

	return val, nil
}

// readUint16 reads a 2-byte 16 bit integer from a ReadSeeker.
// It returns the 16-bit value and an error.
func ReadUint16(rd io.Reader) (uint16, error) {
	b, err := ReadN(2, rd)

	if err != nil {
		return 0, err
	}

	var val uint16 = 0x00
	val |= uint16(b[1]) << 0
	val |= uint16(b[0]) << 8

	return val, nil
}

// parseUint7 parses a 7-bit bit integer from a byte, ignoring the high bit.
func ParseUint7(b byte) uint8 {
	return b & 0x7f
}

// parseTwoUint7 parses two 7-bit bit integer stored in two bytes, ignoring the high bit in each.
func ParseTwoUint7(b1, b2 byte) (uint8, uint8) {
	return (b1 & 0x7f), (b2 & 0x7f)
}

func ParsePitchWheelVals(b1 byte, b2 byte) (relative int16, absolute uint16) {
	var val uint16 = 0

	val = uint16((b2)&0x7f) << 7
	val |= uint16(b1) & 0x7f

	// Turn into a signed value relative to the centre.
	relative = int16(val) - 0x2000

	return relative, val
}

func VlqDecode(b []byte) uint32 {
	return uint32(decodeVLQ(b).Int64())
}

/*
func decodeVarLength(rd io.Reader) (uint32, error) {
	var b uint8
	var bf []byte
	err := binary.Read(rd, binary.BigEndian, &b)


	// every but the last byte has the bit 7 set
	for err == nil && hasBitU8(b, 7) {
		bf = append(bf, b)
		err = binary.Read(rd, binary.BigEndian, &b)
	}

	if err != nil {
		return 0, err
	}

	bf = append(bf, b)

	return vlqDecode(bf), nil
}
*/

// RoundFloat rounds the given float by the given decimals after the dot
func roundFloat(x float64, decimals int) float64 {
	// return roundFloat(x, numDig(x)+decimals)
	frep := strconv.FormatFloat(x, 'f', decimals, 64)
	f, _ := strconv.ParseFloat(frep, 64)
	return f
}

func floatToInt(x float64) int {
	return int(roundFloat(x, 0))
}

// takes a 14bit int and pads it to 16 bit like in the specs for e.g. pitchbend
func MsbLsbSigned(n int16) uint16 {
	if n > 8191 {
		panic("n must not overflow 14bits (max 8191)")
	}
	if n < -8191 {
		panic("n must not overflow 14bits (min -8191)")
	}
	return MsbLsbUnsigned(uint16(n + 8192))
}

// takes a 14bit uint and pads it to 16 bit like in the specs for e.g. pitchbend
func MsbLsbUnsigned(n uint16) uint16 {
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

// limit the largest possible value to int32
/*
The largest number which is allowed is 0FFFFFFF so that the variable-length representations must fit in 32
bits in a routine to write variable-length numbers. Theoretically, larger numbers are possible, but 2 x 10 8
96ths of a beat at a fast tempo of 500 beats per minute is four days, long enough for any delta-time!
*/
func VlqEncode(i uint32) []byte {
	return encodeVLQ(big.NewInt(int64(i)))
}

var _ = fmt.Sprintf

// helpers from http://stackoverflow.com/questions/23192262/how-would-you-set-and-clear-a-single-bit-in-go

/*
Here's a function to set a bit. First, shift the number 1 the specified number of spaces
in the integer (so it becomes 0010, 0100, etc). Then OR it with the original input. This
leaves the other bits unaffected but will always set the target bit to 1.
*/
// Sets the bit at pos in the integer n.
func setBit(n int, pos uint) int {
	n |= (1 << pos)
	return n
}

/*
Here's a function to clear a bit. First shift the number 1 the specified number of spaces
in the integer (so it becomes 0010, 0100, etc). Then flip every bit in the mask with
the ^ operator (so 0010 becomes 1101). Then use a bitwise AND, which doesn't touch the
numbers AND'ed with 1, but which will unset the value in the mask which is set to 0.
*/
// Clears the bit at pos in n.
func clearBit(n int, pos uint) int {
	mask := ^(1 << pos)
	n &= mask
	return n
}

/*
Finally here's a function to check whether a bit is set. Shift the number 1 the specified
number of spaces (so it becomes 0010, 0100, etc) and then AND it with the target number.
If the resulting number is greater than 0 (it'll be 1, 2, 4, 8, etc) then the bit is set.
*/
func hasBit(n int, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}

/*
try to adopt them for int16
*/
func setBit16(n int16, pos uint16) int16 {
	n |= (1 << pos)
	return n
}

func clearBit16(n int16, pos uint16) int16 {
	mask := ^(int16(1) << pos)
	n &= mask
	return n
}

func hasBit16(n int16, pos uint16) bool {
	val := n & (1 << pos)
	return (val > 0)
}

/*
try to adopt them for int32
*/
func setBit32(n int32, pos uint32) int32 {
	n |= (1 << pos)
	return n
}

func clearBit32(n int32, pos uint32) int32 {
	mask := ^(int32(1) << pos)
	n &= mask
	return n
}

func hasBit32(n int32, pos uint32) bool {
	val := n & (1 << pos)
	return (val > 0)
}

/*
try to adopt them for uint16
*/
func setBitU16(n uint16, pos uint16) uint16 {
	n |= (1 << pos)
	return n
}

func clearBitU16(n uint16, pos uint16) uint16 {
	mask := ^(uint16(1) << pos)
	n &= mask
	return n
}

func hasBitU16(n uint16, pos uint16) bool {
	val := n & (1 << pos)
	return (val > 0)
}

/*
try to adopt them for uint8
*/
func setBitU8(n uint8, pos uint8) uint8 {
	n |= (1 << pos)
	return n
}

func ClearBitU8(n uint8, pos uint8) uint8 {
	mask := ^(uint8(1) << pos)
	n &= mask
	return n
}

func hasBitU8(n uint8, pos uint8) bool {
	val := n & (1 << pos)
	return (val > 0)
}

/*
func i32To3bytes(i int32) [3]byte {
	bt := [3]byte{}
	bt[0] = i & 0xff000000 >> 24
	bt[1] = i & 0x00ff0000 >> 16
	bt[2] = i & 0x0000ff00 >> 8
	return bt
}
*/

/*
CP1 = (CurrentPosition & 0xff000000UL) >> 24;
    CP2 = (CurrentPosition & 0x00ff0000UL) >> 16;
    CP3 = (CurrentPosition & 0x0000ff00UL) >>  8;
    CP4 = (CurrentPosition & 0x000000ffUL)      ;
*/

// readVarLength reads a variable length value from a ReadSeeker.
// It returns the [up to] 32-bit value and an error.
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
		return result, ErrUnexpectedEOF
	}

	return result, nil
}

// readByte reads one byte from the reader and returns it
func ReadByte(rd io.Reader) (byte, error) {
	b, err := ReadN(1, rd)

	if err != nil {
		return 0, err
	}

	return b[0], nil
}

// read2Bytes reads two bytes from the reader and returns them
func read2Bytes(rd io.Reader) (b1, b2 byte, err error) {
	b, err := ReadN(2, rd)

	if err != nil {
		return 0, 0, err
	}

	return b[0], b[1], nil
}

func ParseStatus(b byte) (messageType uint8, messageChannel uint8) {
	messageType = (b & 0xF0) >> 4
	messageChannel = b & 0x0F
	return
}

// readVarLengthData reads data that is prefixed by a varLength that tells the length of the data
func ReadVarLengthData(reader io.Reader) ([]byte, error) {
	length, err := ReadVarLength(reader)

	if err != nil {
		return []byte{}, err
	}

	var buffer []byte = make([]byte, length)

	num, err := reader.Read(buffer)

	// If we couldn't read the entire expected-length buffer, that's a problem.
	if num != int(length) {
		return []byte{}, ErrUnexpectedEOF
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return []byte{}, err
	}

	return buffer, nil
}

func ReadText(rd io.Reader) (string, error) {
	b, err := ReadVarLengthData(rd)

	if err != nil {
		return "", err
	}

	// TODO: Data should be ASCII but might go up to 0xFF.
	// What will Go do? Try and decode UTF-8?
	return string(b), nil
}

/*
func hasBitU8(n uint8, pos uint8) bool {
	val := n & (1 << pos)
	return (val > 0)
}
*/
