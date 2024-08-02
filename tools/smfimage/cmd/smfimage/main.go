package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/metakeule/observe/lib/runfunc"
	"gitlab.com/golang-utils/config/v2"
	"gitlab.com/gomidi/midi/tools/smfimage"
)

var (
	cfg   = config.New("smfimage", 0, 16, 0, "converts a Standard MIDI File (SMF) to an image", config.AsciiArt("smfimage"))
	inArg = cfg.LastString("midifile", "the SMF input file or a glob/pattern - then all matching MIDI files are affected (only the extensions .mid, .midi, .MID and .MIDI are taken into account). Output file is the same name as the input file but with the extension .png", config.Required())
	// outArg        = cfg.NewString("out", "the PNG output file (defaults to [in].png)", config.Shortflag('o'))
	showTracksArg   = cfg.Bool("names", "don't generate output, just show the tracknames", config.Shortflag('n'))
	orderArg        = cfg.String("order", "order of tracks: comma separated list of track ids, e.g. '3,2,0'")
	skipArg         = cfg.String("skip", "skip tracks: comma separated list of track channels (counting from 0), e.g. '3,2,0'")
	baseArg         = cfg.String("base", "base note. -1=most used note, c, c#/db, d, etc.", config.Shortflag('b'), config.Default("-1"))
	watchArg        = cfg.Bool("watch", "watching the file and export the image on every change", config.Shortflag('w'))
	sleepingTimeArg = cfg.Int("sleep", "sleeping time between invocations (in milliseconds)", config.Default(10))

	heightArg      = cfg.Int("height", "height of a 32th in pixel.", config.Default(8))
	trackborderArg = cfg.Int("border", "border of a track in pixel.", config.Default(4))
	widthArg       = cfg.Int("width", "width of a 32th in pixel.", config.Default(4))

	verboseArg     = cfg.Bool("verbose", "verbose output", config.Shortflag('v'))
	backgroundArg  = cfg.String("background", "set the background color. available colors are: black, white and transparent", config.Default("transparent"))
	beatsInGridArg = cfg.Bool("beats", "show beats in grid", config.Default(false))
	noBarLinesArg  = cfg.Bool("hidelines", "hide vertical bar lines", config.Default(false))
	monochromeArg  = cfg.Bool("monochrome", "use just one color", config.Default(false))
	curveArg       = cfg.Bool("curve", "use curve mode", config.Default(false))
	oneBarArg      = cfg.Bool("onebar", "write curve on one bar", config.Default(false))

	overviewArg = cfg.Bool("overview", "write harmonic overview above all tracks", config.Default(false))
	// radiusArg      = cfg.NewFloat32("radius", "radius of the curve", config.Default(float32(180)))

	SIGNAL_CHANNEL = make(chan os.Signal, 10)
)

func main() {
	err := run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func runFile(dir, file string) error {
	return _run()
}

func watch() error {

	r := runfunc.New(
		inArg.Get(),
		runFile,
		runfunc.Sleep(time.Millisecond*time.Duration(int(sleepingTimeArg.Get()))),
	)

	errors := make(chan error, 1)
	stopped, err := r.Run(errors)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case e := <-errors:
				fmt.Fprintf(os.Stderr, "error: %s\n", e)
			}
		}
	}()

	// ctrl+c
	<-SIGNAL_CHANNEL
	fmt.Println("\n--interrupted!")
	stopped.Kill()

	os.Exit(0)
	return nil
}

// returns infile -> outfile
func getFiles(globPattern string) (map[string]string, error) {
	d := filepath.Dir(globPattern)
	abs, err1 := filepath.Abs(d)

	if err1 != nil {
		return nil, fmt.Errorf("invalid directory: %#v", d)
	}

	di, err2 := os.Stat(abs)

	if err2 != nil {
		return nil, fmt.Errorf("directory does not exist: %#v", abs)
	}

	if !di.IsDir() {
		return nil, fmt.Errorf("is not a directory: %#v", abs)
	}

	var m = map[string]string{}
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, err
	}

	for _, mt := range matches {
		n := filepath.Base(mt)
		ext := filepath.Ext(n)

		switch ext {
		case ".mid", ".midi", ".MID", ".MIDI":
			m[mt] = strings.TrimSuffix(n, ext) + ".png"
		default:
			// do nothing
		}
	}

	return m, nil

}

func parseInts(orderStr string) (order []int, err error) {
	s := strings.Split(orderStr, ",")

	var i int
	for _, ss := range s {
		i, err = strconv.Atoi(ss)
		if err != nil {
			return
		}
		order = append(order, i)
	}
	return
}

func run() (err error) {
	err = cfg.Run()

	if err != nil {
		return
	}

	if watchArg.Get() {
		return watch()
	}
	return _run()
}

func _run() (err error) {

	var (
		names       []string
		order       []int
		skip        []int
		options     []smfimage.Option
		files       map[string]string
		verbose     bool
		background  string
		beatsInGrid bool
		noBarLines  bool
		monochrome  bool
		curve       bool
		oneBar      bool
		overview    bool
	)

	for {
		if err != nil {
			break
		}

		verbose = verboseArg.Get()
		background = backgroundArg.Get()
		beatsInGrid = beatsInGridArg.Get()
		noBarLines = noBarLinesArg.Get()
		monochrome = monochromeArg.Get()
		curve = curveArg.Get()
		oneBar = oneBarArg.Get()
		overview = overviewArg.Get()

		if showTracksArg.Get() {
			files, err = getFiles(inArg.Get())

			if err != nil {
				break
			}

			for in := range files {
				names, err = smfimage.SMF2PNG(in, "")

				if err != nil {
					break
				}

				fmt.Fprintf(os.Stdout, "\n\nFile: %#v\n", in)

				for i, n := range names {
					if n != "" {
						fmt.Fprintf(os.Stdout, "Track(s) %s on channel %v\n", n, i)
					}
				}
			}

			return
		}

		// options = append(options, smfimage.BaseNote(smfimage.Note(int(baseArg.Get()))))
		options = append(options, smfimage.Height(int(heightArg.Get())))
		options = append(options, smfimage.Width(int(widthArg.Get())))
		options = append(options, smfimage.TrackBorder(int(trackborderArg.Get())))
		// options = append(options, smfimage.Radius(float64(radiusArg.Get())))

		if baseArg.IsSet() && strings.TrimSpace(baseArg.Get()) != "-1" {
			num := smfimage.NoteToNumber(baseArg.Get())
			if num == -1 {
				err = fmt.Errorf("unknown note: %q", baseArg.Get())
				return
			}
			options = append(options, smfimage.BaseNote(smfimage.Note(num)))
		}

		if curve {
			options = append(options, smfimage.Curve())
		}

		if oneBar {
			options = append(options, smfimage.SingleBar())
		}

		if monochrome {
			options = append(options, smfimage.Monochrome())
		}

		if noBarLines {
			options = append(options, smfimage.NoBarLines())
		}

		if beatsInGrid {
			options = append(options, smfimage.BeatsInGrid())
		}

		if overview {
			options = append(options, smfimage.Overview())
		}

		if verbose {
			options = append(options, smfimage.Verbose())
		}

		switch background {
		case "white", "black", "transparent", "":
			options = append(options, smfimage.Background(background))
		default:
			err = fmt.Errorf("unknown background color: %s", background)
		}

		if err != nil {
			break
		}

		if o := orderArg.Get(); o != "" {
			order, err = parseInts(o)
			if err != nil {
				break
			}

			options = append(options, smfimage.TrackOrder(order...))
		}

		if s := skipArg.Get(); s != "" {
			skip, err = parseInts(s)
			if err != nil {
				break
			}

			options = append(options, smfimage.SkipTracks(skip...))
		}

		files, err = getFiles(inArg.Get())
		if err != nil {
			break
		}

		for in, out := range files {
			if verbose {
				fmt.Fprintf(os.Stdout, "processing %s.........", filepath.Base(in))
			}

			_, err = smfimage.SMF2PNG(in, out, options...)
			if err != nil {
				break
			}
			if verbose {
				fmt.Fprintf(os.Stdout, "done\n")
			}
		}

		break // leave at the end
	}

	return

}
