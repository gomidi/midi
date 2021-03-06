package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func main() {

	ch := midi.Channel(2)
	key := uint8(60)
	velocity := uint8(125)

	n1 := ch.NoteOn(key, velocity)
	n1Off := ch.NoteOff(key)
	l := midi.MetaLyric("hello world")

	s := smf.NewSMF1()
	tr := smf.NewTrack()
	tr.Add(0, n1)
	tr.Add(0, l)
	tr.Add(960, n1Off)

	s.AddAndClose(0, tr)
	s.WriteFile("./test.mid")

	fmt.Printf("%v: %q\n", n1, midi.GetMessageType(n1))
	fmt.Printf("%v: %q\n", n1Off, midi.GetMessageType(n1Off))
	fmt.Printf("%q: %q\n", midi.NewMessage(l).Text(), midi.GetMessageType(l))

	if !midi.GetMessageType(n1).IsAllOf(midi.Channel2Msg, midi.ChannelMsg, midi.NoteOnMsg) {
		println("type is invalid")
	}

	var n1m midi.Message
	n1m = midi.NewMessage(n1)
	n1m.Type = midi.GetMessageType(n1)

	c, k, v := n1m.Channel(), n1m.Key(), n1m.Velocity()

	if uint8(ch) != uint8(c) {
		println("channel does not match ", c)
	}

	if key != uint8(k) {
		println("key does not match ", k)
	}

	if velocity != uint8(v) {
		println("velocity does not match ", v)
	}

	println("OK")
}
