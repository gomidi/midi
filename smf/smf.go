package smf

import (
	"fmt"

	"github.com/gomidi/midi"
)

// Writer writes midi messages to a standard midi file (SMF)
type Writer interface {
	// Writer is also a midi.Writer that writes midi messages
	midi.Writer

	// SetDelta sets a time distance between the last written and the next midi message in ticks.
	// The meaning of a tick depends on the time format that is set in the header of the SMF file.
	SetDelta(ticks uint32)
}

// Reader reads midi messages from a standard midi file (SMF)
type Reader interface {
	// Reader is also a midi.Reader that reads midi messages
	midi.Reader

	// ReadHeader reads the header of SMF file
	// If it is not called, the first call to Read will implicitely read the header.
	// However to get the header information, ReadHeader must be called (which may also happen after the first message read)
	ReadHeader() (Header, error)

	// Delta returns the time distance between the last read midi message and the message before in ticks.
	// The meaning of a tick depends on the time format that is set in the header of the SMF file.
	Delta() (ticks uint32)

	// Track returns the number of the track of the last read midi message (starting with 0)
	// It returns -1 if no message has been read yet.
	Track() int16
}

type Header struct {
	Format
	NumTracks uint16
	TimeFormat
}

/*
// Header represents the header of a SMF file
type Header interface {
	// Format returns the SMF format (0 = SingleTrack, 1 = MultiTrack, 2 = SequentialTracks)
	Format() Format

	// TimeFormat returns the time format (MetricResolution or TimeCode)
	// To get the value, type cast to MetricResolution or TimeCode
	TimeFormat() TimeFormat

	// NumTracks returns the number of tracks as defined inside the SMF header. It should be the same
	// as the real number of tracks in the file, although there is no guaranty.
	NumTracks() uint16
}
*/

const (
	// SMF0 represents the singletrack SMF format (0)
	SMF0 = format(0)

	// SMF1 represents the multitrack SMF format (1)
	SMF1 = format(1)

	// SMF2 represents the sequential track SMF format (2)
	SMF2 = format(2)
)

var (
	_ TimeFormat = MetricResolution(0)
	_ TimeFormat = TimeCode{}
)

type TimeCode struct {
	FramesPerSecond uint8
	SubFrames       uint8
}

func (t TimeCode) String() string {

	switch t.FramesPerSecond {
	case 29:
		return fmt.Sprintf("SMPTE30DropFrame %v subframes", t.SubFrames)
	default:
		return fmt.Sprintf("SMPTE%v %v subframes", t.FramesPerSecond, t.SubFrames)
	}

}

func (t TimeCode) timeformat() {}

// SMPTE24 returns a SMPTE24 TimeCode with the given subframes
func SMPTE24(subframes uint8) TimeCode {
	return TimeCode{24, subframes}
}

// SMPTE25 returns a SMPTE25 TimeCode with the given subframes
func SMPTE25(subframes uint8) TimeCode {
	return TimeCode{25, subframes}
}

// SMPTE30DropFrame returns a SMPTE30 drop frame TimeCode with the given subframes
func SMPTE30DropFrame(subframes uint8) TimeCode {
	return TimeCode{29, subframes}
}

// SMPTE30 returns a SMPTE30 TimeCode with the given subframes
func SMPTE30(subframes uint8) TimeCode {
	return TimeCode{30, subframes}
}

// MetricResolution represents the "ticks per quarter note" (metric) time format
// It defaults to 960 (i.e. 0 is treated as if it where 960 ticks per quarter note)
type MetricResolution uint16

// Ticks returns the ticks for a quarter note (defaults to 960)
func (q MetricResolution) Ticks() uint16 {
	if uint16(q) == 0 {
		return 960 // default
	}
	return uint16(q)
}

func (q MetricResolution) div(d float64) uint32 {
	return uint32(roundFloat(float64(q.Ticks())/d, 0))
}

// N4 returns the ticks for a quarter note
func (q MetricResolution) N4() uint32 {
	return uint32(q.Ticks())
}

// N8 returns the ticks for a 8th note
func (q MetricResolution) N8() uint32 {
	return q.div(2)
}

// N16 returns the ticks for a 16th note
func (q MetricResolution) N16() uint32 {
	return q.div(4)
}

// N32 returns the ticks for a 32th note
func (q MetricResolution) N32() uint32 {
	return q.div(8)
}

// N64 returns the ticks for a 64th note
func (q MetricResolution) N64() uint32 {
	return q.div(16)
}

// N128 returns the ticks for a 128th note
func (q MetricResolution) N128() uint32 {
	return q.div(32)
}

// N256 returns the ticks for a 256th note
func (q MetricResolution) N256() uint32 {
	return q.div(64)
}

// N512 returns the ticks for a 512th note
func (q MetricResolution) N512() uint32 {
	return q.div(128)
}

// N1024 returns the ticks for a 1024th note
func (q MetricResolution) N1024() uint32 {
	return q.div(256)
}

// N2 returns the ticks for a half note
func (q MetricResolution) N2() uint32 {
	return q.N4() * 2
}

// String returns the string representation of the quarter note resolution
func (q MetricResolution) String() string {
	return fmt.Sprintf("%v MetricResolution", q.Ticks())
}

func (q MetricResolution) timeformat() {}

// Format is the common interface of all SMF file formats
type Format interface {
	String() string
	Number() uint16
	smfformat() // make the implementation exclusive to this package
}

// TimeFormat is the common interface of all SMF time formats
type TimeFormat interface {
	String() string
	timeformat() // make the implementation exclusive to this package
}

// format is an implementation of Format
type format uint16

func (f format) Number() uint16 {
	return uint16(f)
}

func (f format) smfformat() {}

func (f format) String() string {
	switch f {
	case SMF0:
		return "SMF0 (singletrack)"
	case SMF1:
		return "SMF1 (multitrack)"
	case SMF2:
		return "SMF2 (sequential tracks)"
	}
	panic("unreachable")
}

// timeformat is an implementation of TimeFormat
type timeformat string

func (t timeformat) String() string { return string(t) }
func (t timeformat) timeformat()    {}
