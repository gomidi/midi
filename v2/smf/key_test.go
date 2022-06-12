package smf

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
