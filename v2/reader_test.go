package midi

import (
	"bytes"
	"fmt"
	"testing"
)

type testreceiver struct {
	bf bytes.Buffer
}

func (r *testreceiver) Receive(msg Message, timestamp int64) {
	fmt.Fprintf(&r.bf, "|%s|", msg)
}

func (r *testreceiver) ReceiveSysEx(b []byte) {
	fmt.Fprintf(&r.bf, "{% X}", b)
}

func (r *testreceiver) ReceiveSysCommon(msg Message, timestamp int64) {
	fmt.Fprintf(&r.bf, "/%s/", msg)
}

func (r *testreceiver) ReceiveRealTime(mtype MsgType, timestamp int64) {
	fmt.Fprintf(&r.bf, "[%s]", mtype)
}

func combine(a ...[]byte) []byte {
	var bf bytes.Buffer

	for _, bt := range a {
		bf.Write(bt)
	}

	return bf.Bytes()
}

func inject(src []byte, pos int, inj []byte) []byte {
	var bf bytes.Buffer
	bf.Write(src[:pos])
	bf.Write(inj)
	bf.Write(src[pos:])
	return bf.Bytes()
}

func injectAndStop(src []byte, pos int, inj []byte) []byte {
	var bf bytes.Buffer
	bf.Write(src[:pos])
	bf.Write(inj)
	return bf.Bytes()
}

func TestReader(t *testing.T) {
	tests := []struct {
		comment  string
		input    []byte
		expected string
	}{
		// TODO test every singe channel, realtime and syscommon message
		{
			"complete sysex message",
			SysEx([]byte{0x41, 0x10, 0x42, 0x12, 0x40, 0x00, 0x7F, 0x00, 0x41}).Data,
			"{41 10 42 12 40 00 7F 00 41}",
		},
		{
			"complete noteon message",
			Channel(2).NoteOn(23, 109).Data,
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|",
		},
		{
			"activesense realtime message",
			Activesense().Data,
			"[ActiveSenseMsg]",
		},
		{
			"complete song select system common message",
			SongSelect(3).Data,
			"/SongSelectMsg song: 3/",
		},
		{
			"sequence of complete messages: sysexmsg, channelmsg, realtimemsg, syscommonmsg",
			combine(
				SysEx([]byte{0x41, 0x10, 0x42, 0x12, 0x40, 0x00, 0x7F, 0x00, 0x41}).Data,
				Channel(2).NoteOn(23, 109).Data,
				Activesense().Data,
				SongSelect(3).Data,
			),
			"{41 10 42 12 40 00 7F 00 41}|Channel2Msg & NoteOnMsg key: 23 velocity: 109|[ActiveSenseMsg]/SongSelectMsg song: 3/",
		},
		{
			"running status (the second message will be 'converted from channel4 noteoff to channel2 noteon, since we cut the status byte)",
			combine(Channel(2).NoteOn(23, 109).Data, Channel(4).NoteOffVelocity(25, 20).Data[1:]),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109||Channel2Msg & NoteOnMsg key: 25 velocity: 20|",
		},
		{
			"discarded sysex message (interrupted via channel message)",
			injectAndStop(SysEx([]byte{0x41, 0x10, 0x42, 0x12, 0x40, 0x00, 0x7F, 0x00, 0x41}).Data, 3, Channel(2).NoteOn(23, 109).Data),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|",
		},
		{
			"discarded sysex message (interrupted via unknown status message)",
			injectAndStop(
				SysEx([]byte{0x41, 0x10, 0x42, 0x12, 0x40, 0x00, 0x7F, 0x00, 0x41}).Data,
				3,
				[]byte{0xF4},
			),
			"",
		},
		{
			"disfunctional sysex message (via injected unknown status message)",
			inject(
				SysEx([]byte{0x41, 0x10}).Data,
				3,
				[]byte{0xF4},
			),
			"",
		},
		{
			"disfunctional sysex message (via injected unknown status message) followed by normal message",
			combine(
				inject(
					SysEx([]byte{0x41, 0x10}).Data,
					3,
					[]byte{0xF4},
				),
				Channel(2).NoteOn(23, 109).Data,
			),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|",
		},
		{
			"complete sysex message with injected realtime message",
			inject(SysEx([]byte{0x41, 0x10, 0x42, 0x12, 0x40, 0x00, 0x7F, 0x00, 0x41}).Data, 3, Activesense().Data),
			"[ActiveSenseMsg]{41 10 42 12 40 00 7F 00 41}",
		},
		{
			"complete channel message with injected realtime message",
			inject(Channel(2).NoteOn(23, 109).Data, 1, Activesense().Data),
			"[ActiveSenseMsg]|Channel2Msg & NoteOnMsg key: 23 velocity: 109|",
		},
		{
			"running status with realtime message in between",
			combine(Channel(2).NoteOn(23, 109).Data, Activesense().Data, Channel(4).NoteOffVelocity(25, 20).Data[1:]),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|[ActiveSenseMsg]|Channel2Msg & NoteOnMsg key: 25 velocity: 20|",
		},
		{
			"running status with injected realtime message",
			combine(Channel(2).NoteOn(23, 109).Data, inject(Channel(4).NoteOffVelocity(25, 20).Data[1:], 1, Activesense().Data)),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|[ActiveSenseMsg]|Channel2Msg & NoteOnMsg key: 25 velocity: 20|",
		},
		{
			"complete syscommon message with injected realtime message",
			inject(SongSelect(3).Data, 1, Activesense().Data),
			"[ActiveSenseMsg]/SongSelectMsg song: 3/",
		},
		{
			"sequence of complete messages: channelmsg, syscommonmsg, channelmsg",
			combine(
				Channel(2).NoteOn(23, 109).Data,
				SongSelect(3).Data,
				Channel(4).NoteOffVelocity(25, 20).Data,
			),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|/SongSelectMsg song: 3/|Channel4Msg & NoteOffMsg key: 25 velocity: 20|",
		},
		{
			"dysfunctional running status because syscommon message resets status byte",
			combine(
				Channel(2).NoteOn(23, 109).Data,
				SongSelect(3).Data,
				Channel(4).NoteOffVelocity(25, 20).Data[1:],
			),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|/SongSelectMsg song: 3/",
		},
		{
			"dysfunctional running status because sysex message resets status byte",
			combine(
				Channel(2).NoteOn(23, 109).Data,
				SysEx([]byte{0x41, 0x10, 0x42, 0x12, 0x40, 0x00, 0x7F, 0x00, 0x41}).Data,
				Channel(4).NoteOffVelocity(25, 20).Data[1:],
			),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|{41 10 42 12 40 00 7F 00 41}",
		},
		{
			"sequence with undefined status message",
			combine(
				Channel(2).NoteOn(23, 109).Data,
				[]byte{0xF4, 0xF5, 0xFD},
				Channel(4).NoteOffVelocity(25, 20).Data,
			),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|[UndefinedMsg]|Channel4Msg & NoteOffMsg key: 25 velocity: 20|",
		},
		{
			"sequence with undefined status message followed by other data",
			combine(
				Channel(2).NoteOn(23, 109).Data,
				[]byte{0xF4, 0x42, 0xF5, 0xFD},
				Channel(4).NoteOffVelocity(25, 20).Data,
			),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|[UndefinedMsg]|Channel4Msg & NoteOffMsg key: 25 velocity: 20|",
		},
		{
			"sequence of complete messages with a random 0xF7 in between",
			combine(
				Channel(2).NoteOn(23, 109).Data,
				[]byte{0xF7},
				Channel(4).NoteOffVelocity(25, 20).Data,
			),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109||Channel4Msg & NoteOffMsg key: 25 velocity: 20|",
		},
		{
			"dysfunctional running status triggered by a random 0xF7 in between",
			combine(
				Channel(2).NoteOn(23, 109).Data,
				[]byte{0xF7},
				Channel(4).NoteOffVelocity(25, 20).Data[1:],
			),
			"|Channel2Msg & NoteOnMsg key: 23 velocity: 109|",
		},
	}

	var rec testreceiver
	pt := newReader(&rec)

	for i, test := range tests {
		rec.bf.Reset()

		pt.Write(test.input, 0)

		if got, expected := rec.bf.String(), test.expected; got != expected {
			t.Errorf("\n\n[%v] // %s\n\n\tbytes:\n\t\t\"% X\"\n\tgot: \n\t\t%q\n\texpected: \n\t\t%q\n\n", i, test.comment, test.input, got, expected)
		}
	}
}
