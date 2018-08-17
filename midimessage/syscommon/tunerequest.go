package syscommon

import (
	"fmt"
	"io"
)

type tuneRequest bool

const (
	// TuneRequest represents a MIDI tune request message
	TuneRequest = tuneRequest(false)
)

func (m tuneRequest) meta() {}

// String represents the MIDI tune request message as a string (for debugging)
func (m tuneRequest) String() string {
	return fmt.Sprintf("%T", m)
}

func (m tuneRequest) readFrom(rd io.Reader) (Message, error) {
	return m, nil
}

func (m tuneRequest) sysCommon() {}

// Raw returns the raw bytes for the message
// TODO test
func (m tuneRequest) Raw() []byte {
	return []byte{byte(0xF6)}
}
