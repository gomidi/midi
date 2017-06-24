package smf

import "github.com/gomidi/midi"

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
	Track() uint16
}

// Header represents the header of a SMF file
type Header interface {
	// Format returns the SMF format (0 = SingleTrack, 1 = MultiTrack, 2 = SequentialTracks)
	Format() Format

	// TimeFormat returns the time format (QuarterNoteTicks or TimeCode) and the value of that format
	// If TimeFormat is QuarterNoteTicks, the value is the ticks per quarter note.
	// If TimeFormat is TimeCode, the value is a raw value that must be unpacked with the help of
	// UnpackTimeCode.
	TimeFormat() (format TimeFormat, value uint16)

	// NumTracks returns the number of tracks as defined inside the SMF header. It should be the same
	// as the real number of tracks in the file, although there is no guaranty.
	NumTracks() uint16
}

// UnpackTimeCode unpacks the raw value returned from Header.TimeFormat if the format is TimeCode
// It returns SMPTE frames per second (29 corresponds to 30 drop frame) and the subframes.
func UnpackTimeCode(raw uint16) (fps, subframes uint8) {
	// bit shifting first byte to second inverting sign
	fps = uint8(int8(byte(raw>>8)) * (-1))

	// taking the second byte
	subframes = byte(raw & uint16(255))
	return
}

const (
	// SingleTrack represents the singletrack SMF format (0)
	SingleTrack = format(0)

	// MultiTrack represents the multitrack SMF format (1)
	MultiTrack = format(1)

	// SequentialTracks represents the sequential track SMF format (2)
	SequentialTracks = format(2)

	// QuarterNoteTicks represents the "ticks per quarter note" (metric) time format
	QuarterNoteTicks = timeformat("QuarterNoteTicks")

	// TimeCode represents the SMTPE/Timecode time formats
	TimeCode = timeformat("TimeCode")
)

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
	case SingleTrack:
		return "SMF0 (SingleTrack)"
	case MultiTrack:
		return "SMF1 (MultiTrack)"
	case SequentialTracks:
		return "SMF2 (SequentialTracks)"
	}
	panic("unreachable")
}

// timeformat is an implementation of TimeFormat
type timeformat string

func (t timeformat) String() string { return string(t) }
func (t timeformat) timeformat()    {}
