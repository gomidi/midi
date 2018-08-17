package syscommon

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

// Raw returns the raw bytes for the message
// TODO Test
func (m SongSelect) Raw() []byte {
	// TODO check - it is a guess
	return []byte{byte(0xF3), byte(m)}
}

// SongSelect represents the MIDI song select system message
type SongSelect uint8

// Number returns the number of the song
func (m SongSelect) Number() uint8 {
	return uint8(m)
}

// String represents the MIDI song select message as a string (for debugging)
func (m SongSelect) String() string {
	return fmt.Sprintf("%T: %v", m, m.Number())
}

func (m SongSelect) sysCommon() {}

// TODO: check
func (m SongSelect) readFrom(rd io.Reader) (Message, error) {

	b, err := midilib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return SongSelect(b), nil
}
