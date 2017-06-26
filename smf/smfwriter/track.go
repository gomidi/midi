package smfwriter

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/internal/runningstatus"
	"github.com/gomidi/midi/internal/vlq"
	// "github.com/gomidi/midi/messages/sysex"
	"io"
)

type track struct {
	chunk         chunk
	runningWriter runningstatus.SMFWriter
}

// <Track Chunk> = <chunk type><length><MTrk event>+
func (t *track) WriteTo(wr io.Writer) (int, error) {
	t.chunk.typ = [4]byte{byte('M'), byte('T'), byte('r'), byte('k')}
	return t.chunk.writeTo(wr)
}

func (t *track) appendToChunk(deltaTime uint32, b []byte) {
	t.chunk.data = append(t.chunk.data, append(vlq.Encode(deltaTime), b...)...)
}

// delta is distance in time to last event in this track (independant of channel)
func (t *track) Add(deltaTime uint32, msg midi.Message) {
	// we have some sort of sysex, so we need to
	// calculate the length of msg[1:]
	// set msg to msg[0] + length of msg[1:] + msg[1:]
	raw := msg.Raw()
	if raw[0] == 0xF0 || raw[0] == 0xF7 {
		//if sys, ok := msg.(sysex.Message); ok {
		b := []byte{raw[0]}
		b = append(b, vlq.Encode(uint32(len(raw)))...)
		if len(raw[1:]) != 0 {
			b = append(b, raw[1:]...)
		}

		t.appendToChunk(deltaTime, b)
		return
	}

	if t.runningWriter != nil {
		t.appendToChunk(deltaTime, t.runningWriter.Write(msg))
		return
	}

	t.appendToChunk(deltaTime, msg.Raw())
}
