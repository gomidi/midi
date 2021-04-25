package midi

import (
	"testing"
)

func TestPitchbend(t *testing.T) {

	tests := []struct {
		in       int16
		expected int16
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
		m := NewMessage(Channel(0).Pitchbend(test.in))

		_, abs := m.Pitch()

		got := abs

		if got != test.expected {
			t.Errorf("Pitchbend(%v).absValue = %v; wanted %v", test.in, got, test.expected)
		}
	}
}
