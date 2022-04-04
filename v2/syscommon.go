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

var syscommMessages = map[byte]Type{
	byteMIDITimingCodeMessage:/* SysCommonMsg.Set(MTCMsg), */ MTCMsg,
	byteSysSongPositionPointer:/* SysCommonMsg.Set(SPPMsg), */ SPPMsg,
	byteSysSongSelect:/* SysCommonMsg.Set(SongSelectMsg), */ SongSelectMsg,
	byteSysTuneRequest:/* SysCommonMsg.Set(TuneMsg), */ TuneMsg,
}

// Tune returns a tune message
func Tune() Message {
	//return NewMessage([]byte{byteSysTuneRequest})
	return []byte{byteSysTuneRequest}
}

// SPP returns a song position pointer message
func SPP(pointer uint16) Message {
	var b = make([]byte, 2)
	b[1] = byte(pointer & 0x7F)
	b[0] = byte((pointer >> 7) & 0x7F)
	//return NewMessage([]byte{byteSysSongPositionPointer, b[0], b[1]})
	return []byte{byteSysSongPositionPointer, b[0], b[1]}
}

// SongSelect returns a song select message
func SongSelect(song uint8) Message {
	// TODO check - it is a guess
	//return NewMessage([]byte{byteSysSongSelect, song})
	return []byte{byteSysSongSelect, song}
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

// MTC returns a timing code message (quarter frame)
func MTC(m uint8) Message {
	// TODO check - it is a guess
	// TODO provide a better abstraction for MTC
	//return NewMessage([]byte{byteMIDITimingCodeMessage, byte(m)})
	return []byte{byteMIDITimingCodeMessage, byte(m)}
}
