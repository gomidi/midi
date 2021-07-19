package main

import (
	"fmt"
	"os"
	"time"

	"gitlab.com/gomidi/midi/v2"

	// include a driver (autoregisters it)
	//_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	_ "gitlab.com/gomidi/midi/v2/drivers/testdrv"
	//_ "gitlab.com/gomidi/midi/v2/drivers/midicatdrv"
	//_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// in order to receive other midi messages than channel messages, you need to defined a receiver type
// for only channel messages, a midi.ReceiverFunc is sufficient
type receiver struct{}

func (r receiver) Receive(msg midi.Message, timestamp int32) {
	fmt.Printf("got %s @%v\n", msg, timestamp)
}

var _ midi.Receiver = receiver{}

// To receive sysex messages, implement the midi.SysExReceiver interface
func (r receiver) ReceiveSysEx(b []byte, timestamp int32) {
	fmt.Printf("got sysex: % X @%v\n", b, timestamp)
}

var _ midi.SysExReceiver = receiver{}

// To receive sys common messages, implement the midi.SysCommonReceiver interface
func (r receiver) ReceiveSysCommon(msg midi.Message, timestamp int32) {
	fmt.Printf("got syscommon: %s @%v\n", msg, timestamp)
}

var _ midi.SysCommonReceiver = receiver{}

// To receive realtime messages, implement the midi.RealtimeReceiver interface
func (r receiver) ReceiveRealtime(mtype midi.MsgType, timestamp int32) {
	fmt.Printf("got realtime: %s @%v\n", mtype, timestamp)
}

var _ midi.RealtimeReceiver = receiver{}

// run this in two terminals. first terminal without args to create the virtual ports and
// second terminal with argument "list" to see the ports.
func main() {
	run()
	os.Exit(0)
}
func run() {

	// always close the driver at the end
	defer midi.CloseDriver()

	if len(os.Args) == 2 && os.Args[1] == "list" {
		fmt.Printf("MIDI IN Ports\n")
		fmt.Printf("\n\n")
		printPorts(midi.InPorts())
		fmt.Printf("MIDI OUT Ports\n")
		printPorts(midi.OutPorts())
		return
	}

	//in, err = drv.OpenVirtualIn("test-virtual-in")
	//sender, err := midi.SenderToPort(midi.FindOutPort("Through Port-0"))

	// to get the port number via name, use midi.FindOutPort
	s, err := midi.SenderToPort(0)
	must(err)

	//out, err = drv.OpenVirtualOut("test-virtual-out")
	//err = midi.ListenToPort(midi.FindInPort("Through Port-0"), receiver{})

	// to get the port number via name, use midi.FindInPort
	err = midi.ListenToPort(0, receiver{})
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
	s.Send(midi.SysEx([]byte{0x41, 0x10, 0x42, 0x12, 0x40, 0x00, 0x7F, 0x00, 0x41}))
	time.Sleep(time.Millisecond)
	s.Send(midi.Activesense())
	time.Sleep(time.Millisecond)
	s.Send(midi.Channel(2).NoteOff(12))
	// F0   41   10   42   12   40007F   00   41   F7

	time.Sleep(time.Second)
}

func printPorts(ports []string) {
	for i, port := range ports {
		fmt.Printf("[%v] %s\n", i, port)
	}
}
