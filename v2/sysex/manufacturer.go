package sysex

type ManufacturerID byte

const (
	ExtendedRange ManufacturerID = 0 // look for the next two bytes

	SequentialCircuits ManufacturerID = 1
	BigBriar           ManufacturerID = 2
	Octave_Plateau     ManufacturerID = 3
	Moog               ManufacturerID = 4
	PassportDesigns    ManufacturerID = 5
	Lexicon            ManufacturerID = 6
	Kurzweil           ManufacturerID = 7
	Fender             ManufacturerID = 8
	Gulbransen         ManufacturerID = 9
	DeltaLabs          ManufacturerID = 0x0A
	SoundComp          ManufacturerID = 0x0B
	GeneralElectro     ManufacturerID = 0x0C
	Techmar            ManufacturerID = 0x0D
	MatthewsResearch   ManufacturerID = 0x0E
	Oberheim           ManufacturerID = 0x10
	PAIA               ManufacturerID = 0x11
	Simmons            ManufacturerID = 0x12
	DigiDesign         ManufacturerID = 0x13
	Fairlight          ManufacturerID = 0x14
	Peavey             ManufacturerID = 0x1B
	JLCooper           ManufacturerID = 0x15
	Lowery             ManufacturerID = 0x16
	Lin                ManufacturerID = 0x17
	Emu                ManufacturerID = 0x18
	BonTempi           ManufacturerID = 0x20
	SIEL               ManufacturerID = 0x21
	SyntheAxe          ManufacturerID = 0x23
	Hohner             ManufacturerID = 0x24
	Crumar             ManufacturerID = 0x25
	Solton             ManufacturerID = 0x26
	JellinghausMs      ManufacturerID = 0x27
	CTS                ManufacturerID = 0x28
	PPG                ManufacturerID = 0x29
	Elka               ManufacturerID = 0x2F
	Cheetah            ManufacturerID = 0x36
	Waldorf            ManufacturerID = 0x3E
	Kawai              ManufacturerID = 0x40
	Roland             ManufacturerID = 0x41
	Korg               ManufacturerID = 0x42
	Yamaha             ManufacturerID = 0x43
	Casio              ManufacturerID = 0x44
	Akai               ManufacturerID = 0x45

	EducationalUse ManufacturerID = 0x7D // not for commercial use

	// Universal SysEx (not manufacturer specific)
	RealTimeID    ManufacturerID = 0x7F
	NonRealTimeID ManufacturerID = 0x7E
)

var manuIDNames = map[ManufacturerID]string{
	ExtendedRange:      "ExtendedRange",
	SequentialCircuits: "SequentialCircuits",
	BigBriar:           "BigBriar",
	Octave_Plateau:     "Octave_Plateau",
	Moog:               "Moog",
	PassportDesigns:    "PassportDesigns",
	Lexicon:            "Lexicon",
	Kurzweil:           "Kurzweil",
	Fender:             "Fender",
	Gulbransen:         "Gulbransen",
	DeltaLabs:          "DeltaLabs",
	SoundComp:          "SoundComp",
	GeneralElectro:     "GeneralElectro",
	Techmar:            "Techmar",
	MatthewsResearch:   "MatthewsResearch",
	Oberheim:           "Oberheim",
	PAIA:               "PAIA",
	Simmons:            "Simmons",
	DigiDesign:         "DigiDesign",
	Fairlight:          "Fairlight",
	Peavey:             "Peavey",
	JLCooper:           "JLCooper",
	Lowery:             "Lowery",
	Lin:                "Lin",
	Emu:                "Emu",
	BonTempi:           "BonTempi",
	SIEL:               "SIEL",
	SyntheAxe:          "SyntheAxe",
	Hohner:             "Hohner",
	Crumar:             "Crumar",
	Solton:             "Solton",
	JellinghausMs:      "JellinghausMs",
	CTS:                "CTS",
	PPG:                "PPG",
	Elka:               "Elka",
	Cheetah:            "Cheetah",
	Waldorf:            "Waldorf",
	Kawai:              "Kawai",
	Roland:             "Roland",
	Korg:               "Korg",
	Yamaha:             "Yamaha",
	Casio:              "Casio",
	Akai:               "Akai",
	EducationalUse:     "EducationalUse",
	RealTimeID:         "RealTimeID",
	NonRealTimeID:      "NonRealTimeID",
}

func (m ManufacturerID) String() string {
	s, has := manuIDNames[m]
	if has {
		return s
	} else {
		return "unknown"
	}
}

// see http://midi.teragonaudio.com/tech/midispec.htm

/*

You'll note that we use only one byte for the Manufacturer ID. And since a midi data byte can't be greater than 0x7F,
that means we have only 127 IDs to dole out to manufacturers. Well, there are more than 127 manufacturers of MIDI products.

To accomodate a greater range of manufacturer IDs, the MMA decided to reserve a manufacturer ID of 0 for a special purpose.
When you see a manufacturer ID of 0, then there will be two more data bytes after this. These two data bytes combine to
make the real manufacturer ID. So, some manufacturers have IDs that are 3 bytes, where the first byte is always 0.
Using this "trick", the range of unique manufacturer IDs is extended to accomodate over 16,000 MIDI manufacturers.

For example, Microsoft's manufacturer ID consists of the 3 bytes 0x00 0x00 0x41. Note that the first byte is 0 to
indicate that the real ID is 0x0041, but is still different than Roland's ID which is only the single byte of 0x41.

A manufacturer must get a registered ID from the MMA if he wants to define his own SysEx messages, or use the
following:

  Educational Use  0x7D

This ID is for educational or development use only, and should never appear in a commercial design.

On the other hand, it is permissible to use another manufacturer's defined SysEx message(s) in your own products.
For example, if the Roland S-770 has a particular SysEx message that you could use verbatim in your own design,
you're free to use that message (and therefore the Roland ID in it). But, you're not allowed to transmit a mutated
version of any Roland message with a Roland ID. Only Roland can develop new messages that contain a Roland ID.

*/

/*
Universal SysEx (not manufacturer specific) 0x7F and 0x7E


A general template for these two IDs was defined. After the ID byte is a SysEx Channel byte. This could be
from 0 to 127 for a total of 128 SysEx channels. So, although "normal" SysEx messages have no MIDI channel like
Voice Category messages do, a Universal SysEx message can be sent on one of 128 SysEx channels. This allows the
musician to set various devices to ignore certain Universal SysEx messages (ie, if the device allows the musician
to set its Base SysEx Channel. Most devices just set their Base Sysex channel to the same number as the Base Channel
for Voice Category messages). On the other hand, a SysEx channel of 127 is actually meant to tell the device to
"disregard the channel and pay attention to this message regardless".

After the SysEx channel, the next two bytes are Sub IDs which tell what the SysEx is for. There are several
Sub IDs defined for particular messages. There is a Sub ID for a Universal SysEx message to set a device's master
volume. (This is different than Volume controller which sets the volume for only one particular MIDI channel).
There is a Sub ID for a Universal SysEx message to set a device's Pitch Wheel bend range. There are a couple of
Sub IDs for some Universal SysEx messages to implement a waveform (sample) dump over MIDI. Etc.

*/
