package smf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"gitlab.com/gomidi/midi/v2/internal/runningstatus"
	vlq "gitlab.com/gomidi/midi/v2/internal/utils"
)

type wrWrapper struct {
	size int64
	wr   io.Writer
}

var _ io.Writer = &wrWrapper{}

func (w *wrWrapper) Write(p []byte) (int, error) {
	s, err := w.wr.Write(p)
	w.size += int64(s)
	return s, err
}

func newWriter(s *SMF, output io.Writer) *writer {
	// setup
	wr := &writer{}
	wr.SMF = s
	wr.output = &wrWrapper{wr: output}
	wr.currentChunk.SetType([4]byte{byte('M'), byte('T'), byte('r'), byte('k')})

	if !wr.SMF.NoRunningStatus {
		wr.runningWriter = runningstatus.NewSMFWriter()
	}
	return wr
}

type writer struct {
	*SMF
	currentChunk    chunk
	output          *wrWrapper
	headerWritten   bool
	tracksProcessed uint16
	deltatime       uint32
	absPos          uint64
	error           error
	runningWriter   runningstatus.SMFWriter
}

func (w *writer) printf(format string, vals ...interface{}) {
	if w.SMF.Logger == nil {
		return
	}

	w.SMF.Logger.Printf("smfwriter: "+format, vals...)
}

func (w *writer) Close() error {
	if cl, is := w.output.wr.(io.WriteCloser); is {
		w.printf("closing output")
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

func (w *writer) Position() uint64 {
	return w.absPos
}

// SetDelta sets the delta time in ticks for the next message(s)
func (w *writer) SetDelta(deltatime uint32) {
	w.deltatime = deltatime
}

// Write writes the message and returns the bytes that have been physically written.
// If a write fails with an error, every following attempt to write will return this first error,
// so de facto writing will be blocked.
func (w *writer) Write(m Message) (err error) {
	if w.error != nil {
		return w.error
	}
	if !w.headerWritten {
		w.error = w.WriteHeader()
	}
	if w.error != nil {
		w.printf("ERROR: writing header before midi message %#v failed: %v", m, w.error)
		w.error = fmt.Errorf("writing header before midi message %#v failed: %v", m, w.error)
		return w.error
	}
	defer func() {
		w.deltatime = 0
	}()

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
	switch tf := w.SMF.TimeFormat.(type) {
	case MetricTicks:
		ticks := tf.Ticks4th()
		if ticks > 32767 {
			ticks = 32767 // 32767 is the largest possible value, since bit 15 must always be 0
		}
		w.printf("writing metric ticks: %v", ticks)
		return binary.Write(wr, binary.BigEndian, uint16(ticks))
	case TimeCode:
		// multiplication with -1 makes sure that bit 15 is set
		err := binary.Write(wr, binary.BigEndian, int8(tf.FramesPerSecond)*-1)
		if err != nil {
			return err
		}
		w.printf("writing time code fps: %v subframes: %v", int8(tf.FramesPerSecond)*-1, tf.SubFrames)
		return binary.Write(wr, binary.BigEndian, tf.SubFrames)
	default:
		//panic(fmt.Sprintf("unsupported TimeFormat: %#v", w.header.TimeFormat))
		w.printf("ERROR: unsupported TimeFormat: %#v", w.SMF.TimeFormat)
		return fmt.Errorf("unsupported TimeFormat: %#v", w.SMF.TimeFormat)
	}
}

// <Header Chunk> = <chunk type><length><format><ntrks><division>
func (w *writer) writeHeader(wr io.Writer) error {
	w.printf("write header")
	var ch chunk
	ch.SetType([4]byte{byte('M'), byte('T'), byte('h'), byte('d')})
	var bf bytes.Buffer

	w.printf("write format %v", w.format)
	binary.Write(&bf, binary.BigEndian, w.format)
	w.printf("write num tracks %v", w.numTracks)
	binary.Write(&bf, binary.BigEndian, w.numTracks)

	err := w.writeTimeFormat(&bf)
	if err != nil {
		w.printf("ERROR: could not write header: %v", err)
		return fmt.Errorf("could not write header: %v", err)
	}

	_, err = ch.Write(bf.Bytes())
	if err != nil {
		w.printf("ERROR: could not write header: %v", err)
		return fmt.Errorf("could not write header: %v", err)
	}

	_, err = ch.WriteTo(wr)
	if err != nil {
		w.printf("ERROR: could not write header: %v", err)
		return fmt.Errorf("could not write header: %v", err)
	}
	w.printf("header written successfully")
	return nil
}

// <Track Chunk> = <chunk type><length><MTrk event>+
func (w *writer) writeChunkTo(wr io.Writer) (err error) {
	_, err = w.currentChunk.WriteTo(wr)

	if err != nil {
		w.printf("ERROR: could not write track %v: %v", w.tracksProcessed+1, err)
		return fmt.Errorf("could not write track %v: %v", w.tracksProcessed+1, err)
	}

	w.printf("track %v successfully written", w.tracksProcessed+1)

	if !w.SMF.NoRunningStatus {
		w.runningWriter = runningstatus.NewSMFWriter()
	}

	// remove the data for the next track
	w.currentChunk.Clear()
	w.deltatime = 0

	w.tracksProcessed++
	if w.numTracks == w.tracksProcessed {
		w.printf("last track written, finished")
		//		err = ErrFinished
	}

	return
}

func (w *writer) appendToChunk(deltaTime uint32, b []byte) {
	w.currentChunk.Write(append(vlq.VlqEncode(deltaTime), b...))
}

// delta is distance in time to last event in this track (independent of the channel)
func (w *writer) addMessage(deltaTime uint32, raw Message) {
	w.absPos += uint64(deltaTime)

	isSysEx := raw[0] == 0xF0 || raw[0] == 0xF7
	if isSysEx {
		// we have some sort of sysex, so we need to
		// calculate the length of msg[1:]
		// set msg to msg[0] + length of msg[1:] + msg[1:]
		if w.runningWriter != nil {
			w.runningWriter.ResetStatus()
		}

		//if sys, ok := msg.(sysex.Message); ok {
		b := []byte{raw[0]}
		b = append(b, vlq.VlqEncode(uint32(len(raw)-1))...)
		if len(raw[1:]) != 0 {
			b = append(b, raw[1:]...)
		}

		w.appendToChunk(deltaTime, b)
		return
	}

	if w.runningWriter != nil {
		w.appendToChunk(deltaTime, w.runningWriter.Write(raw))
		return
	}

	w.appendToChunk(deltaTime, raw)
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
