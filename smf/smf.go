package smf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"lib"
)

/*
   A MIDI message is made up of an eight-bit status byte which is generally followed by one or two data bytes.

   MIDI message (status byte + 1-2 data bytes)
      |
      -------- Channel Message (channel number included in status byte) 1000 | 1001 | 1010 | 1011 | 1100 | 1101 | 1110
      |            |
      |            ------ Channel Voice Message (musical performance)
      |            |
      |            ------ Mode Message  (how to does instr. respond to Channel Voice message)
      |                   (1011nnnn Channel Mode Message)
      |
      ---------System Message (no channel number in status byte), all beginning with 1111
                   |
                   ------ System Common Messages
                   |
                   ------ System Real Time Messages
                   |
                   ------ System Exclusive Messages (F0, F7)

   There are a number of different types of MIDI messages. At the highest level, MIDI messages are classified
   as being either Channel Messages or System Messages.

   Channel messages are those which apply to a specific
   Channel, and the Channel number is included in the status byte for these messages. System messages are not
   Channel specific, and no Channel number is indicated in their status bytes.

   Channel Messages may be further classified as being either Channel Voice Messages, or Mode Messages.
   Channel Voice Messages carry musical performance data, and these messages comprise most of the traffic in
   a typical MIDI data stream. Channel Mode messages affect the way a receiving instrument will respond to the
   Channel Voice messages.

   MIDI System Messages are classified as being System Common Messages, System Real Time Messages, or
   System Exclusive Messages. System Common messages are intended for all receivers in the system. System
   Real Time messages are used for synchronisation between clock-based MIDI components. System Exclusive
   messages include a Manufacturer's Identification (ID) code, and are used to transfer any number of data bytes
   in a format specified by the referenced manufacturer.
*/

/*
			   read in the next byte (uint8)

			   if it is FF -> meta event
			   if it is F0 or F7 -> sysex event
			   else ->
			      System Common Message F0-FF

			      F0 1111 0000 sysex event
			      F7 1111 0111 sysex event
						FF 1111 1111 meta event

			   		channel voice message D7-D0

			   		D0 1101 0000
	          D7 1101 0111


*/

const (
	SingleTrack      = format(0)
	MultiTrack       = format(1)
	SequentialTracks = format(2)

	QuarterNoteTicks = timeformat("QuarterNoteTicks")
	TimeCodeTicks    = timeformat("TimeCodeTicks")
)

type format uint16

func (f format) String() string {

	switch f {
	case SingleTrack:
		return "SingleTrack"
	case MultiTrack:
		return "MultiTrack"
	case SequentialTracks:
		return "SequentialTracks"
	}

	panic("unreachable")

}

type timeformat string

func (t timeformat) String() string {
	return string(t)
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

// A chunk header
type chunkHeader struct {
	typ    string
	length uint32
}

// parseChunkHeader parses a chunk header from a ReadSeeker.
// It returns the ChunkHeader struct as a value and an error.
func (c *chunkHeader) readFrom(rd io.Reader) error {
	// fmt.Println("Parse Chunk Header")
	b, err := lib.ReadN(4, rd)

	if err != nil {
		return err
	}

	c.length, err = lib.ReadUint32(rd)
	c.typ = string(b)

	// parseUint32 might return an error.
	return err
}

// Header data
type mThdData struct {
	format    format
	numTracks uint16

	// One of MetricalTimeFormat or TimeCodeTimeFormat
	//timeFormat uint
	timeFormat timeformat

	// Used if TimeCodeTimeFormat
	// Currently data is not un-packed.
	timeFormatData uint16

	// Used if MetricalTimeFormat
	quarterNoteTicks uint16
}

func (p mThdData) Format() format {
	return p.format
}

func (p mThdData) NumTracks() uint16 {
	return p.numTracks
}

func (p mThdData) TimeFormat() (timeformat, uint16) {
	if p.timeFormat == QuarterNoteTicks {
		return p.timeFormat, p.quarterNoteTicks
	}

	return TimeCodeTicks, p.timeFormatData
}

// parseHeaderData parses SMF-header chunk header data.
// It returns the ChunkHeader struct as a value and an error.
func (h *mThdData) readFrom(reader io.Reader) error {
	// Format
	_format, err := lib.ReadUint16(reader)

	if err != nil {
		return err
	}

	// Should be one of 0, 1, 2
	if _format > 2 {
		return ErrUnsupportedSMFFormat
	}

	h.format = format(uint8(_format))

	// Num tracks
	h.numTracks, err = lib.ReadUint16(reader)

	if err != nil {
		return err
	}
	// Division
	var division uint16
	division, err = lib.ReadUint16(reader)

	// "If bit 15 of <division> is zero, the bits 14 thru 0 represent the number
	// of delta time "ticks" which make up a quarter-note."
	if division&0x8000 == 0x0000 {
		h.quarterNoteTicks = division & 0x7FFF
		//h.timeFormat = metricalTimeFormat
		h.timeFormat = QuarterNoteTicks
	} else {
		// TODO: Can't be bothered to implement this bit just now.
		// If you want it, write it!
		h.timeFormatData = division & 0x7FFF
		//h.timeFormat = timeCodeTimeFormat
		h.timeFormat = TimeCodeTicks
	}

	return err
}

type tickHeader interface {
	Write(io.Writer) error
	Format() string
	Ticks() uint16
}

type smpteFrames struct {
	codeFormat    int8
	ticksPerFrame int8
}

var _ tickHeader = smpteFrames{}
var _ tickHeader = quarterNote(0)

func (f smpteFrames) Ticks() uint16 {
	return uint16(f.ticksPerFrame)
}

func (f smpteFrames) Format() string {
	if f.codeFormat == 29 {
		return "SMPTE-30-DropFrame"
	}
	return fmt.Sprintf("SMPTE-%v", f.codeFormat)
}

func (f smpteFrames) Write(w io.Writer) error {
	// multiplication with -1 makes sure that bit 15 is set
	err := binary.Write(w, binary.BigEndian, f.codeFormat*-1)
	if err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, f.ticksPerFrame)
}

type quarterNote uint16

func (q quarterNote) Write(w io.Writer) error {
	if q > 32767 {
		q = 32767 // 32767 is the largest possible value, since bit 15 must always be 0
	}
	return binary.Write(w, binary.BigEndian, uint16(q))
}

func (q quarterNote) Format() string {
	return "QuarterNote"
}

func (q quarterNote) Ticks() uint16 {
	return uint16(q)
}

type header struct {
	chunk      chunk
	MidiFormat format
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
