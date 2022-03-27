package midi

/*

System Common Message   Status Byte      Number of Data Bytes
---------------------   -----------      --------------------
MIDI Timing Code            F1                   1
Song Position Pointer       F2                   2
Song Select                 F3                   1
Tune Request                F6                  None

*/

const (
	byteMIDITimingCodeMessage  = byte(0xF1)
	byteSysSongPositionPointer = byte(0xF2)
	byteSysSongSelect          = byte(0xF3)
	byteSysTuneRequest         = byte(0xF6)
)

var syscommMessages = map[byte]MsgType{
	byteMIDITimingCodeMessage:/* SysCommonMsg.Set(MTCMsg), */ MTC,
	byteSysSongPositionPointer:/* SysCommonMsg.Set(SPPMsg), */ SPP,
	byteSysSongSelect:/* SysCommonMsg.Set(SongSelectMsg), */ SongSelect,
	byteSysTuneRequest:/* SysCommonMsg.Set(TuneMsg), */ Tune,
}

// Tune returns a MIDI tune message
func NewTune() Msg {
	return NewMsg([]byte{byteSysTuneRequest})
}

// SPP returns a MIDI song position pointer message
func NewSPP(pointer uint16) Msg {
	var b = make([]byte, 2)
	b[1] = byte(pointer & 0x7F)
	b[0] = byte((pointer >> 7) & 0x7F)
	return NewMsg([]byte{byteSysSongPositionPointer, b[0], b[1]})
}

// SongSelect returns a MIDI song select message
func NewSongSelect(song uint8) Msg {
	// TODO check - it is a guess
	return NewMsg([]byte{byteSysSongSelect, song})
}

/*
MTC Quarter Frame

These are the MTC (i.e. SMPTE based) equivalent of the F8 Timing Clock messages,
though offer much higher resolution, as they are sent at a rate of 96 to 120 times
a second (depending on the SMPTE frame rate). Each Quarter Frame message provides
partial timecode information, 8 sequential messages being required to fully
describe a timecode instant. The reconstituted timecode refers to when the first
partial was received. The most significant nibble of the data byte indicates the
partial (aka Message Type).

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

// MTC represents a MIDI timing code message (quarter frame)
func NewMTC(m uint8) Msg {
	// TODO check - it is a guess
	// TODO provide a better abstraction for MTC
	return NewMsg([]byte{byteMIDITimingCodeMessage, byte(m)})
}
