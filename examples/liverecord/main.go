package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/meta"
	"gitlab.com/gomidi/midi/midimessage/meta/meter"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/smf"
	"gitlab.com/gomidi/midi/smf/smfwriter"
	"gitlab.com/gomidi/midi/writer"

	// replace with e.g. "gitlab.com/gomidi/rtmididrv" for real midi connections
	driver "gitlab.com/gomidi/midi/testdrv"
)

type timedMsg struct {
	deltaMicrosecs int64
	data           []byte
}

func main() {
	// you would take a real driver here e.g. rtmididrv.New()
	drv := driver.New("fake cables: messages written to output port 0 are received on input port 0")

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs() // just for the test with the fake cables
	must(err)

	inPort := ins[0]   // set the right midi in port
	outPort := outs[0] // just for testing with fake cables
	must(inPort.Open())
	must(outPort.Open())

	defer inPort.Close()
	defer outPort.Close()

	// here comes the meat

	var inbf bytes.Buffer
	var outbf bytes.Buffer

	resolution := smf.MetricTicks(1920)
	bpm := 120.00

	wr := writer.NewSMF(&outbf, 1, smfwriter.TimeFormat(resolution))
	wr.WriteHeader()
	wr.Write(meta.FractionalBPM(bpm)) // set the initial bpm
	wr.Write(meter.M3_4())            // set the meter if needed

	defer func() {
		wr.Write(meta.EndOfTrack)
		ioutil.WriteFile("recorded.mid", outbf.Bytes(), 0644)
	}()

	rd := midireader.New(&inbf, nil)

	ch := make(chan timedMsg)
	stop := make(chan bool)
	bpmCh := make(chan float64) // allows to change bpm on the fly

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			select {

			case bpm = <-bpmCh: // change the bpm
				wr.Write(meta.FractionalBPM(bpm))

			case tm := <-ch:
				deltaticks := resolution.FractionalTicks(bpm, time.Duration(tm.deltaMicrosecs)*time.Microsecond)

				wr.SetDelta(deltaticks)
				inbf.Write(tm.data)
				msg, _ := rd.Read()
				wr.Write(msg)
				go func(m midi.Message, ticks uint32, msec int64) {
					fmt.Printf("msg: %s deltaticks: %v dur: %vµsec\n", m, ticks, msec)
				}(msg, deltaticks, tm.deltaMicrosecs)

			case <-stop:
				wg.Done()
				return
			}
		}
	}()

	inPort.SetListener(func(data []byte, deltaMicrosecs int64) {
		if len(data) == 0 {
			return
		}
		ch <- timedMsg{data: data, deltaMicrosecs: deltaMicrosecs}
	})

	// check with fake cables
	livewr := writer.New(outPort)

	livewr.SetChannel(0) // is midi channel 1
	writer.NoteOn(livewr, 60, 100)
	time.Sleep(1 * time.Second)
	writer.NoteOff(livewr, 60)

	M
	bpmCh <- 60.00                    // change BPM
	time.Sleep(10 * time.Millisecond) // give it a bit time
	writer.NoteOn(livewr, 62, 90)
	time.Sleep(1 * time.Second)
	writer.NoteOff(livewr, 62)

	inPort.StopListening()
	stop <- true
	wg.Wait()
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
