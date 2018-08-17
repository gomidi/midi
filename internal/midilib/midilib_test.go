package midilib

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gomidi/midi/internal/vlq"
	"testing"
)

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
		var b = vlq.Encode(test.num)

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
