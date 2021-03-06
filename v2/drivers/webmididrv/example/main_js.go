// +build js,wasm

package main

import (
	"bytes"
	"fmt"
	"syscall/js"
	"time"

	"gitlab.com/gomidi/midi/v2"
	//"gitlab.com/gomidi/midi/midiwriter"
	//"gitlab.com/gomidi/midi/reader"
	//"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/midi/v2/drivers/webmididrv"
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
	drv, err := driver.New()
	must(err)

	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	var bf bytes.Buffer

	for _, in := range ins {
		fmt.Fprintf(&bf, "found MIDI in port: %v: %s<br />", in.Number(), in.String())
	}

	printMessage(bf.String())

	outs, err := drv.Outs()
	must(err)

	bf.Reset()

	for _, out := range outs {
		fmt.Fprintf(&bf, "found MIDI out port: %v: %s<br />", out.Number(), out.String())
	}

	printMessage(bf.String())

	in := ins[0]
	err = in.Open()
	must(err)

	out := outs[0]
	err = out.Open()
	must(err)

	recv := midi.NewReceiver(func(msg midi.Message, deltamicrosecs int64) {
		printMessage(fmt.Sprintf("got: %s<br />", msg))
	}, nil)

	/*
		rd := reader.New(
			reader.NoLogger(),
			reader.Each(func(_ *reader.Position, msg midi.Message) {
				printMessage(fmt.Sprintf("got: %s<br />", msg))
			}),
		)

		rd.ListenTo(in)
	*/
	in.SendTo(recv)

	// Running status is not allowed according to the specs.
	//wr := writer.New(out, midiwriter.NoRunningStatus())

	channel := midi.Channel(3)
	key := uint8(60)
	velocity := uint8(120)

	printMessage(fmt.Sprintf("send: NoteOn key: %v veloctiy: %v on channel %v<br />", key, velocity, channel))

	// do some writing: if you are using a loopback midi device on your os, you will see
	// this messages in the browser window
	//wr.SetChannel(channel)
	//writer.NoteOn(wr, key, velocity)
	out.Send(channel.NoteOn(key, velocity))
	time.Sleep(time.Second)
	printMessage(fmt.Sprintf("send: NoteOff key: %v on channel %v<br />", key, channel))
	//writer.NoteOff(wr, key)
	out.Send(channel.NoteOff(key))

	// stay alive
	ch := make(chan bool)
	<-ch
}

func must(err error) {
	if err != nil {
		printMessage(fmt.Sprintf("ERROR: %s", err.Error()))
	}
}
