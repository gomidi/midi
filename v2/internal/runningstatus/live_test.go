package runningstatus

import (
	"bytes"
	"fmt"
	"testing"

	"gitlab.com/gomidi/midi/v2"
)

func msgs(m ...midi.Message) []midi.Message {
	return m
}

func TestSMFWriter(t *testing.T) {

	tests := []struct {
		input    []midi.Message
		expected string
	}{
		// the following examples are taken from the SMF format spec
		{
			// single message - no running status
			msgs(midi.NoteOn(2, 48, 96)),
			"92 30 60",
		},
		{
			// single message - no running status
			msgs(midi.NoteOn(2, 60, 96)),
			"92 3C 60",
		},
		{
			// running status (same channel, same message type)
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOn(2, 60, 96),
			),
			"92 30 60" +
				" 3C 60", // running status
		},
		{
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOn(2, 60, 96), // running status
				midi.NoteOn(1, 67, 64), // no running status (channel change)
			),
			"92 30 60" +
				" 3C 60" + // running status
				" 91 43 40", // no running status (channel change)
		},
		{
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOn(2, 60, 96), // running status
				midi.NoteOn(1, 67, 64), // no running status (channel change)
				midi.NoteOn(0, 76, 32), // no running status (channel change)
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
				midi.NoteOn(2, 48, 96),
				midi.NoteOn(2, 48, 0), // noteoff simulation
			),
			"92 30 60" +
				" 30 00", // running status
		},
		{
			// running status (same channel, same message type)
			msgs(
				midi.ControlChange(2, 48, 96),
				midi.ControlChange(2, 58, 0),
			),
			"B2 30 60" +
				" 3A 00", // running status
		},
		{
			// no running status (same channel, different message type)
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOff(2, 48),
			),
			"92 30 60" +
				" 82 30 00", // no running status
		},
		{
			// no running status (same channel, different message type)
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOffVelocity(2, 48, 96),
			),
			"92 30 60" +
				" 82 30 60", // no running status
		},
		{
			// running status (same channel, same message type)
			msgs(
				midi.ControlChange(2, 48, 96),
				midi.PolyAfterTouch(2, 58, 0),
			),
			"B2 30 60 " +
				"A2 3A 00", // running status
		},
	}

	for _, test := range tests {
		// var bf bytes.Buffer
		wr := NewSMFWriter()

		var bf bytes.Buffer
		var input bytes.Buffer

		// _ = wr.Write

		for _, m := range test.input {
			bf.Write(wr.Write(m))
			input.WriteString(m.String() + "\n")
		}

		if got, want := fmt.Sprintf("% X", bf.Bytes()), test.expected; got != want {
			t.Errorf("NewSMFWriter().Write(%s) = %v; want %v", input.String(), got, want)
		}
	}

}

func TestLiveWriter(t *testing.T) {

	tests := []struct {
		input    []midi.Message
		expected string
	}{
		// the following examples are taken from the SMF format spec
		{
			// single message - no running status
			msgs(midi.NoteOn(2, 48, 96)),
			"92 30 60",
		},
		{
			// single message - no running status
			msgs(midi.NoteOn(2, 60, 96)),
			"92 3C 60",
		},
		{
			// running status (same channel, same message type)
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOn(2, 60, 96),
			),
			"92 30 60" +
				" 3C 60", // running status
		},
		{
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOn(2, 60, 96), // running status
				midi.NoteOn(1, 67, 64), // no running status (channel change)
			),
			"92 30 60" +
				" 3C 60" + // running status
				" 91 43 40", // no running status (channel change)
		},
		{
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOn(2, 60, 96), // running status
				midi.NoteOn(1, 67, 64), // no running status (channel change)
				midi.NoteOn(0, 76, 32), // no running status (channel change)
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
				midi.NoteOn(2, 48, 96),
				midi.NoteOn(2, 48, 0), // noteoff simulation
			),
			"92 30 60" +
				" 30 00", // running status
		},
		{
			// running status (same channel, same message type)
			msgs(
				midi.ControlChange(2, 48, 96),
				midi.ControlChange(2, 58, 0),
			),
			"B2 30 60" +
				" 3A 00", // running status
		},
		{
			// no running status (same channel, different message type)
			msgs(
				midi.NoteOn(2, 48, 96),
				midi.NoteOff(2, 48),
			),
			"92 30 60" +
				" 82 30 00", // no running status
		},
		{
			// no running status (same channel, different message type)
			msgs(
				midi.NoteOn(2, 48, 96),
				// NoteOffPedantic creates a "real" noteoff message with the given velocity,
				// that is a different message type than noteon, so running status is not active
				midi.NoteOffVelocity(2, 48, 96),
			),
			"92 30 60" +
				" 82 30 60", // no running status
		},
		{
			// running status (same channel, same message type)
			msgs(
				midi.ControlChange(2, 48, 96),
				midi.PolyAfterTouch(2, 58, 0),
			),
			"B2 30 60 " +
				"A2 3A 00", // running status
		},
	}

	for _, test := range tests {
		var bf bytes.Buffer
		wr := NewLiveWriter(&bf)

		var input bytes.Buffer

		_ = wr.Write

		for _, m := range test.input {
			wr.Write(m.Bytes())
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

		// own examples
		{
			_bytes(0xB2, 0x30, 0x60, 0x3A, 0x00),
			"B2 B2 B2 B2 B2", // running status
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

func TestSMFReader(t *testing.T) {

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

		// own examples
		{
			_bytes(0xB2, 0x30, 0x60, 0x3A, 0x00),
			"B2 B2 B2 B2 B2", // running status
		},
	}

	for _, test := range tests {
		var res = make([]byte, len(test.input))
		rd := NewSMFReader()

		for i, b := range test.input {
			r, _ := rd.Read(b)
			res[i] = r
		}

		if got, want := fmt.Sprintf("% X", res), test.expected; got != want {
			t.Errorf("NewSMFReader().Read(% X) = %v; want %v", test.input, got, want)
		}
	}

}
