package smfwriter

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gomidi/midi/smf"
	"io"
)

type header struct {
	chunk      chunk
	MidiFormat smf.Format
	NumTracks  uint16
	TimeFormat smf.TimeFormat
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
func (hc *header) writeTimeFormat(wr io.Writer) error {
	switch tf := hc.TimeFormat.(type) {
	case smf.QuarterNoteTicks:
		ticks := tf.Ticks()
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
		panic(fmt.Sprintf("unsupported TimeFormat: %#v", hc.TimeFormat))
	}
}

// <Header Chunk> = <chunk type><length><format><ntrks><division>
func (hc *header) WriteTo(wr io.Writer) (int, error) {
	hc.chunk.typ = [4]byte{byte('M'), byte('T'), byte('h'), byte('d')}
	var bf bytes.Buffer
	if hc.NumTracks == 0 {
		hc.NumTracks = 1
	}
	if hc.MidiFormat == nil {
		if hc.NumTracks == 1 {
			hc.MidiFormat = smf.SMF0
		} else {
			hc.MidiFormat = smf.SMF1
		}
	}
	binary.Write(&bf, binary.BigEndian, hc.MidiFormat.Number())
	binary.Write(&bf, binary.BigEndian, hc.NumTracks)

	err := hc.writeTimeFormat(&bf)
	if err != nil {
		return bf.Len(), err
	}

	hc.chunk.data = bf.Bytes()

	return hc.chunk.writeTo(wr)
}
