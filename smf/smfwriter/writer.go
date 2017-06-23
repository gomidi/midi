package smfwriter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gomidi/midi"
	// "github.com/gomidi/midiwriter"
	"lib"

	"github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
)

type resSmpteFrames struct {
	CodeFormat    int8
	TicksPerFrame int8
}

type tickHeader interface {
	Write(io.Writer) error
	Format() string
	Ticks() uint16
}

var _ tickHeader = resSmpteFrames{}
var _ tickHeader = resQuarterNote(0)

func (f resSmpteFrames) Ticks() uint16 {
	return uint16(f.TicksPerFrame)
}

func (f resSmpteFrames) Format() string {
	if f.CodeFormat == 29 {
		return "SMPTE-30-DropFrame"
	}
	return fmt.Sprintf("SMPTE-%v", f.CodeFormat)
}

func (f resSmpteFrames) Write(w io.Writer) error {
	// multiplication with -1 makes sure that bit 15 is set
	err := binary.Write(w, binary.BigEndian, f.CodeFormat*-1)
	if err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, f.TicksPerFrame)
}

type resQuarterNote uint16

func (q resQuarterNote) Write(w io.Writer) error {
	if q > 32767 {
		q = 32767 // 32767 is the largest possible value, since bit 15 must always be 0
	}
	return binary.Write(w, binary.BigEndian, uint16(q))
}

func (q resQuarterNote) Format() string {
	return "ResQuarterNote"
}

func (q resQuarterNote) Ticks() uint16 {
	return uint16(q)
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
TODO
defaults:
	- store sysex if they are written
	- ignore realtime and System Common or System Real-Time messages if they are written
	- take deltas set via SetDelta
	- no quantization/rounding
options:
  - ignore incoming sysex
  - store incoming realtime and System Common or System Real-Time messages (escape them inside a sysex)
  - ignore SetDelta and measure time instead (autogenerating delta on the way)
  - allow delta quantization/rounding

*/

type Option func(*writer)

// New returns a Writer
// If no options are passed, a single track midi file is written (SMT0).
// Each track must be finished with a meta.EndOfTrack message.
func New(dest io.Writer, opts ...Option) smf.Writer {
	return newWriter(dest, opts...)
}

type header struct {
	chunk      chunk
	MidiFormat uint16
	NumTracks  uint16
	TickHeader tickHeader
}

// <Header Chunk> = <chunk type><length><format><ntrks><division>
func (hc *header) WriteTo(wr io.Writer) (int, error) {
	hc.chunk.typ = [4]byte{byte('M'), byte('T'), byte('h'), byte('d')}
	var bf bytes.Buffer
	binary.Write(&bf, binary.BigEndian, hc.MidiFormat)
	binary.Write(&bf, binary.BigEndian, hc.NumTracks)
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

	err := hc.TickHeader.Write(&bf)
	if err != nil {
		return bf.Len(), err
	}

	/*
	   // uncommented to be possibly "future proof"
	   	if bf.Len() > 6 {
	   		panic(fmt.Sprintf("data too large for header, maxLength: 6, current length: %d", bf.Len()))
	   	}
	*/

	hc.chunk.data = bf.Bytes()

	//hc.chunk.Data = bf.Bytes()
	return hc.chunk.writeTo(wr)
}

type track struct {
	chunk chunk
}

// <Track Chunk> = <chunk type><length><MTrk event>+
func (t *track) WriteTo(wr io.Writer) (int, error) {
	t.chunk.typ = [4]byte{byte('M'), byte('T'), byte('r'), byte('k')}
	return t.chunk.writeTo(wr)
}

// delta is distance in time to last event in this track (independant of channel)
func (t *track) Add(deltaTime uint32, ev []byte) {
	t.chunk.data = append(t.chunk.data, append(lib.VlqEncode(deltaTime), ev...)...)
}

type chunk struct {
	typ  [4]byte
	data []byte
}

func (c *chunk) Type() string {
	var bf bytes.Buffer
	bf.WriteByte(c.typ[0])
	bf.WriteByte(c.typ[1])
	bf.WriteByte(c.typ[2])
	bf.WriteByte(c.typ[3])
	return bf.String()
	//return fmt.Sprintf("%s%s%s%s", c.typ[0], c.typ[1], c.typ[2], c.typ[3])
}

func (c *chunk) writeTo(wr io.Writer) (int, error) {
	length := int32(len(c.data))
	var bf bytes.Buffer
	bf.WriteByte(c.typ[0])
	bf.WriteByte(c.typ[1])
	bf.WriteByte(c.typ[2])
	bf.WriteByte(c.typ[3])
	binary.Write(&bf, binary.BigEndian, length)
	bf.Write(c.data)
	return wr.Write(bf.Bytes())
}

func QuarterNoteTicks(ticks uint16) Option {
	return func(e *writer) {
		e.header.TickHeader = resQuarterNote(ticks)
	}
}

func SMPTE24(ticksPerFrame int8) Option {
	return func(e *writer) {
		e.header.TickHeader = resSmpteFrames{24, ticksPerFrame}
	}
}

func SMPTE25(ticksPerFrame int8) Option {
	return func(e *writer) {
		e.header.TickHeader = resSmpteFrames{25, ticksPerFrame}
	}
}

func SMPTE30DropFrame(ticksPerFrame int8) Option {
	return sMPTE29(ticksPerFrame)
}

func sMPTE29(ticksPerFrame int8) Option {
	return func(e *writer) {
		e.header.TickHeader = resSmpteFrames{29, ticksPerFrame}
	}
}

func SMPTE30(ticksPerFrame int8) Option {
	return func(e *writer) {
		e.header.TickHeader = resSmpteFrames{30, ticksPerFrame}
	}
}

func NumTracks(ntracks uint16) Option {
	return func(e *writer) {
		e.header.NumTracks = ntracks
	}
}

func SMF2() Option {
	return func(e *writer) {
		e.header.MidiFormat = smf.SequentialTracks.Number()
	}
}

func SMF1() Option {
	return func(e *writer) {
		e.header.MidiFormat = smf.MultiTrack.Number()
	}
}

func SMF0() Option {
	return func(e *writer) {
		e.header.MidiFormat = smf.SingleTrack.Number()
	}
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

type writer struct {
	header *header
	wr     io.Writer
	qticks uint16
	// AutoEndTrack    bool
	writeHeader     bool
	currentTrack    *track
	tracksProcessed uint16
	deltatime       uint32
}

// WriteTo writes a midi file to writer
// Pass NumTracks to write multiple tracks (SMF1), otherwise everything will be written
// into a single track (SMF0). However SMF1 can also be enforced with a single track by passing SMF1 as an option
func newWriter(wr io.Writer, opts ...Option) *writer {
	enc := &writer{
		header: &header{
		// MidiFormat: format(10), // not existing, only for checking if it is undefined to be able to set the default
		},
		writeHeader:  true,
		wr:           wr,
		currentTrack: &track{},
	}

	for _, opt := range opts {
		opt(enc)
	}

	if enc.header.NumTracks == 0 {
		enc.header.NumTracks = 1
	}

	// if midiformat is undefined (see above), i.e. not set via options
	// set the default, which is format 0 for one track and format 1 for multitracks
	// if enc.header.MidiFormat == format(10) {
	if enc.header.NumTracks > 1 {
		enc.header.MidiFormat = smf.MultiTrack.Number()
	}
	// }

	if enc.header.TickHeader == nil {
		enc.header.TickHeader = resQuarterNote(960)
	}

	if qn, is := enc.header.TickHeader.(resQuarterNote); is {
		enc.qticks = qn.Ticks()
	}

	return enc
}

func (e *writer) SetDelta(deltatime uint32) {
	e.deltatime = deltatime
}

// WriteEvent writes the header on the first call, if e.writeHeader is true
// in realtime mode, no header and no track is written, instead each event is
// written as is to the output writer until an end of track event had come
// then io.EOF is returned
// WriteEvent returns any writing error or io.EOF if the last track has been written
func (e *writer) Write(m midi.Message) (err error) {
	defer func() {
		e.deltatime = 0
	}()

	if e.writeHeader {
		err = e.WriteHeader()
		if err != nil {
			return err
		}
		e.writeHeader = false
	}
	// fmt.Printf("%T\n", ev)
	if m == meta.EndOfTrack {
		e.currentTrack.Add(e.deltatime, m.Raw())
		_, err = e.currentTrack.WriteTo(e.wr)
		e.tracksProcessed++
		if e.header.NumTracks == e.tracksProcessed {
			return io.EOF
		}
		e.currentTrack = &track{}
		return nil
	}
	e.currentTrack.Add(e.deltatime, m.Raw())
	return nil
}

func (e *writer) WriteHeader() (err error) {
	_, err = e.header.WriteTo(e.wr)
	return
}
