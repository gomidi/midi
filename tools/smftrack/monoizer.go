package smftrack

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

func Monoize(in io.Reader, out io.Writer, trackNos []int) error {
	var isAffectedTrack = map[int]bool{}

	for _, tr := range trackNos {
		isAffectedTrack[tr] = true
	}

	bts, err := ioutil.ReadAll(in)

	if err != nil {
		return fmt.Errorf("could not read SMF %v\n", err.Error())
	}

	rd, err1 := smf.ReadFrom(bytes.NewBuffer(bts))
	if err1 != nil {
		return fmt.Errorf("could not read SMF %v\n", err1.Error())
	}

	wr := smf.New()
	wr.TimeFormat = rd.TimeFormat

	var currentTrack smf.Track
	var currentMonoizedTrack = newMonoizedTrack()

	for trno, track := range rd.Tracks {
		for _, ev := range track {
			if !isAffectedTrack[trno] {
				currentTrack.Add(ev.Delta, ev.Message)
			} else {
				currentMonoizedTrack.WriteMessage(ev.Delta, ev.Message)
			}
		}

		if !isAffectedTrack[trno] {
			currentTrack.Close(0)
			wr.Add(currentTrack)
		} else {
			currentMonoizedTrack.track.Close(0)
			wr.Add(*currentMonoizedTrack.track)
		}
		currentTrack = smf.Track{}
		currentMonoizedTrack = newMonoizedTrack()
	}

	_, err = wr.WriteTo(out)
	return err
}

func MonoizeTracks(inFile, outFile string, trackNos []int) error {

	out, err := os.Create(outFile)

	if err != nil {
		return fmt.Errorf("can't create output file %q", outFile)
	}

	defer out.Close()

	in, err2 := os.Open(inFile)

	if err2 != nil {
		return fmt.Errorf("can't open input file %q", inFile)
	}

	defer in.Close()

	return Monoize(in, out, trackNos)
}

func newMonoizedTrack() *monoizedTrack {
	return &monoizedTrack{
		track:     &smf.Track{},
		lastNotes: map[uint8]int8{},
	}
}

type monoizedTrack struct {
	lastPosition uint64
	lastNotes    map[uint8]int8
	track        *smf.Track
}

func (m *monoizedTrack) WriteMessage(delta uint32, msg smf.Message) {
	var channel, val1, val2 uint8

	switch {
	case msg.GetNoteStart(&channel, &val1, &val2):
		if m.lastNotes[channel] > 0 {
			m.track.Add(delta, smf.Message(midi.NoteOff(channel, uint8(m.lastNotes[channel]))))
			m.track.Add(0, msg)
		} else {
			m.track.Add(delta, msg)
		}
		m.lastNotes[channel] = int8(val1)
		m.lastPosition = m.lastPosition + uint64(delta)
	case msg.GetNoteEnd(&channel, &val1):
		if m.lastNotes[channel] > 0 {
			m.track.Add(delta, smf.Message(midi.NoteOff(channel, uint8(m.lastNotes[channel]))))
			m.lastNotes[channel] = 0
			m.lastPosition = m.lastPosition + uint64(delta)
		} else {
			// do nothing
		}
	default:
		m.track.Add(delta, msg)
		m.lastPosition = m.lastPosition + uint64(delta)
	}
}
