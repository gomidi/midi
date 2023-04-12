package smf

import (
	"testing"

	"gitlab.com/gomidi/midi/v2"
)

func writeFile(file string, sig Message) {
	s := New()
	var t Track
	t.Add(0, sig)
	t.Add(400, midi.NoteOn(0, 64, 33))
	t.Add(400, midi.NoteOff(0, 64))
	t.Close(0)
	s.Add(t)
	s.WriteFile(file)
}

func TestKeys(t *testing.T) {

	tests := []struct {
		sig      func() Message
		file     string
		expected string
	}{
		{CMaj, "c-maj.mid", "CMaj"},

		{GMaj, "g-maj.mid", "GMaj"},
		{DMaj, "d-maj.mid", "DMaj"},
		{AMaj, "a-maj.mid", "AMaj"},
		{EMaj, "e-maj.mid", "EMaj"},
		{BMaj, "h-maj.mid", "BMaj"},
		{FsharpMaj, "fis-maj.mid", "FsharpMaj"},

		{FMaj, "f-maj.mid", "FMaj"},
		{BbMaj, "b-maj.mid", "BbMaj"},
		{EbMaj, "es-maj.mid", "EbMaj"},
		{AbMaj, "as-maj.mid", "AbMaj"},
		{DbMaj, "des-maj.mid", "DbMaj"},
		{GbMaj, "ges-maj.mid", "GbMaj"},

		{AMin, "a-min.mid", "AMin"},

		{BMin, "h-min.mid", "BMin"},
		{CsharpMin, "cis-min.mid", "CsharpMin"},
		{DsharpMin, "dis-min.mid", "DsharpMin"},
		{EMin, "e-min.mid", "EMin"},
		{FsharpMin, "fis-min.mid", "FsharpMin"},
		{GsharpMin, "gis-min.mid", "GsharpMin"},

		{DMin, "d-min.mid", "DMin"},
		{GMin, "g-min.mid", "GMin"},
		{CMin, "c-min.mid", "CMin"},
		{FMin, "f-min.mid", "FMin"},
		{BbMin, "b-min.mid", "BbMin"},
		{EbMin, "es-min.mid", "EbMin"},
	}

	for _, test := range tests {
		// writeFile(test.file, test.sig)
		var k Key
		test.sig().GetMetaKey(&k)

		if got, want := k.String(), test.expected; got != want {
			t.Errorf("%#v = %v; want %v", test.file, got, want)
		}
		/*
		 */
	}

}
