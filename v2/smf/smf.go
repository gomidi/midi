package smf

import (
	"fmt"
	"io"
	"os"
)

// New returns a SMF file of format type 0 (single track), that becomes type 1 (multi track), if you add tracks
func New() *SMF {
	return newSMF(0)
}

// NewSMF1 returns a SMF file of format type 1 (multi track)
func NewSMF1() *SMF {
	return newSMF(1)
}

// NewSMF2 returns a SMF file of format type 2 (multi sequence)
func NewSMF2() *SMF {
	return newSMF(2)
}

func newSMF(format uint16) *SMF {
	s := &SMF{
		format: format,
	}
	s.TimeFormat = MetricTicks(960)
	return s
}

type SMF struct {
	//Header       SMFHeader
	// Format is the SMF file format: SMF0, SMF1 or SMF2.
	format uint16
	//Format

	// NumTracks is the number of tracks (0 indicates that the number is not set yet).
	numTracks uint16

	tracks []*Track

	// TimeFormat is the time format (either MetricTicks or TimeCode).
	//	timeFormat TimeFormat

	tempoChanges TempoChanges

	finished bool

	//opts []Option
	//Config Config

	NoRunningStatus bool
	Logger          Logger
	TimeFormat      TimeFormat
}

func (s *SMF) TempoChanges() TempoChanges {
	return s.tempoChanges
}

func (s *SMF) Tracks() []*Track {
	return s.tracks
}

func (s *SMF) NumTracks() uint16 {
	return uint16(len(s.tracks))
}

// WriteFile creates file, calls callback with a writer and closes file.
//
// WriteFile makes sure that the data of the last track is written by sending
// an meta.EndOfTrack message after callback has been run.
//
// For single track (SMF0) files this makes sense since no meta.EndOfTrack message
// must then be send from callback (although it does not harm).
//
// For multitrack files however there must be sending of meta.EndOfTrack anyway,
// so it is better practise to send it after each track (including the last one).
// The options and their defaults are the same as for New and they are documented
// at the corresponding option.
// The callback may call the given writer to write messages. If any of this write
// results in an error, the file won't be written and the error is returned.
// Only a successful write will manifest itself in the file being created.
//func (s *SMF) WriteFile(file string, options ...Option) error {

//var s io.WriterTo = &smf{}
func (s *SMF) WriteFile(file string) error {
	f, err := os.Create(file)

	if err != nil {
		return fmt.Errorf("writing midi file failed: could not create file %#v", file)
	}

	//err = s.WriteTo(f)
	err = s.WriteTo(f)
	f.Close()

	if err != nil {
		os.Remove(file)
		return fmt.Errorf("writing to midi file %#v failed: %v", file, err)
	}

	return nil
}

func (s *SMF) WriteTo(f io.Writer) (err error) {
	s.numTracks = uint16(len(s.tracks))
	if s.numTracks == 0 {
		return fmt.Errorf("no track added")
	}
	if s.numTracks > 1 && s.format == 0 {
		s.format = 1
	}
	//wr := newWriter(f, options...)
	//fmt.Printf("numtracks: %v\n", s.numTracks)
	wr := newWriter(s, f)
	err = wr.WriteHeader()
	if err != nil {
		return fmt.Errorf("could not write header: %v", err)
	}

	for _, t := range s.tracks {
		t.Close(0) // just to be sure
		for _, ev := range t.Events {
			//fmt.Printf("written ev: %v\n ", ev)
			wr.SetDelta(ev.Delta)
			err = wr.Write(ev.Data)
			if err != nil {
				break
			}
		}

		err = wr.writeChunkTo(wr.output)

		if err != nil {
			break
		}
	}

	return
}

// AddAndClose closes the given track at deltatime and adds it to the smf
func (s *SMF) AddAndClose(deltatime uint32, t *Track) {
	t.Close(deltatime)
	s.tracks = append(s.tracks, t)
}

//var ErrFinished = errors.New("SMF action finished successfully")

func (s SMF) Format() uint16 {
	return s.format
}
