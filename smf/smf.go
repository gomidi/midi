package smf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/gomidi/midi/internal/midilib"

	"github.com/gomidi/midi"
)

var (
	_ midi.Writer = Writer(nil)
	_ midi.Reader = Reader(nil)
)

// Writer writes midi messages to a standard midi file (SMF)
// Writer is also a midi.Writer
type Writer interface {

	// Header returns the header
	Header() Header

	// WriteHeader writes the midi header
	// If WriteHeader was not called before the first run of Write,
	// it will implicitly be called when calling Write.
	WriteHeader() error

	// Write writes a midi message to the SMF file.
	//
	// Due to the nature of SMF files there is some maybe surprising behavior.
	// - If the header has not been written yet, it will be written before writing the first message.
	// - The first message will be written to track 0 which will be implicitly created.
	// - All messages of a track will be buffered inside the track and only be written if an EndOfTrack
	//   message is written.
	// - The number of tracks that are written will never execeed the NumTracks that have been defined when creating the writer.
	//   If the last track has been written, io.EOF will be returned. (Also for any further attempt to write).
	// - It is the responsibility of the caller to make sure the provided NumTracks (which defaults to 1) is not
	//   larger as the number of tracks in the file.
	// Any error stops the writing, is tracked and prohibits further writing.
	// At the end smf.ErrFinished will be returned
	Write(midi.Message) error

	// SetDelta sets a time distance between the last written and the following message in ticks.
	// The meaning of a tick depends on the time format that is set in the header of the SMF file.
	// Use Header.TimeFormat.(MetricTicks).Ticks*th() to get the ticks of quarter notes etc.
	SetDelta(ticks uint32)
}

// Reader reads midi messages from a standard midi file (SMF)
// Reader is also a midi.Reader
type Reader interface {

	// ReadHeader reads the header of the SMF file. If Header is called before ReadHeader, it will panic.
	// ReadHeader is also implicitly called with the first call of Read() (if it has not been run before)
	ReadHeader() error

	// Read reads a MIDI message from a SMF file.
	// any error will be tracked and stops reading and prevents any other attempt to read.
	// At the end, smf.ErrFinished will be returned.
	Read() (midi.Message, error)

	// Header returns the header of SMF file
	// if the header is not yet read, it will be read before
	Header() Header

	// Delta returns the time distance between the last read midi message and the message before in ticks.
	// The meaning of a tick depends on the time format that is set in the header of the SMF file.
	// Use Header.TimeFormat.(MetricTicks).In64ths() to convert the ticks to 64ths notes. From there you may calculate
	// 16ths,quaver,4ths etc by dividing by 4,8,16 etc.
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

func (h Header) String() string {
	return fmt.Sprintf("<Format: %v, NumTracks: %v, TimeFormat: %v>", h.Format, h.NumTracks, h.TimeFormat)
}

const (
	// SMF0 represents the singletrack SMF format (0)
	SMF0 = format(0)

	// SMF1 represents the multitrack SMF format (1)
	SMF1 = format(1)

	// SMF2 represents the sequential track SMF format (2)
	SMF2 = format(2)
)

// Chunk is a chunk of a SMF file
type Chunk struct {
	typ  []byte // must always be 4 bytes long, to avoid conversions everytime, we take []byte here instead of [4]byte
	data []byte
}

// Len returns the length of the chunk body
func (c *Chunk) Len() int {
	return len(c.data)
}

// SetType sets the type of the chunk
func (c *Chunk) SetType(typ [4]byte) {
	c.typ = make([]byte, 4)
	c.typ[0] = typ[0]
	c.typ[1] = typ[1]
	c.typ[2] = typ[2]
	c.typ[3] = typ[3]
}

// Type returns the type of the chunk (from the header)
func (c *Chunk) Type() string {
	var bf bytes.Buffer
	bf.Write(c.typ)
	return bf.String()
}

// Clear removes all data but keeps the type
func (c *Chunk) Clear() {
	c.data = nil
}

// WriteTo writes the content of the chunk to the given writer
func (c *Chunk) WriteTo(wr io.Writer) (int64, error) {
	if len(c.typ) != 4 {
		return 0, fmt.Errorf("chunk header not set properly")
	}

	var bf bytes.Buffer
	bf.Write(c.typ)
	binary.Write(&bf, binary.BigEndian, int32(c.Len()))
	bf.Write(c.data)
	n, err := wr.Write(bf.Bytes())
	if err != nil {
		return int64(n), fmt.Errorf("could not write chunk: %v", err)
	}
	return int64(n), nil
}

// ReadHeader reads the header from the given reader
// returns the length of the following body
// for errors, length of 0 is returned
func (c *Chunk) ReadHeader(rd io.Reader) (length uint32, err error) {
	c.typ, err = midilib.ReadNBytes(4, rd)

	if err != nil {
		c.typ = nil
		return
	}

	return midilib.ReadUint32(rd)
}

// Write writes the given bytes to the body of the chunk
func (c *Chunk) Write(b []byte) (int, error) {
	c.data = append(c.data, b...)
	return len(b), nil
}

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

// In64ths returns the deltaTicks in 64th notes.
// To get 32ths, divide result by 2
// To get 16ths, divide result by 4
// To get 8ths, divide result by 8
// To get 4ths, divide result by 16
func (q MetricTicks) In64ths(deltaTicks uint32) uint32 {
	return (deltaTicks * 16) / uint32(q)
}

// Duration returns the time.Duration for a number of ticks at a certain tempo (in BPM)
func (q MetricTicks) Duration(tempoBPM uint32, deltaTicks uint32) time.Duration {
	// (60000 / T) * (d / R) = D[ms]
	//	durQnMilli := 60000 / float64(tempoBPM)
	//	_4thticks := float64(deltaTicks) / float64(uint16(q))
	res := 60000000000 * float64(deltaTicks) / (float64(tempoBPM) * float64(uint16(q)))
	//fmt.Printf("what: %vns\n", res)
	return time.Duration(int64(math.Round(res)))
	//	return time.Duration(roundFloat(durQnMilli*_4thticks, 0)) * time.Millisecond
}

// FractionalDuration returns the time.Duration for a number of ticks at a certain tempo (in fractional BPM)
func (q MetricTicks) FractionalDuration(fractionalBPM float64, deltaTicks uint32) time.Duration {
	// (60000 / T) * (d / R) = D[ms]
	//	durQnMilli := 60000 / float64(tempoBPM)
	//	_4thticks := float64(deltaTicks) / float64(uint16(q))
	res := 60000000000 * float64(deltaTicks) / (fractionalBPM * float64(uint16(q)))
	//fmt.Printf("what: %vns\n", res)
	return time.Duration(int64(math.Round(res)))
	//	return time.Duration(roundFloat(durQnMilli*_4thticks, 0)) * time.Millisecond
}

// Ticks returns the ticks for a given time.Duration at a certain tempo (in BPM)
func (q MetricTicks) Ticks(tempoBPM uint32, d time.Duration) (ticks uint32) {
	// d = (D[ms] * R * T) / 60000
	ticks = uint32(roundFloat((float64(d.Nanoseconds())/1000000*float64(uint16(q))*float64(tempoBPM))/60000, 0))
	return ticks
}

// FractionalTicks returns the ticks for a given time.Duration at a certain tempo (in fractional BPM)
func (q MetricTicks) FractionalTicks(fractionalBPM float64, d time.Duration) (ticks uint32) {
	// d = (D[ms] * R * T) / 60000
	ticks = uint32(roundFloat((float64(d.Nanoseconds())/1000000*float64(uint16(q))*fractionalBPM)/60000, 0))
	return ticks
}

func (q MetricTicks) div(d float64) uint32 {
	return uint32(roundFloat(float64(q.Number())/d, 0))
}

// Number returns the number of the metric ticks (ticks for a quarter note, defaults to 960)
func (q MetricTicks) Number() uint16 {
	if q == 0 {
		return 960 // default
	}
	return uint16(q)
}

// Ticks4th returns the ticks for a quarter note
func (q MetricTicks) Ticks4th() uint32 {
	return uint32(q.Number())
}

// Ticks8th returns the ticks for a quaver note
func (q MetricTicks) Ticks8th() uint32 {
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

// String returns the string representation of the quarter note resolution
func (q MetricTicks) String() string {
	return fmt.Sprintf("%v MetricTicks", q.Ticks4th())
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
