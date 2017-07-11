package syscommon

import (
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/midilib"
)

func (m SongSelect) Pos() uint16 {
	return uint16(m)
}

// TODO Test
func (m SongSelect) Raw() []byte {
	// TODO check - it is a guess
	return []byte{byte(0xF3), byte(m)}
}

type SongSelect uint8

func (m SongSelect) Number() uint8 {
	return uint8(m)
}

func (m SongSelect) String() string {
	return fmt.Sprintf("%T: %v", m, uint8(m))
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
