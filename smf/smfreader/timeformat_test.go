package smfreader

import (
	"bytes"
	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfwriter"
	"testing"
)

func TestTimeCode(t *testing.T) {

	tests := []struct {
		fps       uint8
		subframes uint8
		option    smfwriter.Option
	}{
		{24, 8, smfwriter.SMPTE24(8)},
		{25, 40, smfwriter.SMPTE25(40)},
		{30, 100, smfwriter.SMPTE30(100)},
		{29, 80, smfwriter.SMPTE30DropFrame(80)},
	}

	for _, test := range tests {

		var bf bytes.Buffer
		_, err := smfwriter.New(&bf, test.option).Write(meta.Tempo(100))

		if err != nil {
			t.Fatalf("can't write smf: %v", err)
		}

		rd := New(bytes.NewReader(bf.Bytes()))

		var header smf.Header

		header, err = rd.ReadHeader()

		if err != nil {
			t.Fatalf("can't write read header: %v", err)
		}

		fm, val := header.TimeFormat()

		if fm != smf.TimeCode {
			t.Fatalf("wrong time format: %s; expected: %s", fm.String(), smf.TimeCode.String())
		}

		fps, subframes := smf.UnpackTimeCode(val)

		if fps != test.fps {
			t.Fatalf("wrong fps: %v; expected: %v", fps, test.fps)
		}

		if subframes != test.subframes {
			t.Fatalf("wrong subframes: %v; expected: %v", subframes, test.subframes)
		}

	}

}
