package midilib

import (
	"github.com/gomidi/midi"
	"io"
)

func ReadNBytes(n int, rd io.Reader) ([]byte, error) {
	var b []byte = make([]byte, n)
	num, err := rd.Read(b)

	if err != nil {
		return nil, err
	}
	if num != n {
		return nil, midi.ErrUnexpectedEOF
	}

	return b, nil
}

func ReadByte(rd io.Reader) (byte, error) {
	b, err := ReadNBytes(1, rd)

	if err != nil {
		return 0, err
	}

	return b[0], nil
}

// ReadUint16 reads a 2-byte 16 bit integer from a ReadSeeker.
// It returns the 16-bit value and an error.
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

// ReadUint24 parse a 3-byte 24 bit integer from a ReadSeeker.
// It returns the 32-bit value and an error.
// TODO TEST
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

// ReadUint32 parse a 4-byte 32 bit integer from a ReadSeeker.
// It returns the 32-bit value and an error.
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
		return result, midi.ErrUnexpectedEOF
	}

	return result, nil
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
		return []byte{}, midi.ErrUnexpectedEOF
	}

	// If there was some other problem, that's also a problem.
	if err != nil {
		return []byte{}, err
	}

	return buffer, nil
}
