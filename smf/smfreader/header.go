package smfreader

import (
	// "fmt"
	"io"
	// "io/ioutil"
	"github.com/gomidi/midi/internal/lib"
	// "github.com/gomidi/midi"
	// "github.com/gomidi/midi/messages/channel"
	// "github.com/gomidi/midi/messages/meta"
	"github.com/gomidi/midi/smf"
)

// Header data
type mThdData struct {
	format    smf.Format
	numTracks uint16

	// One of MetricalTimeFormat or TimeCodeTimeFormat
	//timeFormat uint
	timeFormat smf.TimeFormat

	// Used if TimeCodeTimeFormat
	// Currently data is not un-packed.
	timeFormatData uint16

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
	_format, err := lib.ReadUint16(reader)

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
		h.timeFormat = smf.QuarterNoteTicks
	} else {
		// TODO: Can't be bothered to implement this bit just now.
		// If you want it, write it!
		h.timeFormatData = division & 0x7FFF
		//h.timeFormat = timeCodeTimeFormat
		h.timeFormat = smf.TimeCode
	}

	return err
}
