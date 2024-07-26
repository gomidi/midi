package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"gitlab.com/golang-utils/config/v2"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
	"gitlab.com/gomidi/midi/v2/sysex"
)

var (
	cfg = config.New("smfsysex", 0, 0, 3, "a tool to read and write sysex data from / to a SMF file",
		config.AsciiArt("smfsysex"))

	argFile = cfg.LastString("file", "the file that is written to / read from", config.Required())
	argRaw  = cfg.Bool("raw", "don't interpret the sysex data as hex string, but use the raw bytes instead", config.Shortflag('r'))

	cmdWrite           = cfg.Command("write", "writes sysex data to a midi file")
	argWriteSysex      = cmdWrite.String("sysex", "the sysex data to be written", config.Required(), config.Shortflag('x'))
	argWriteTrackname  = cmdWrite.String("trackname", "the name of the track", config.Shortflag('n'))
	argWriteInstrument = cmdWrite.String("instrument", "the name of the instrument", config.Shortflag('i'))
	argWriteAdd        = cmdWrite.Bool("add", "add track to existing file, instead of writing a new file", config.Shortflag('a'))

	cmdRead        = cfg.Command("read", "reads sysex data")
	argReadVerbose = cmdRead.Bool("verbose", "be verbose about printing", config.Shortflag('v'))
	argTrackRead   = cmdRead.Int("track", "number of the track to read from", config.Required(), config.Shortflag('t'), config.Default(1))

	cmdExample = cfg.Command("example", "write example sysex data")
)

func main() {

	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	err := cfg.Run()

	if err != nil {
		//fmt.Fprintf(os.Stdout, cfg.Usage()+"\n")
		return err
	}

	switch cfg.ActiveCommand() {
	case cmdWrite:
		return writeSysEx(
			argFile.Get(),
			argWriteSysex.Get(),
			argRaw.Get(),
			argWriteTrackname.Get(),
			argWriteInstrument.Get(),
			argWriteAdd.Get(),
		)
	case cmdRead:
		trno := int(argTrackRead.Get()) - 1
		return readSysEx(argFile.Get(), trno, argRaw.Get(), argReadVerbose.Get())
	case cmdExample:
		return writeExample(argFile.Get())
	default:
		return showTracks(argFile.Get())
	}
}

func parseSysEx(s string) ([]byte, error) {
	s = strings.ReplaceAll(s, " ", "")
	data, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// TODO allow it to read from stdin as well, dann muss der sysex string natürlich nicht gelesen werden
func writeSysEx(file string, s string, raw bool, tname string, instr string, add bool) (err error) {
	var bt []byte

	if raw {
		bt = []byte(s)
	} else {
		bt, err = parseSysEx(s)
		if err != nil {
			return err
		}
	}

	var mf *smf.SMF

	if add {
		mf, err = smf.ReadFile(file)
		if err != nil {
			return err
		}
	} else {
		mf = smf.NewSMF1()
	}

	var t smf.Track
	if tname != "" {
		t.Add(0, smf.MetaTrackSequenceName(tname))
	}
	if instr != "" {
		t.Add(0, smf.MetaInstrument(instr))
	}
	t.Add(0, midi.SysEx(bt))
	t.Close(0)
	err = mf.Add(t)
	if err != nil {
		return err
	}
	return mf.WriteFile(file)
}

func writeExample(file string) error {
	/*
		if len(bt) == 0 {
			return fmt.Errorf("no bytes given")
		}
	*/

	var mf = smf.NewSMF1()
	var t1 smf.Track
	t1.Add(0, smf.MetaTrackSequenceName("identity request"))
	t1.Add(0, smf.MetaInstrument("standard"))
	t1.Add(0, midi.SysEx(sysex.IdentityRequest(1)))
	t1.Close(0)
	err := mf.Add(t1)
	if err != nil {
		return err
	}
	var t2 smf.Track
	t2.Add(0, smf.MetaTrackSequenceName("GM system on"))
	t2.Add(0, smf.MetaInstrument("GM standard instrument"))
	t2.Add(0, midi.SysEx(sysex.GMSystem(1, true)))
	t2.Close(0)
	err = mf.Add(t2)
	if err != nil {
		return err
	}
	return mf.WriteFile(file)
}

func readSysEx(file string, trackNo int, raw bool, verbose bool) error {

	mf, err := smf.ReadFile(file)

	if err != nil {
		return err
	}

	for i, tr := range mf.Tracks {
		if i != trackNo {
			continue
		}
		var sysx []byte
		var trackname string
		var instrname string

		for _, ev := range tr {
			var bt []byte

			switch {
			case ev.Message.GetMetaTrackName(&trackname):
			case ev.Message.GetMetaInstrument(&instrname):
			case ev.Message.GetSysEx(&bt):
				sysx = append(sysx, bt...)
			}
		}

		if verbose {
			fmt.Fprintf(os.Stdout, "track number: %v\n", i+1)
			if trackname != "" {
				fmt.Fprintf(os.Stdout, "name: %s\n", trackname)
			}
			if instrname != "" {
				fmt.Fprintf(os.Stdout, "instrument: %s\n", instrname)
			}
		}

		if len(sysx) > 0 {
			if raw {
				fmt.Fprintf(os.Stdout, "%s", sysx)
			} else {
				fmt.Fprintf(os.Stdout, "% X", sysx)
			}
			//fmt.Fprintf(os.Stdout, "%s", sysx)
		}
	}

	return nil
}

func showTracks(file string) error {

	mf, err := smf.ReadFile(file)

	if err != nil {
		return err
	}

	var tracksWithSysex = map[int][2]string{}

	for i, tr := range mf.Tracks {
		var sysx bool
		var trackname string
		var instrname string
		for _, ev := range tr {
			var bt []byte
			switch {
			case ev.Message.GetMetaTrackName(&trackname):
			case ev.Message.GetMetaInstrument(&instrname):
			case ev.Message.GetSysEx(&bt):
				sysx = true
			}
		}

		if sysx {
			tracksWithSysex[i] = [2]string{trackname, instrname}
		}
	}

	fmt.Fprintf(os.Stdout, "The following tracks have sysex data:\n")

	//for k, v := range tracksWithSysex {
	for k := 0; k < len(tracksWithSysex); k++ {
		v := tracksWithSysex[k]
		var nme string
		switch {
		case v[0] != "" && v[1] != "":
			nme = v[0] + " / " + v[1]
		case v[0] != "":
			nme = v[0]
		case v[1] != "":
			nme = v[1]
		}
		fmt.Fprintf(os.Stdout, "[%v] %s\n", k+1, nme)
	}

	return nil
}
