package mid

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/smf"
	"gitlab.com/gomidi/midi/smf/smfwriter"
)

func msgs(m ...midi.Message) []midi.Message {
	return m
}

func TestRPN_NRPN_JustLSBCallback(t *testing.T) {
	tests := []struct {
		write       func(w *Writer)
		description string
		expected    string
	}{
		// RPN
		{
			func(w *Writer) { w.RPN(0, 0, 23, 0) },
			"simple RPN without LSB",
			"RPN(0,0) LSB on channel 0: 0 | ",
		},
		{
			func(w *Writer) { w.SetChannel(1); w.RPN(0, 0, 23, 0) },
			"simple RPN without LSB channel 1",
			"RPN(0,0) LSB on channel 1: 0 | ",
		},
		{
			func(w *Writer) { w.RPN(1, 1, 43, 54) },
			"simple RPN with LSB",
			"RPN(1,1) LSB on channel 0: 54 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"simple RPN without LSB and without reset",
			"",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"simple RPN (reverse order) without LSB and without reset",
			"",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
				w.Write(channel.Channel0.ControlChange(6, 53))
			},
			"RPN without LSB and without reset and with two MSB values",
			"",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
				w.Write(channel.Channel0.ControlChange(38, 53))
				w.Write(channel.Channel0.ControlChange(6, 23))
			},
			"RPN without reset and with a MSB followed by a LSB, followed by a MSB values",
			"RPN(2,1) LSB on channel 0: 53 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"CC6 value without preceding RPN/NRPN definition",
			"CC6 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel1.ControlChange(38, 13))
			},
			"CC38 value without preceding RPN/NRPN definition",
			"CC38 on channel 1: 13 | ",
		},
		{
			func(w *Writer) {
				w.RPN(1, 2, 23, 0)
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"CC6 value after RPN and reset",
			"RPN(1,2) LSB on channel 0: 0 | CC6 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.RPN(1, 2, 23, 0)
				w.Write(channel.Channel0.ControlChange(38, 13))
			},
			"CC38 value after RPN and reset",
			"RPN(1,2) LSB on channel 0: 0 | CC38 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(101, 2))
				w.Write(channel.Channel2.ControlChange(100, 1))
				w.Write(channel.Channel1.ControlChange(6, 13))
				w.Write(channel.Channel2.ControlChange(6, 23))
			},
			"CC6 value on channel 1 without preceding RPN and on channel 2 with precedding RPN",
			"CC6 on channel 1: 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(101, 2))
				w.Write(channel.Channel2.ControlChange(100, 1))
				w.Write(channel.Channel2.ControlChange(6, 23))
				w.Write(channel.Channel1.ControlChange(6, 13))
				w.Write(channel.Channel2.ControlChange(38, 43))
			},
			"CC6 value on channel 1 without preceding RPN and on channel 2 with precedding RPN and LSB",
			"CC6 on channel 1: 13 | RPN(2,1) LSB on channel 2: 43 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(101, 2))
				w.Write(channel.Channel2.ControlChange(98, 1))
				w.Write(channel.Channel2.ControlChange(6, 23))
			},
			"wrong RPN message, followed by a CC6 value on channel 2",
			"CC6 on channel 2: 23 | ",
		},

		// NRPN
		{
			func(w *Writer) { w.NRPN(0, 0, 23, 0) },
			"simple NRPN without LSB",
			"NRPN(0,0) LSB on channel 0: 0 | ",
		},
		{
			func(w *Writer) { w.SetChannel(1); w.NRPN(0, 0, 23, 0) },
			"simple NRPN without LSB channel 1",
			"NRPN(0,0) LSB on channel 1: 0 | ",
		},
		{
			func(w *Writer) { w.NRPN(1, 1, 43, 54) },
			"simple NRPN with LSB",
			"NRPN(1,1) LSB on channel 0: 54 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"simple NRPN without LSB and without reset",
			"",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"simple NRPN (reverse order) without LSB and without reset",
			"",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
				w.Write(channel.Channel0.ControlChange(6, 53))
			},
			"NRPN without LSB and without reset and with two MSB values",
			"",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
				w.Write(channel.Channel0.ControlChange(38, 53))
				w.Write(channel.Channel0.ControlChange(6, 23))
			},
			"NRPN without reset and with a MSB followed by a LSB, followed by a MSB values",
			"NRPN(2,1) LSB on channel 0: 53 | ",
		},
		{
			func(w *Writer) {
				w.NRPN(1, 2, 23, 0)
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"CC6 value after NRPN and reset",
			"NRPN(1,2) LSB on channel 0: 0 | CC6 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.NRPN(1, 2, 23, 0)
				w.Write(channel.Channel0.ControlChange(38, 13))
			},
			"CC38 value after NRPN and reset",
			"NRPN(1,2) LSB on channel 0: 0 | CC38 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(99, 2))
				w.Write(channel.Channel2.ControlChange(98, 1))
				w.Write(channel.Channel1.ControlChange(6, 13))
				w.Write(channel.Channel2.ControlChange(6, 23))
			},
			"CC6 value on channel 1 without preceding NRPN and on channel 2 with precedding NRPN",
			"CC6 on channel 1: 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(99, 2))
				w.Write(channel.Channel2.ControlChange(98, 1))
				w.Write(channel.Channel2.ControlChange(6, 23))
				w.Write(channel.Channel1.ControlChange(6, 13))
				w.Write(channel.Channel2.ControlChange(38, 43))
			},
			"CC6 value on channel 1 without preceding NRPN and on channel 2 with precedding NRPN and LSB",
			"CC6 on channel 1: 13 | NRPN(2,1) LSB on channel 2: 43 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(99, 2))
				w.Write(channel.Channel2.ControlChange(100, 1))
				w.Write(channel.Channel2.ControlChange(6, 23))
			},
			"wrong NRPN message, followed by a CC6 value on channel 2",
			"CC6 on channel 2: 23 | ",
		},
	}

	for _, test := range tests {
		var bf bytes.Buffer
		var out bytes.Buffer

		rd := NewReader(NoLogger())

		rd.Msg.Channel.ControlChange.Each = func(p *Position, ch, cc, val uint8) {
			fmt.Fprintf(&out, "CC%v on channel %v: %v | ", cc, ch, val)
		}

		rd.Msg.Channel.ControlChange.RPN.LSB = func(p *Position, ch, ident1, ident2, lsb uint8) {
			fmt.Fprintf(&out, "RPN(%v,%v) LSB on channel %v: %v | ", ident1, ident2, ch, lsb)
		}

		rd.Msg.Channel.ControlChange.NRPN.LSB = func(p *Position, ch, ident1, ident2, lsb uint8) {
			fmt.Fprintf(&out, "NRPN(%v,%v) LSB on channel %v: %v | ", ident1, ident2, ch, lsb)
		}

		wr := NewWriter(&bf)
		test.write(wr)

		rd.ReadAllFrom(&bf)

		if got, want := out.String(), test.expected; got != want {
			t.Errorf("%#v\n\tgot  %#v\n\twant %#v", test.description, got, want)
		}

	}

}

func TestRPN_NRPN(t *testing.T) {
	tests := []struct {
		write       func(w *Writer)
		description string
		expected    string
	}{
		// RPN
		{
			func(w *Writer) { w.RPN(0, 0, 23, 0) },
			"simple RPN without LSB",
			"RPN(0,0) on channel 0 MSB 23 | RPN(0,0) LSB on channel 0: 0 | RPN RESET on channel 0 | ",
		},
		{
			func(w *Writer) { w.SetChannel(1); w.RPN(0, 0, 23, 0) },
			"simple RPN without LSB channel 1",
			"RPN(0,0) on channel 1 MSB 23 | RPN(0,0) LSB on channel 1: 0 | RPN RESET on channel 1 | ",
		},
		{
			func(w *Writer) { w.RPN(1, 1, 43, 54) },
			"simple RPN with LSB",
			"RPN(1,1) on channel 0 MSB 43 | RPN(1,1) LSB on channel 0: 54 | RPN RESET on channel 0 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(96, 0))
			},
			"simple RPN with increment",
			"RPN(2,1) Increment on channel 0 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(97, 0))
			},
			"simple RPN with decrement",
			"RPN(2,1) Decrement on channel 0 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"simple RPN without LSB and without reset",
			"RPN(2,1) on channel 0 MSB 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"simple RPN (reverse order) without LSB and without reset",
			"RPN(2,1) on channel 0 MSB 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
				w.Write(channel.Channel0.ControlChange(6, 53))
			},
			"RPN without LSB and without reset and with two MSB values",
			"RPN(2,1) on channel 0 MSB 13 | RPN(2,1) on channel 0 MSB 53 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(101, 2))
				w.Write(channel.Channel0.ControlChange(100, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
				w.Write(channel.Channel0.ControlChange(38, 53))
				w.Write(channel.Channel0.ControlChange(6, 23))
			},
			"RPN without reset and with a MSB followed by a LSB, followed by a MSB values",
			"RPN(2,1) on channel 0 MSB 13 | RPN(2,1) LSB on channel 0: 53 | RPN(2,1) on channel 0 MSB 23 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"CC6 value without preceding RPN/NRPN definition",
			"CC6 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel1.ControlChange(38, 13))
			},
			"CC38 value without preceding RPN/NRPN definition",
			"CC38 on channel 1: 13 | ",
		},
		{
			func(w *Writer) {
				w.RPN(1, 2, 23, 0)
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"CC6 value after RPN and reset",
			"RPN(1,2) on channel 0 MSB 23 | RPN(1,2) LSB on channel 0: 0 | RPN RESET on channel 0 | CC6 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.RPN(1, 2, 23, 0)
				w.Write(channel.Channel0.ControlChange(38, 13))
			},
			"CC38 value after RPN and reset",
			"RPN(1,2) on channel 0 MSB 23 | RPN(1,2) LSB on channel 0: 0 | RPN RESET on channel 0 | CC38 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(101, 2))
				w.Write(channel.Channel2.ControlChange(100, 1))
				w.Write(channel.Channel1.ControlChange(6, 13))
				w.Write(channel.Channel2.ControlChange(6, 23))
			},
			"CC6 value on channel 1 without preceding RPN and on channel 2 with precedding RPN",
			"CC6 on channel 1: 13 | RPN(2,1) on channel 2 MSB 23 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(101, 2))
				w.Write(channel.Channel2.ControlChange(100, 1))
				w.Write(channel.Channel2.ControlChange(6, 23))
				w.Write(channel.Channel1.ControlChange(6, 13))
				w.Write(channel.Channel2.ControlChange(38, 43))
			},
			"CC6 value on channel 1 without preceding RPN and on channel 2 with precedding RPN and LSB",
			"RPN(2,1) on channel 2 MSB 23 | CC6 on channel 1: 13 | RPN(2,1) LSB on channel 2: 43 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(101, 2))
				w.Write(channel.Channel2.ControlChange(98, 1))
				w.Write(channel.Channel2.ControlChange(6, 23))
			},
			"wrong RPN message, followed by a CC6 value on channel 2",
			"CC6 on channel 2: 23 | ",
		},

		// NRPN
		{
			func(w *Writer) { w.NRPN(0, 0, 23, 0) },
			"simple NRPN without LSB",
			"NRPN(0,0) on channel 0 MSB 23 | NRPN(0,0) LSB on channel 0: 0 | NRPN RESET on channel 0 | ",
		},
		{
			func(w *Writer) { w.SetChannel(1); w.NRPN(0, 0, 23, 0) },
			"simple NRPN without LSB channel 1",
			"NRPN(0,0) on channel 1 MSB 23 | NRPN(0,0) LSB on channel 1: 0 | NRPN RESET on channel 1 | ",
		},
		{
			func(w *Writer) { w.NRPN(1, 1, 43, 54) },
			"simple NRPN with LSB",
			"NRPN(1,1) on channel 0 MSB 43 | NRPN(1,1) LSB on channel 0: 54 | NRPN RESET on channel 0 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(96, 0))
			},
			"simple NRPN with increment",
			"NRPN(2,1) Increment on channel 0 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(97, 0))
			},
			"simple NRPN with decrement",
			"NRPN(2,1) Decrement on channel 0 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"simple NRPN without LSB and without reset",
			"NRPN(2,1) on channel 0 MSB 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"simple NRPN (reverse order) without LSB and without reset",
			"NRPN(2,1) on channel 0 MSB 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
				w.Write(channel.Channel0.ControlChange(6, 53))
			},
			"NRPN without LSB and without reset and with two MSB values",
			"NRPN(2,1) on channel 0 MSB 13 | NRPN(2,1) on channel 0 MSB 53 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel0.ControlChange(99, 2))
				w.Write(channel.Channel0.ControlChange(98, 1))
				w.Write(channel.Channel0.ControlChange(6, 13))
				w.Write(channel.Channel0.ControlChange(38, 53))
				w.Write(channel.Channel0.ControlChange(6, 23))
			},
			"NRPN without reset and with a MSB followed by a LSB, followed by a MSB values",
			"NRPN(2,1) on channel 0 MSB 13 | NRPN(2,1) LSB on channel 0: 53 | NRPN(2,1) on channel 0 MSB 23 | ",
		},
		{
			func(w *Writer) {
				w.NRPN(1, 2, 23, 0)
				w.Write(channel.Channel0.ControlChange(6, 13))
			},
			"CC6 value after NRPN and reset",
			"NRPN(1,2) on channel 0 MSB 23 | NRPN(1,2) LSB on channel 0: 0 | NRPN RESET on channel 0 | CC6 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.NRPN(1, 2, 23, 0)
				w.Write(channel.Channel0.ControlChange(38, 13))
			},
			"CC38 value after NRPN and reset",
			"NRPN(1,2) on channel 0 MSB 23 | NRPN(1,2) LSB on channel 0: 0 | NRPN RESET on channel 0 | CC38 on channel 0: 13 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(99, 2))
				w.Write(channel.Channel2.ControlChange(98, 1))
				w.Write(channel.Channel1.ControlChange(6, 13))
				w.Write(channel.Channel2.ControlChange(6, 23))
			},
			"CC6 value on channel 1 without preceding NRPN and on channel 2 with precedding NRPN",
			"CC6 on channel 1: 13 | NRPN(2,1) on channel 2 MSB 23 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(99, 2))
				w.Write(channel.Channel2.ControlChange(98, 1))
				w.Write(channel.Channel2.ControlChange(6, 23))
				w.Write(channel.Channel1.ControlChange(6, 13))
				w.Write(channel.Channel2.ControlChange(38, 43))
			},
			"CC6 value on channel 1 without preceding NRPN and on channel 2 with precedding NRPN and LSB",
			"NRPN(2,1) on channel 2 MSB 23 | CC6 on channel 1: 13 | NRPN(2,1) LSB on channel 2: 43 | ",
		},
		{
			func(w *Writer) {
				w.Write(channel.Channel2.ControlChange(99, 2))
				w.Write(channel.Channel2.ControlChange(100, 1))
				w.Write(channel.Channel2.ControlChange(6, 23))
			},
			"wrong NRPN message, followed by a CC6 value on channel 2",
			"CC6 on channel 2: 23 | ",
		},
	}

	for _, test := range tests {
		var bf bytes.Buffer
		var out bytes.Buffer

		rd := NewReader(NoLogger())

		rd.Msg.Channel.ControlChange.Each = func(p *Position, ch, cc, val uint8) {
			fmt.Fprintf(&out, "CC%v on channel %v: %v | ", cc, ch, val)
		}

		rd.Msg.Channel.ControlChange.RPN.MSB = func(p *Position, ch, ident1, ident2, msb uint8) {
			fmt.Fprintf(&out, "RPN(%v,%v) on channel %v MSB %v | ", ident1, ident2, ch, msb)
		}

		rd.Msg.Channel.ControlChange.RPN.LSB = func(p *Position, ch, ident1, ident2, lsb uint8) {
			fmt.Fprintf(&out, "RPN(%v,%v) LSB on channel %v: %v | ", ident1, ident2, ch, lsb)
		}

		rd.Msg.Channel.ControlChange.RPN.Increment = func(p *Position, ch, ident1, ident2 uint8) {
			fmt.Fprintf(&out, "RPN(%v,%v) Increment on channel %v | ", ident1, ident2, ch)
		}

		rd.Msg.Channel.ControlChange.RPN.Decrement = func(p *Position, ch, ident1, ident2 uint8) {
			fmt.Fprintf(&out, "RPN(%v,%v) Decrement on channel %v | ", ident1, ident2, ch)
		}

		rd.Msg.Channel.ControlChange.RPN.Reset = func(p *Position, ch uint8) {
			fmt.Fprintf(&out, "RPN RESET on channel %v | ", ch)
		}

		rd.Msg.Channel.ControlChange.NRPN.MSB = func(p *Position, ch, ident1, ident2, msb uint8) {
			fmt.Fprintf(&out, "NRPN(%v,%v) on channel %v MSB %v | ", ident1, ident2, ch, msb)
		}

		rd.Msg.Channel.ControlChange.NRPN.Reset = func(p *Position, ch uint8) {
			fmt.Fprintf(&out, "NRPN RESET on channel %v | ", ch)
		}

		rd.Msg.Channel.ControlChange.NRPN.LSB = func(p *Position, ch, ident1, ident2, lsb uint8) {
			fmt.Fprintf(&out, "NRPN(%v,%v) LSB on channel %v: %v | ", ident1, ident2, ch, lsb)
		}

		rd.Msg.Channel.ControlChange.NRPN.Increment = func(p *Position, ch, ident1, ident2 uint8) {
			fmt.Fprintf(&out, "NRPN(%v,%v) Increment on channel %v | ", ident1, ident2, ch)
		}

		rd.Msg.Channel.ControlChange.NRPN.Decrement = func(p *Position, ch, ident1, ident2 uint8) {
			fmt.Fprintf(&out, "NRPN(%v,%v) Decrement on channel %v | ", ident1, ident2, ch)
		}

		wr := NewWriter(&bf)
		test.write(wr)

		rd.ReadAllFrom(&bf)

		if got, want := out.String(), test.expected; got != want {
			t.Errorf("%#v\n\tgot  %#v\n\twant %#v", test.description, got, want)
		}

	}

}

func TestTimeAt(t *testing.T) {
	mt := smf.MetricTicks(960)
	twobars := mt.Ticks4th() * 8

	tests := []struct {
		tempo1          float64
		tempo2          float64
		absPos          uint32
		durationSeconds int64
	}{

		{120, 120, twobars * 2, 8},
		{120, 120, twobars, 4},
		{120, 120, twobars / 2, 2},

		{60, 60, twobars * 2, 16},
		{60, 60, twobars, 8},
		{60, 60, twobars / 2, 4},

		{120, 60, twobars * 2, 12},
		{120, 30, twobars * 2, 20},
		{120, 30, twobars * 3, 36},

		{120, 30, twobars, 4},
	}

	for _, test := range tests {
		var bf bytes.Buffer

		wr := NewSMF(&bf, 1, smfwriter.TimeFormat(mt))
		wr.TempoBPM(test.tempo1)
		wr.SetDelta(twobars)
		wr.TempoBPM(test.tempo2)
		wr.SetDelta(twobars)
		wr.SetDelta(twobars * 8)
		wr.NoteOn(64, 120)
		wr.SetDelta(twobars)
		wr.NoteOff(64)
		wr.EndOfTrack()

		h := NewReader(NoLogger())
		h.ReadAllSMF(&bf)
		d := *h.TimeAt(uint64(test.absPos))
		// ms := int64(d / time.Millisecond)

		if got, want := int64(d/time.Second), test.durationSeconds; got != want {
			t.Errorf("tempo1, tempo2 = %v, %v; TimeAt(%v) = %v; want %v", test.tempo1, test.tempo2, test.absPos, got, want)
		}
	}

}
