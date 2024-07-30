package version

import (
	"sort"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		error   bool
		version *Version
	}{
		{ // 0
			"v1.0.1",
			false,
			&Version{1, 0, 1},
		},
		{ // 1
			"x1.0.1",
			true,
			nil,
		},
		{ // 2
			"1.0.1",
			false,
			&Version{1, 0, 1},
		},
		{ // 3
			"1.0.1.3",
			true,
			nil,
		},
		{ // 4
			"-1.0.1",
			true,
			nil,
		},
		{ // 5
			"1.0",
			false,
			&Version{1, 0, 0},
		},
		{ // 6
			"1.0.0",
			false,
			&Version{1, 0, 0},
		},
		{ // 7
			"0.0.0",
			true,
			nil,
		},
		{ // 8
			"0.0",
			true,
			nil,
		},
		{ // 9
			"0",
			true,
			nil,
		},
		{ // 10
			"0.0.1",
			false,
			&Version{0, 0, 1},
		},
		{ // 11
			"1,5.0.1",
			true,
			nil,
		},
		{ // 12
			"5",
			false,
			&Version{5, 0, 0},
		},
	}

	for i, test := range tests {
		v, err := Parse(test.input)
		if test.error && err == nil {
			t.Errorf("[%v] Parse(%q) should return error, but does not", i, test.input)
		}
		if !test.error && err != nil {
			t.Errorf("[%v] Parse(%q) should not return error, but does", i, test.input)
		}

		if err != nil {
			continue
		}

		if !v.Equals(*test.version) {
			t.Errorf("[%v] Parse(%q) = %v // expected %v", i, test.input, v, test.version)
		}

		/*
			if v.String() != test.input {
				t.Errorf("[%v] Parse(%q).String() = %q // expected %q", i, test.input, v.String(), test.input)
			}
		*/
	}
}

func TestEqualsMajor(t *testing.T) {

	tests := []struct {
		a        Version
		b        Version
		expected bool
	}{
		{
			Version{0, 0, 1},
			Version{0, 1, 0},
			true,
		},
		{
			Version{0, 0, 1},
			Version{1, 1, 0},
			false,
		},
		{
			Version{1, 2, 1},
			Version{1, 1, 0},
			true,
		},
	}

	for i, test := range tests {
		got := test.a.EqualsMajor(test.b)

		if got != test.expected {
			t.Errorf("[%v] %v.EqualsMajor(%v) = %v // expected: %v", i, test.a, test.b, got, test.expected)
		}
	}
}

func TestEqualsMinor(t *testing.T) {

	tests := []struct {
		a        Version
		b        Version
		expected bool
	}{
		{
			Version{0, 0, 1},
			Version{0, 1, 0},
			false,
		},
		{
			Version{0, 0, 1},
			Version{0, 0, 5},
			true,
		},
		{
			Version{1, 2, 0},
			Version{2, 2, 0},
			false,
		},
	}

	for i, test := range tests {
		got := test.a.EqualsMinor(test.b)

		if got != test.expected {
			t.Errorf("[%v] %v.EqualsMinor(%v) = %v // expected: %v", i, test.a, test.b, got, test.expected)
		}
	}
}

func TestLess(t *testing.T) {

	tests := []struct {
		a        Version
		b        Version
		expected bool
	}{
		{
			Version{0, 0, 1},
			Version{0, 1, 0},
			true,
		},
		{
			Version{1, 1, 0},
			Version{1, 0, 1},
			false,
		},
		{
			Version{1, 2, 1},
			Version{1, 2, 1},
			false,
		},
		{
			Version{1, 2, 1},
			Version{2, 2, 1},
			true,
		},
	}

	for i, test := range tests {
		got := test.a.Less(test.b)

		if got != test.expected {
			t.Errorf("[%v] %v.Less(%v) = %v // expected: %v", i, test.a, test.b, got, test.expected)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{ // 0
			"v1.0.1",
			"1.0.1",
		},
		{ // 1
			"1.0.1",
			"1.0.1",
		},
		{ // 2
			"1.0",
			"1.0.0",
		},
		{ // 3
			"1.0.0",
			"1.0.0",
		},
		{ // 4
			"0.0.1",
			"0.0.1",
		},
		{ // 5
			"5",
			"5.0.0",
		},
	}

	for i, test := range tests {
		v, err := Parse(test.input)
		if err != nil {
			t.Fatalf("[%v] Parse(%q) should not return error, but does", i, test.input)
		}

		if v.String() != test.expected {
			t.Errorf("[%v] Parse(%q).String() = %q // expected %q", i, test.input, v.String(), test.expected)
		}
	}
}

func TestVersions(t *testing.T) {
	var vs = Versions{&Version{1, 0, 2}, &Version{0, 4, 2}, &Version{1, 0, 1}}
	sort.Sort(vs)

	var bd strings.Builder

	for _, v := range vs {
		bd.WriteString(v.String() + "/")
	}

	got := bd.String()

	expected := `0.4.2/1.0.1/1.0.2/`

	if got != expected {
		t.Errorf("Sort(Versions) returned %q ; expected: %q", got, expected)
	}

	first := vs.First().String()
	expectedFirst := "0.4.2"

	if first != expectedFirst {
		t.Errorf("Versions.First() returned %q ; expected: %q", first, expectedFirst)
	}

	last := vs.Last().String()
	expectedLast := "1.0.2"

	if last != expectedLast {
		t.Errorf("Versions.Last() returned %q ; expected: %q", last, expectedLast)
	}
}

func TestVersionsEmpty(t *testing.T) {
	var vs Versions

	if vs.First() != nil {
		t.Errorf("(Versions{}).First() must be nil, but is %v", vs.First())
	}

	if vs.Last() != nil {
		t.Errorf("(Versions{}).Last() must be nil, but is %v", vs.Last())
	}
}
