package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
)

func TestKeyFromSharpsOrFlats(t *testing.T) {
	var tests = []struct {
		sharpsOrFlats int8
		mode          uint8
		key           byte
	}{
		{-4, 1, 5}, /* 4 flats minor */
		{-4, 0, 8}, /* 4 flats major */
		{3, 1, 6},  /* 3 sharps minor */
		{3, 0, 9},  /* 3 sharps major */
	}

	for _, test := range tests {
		key := KeyFromSharpsOrFlats(test.sharpsOrFlats, test.mode)

		if key != test.key {
			t.Errorf("KeyFromSharpsOrFlats(%v,%v) = %X; want %X", test.sharpsOrFlats, test.mode, key, test.key)
		}
	}
}

func TestParseStatus(t *testing.T) {
	var tests = []struct {
		byte           byte
		messageType    uint8
		messageChannel uint8
	}{
		{0xF0, 15, 0},  /* sysex */
		{0xF7, 15, 7},  /* sysex */
		{0xFF, 15, 15}, /* meta */
		{0xF8, 15, 8},  /* reatime MTC */
		{0xC0, 12, 0},  /* prog change chan 0 */
		{0x92, 9, 2},   /* note on chan2 */
		{0x81, 8, 1},   /* note off chan 1 */

	}

	for _, test := range tests {
		typ, ch := ParseStatus(test.byte)

		if typ != test.messageType || ch != test.messageChannel {
			t.Errorf("ParseStatus(%X) = %v,%v; want %v,%v", test.byte, typ, ch, test.messageType, test.messageChannel)
		}
	}

}

func TestParsePitchWheelVals(t *testing.T) {
	var tests = []struct {
		byte0    byte
		byte1    byte
		relative int16
		absolute uint16
	}{
		{0x00, 0x00, -8192, 0},
		{0xFF, 0xFF, 8191, 16383},
		{0xFF, 0xBF, -1, 8191},
		{0xF3, 0xF4, 6771, 14963},
	}

	for _, test := range tests {
		rel, abs := ParsePitchWheelVals(test.byte0, test.byte1)

		if rel != test.relative || abs != test.absolute {
			t.Errorf("ParsePitchWheelVals(%X, %X) = %v,%v; want %v,%v", test.byte0, test.byte1, rel, abs, test.relative, test.absolute)
		}
	}

}

func TestMsbLsbUnsigned(t *testing.T) {
	var tests = []struct {
		num      uint16
		descr    string
		expected string
	}{
		{0, "msbLsbUnsigned(0)", "0"},
		{8192, "msbLsbUnsigned(8192)", "1000000"},
		{16383, "msbLsbUnsigned(16383)", "111111101111111"},
	}

	for _, test := range tests {
		var b = MsbLsbUnsigned(test.num)

		if got, want := fmt.Sprintf("%b", b), test.expected; got != want {
			t.Errorf("%s = %s; want %s", test.descr, got, want)
		}
	}

}

func TestVLQ(t *testing.T) {
	var tests = []struct {
		num   uint32
		bytes []byte
	}{
		// tests according to MIDI SMF spec
		{0x40, []byte{0x40}},
		{0x7F, []byte{0x7F}},
		{0x80, []byte{0x81, 0x00}},
		{0x2000, []byte{0xC0, 0x00}},
		{0x3FFF, []byte{0xFF, 0x7F}},
		{0x4000, []byte{0x81, 0x80, 0x00}},
		{0x100000, []byte{0xC0, 0x80, 0x00}},
		{0x1FFFFF, []byte{0xFF, 0xFF, 0x7F}},
		{0x200000, []byte{0x81, 0x80, 0x80, 0x00}},
		{0x8000000, []byte{0xC0, 0x80, 0x80, 0x00}},
		{0xFFFFFFF, []byte{0xFF, 0xFF, 0xFF, 0x7F}},
	}

	for _, test := range tests {
		var b = VlqEncode(test.num)

		var bf bytes.Buffer
		binary.Write(&bf, binary.BigEndian, b)
		res, _ := ReadVarLength(&bf)

		if got, want := res, test.num; got != want {
			t.Errorf("ReadVarLength(%#v) = %d; want %d", test.bytes, got, want)
		}
	}

}

func TestLibBits(t *testing.T) {

	tests := []struct {
		result   interface{}
		descr    string
		expected string
	}{
		{clearBitU16(8191, 3), "clearBitU16(8192,2)", "1111111110111"},
		{clearBitU16(50, 1), "clearBitU16(50,1)", "110000"},
		{ClearBitU8(50, 1), "ClearBitU8(50, 1)", "110000"},
		{ClearBitU8(50, 4), "ClearBitU8(50, 4)", "100010"},
		{hasBitU8(50, 4), "hasBitU8(50, 4)", "%!b(bool=true)"},
		{hasBitU8(50, 3), "hasBitU8(50, 3)", "%!b(bool=false)"},
		{IsStatusByte(50), "IsStatusByte(50)", "%!b(bool=false)"},
		{IsStatusByte(128), "IsStatusByte(128)", "%!b(bool=true)"},
		{IsChannelMessage(128), "IsChannelMessage(128)", "%!b(bool=true)"},
		{IsChannelMessage(121), "IsChannelMessage(121)", "%!b(bool=false)"},
	}

	for _, test := range tests {

		if got, want := fmt.Sprintf("%b", test.result), test.expected; got != want {
			t.Errorf("%s = %#v; want %#v", test.descr, got, want)
		}
	}

}
