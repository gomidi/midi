package cc

import (
	"github.com/gomidi/midi/midimessage/channel"
)

func newControlCaseSwitch(ch channel.Channel, controller uint8, on bool) channel.ControlChange {
	var value uint8 // start with off = 0
	if on {
		value = 127 // on = 127
	}
	return ch.ControlChange(controller, value)
}

/* http://www.somascape.org/midi/tech/spec.html
RPNs and NRPNs

Controller numbers 98 & 99 (Non-Registered Parameter Number, LSB & MSB), and 100 & 101 (Registered Parameter Number, LSB & MSB), in conjunction with Controller numbers 6 & 38 (Data Entry, MSB & LSB), 96 (Data Increment), and 97 (Data Decrement) extend the number of controllers available via MIDI.

Their use involves selecting the parameter number to be edited using Controllers 98 & 99 or 100 & 101, and then adjusting the value for that parameter using Controller number 6/38, 96, or 97. Controllers 6/38 would be used to set a specific value, whereas Controllers 96 and 97 would be used to nudge the current value up or down, respectively (with the Data2 byte specifying the step size, 1-127).

Note that two 7-bit values are used both to specify the RPN/NRPN itself, and its value. I.e. there are 16,384 possible RPNs (and 16,384 NRPNs) each of which can be set with 14-bit precision. The MSB and LSB data values need not represent a single 0-16,383 range, and can be used to represent different quantities, as in the case of RPN 00,00 (Pitch Bend Sensitivity) where the MSB represents 'semitones' and the LSB represents 'cents'.

To calculate the full number from its LSB and MSB : N = MSB * 128 + LSB

Once the required RPN/NRPN has been specified (using CC98 & CC99, or CC100 & CC101) the value needs to be specified. There are several options :

    High resolution values (i.e. ranging 0-16,383) require two messages :
        A CC6 Data Entry message, where the Data2 byte contains the MSB of the value.
        A CC38 Data Entry message, where the Data2 byte contains the LSB of the value.
    Low resolution values (i.e. 0-127) can be sent in one message. It could be either a CC6 or a CC38 Data Entry message.
    A CC96 Data Increment message can be sent. Its Data2 byte holds a value (0-127) to be added to the current parameter value.
    A CC97 Data Decrement message can be sent. Its Data2 byte holds a value (0-127) to be subtracted from the current parameter value.

None of the Non-Registered Parameter Numbers have been assigned specific functions. They may be used for different functions by different manufacturers, and are thus manufacturer-specific.

Registered Parameter Numbers are those which have been assigned some particular function by the MIDI Manufacturers Association (MMA) and the Japan MIDI Standards Committee (JMSC). The following RPNs are currently defined :
RPN	Description
MSB	LSB
00	00	Pitch Bend Sensitivity
The coarse adjustment (Controller 6) sets the number of semitones.
The fine adjustment (Controller 38) sets the number of cents.
00	01	Channel Fine Tuning
Use controllers 6 and 38 to set MSB and LSB :
00 00 = -100 cents; 40 00 = A440; 7F 7F = +100 cents. 	See also the Real Time Universal System Exclusive Master Tuning messages.
00	02	Channel Coarse Tuning
Coarse adjustment only, using controller 6 :
00 = -64 semitones; 40 = A440; 7F = +63 semitones.
00	03	Select Tuning Program 	See also the Real Time Universal System Exclusive MIDI Tuning Standard messages.
00	04	Select Tuning Bank
00	05	Modulation Depth Range
Used to scale the effective range of Controller 1 (Modulation Wheel).
7F	7F	Null Function
Used to cancel a RPN or NRPN. After it is received, any additional value updates received should no longer be applied to the previously selected RPN or NRPN.

It is recommended that the Null Function (RPN 7F,7F) should be sent immediately after a RPN or NRPN and its value are sent.

Example usage :

To set the Pitch Bend Sensitivity on channel 'n' to +/- 7 semitones (ie +/- a fifth) :
11 bytes 	Bn 64 00 65 00 06 07 64 7F 65 7F
Bn 64 00	RPN LSB = 00 	Select RPN : Pitch Bend Sensitivity
65 00	RPN MSB = 00 (running status in effect)
06 07	Data Entry MSB = 07 (running status in effect) 	Coarse adjustment (semitones)
64 7F	RPN LSB = 7F (running status in effect) 	Null Function (Cancel RPN)
65 7F	RPN MSB = 7F (running status in effect)

Or, to set Tuning Program 'tt' on channel 'n' :
11 bytes 	Bn 64 03 65 00 06 tt 64 7F 65 7F
Bn 64 03	RPN LSB = 03 	Select RPN : Select Tuning Program
65 00	RPN MSB = 00 (running status in effect)
06 tt	Data Entry MSB = tt (running status in effect) 	Coarse adjustment
64 7F	RPN LSB = 7F (running status in effect) 	Null Function (Cancel RPN)
65 7F	RPN MSB = 7F (running status in effect)

Or, to increment (by 1) the current Tuning Bank on channel 'n' :
11 bytes 	Bn 64 04 65 00 60 01 64 7F 65 7F
Bn 64 04	RPN LSB = 04 	Select RPN : Select Tuning Bank
65 00	RPN MSB = 00 (running status in effect)
60 01	Data Increment (running status in effect) 	Increment (by 1)
64 7F	RPN LSB = 7F (running status in effect) 	Null Function (Cancel RPN)
65 7F	RPN MSB = 7F (running status in effect)

When sending successive RPN (or NRPN) messages, the standard allows the omission of CC100 & CC101 (or CC98 & CC99) if their value is unchanged.

For example, to send a pair of Channel Coarse Tuning and Channel Fine Tuning RPN messages (with a low-res and high-res data value, respectively) :
17 bytes 	Bn 64 02 65 00 06 TM 64 01 06 tM 26 tL 64 7F 65 7F
Bn 64 02	RPN LSB = 02 	Select RPN : Channel Coarse Tuning
65 00	RPN MSB = 00 (running status in effect)
06 TM	Data Entry, MSB (running status in effect) 	Coarse adjustment (of 'coarse tuning')
64 01	RPN LSB = 01 (running status in effect) 	Select RPN : Channel Fine Tuning (RPN MSB = 00 is still in effect)
06 tM	Data Entry, MSB (running status in effect) 	Coarse adjustment (of 'fine tuning')
26 tL	Data Entry, LSB (running status in effect) 	Fine adjustment (of 'fine tuning')
64 7F	RPN LSB = 7F (running status in effect) 	Null Function (Cancel RPN)
65 7F	RPN MSB = 7F (running status in effect)


Controller Numbers

The controller number is the first data byte (0ccccccc) following a Controller Change status byte (Bn). The 128 available controller numbers are split into four groups :
0-63	High resolution continuous controllers (0-31 = MSB; 32-63 = LSB)
64-69	Switches
70-119	Low resolution continuous controllers
120-127	Channel Mode messages

Note that for switches, the second data byte (0vvvvvvv) is either 0 (Off) or 127 (On), and that for high and low resolution continuous controllers, the second data byte takes the range 0-127.

The high resolution continuous controllers are divided into MSB and LSB values, providing a maximum of 14-bit resolution. If only 7-bit resolution is needed for a specific controller, only the MSB is used - it is not necessary to send the LSB. If the full resolution is required, then the MSB should be sent first, followed by the LSB. If only the LSB has changed in value, the LSB may be sent without re-sending the MSB.

The controller numbers missing from the following list (3, 9, 14, 15, 20-31, 35, 41, 46, 47, 52-63, 85-87, 89, 90 and 102-119) are currently undefined.
High resolution continuous controllers (MSB)
0	Bank Select   (Detail)
1	Modulation Wheel
2	Breath Controller
4	Foot Controller
5	Portamento Time
6	Data Entry   (used with RPNs/NRPNs)
7	Channel Volume
8	Balance
10	Pan
11	Expression Controller
12	Effect Control 1
13	Effect Control 2
16	Gen Purpose Controller 1
17	Gen Purpose Controller 2
18	Gen Purpose Controller 3
19	Gen Purpose Controller 4
High resolution continuous controllers (LSB)
32	Bank Select
33	Modulation Wheel
34	Breath Controller
36	Foot Controller
37	Portamento Time
38	Data Entry
39	Channel Volume
40	Balance
42	Pan
43	Expression Controller
44	Effect Control 1
45	Effect Control 2
48	General Purpose Controller 1
49	General Purpose Controller 2
50	General Purpose Controller 3
51	General Purpose Controller 4
Switches
64	Sustain On/Off
65	Portamento On/Off
66	Sostenuto On/Off
67	Soft Pedal On/Off
68	Legato On/Off
69	Hold 2 On/Off
Low resolution continuous controllers
70	Sound Controller 1   (TG: Sound Variation;   FX: Exciter On/Off)
71	Sound Controller 2   (TG: Harmonic Content;   FX: Compressor On/Off)
72	Sound Controller 3   (TG: Release Time;   FX: Distortion On/Off)
73	Sound Controller 4   (TG: Attack Time;   FX: EQ On/Off)
74	Sound Controller 5   (TG: Brightness;   FX: Expander On/Off)
75	Sound Controller 6   (TG: Decay Time;   FX: Reverb On/Off)
76	Sound Controller 7   (TG: Vibrato Rate;   FX: Delay On/Off)
77	Sound Controller 8   (TG: Vibrato Depth;   FX: Pitch Transpose On/Off)
78	Sound Controller 9   (TG: Vibrato Delay;   FX: Flange/Chorus On/Off)
79	Sound Controller 10   (TG: Undefined;   FX: Special Effects On/Off)
80	General Purpose Controller 5
81	General Purpose Controller 6
82	General Purpose Controller 7
83	General Purpose Controller 8
84	Portamento Control (PTC)   (0vvvvvvv is the source Note number)   (Detail)
88	High Resolution Velocity Prefix
91	Effects 1 Depth (Reverb Send Level)
92	Effects 2 Depth (Tremolo Depth)
93	Effects 3 Depth (Chorus Send Level)
94	Effects 4 Depth (Celeste Depth)
95	Effects 5 Depth (Phaser Depth)
RPNs / NRPNs - (Detail)
96	Data Increment
97	Data Decrement
98	Non Registered Parameter Number (LSB)
99	Non Registered Parameter Number (MSB)
100	Registered Parameter Number (LSB)
101	Registered Parameter Number (MSB)
Channel Mode messages - (Detail)
120	All Sound Off
121	Reset All Controllers
122	Local Control On/Off
123	All Notes Off
124	Omni Mode Off (also causes ANO)
125	Omni Mode On (also causes ANO)
126	Mono Mode On (Poly Off; also causes ANO)
127	Poly Mode On (Mono Off; also causes ANO)
*/

/*
from http://midi.teragonaudio.com/tech/midispec.htm

Many (continuous) controller numbers are coarse adjustments, and have a respective fine adjustment controller number.
For example, controller #1 is the coarse adjustment for Modulation Wheel. Using this controller number in a message, a
device's Modulation Wheel can be adjusted in large (coarse) increments (ie, 128 steps). If finer adjustment (from a coarse setting)
needs to be made, then controller #33 is the fine adjust for Modulation Wheel. For controllers that have coarse/fine pairs of numbers,
there is thus a 14-bit resolution to the range. In other words, the Modulation Wheel can be set from 0x0000 to 0x3FFF (ie, one of 16,384 values).
For this 14-bit value, bits 7 to 13 are the coarse adjust, and bits 0 to 6 are the fine adjust. For example, to set the Modulation Wheel to 0x2005,
first you have to break it up into 2 bytes (as is done with Pitch Wheel messages). Take bits 0 to 6 and put them in a byte that is the fine adjust.
Take bits 7 to 13 and put them right-justified in a byte that is the coarse adjust. Assuming a MIDI channel of 0, here's the coarse and fine
Mod Wheel controller messages that a device would receive (coarse adjust first):

0xB0 0x01 0x40
Controller on chan 0, Mod Wheel coarse, bits 7 to 13 of 14-bit
value right-justified (with high bit clear).

0xB0 0x33 0x05
Controller on chan 0, Mod Wheel fine, bits 0 to 6 of 14-bit
value (with high bit clear).

Some devices do not implement fine adjust counterparts to coarse controllers. For example, some devices do not implement controller #33 for Mod Wheel
fine adjust. Instead the device only recognizes and responds to the Mod Wheel coarse controller number (#1). It is perfectly acceptable for devices
to only respond to the coarse adjustment for a controller if the device desires 7-bit (rather than 14-bit) resolution. The device should ignore that
controller's respective fine adjust message. By the same token, if it's only desirable to make fine adjustments to the Mod Wheel without changing its
current coarse setting (or vice versa), a device can be sent only a controller #33 message without a preceding controller #1 message (or vice versa).
Thus, if a device can respond to both coarse and fine adjustments for a particular controller (ie, implements the full 14-bit resolution), it should be
able to deal with either the coarse or fine controller message being sent without its counterpart following. The same holds true for other continuous
(ie, coarse/fine pairs of) controllers.

Note: In most MIDI literature, the coarse adjust is referred to with the designation "MSB" and the fine adjust is referred to with the designation "LSB".
I prefer the terms "coarse" and "fine".
*/

// coarse means MSB, fine means LSB
func BankSelectMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		0, value)
}
func ModulationWheelMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		1, value)
}
func BreathControllerMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		2, value)
}
func FootPedalMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		4, value)
}
func PortamentoTimeMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		5, value)
}
func DataEntryMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		6, value)
}
func VolumeMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		7, value)
}
func BalanceMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		8, value)
}
func PanPositionMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		10, value)
}
func ExpressionMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		11, value)
}
func EffectControl1MSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		12, value)
}
func EffectControl2MSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		13, value)
}
func GeneralPurposeSlider1(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		16, value)
}
func GeneralPurposeSlider2(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		17, value)
}
func GeneralPurposeSlider3(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		18, value)
}
func GeneralPurposeSlider4(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		19, value)
}
func BankSelectLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		32, value)
}
func ModulationWheelLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		33, value)
}
func BreathControllerLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		34, value)
}
func FootPedalLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		36, value)
}
func PortamentoTimeLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		37, value)
}
func DataEntryLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		38, value)
}
func VolumeLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		39, value)
}
func BalanceLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		40, value)
}
func PanPositionLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		42, value)
}
func ExpressionLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		43, value)
}
func EffectControl1LSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		44, value)
}
func EffectControl2LSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		45, value)
}
func HoldPedal_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		64, on)
}
func Portamento_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		65, on)
}
func SustenutoPedal_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		66, on)
}
func SoftPedal_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		67, on)
}
func LegatoPedal_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		68, on)
}
func Hold2Pedal_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		69, on)
}
func SoundVariation(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		70, value)
}
func SoundTimbre(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		71, value)
}
func SoundReleaseTime(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		72, value)
}
func SoundAttackTime(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		73, value)
}
func SoundBrightness(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		74, value)
}
func SoundControl6(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		75, value)
}
func SoundControl7(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		76, value)
}
func SoundControl8(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		77, value)
}
func SoundControl9(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		78, value)
}
func SoundControl10(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		79, value)
}
func GeneralPurposeButton1_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		80, on)
}
func GeneralPurposeButton2_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		81, on)
}
func GeneralPurposeButton3_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		82, on)
}
func GeneralPurposeButton4_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		83, on)
}
func EffectsLevel(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		91, value)
}
func TremuloLevel(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		92, value)
}
func ChorusLevel(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		93, value)
}
func CelesteLevel(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		94, value)
}
func PhaserLevel(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		95, value)
}
func DataButtonIncrement(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		96, value)
}
func DataButtonDecrement(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		97, value)
}
func NonRegisteredParameterLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		98, value)
}
func NonRegisteredParameterMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		99, value)
}
func RegisteredParameterLSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		100, value)
}
func RegisteredParameterMSB(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		101, value)
}
func AllSoundOff(ch channel.Channel) channel.ControlChange {
	return ch.ControlChange(
		120, 0)
}
func AllControllersOff(ch channel.Channel) channel.ControlChange {
	return ch.ControlChange(
		121, 0)
}
func LocalKeyboard_switch(ch channel.Channel, on bool) channel.ControlChange {
	return newControlCaseSwitch(ch,
		122, on)
}
func AllNotesOff(ch channel.Channel) channel.ControlChange {
	return ch.ControlChange(
		123, 0)
}
func OmniModeOff(ch channel.Channel) channel.ControlChange {
	return ch.ControlChange(
		124, 0)
}
func OmniModeOn(ch channel.Channel) channel.ControlChange {
	return ch.ControlChange(
		125, 0)
}
func MonoOperation(ch channel.Channel, value uint8) channel.ControlChange {
	return ch.ControlChange(
		126, value)
}
func PolyOperation(ch channel.Channel) channel.ControlChange {
	return ch.ControlChange(
		127, 0)
}
