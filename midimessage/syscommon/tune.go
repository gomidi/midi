package syscommon

import (
	"io"
)

type tune bool

const (
	// Tune represents a MIDI tune request message
	Tune = tune(false)
)

func (m tune) meta() {}

// String represents the MIDI tune request message as a string (for debugging)
func (m tune) String() string {
	return "syscommon.Tune"
}

func (m tune) readFrom(rd io.Reader) (Message, error) {
	return m, nil
}

func (m tune) sysCommon() {}

// Raw returns the raw bytes for the message
// TODO test
func (m tune) Raw() []byte {
	return []byte{byte(0xF6)}
}
