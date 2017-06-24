package sysex

import (
	"bytes"
	// "bytes"
	// "encoding/binary"
	"fmt"
	"github.com/gomidi/midi/internal/lib"
	"github.com/gomidi/midi/messages/realtime"
	"io"
)

/* http://www.somascape.org/midi/tech/spec.html

Universal System Exclusive ID Numbers

This includes the Non Real Time (7E) and Real Time (7F) ID extensions mentioned above in note 5.

The generalised format for both is as follows :
F0 <ID number> <Device ID> <Sub ID#1> <Sub ID#2> . . . . F7

The Device ID (also referred to as the channel number) is generally used to specify a discrete physical device, however complex devices (e.g. computers with a number of different MIDI expansion cards) may have more than one Device ID. A value of 7F is used to specify all devices.

From one to sixteen virtual devices may be accessed at each Device ID by use of the normal MIDI channel numbers, within the capabilities of the device.

The Non-Commercial Universal System Exclusive ID (7D) is not detailed here.
Sub ID#1	Sub ID#2	Description
Non Real Time Universal Sys Ex   (Sys Ex ID number = 7E)
00	—	Unused
01	(not used)	Sample Dump Header - (Detail)
02	(not used)	Sample Dump Data Packet - (Detail)
03	(not used)	Sample Dump Request - (Detail)
04	nn	MTC Cueing - (Detail)
00	Special
01	Punch In Points
02	Punch Out Points
03	Delete Punch In Point
04	Delete Punch Out Point
05	Event Start Point
06	Event Stop Point
07	Event Start Points with additional info
08	Event Stop Points with additional info
09	Delete Event Start Point
0A	Delete Event Stop Point
0B	Cue Points
0C	Cue Points with additional info
0D	Delete Cue Point
0E	Event name in additional info
05	nn	Sample Dump Extensions
01	Loop Point Transmission - (Detail)
02	Loop Point Request - (Detail)
03	Sample Name Transmission - (Detail)
04	Sample Name Request - (Detail)
05	Extended Dump Header - (Detail)
06	Extended Loop Point Transmission - (Detail)
07	Extended Loop Point Request - (Detail)
06	nn	General System Information
01	Device Identity Request - (Detail)
02	Device Identity Reply - (Detail)
07	nn	File Dump
01	Header - (Detail)
02	Data Packet - (Detail)
03	Request - (Detail)
08	nn	MIDI Tuning Standard
00	Bulk Tuning Dump Request - (Detail)
01	Bulk Tuning Dump Reply - (Detail)
09	nn	General MIDI (GM) System
01	Enable - (Detail)
02	Disable - (Detail)
0A	nn	Down-Loadable Sounds (DLS) System
01	Enable - (Detail)
02	Disable - (Detail)
7B	(not used)	End of File - (Detail)
7C	(not used)	Wait - (Detail)
7D	(not used)	Cancel - (Detail)
7E	(not used)	NAK - (Detail)
7F	(not used)	ACK - (Detail)
Real Time Universal Sys Ex   (Sys Ex ID number = 7F)
00	—	Unused
01	nn	MIDI Time Code (MTC)
01	Full Message - (Detail)
02	User Bits - (Detail)
03	nn	Notation Information
01	Bar Marker - (Detail)
02	Time Signature (immediate) - (Detail)
42	Time Signature (delayed) - (Detail)
04	nn	Device Control
01	Master Volume - (Detail)
02	Master Balance - (Detail)
03	Master Fine Tuning - (Detail)
04	Master Coarse Tuning - (Detail)
05	nn	MTC Cueing - (Detail)
00	Special
01	Punch In Points
02	Punch Out Points
05	Event Start Point
06	Event Stop Point
07	Event Start Points with additional info
08	Event Stop Points with additional info
0B	Cue Points
0C	Cue Points with additional info
0E	Event name in additional info
06	nn	MIDI Machine Control (MMC) Commands - (Detail)
01	Stop
02	Play
03	Deferred Play
04	Fast Forward
05	Rewind
06	Record Strobe (Punch In)
07	Record Exit (Punch Out)
08	Record Pause
09	Pause
0A	Eject
0B	Chase
0C	Command Error Reset
0D	MMC Reset
44	Locate / Go To - (Detail)
47	Shuttle - (Detail)
08	nn	MIDI Tuning Standard
02	Single Note Tuning Change - (Detail)


General System Information
Device Identity Request

This message is sent to request the identity of the receiving device.
F0 7E id 06 01 F7
F0 7E	Universal Non Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
06	Sub ID#1 = General System Information
01	Sub ID#2 = Device Identity Request
F7	EOX
Device Identity Reply
F0 7E id 06 02 mm ff ff dd dd ss ss ss ss F7
F0 7E	Universal Non Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
06	Sub ID#1 = General System Information
02	Sub ID#2 = Device Identity message
mm	Manufacturers System Exclusive ID code.
If mm = 00, then the message is extended by 2 bytes to accomodate the additional manufacturers ID code.
ff ff	Device family code (14 bits, LSB first)
dd dd	Device family member code (14 bits, LSB first)
ss ss ss ss	Software revision level (the format is device specific)
F7	EOX

File Dump

This standard is similar to the Sample Dump Standard (SDS) but applies to file data rather than to sample data.

It comprises three similarly named messages (Dump Request, Dump Header and Data Packet), which as with the SDS are used in combination with the five Handshaking messages (EOF, ACK, NAK, Wait and Cancel) in accordance with the Dump Procedure.
File Dump Header
F0 7E id 07 01 ss <type> <length> <name> F7
F0 7E	Universal Non Real Time Sys Ex header
id	Device ID of requester / receiver
07	Sub ID#1 = File Dump
01	Sub ID#2 = Header
ss	Device ID of sender
<type>	Type of file (4 7-bit ASCII bytes) :
Type	DOS extension	Description
'MIDI'	.MID	MIDI File
'MIEX'	.MEX	MIDIEX File
'ESEQ'	.ESQ	ESEQ File
'TEXT'	.TXT	7-bit ASCII Text File
'BIN '	.BIN	Binary File (e.g. MS-DOS)
'MAC '	.MAC	Macintosh File (MacBinary header)
<length>	File length (4 7-bit bytes, LSB first)
<name>	Filename (7-bit ASCII bytes; as many as needed)
F7	EOX
File Dump Data Packet
F0 7E id 07 02 pp bb <data> <checksum> F7
F0 7E	Universal Non Real Time Sys Ex header
id	Device ID of receiver
07	Sub ID#1 = File Dump
02	Sub ID#2 = Data Packet
pp	Packet count
bb	Packet size (number of encoded data bytes - 1)
<data>	Data (encoded as below)
<checksum>	XOR of all bytes following the SOX up to the checksum byte
F7	EOX

Encoding of data :

The 8-bit file data needs to be converted to 7-bit form, with the result that every 7 bytes of file data translates to 8 bytes in the MIDI stream.

For each group of 7 bytes (of file data) the top bit from each is used to construct an eigth byte, which is sent first. So :
AAAAaaaa BBBBbbbb CCCCcccc DDDDdddd EEEEeeee FFFFffff GGGGgggg

becomes :
0ABCDEFG 0AAAaaaa 0BBBbbbb 0CCCcccc 0DDDdddd 0EEEeeee 0FFFffff 0GGGgggg

The final group may have less than 7 bytes, and is coded as follows (e.g. with 3 bytes in the final group) :
0ABC0000 0AAAaaaa 0BBBbbbb 0CCCcccc
File Dump Request
F0 7E id 07 03 ss <type> <name> F7
F0 7E	Universal Non Real Time Sys Ex header
id	Device ID of request destination (file sender)
07	Sub ID#1 = File Dump
03	Sub ID#2 = Request
ss	Device ID of requester (file receiver)
<type>	Type of file (4 7-bit ASCII bytes) :
Type	DOS extension	Description
'MIDI'	.MID	MIDI File
'MIEX'	.MEX	MIDIEX File
'ESEQ'	.ESQ	ESEQ File
'TEXT'	.TXT	7-bit ASCII Text File
'BIN '	.BIN	Binary File (e.g. MS-DOS)
'MAC '	.MAC	Macintosh File (MacBinary header)
<name>	Filename (7-bit ASCII bytes; as many as needed)
F7	EOX


MIDI Tuning Standard

This standard requires that each of the 128 MIDI note numbers can be tunable to any frequency within the instrument's range. Additionally it allows provision for up to 128 tuning programs.

An instrument which supports MIDI Tuning may have less than the full complement of tuning programs.

The specification provides a frequency resolution somewhat finer than one-hundredth of a cent, which should be fine enough for most needs. Instruments supporting MIDI Tuning need not necessarily provide this fine a resolution – the specification merely permits the transfer of tuning data at any resolution up to this limit.

See also the Real Time MIDI Tuning messages and Registered Parameter Numbers 3 (Select Tuning Program) and 4 (Select Tuning Bank).
Bulk Tuning Dump Request
F0 7E id 08 00 tt F7
F0 7E	Universal Non Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
08	Sub ID#1 = MIDI Tuning Standard
00	Sub ID#2 = Bulk Tuning Dump Request
tt	Tuning Program number (0-127)
F7	EOX
Bulk Tuning Dump Reply

This message comprises frequency data in the 3-byte format outlined below, for all 128 MIDI note numbers, in sequence from note 0 (sent first) to note 127 (sent last).
F0 7E id 08 01 tt <name> [ xx yy zz ] ... <checksum> F7
F0 7E	Universal Non Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
08	Sub ID#1 = MIDI Tuning Standard
01	Sub ID#2 = Bulk Tuning Dump Reply
tt	Tuning Program number (0-127)
<name>	Name (16 7-bit ASCII bytes)
[ xx yy zz ]	Frequency data for one note (repeated 128 times).   See below.
<checksum>	Checksum
F7	EOX

Frequency data :

The 3-byte frequency data has the following format :
0xxxxxxx 0abcdefg 0hijklmn

Where :
xxxxxxx = semitone (MIDI note number);
abcdefghijklmn = fraction of a semitone, in 0.0061 cent units.

The frequency range starts at MIDI note 0, C = 8.1758 Hz, and extends above MIDI note 127, G = 12543.875 Hz. The first byte of the frequency data word specifies the nearest equal-tempered semitone below the required frequency. The remaining two bytes specify the fraction of 100 cents above the semitone at which the required frequency lies.

The frequency data of 7F 7F 7F has special significance, and indicates a no change situation. I.e. when an instrument receives these 3 bytes as frequency data, it should make no change to its stored frequency data for that MIDI key number.


General MIDI (GM)
Enable
F0 7E id 09 01 F7
F0 7E	Universal Non Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
09	Sub ID#1 = General MIDI (GM)
01	Sub ID#2 = Enable
F7	EOX
Disable
F0 7E id 09 02 F7
F0 7E	Universal Non Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
09	Sub ID#1 = General MIDI (GM)
02	Sub ID#2 = Disable
F7	EOX


Device Control

These four messages each have a channel-based equivalent.
Device control	Channel-based
Master Volume	Channel Volume (CC 7)
Master Balance	Channel Balance (CC 8)
Master Fine Tuning	Channel Fine Tuning (RPN 1)
Master Coarse Tuning	Channel Coarse Tuning (RPN 2)
Master Volume
F0 7F id 04 01 vv vv F7
F0 7F	Universal Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
04	Sub ID#1 = Device Control message
01	Sub ID#2 = Master Volume
vv vv	Volume (LSB first) : 00 00 = volume off
F7	EOX
Master Balance
F0 7F id 04 02 bb bb F7
F0 7F	Universal Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
04	Sub ID#1 = Device Control message
02	Sub ID#2 = Master Balance
bb bb	Balance (LSB first) : 00 00 = hard left, 7F 7F = hard right
F7	EOX
Master Fine Tuning
F0 7F id 04 03 tt tt F7
F0 7F	Universal Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
04	Sub ID#1 = Device Control message
03	Sub ID#2 = Master Fine Tuning
tt tt	Fine Tuning (LSB first). Displacement in cents from A440
LSB	MSB	Displacement (cents)
00	00	100 / 8192 * (-8192)
00	40	100 / 8192 * 0
7F	7F	100 / 8192 * (+8191)
F7	EOX

The total fine tuning displacement in cents from A440 for each MIDI channel is the summation of the displacement of this Master Fine Tuning and the displacement of the Channel Fine Tuning RPN.
Master Coarse Tuning
F0 7F id 04 04 00 tt F7
F0 7F	Universal Real Time Sys Ex header
id	ID of target device   (default = 7F = All devices)
04	Sub ID#1 = Device Control message
04	Sub ID#2 = Master Coarse Tuning
00 tt	Coarse Tuning (LSB first, though the LSB is always zero)
LSB	MSB	Displacement (cents)
00	00	100 * (-64)
00	40	100 * 0
00	7F	100 * (+63)
F7	EOX

The total coarse tuning displacement in cents from A440 for each MIDI channel is the summation of the displacement of this Master Coarse Tuning and the displacement of the Channel Coarse Tuning RPN.

MIDI Machine Control (MMC) Commands
'Single-byte' MMC Commands

These messages require no additional data beyond the command code (Sub ID#2) itself.
F0 7F id 06 cc F7
F0 7F	Universal Real Time Sys Ex header
id	ID of target device
06	Sub ID#1 = MMC Command
cc	Sub ID#2 = Command (see the following table)
F7	EOX
cc	Description
01	Stop 	Instructs the receiving device to immediately cease playback.
02	Play 	Instructs the receiving device to immediately begin playback.
03	Deferred Play 	Instructs the receiving device to begin playback. If the device is busy (e.g. it has been instructed to locate the playhead to a specific position, but hasn't yet got there) then playback will begin when the device is ready.
04	Fast Forward 	Instructs the receiving device to immediately enter fast forward mode.
05	Rewind 	Instructs the receiving device to immediately enter rewind mode.
06	Record Strobe 	Instructs the receiving device to immediately start recording (Punch In).
07	Record Exit 	Instructs the receiving device to immediately stop recording (Punch Out).
08	Record Pause 	Instructs the receiving device to immediately enter record ready mode.
09	Pause 	Instructs the receiving device to immediately enter paused mode.
0A	Eject
0B	Chase
0C	Command Error Reset
0D	MMC Reset 	Resets the receiving device to its default / start-up state.
Locate / Go To

This message moves the receiving device's playhead to the specified SMPTE location.
F0 7F id 06 44 06 01 hr mn sc fr sf F7
F0 7F	Universal Real Time Sys Ex header
id	ID of target device
06	Sub ID#1 = MMC Command
44	Sub ID#2 = Locate / Go To
06	number of data bytes that follow
01
hr	Hours and Type : 0yyzzzzz
yy = Type: 00 = 24 fps, 01 = 25 fps, 10 = 30 fps (drop frame), 11 = 30 fps (non-drop frame)
zzzzz = Hours (0-23)
mn	Minutes (0-59)
sc	Seconds (0-59)
fr	SMPTE frame number (0-29)
sf	SMPTE sub-frame number (0-99)
F7	EOX
Shuttle

This message provides bi-directional shuttling of the receiving device's playhead position.
F0 7F id 06 47 03 sh sm sl F7
F0 7F	Universal Real Time Sys Ex header
id	ID of target device
06	Sub ID#1 = MMC Command
47	Sub ID#2 = Shuttle
03	number of data bytes that follow
sh sm sl	Shuttle direction and speed. Bit 6 of sh gives direction (0 = forward, 1 = backward)
F7	EOX
B

MIDI Tuning Standard

See also the Non Real Time MIDI Tuning messages.
Single Note Tuning Change

This message enables retuning of individual MIDI note numbers to new frequencies in real time as a performance control. It also allows multiple changes to be made using a single message.

Note that if the note being retuned is currently sounding (within a tone generator), the note should be immediately retuned as it continues to sound with no glitching, re-triggering or other audible artifacts.

See also Registered Parameter Numbers 3 (Select Tuning Program) and 4 (Select Tuning Bank), which allow the selection of predefined tunings.
F0 7F id 08 02 tt nn [ kk xx yy zz ] ... F7
F0 7F	Universal Real Time Sys Ex header
id	ID of target device
08	Sub ID#1 = MIDI Tuning Standard
02	Sub ID#2 = Single Note Tuning Change
tt	Tuning Program number (0-127)
nn	Number of changes; 1 change = 1 set of [ kk xx yy zz ]
[ kk	MIDI Key number
xx yy zz ]	Frequency data for key 'kk' (repeated 'nn' times).   See below.
F7	EOX

Frequency data :

The 3-byte frequency data has the following format :
0xxxxxxx 0abcdefg 0hijklmn

Where :
xxxxxxx = semitone (MIDI note number);
abcdefghijklmn = fraction of a semitone, in 0.0061 cent units.

The frequency range starts at MIDI note 0, C = 8.1758 Hz, and extends above MIDI note 127, G = 12543.875 Hz. The first byte of the frequency data word specifies the nearest equal-tempered semitone below the required frequency. The remaining two bytes specify the fraction of 100 cents above the semitone at which the required frequency lies.

*/

// if canary >= 0xF0 && canary <= 0xF7 {
const (
	byteSysExStart = byte(0xF0)
	byteSysExEnd   = byte(0xF7)
)

type Message interface {
	String() string
	Raw() []byte
	// readFrom(io.Reader) (Message, error)
}

/*
   Furthermore, although the 0xF7 is supposed to mark the end of a SysEx message, in fact, any status
   (except for Realtime Category messages) will cause a SysEx message to be
   considered "done" (ie, actually "aborted" is a better description since such a scenario
   indicates an abnormal MIDI condition). For example, if a 0x90 happened to be sent sometime
   after a 0xF0 (but before the 0xF7), then the SysEx message would be considered
   aborted at that point. It should be noted that, like all System Common messages,
   SysEx cancels any current running status. In other words, the next Voice Category
   message (after the SysEx message) must begin with a Status.
*/

// ReadLive reads a sysex "over the wire", "in live mode", "as a stream" - you name it -
// opposed to reading a sysex from a SMF standard midi file
// the sysex has already been started (0xF0 has been read)
// we need a realtime.Reader here, since realtime messages must be handled (or ignored from the viewpoit of sysex)
func ReadLive(rd realtime.Reader) (sys SysEx, status byte, err error) {
	var b byte
	var bf bytes.Buffer
	// read byte by byte
	for {
		b, err = lib.ReadByte(rd)
		if err != nil {
			break
		}

		// the normal way to terminate
		if b == byte(0xF7) {
			sys = SysEx(bf.Bytes())
			return
		}

		// not so elegant way to terminate by sending a new status
		if lib.IsStatusByte(b) {
			sys = SysEx(bf.Bytes())
			status = b
			return
		}

		bf.WriteByte(b)
	}

	// any error, especially io.EOF is considered a failure.
	// however return the sysex that had been received so far back to the user
	// and leave him to decide what to do.
	sys = SysEx(bf.Bytes())
	return
}

/*
	F0 <length> <bytes to be transmitted after F0>

	The length is stored as a variable-length quantity. It specifies the number of bytes which follow it, not
	including the F0 or the length itself. For instance, the transmitted message F0 43 12 00 07 F7 would be stored
	in a MIDI File as F0 05 43 12 00 07 F7. It is required to include the F7 at the end so that the reader of the
	MIDI File knows that it has read the entire message.
*/

/*
	   Another form of sysex event is provided which does not imply that an F0 should be transmitted. This may be
	   used as an "escape" to provide for the transmission of things which would not otherwise be legal, including
	   system realtime messages, song pointer or select, MIDI Time Code, etc. This uses the F7 code:

	   F7 <length> <all bytes to be transmitted>

	   Unfortunately, some synthesiser manufacturers specify that their system exclusive messages are to be
	   transmitted as little packets. Each packet is only part of an entire syntactical system exclusive message, but
	   the times they are transmitted are important. Examples of this are the bytes sent in a CZ patch dump, or the
	   FB-01's "system exclusive mode" in which microtonal data can be transmitted. The F0 and F7 sysex events
	   may be used together to break up syntactically complete system exclusive messages into timed packets.
	   An F0 sysex event is used for the first packet in a series -- it is a message in which the F0 should be
	   transmitted. An F7 sysex event is used for the remainder of the packets, which do not begin with F0. (Of
	   course, the F7 is not considered part of the system exclusive message).
	   A syntactic system exclusive message must always end with an F7, even if the real-life device didn't send one,
	   so that you know when you've reached the end of an entire sysex message without looking ahead to the next
	   event in the MIDI File. If it's stored in one complete F0 sysex event, the last byte must be an F7. There also
	   must not be any transmittable MIDI events in between the packets of a multi-packet system exclusive
	   message. This principle is illustrated in the paragraph below.

			Here is a MIDI File of a multi-packet system exclusive message: suppose the bytes F0 43 12 00 were to be
			sent, followed by a 200-tick delay, followed by the bytes 43 12 00 43 12 00, followed by a 100-tick delay,
			followed by the bytes 43 12 00 F7, this would be in the MIDI File:

			F0 03 43 12 00						|
			81 48											| 200-tick delta time
			F7 06 43 12 00 43 12 00   |
			64												| 100-tick delta time
			F7 04 43 12 00 F7         |

			When reading a MIDI File, and an F7 sysex event is encountered without a preceding F0 sysex event to start a
			multi-packet system exclusive message sequence, it should be presumed that the F7 event is being used as an
			"escape". In this case, it is not necessary that it end with an F7, unless it is desired that the F7 be transmitted.
*/

// even better readable: (from http://www.somascape.org/midi/tech/mfile.html#sysex)

/*
			SysEx events

There are a couple of ways of encoding System Exclusive messages. The normal method is to encode them as a single event, though it is also possible to split messages into separate packets (continuation events). A third form (an escape sequence) is used to wrap up arbitrary bytes that could not otherwise be included in a MIDI file.
Single (complete) SysEx messages

F0 length message

length is a variable length quantity (as used to represent delta-times) which specifies the number of bytes in the following message.
message is the remainder of the system exclusive message, minus the initial 0xF0 status byte.

Thus, it is just like a normal system exclusive message, though with the additional length parameter.

Note that although the terminal 0xF7 is redundant (strictly speaking, due to the use of a length parameter) it must be included.
Example

The system exclusive message :
F0 7E 00 09 01 F7

would be encoded (without the preceding delta-time) as :
F0 05 7E 00 09 01 F7

(In case you're wondering, this is a General MIDI Enable message.)
SysEx messages sent as packets - Continuation events

Some older MIDI devices, with slow onboard processors, cannot cope with receiving a large amount of data en-masse, and require large system exclusive messages to be broken into smaller packets, interspersed with a pause to allow the receiving device to process a packet and be ready for the next one.

This approach can of course be used with the method described above, i.e. with each packet being a self-contained system exclusive message (i.e. each starting with 0xF0 and ending with 0xF7).

Unfortunately, some manufacturers (notably Casio) have chosen to bend the standard, and rather than sending the packets as self-contained system exclusive messages, they act as though running status applied to system exclusive messages (which it doesn't - or at least it shouldn't).

What Casio do is this : the first packet has an initial 0xF0 byte but doesn't have a terminal 0xF7. The last packet doesn't have an initial 0xF0 but does have a terminal 0xF7. All intermediary packets have neither. No unrelated events should occur between these packets. The idea is that all the packets can be stitched together at the receiving device to create a single system exclusive message.

Putting this into a MIDI file, the first packet uses the 0xF0 status, whereas the second and subsequent packets use the 0xF7 status. This use of the 0xF7 status is referred to as a continuation event.
Example

A 3-packet message :
F0 43 12 00
43 12 00 43 12 00
43 12 00 F7

with a 200-tick delay between the first two, and a 100-tick delay between the final two, would be encoded (without the initial delta-time, before the first packet) :
F0 03 43 12 00 	first packet (the 4 bytes F0,43,12,00 are transmitted)
81 48 	200-tick delta-time
F7 06 43 12 00 43 12 00 	second packet (the 6 bytes 43,12,00,43,12,00 are transmitted)
64 	100-tick delta-time
F7 04 43 12 00 F7 	third packet (the 4 bytes 43,12,00,F7 are transmitted)

See the note below regarding distinguishing packets and escape sequences (which both use the 0xF7 status).
Escape sequences

F7 length bytes

length is a variable length quantity which specifies the number of bytes in bytes.

This has nothing to do with System Exclusive messages as such, though it does use the 0xF7 status. It provides a way of including bytes that could not otherwise be included within a MIDI file, e.g. System Common and System Real Time messages (Song Position Pointer, MTC, etc).

Note that Escape sequences do not have a terminal 0xF7 byte.
Example

The Song Select System Common message :
F3 01

would be encoded (without the preceding delta-time) as :
F7 02 F3 01

You are not restricted to single messages per escape sequence - any arbitrary collection of bytes may be included in a single sequence.

Note Parsing the 0xF7 status byte

When an event with an 0xF7 status byte is encountered whilst reading a MIDI file, its interpretation (SysEx packet or escape sequence) is determined as follows :

    When an event with 0xF0 status but lacking a terminal 0xF7 is encountered, then this is the first of a Casio-style multi-packet message, and a flag (boolean variable) should be set to indicate this.

    If an event with 0xF7 status is encountered whilst this flag is set, then this is a continuation event (a system exclusive packet, one of many).
    If this event has a terminal 0xF7, then it is the last packet and flag should be cleared.

    If an event with 0xF7 status is encountered whilst flag is clear, then this event is an escape sequence.

Naturally, the flag should be initialised clear prior to reading each track of a MIDI file.

*/

func ReadSMF(startcode byte, rd io.Reader) (sys SysEx, err error) {
	/*
		what this means to us is relatively simple:
		we read the data after the startcode based of the following length
		and return the sysex chunk with the start code.
		If it ends with F7 or not, is not our business (the device has to deal with it).
		Also, if there are multiple sysexes belonging to each other yada yada.
	*/

	switch startcode {
	case 0xF0, 0xF7:
		var data []byte
		data, err = lib.ReadVarLengthData(rd)

		if err != nil {
			return nil, err
		}

		sys = append(sys, startcode)
		sys = append(sys, data...)

	default:
		panic("sysex in SMF must start with F0 or F7")
	}

	return

}

var _ Message = SysEx([]byte{})

type SysEx []byte

func (m SysEx) Bytes() []byte {
	return []byte(m)
}

/*
// TODO: implement
func (m SysEx) readFrom(rd io.Reader) (Message, error) {
	return m, nil
}
*/

func (m SysEx) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

func (m SysEx) Len() int {
	return len(m)
}

func (m SysEx) Raw() []byte {
	var b = []byte{0xF0}
	b = append(b, []byte(m)...)
	b = append(b, 0xF7)
	return b
}
