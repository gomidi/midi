package writer

import (
	"fmt"
	"io"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/cc"
	"gitlab.com/gomidi/midi/midimessage/channel"
	"gitlab.com/gomidi/midi/midiwriter"
)

// Writer writes live MIDI data. Its methods must not be called concurrently
type Writer struct {
	wr              midi.Writer
	channel         uint8
	noteState       [16][128]bool
	noConsolidation bool
}

type ChannelWriter interface {
	Channel() uint8 /* 0-15 */
	SetChannel(no uint8 /* 0-15 */)
	midi.Writer
}

var _ midi.Writer = &Writer{}
var _ ChannelWriter = &Writer{}

func Wrap(wr midi.Writer) *Writer {
	return &Writer{wr: wr, channel: 0}
}

// New creates and new Writer for writing of "live" MIDI data ("over the wire")
// By default it makes no use of the running status.
func New(dest io.Writer, options ...midiwriter.Option) *Writer {
	options = append(
		[]midiwriter.Option{
			midiwriter.NoRunningStatus(),
		}, options...)

	return Wrap(midiwriter.New(dest, options...))
}

// SetChannel sets the channel for the following midi messages
// Channel numbers are counted from 0 to 15 (MIDI channel 1 to 16).
// The initial channel number is 0.
func (w *Writer) SetChannel(no uint8 /* 0-15 */) {
	w.channel = no
}

func (w *Writer) Channel() uint8 /* 0-15 */ {
	return w.channel
}

// silentium for a single channel
func silentium(w *Writer, ch uint8, force bool) (err error) {
	for key, val := range w.noteState[ch] {
		if force || val {
			err = w.Write(channel.Channel(ch).NoteOff(uint8(key)))
		}
	}

	if force {
		err = w.Write(channel.Channel(ch).ControlChange(cc.AllNotesOff, 0))
		err = w.Write(channel.Channel(ch).ControlChange(cc.AllSoundOff, 0))
	}
	return
}

// Silence sends a  note off message for every running note
// on the given channel. If the channel is -1, every channel is affected.
// If note consolidation is switched off, notes aren't tracked and therefor
// every possible note will get a note off. This could also be enforced by setting
// force to true. If force is true, additionally the cc messages AllNotesOff (123) and AllSoundOff (120)
// are send.
// If channel is > 15, the method panics.
// The error or not error of the last write is returned.
// (i.e. it does not return on the first error, but tries everything instead to make it silent)
func (w *Writer) Silence(ch int8, force bool) (err error) {
	if ch > 15 {
		panic("invalid channel number")
	}
	if w.noConsolidation {
		force = true
	}

	// single channel
	if ch >= 0 {
		err = silentium(w, uint8(ch), force)
		// set note states for the channel
		w.noteState[ch] = [128]bool{}
		return
	}

	// all channels
	for c := 0; c < 16; c++ {
		err = silentium(w, uint8(c), force)
	}
	// reset all note states
	w.noteState = [16][128]bool{}
	return
}

// ConsolidateNotes enables/disables the midi note consolidation (default: enabled)
// When enabled, midi note on/off messages are consolidated, that means the on/off state of
// every possible note on every channel is tracked and note on messages are only
// written, if the corresponding note is off and vice versa. Note on messages
// with a velocity of 0 behave the same way as note offs.
// The tracking of the notes is not cross SMF tracks, i.e. a meta.EndOfTrack
// message will reset the tracking.
// The consolidation should prevent "hanging" notes in most cases.
// If on is true, the note will be started tracking again (fresh state), assuming no note is currently running.
func (w *Writer) ConsolidateNotes(on bool) {
	if on {
		w.noteState = [16][128]bool{}
	}
	w.noConsolidation = !on
}

// Write writes the given midi.Message. By default, midi notes are consolidated (see ConsolidateNotes method)
func (w *Writer) Write(msg midi.Message) error {
	if w.noConsolidation {
		return w.wr.Write(msg)
	}
	switch m := msg.(type) {
	case channel.NoteOn:
		if m.Velocity() > 0 && w.noteState[m.Channel()][m.Key()] {
			return fmt.Errorf("can't write %s. note already running.", msg)
		}
		if m.Velocity() == 0 && !w.noteState[m.Channel()][m.Key()] {
			return fmt.Errorf("can't write %s. note is not running.", msg)
		}
		w.noteState[m.Channel()][m.Key()] = m.Velocity() > 0
	case channel.NoteOff:
		if !w.noteState[m.Channel()][m.Key()] {
			return fmt.Errorf("can't write %s. note is not running.", msg)
		}
		w.noteState[m.Channel()][m.Key()] = false
	case channel.NoteOffVelocity:
		if !w.noteState[m.Channel()][m.Key()] {
			return fmt.Errorf("can't write %s. note is not running.", msg)
		}
		w.noteState[m.Channel()][m.Key()] = false
	}
	return w.wr.Write(msg)
}
