package midicat

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	var bf bytes.Buffer

	bf.WriteString("23 B0E4FF\n3455 B3EEF5\n")

	out1, delta1, err1 := ReadAndConvert(&bf)

	if err1 != nil {
		t.Errorf("error on first read: %s", err1.Error())
		return
	}

	got1 := fmt.Sprintf("[%v] % X", delta1, out1)
	expected1 := "[23] B0 E4 FF"

	if got1 != expected1 {
		t.Errorf("read1: got %q, expected %q", got1, expected1)
	}

	out2, delta2, err2 := ReadAndConvert(&bf)

	if err2 != nil {
		t.Errorf("error on second read: %s", err2.Error())
		return
	}

	got2 := fmt.Sprintf("[%v] % X", delta2, out2)
	expected2 := "[3455] B3 EE F5"

	if got2 != expected2 {
		t.Errorf("read2: got %q, expected %q", got2, expected2)
	}

}
