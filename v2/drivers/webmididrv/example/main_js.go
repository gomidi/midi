// +build js,wasm

package main

import (
	"bytes"
	"fmt"
	"syscall/js"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/webmididrv"
)

/*
to build, run

GOOS=js GOARCH=wasm go build -o main.wasm main_js.go
*/

func printMessage(message string) {
	document := js.Global().Get("document")
	p := document.Call("createElement", "p")
	p.Set("innerHTML", message)
	document.Get("body").Call("appendChild", p)
}

func main() {
	defer midi.CloseDriver()

	ins := midi.InPorts()
	var bf bytes.Buffer

	for i, in := range ins {
		fmt.Fprintf(&bf, "found MIDI in port: %v: %s<br />", i, in)
	}

	printMessage(bf.String())

	outs := midi.OutPorts()

	bf.Reset()

	for i, out := range outs {
		fmt.Fprintf(&bf, "found MIDI out port: %v: %s<br />", i, out)
	}

	printMessage(bf.String())

	s, err := midi.SenderToPort(0)
	must(err)

	var o midi.ListenOptions

	stop, err := midi.ListenToPort(0, midi.ReceiverFunc(func(msg midi.Message, timestamp int32) {
		printMessage(fmt.Sprintf("got: %s<br />", msg))
	}), o)
	must(err)

	channel := midi.Channel(3)
	key := uint8(60)
	velocity := uint8(120)

	printMessage(fmt.Sprintf("send: NoteOn key: %v veloctiy: %v on channel %v<br />", key, velocity, channel))

	// do some writing: if you are using a loopback midi device on your os, you will see
	// this messages in the browser window
	s.Send(channel.NewNoteOn(key, velocity))
	time.Sleep(time.Second)
	printMessage(fmt.Sprintf("send: NoteOff key: %v on channel %v<br />", key, channel))
	s.Send(channel.NewNoteOff(key))

	// stay alive
	ch := make(chan bool)
	<-ch
	stop()
}

func must(err error) {
	if err != nil {
		printMessage(fmt.Sprintf("ERROR: %s", err.Error()))
	}
}
