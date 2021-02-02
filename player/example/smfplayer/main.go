package main

import (
	"fmt"
	"os"
	"os/signal"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/player"

	// driver "gitlab.com/gomidi/rtmididrv"
	driver "gitlab.com/gomidi/midicatdrv"
	"gitlab.com/metakeule/config"
)

var (
	cfg      = config.MustNew("smfplayer", "0.0.1", "a simple SMF player")
	fileArg  = cfg.LastString("file", "MIDI file that should be played", config.Required)
	outArg   = cfg.NewInt32("out", "number of the MIDI output port", config.Shortflag('o'), config.Required)
	listCmd  = cfg.MustCommand("list", "list MIDI out ports").Relax("out").Relax("file")
	sigchan  = make(chan os.Signal, 10)
	finished = make(chan bool, 1)
	stop     = make(chan bool, 1)
)

func abort(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %q", err.Error())
		os.Exit(1)
	}
}

func main() {
	err := cfg.Run()

	if err != nil {
		fmt.Fprintln(os.Stdout, cfg.Usage())
		abort(err)
	}

	drv, err := driver.New()
	abort(err)

	defer drv.Close()

	if cfg.ActiveCommand() == listCmd {
		err = printMIDIPorts(drv)
		abort(err)
		return
	}

	outs, err := drv.Outs()
	abort(err)

	pl, err := player.SMF(fileArg.Get())
	abort(err)

	// listen for ctrl+c
	go signal.Notify(sigchan, os.Interrupt)

	go func() {
		// interrupt has happend
		<-sigchan
		fmt.Println("\n--interrupted!")
		// stop the playing, triggered via ctrl+c
		stop <- true
	}()

	out := outs[int(outArg.Get())]
	err = out.Open()
	abort(err)

	fmt.Fprintf(os.Stdout, "using MIDI out port %q\n", out)

	pl.PlayAll(out, stop, finished)

	<-finished
	// now playing is done
}

func printMIDIPorts(drv midi.Driver) error {
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
