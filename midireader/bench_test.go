package midireader

import (
	"bytes"
	"io"
	"testing"

	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midiwriter"
)

type testreader struct {
	ptr   int
	bytes []byte
}

func (r *testreader) Read(b []byte) (int, error) {

	for i := 0; i < len(b); i++ {
		if r.ptr >= len(r.bytes) {
			r.ptr = 0
		}

		b[i] = r.bytes[r.ptr]

		r.ptr++
	}

	return len(b), nil
}

func sameChannel() io.Reader {
	var bf bytes.Buffer

	wr := midiwriter.New(&bf)

	var (
		m1 = channel.Channel1.NoteOn(20, 100)
		m2 = channel.Channel1.NoteOff(20)
		m3 = channel.Channel1.NoteOn(23, 70)
		m4 = channel.Channel1.NoteOff(23)
	)

	wr.Write(m1)
	wr.Write(m2)
	wr.Write(m3)
	wr.Write(m4)

	return &testreader{0, bf.Bytes()}
}

func alternatingChannel() io.Reader {
	var bf bytes.Buffer

	wr := midiwriter.New(&bf)

	var (
		m1 = channel.Channel1.NoteOn(20, 100)
		m2 = channel.Channel4.NoteOn(23, 70)
		m3 = channel.Channel1.NoteOff(20)
		m4 = channel.Channel4.NoteOff(23)
	)

	wr.Write(m1)
	wr.Write(m2)
	wr.Write(m3)
	wr.Write(m4)

	return &testreader{0, bf.Bytes()}
}

// BenchmarkNoteOnOffSameChannel1000 reads 1000 channel messages per iteration
// which are noteon or noteoff messages on the same channel.
// running status is used, since they all have the same status byte
// messages are not handled
func BenchmarkNoteOnOffSameChannel1000(b *testing.B) {
	b.StopTimer()

	src := sameChannel()
	rd := New(src, nil)

	var err error

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_, err = rd.Read()
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	}

}

// BenchmarkNoteOnOffAlternatingChannel1000 reads 1000 channel messages per iteration
// which are noteon or noteoff messages alternating on different channels.
// therefore no running status
func BenchmarkNoteOnOffAlternatingChannel1000(b *testing.B) {
	b.StopTimer()

	src := alternatingChannel()
	rd := New(src, nil)

	var err error

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			_, err = rd.Read()
			if err != nil {
				b.Fatalf("Error: %v", err)
			}
		}
	}

}
