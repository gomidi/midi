package smfreader

import (
	"github.com/gomidi/midi/internal/midilib"
	"github.com/gomidi/midi/smf"
	"io"
)

// Header data
type mThdData struct {
	format    smf.Format
	numTracks uint16

	// One of MetricalTimeFormat or TimeCodeTimeFormat
	//timeFormat uint
	timeFormat smf.TimeFormat

	// Used if TimeCodeTimeFormat
	// the raw timeformat data
	// unpack it with UnpackTimeCode
	timeFormatData uint16

	timeCodeSubFrames uint8

	// Used if MetricalTimeFormat
	quarterNoteTicks uint16
}

func (p mThdData) Format() smf.Format {
	return p.format
}

func (p mThdData) NumTracks() uint16 {
	return p.numTracks
}

func (p mThdData) TimeFormat() (smf.TimeFormat, uint16) {
	if p.timeFormat == smf.QuarterNoteTicks {
		return p.timeFormat, p.quarterNoteTicks
	}

	return smf.TimeCode, p.timeFormatData
}

// parseHeaderData parses SMF-header chunk header data.
// It returns the ChunkHeader struct as a value and an error.
func (h *mThdData) readFrom(reader io.Reader) error {
	// Format
	_format, err := midilib.ReadUint16(reader)

	if err != nil {
		return err
	}

	switch _format {
	case 0:
		h.format = smf.SingleTrack
	case 1:
		h.format = smf.MultiTrack
	case 2:
		h.format = smf.SequentialTracks
	default:
		return ErrUnsupportedSMFFormat
	}

	// Num tracks
	h.numTracks, err = midilib.ReadUint16(reader)

	if err != nil {
		return err
	}
	// Division
	var division uint16
	division, err = midilib.ReadUint16(reader)

	// "If bit 15 of <division> is zero, the bits 14 thru 0 represent the number
	// of delta time "ticks" which make up a quarter-note."
	if division&0x8000 == 0x0000 {
		h.quarterNoteTicks = division & 0x7FFF
		//h.timeFormat = metricalTimeFormat
		h.timeFormat = smf.QuarterNoteTicks
	} else {
		// TODO: Can't be bothered to implement this bit just now.
		// If you want it, write it!
		// h.timeFormatData = division & 0x7FFF
		h.timeFormatData = division
		// h.smpteFPS = uint8(int8(byte(division>>8)) * (-1)) // bit shifting first byte to second inverting sign
		// h.timeCodeSubFrames = byte(division & uint16(255)) // taking the second byte

		//h.timeFormatData = division
		//h.timeFormat = timeCodeTimeFormat
		h.timeFormat = smf.TimeCode
	}

	/*
			The last two bytes indicate how many Pulses (i.e. clocks) Per Quarter Note
			(abbreviated as PPQN) resolution the time-stamps are based upon, Division.
			For example, if your sequencer has 96 ppqn, this field would be (in hex):

		00 60

		Alternately, if the first byte of Division is negative, then this represents
		the division of a second that the time-stamps are based upon. The first byte
		will be -24, -25, -29, or -30, corresponding to the 4 SMPTE standards
		representing frames per second. The second byte (a positive number)
		is the resolution within a frame (ie, subframe). Typical values may
		be 4 (MIDI Time Code), 8, 10, 80 (SMPTE bit resolution), or 100.

		You can specify millisecond-based timing by the data bytes of -25 and 40 subframes.
	*/

	/* http://www.somascape.org/midi/tech/mfile.html

	tickdiv : specifies the timing interval to be used, and whether timecode (Hrs.Mins.Secs.Frames) or metrical (Bar.Beat) timing is to be used. With metrical timing, the timing interval is tempo related, whereas with timecode the timing interval is in absolute time, and hence not related to tempo.

	    Bit 15 (the top bit of the first byte) is a flag indicating the timing scheme in use :

	    Bit 15 = 0 : metrical timing
	    Bits 0 - 14 are a 15-bit number indicating the number of sub-divisions of a quarter note (aka pulses per quarter note, ppqn). A common value is 96, which would be represented in hex as 00 60. You will notice that 96 is a nice number for dividing by 2 or 3 (with further repeated halving), so using this value for tickdiv allows triplets and dotted notes right down to hemi-demi-semiquavers to be represented.

	    Bit 15 = 1 : timecode
	    Bits 8 - 15 (i.e. the first byte) specifies the number of frames per second (fps),
	    and will be one of the four SMPTE standards - 24, 25, 29 or 30, though expressed as a negative value
	    (using 2's complement notation), as follows :
	    fps	Representation (hex)
	    24 E8
	    25 E7
	    29 E3
	    30 E2


	    Bits 0 - 7 (the second byte) specifies the sub-frame resolution, i.e. the number of sub-divisions of a frame.
	    Typical values are 4 (corresponding to MIDI Time Code), 8, 10, 80 (corresponding to SMPTE bit resolution), or 100.

	    A timing resolution of 1 ms can be achieved by specifying 25 fps and 40 sub-frames, which would be encoded in hex as  E7 28.

	A complete MThd chunk thus contains 14 bytes (including the 8 byte header).
	Example
	Data (hex)	Interpretation
	4D 54 68 64 	identifier, the ascii chars 'MThd'
	00 00 00 06 	chunklen, 6 bytes of data follow . . .
	00 01 	format = 1
	00 11 	ntracks = 17
	00 60 	tickdiv = 96 ppqn, metrical time

	*/

	return err
}
