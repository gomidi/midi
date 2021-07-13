package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2"
	//"gitlab.com/gomidi/midi/reader"
	//"gitlab.com/gomidi/midi/writer"
	driver "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

/*
a better way would be:
*/

type receiver struct {
}

func (r *receiver) Receive(msg midi.Message, deltamicro int64) {
	// do something
}

func (r *receiver) ReceiveSysEx(data []byte) {
	// do something
}

func (r *receiver) ReceiveSysCommon(msg midi.Message, deltamicro int64) {
	// do something
}

func (r *receiver) ReceiveRealTime(typ midi.MsgType, deltamicro int64) {
	// do something
}

/*
and then in.SendTo(recv)
*/

var _ midi.Receiver = &receiver{}
var _ midi.SysExReceiver = &receiver{}
var _ midi.SysCommonReceiver = &receiver{}
var _ midi.RealtimeReceiver = &receiver{}

var (
	portsMx sync.Mutex
	drv     midi.Driver

	inPorts  = map[int]midi.In{}
	outPorts = map[int]midi.Out{}
)

func init() {
	var err error
	drv, err = driver.New()
	if err != nil {
		panic("can't initialize driver")
	}
}

func main() {
	// make sure to close all open ports at the end
	defer drv.Close()

	var ww = make(chan int, 10)

	go func() {
		for {
			go checkPorts()
			time.Sleep(time.Second * 1)
		}
	}()

	// interrupt with ctrl+c
	<-ww
}

func greet(out midi.Out) {
	out.Open()
	ch := midi.Channel(0)
	//wr := writer.New(out)
	time.Sleep(time.Millisecond * 200)
	out.Send(ch.NoteOn(60, 100))
	time.Sleep(time.Nanosecond)
	out.Send(ch.NoteOff(60))
	time.Sleep(time.Nanosecond)
	ch = midi.Channel(1)
	out.Send(ch.NoteOn(70, 100))
	time.Sleep(time.Nanosecond)
	out.Send(ch.NoteOff(70))
	time.Sleep(time.Second * 1)
}

func listen(in midi.In) {
	in.Open()
	recv := midi.ReceiverFunc(func(msg midi.Message, deltatime int64) {
		fmt.Printf("got message %s from in port %s\n", msg.String(), in.String())
	})
	in.SendTo(recv)
	//rd.ListenTo(in)
}

func checkPorts() {
	//fmt.Println("...")
	portsMx.Lock()
	ins, _ := drv.Ins()

	for _, in := range ins {
		if strings.Contains(in.String(), "Client") {
			continue
		}
		if inPorts[in.Number()] != nil {
			if inPorts[in.Number()].String() != in.String() {
				inPorts[in.Number()].StopListening()
				inPorts[in.Number()].Close()
				fmt.Printf("closing in port: [%v] %s\n", in.Number(), inPorts[in.Number()].String())
				inPorts[in.Number()] = in
				fmt.Printf("new in port: [%v] %s\n", in.Number(), in.String())
				go listen(in)
			} else {
				continue
			}
		} else {
			inPorts[in.Number()] = in
			fmt.Printf("new in port: [%v] %s\n", in.Number(), in.String())
			go listen(in)
		}
	}

	outs, _ := drv.Outs()

	for _, out := range outs {
		if strings.Contains(out.String(), "Client") {
			continue
		}
		if outPorts[out.Number()] != nil {
			if outPorts[out.Number()].String() != out.String() {
				outPorts[out.Number()].Close()
				fmt.Printf("closing out port: [%v] %s\n", out.Number(), outPorts[out.Number()].String())
				outPorts[out.Number()] = out
				fmt.Printf("new out port: [%v] %s\n", out.Number(), out.String())
				go greet(out)
			} else {
				continue
			}
		} else {
			fmt.Printf("new out port: [%v] %s\n", out.Number(), out.String())
			outPorts[out.Number()] = out
			go greet(out)
		}
	}
	portsMx.Unlock()
}
