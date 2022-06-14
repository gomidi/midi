package smf

import (
	"reflect"
	"testing"
)

func TestTempo(t *testing.T) {
	bt := []byte{0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20} // 120 BPM

	tm := Message(bt)

	var bpm float64
	if tm.GetMetaTempo(&bpm) {
		if bpm != 120 {
			t.Errorf("wrong tempo: wanted 120, got: %v", bpm)
		}
	} else {
		t.Fatalf("is no proper tempo message")
	}
}

func TestTempo2(t *testing.T) {
	tt := MetaTempo(120)

	if got, want := tt.Bytes(), []byte{0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20}; !reflect.DeepEqual(got, want) {
		t.Errorf("got % X wanted: % X", got, want)
	}
}
