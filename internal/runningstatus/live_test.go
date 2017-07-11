package runningstatus

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/gomidi/midi"
	"github.com/gomidi/midi/messages/channel"
)

func msgs(m ...midi.Message) []midi.Message {
	return m
}

func TestLiveWriter(t *testing.T) {

	tests := []struct {
		input    []midi.Message
		expected string
	}{
		// the following examples are taken from the SMF format spec
		{
			// single message - no running status
			msgs(channel.Ch2.NoteOn(48, 96)),
			"92 30 60",
		},
		{
			// single message - no running status
			msgs(channel.Ch2.NoteOn(60, 96)),
			"92 3C 60",
		},
		{
			// running status (same channel, same message type)
			msgs(
				channel.Ch2.NoteOn(48, 96),
				channel.Ch2.NoteOn(60, 96),
			),
			"92 30 60" +
				" 3C 60", // running status
		},
		{
			msgs(
				channel.Ch2.NoteOn(48, 96),
				channel.Ch2.NoteOn(60, 96), // running status
				channel.Ch1.NoteOn(67, 64), // no running status (channel change)
			),
			"92 30 60" +
				" 3C 60" + // running status
				" 91 43 40", // no running status (channel change)
		},
		{
			msgs(
				channel.Ch2.NoteOn(48, 96),
				channel.Ch2.NoteOn(60, 96), // running status
				channel.Ch1.NoteOn(67, 64), // no running status (channel change)
				channel.Ch0.NoteOn(76, 32), // no running status (channel change)
			),
			"92 30 60" +
				" 3C 60" + // running status
				" 91 43 40" + // no running status (channel change)
				" 90 4C 20", // no running status (channel change)
		},

		// own variations
		{
			// running status (same channel, same message type)
			msgs(
				channel.Ch2.NoteOn(48, 96),
				channel.Ch2.NoteOn(48, 0), // noteoff simulation
			),
			"92 30 60" +
				" 30 00", // running status
		},
		{
			// running status (same channel, same message type)
			msgs(
				channel.Ch2.NoteOn(48, 96),
				// noteoff should by default result in noteon message with velocity 0,
				// so that running status is active
				channel.Ch2.NoteOff(48),
			),
			"92 30 60" +
				" 30 00", // running status
		},
		{
			// no running status (same channel, different message type)
			msgs(
				channel.Ch2.NoteOn(48, 96),
				// NoteOffPedantic creates a "real" noteoff message with the given velocity,
				// that is a different message type than noteon, so running status is not active
				channel.Ch2.NoteOffVelocity(48, 96),
			),
			"92 30 60" +
				" 82 30 60", // no running status
		},
	}

	for _, test := range tests {
		var bf bytes.Buffer
		wr := NewLiveWriter(&bf)

		var input bytes.Buffer

		_ = wr.Write

		for _, m := range test.input {
			wr.Write(m.Raw())
			input.WriteString(m.String() + "\n")
		}

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("NewLiveWriter().Write(%s) = %v; want %v", input.String(), got, want)
		}
	}

}

func _bytes(b ...byte) []byte {
	return b
}

func TestLiveReader(t *testing.T) {

	tests := []struct {
		input    []byte
		expected string
	}{
		// the following examples are taken from the SMF format spec
		{
			// single message - no running status
			_bytes(0x92),
			"92",
		},
		{
			// running status (same channel, same message type)
			_bytes(0x92, 0x3C),
			"92 92", // running status
		},
		{
			_bytes(0x92, 0x3C, 0x91),
			"92 92 91", // running status; no running status
		},
		{
			_bytes(0x92, 0x3C, 0x91, 0x90),
			"92 92 91 90", // running status; no running status; no running status
		},
		{
			_bytes(0x92, 0x3C, 0x91, 0x90, 0x3C),
			"92 92 91 90 90", // running status; no running status; no running status; running status
		},
	}

	for _, test := range tests {
		var res = make([]byte, len(test.input))
		rd := NewLiveReader()

		for i, b := range test.input {
			r, _ := rd.Read(b)
			res[i] = r
		}

		if got, want := fmt.Sprintf("% X", res), test.expected; got != want {
			t.Errorf("NewLiveReader().Read(% X) = %v; want %v", test.input, got, want)
		}
	}

}
