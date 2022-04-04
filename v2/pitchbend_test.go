package midi

import (
	"testing"
)

func TestPitchbend(t *testing.T) {

	tests := []struct {
		in       int16
		expected uint16
	}{
		{
			in:       0,
			expected: 8192,
		},
		{
			in:       PitchHighest,
			expected: 16383,
		},
		{
			in:       PitchHighest + 1,
			expected: 16383,
		},
		{
			in:       PitchLowest,
			expected: 0,
		},
		{
			in:       PitchLowest - 1,
			expected: 0,
		},
	}

	for _, test := range tests {
		m := Pitchbend(0, test.in)

		var got uint16
		var ch uint8
		Message(m).GetPitchBend(&ch, nil, &got)

		if got != test.expected {
			t.Errorf("Pitchbend(%v).absValue = %v; wanted %v", test.in, got, test.expected)
		}
	}
}
