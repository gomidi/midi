package syscommon

import "io"

/*

System Common Message   Status Byte      Number of Data Bytes
---------------------   -----------      --------------------
MIDI Timing Code            F1                   1
Song Position Pointer       F2                   2
Song Select                 F3                   1
Tune Request                F6                  None

*/

// Message is a System Common Message
type Message interface {
	String() string
	Raw() []byte
	readFrom(io.Reader) (Message, error)
	sysCommon()
}

var (
	_ Message = SongPositionPointer(0)
	_ Message = SongSelect(0)
	_ Message = TuneRequest
	//	_ Message = Undefined4(0)
	//	_ Message = Undefined5(0)
	_ Message = MIDITimingCode(0)
)
