package midi

import (
	"gitlab.com/gomidi/midi/v2/drivers"
)

func newWriter(out drivers.Out) *writer {
	return &writer{
		out: out,
	}
}

type writer struct {
	out drivers.Out
}

func (w *writer) Send(msg Message) error {
	return w.out.Send(msg.Data)
}

var _ Sender = &writer{}

func SenderToPort(portnumber int) (Sender, error) {
	out, err := drivers.OutByNumber(portnumber)
	if err != nil {
		return nil, err
	}
	return newWriter(out), nil
}
