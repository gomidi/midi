package sysex

import "bytes"

type NonRealtime struct {
	Channel byte
	SubID1  byte
	SubID2  byte
}

func (r NonRealtime) SysEx() []byte {
	var bf bytes.Buffer

	bf.WriteByte(0xF0)
	bf.WriteByte(0x7E)
	bf.WriteByte(r.Channel)
	bf.WriteByte(r.SubID1)
	bf.WriteByte(r.SubID2)

	bf.WriteByte(0xF7)
	return bf.Bytes()

}

func GMSystem(channel byte, enable bool) []byte {
	var n NonRealtime
	n.Channel = channel
	n.SubID1 = 0x09
	if enable {
		n.SubID2 = 0x01
	} else {
		n.SubID2 = 0x00
	}

	return n.SysEx()
}

/*

GM System Enable/Disable

This Universal SysEx message enables or disables the General MIDI mode of a sound module. Some devices have
built-in GM modules or GM Patch Sets in addition to non-GM Patch Sets or non-GM modes of operation. When GM
is enabled, it replaces any non-GM Patch Set or non-GM mode with a GM mode/patch set. This allows a device to have
modes or Patch Sets that go beyond the limits of GM, and yet, still have the capability to be switched into a
GM-compliant mode when desirable.

0xF0  SysEx
0x7E  Non-Realtime
0x7F  The SysEx channel. Could be from 0x00 to 0x7F.
      Here we set it to "disregard channel".
0x09  Sub-ID -- GM System Enable/Disable
0xNN  Sub-ID2 -- NN=00 for disable, NN=01 for enable
0xF7  End of SysEx

It is best to respond as quickly as possible to this message, and to be ready to accept incoming note (and other)
messages soon after, as this message may be included at the start of a General MIDI file to ensure playback by a
GM module. Most modules are fully setup in GM mode by 100 milliseconds after receiving the GM System Enable message.

When GM Mode is first enabled, a device should assume that the "Grand Piano" patch (ie, the first GM patch) is the
currently selected patch upon all 16 MIDI channels. The device should also internally reset all controllers and
assume the power-up state described in the General MIDI Specification.

While GM mode is enabled, a device should also ignore Bank Select messages (since GM does not have more than one
bank of patches). Only when the GM Disable message is received (with Sub-ID2 = 0 to disable GM mode) will a device
then respond to Bank Select messages (and knock itself out of GM mode).

*/

func IdentityRequest(channel byte) []byte {
	var n NonRealtime
	n.Channel = channel
	n.SubID1 = 0x06
	n.SubID2 = 0x01
	return n.SysEx()
}

func IdentityReply(channel byte, manuID ManufacturerID, familycode [2]byte, modelnumber [2]byte, version [4]byte) []byte {
	var bf bytes.Buffer
	var n NonRealtime
	n.Channel = channel
	n.SubID1 = 0x06
	n.SubID2 = 0x02
	bt := n.SysEx()
	bf.Write(bt[:len(bt)-1]) // strip the 0xF7
	bf.WriteByte(byte(manuID))
	bf.WriteByte(familycode[0])
	bf.WriteByte(familycode[1])
	bf.WriteByte(modelnumber[0])
	bf.WriteByte(modelnumber[1])
	bf.WriteByte(version[0])
	bf.WriteByte(version[1])
	bf.WriteByte(version[2])
	bf.WriteByte(version[3])
	bf.WriteByte(0xF7)

	return bf.Bytes()
}

/*

Identity Request / Reply

Sometimes, a device may wish to know what other devices are connected to it. For example, a Patch Editor software running on a computer may wish to know what devices are connected to the computer's MIDI port, so that the software can configure itself to accept dumps from those devices.

The Identity Request Universal Sysex message can be sent by the Patch Editor software. When this message is received by some device connected to the computer, that device will respond by sending an Identity Reply Universal Sysex message back to the computer. The Patch Editor can then examine the information in the Identity Reply message to determine what make and model device is connected to the computer. Each device that understands the Identity Request will reply with its own Identity Reply message.

Here is the Identity Request message:

0xF0  SysEx
0x7E  Non-Realtime
0x7F  The SysEx channel. Could be from 0x00 to 0x7F.
      Here we set it to "disregard channel".
0x06  Sub-ID -- General Information
0x01  Sub-ID2 -- Identity Request
0xF7  End of SysEx

Here is the Identity Reply message:

0xF0  SysEx
0x7E  Non-Realtime
0x7F  The SysEx channel. Could be from 0x00 to 0x7F.
      Here we set it to "disregard channel".
0x06  Sub-ID -- General Information
0x02  Sub-ID2 -- Identity Reply
0xID  Manufacturer's ID
0xf1  The f1 and f2 bytes make up the family code. Each
0xf2  manufacturer assigns different family codes to his products.
0xp1  The p1 and p2 bytes make up the model number. Each
0xp2  manufacturer assigns different model numbers to his products.
0xv1  The v1, v2, v3 and v4 bytes make up the version number.
0xv2
0xv3
0xv4
0xF7  End of SysEx


*/
