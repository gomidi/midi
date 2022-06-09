package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	"gitlab.com/gomidi/midi/v2/smf"
)

func record(in drivers.In, ticks smf.MetricTicks, bpm float64, send func(msg midi.Message) error) (stop func() smf.Track) {
	var tr smf.Track
	var absmillisec int32

	_stop, err := midi.ListenTo(in, func(msg midi.Message, absms int32) {
		send(msg)
		deltams := absms - absmillisec
		absmillisec = absms
		delta := ticks.Ticks(bpm, time.Duration(deltams)*time.Millisecond)
		tr.Add(delta, msg)
	})

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil
	}

	return func() smf.Track {
		_stop()
		tr.Close(0)
		return tr
	}
}

func main() {
	defer midi.CloseDriver()

	in, err := midi.FindInPort("VMPK")
	if err != nil {
		fmt.Println("can't find in port")
		return
	}
	out, err := midi.FindOutPort("qsynth")
	if err != nil {
		fmt.Println("can't find out port")
		return
	}

	send, _ := midi.SendTo(out)

	var s smf.SMF
	var ticks = smf.MetricTicks(960)
	s.TimeFormat = ticks
	var bpm float64 = 120.00

	var tempoTrack smf.Track
	var tick = midi.NoteOn(9, 60, 120)
	var tickOff = midi.NoteOff(9, 60)
	var tock = midi.NoteOn(9, 50, 100)
	var tockOff = midi.NoteOff(9, 50)

	// 2 bars ticks and tocks on quarter notes
	tempoTrack.Add(0, tick)
	for i := 0; i < 7; i++ {
		if (i+1)%4 == 0 {
			tempoTrack.Add(ticks.Ticks4th(), tockOff)
			tempoTrack.Add(9, tick)
		} else {
			tempoTrack.Add(ticks.Ticks4th(), tickOff)
			tempoTrack.Add(9, tock)
		}
	}

	tempoTrack.Add(ticks.Ticks4th(), tockOff)
	tempoTrack.Close(0)
	s.Add(tempoTrack)

	var bf bytes.Buffer
	s.WriteTo(&bf)

	player := smf.ReadTracksFrom(bytes.NewReader(bf.Bytes()))

	sigchan := make(chan os.Signal, 10)

	// listen for ctrl+c
	go signal.Notify(sigchan, os.Interrupt)

	go func() {
		for {
			stop := record(in, ticks, bpm, send)
			player.Play(out)
			rec := stop()
			if !rec.IsEmpty() {
				s.Add(rec)
				print(".")
			}
			send(midi.ControlChange(0, midi.AllNotesOff, midi.On))
			send(midi.ControlChange(9, midi.AllNotesOff, midi.On))
			bf.Reset()
			s.WriteTo(&bf)
			player = smf.ReadTracksFrom(bytes.NewReader(bf.Bytes()))
		}
	}()

	// interrupt has happend
	<-sigchan
	fmt.Println("\n--interrupted!")
	send(midi.ControlChange(0, midi.AllNotesOff, midi.On))
	send(midi.ControlChange(9, midi.AllNotesOff, midi.On))
	s.WriteFile("recorded.mid")
}
