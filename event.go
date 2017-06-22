package midi

import "fmt"

type Event interface {
	String() string
	Raw() []byte
}

type UnknownEvent []byte

func (u UnknownEvent) String() string {
	return fmt.Sprintf("%T, % X", u, []byte(u))
}

func (u UnknownEvent) Bytes() []byte {
	return []byte(u)
}

func (u UnknownEvent) Raw() []byte {
	panic("unsupported")
}

var _ Event = UnknownEvent(nil)
