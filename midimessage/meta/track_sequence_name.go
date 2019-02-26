package meta

import (
	"fmt"
	"io"
)

// TrackSequenceName represents a MIDI track name/sequence name message
// If in a format 0 track, or the first track in a format 1 file, the name of the sequence. Otherwise, the name of the track.
type TrackSequenceName string

// String represents the MIDI track/sequence name message as a string (for debugging)
func (m TrackSequenceName) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}
func (m TrackSequenceName) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return TrackSequenceName(text), nil

}

// Text returns the name
func (m TrackSequenceName) Text() string {
	return string(m)
}

func (m TrackSequenceName) meta() {}

// Raw returns the raw bytes for the message
func (m TrackSequenceName) Raw() []byte {
	return (&metaMessage{
		Typ:  byteTrackSequenceName,
		Data: []byte(m),
	}).Bytes()
}
