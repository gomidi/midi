package main

import (
	"fmt"
	"os"

	"gitlab.com/golang-utils/config/v2"
	"gitlab.com/gomidi/midi/tools/smftrack"
	"gitlab.com/gomidi/midi/v2/smf"
)

var (
	cfg    = config.New("smf1to0", 0, 0, 2, "converts a SMF1 MIDI file to a SMF0 file", config.AsciiArt("smf1to0"))
	inArg  = cfg.String("in", "SMF1 file that is the input", config.Required(), config.Shortflag('i'))
	outArg = cfg.String("out", "name of the resulting SMF0 file", config.Required(), config.Shortflag('o'))
)

func main() {
	err := run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
		os.Exit(1)
	}
}

func run() (err error) {
	err = cfg.Run()
	if err != nil {
		return err
	}

	return convertSMF1To0(inArg.Get(), outArg.Get())
}

func convertSMF1To0(inFile, outFile string) (err error) {
	var out *os.File
	src, err := smf.ReadFile(inFile)
	if err != nil {
		return err
	}
	if hf := src.Format(); hf != 1 {
		return fmt.Errorf("file %q is not a SMF1 file, but %v", inFile, hf)
	}

	out, err = os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = (smftrack.SMF1{}).ToSMF0(src, out)
	return
}
