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

	"gitlab.com/golang-utils/config/v2"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	lib "gitlab.com/gomidi/midi/v2/drivers/midicat"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

var (
	cfg = config.New("midicat", VERSION, "midicat transfers MIDI data between midi ports and stdin/stdout")

	argPortNum  = cfg.Int("index", "index of the midi port. Only specify either the index or the name. If neither is given, the first port is used.", config.Shortflag('i'))
	argPortName = cfg.String("name", "name of the midi port. Only specify either the index or the name. If neither is given, the first port is used.")
	argJson     = cfg.Bool("json", "return the list in JSON format")

	cmdIn  = cfg.Command("in", "read midi from an in port and print it to stdout").Skip("json")
	cmdOut = cfg.Command("out", "read midi from stdin and print it to an out port").Skip("json")

	cmdIns  = cfg.Command("ins", "show the available midi in ports").SkipAllBut("json")
	cmdOuts = cfg.Command("outs", "show the available midi out ports").SkipAllBut("json")

	cmdLog      = cfg.Command("log", "pass the midi from stdin to stdout while logging it to stderr").SkipAllBut()
	argLogNoOut = cmdLog.Bool("nopass", "don't pass the midi to stdout")

	shouldStop = make(chan bool, 1)
	didStop    = make(chan bool, 1)
)

func main() {
	defer midi.CloseDriver()

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

	/*
		drv, err := rtmididrv.New()
		if err != nil {
			return err
		}
	*/

	switch cfg.ActiveCommand() {
	case cmdIns:
		if argJson.Get() {
			return showInJson()
		} else {
			return showInPorts()
		}
	case cmdOuts:
		if argJson.Get() {
			return showOutJson()
		} else {
			return showOutPorts()
		}
	case cmdIn:
		return runIn()
	case cmdOut:
		return runOut()
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

		b, abstime, err := lib.ReadAndConvert(os.Stdin)
		if err == io.EOF {
			break
		}

		if err != nil {
			logMsg("midicat log: could not read from stdin: %s\n", string(b), err.Error())
		}

		if !argLogNoOut.Get() {
			_, werr := fmt.Fprintf(os.Stdout, "%d %X\n", abstime, b)
			os.Stdout.Sync()
			_ = werr
		}

		//logBuffer.Write(b)
		//msg, merr := logrd.Read()
		msg := midi.Message(b)
		//logBuffer.Reset()

		/*
			if merr != nil {
				logMsg("midicat log: could understand % X: %s\n", b, merr.Error())
			} else {
				//logMsg("%s\n", msg)
			}
		*/
		fmt.Fprintf(os.Stderr, "%vms %s # ", abstime, msg.String())
		runtime.Gosched()
	}
	return nil
}

type timestampedMsg struct {
	absmillisec int32
	msg         []byte
}

func runIn() (err error) {
	var in drivers.In

	switch {
	case argPortNum.IsSet():
		in, err = midi.InPort(int(argPortNum.Get()))
	case argPortName.IsSet():
		in, err = midi.FindInPort(argPortName.Get())
	default:
		in, err = midi.InPort(0)
	}

	if err != nil {
		return
	}

	var msgChan = make(chan timestampedMsg, 1)
	var stopChan = make(chan bool, 1)
	var stoppedChan = make(chan bool, 1)

	recv := func(msg midi.Message, absmillisec int32) {
		//fmt.Printf("got message %s from in port %v\n", msg.String(), in)

		msgChan <- timestampedMsg{absmillisec: absmillisec, msg: msg.Bytes()}
	}

	go func() {
		for {
			select {
			case m := <-msgChan:
				_, werr := fmt.Fprintf(os.Stdout, "%d %X\n", m.absmillisec, m.msg)
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

	var stop func()

	go func() {
		stop, err = midi.ListenTo(in, recv, midi.UseActiveSense(), midi.UseSysEx(), midi.UseTimeCode())

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
	stop()
	//in.StopListening()
	stopChan <- true
	<-stoppedChan

	return nil
}

func runOut() (err error) {
	var out drivers.Out

	switch {
	case argPortNum.IsSet():
		out, err = midi.OutPort(int(argPortNum.Get()))
	case argPortName.IsSet():
		out, err = midi.FindOutPort(argPortName.Get())
	default:
		out, err = midi.OutPort(0)
	}

	if err != nil {
		return err
	}

	send, err := midi.SendTo(out)

	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		//var lastAbstime int32
		for {
			//b, abstime, err := lib.ReadAndConvert(os.Stdin)
			b, _, err := lib.ReadAndConvert(os.Stdin)

			if err == io.EOF {
				break
			}

			if err != nil {
				logMsg("midicat out: error %s\n", err.Error())
				continue
			}

			/*
				delta := abstime - lastAbstime

				if delta > 0 {
					time.Sleep(time.Millisecond * time.Duration(delta))
				}

				lastAbstime = abstime
			*/
			werr := send(b)

			if werr != nil {
				logMsg("midicat out: could not write % X to port %q: %s\n", b, out, werr.Error())
			}
			runtime.Gosched()
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}

func showInJson() error {
	var portm = map[int]string{}

	for num, port := range midi.GetInPorts() {
		portm[num] = port.String()
	}

	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(portm)
}

func showInPorts() error {
	fmt.Fprintln(os.Stdout, "MIDI inputs")

	for num, in := range midi.GetInPorts() {
		fmt.Fprintf(os.Stdout, "[%v] %s\n", num, in.String())
	}

	return nil
}

func showOutJson() error {
	var portm = map[int]string{}

	for num, port := range midi.GetOutPorts() {
		portm[num] = port.String()
	}

	enc := json.NewEncoder(os.Stdout)
	return enc.Encode(portm)
}

func showOutPorts() error {
	fmt.Fprintln(os.Stdout, "MIDI outputs")

	for num, out := range midi.GetOutPorts() {
		fmt.Fprintf(os.Stdout, "[%v] %s\n", num, out.String())
	}

	return nil
}
