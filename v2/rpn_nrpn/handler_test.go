package rpn_nrpn

import (
	"bytes"
	"fmt"
	"testing"

	"gitlab.com/gomidi/midi/v2"
)

type testHandler struct {
	handler Handler
	bf      bytes.Buffer
}

func (t *testHandler) String() string {
	return t.bf.String()
}

func (t *testHandler) handleCC(channel, cc, val uint8) {
	fmt.Fprintf(&t.bf, "CC: ch %v cc: %v val: %v\n", channel, cc, val)
}

func (t *testHandler) reset() {
	t.bf.Reset()
	t.handler.NRPN.MSB = func(channel, typ1, typ2, msbval uint8) {
		fmt.Fprintf(&t.bf, "NRPN.MSB: ch %v t1: %v t2: %v val: %v\n", channel, typ1, typ2, msbval)
	}

	t.handler.NRPN.LSB = func(channel, typ1, typ2, lsbval uint8) {
		fmt.Fprintf(&t.bf, "NRPN.LSB: ch %v t1: %v t2: %v val: %v\n", channel, typ1, typ2, lsbval)
	}
	t.handler.NRPN.Increment = func(channel, typ1, typ2 uint8) {
		fmt.Fprintf(&t.bf, "NRPN.Increment: ch %v t1: %v t2: %v\n", channel, typ1, typ2)
	}
	t.handler.NRPN.Decrement = func(channel, typ1, typ2 uint8) {
		fmt.Fprintf(&t.bf, "NRPN.Decrement: ch %v t1: %v t2: %v\n", channel, typ1, typ2)
	}
	t.handler.NRPN.Reset = func(channel uint8) {
		fmt.Fprintf(&t.bf, "NRPN.Reset: ch %v\n", channel)
	}

	t.handler.RPN.MSB = func(channel, typ1, typ2, msbval uint8) {
		fmt.Fprintf(&t.bf, "RPN.MSB: ch %v t1: %v t2: %v val: %v\n", channel, typ1, typ2, msbval)
	}
	t.handler.RPN.LSB = func(channel, typ1, typ2, lsbval uint8) {
		fmt.Fprintf(&t.bf, "RPN.LSB: ch %v t1: %v t2: %v val: %v\n", channel, typ1, typ2, lsbval)
	}
	t.handler.RPN.Increment = func(channel, typ1, typ2 uint8) {
		fmt.Fprintf(&t.bf, "RPN.Increment: ch %v t1: %v t2: %v\n", channel, typ1, typ2)
	}
	t.handler.RPN.Decrement = func(channel, typ1, typ2 uint8) {
		fmt.Fprintf(&t.bf, "RPN.Decrement: ch %v t1: %v t2: %v\n", channel, typ1, typ2)
	}
	t.handler.RPN.Reset = func(channel uint8) {
		fmt.Fprintf(&t.bf, "RPN.Reset: ch %v\n", channel)
	}

}

func TestValidHandlers(t *testing.T) {
	var th testHandler

	tests := []struct {
		messages []midi.Message
		//	expectedHandled bool
		expected string
	}{
		{RPNReset(4), "RPN.Reset: ch 4\n"},
		{NRPNReset(3), "NRPN.Reset: ch 3\n"},
		{RPNIncrement(12, 44, 75), "RPN.Increment: ch 12 t1: 44 t2: 75\nRPN.Reset: ch 12\n"},
		{NRPNIncrement(2, 34, 45), "NRPN.Increment: ch 2 t1: 34 t2: 45\nNRPN.Reset: ch 2\n"},
		{RPNDecrement(12, 44, 75), "RPN.Decrement: ch 12 t1: 44 t2: 75\nRPN.Reset: ch 12\n"},
		{NRPNDecrement(2, 34, 45), "NRPN.Decrement: ch 2 t1: 34 t2: 45\nNRPN.Reset: ch 2\n"},

		{RPN(10, 80, 70, 60, 50), "RPN.MSB: ch 10 t1: 80 t2: 70 val: 60\nRPN.LSB: ch 10 t1: 80 t2: 70 val: 50\nRPN.Reset: ch 10\n"},
		{append(RPN(10, 80, 70, 60, 50), midi.ControlChange(2, 12, 123)),
			"RPN.MSB: ch 10 t1: 80 t2: 70 val: 60\nRPN.LSB: ch 10 t1: 80 t2: 70 val: 50\nRPN.Reset: ch 10\nCC: ch 2 cc: 12 val: 123\n"},

		{NRPN(11, 80, 70, 60, 50), "NRPN.MSB: ch 11 t1: 80 t2: 70 val: 60\nNRPN.LSB: ch 11 t1: 80 t2: 70 val: 50\nNRPN.Reset: ch 11\n"},
		{append(NRPN(11, 80, 70, 60, 50), midi.ControlChange(2, 12, 123)),
			"NRPN.MSB: ch 11 t1: 80 t2: 70 val: 60\nNRPN.LSB: ch 11 t1: 80 t2: 70 val: 50\nNRPN.Reset: ch 11\nCC: ch 2 cc: 12 val: 123\n"},
	}

	for i, test := range tests {
		th.reset()

		for _, msg := range test.messages {
			var ch, cc, val uint8
			if msg.GetControlChange(&ch, &cc, &val) {
				used := th.handler.AddMessage(ch, cc, val)
				if !used {
					th.handleCC(ch, cc, val)
				}
			}
		}

		got := th.String()
		expected := test.expected

		if got != expected {
			t.Errorf("[%v]\ngot:\n%s\nexpected:\n%s", i, got, expected)
		}

	}
}
