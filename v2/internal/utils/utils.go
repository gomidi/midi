package utils

import (
	"fmt"
	"io"
)

const (
	vlqContinue = 128
	vlqMask     = 127
)

func hasBitU8(n uint8, pos uint8) bool {
	val := n & (1 << pos)
	return (val > 0)
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

// IsChannelMessage returns if the given byte is a channel message
func IsChannelMessage(b uint8) bool {
	return !hasBitU8(b, 6)
}

// IsStatusByte returns if the given byte is a status byte
func IsStatusByte(b uint8) bool {
	return hasBitU8(b, 7)
}

/*
// ReadByte reads a byte from the reader
func ReadByte(rd io.Reader) (byte, error) {
	b, err := ReadNBytes(1, rd)

	if err != nil {
		return 0, err
	}

	return b[0], nil
}
*/

// ParseStatus parses the status byte and returns type and channel
//
// This is a slightly modified variant of the readStatusByte function
// from Joe Wass. See the file midi_functions.go for the original.
func ParseStatus(b byte) (messageType uint8, messageChannel uint8) {
	messageType = (b & 0xF0) >> 4
	messageChannel = b & 0x0F
	return
}

// ReadByte reads a byte from the reader
func ReadByte(rd io.Reader) (byte, error) {
	b, err := ReadNBytes(1, rd)

	if err != nil {
		return 0, err
	}

	return b[0], nil
}

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

// ParseUint16 converts 2 bytes to a 16 bit integer
// This is a slightly modified variant of the parseUint16 function
// from Joe Wass. See the file midi_functions.go for the original.
func ParseUint16(b1, b2 byte) uint16 {

	var val uint16 = 0x00
	val |= uint16(b2) << 0
	val |= uint16(b1) << 8

	return val
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

// ErrUnexpectedEOF is returned, when an unexspected end of file is reached.
var ErrUnexpectedEOF = fmt.Errorf("Unexpected End of File found.")

func ReadText(rd io.Reader) (string, error) {
	b, err := ReadVarLengthData(rd)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

// limit the largest possible value to int32
/*
The largest number which is allowed is 0FFFFFFF so that the variable-length representations must fit in 32
bits in a routine to write variable-length numbers. Theoretically, larger numbers are possible, but 2 x 10 8
96ths of a beat at a fast tempo of 500 beats per minute is four days, long enough for any delta-time!
*/

// Variable-Length Quantity (VLQ) is an way of representing arbitrary
// see https://blogs.infosupport.com/a-primer-on-vlq/
// we use the variant of the midi-spec
// stolen and converted to go from https://github.com/dvberkel/VLQKata/blob/master/src/main/java/nl/dvberkel/kata/Kata.java#L12

// Encode encodes the given value as variable length quantity
func VlqEncode(n uint32) (out []byte) {
	var quo, rem uint32
	quo = n / vlqContinue
	rem = n % vlqContinue

	out = append(out, byte(rem))

	for quo > 0 {
		out = append(out, byte(quo)|vlqContinue)
		quo = quo / vlqContinue
		// rem = quo % vlqContinue
	}

	reverse(out)
	return
}

// stolen from http://stackoverflow.com/questions/19239449/how-do-i-reverse-an-array-in-go
func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

// Decode decodes a variable length quantity
func VlqDecode(source []byte) (num uint32) {

	for i := 0; i < len(source); i++ {
		var n = uint32(source[i] & vlqMask)
		for (source[i] & vlqContinue) != 0 {
			i++
			n *= 128
			n += uint32(source[i] & vlqMask)
		}
		num += n
	}

	return
}

// ClearBitU8 clears the bit at position pos within n
func ClearBitU8(n uint8, pos uint8) uint8 {
	mask := ^(uint8(1) << pos)
	n &= mask
	return n
}

func ParseTwoUint7(b1, b2 byte) (uint8, uint8) {
	return (b1 & 0x7f), (b2 & 0x7f)
}

func ParseUint7(b byte) uint8 {
	return b & 0x7f
}

func ParsePitchWheelVals(b1 byte, b2 byte) (relative int16, absolute uint16) {
	var val uint16

	val = uint16((b2)&0x7f) << 7
	val |= uint16(b1) & 0x7f

	// Turn into a signed value relative to the centre.
	relative = int16(val) - 0x2000

	return relative, val
}

// MsbLsbSigned returns the uint16 for a signed MSB LSB message combination
func MsbLsbSigned(n int16) uint16 {

	//		if n > 8191 {
	//			panic("n must not overflow 14bits (max 8191)")
	//		}
	//		if n < -8191 {
	//			panic("n must not overflow 14bits (min -8191)")
	//		}

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

func clearBitU16(n uint16, pos uint16) uint16 {
	mask := ^(uint16(1) << pos)
	n &= mask
	return n
}

/*
func ParseStatus(b byte) (messageType uint8, messageChannel uint8) {
	messageType = (b & 0xF0) >> 4
	messageChannel = b & 0x0F
	return
}
*/

/*
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
*/

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
		return []byte{}, ErrUnexpectedEOF
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return []byte{}, err
	}

	return buffer, nil
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
		return result, ErrUnexpectedEOF
	}

	return result, nil
}
