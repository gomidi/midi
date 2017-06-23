package lib

import (
	"bytes"
	"testing"
)

func TestVarLength(t *testing.T) {

	tests := []struct {
		input    []byte
		expected uint32
	}{
		{[]byte{0x05}, 5},
	}

	// F0 05 43 12 00 07 F7

	for _, test := range tests {
		res, err := ReadVarLength(bytes.NewReader(test.input))

		//res := VlqDecode(test.input)

		if err != nil {
			t.Fatalf("ReadVarLength(% X) ; error: %v", test.input, err)
		}

		if got, want := res, test.expected; got != want {
			t.Errorf("ReadVarLength(% X) = %v; want %v", test.input, got, want)
		}
	}

}
