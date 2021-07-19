package midi

import (
	"fmt"
	"os"

	"gitlab.com/gomidi/midi/v2/drivers"
)

func CloseDriver() {
	drivers.Close()
}

// Sender sends MIDI messages.
type Sender interface {
	// Send sends the given MIDI message and returns any error.
	Send(msg Message) error
}

type ReceiverFunc func(msg Message, absdecimillisec int32)

func (r ReceiverFunc) Receive(msg Message, absdecimillisec int32) {
	r(msg, absdecimillisec)
}

// Receiver receives MIDI messages.
type Receiver interface {
	// Receive receives a MIDI message. deltamicrosec is the delta to the previous note in microseconds (^-4)
	// max int32  = 2147483647 / 10 (ms) / 1000 (sec) / 60 (min) / 60 (hours) / 24 = 2,48 days which is long enough IMHO for a midi recording
	Receive(msg Message, absdecimillisec int32)

	// println(big.NewRat(math.MaxInt64,1000 /* milliseconds */ *1000 /* seconds */ *60 /* minutes */ *60 /* hours */ *24 /* days */ *365 /* years */).FloatString(0))
	// output: 292471
	// => a ascending timestamp based on microseconds would wrap after 292471 years
	// so absolute timestamp should be preferred

	/*
		I would prefer decimillisecs (dmsec) of absolute time  (10^-4 secs) with uint32:

		max uint32 = 4294967295 / 10 (ms) / 1000 (sec) / 60 (min) / 60 (hours) / 24 = 4,9 days which is long enough IMHO for a midi recording

		max int32  = 2147483647 / 10 (ms) / 1000 (sec) / 60 (min) / 60 (hours) / 24 = 2,48 days which is long enough IMHO for a midi recording
	*/
}

type SysExReceiver interface {
	Receiver
	ReceiveSysEx(data []byte, absdecimillisec int32)
}

type RealtimeReceiver interface {
	Receiver
	ReceiveRealtime(typ MsgType, absdecimillisec int32)
}

type SysCommonReceiver interface {
	Receiver
	ReceiveSysCommon(msg Message, absdecimillisec int32)
}

func InPorts() []string {
	ins, err := drivers.Ins()

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get midi in ports: %s\n", err.Error())
		return nil
	}

	res := make([]string, len(ins))

	for i, in := range ins {
		res[i] = in.String()
	}

	return res
}

func OutPorts() []string {
	outs, err := drivers.Outs()

	if err != nil {
		fmt.Fprintf(os.Stderr, "can't get midi out ports: %s\n", err.Error())
		return nil
	}

	res := make([]string, len(outs))

	for i, out := range outs {
		res[i] = out.String()
	}

	return res
}
