package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gitlab.com/golang-utils/config/v2"
	"gitlab.com/gomidi/midi/tools/smftrack"
)

var (
	cfg       = config.New("smftool", 0, 0, 2, "tool for SMF MIDI file manipulation", config.AsciiArt("smftool"))
	argIn     = cfg.String("in", "SMF file that is the input", config.Required(), config.Shortflag('i'))
	argOut    = cfg.String("out", "name of the resulting SMF file", config.Required(), config.Shortflag('o'))
	cmdMono   = cfg.Command("mono", "make track monophon")
	argTracks = cmdMono.String("tracks", "track numbers, separated by commata", config.Shortflag('t'))
	//cmdDistribute = cfg.MustCommand("distribute", "distribute events across tracks, so that the first track contains tempo and timesignature changes and the other tracks contain the MIDI data for one channel per track")
)

func main() {
	err := run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
		fmt.Fprintf(os.Stderr, cfg.Usage())
		os.Exit(1)
	}

	os.Exit(0)
}

func run() (err error) {
	err = cfg.Run()
	if err != nil {
		return err
	}

	cmd := cfg.ActiveCommand()

	inFile, _ := filepath.Abs(argIn.Get())
	outFile, _ := filepath.Abs(argOut.Get())

	switch cmd {
	case cmdMono:
		return runMono(inFile, outFile)
		//	case cmdDistribute:
		//		return runDistribute(inFile, outFile)
	default:
		return fmt.Errorf("Unknown command: %q\n", cmd.CommmandName())
	}
}

func runMono(inFile, outFile string) error {
	if !argTracks.IsSet() {
		tracks, err := showTracks(inFile)
		if err != nil {
			return err
		}
		fmt.Printf("track argument missing, available tracks:\n%s\n", tracks)
		return nil
	}

	tracks := argTracks.Get()

	trs := strings.Split(tracks, ",")

	var trackNos []int

	for _, ts := range trs {
		ti, err := strconv.Atoi(strings.TrimSpace(ts))

		if err != nil {
			return fmt.Errorf("invalid track number %q", strings.TrimSpace(ts))
		}

		trackNos = append(trackNos, ti)
	}

	return smftrack.MonoizeTracks(inFile, outFile, trackNos)
}

/*
func runDistribute(inFile, outFile string) error {
	return nil
}
*/

func showTracks(inFile string) (string, error) {
	tracks, err := smftrack.TracksInfos(inFile)

	if err != nil {
		return "", err
	}

	var bf bytes.Buffer

	for _, tr := range tracks {
		if len(tr.Channels) > 0 {
			bf.WriteString(tr.String() + "\n")
		}
	}

	return bf.String(), nil
}
