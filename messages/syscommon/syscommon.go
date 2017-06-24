package syscommon

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gomidi/midi/internal/lib"
)

var (
	_ Message = SongPositionPointer(0)
	_ Message = SongSelect(0)
	_ Message = TuneRequest
	_ Message = Undefined4(0)
	_ Message = Undefined5(0)
	_ Message = MIDITimingCode(0)
)

type Message interface {
	String() string
	Raw() []byte
	readFrom(io.Reader) (Message, error)
	sysCommon()
}

/*
MTC Quarter Frame

These are the MTC (i.e. SMPTE based) equivalent of the F8 Timing Clock messages, though offer much higher resolution, as they are sent at a rate of 96 to 120 times a second (depending on the SMPTE frame rate). Each Quarter Frame message provides partial timecode information, 8 sequential messages being required to fully describe a timecode instant. The reconstituted timecode refers to when the first partial was received. The most significant nibble of the data byte indicates the partial (aka Message Type).
Partial	Data byte	Usage
1	0000 bcde	Frame number LSBs 	abcde = Frame number (0 to frameRate-1)
2	0001 000a	Frame number MSB
3	0010 cdef	Seconds LSBs 	abcdef = Seconds (0-59)
4	0011 00ab	Seconds MSBs
5	0100 cdef	Minutes LSBs 	abcdef = Minutes (0-59)
6	0101 00ab	Minutes MSBs
7	0110 defg	Hours LSBs 	ab = Frame Rate (00 = 24, 01 = 25, 10 = 30drop, 11 = 30nondrop)
cdefg = Hours (0-23)
8	0111 0abc	Frame Rate, and Hours MSB
*/

/*


  B.2. System common messages:

    contain the following unrelated messages

System Common Message   Status Byte      Number of Data Bytes
---------------------   -----------      --------------------
MIDI Timing Code            F1                   1
Song Position Pointer       F2                   2
Song Select                 F3                   1
Tune Request                F6                  None

*/

type MIDITimingCode uint8

func (m MIDITimingCode) QuarterFrame() uint8 {
	return uint8(m)
}

func (m MIDITimingCode) readFrom(rd io.Reader) (Message, error) {
	b, err := lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return MIDITimingCode(b), nil
}

func (m MIDITimingCode) sysCommon() {}

func (m MIDITimingCode) String() string {
	return fmt.Sprintf("%T: %v", m, uint8(m))
}

func (m MIDITimingCode) Raw() []byte {
	// TODO check - it is a guess
	return []byte{byte(0xF1), byte(m)}
}

/*
	statusSysEx                     = byte(0xF0)
	statusMIDITimingCodeMessage     = byte(0xF1)
	statusSongPositionPointer       = byte(0xF2)
	statusSongSelect                = byte(0xF3)
	statusTuneRequest               = byte(0xF6)
*/

// if canary >= 0xF0 && canary <= 0xF7 {
const (
	byteMIDITimingCodeMessage  = byte(0xF1)
	byteSysSongPositionPointer = byte(0xF2)
	byteSysSongSelect          = byte(0xF3)
	byteSysTuneRequest         = byte(0xF6)
	// byteSysEx                     = byte(0xF0)
)

var systemMessages = map[byte]Message{
	byteMIDITimingCodeMessage:  MIDITimingCode(0),
	byteSysSongPositionPointer: SongPositionPointer(0),
	byteSysSongSelect:          SongSelect(0),
	byteSysTuneRequest:         TuneRequest,
	0xF4:                       Undefined4(0), // unused (ignore them)
	0xF5:                       Undefined5(0), // unused (ignore them)
}

// Reader read a syscommon
type Reader interface {
	// Read reads a single syscommon message.
	// It may just be called once per Reader. A second call returns io.EOF
	Read() (Message, error)
}

// NewReader returns a reader that can read a single syscommon message
// Read may just be called once per Reader. A second call returns io.EOF
func NewReader(input io.Reader, typ byte) Reader {
	return &reader{input, typ, false}
}

type reader struct {
	input io.Reader
	typ   byte
	done  bool
}

// Read may just be called once per Reader. A second call returns io.EOF
func (r *reader) Read() (msg Message, err error) {
	if r.done {
		return nil, io.EOF
	}
	return Dispatch(r.typ).readFrom(r.input)
}

func (m SongPositionPointer) readFrom(rd io.Reader) (Message, error) {

	bt, err := lib.ReadN(2, rd)

	if err != nil {
		return nil, err
	}

	if len(bt) != 2 {
		err = lib.UnexpectedMessageLengthError("SongPositionPointer expected length 2")
		return nil, err
	}

	// TODO: check if it is correct
	val := uint16((bt[1])&0x7f) << 7
	val |= uint16(bt[0]) & 0x7f

	return SongPositionPointer(val), nil
}

type SongPositionPointer uint16

func (m SongPositionPointer) Number() uint16 {
	return uint16(m)
}

func (m SongPositionPointer) String() string {
	return fmt.Sprintf("%T: %v", m, uint16(m))
}

// TODO test
func (m SongPositionPointer) Raw() []byte {
	// TODO check - it is totally a guess at the moment

	r := lib.MsbLsbUnsigned(uint16(m))

	var bf bytes.Buffer
	//	binary.Write(&bf, binary.BigEndian, uint16(change))
	binary.Write(&bf, binary.BigEndian, 0xF2)

	binary.Write(&bf, binary.BigEndian, r)
	return bf.Bytes()
}
func (m SongPositionPointer) sysCommon() {}

func (m SongSelect) Pos() uint16 {
	return uint16(m)
}

// TODO Test
func (m SongSelect) Raw() []byte {
	// TODO check - it is a guess
	return []byte{byte(0xF3), byte(m)}
}

type SongSelect uint8

func (m SongSelect) Number() uint8 {
	return uint8(m)
}

func (m SongSelect) String() string {
	return fmt.Sprintf("%T: %v", m, uint8(m))
}

func (m SongSelect) sysCommon() {}

// TODO: check
func (m SongSelect) readFrom(rd io.Reader) (Message, error) {

	b, err := lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return SongSelect(b), nil
}

type tuneRequest bool

func Dispatch(b byte) Message {
	return systemMessages[b]
}

const (
	TuneRequest = tuneRequest(false)
)

func (m tuneRequest) meta() {}

func (m tuneRequest) String() string {
	return fmt.Sprintf("%T", m)
}

func (m tuneRequest) readFrom(rd io.Reader) (Message, error) {
	return m, nil
}

func (m tuneRequest) sysCommon() {}

// TODO test
func (m tuneRequest) Raw() []byte {
	return []byte{byte(0xF6)}
}

type Undefined4 int

func (m Undefined4) String() string {
	return fmt.Sprintf("%T", m)
}

func (m Undefined4) readFrom(rd io.Reader) (Message, error) {
	return m, nil
}

func (m Undefined4) sysCommon() {}

func (m Undefined4) Raw() []byte {
	return nil
}

type Undefined5 int

func (m Undefined5) String() string {
	return fmt.Sprintf("%T", m)
}

func (m Undefined5) meta() {}

func (m Undefined5) readFrom(rd io.Reader) (Message, error) {
	return m, nil
}

func (m Undefined5) sysCommon() {}

func (m Undefined5) Raw() []byte {
	return nil
}
