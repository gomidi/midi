package utils

import (
	"fmt"
	"testing"
)

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

func TestVlqEncode(t *testing.T) {
	for _, test := range tests {
		var b = VlqEncode(test.num)

		if got, want := fmt.Sprintf("%X", b), fmt.Sprintf("%X", test.bytes); got != want {
			t.Errorf("Encode(%#v) = %#v; want %#v", test.num, got, want)
		}
	}

}

func TestVlqDecode(t *testing.T) {
	for _, test := range tests {
		var b = VlqDecode(test.bytes)

		if got, want := b, test.num; got != want {
			t.Errorf("Decode(%#v) = %d; want %d", test.bytes, got, want)
		}
	}

}
