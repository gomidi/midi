/*

midicat is a program that transfers MIDI data between midi ports and stdin/stdout.
The idea is, that you can have midi libraries that do not depend on c (or CGO in the case of go)
and still might want to use some midi to ports. But maybe it is just an option that is not
used much and we don't want to bother the other users with a c/CGO dependency.

example

midicat in -i=10 | midicat log | midicat out -i=11

(routes midi from midi in port 10 to midi out port 11 while logging the parsed messages in readable way to stderr)

*/

package main

import (
	//	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"gitlab.com/gomidi/midi/v2/tools/midicat/lib"
	"gitlab.com/metakeule/config"
)

var (
	cfg = config.MustNew("midicat", VERSION, "midicat transfers MIDI data between midi ports and stdin/stdout")

	argPortNum  = cfg.NewInt32("index", "index of the midi port. Only specify either the index or the name. If neither is given, the first port is used.", config.Shortflag('i'))
	argPortName = cfg.NewString("name", "name of the midi port. Only specify either the index or the name. If neither is given, the first port is used.")
	argJson     = cfg.NewBool("json", "return the list in JSON format")

	cmdIn  = cfg.MustCommand("in", "read midi from an in port and print it to stdout").Skip("json")
	cmdOut = cfg.MustCommand("out", "read midi from stdin and print it to an out port").Skip("json")

	cmdIns  = cfg.MustCommand("ins", "show the available midi in ports").SkipAllBut("json")
	cmdOuts = cfg.MustCommand("outs", "show the available midi out ports").SkipAllBut("json")

	cmdLog      = cfg.MustCommand("log", "pass the midi from stdin to stdout while logging it to stderr").SkipAllBut()
	argLogNoOut = cmdLog.NewBool("nopass", "don't pass the midi to stdout")

	shouldStop = make(chan bool, 1)
	didStop    = make(chan bool, 1)
)

func main() {
	err := run()

	if err != nil {
		fmt.Fprintln(os.Stderr, cfg.Usage())
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	runtime.GOMAXPROCS(1)
	err := cfg.Run()

	if err != nil {
		return err
	}

	if cfg.ActiveCommand() == cmdLog {
		return log()
	}

	drv, err := rtmididrv.New()
	if err != nil {
		return err
	}

	switch cfg.ActiveCommand() {
	case cmdIns:
		if argJson.Get() {
			return showInJson(drv)
		} else {
			return showInPorts(drv)
		}
	case cmdOuts:
		if argJson.Get() {
			return showOutJson(drv)
		} else {
			return showOutPorts(drv)
		}
	case cmdIn:
		return runIn(drv)
	case cmdOut:
		return runOut(drv)
	default:
		return fmt.Errorf("[command] missing")
	}
}

func logRealTime(rt midi.Message) {
	fmt.Fprintf(os.Stderr, "%s\n", rt)
}

func logMsg(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s\n", args...)
}

func log() error {
	//var logBuffer bytes.Buffer
	//logrd := midireader.New(&logBuffer, logRealTime)
	for {

		b, err := lib.ReadAndConvert(os.Stdin)
		if err == io.EOF {
			break
		}

		if err != nil {
			logMsg("midicat log: could not read from stdin: %s\n", string(b), err.Error())
		}

		if !argLogNoOut.Get() {
			_, werr := fmt.Fprintf(os.Stdout, "%X\n", b)
			os.Stdout.Sync()
			_ = werr
		}

		//logBuffer.Write(b)
		//msg, merr := logrd.Read()
		msg := midi.NewMessage(b)
		//logBuffer.Reset()

		/*
			if merr != nil {
				logMsg("midicat log: could understand % X: %s\n", b, merr.Error())
			} else {
				//logMsg("%s\n", msg)
			}
		*/
		fmt.Fprintln(os.Stderr, msg.String())
		runtime.Gosched()
	}
	return nil
}

func runIn(drv drivers.Driver) (err error) {
	defer drv.Close()
	var in drivers.In

	switch {
	case argPortNum.IsSet():
		in, err = drivers.InByNumber(int(argPortNum.Get()))
	case argPortName.IsSet():
		in, err = drivers.InByName(argPortName.Get())
	default:
		in, err = drivers.InByNumber(0)
	}

	if err != nil {
		return err
	}

	err = in.Open()

	if err != nil {
		return err
	}

	var msgChan = make(chan []byte, 1)
	var stopChan = make(chan bool, 1)
	var stoppedChan = make(chan bool, 1)

	recv := midi.ReceiverFunc(func(msg midi.Message, absdecimillisec int32) {
		//fmt.Printf("got message %s from in port %s\n", msg.String(), in.String())
		msgChan <- msg.Data
	})

	go func() {
		for {
			select {
			case msg := <-msgChan:
				_, werr := fmt.Fprintf(os.Stdout, "%X\n", msg)
				if werr != nil {
					logMsg("midicat in: error while writing: %s\n", werr.Error())
				}
				os.Stdout.Sync()

			case <-stopChan:
				stoppedChan <- true
				return
			}
		}
	}()

	//var stop func()

	go func() {
		err = midi.ListenToPort(in.Number(), recv)
		//err = in.SendTo(recv)

		if err != nil {
			stopChan <- true
			<-stoppedChan
			logMsg("midicat in: could not start listener %s\n", err.Error())
		}
	}()

	sigchan := make(chan os.Signal, 10)

	// listen for ctrl+c
	go signal.Notify(sigchan, os.Interrupt)

	// interrupt has happend
	<-sigchan
	in.StopListening()
	stopChan <- true
	<-stoppedChan

	return nil
}

func runOut(drv drivers.Driver) (err error) {
	defer drv.Close()

	var out drivers.Out

	switch {
	case argPortNum.IsSet():
		out, err = drivers.OutByNumber(int(argPortNum.Get()))
	case argPortName.IsSet():
		out, err = drivers.OutByName(argPortName.Get())
	default:
		out, err = drivers.OutByNumber(0)
	}

	if err != nil {
		return err
	}

	err = out.Open()

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for {
			b, err := lib.ReadAndConvert(os.Stdin)

			if err == io.EOF {
				break
			}

			if err != nil {
				logMsg("midicat out: error %s\n", err.Error())
				continue
			}

			werr := out.Send(b)

			if werr != nil {
				logMsg("midicat out: could not write % X to port %q: %s\n", b, out.String(), werr.Error())
			}
			runtime.Gosched()
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func showInJson(drv drivers.Driver) error {
	defer drv.Close()
	ports, err := drv.Ins()

	if err != nil {
		return err
	}

	var portm = map[int]string{}

	for _, port := range ports {
		portm[port.Number()] = port.String()
	}

	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(portm)
}

func showInPorts(drv drivers.Driver) error {
	defer drv.Close()
	ins, err := drv.Ins()

	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, "MIDI inputs")

	for _, in := range ins {
		fmt.Fprintf(os.Stdout, "[%v] %s\n", in.Number(), in.String())
	}

	return nil
}

func showOutJson(drv drivers.Driver) error {
	defer drv.Close()
	ports, err := drv.Outs()

	if err != nil {
		return err
	}

	var portm = map[int]string{}

	for _, port := range ports {
		portm[port.Number()] = port.String()
	}

	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(portm)
}

func showOutPorts(drv drivers.Driver) error {
	defer drv.Close()
	outs, err := drv.Outs()

	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, "MIDI outputs")

	for _, out := range outs {
		fmt.Fprintf(os.Stdout, "[%v] %s\n", out.Number(), out.String())
	}

	return nil
}
