package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"

	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"gitlab.com/metakeule/config"
)

var (
	cfg      = config.MustNew("midispy", "2.0.0", "spy on the MIDI data that is sent from a device to another.")
	inArg    = cfg.NewInt32("in", "number of the input device", config.Required, config.Shortflag('i'))
	outArg   = cfg.NewInt32("out", "number of the output device", config.Shortflag('o'))
	noLogArg = cfg.NewBool("nolog", "don't log, just connect in and out", config.Shortflag('n'))
	shortArg = cfg.NewBool("short", "log the short way", config.Shortflag('s'))
	listCmd  = cfg.MustCommand("list", "list devices").Relax("in").Relax("out")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	defer midi.CloseDriver()

	if err := cfg.Run(); err != nil {
		listMIDIInDevices()
		return err
	}

	if cfg.ActiveCommand() == listCmd {
		listMIDIDevices()
		return nil
	}

	err = startSpying(!noLogArg.Get())

	if err != nil {
		return err
	}

	sigchan := make(chan os.Signal, 10)

	// listen for ctrl+c
	go signal.Notify(sigchan, os.Interrupt)

	// interrupt has happend
	<-sigchan
	fmt.Println("\n--interrupted!")

	return nil
}

func listMIDIDevices() {
	listMIDIInDevices()

	fmt.Print("\n--- MIDI output ports ---\n\n")

	for num, port := range midi.OutPorts() {
		fmt.Printf("[%d] %#v\n", num, port)
	}

	return
}

func listMIDIInDevices() {
	fmt.Print("\n--- MIDI input ports ---\n\n")

	for num, port := range midi.InPorts() {
		fmt.Printf("[%d] %#v\n", num, port)
	}
}

func startSpying(shouldlog bool) error {

	in := inArg.Get()

	inPort, err := drivers.InByNumber(int(in))
	if err != nil {
		return err
	}

	err = inPort.Open()

	if err != nil {
		return err
	}

	var outPort drivers.Out = nil
	var logfn func(...interface{})

	if outArg.IsSet() {

		out := outArg.Get()
		outPort, err = drivers.OutByNumber(int(out))
		if err != nil {
			return err
		}

		err = outPort.Open()

		if err != nil {
			return err
		}

		fmt.Printf("[%d] %#v\n->\n[%d] %#v\n-----------------------\n",
			inPort.Number(), inPort.String(), outPort.Number(), outPort.String())
		logfn = logger(in, out)
	} else {
		fmt.Printf("[%d] %#v\n-----------------------\n",
			inPort.Number(), inPort.String())
		logfn = logger(in, 0)
	}

	recv := func(m midi.Message, absmillisec int32) {}

	if shouldlog {
		recv = func(m midi.Message, absmillisec int32) {
			logfn(m)
		}
	}

	return Run(inPort, outPort, recv)
}

func logger(in, out int32) func(...interface{}) {
	if shortArg.Get() {
		return func(v ...interface{}) {
			fmt.Println(v...)
		}
	}
	if outArg.IsSet() {
		l := log.New(os.Stdout, fmt.Sprintf("[%d->%d] ", in, out), log.Lmicroseconds)
		return l.Println
	}

	l := log.New(os.Stdout, fmt.Sprintf("[%d] ", in), log.Lmicroseconds)
	return l.Println
}
