package syscommon

import (
	"fmt"
	"io"
)

type tuneRequest bool

const (
	TuneRequest = tuneRequest(false)
)

func (m tuneRequest) meta() {}

func (m tuneRequest) String() string {
	return fmt.Sprintf("%T", m)
}

func (m tuneRequest) readFrom(rd io.Reader) (Message, error) {
	return m, nil
}

func (m tuneRequest) sysCommon() {}

// TODO test
func (m tuneRequest) Raw() []byte {
	return []byte{byte(0xF6)}
}
