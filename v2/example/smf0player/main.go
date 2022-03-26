package main

import (
	"fmt"
	"sort"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"

	// _ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregisters driver
	_ "gitlab.com/gomidi/midi/v2/drivers/portmididrv" // autoregisters driver
	"gitlab.com/gomidi/midi/v2/smf"
)

type playEvent struct {
	absTime int64
	sleep   time.Duration
	data    []byte
	//bytes   []byte
	out     drivers.Out
	trackNo int
	str     string
}

type player []playEvent

func (p player) Swap(a, b int) {
	p[a], p[b] = p[b], p[a]
}

func (p player) Less(a, b int) bool {
	return p[a].absTime < p[b].absTime
}

func (p player) Len() int {
	return len(p)
}

func run() error {
	/*
		outs := midi.OutPorts()
		for _, o := range outs {
			fmt.Printf("out: %s\n", o)
		}
	*/
	//out, err := drivers.OutByName("FLUID Synth")
	out, err := drivers.OutByName("qsynth")
	if err != nil {
		return err
	}

	defer out.Close()

	var pl player
	var absTime time.Duration = 0
	var currentTrack int

	// single track playing
	// for multitrack we would have to collect the tracks events first and properly synchronize playback
	//_, err = smf.ReadTracks("Prelude4.mid", 2).
	_, err = smf.ReadTracks("Prelude4.mid", 1, 2, 3, 4, 5, 6, 7).
		//_, err = smf.ReadTracks("VOYAGER.MID", 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20).
		//Only(midi.NoteOnMsg, midi.NoteOffMsg).
		//Only(midi.NoteOnMsg, midi.NoteOffMsg, midi.MetaMsgType).
		//Only(midi.NoteMsg, midi.ControlChangeMsg, midi.ProgramChangeMsg).
		//Only(midi.NoteMsg, midi.ControlChangeMsg, midi.ProgramChangeMsg).
		//Only(midi.MetaMsg).
		// TODO: respect tempo changes
		Do(
			func(trackNo int, msg midi.Message, delta int64, deltamicroSec int64) {
				if currentTrack != trackNo {
					absTime = 0
					currentTrack++
				}
				if msg.Type().Kind() == midi.MetaMsg {
					fmt.Printf("[%v] %s (%s)\n", trackNo, msg.String(), msg.Type())
				}
				if mm, ok := msg.(midi.Msg); ok {
					//time.Sleep(time.Microsecond * time.Duration(deltamicroSec))

					absTime = absTime + (time.Microsecond * time.Duration(deltamicroSec))
					pl = append(pl, playEvent{
						absTime: absTime.Microseconds(),
						data:    mm.Data,
						out:     out,
						trackNo: trackNo,
						str:     msg.String(),
						//						bytes:   msg.Bytes(),
					})
					_ = mm
					//out.Send(mm.Data)
				}
			},
		)

	sort.Sort(pl)

	if err != nil {
		return err
	}

	var last time.Duration = 0

	for i, _ := range pl {
		last = play(last, pl[i])
	}

	return nil
}

func init() {
	//m := midi.NewMsg(0xC4, 0x31, 0x00)
	//println(m.String())
}

func play(last time.Duration, p playEvent) time.Duration {
	current := (time.Microsecond * time.Duration(p.absTime))
	diff := current - last
	//fmt.Printf("sleeping %s\n", diff)
	time.Sleep(diff)
	//fmt.Printf("[%v] %q % X\n", p.trackNo, p.str, p.data)
	fmt.Printf("[%v] %q\n", p.trackNo, p.str)
	p.out.Send(p.data)
	return current
}

func main() {
	err := run()

	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
