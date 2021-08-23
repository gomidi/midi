package sysex

import "bytes"

type Realtime struct {
	Channel byte
	SubID1  byte
	SubID2  byte
}

func (r Realtime) SysEx() []byte {
	var bf bytes.Buffer

	bf.WriteByte(0xF0)
	bf.WriteByte(0x7F)
	bf.WriteByte(r.Channel)
	bf.WriteByte(r.SubID1)
	bf.WriteByte(r.SubID2)

	bf.WriteByte(0xF7)
	return bf.Bytes()
}

const EveryChannel = 0x7F

func MasterVolume(channel byte, vol uint16) []byte {
	var r Realtime
	r.Channel = channel
	r.SubID1 = 0x04
	r.SubID2 = 0x01

	/*
		TODO: parse the bits 0 to 6 and 7 to 13 of a 14-bit volume
	*/

	return r.SysEx()
}

/*

Master Volume

This Universal SysEx message adjusts a device's master volume. Remember that in a multitimbral device,
the Volume controller messages are used to control the volumes of the individual Parts. So, we need some
message to control Master Volume. Here it is.

0xF0  SysEx
0x7F  Realtime
0x7F  The SysEx channel. Could be from 0x00 to 0x7F.
      Here we set it to "disregard channel".
0x04  Sub-ID -- Device Control
0x01  Sub-ID2 -- Master Volume
0xLL  Bits 0 to 6 of a 14-bit volume
0xMM  Bits 7 to 13 of a 14-bit volume
0xF7  End of SysEx


*/
