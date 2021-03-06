package smf

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gitlab.com/gomidi/midi/v2/internal/runningstatus"
	vlq "gitlab.com/gomidi/midi/v2/internal/utils"

	"encoding/binary"
	//	"gitlab.com/gomidi/midi/midimessage/meta"
	//"gitlab.com/gomidi/midi/smf"
)

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

		err = wr.writeTrackTo(wr.output)

		if err != nil {
			break
		}
	}

	return
}

/*
// New returns a Writer
//
// The writer just uses an io.Writer. It is the responsibility of the caller to open and close any file where appropriate.
//
// For the documentation of the Write and the SetDelta method, consult the documentation for smf.Writer.
//
// The options and their defaults are documented at the corresponding option.
// When New returns, the header has already been written to dest.
// Any error that happened during the header writing is returned. In this case writer is nil.
func New(dest io.Writer, opts ...Option) smf.Writer {
	return newWriter(dest, opts...)
}
*/

func newSMF(format uint16) *SMF {
	s := &SMF{
		format: format,
	}
	s.TimeFormat = MetricTicks(960)
	return s
}

// AddAndClose closes the given track at deltatime and adds it to the smf
func (s *SMF) AddAndClose(deltatime uint32, t *Track) {
	t.Close(deltatime)
	s.tracks = append(s.tracks, t)
}

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

type writer struct {
	*SMF
	//header          smf.Header
	track           Chunk
	output          io.Writer
	headerWritten   bool
	tracksProcessed uint16
	deltatime       uint32
	absPos          uint64
	//noRunningStatus bool
	error error
	//logger          logger
	runningWriter runningstatus.SMFWriter
}

func (w *writer) printf(format string, vals ...interface{}) {
	if w.SMF.Logger == nil {
		return
	}

	w.SMF.Logger.Printf("smfwriter: "+format, vals...)
}

func (w *writer) Close() error {
	if cl, is := w.output.(io.WriteCloser); is {
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

//
//func newWriter(s *SMF, output io.Writer, opts ...Option) *writer {
func newWriter(s *SMF, output io.Writer) *writer {

	// setup
	wr := &writer{}
	wr.SMF = s
	wr.output = output
	wr.track.SetType([4]byte{byte('M'), byte('T'), byte('r'), byte('k')})

	// defaults
	/*
		wr.header.TimeFormat = smf.MetricTicks(0) // take the default, based on smf package (should be 960)
		wr.header.NumTracks = 1
		wr.header.Format = smf.SMF0
	*/

	// overrides with options

	/*
		for _, opt := range opts {
			opt(wr)
		}
	*/

	if !wr.SMF.NoRunningStatus {
		wr.runningWriter = runningstatus.NewSMFWriter()
	}

	return wr
}

func (w *writer) Position() uint64 {
	return w.absPos
}

// SetDelta sets the delta time in ticks for the next message(s)
func (w *writer) SetDelta(deltatime uint32) {
	w.deltatime = deltatime
}

/*
// Header returns the smf.Header of the file
func (w *writer) Header() smf.Header {
	return w.header
}
*/

// Write writes the message and returns the bytes that have been physically written.
// If a write fails with an error, every following attempt to write will return this first error,
// so de facto writing will be blocked.
func (w *writer) Write(m []byte) (err error) {
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

	/*
		if w.numTracks == w.tracksProcessed {
			w.printf("last track written, finished")
			w.error = ErrFinished
			return w.error
		}

		if midi2.GetMessageType(m) == midi2.MetaEndOfTrackMsg {
			w.addMessage(w.deltatime, m)
			err = w.writeTrackTo(w.output)
			if err != nil {
				w.error = err
			}
			return
		}
	*/
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
	var ch Chunk
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
func (w *writer) writeTrackTo(wr io.Writer) (err error) {
	_, err = w.track.WriteTo(wr)

	if err != nil {
		w.printf("ERROR: could not write track %v: %v", w.tracksProcessed+1, err)
		return fmt.Errorf("could not write track %v: %v", w.tracksProcessed+1, err)
	}

	w.printf("track %v successfully written", w.tracksProcessed+1)

	if !w.SMF.NoRunningStatus {
		w.runningWriter = runningstatus.NewSMFWriter()
	}

	// remove the data for the next track
	w.track.Clear()
	w.deltatime = 0

	w.tracksProcessed++
	if w.numTracks == w.tracksProcessed {
		w.printf("last track written, finished")
		//		err = ErrFinished
	}

	return
}

func (w *writer) appendToChunk(deltaTime uint32, b []byte) {
	w.track.Write(append(vlq.VlqEncode(deltaTime), b...))
	//t.track.data = append(t.track.data, append(vlq.Encode(deltaTime), b...)...)
}

// delta is distance in time to last event in this track (independent of the channel)
//func (w *writer) addMessage(deltaTime uint32, msg midi.Message) {
func (w *writer) addMessage(deltaTime uint32, raw []byte) {
	//w.printf("adding message deltaTime %v and message %s", deltaTime, msg)
	w.absPos += uint64(deltaTime)
	// we have some sort of sysex, so we need to
	// calculate the length of msg[1:]
	// set msg to msg[0] + length of msg[1:] + msg[1:]
	//raw := msg.Raw()
	if raw[0] == 0xF0 || raw[0] == 0xF7 {
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
