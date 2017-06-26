package smfwriter

import (
	"github.com/gomidi/midi/internal/vlq"
	"io"
)

type track struct {
	chunk chunk
}

// <Track Chunk> = <chunk type><length><MTrk event>+
func (t *track) WriteTo(wr io.Writer) (int, error) {
	t.chunk.typ = [4]byte{byte('M'), byte('T'), byte('r'), byte('k')}
	return t.chunk.writeTo(wr)
}

// delta is distance in time to last event in this track (independant of channel)
func (t *track) Add(deltaTime uint32, msg []byte) {
	// we have some sort of sysex, so we need to
	// calculate the length of msg[1:]
	// set msg to msg[0] + length of msg[1:] + msg[1:]
	if msg[0] == 0xF0 || msg[0] == 0xF7 {
		b := []byte{msg[0]}
		b = append(b, vlq.Encode(uint32(len(msg[1:])))...)
		if len(msg[1:]) != 0 {
			b = append(b, msg[1:]...)
		}
		msg = b
	}

	t.chunk.data = append(t.chunk.data, append(vlq.Encode(deltaTime), msg...)...)
}
