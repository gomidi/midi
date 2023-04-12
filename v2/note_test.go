package midi

import (
	"testing"
)

func TestNoteIs(t *testing.T) {
	var tests = []struct {
		noteA uint8
		noteB uint8
		equal bool
	}{
		{C(5), C(5), true},
		{C(5), C(8), true},
		{C(5), C(50), true},
		{C(5), D(5), false},
	}

	for i, test := range tests {

		got := Note(test.noteA).Is(Note(test.noteB))
		expected := test.equal

		if got != expected {
			t.Errorf("[%v] %s.Is(%s) = %v // expected: %v", i, Note(test.noteA).String(), Note(test.noteB).String(), got, expected)
		}
	}
}

func TestTranspose(t *testing.T) {
	var tests = []struct {
		note     uint8
		expected uint8
		interval Interval
	}{
		{C(5), Db(5), MinorSecond},
		{C(5), D(5), MajorSecond},
		{C(5), Eb(5), MinorThird},
		{C(5), E(5), MajorThird},
		{C(5), F(5), Fourth},
		{C(5), Gb(5), Tritone},
		{C(5), G(5), Fifth},
		{C(5), Ab(5), MinorSixth},
		{C(5), A(5), MajorSixth},
		{C(5), Bb(5), MinorSeventh},
		{C(5), B(5), MajorSeventh},
		{C(5), C(6), Octave},
		{C(5), Db(6), MinorNinth},
		{C(5), D(6), MajorNinth},
		{C(5), Eb(6), MinorTenth},
		{C(5), E(6), MajorTenth},
		{C(5), F(6), Eleventh},
		{C(5), Gb(6), DiminishedTwelfth},
		{C(5), G(6), Twelfth},
		{C(5), Ab(6), MinorThirteenth},
		{C(5), A(6), MajorThirteenth},
		{C(5), Bb(6), MinorFourteenth},
		{C(5), B(6), MajorFourteenth},
		{C(5), C(7), DoubleOctave},

		{C(5), B(4), -MinorSecond},
		{C(5), Bb(4), -MajorSecond},
		{C(5), A(4), -MinorThird},
		{C(5), Ab(4), -MajorThird},
		{C(5), G(4), -Fourth},
		{C(5), Gb(4), -Tritone},
		{C(5), F(4), -Fifth},
		{C(5), E(4), -MinorSixth},
		{C(5), Eb(4), -MajorSixth},
		{C(5), D(4), -MinorSeventh},
		{C(5), Db(4), -MajorSeventh},
		{C(5), C(4), -Octave},
		{C(5), B(3), -MinorNinth},
		{C(5), Bb(3), -MajorNinth},
		{C(5), A(3), -MinorTenth},
		{C(5), Ab(3), -MajorTenth},
		{C(5), G(3), -Eleventh},
		{C(5), Gb(3), -DiminishedTwelfth},
		{C(5), F(3), -Twelfth},
		{C(5), E(3), -MinorThirteenth},
		{C(5), Eb(3), -MajorThirteenth},
		{C(5), D(3), -MinorFourteenth},
		{C(5), Db(3), -MajorFourteenth},
		{C(5), C(3), -DoubleOctave},

		{G(5), D(6), Fifth},
		{A(5), E(6), Fifth},
		{B(5), Gb(6), Fifth},

		{C(5), D(2), -MinorSeventh - 24},
		{C(5), G(7), Fifth + 24},
	}

	for i, test := range tests {

		got := Note(Note(test.note).Transpose(test.interval)).String()
		expected := Note(test.expected).String()

		if got != expected {
			t.Errorf("[%v] %s.Transpose(%s) = %s // expected: %s", i, Note(test.note).String(), test.interval.String(), got, expected)
		}
	}
}

func TestInterval(t *testing.T) {
	var tests = []struct {
		noteA    uint8
		noteB    uint8
		expected Interval
	}{
		{C(5), Db(5), MinorSecond},
		{C(5), D(5), MajorSecond},
		{C(5), Eb(5), MinorThird},
		{C(5), E(5), MajorThird},
		{C(5), F(5), Fourth},
		{C(5), Gb(5), Tritone},
		{C(5), G(5), Fifth},
		{C(5), Ab(5), MinorSixth},
		{C(5), A(5), MajorSixth},
		{C(5), Bb(5), MinorSeventh},
		{C(5), B(5), MajorSeventh},
		{C(5), C(6), Octave},
		{C(5), Db(6), MinorNinth},
		{C(5), D(6), MajorNinth},
		{C(5), Eb(6), MinorTenth},
		{C(5), E(6), MajorTenth},
		{C(5), F(6), Eleventh},
		{C(5), Gb(6), DiminishedTwelfth},
		{C(5), G(6), Twelfth},
		{C(5), Ab(6), MinorThirteenth},
		{C(5), A(6), MajorThirteenth},
		{C(5), Bb(6), MinorFourteenth},
		{C(5), B(6), MajorFourteenth},
		{C(5), C(7), DoubleOctave},

		{C(5), B(4), -MinorSecond},
		{C(5), Bb(4), -MajorSecond},
		{C(5), A(4), -MinorThird},
		{C(5), Ab(4), -MajorThird},
		{C(5), G(4), -Fourth},
		{C(5), Gb(4), -Tritone},
		{C(5), F(4), -Fifth},
		{C(5), E(4), -MinorSixth},
		{C(5), Eb(4), -MajorSixth},
		{C(5), D(4), -MinorSeventh},
		{C(5), Db(4), -MajorSeventh},
		{C(5), C(4), -Octave},
		{C(5), B(3), -MinorNinth},
		{C(5), Bb(3), -MajorNinth},
		{C(5), A(3), -MinorTenth},
		{C(5), Ab(3), -MajorTenth},
		{C(5), G(3), -Eleventh},
		{C(5), Gb(3), -DiminishedTwelfth},
		{C(5), F(3), -Twelfth},
		{C(5), E(3), -MinorThirteenth},
		{C(5), Eb(3), -MajorThirteenth},
		{C(5), D(3), -MinorFourteenth},
		{C(5), Db(3), -MajorFourteenth},
		{C(5), C(3), -DoubleOctave},

		{G(5), D(6), Fifth},
		{A(5), E(6), Fifth},
		{B(5), Gb(6), Fifth},

		{C(5), D(2), -MinorSeventh},
		{C(5), G(7), Fifth},
	}

	for i, test := range tests {

		got := Note(test.noteA).Interval(Note(test.noteB)).String()
		expected := test.expected.String()

		if got != expected {
			t.Errorf("[%v] %s to %s = %s // expected: %s", i, Note(test.noteA).String(), Note(test.noteB).String(), got, expected)
		}
	}
}

func TestNote(t *testing.T) {

	var tests = []struct {
		note uint8
		str  string
	}{
		{C(0), "C0"},
		{C(1), "C1"},
		{C(2), "C2"},
		{C(3), "C3"},
		{C(4), "C4"},
		{C(5), "C5"},
		{C(6), "C6"},
		{C(7), "C7"},
		{C(8), "C8"},
		{C(9), "C9"},
		{C(10), "C10"},

		{Db(0), "Db0"},
		{Db(1), "Db1"},
		{Db(2), "Db2"},
		{Db(3), "Db3"},
		{Db(4), "Db4"},
		{Db(5), "Db5"},
		{Db(6), "Db6"},
		{Db(7), "Db7"},
		{Db(8), "Db8"},
		{Db(9), "Db9"},
		{Db(10), "Db10"},

		{D(0), "D0"},
		{D(1), "D1"},
		{D(2), "D2"},
		{D(3), "D3"},
		{D(4), "D4"},
		{D(5), "D5"},
		{D(6), "D6"},
		{D(7), "D7"},
		{D(8), "D8"},
		{D(9), "D9"},
		{D(10), "D10"},

		{Eb(0), "Eb0"},
		{Eb(1), "Eb1"},
		{Eb(2), "Eb2"},
		{Eb(3), "Eb3"},
		{Eb(4), "Eb4"},
		{Eb(5), "Eb5"},
		{Eb(6), "Eb6"},
		{Eb(7), "Eb7"},
		{Eb(8), "Eb8"},
		{Eb(9), "Eb9"},
		{Eb(10), "Eb10"},

		{G(0), "G0"},
		{G(1), "G1"},
		{G(2), "G2"},
		{G(3), "G3"},
		{G(4), "G4"},
		{G(5), "G5"},
		{G(6), "G6"},
		{G(7), "G7"},
		{G(8), "G8"},
		{G(9), "G9"},
		{G(10), "G10"},

		{Ab(0), "Ab0"},
		{Ab(1), "Ab1"},
		{Ab(2), "Ab2"},
		{Ab(3), "Ab3"},
		{Ab(4), "Ab4"},
		{Ab(5), "Ab5"},
		{Ab(6), "Ab6"},
		{Ab(7), "Ab7"},
		{Ab(8), "Ab8"},
		{Ab(9), "Ab9"},
		{Ab(10), "Ab9"},

		{B(0), "B0"},
		{B(1), "B1"},
		{B(2), "B2"},
		{B(3), "B3"},
		{B(4), "B4"},
		{B(5), "B5"},
		{B(6), "B6"},
		{B(7), "B7"},
		{B(8), "B8"},
		{B(9), "B9"},
		{B(10), "B9"},
	}

	for i, test := range tests {
		//fmt.Printf("%v\n", test.note)
		if test.note > 127 {
			t.Errorf("note in test %v is too large: %v", i, test.note)
		}

		res := Note(test.note).String()
		exp := test.str

		if res != exp {
			t.Errorf("expected: %q, but got %q", exp, res)
		}
	}
}
