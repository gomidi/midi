package meta

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTempo(t *testing.T) {
	bt := []byte{0x03, 0x07, 0xA1, 0x20} // 120 BPM

	var tm Tempo
	tt, err := tm.readFrom(bytes.NewBuffer(bt))

	if err != nil {
		t.Fatalf(err.Error())
	}

	ttt := tt.(Tempo)

	if ttt.BPM() != 120 {
		t.Errorf("wrong tempo: wanted 120, got: %v", ttt.BPM())
	}
}

func TestTempo2(t *testing.T) {
	tt := BPM(120)

	if got, want := tt.Raw(), []byte{0xFF, 0x51, 0x03, 0x07, 0xA1, 0x20}; !reflect.DeepEqual(got, want) {
		t.Errorf("got % X wanted: % X", got, want)
	}
}
