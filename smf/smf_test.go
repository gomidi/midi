package smf

import (
	"testing"
)

func TestUnpackTimeCode(t *testing.T) {

	tests := []struct {
		input     uint16
		fps       uint8
		subframes uint8
	}{
		{0xE808, 24, 8},
		{0xE728, 25, 40},
		{0xE264, 30, 100},
		{0xE350, 29, 80},
	}

	for _, test := range tests {
		fps, subframes := UnpackTimeCode(test.input)

		if got, want := fps, test.fps; got != want {
			t.Errorf("UnpackTimeCode(% X) [fps] = %v; want %v", test.input, got, want)
		}

		if got, want := subframes, test.subframes; got != want {
			t.Errorf("UnpackTimeCode(% X) [subframes] = %v; want %v", test.input, got, want)
		}
	}

}
