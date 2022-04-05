package gm

import (
	"gitlab.com/gomidi/midi/v2"
)

// GMProgram is a shortcut to write GM bank select control change message followed
// by a program change.
func GMProgram(ch, prog uint8) (msgs []midi.Message) {
	//c := channel.Channel(ch)
	msgs = append(msgs, midi.ControlChange(ch, midi.BankSelectMSB, 0))
	msgs = append(msgs, midi.ProgramChange(ch, prog))
	return
}

// Reset writes a kind of somewhat homegrown GM/GS reset message.
// The idea is inspired by http://www.artandscienceofsound.com/article/standardmidifiles.
// The following messages will be written to the writer on the given channel:
/*
     cc bank select 0
	 program change prog
	 cc all controllers off
	 cc volume 100
	 cc expression 127
	 cc hold pedal 0
	 cc pan position 64
*/
func Reset(ch, prog uint8) []midi.Message {
	return []midi.Message{
		midi.ControlChange(ch, midi.BankSelectMSB, 0),
		midi.ProgramChange(ch, prog),
		midi.ControlChange(ch, midi.AllControllersOff, 0),
		midi.ControlChange(ch, midi.VolumeMSB, 100),
		midi.ControlChange(ch, midi.ExpressionMSB, 127),
		midi.ControlChange(ch, midi.HoldPedalSwitch, 0),
		midi.ControlChange(ch, midi.PanPositionMSB, 64),
	}
}

/*
default of 2 semitone
Pitch Bend Range can be set by sending MIDI controller messages. Specifically, you do it with Registered Parameters (cc# 100 and 101).

On the MIDI channel in question, you need to send:
MIDI cc100 	= 0
MIDI cc101 	= 0
MIDI cc6 	= value of desired bend range (in semitones)

Example: Lets say you want to set the bend range to 2 semi-tones. First you send cc# 100 with a value of 0; then cc#101 with a value of 0. This turns on reception for setting pitch bend with the Data controller (#6). Then you send cc# 6 with a value of 2 (in semitones; this will give you a whole step up and a whole step down from the center).

Once you have set the bend range the way you want, then you send controller 100 or 101 with a value of 127 so that any further messages of controller 6 (which you might be using for other stuff) won't change the bend range.
*/

/*
from http://www.artandscienceofsound.com/article/standardmidifiles

Depending upon the application you are using to create the file in the first place, header information may automatically be saved from within parameters set in the application, or may need to be placed in a ‘set-up’ bar before the music data commences.

Either way, information that should be considered includes:

GM/GS Reset message

Per MIDI Channel
Bank Select (0=GM) / Program Change #
Reset All Controllers (not all devices may recognize this command so you may prefer to zero out or reset individual controllers)
Initial Volume (CC7) (standard level = 100)
Expression (CC11) (initial level set to 127)
Hold pedal (0 = off)
Pan (Center = 64)
Modulation (0)
Pitch bend range
Reverb (0 = off)
Chorus level (0 = off)
*/
