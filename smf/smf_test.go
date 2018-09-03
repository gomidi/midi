package smf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"time"
)

func TestChunk(t *testing.T) {
	var bf bytes.Buffer
	var ch Chunk
	ch.SetType([4]byte{byte('M'), byte('T'), byte('h'), byte('d')})

	var hd Header
	hd.Format = SMF0
	hd.NumTracks = 1
	hd.TimeFormat = MetricTicks(1920)

	binary.Write(&bf, binary.BigEndian, hd.Format.Type())
	binary.Write(&bf, binary.BigEndian, hd.NumTracks)

	var tf = MetricTicks(1920)
	ticks := tf.Ticks4th()
	binary.Write(&bf, binary.BigEndian, uint16(ticks))

	ch.Write(bf.Bytes())

	var out bytes.Buffer
	ch.WriteTo(&out)

	expected := "4D 54 68 64 00 00 00 06 00 00 00 01 07 80"
	if got, want := fmt.Sprintf("% X", out.Bytes()), expected; got != want {
		t.Errorf("Chunk.Bytes() = %#v; want %#v", got, want)
	}

	ch = Chunk{}
	ch.ReadHeader(&out)

	expected = "MThd"

	if got, want := ch.Type(), expected; got != want {
		t.Errorf("Chunk.Type() = %#v; want %#v", got, want)
	}
}

func TestMetricTicks(t *testing.T) {
	var tf = MetricTicks(1920)

	expected := uint32(tf.Number()) / 2

	if got, want := tf.Ticks8th(), expected; got != want {
		t.Errorf("MetricTicks(1920).Ticks8th() ==  = %#v; want %#v", got, want)
	}

}

func TestHeaderSMF2(t *testing.T) {
	var h Header
	h.Format = SMF2
	h.NumTracks = 10
	h.TimeFormat = SMPTE30DropFrame(100)
	expected := "<Format: SMF2 (sequential tracks), NumTracks: 10, TimeFormat: SMPTE30DropFrame 100 subframes>"

	if got, want := h.String(), expected; got != want {
		t.Errorf("Header.String() = %#v; want %#v", got, want)
	}

}

func TestHeaderSMF1(t *testing.T) {
	var h Header
	h.Format = SMF1
	h.NumTracks = 10
	h.TimeFormat = SMPTE30(100)
	expected := "<Format: SMF1 (multitrack), NumTracks: 10, TimeFormat: SMPTE30 100 subframes>"

	if got, want := h.String(), expected; got != want {
		t.Errorf("Header.String() = %#v; want %#v", got, want)
	}

}

func TestHeaderSMF0(t *testing.T) {
	var h Header
	h.Format = SMF0
	h.NumTracks = 1
	h.TimeFormat = MetricTicks(0)
	expected := "<Format: SMF0 (singletrack), NumTracks: 1, TimeFormat: 960 MetricTicks>"

	if got, want := h.String(), expected; got != want {
		t.Errorf("Header.String() = %#v; want %#v", got, want)
	}

}

func TestTicksNoteLengths(t *testing.T) {
	var resolution MetricTicks

	tests := []struct {
		description string
		got         uint32
		expected    uint32
	}{
		{"4th", resolution.Ticks4th(), uint32(resolution.Number())},
		{"8th", resolution.Ticks8th(), 480},
		{"16th", resolution.Ticks16th(), 240},
		{"32th", resolution.Ticks32th(), 120},
		{"64th", resolution.Ticks64th(), 60},
		{"128th", resolution.Ticks128th(), 30},
		{"256th", resolution.Ticks256th(), 15},
		{"512th", resolution.Ticks512th(), 8},
		{"1024th", resolution.Ticks1024th(), 4},
	}

	for _, test := range tests {

		if got, want := test.got, test.expected; got != want {
			t.Errorf(
				"%s = %v; want %v",
				test.description,
				got,
				want,
			)
		}
	}

}

func TestTicksDuration(t *testing.T) {
	tests := []struct {
		resolution MetricTicks
		tempo      uint32
		deltaTicks uint32
		duration   time.Duration
	}{
		{96, 120, 96, 500 * time.Millisecond},
		{96, 120, 48, 250 * time.Millisecond},
		{96, 120, 192, 1000 * time.Millisecond},
		{90, 240, 90, 250 * time.Millisecond},
	}

	for _, test := range tests {

		if got, want := test.resolution.Duration(test.tempo, test.deltaTicks), test.duration; got != want {
			t.Errorf(
				"MetricTicks(%v).Duration(%v, %v) = %s; want %s",
				uint16(test.resolution),
				test.tempo,
				test.deltaTicks,
				got,
				want,
			)
		}

		if got, want := test.resolution.Ticks(test.tempo, test.duration), test.deltaTicks; got != want {
			t.Errorf(
				"MetricTicks(%v).Ticks(%v, %v) = %v; want %v",
				uint16(test.resolution),
				test.tempo,
				test.duration,
				got,
				want,
			)
		}
	}

}
