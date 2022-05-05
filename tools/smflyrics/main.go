package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/gomidi/midi/v2/smf"
	"gitlab.com/metakeule/config"
)

var (
	cfg = config.MustNew("smflyrics", "1.6.1",
		"extracts lyrics from a SMF file, tracks are separated by an empty line")

	argFile = cfg.LastString("file",
		"the SMF file that is read in",
		config.Required)

	argTrack = cfg.NewInt32("track",
		"the track from which the lyrics are taken. -1 means all tracks, 0 is the first, 1 the second etc",
		config.Shortflag('t'), config.Default(int32(-1)))

	argIncludeText = cfg.NewBool("text",
		"include free text entries in the SMF file. Text is surrounded by doublequotes")

	argJson = cfg.NewBool("json", "output json format",
		config.Shortflag('j'))
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}

func run() (err error) {
	err = cfg.Run()
	if err != nil {
		fmt.Fprintln(os.Stdout, cfg.Usage())
		return err
	}

	mf, err := smf.ReadFile(argFile.Get())

	if err != nil {
		return err
	}

	j := newJson(mf)
	j.ReadSMF()

	if argJson.Get() {
		bt, err := json.MarshalIndent(j, "", " ")
		if err != nil {
			return err
		}
		fmt.Println(string(bt))
		return nil
	}

	return nil
}

func newJson(mf *smf.SMF) *Json {
	var j = &Json{}
	j.Tracks = []*Track{&Track{}}
	j.mf = mf
	j.File = argFile.Get()
	j.includeText = argIncludeText.Get()
	j.requestedTrack = int(argTrack.Get())
	j.printText = !argJson.Get()

	return j
}

type Json struct {
	Tracks         []*Track `json:"tracks,omitempty"`
	File           string   `json:"file"`
	current        int
	includeText    bool
	requestedTrack int // if < 0 : all tracks
	trackNo        int
	mf             *smf.SMF
	printText      bool
}

func (p *Json) ReadSMF() {
	for _, tr := range p.mf.Tracks {
		if p.requestedTrack < 0 || int(p.requestedTrack) == p.trackNo {
			p.ReadTrack(tr)
		}
		p.trackNo++
	}
}

func (t *Json) ReadTrack(tr smf.Track) {
	var text string
	for _, ev := range tr {
		if !ev.Message.IsMeta() {
			continue
		}
		msg := ev.Message

		switch {
		case msg.GetMetaLyric(&text):
			t.Tracks[t.current].addLyric(text)
			if t.printText {
				fmt.Print(text + " ")
			}
		case t.includeText && msg.GetMetaText(&text):
			t.Tracks[t.current].addText(text)
			if t.printText {
				fmt.Printf(" %q ", text)
			}
		case msg.GetMetaTrackName(&text):
			t.Tracks[t.current].Name = text
			if t.printText {
				fmt.Println(fmt.Sprintf("[track: %v]\n", text))
			}
		case msg.GetMetaInstrument(&text):
			t.Tracks[t.current].Instrument = text
			if t.printText {
				fmt.Println(fmt.Sprintf("[instrument: %v]\n", text))
			}
		case msg.GetMetaProgramName(&text):
			t.Tracks[t.current].Program = text
			if t.printText {
				fmt.Println(fmt.Sprintf("[program: %v]\n", text))
			}
		}
	}

	t.Tracks = append(t.Tracks, &Track{No: t.trackNo + 1})
	t.current = len(t.Tracks) - 1

	if t.printText {
		fmt.Printf("\n\n------------------------\n\n")
	}
}

type Track struct {
	Program    string   `json:"program,omitempty"`
	Name       string   `json:"name,omitempty"`
	Instrument string   `json:"instrument,omitempty"`
	Texts      []string `json:"texts,omitempty"`
	Lyrics     []string `json:"lyrics,omitempty"`
	No         int      `json:"no"`
}

func (t *Track) addLyric(l string) {
	t.Lyrics = append(t.Lyrics, l)
}

func (t *Track) addText(tx string) {
	t.Texts = append(t.Texts, tx)
}
