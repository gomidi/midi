package smf

import (
	"fmt"

	"github.com/gomidi/midi"
)

var (
	_ midi.Writer = Writer(nil)
	_ midi.Reader = Reader(nil)
)

// Writer writes midi messages to a standard midi file (SMF)
// Writer is also a midi.Writer
type Writer interface {

	// Write writes a midi message to the SMF file.
	//
	// Due to the nature of SMF files there is some maybe surprising behavior.
	// - If the header has not been written yet, it will be written before writing the first message.
	// - The first message will be written to track 0 which will be implicetly created.
	// - All messages of a track will be buffered inside the track and only be written if an EndOfTrack
	//   message is written.
	// - The number of tracks that are written will never execeed the NumTracks that have been defined when creating the writer.
	//   If the last track has been written, io.EOF will be returned. (Also for any further attempt to write).
	// - It is the responsability of the caller to make sure the provided NumTracks (which defaults to 1) is not
	//   larger as the number of tracks in the file.
	// Keep the above in mind when examinating the written nbytes that are returned. They reflect the number of bytes
	// that have been physically written at that point in time.
	Write(midi.Message) (nBytes int, err error)

	// SetDelta sets a time distance between the last written and the following message in ticks.
	// The meaning of a tick depends on the time format that is set in the header of the SMF file.
	SetDelta(ticks uint32)
}

// Reader reads midi messages from a standard midi file (SMF)
// Reader is also a midi.Reader
type Reader interface {

	// Read reads a MIDI message from a SMF file.
	Read() (midi.Message, error)

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

// Header represents the header of a SMF file.
type Header struct {

	// Format is the SMF file format: SMF0, SMF1 or SMF2
	Format

	// NumTracks is the number of tracks (always > 0)
	NumTracks uint16

	// TimeFormat is the time format (either MetricTicks or TimeCode)
	TimeFormat
}

const (
	// SMF0 represents the singletrack SMF format (0)
	SMF0 = format(0)

	// SMF1 represents the multitrack SMF format (1)
	SMF1 = format(1)

	// SMF2 represents the sequential track SMF format (2)
	SMF2 = format(2)
)

var (
	_ TimeFormat = MetricTicks(0)
	_ TimeFormat = TimeCode{}
	_ Format     = SMF0
)

// TimeCode is the SMPTE time format.
// It can be comfortable created with the SMPTE* functions.
type TimeCode struct {
	FramesPerSecond uint8
	SubFrames       uint8
}

// String represents the TimeCode as a string.
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

// MetricTicks represents the "ticks per quarter note" (metric) time format
// It defaults to 960 (i.e. 0 is treated as if it where 960 ticks per quarter note)
type MetricTicks uint16

// Ticks returns the ticks for a quarter note (defaults to 960)
func (q MetricTicks) Ticks() uint16 {
	if uint16(q) == 0 {
		return 960 // default
	}
	return uint16(q)
}

func (q MetricTicks) div(d float64) uint32 {
	return uint32(roundFloat(float64(q.Ticks())/d, 0))
}

// TicksQuarter returns the ticks for a quarter note
func (q MetricTicks) TicksQuarter() uint32 {
	return uint32(q.Ticks())
}

// TicksQuaver returns the ticks for a quaver note
func (q MetricTicks) TicksQuaver() uint32 {
	return q.div(2)
}

// Ticks16th returns the ticks for a 16th note
func (q MetricTicks) Ticks16th() uint32 {
	return q.div(4)
}

// Ticks32th returns the ticks for a 32th note
func (q MetricTicks) Ticks32th() uint32 {
	return q.div(8)
}

// Ticks64th returns the ticks for a 64th note
func (q MetricTicks) Ticks64th() uint32 {
	return q.div(16)
}

// Ticks128th returns the ticks for a 128th note
func (q MetricTicks) Ticks128th() uint32 {
	return q.div(32)
}

// Ticks256th returns the ticks for a 256th note
func (q MetricTicks) Ticks256th() uint32 {
	return q.div(64)
}

// Ticks512th returns the ticks for a 512th note
func (q MetricTicks) Ticks512th() uint32 {
	return q.div(128)
}

// Ticks1024th returns the ticks for a 1024th note
func (q MetricTicks) Ticks1024th() uint32 {
	return q.div(256)
}

// TicksHalf returns the ticks for a half note
func (q MetricTicks) TicksHalf() uint32 {
	return q.TicksQuarter() * 2
}

// TicksWhole returns the ticks for a whole note
func (q MetricTicks) TicksWhole() uint32 {
	return q.TicksQuarter() * 4
}

// String returns the string representation of the quarter note resolution
func (q MetricTicks) String() string {
	return fmt.Sprintf("%v MetricTicks", q.Ticks())
}

func (q MetricTicks) timeformat() {}

// Format is the common interface of all SMF file formats
type Format interface {

	// String returns the string representation of the SMF format.
	String() string

	// Type returns the type of the SMF file: 0 for SMF0, 1 for SMF1 and 2 for SMF2
	Type() uint16

	smfformat() // make the implementation exclusive to this package
}

// TimeFormat is the common interface of all SMF time formats
type TimeFormat interface {
	String() string
	timeformat() // make the implementation exclusive to this package
}

// format is an implementation of Format
type format uint16

func (f format) Type() uint16 {
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
