package meta

import (
	"fmt"
	"io"
)

/*
	http://midi.teragonaudio.com/tech/midifile/port.htm

	   Device (Port) Name

	   FF 09 len text

	   The name of the MIDI device (port) where the track is routed.
	   This replaces the "MIDI Port" meta-Event which some sequencers
	   formally used to route MIDI tracks to various MIDI ports
	   (in order to support more than 16 MIDI channels).

	   For example, assume that you have a MIDI interface that has 4 MIDI output ports.
	   They are listed as "MIDI Out 1", "MIDI Out 2", "MIDI Out 3", and "MIDI Out 4".
	   If you wished a particular MTrk to use "MIDI Out 1" then you would put a
	   Port Name meta-event at the beginning of the MTrk, with "MIDI Out 1" as the text.

	   All MIDI events that occur in the MTrk, after a given Port Name event, will be
	   routed to that port.

	   In a format 0 MIDI file, it would be permissible to have numerous Port Name events
	   intermixed with MIDI events, so that the one MTrk could address numerous ports.
	   But that would likely make the MIDI file much larger than it need be. The Port Name
	   event is useful primarily in format 1 MIDI files, where each MTrk gets routed to
	   one particular port.

	   Note that len could be a series of bytes since it is expressed as a variable length quantity.
*/

// Device represents a MIDI device point message
type Device string

// String represents the device point MIDI message as a string (for debugging)
func (m Device) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

func (m Device) meta() {}

func (m Device) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)
	if err != nil {
		return nil, err
	}

	return Device(text), nil
}

// Raw returns the raw MIDI data
func (m Device) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteDevicePort),
		Data: []byte(m),
	}).Bytes()
}

// Text returns the name of the device port
func (m Device) Text() string {
	return string(m)
}
