package midi

import (
	"fmt"
	"strings"
	"testing"
)

func TestResetChannel(t *testing.T) {

	msgs := ResetChannel(2, 5, 7)

	var bd strings.Builder

	for _, msg := range msgs {
		bd.WriteString(fmt.Sprintf("%s\n", msg.String()))
	}

	expected := strings.TrimSpace(`
ControlChange channel: 2 controller: 0 value: 5
ProgramChange channel: 2 program: 7
ControlChange channel: 2 controller: 121 value: 0
ControlChange channel: 2 controller: 7 value: 100
ControlChange channel: 2 controller: 11 value: 127
ControlChange channel: 2 controller: 64 value: 0
ControlChange channel: 2 controller: 10 value: 64
`)

	got := strings.TrimSpace(bd.String())

	if got != expected {
		t.Errorf("got: \n%s\nexpected:\n%s\n", got, expected)
	}

}

func TestSilenceChannelSingle(t *testing.T) {

	msgs := SilenceChannel(4)

	var bd strings.Builder

	for _, msg := range msgs {
		bd.WriteString(fmt.Sprintf("%s\n", msg.String()))
	}

	expected := strings.TrimSpace(`
ControlChange channel: 4 controller: 123 value: 0
ControlChange channel: 4 controller: 120 value: 0
`)

	got := strings.TrimSpace(bd.String())

	if got != expected {
		t.Errorf("got: \n%s\nexpected:\n%s\n", got, expected)
	}

}

func TestSilenceChannelAll(t *testing.T) {

	msgs := SilenceChannel(-1)

	var bd strings.Builder

	for _, msg := range msgs {
		bd.WriteString(fmt.Sprintf("%s\n", msg.String()))
	}

	expected := strings.TrimSpace(`
ControlChange channel: 0 controller: 123 value: 0
ControlChange channel: 0 controller: 120 value: 0
ControlChange channel: 1 controller: 123 value: 0
ControlChange channel: 1 controller: 120 value: 0
ControlChange channel: 2 controller: 123 value: 0
ControlChange channel: 2 controller: 120 value: 0
ControlChange channel: 3 controller: 123 value: 0
ControlChange channel: 3 controller: 120 value: 0
ControlChange channel: 4 controller: 123 value: 0
ControlChange channel: 4 controller: 120 value: 0
ControlChange channel: 5 controller: 123 value: 0
ControlChange channel: 5 controller: 120 value: 0
ControlChange channel: 6 controller: 123 value: 0
ControlChange channel: 6 controller: 120 value: 0
ControlChange channel: 7 controller: 123 value: 0
ControlChange channel: 7 controller: 120 value: 0
ControlChange channel: 8 controller: 123 value: 0
ControlChange channel: 8 controller: 120 value: 0
ControlChange channel: 9 controller: 123 value: 0
ControlChange channel: 9 controller: 120 value: 0
ControlChange channel: 10 controller: 123 value: 0
ControlChange channel: 10 controller: 120 value: 0
ControlChange channel: 11 controller: 123 value: 0
ControlChange channel: 11 controller: 120 value: 0
ControlChange channel: 12 controller: 123 value: 0
ControlChange channel: 12 controller: 120 value: 0
ControlChange channel: 13 controller: 123 value: 0
ControlChange channel: 13 controller: 120 value: 0
ControlChange channel: 14 controller: 123 value: 0
ControlChange channel: 14 controller: 120 value: 0
ControlChange channel: 15 controller: 123 value: 0
ControlChange channel: 15 controller: 120 value: 0
`)

	got := strings.TrimSpace(bd.String())

	if got != expected {
		t.Errorf("got: \n%s\nexpected:\n%s\n", got, expected)
	}

}
