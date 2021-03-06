package midi

type Filter []MessageType

func (f Filter) Sender(target Sender) Sender {
	return &filteringSender{f, target}
}

func (f Filter) Receiver(rec Receiver) Receiver {
	return &filteringReceiver{f, rec}
}

type filteringSender struct {
	Filter
	Sender
}

// Send sends the midi message and silently filters the unwanted out
func (f *filteringSender) Send(m Message) error {
	if m.Type.IsOneOf(f.Filter...) {
		return f.Send(m)
	}
	return nil
}

type filteringReceiver struct {
	Filter
	Receiver
}

// Read reads the midi message and silently filters the unwanted out
func (f *filteringReceiver) Receive(m Message, deltamicrosec int64) {
	if m.Type.IsOneOf(f.Filter...) {
		f.Receiver.Receive(m, deltamicrosec)
	}
	return
}

// Sender sends MIDI messages.
type Sender interface {
	// Send sends the given MIDI message and returns any error.
	Send(Message) error
}

// SenderTo sends MIDI messages.
type SenderTo interface {
	// SendTo sends MIDI messages to the given receiver.
	SendTo(Receiver)
}

// Receiver receives MIDI messages.
type Receiver interface {
	// Receive receives a MIDI message. deltamicrosec is the delta to the previous note in microseconds (^-6)
	Receive(m Message, deltamicrosec int64)
}
