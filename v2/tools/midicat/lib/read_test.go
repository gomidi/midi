package lib

import (
	"bytes"
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	var bf bytes.Buffer

	bf.WriteString("B0E4FF\nB3EEF5\n")

	out1, err1 := ReadAndConvert(&bf)

	if err1 != nil {
		t.Errorf("error on first read: %s", err1.Error())
		return
	}

	got1 := fmt.Sprintf("% X", out1)
	expected1 := "B0 E4 FF"

	if got1 != expected1 {
		t.Errorf("read1: got %q, expected %q", got1, expected1)
	}

	out2, err2 := ReadAndConvert(&bf)

	if err2 != nil {
		t.Errorf("error on second read: %s", err2.Error())
		return
	}

	got2 := fmt.Sprintf("% X", out2)
	expected2 := "B3 EE F5"

	if got2 != expected2 {
		t.Errorf("read2: got %q, expected %q", got2, expected2)
	}

}
