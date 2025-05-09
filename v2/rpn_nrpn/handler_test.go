package rpn_nrpn

import (
	"bytes"
	"fmt"
	"strings"
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
	t.handler.NRPN.MSB = func(channel, typ1, typ2, msbval uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "NRPN.MSB: ch %v t1: %v t2: %v val: %v\n", channel, typ1, typ2, msbval)
		return true
	}

	t.handler.NRPN.LSB = func(channel, typ1, typ2, lsbval uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "NRPN.LSB: ch %v t1: %v t2: %v val: %v\n", channel, typ1, typ2, lsbval)
		return true
	}
	t.handler.NRPN.Increment = func(channel, typ1, typ2 uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "NRPN.Increment: ch %v t1: %v t2: %v\n", channel, typ1, typ2)
		return true
	}
	t.handler.NRPN.Decrement = func(channel, typ1, typ2 uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "NRPN.Decrement: ch %v t1: %v t2: %v\n", channel, typ1, typ2)
		return true
	}
	t.handler.NRPN.Reset = func(channel uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "NRPN.Reset: ch %v\n", channel)
		return true
	}

	t.handler.RPN.MSB = func(channel, typ1, typ2, msbval uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "RPN.MSB: ch %v t1: %v t2: %v val: %v\n", channel, typ1, typ2, msbval)
		return true
	}
	t.handler.RPN.LSB = func(channel, typ1, typ2, lsbval uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "RPN.LSB: ch %v t1: %v t2: %v val: %v\n", channel, typ1, typ2, lsbval)
		return true
	}
	t.handler.RPN.Increment = func(channel, typ1, typ2 uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "RPN.Increment: ch %v t1: %v t2: %v\n", channel, typ1, typ2)
		return true
	}
	t.handler.RPN.Decrement = func(channel, typ1, typ2 uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "RPN.Decrement: ch %v t1: %v t2: %v\n", channel, typ1, typ2)
		return true
	}
	t.handler.RPN.Reset = func(channel uint8) (handled bool) {
		fmt.Fprintf(&t.bf, "RPN.Reset: ch %v\n", channel)
		return true
	}

}

func TestValidHandlers(t *testing.T) {
	var th testHandler

	tests := []struct {
		messages []midi.Message
		//	expectedHandled bool
		expected string
	}{
		{RPNReset(4), `
CC: ch 4 cc: 101 val: 127
RPN.Reset: ch 4
`},
		{NRPNReset(3), `
CC: ch 3 cc: 99 val: 127
NRPN.Reset: ch 3
`},
		{RPNIncrement(12, 44, 75), `
CC: ch 12 cc: 101 val: 44
CC: ch 12 cc: 100 val: 75
RPN.Increment: ch 12 t1: 44 t2: 75
`},
		{NRPNIncrement(2, 34, 45), `
CC: ch 2 cc: 99 val: 34
CC: ch 2 cc: 98 val: 45
NRPN.Increment: ch 2 t1: 34 t2: 45
`},
		{RPNDecrement(12, 44, 75), `
CC: ch 12 cc: 101 val: 44
CC: ch 12 cc: 100 val: 75
RPN.Decrement: ch 12 t1: 44 t2: 75
`},
		{NRPNDecrement(2, 34, 45), `
CC: ch 2 cc: 99 val: 34
CC: ch 2 cc: 98 val: 45
NRPN.Decrement: ch 2 t1: 34 t2: 45
`},

		{RPN(10, 80, 70, 60, 50), `
CC: ch 10 cc: 101 val: 80
CC: ch 10 cc: 100 val: 70
RPN.MSB: ch 10 t1: 80 t2: 70 val: 60
RPN.LSB: ch 10 t1: 80 t2: 70 val: 50
`},
		{append(RPN(10, 80, 70, 60, 50), midi.ControlChange(2, 12, 123)),
			`
CC: ch 10 cc: 101 val: 80
CC: ch 10 cc: 100 val: 70
RPN.MSB: ch 10 t1: 80 t2: 70 val: 60
RPN.LSB: ch 10 t1: 80 t2: 70 val: 50
CC: ch 2 cc: 12 val: 123
`},

		{NRPN(11, 80, 70, 60, 50), `
CC: ch 11 cc: 99 val: 80
CC: ch 11 cc: 98 val: 70
NRPN.MSB: ch 11 t1: 80 t2: 70 val: 60
NRPN.LSB: ch 11 t1: 80 t2: 70 val: 50
`},
		{append(NRPN(11, 80, 70, 60, 50), midi.ControlChange(2, 12, 123)),
			`
CC: ch 11 cc: 99 val: 80
CC: ch 11 cc: 98 val: 70
NRPN.MSB: ch 11 t1: 80 t2: 70 val: 60
NRPN.LSB: ch 11 t1: 80 t2: 70 val: 50
CC: ch 2 cc: 12 val: 123
`},
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

		got := strings.TrimSpace(th.String())
		expected := strings.TrimSpace(test.expected)

		if got != expected {
			t.Errorf("[%v]\ngot:\n%s\nexpected:\n%s", i, got, expected)
		}

	}
}
