package drivertest

import (
	"bytes"
	"fmt"
	"testing"
)

func TestMkWithRunningStatus(t *testing.T) {
	messages := mkWithRunningStatus()
	bts := bytes.Join(messages, []byte{})

	got := fmt.Sprintf("% X", bts)

	expected := `92 41 78 37 78 41 00 91 41 14 92 37 00 91 41 00`

	if got != expected {
		t.Errorf("\nexpected: %q\n     got: %q", expected, got)
	}

}

func TestMkSysex(t *testing.T) {
	got := fmt.Sprintf("% X", mkSysex())

	expected := `F0 7E 02 09 01 F7`

	if got != expected {
		t.Errorf("\nexpected: %q\n     got: %q", expected, got)
	}

}
