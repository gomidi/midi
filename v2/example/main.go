package main

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func main() {

	channel := midi.Channel(2)
	key := uint8(60)
	velocity := uint8(125)

	n1 := channel.NoteOn(key, velocity)
	n1Off := channel.NoteOff(key)
	l := midi.MetaLyric("hello world")

	s := smf.NewSMF1()
	tr := smf.NewTrack()
	tr.Add(0, n1)
	tr.Add(0, l)
	tr.Add(960, n1Off)

	s.AddAndClose(0, tr)

	/*
		s.WriteTo(0, n1, 0)
		s.WriteTo(0, l, 0)
		s.WriteTo(0, n1Off, 960)
	*/
	//s.WriteFile("./test.mid")

	fmt.Printf("%v: %q\n", n1, midi.GetMessageType(n1))
	fmt.Printf("%v: %q\n", n1Off, midi.GetMessageType(n1Off))
	fmt.Printf("%v: %q\n", l, midi.GetMessageType(l))

	if !midi.GetMessageType(n1).IsAllOf(midi.Channel2, midi.ChannelMsg, midi.NoteOnMsg) {
		println("type is invalid")
	}

	var n1m midi.Message
	n1m.Data = n1
	n1m.Type = midi.GetMessageType(n1)

	ch, k, v := n1m.Channel(), n1m.Key(), n1m.Velocity()

	if uint8(channel) != uint8(ch) {
		println("channel does not match ", ch)
	}

	if key != uint8(k) {
		println("key does not match ", k)
	}

	if velocity != uint8(v) {
		println("velocity does not match ", v)
	}

	println("OK")
}
