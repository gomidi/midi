package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"

	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// in order to receive sysex messages or get errors, you need to define a receiver type
// otherwise a midi.ReceiverFunc is sufficient
type receiver struct{}

func (r receiver) Receive(msg midi.Message, timestamp int32) {
	fmt.Printf("got %s @%v\n", msg, timestamp)
}

var _ midi.Receiver = receiver{}

// To receive sysex messages, implement the midi.SysExReceiver interface
func (r receiver) OnSysEx(b []byte, timestamp int32) {
	fmt.Printf("got sysex: '% X' @%v\n", b, timestamp)
}

var _ midi.SysExReceiver = receiver{}

// To receive errors, implement the midi.ErrorReceiver interface
func (r receiver) OnError(err error) {
	if err == midi.ErrListenStopped {
		fmt.Println("stopped listening")
	} else {
		fmt.Printf("error: %s\n", err.Error())
	}
}

var _ midi.ErrorReceiver = receiver{}

func main() {
	run()
	os.Exit(0)
}

func run() {

	// always close the driver at the end
	defer midi.CloseDriver()

	if len(os.Args) == 2 && os.Args[1] == "list" {
		fmt.Printf("MIDI IN Ports\n")
		printPorts(midi.InPorts())
		fmt.Printf("\n\n")
		fmt.Printf("MIDI OUT Ports\n")
		printPorts(midi.OutPorts())
		return
	}

	// to get the port number via name, use midi.FindOutPort
	s, err := midi.SenderToPort(0)
	must(err)

	// to get the port number via name, use midi.FindInPort
	stop, err := midi.ListenToPort(0, receiver{}, midi.ListenOptions{ActiveSense: true, TimeCode: true})
	must(err)

	time.Sleep(time.Millisecond)
	s.Send(midi.Channel(2).NoteOn(12, 34))
	time.Sleep(time.Millisecond)
	s.Send(midi.Activesense())
	time.Sleep(time.Millisecond)
	s.Send(midi.Channel(2).NoteOff(12))
	time.Sleep(time.Millisecond)
	s.Send(midi.Tune())
	time.Sleep(time.Millisecond)
	s.Send(midi.Activesense())
	time.Sleep(time.Millisecond)
	s.Send(midi.SysEx([]byte{0x41, 0x10, 0x42, 0x12, 0x40, 0x00, 0x7F, 0x00, 0x41}))

	time.Sleep(time.Millisecond)
	stop()
}

func printPorts(ports []string) {
	for i, port := range ports {
		fmt.Printf("[%v] %s\n", i, port)
	}
}
