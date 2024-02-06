package mmc

import (
	"bytes"
	"fmt"
)

// see https://en.wikipedia.org/wiki/MIDI_Machine_Control

/*
MIDI Universal Real Time SysEx Commands

All numbers are in hexadecimal notation. SysEx message format:

F0, 7F, nn, sub-ID, data, F7

nn = channel number, 00 to 7F; 7F = global
sub-IDs:
01 = Long Form MTC
02 = MIDI Show Control
03 = Notation Information
04 = Device Control
05 = Real Time MTC Cueing
06 = MIDI Machine Control Command
07 = MIDI Machine Control Response
08 = Single Note Retune
*/

type Command byte

/*
01 Stop
02 Play
03 Deferred Play (play after no longer busy)
04 Fast Forward
05 Rewind
06 Record Strobe (AKA [[Punch in/out|Punch In]])
07 Record Exit (AKA [[Punch out (music)|Punch out]])
08 Record Pause
09 Pause (pause playback)
0A Eject (disengage media container from MMC device)
0B Chase
0D MMC Reset (to default/startup state)
40 Write (AKA Record Ready, AKA Arm Tracks)

	parameters: <length1> 4F <length2> <track-bitmap-bytes>

44 Goto (AKA Locate)

	parameters: <length>=06 01 <hours> <minutes> <seconds> <frames> <subframes>

47 Shuttle

	parameters: <length>=03 <sh> <sm> <sl> (MIDI Standard Speed codes)
*/
const (
	StopCmd         Command = 0x01
	PlayCmd         Command = 0x02
	DeferredPlayCmd Command = 0x03
	FastForwardCmd  Command = 0x04
	RewindCmd       Command = 0x05
	RecordStrobeCmd Command = 0x06
	PunchInCmd      Command = 0x06
	RecordExitCmd   Command = 0x07
	PunchOutCmd     Command = 0x07
	RecordPauseCmd  Command = 0x08
	PauseCmd        Command = 0x09
	EjectCmd        Command = 0x0A
	ChaseCmd        Command = 0x0B
	ResetCmd        Command = 0x0D
	WriteCmd        Command = 0x40
	RecordReadyCmd  Command = 0x40
	ArmTrackCmd     Command = 0x40
	GoToCmd         Command = 0x44
	LocateCmd       Command = 0x44
	ShuttleCmd      Command = 0x47
)

func (c Command) String() string {
	switch byte(c) {
	case 0x01:
		return "StopCmd"
	case 0x02:
		return "PlayCmd"
	case 0x03:
		return "DeferredPlayCmd"
	case 0x04:
		return "FastForward"
	case 0x05:
		return "RewindCmd"
	case 0x06:
		return "RecordStrobeCmd/PunchInCmd"
	case 0x07:
		return "RecordExitCmd/PunchOutCmd"
	case 0x08:
		return "RecordPauseCmd"
	case 0x09:
		return "PauseCmd"
	case 0x0A:
		return "EjectCmd"
	case 0x0B:
		return "ChaseCmd"
	case 0x0D:
		return "ResetCmd"
	case 0x40:
		return "WriteCmd/RecordReadyCmd/ArmTrackCmd"
	case 0x44:
		return "GotoCmd/LocateCmd"
	case 0x47:
		return "ShuttleCmd"
	default:
		return "unknownCmd"
	}
}

type Message struct {
	DeviceID byte
	Command
	IsResponse bool
	Data       []byte
}

func (g *Message) Parse(bt []byte) error {
	if len(bt) < 5 {
		return fmt.Errorf("wrong length: %v (must be >= 5)", len(bt))
	}

	if bt[0] != 0xF0 {
		return fmt.Errorf("wrong byte 0")
	}

	if bt[1] != 0x7F {
		return fmt.Errorf("wrong byte 1")
	}

	if bt[len(bt)-1] != 0xF7 {
		return fmt.Errorf("wrong last byte")
	}

	g.DeviceID = bt[2]

	switch bt[3] {
	case 0x06:
		g.IsResponse = false
		if len(bt) < 7 {
			return fmt.Errorf("wrong length for command: %v (must be >= 7)", len(bt))
		}
		g.Command = Command(bt[4])
		if bt[4] >= 0x40 {
			if len(bt) < 8 {
				return fmt.Errorf("wrong length for %s command: %v (must be >= 8)", g.Command.String(), len(bt))
			}
			g.Data = bt[5 : len(bt)-2]
		}
	case 0x07:
		g.IsResponse = true
		if len(bt) > 5 {
			g.Data = bt[4 : len(bt)-2]
		} else {
			g.Data = nil
		}
	}

	return nil
}

func (m Message) String() string {
	return fmt.Sprintf("MMC device: %v command: %s", m.DeviceID, m.Command.String())
}

func (m Message) SysEx() []byte {
	var bf bytes.Buffer
	bf.WriteByte(0xF0)
	bf.WriteByte(0x7F)
	devID := m.DeviceID
	if devID == 0 || devID > 127 {
		devID = 127
	}
	bf.WriteByte(devID)
	bf.WriteByte(0x06)
	bf.WriteByte(byte(m.Command))
	bf.WriteByte(0xF7)

	return bf.Bytes()
}

type ArmTrack struct {
}

type GoTo struct {
	DeviceID byte
	Hour     byte
	Minute   byte
	Second   byte
	Frame    byte
	SubFrame byte
}

func (g GoTo) SysEx() []byte {
	return []byte{0xF0, 0x7F, g.DeviceID, 0x06, 0x44, 0x06, 0x01, g.Hour, g.Minute, g.Second, g.Frame, g.SubFrame, 0xF7}
}

func (g *GoTo) Parse(bt []byte) error {
	if len(bt) != 13 {
		return fmt.Errorf("wrong length: %v (must be 13)", len(bt))
	}

	if bt[0] != 0xF0 {
		return fmt.Errorf("wrong byte 0")
	}

	if bt[1] != 0x7F {
		return fmt.Errorf("wrong byte 1")
	}

	g.DeviceID = bt[2]

	if bt[3] != 0x06 {
		return fmt.Errorf("wrong byte 3")
	}

	if bt[4] != 0x44 {
		return fmt.Errorf("wrong byte 4")
	}

	if bt[5] != 0x06 {
		return fmt.Errorf("wrong byte 5")
	}

	if bt[6] != 0x01 {
		return fmt.Errorf("wrong byte 6")
	}

	g.Hour = bt[7]
	g.Minute = bt[8]
	g.Second = bt[9]
	g.Frame = bt[10]
	g.SubFrame = bt[11]

	if bt[12] != 0xF7 {
		return fmt.Errorf("wrong byte 12")
	}

	return nil
}

type Identity struct {
	Channel byte
}

func (i Identity) SysEx() []byte {
	return []byte{0xF0, 0x7E, i.Channel, 0x06, 0x01, 0xF7}
}

func (i Identity) Parse(bt []byte) error {
	//return []byte{0xF0, 0x7E, i.Channel, 0x06, 0x01, 0xF7}
	if len(bt) != 6 {
		return fmt.Errorf("wrong length: %v (must be 6)", len(bt))
	}

	if bt[0] != 0xF0 {
		return fmt.Errorf("wrong byte 0")
	}

	if bt[1] != 0x7E {
		return fmt.Errorf("wrong byte 1")
	}

	i.Channel = bt[2]

	if bt[3] != 0x06 {
		return fmt.Errorf("wrong byte 3")
	}

	if bt[4] != 0x01 {
		return fmt.Errorf("wrong byte 4")
	}

	if bt[5] != 0xF7 {
		return fmt.Errorf("wrong byte 5")
	}

	return nil
}
