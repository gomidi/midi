package smfreader

import (
	"bytes"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfwriter"
	"testing"
)

func TestParseTimeCode(t *testing.T) {

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
		timecode := parseTimeCode(test.input)

		if got, want := timecode.FramesPerSecond, test.fps; got != want {
			t.Errorf("parseTimeCode(% X) [fps] = %v; want %v", test.input, got, want)
		}

		if got, want := timecode.SubFrames, test.subframes; got != want {
			t.Errorf("parseTimeCode(% X) [subframes] = %v; want %v", test.input, got, want)
		}
	}

}

func TestTimeCode(t *testing.T) {

	tests := []struct {
		fps       uint8
		subframes uint8
		format    smf.TimeFormat
	}{
		{24, 8, smf.SMPTE24(8)},
		{25, 40, smf.SMPTE25(40)},
		{30, 100, smf.SMPTE30(100)},
		{29, 80, smf.SMPTE30DropFrame(80)},
	}

	for _, test := range tests {

		var bf bytes.Buffer
		_, err := smfwriter.New(&bf, smfwriter.TimeFormat(test.format)).Write(meta.Tempo(100))

		if err != nil {
			t.Fatalf("can't write smf: %v", err)
		}

		rd := New(bytes.NewReader(bf.Bytes()))

		var header smf.Header

		header, err = rd.ReadHeader()

		if err != nil {
			t.Fatalf("can't write read header: %v", err)
		}

		format := header.TimeFormat()

		tc, isTC := format.(smf.TimeCode)

		if !isTC {
			t.Fatalf("wrong time format: %#v; expected TimeCode", format)
		}

		if tc.FramesPerSecond != test.fps {
			t.Fatalf("wrong fps: %v; expected: %v", tc.FramesPerSecond, test.fps)
		}

		if tc.SubFrames != test.subframes {
			t.Fatalf("wrong subframes: %v; expected: %v", tc.SubFrames, test.subframes)
		}

	}

}
