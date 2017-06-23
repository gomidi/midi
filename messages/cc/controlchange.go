package cc

import (
	"github.com/gomidi/midi/messages/channel"
)

func newControlCaseSwitch(ch channel.Channel, controller uint8, on bool) channel.ControlChange {
	var value uint8 // start with off = 0
	if on {
		value = 127 // on = 127
	}
	return ch.ControlChange(controller, value)
}

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
