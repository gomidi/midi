package key

import (
	"os"
	"testing"

	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/midi/smf/smfwriter"
)

func writeFile(file string, sig meta.Key) {
	f, _ := os.Create(file)
	wr := smfwriter.New(f, smfwriter.NumTracks(1))
	wr.WriteHeader()
	wr.Write(sig)
	wr.SetDelta(400)
	wr.Write(channel.Channel0.NoteOn(64, 33))
	wr.SetDelta(400)
	wr.Write(channel.Channel0.NoteOff(64))
	wr.Write(meta.EndOfTrack)
	f.Close()
}

func TestKeys(t *testing.T) {

	tests := []struct {
		sig      func() meta.Key
		file     string
		expected string
	}{
		{CMaj, "c-maj.mid", "C maj."},

		{GMaj, "g-maj.mid", "G maj."},
		{DMaj, "d-maj.mid", "D maj."},
		{AMaj, "a-maj.mid", "A maj."},
		{EMaj, "e-maj.mid", "E maj."},
		{BMaj, "h-maj.mid", "B maj."},
		{FSharpMaj, "fis-maj.mid", "F♯ maj."},

		{FMaj, "f-maj.mid", "F maj."},
		{BFlatMaj, "b-maj.mid", "B♭ maj."},
		{EFlatMaj, "es-maj.mid", "E♭ maj."},
		{AFlatMaj, "as-maj.mid", "A♭ maj."},
		{DFlatMaj, "des-maj.mid", "D♭ maj."},
		{GFlatMaj, "ges-maj.mid", "G♭ maj."},

		{AMin, "a-min.mid", "A min."},

		{BMin, "h-min.mid", "B min."},
		{CSharpMin, "cis-min.mid", "C♯ min."},
		{DSharpMin, "dis-min.mid", "D♯ min."},
		{EMin, "e-min.mid", "E min."},
		{FSharpMin, "fis-min.mid", "F♯ min."},
		{GSharpMin, "gis-min.mid", "G♯ min."},

		{DMin, "d-min.mid", "D min."},
		{GMin, "g-min.mid", "G min."},
		{CMin, "c-min.mid", "C min."},
		{FMin, "f-min.mid", "F min."},
		{BFlatMin, "b-min.mid", "B♭ min."},
		{EFlatMin, "es-min.mid", "E♭ min."},
	}

	for _, test := range tests {
		// writeFile(test.file, test.sig)
		if got, want := test.sig().Text(), test.expected; got != want {
			t.Errorf("%#v = %v; want %v", test.file, got, want)
		}
		/*
		 */
	}

}
