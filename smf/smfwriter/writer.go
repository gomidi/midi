package smfwriter

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/gomidi/midi/internal/runningstatus"
	"github.com/gomidi/midi/internal/vlq"

	"encoding/binary"

	"github.com/gomidi/midi"

	"github.com/gomidi/midi/midimessage/meta"
	"github.com/gomidi/midi/smf"
)

// WriteFile creates file, calls callback with a writer and closes file
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
func WriteFile(file string, callback func(smf.Writer), options ...Option) error {
	f, err := os.Create(file)

	if err != nil {
		return fmt.Errorf("writing midi file failed: could not create file %#v", file)
	}

	defer func() {
		f.Close()
	}()

	wr := newWriter(f, options...)
	err = wr.WriteHeader()
	if err != nil {
		f.Close()
		os.Remove(file)
		return fmt.Errorf("could not write header to midi file %#v: %v", file, err)
	}
	callback(wr)

	if wr.error != nil && wr.error != smf.ErrFinished {
		f.Close()
		os.Remove(file)
		return fmt.Errorf("writing of midi file %#v aborted due to error returned from callback: %v", file, wr.error)
	}

	// make sure the data of the last track is written
	err = wr.Write(meta.EndOfTrack)

	if err != nil && err != smf.ErrFinished {
		f.Close()
		os.Remove(file)
		return fmt.Errorf("could not write end of track message to midi file %#v", file)
	}

	err = f.Close()
	if err != nil {
		os.Remove(file)
		return fmt.Errorf("could not close midi file %#v", file)
	}

	return nil
}

// New returns a Writer
//
// The writer just uses an io.Writer..It is the responsibility of the caller to open and close any file where appropriate.
//
// For the documentation of the Write and the SetDelta method, consult the documentation for smf.Writer.
//
// The options and their defaults are documented at the corresponding option.
// When New returns, the header has already been written to dest.
// Any error that happened during the header writing is returned. In this case writer is nil.
func New(dest io.Writer, opts ...Option) smf.Writer {
	return newWriter(dest, opts...)
}

type writer struct {
	header          smf.Header
	track           smf.Chunk
	output          io.Writer
	headerWritten   bool
	tracksProcessed uint16
	deltatime       uint32
	noRunningStatus bool
	error           error
	runningWriter   runningstatus.SMFWriter
}

func (w *writer) Close() error {
	if cl, is := w.output.(io.WriteCloser); is {
		return cl.Close()
	}
	return nil
}

func (w *writer) WriteHeader() error {
	if w.headerWritten {
		return w.error
	}
	err := w.writeHeader(w.output)
	w.headerWritten = true

	if err != nil {
		w.error = err
	}

	return err
}

//
func newWriter(output io.Writer, opts ...Option) *writer {

	// setup
	wr := &writer{}
	wr.output = output
	wr.track.SetType([4]byte{byte('M'), byte('T'), byte('r'), byte('k')})

	// defaults
	wr.header.TimeFormat = smf.MetricTicks(0) // take the default, based on smf package (should be 960)
	wr.header.NumTracks = 1
	wr.header.Format = smf.SMF0

	// overrides with options
	for _, opt := range opts {
		opt(wr)
	}

	if !wr.noRunningStatus {
		wr.runningWriter = runningstatus.NewSMFWriter()
	}

	// if midiformat is undefined (see above), i.e. not set via options
	// set the default, which is format 0 for one track and format 1 for multitracks
	// if wr.header.MidiFormat == format(10) {
	if wr.header.Format != smf.SMF2 && wr.header.NumTracks > 1 {
		wr.header.Format = smf.SMF1
	}

	return wr
}

// SetDelta sets the delta time in ticks for the next message(s)
func (w *writer) SetDelta(deltatime uint32) {
	w.deltatime = deltatime
}

// Header returns the smf.Header of the file
func (w *writer) Header() smf.Header {
	return w.header
}

// Write writes the message and returns the bytes that have been physically written.
// If a write fails with an error, every following attempt to write will return this first error,
// so de facto writing will be blocked.
func (w *writer) Write(m midi.Message) (err error) {
	if w.error != nil {
		return w.error
	}
	if !w.headerWritten {
		w.error = w.WriteHeader()
	}
	if w.error != nil {
		w.error = fmt.Errorf("writing header before midi message %#v failed: %v", m, w.error)
		return w.error
	}
	defer func() {
		w.deltatime = 0
	}()

	if w.header.NumTracks == w.tracksProcessed {
		w.error = smf.ErrFinished
		return w.error
	}

	if m == meta.EndOfTrack {
		w.addMessage(w.deltatime, m)
		err = w.writeTrackTo(w.output)
		if err != nil {
			w.error = err
		}
		return
	}
	w.addMessage(w.deltatime, m)
	return
}

/*

					| time type            | bit 15 | bits 14 thru 8        | bits 7 thru 0   |
					-----------------------------------------------------------------------------
				  | metrical time        |      0 |         ticks per quarter-note          |
				  | time-code-based time |      1 | negative SMPTE format | ticks per frame |

		If bit 15 of <division> is zero, the bits 14 thru 0 represent the number of delta time "ticks" which make up a
		quarter-note. For instance, if division is 96, then a time interval of an eighth-note between two events in the
		file would be 48.

		If bit 15 of <division> is a one, delta times in a file correspond to subdivisions of a second, in a way
		consistent with SMPTE and MIDI Time Code. Bits 14 thru 8 contain one of the four values -24, -25, -29, or
		-30, corresponding to the four standard SMPTE and MIDI Time Code formats (-29 corresponds to 30 drop
		frame), and represents the number of frames per second. These negative numbers are stored in two's
		compliment form. The second byte (stored positive) is the resolution within a frame: typical values may be 4
		(MIDI Time Code resolution), 8, 10, 80 (bit resolution), or 100. This stream allows exact specifications of
		time-code-based tracks, but also allows millisecond-based tracks by specifying 25 frames/sec and a resolution
		of 40 units per frame. If the events in a file are stored with a bit resolution of thirty-frame time code, the
		division word would be E250 hex. (=> 1110001001010000 or 57936)


	/* unit of time for delta timing. If the value is positive, then it represents the units per beat.
	For example, +96 would mean 96 ticks per beat. If the value is negative, delta times are in SMPTE compatible units.
*/
func (w *writer) writeTimeFormat(wr io.Writer) error {
	switch tf := w.header.TimeFormat.(type) {
	case smf.MetricTicks:
		ticks := tf.Ticks4th()
		if ticks > 32767 {
			ticks = 32767 // 32767 is the largest possible value, since bit 15 must always be 0
		}
		return binary.Write(wr, binary.BigEndian, uint16(ticks))
	case smf.TimeCode:
		// multiplication with -1 makes sure that bit 15 is set
		err := binary.Write(wr, binary.BigEndian, int8(tf.FramesPerSecond)*-1)
		if err != nil {
			return err
		}
		return binary.Write(wr, binary.BigEndian, tf.SubFrames)
	default:
		//panic(fmt.Sprintf("unsupported TimeFormat: %#v", w.header.TimeFormat))
		return fmt.Errorf("unsupported TimeFormat: %#v", w.header.TimeFormat)
	}
}

// <Header Chunk> = <chunk type><length><format><ntrks><division>
func (w *writer) writeHeader(wr io.Writer) error {
	var ch smf.Chunk
	ch.SetType([4]byte{byte('M'), byte('T'), byte('h'), byte('d')})
	var bf bytes.Buffer

	binary.Write(&bf, binary.BigEndian, w.header.Format.Type())
	binary.Write(&bf, binary.BigEndian, w.header.NumTracks)

	err := w.writeTimeFormat(&bf)
	if err != nil {
		return fmt.Errorf("could not write header: %v", err)
	}

	_, err = ch.Write(bf.Bytes())
	if err != nil {
		return fmt.Errorf("could not write header: %v", err)
	}

	_, err = ch.WriteTo(wr)
	if err != nil {
		return fmt.Errorf("could not write header: %v", err)
	}
	return nil
}

// <Track Chunk> = <chunk type><length><MTrk event>+
func (w *writer) writeTrackTo(wr io.Writer) (err error) {
	_, err = w.track.WriteTo(wr)

	if err != nil {
		return fmt.Errorf("could not write track %v: %v", w.tracksProcessed+1, err)
	}

	if !w.noRunningStatus {
		w.runningWriter = runningstatus.NewSMFWriter()
	}

	// remove the data for the next track
	w.track.Clear()
	w.deltatime = 0

	w.tracksProcessed++
	if w.header.NumTracks == w.tracksProcessed {
		err = smf.ErrFinished
	}

	return
}

func (w *writer) appendToChunk(deltaTime uint32, b []byte) {
	w.track.Write(append(vlq.Encode(deltaTime), b...))
	//t.track.data = append(t.track.data, append(vlq.Encode(deltaTime), b...)...)
}

// delta is distance in time to last event in this track (independent of the channel)
func (w *writer) addMessage(deltaTime uint32, msg midi.Message) {
	// we have some sort of sysex, so we need to
	// calculate the length of msg[1:]
	// set msg to msg[0] + length of msg[1:] + msg[1:]
	raw := msg.Raw()
	if raw[0] == 0xF0 || raw[0] == 0xF7 {
		if w.runningWriter != nil {
			w.runningWriter.ResetStatus()
		}

		//if sys, ok := msg.(sysex.Message); ok {
		b := []byte{raw[0]}
		b = append(b, vlq.Encode(uint32(len(raw)-1))...)
		if len(raw[1:]) != 0 {
			b = append(b, raw[1:]...)
		}

		w.appendToChunk(deltaTime, b)
		return
	}

	if w.runningWriter != nil {
		w.appendToChunk(deltaTime, w.runningWriter.Write(msg))
		return
	}

	w.appendToChunk(deltaTime, msg.Raw())
}

/*
from http://www.artandscienceofsound.com/article/standardmidifiles

Depending upon the application you are using to create the file in the first place, header information may automatically be saved from within parameters set in the application, or may need to be placed in a ‘set-up’ bar before the music data commences.

Either way, information that should be considered includes:

GM/GS Reset message

Per MIDI Channel
Bank Select (0=GM) / Program Change #
Reset All Controllers (not all devices may recognize this command so you may prefer to zero out or reset individual controllers)
Initial Volume (CC7) (standard level = 100)
Expression (CC11) (initial level set to 127)
Hold pedal (0 = off)
Pan (Center = 64)
Modulation (0)
Pitch bend range
Reverb (0 = off)
Chorus level (0 = off)

System Exclusive data

If RPNs or more detailed controller messages are being employed in the file these should also be reset or normalized in the header.

If you are inputting header data yourself it is advisable not to clump all such information together but rather space it out in intervals of 5-10 ticks. Certainly if a file is designed to be looped, having too much data play simultaneously will cause most playback devices to ‘choke, ’ and throw off your timing.
*/

/*
from http://www.artandscienceofsound.com/article/standardmidifiles

Depending upon the application you are using to create the file in the first place, header information may automatically be saved from within parameters set in the application, or may need to be placed in a ‘set-up’ bar before the music data commences.

Either way, information that should be considered includes:

GM/GS Reset message

Per MIDI Channel
Bank Select (0=GM) / Program Change #
Reset All Controllers (not all devices may recognize this command so you may prefer to zero out or reset individual controllers)
Initial Volume (CC7) (standard level = 100)
Expression (CC11) (initial level set to 127)
Hold pedal (0 = off)
Pan (Center = 64)
Modulation (0)
Pitch bend range
Reverb (0 = off)
Chorus level (0 = off)

System Exclusive data

If RPNs or more detailed controller messages are being employed in the file these should also be reset or normalized in the header.

If you are inputting header data yourself it is advisable not to clump all such information together but rather space it out in intervals of 5-10 ticks. Certainly if a file is designed to be looped, having too much data play simultaneously will cause most playback devices to ‘choke, ’ and throw off your timing.
*/
