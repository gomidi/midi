package smf

import (
	"gitlab.com/gomidi/midi/v2"
)

const (
	MetaType midi.Type = -5
)

const (

	// MetaChannelMsg is a MIDI channel meta message (which is a MetaMsg).
	// TODO add method to Message to get the channel number and document it.
	MetaChannel midi.Type = 70 + iota

	// MetaCopyrightMsg is a MIDI copyright meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaCopyright

	// MetaCuepointMsg is a MIDI cuepoint meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaCuepoint

	// MetaDeviceMsg is a MIDI device meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaDevice

	// MetaEndOfTrackMsg is a MIDI end of track meta message (which is a MetaMsg).
	MetaEndOfTrack

	// MetaInstrumentMsg is a MIDI instrument meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaInstrument

	// MetaKeySigMsg is a MIDI key signature meta message (which is a MetaMsg).
	// TODO add method to Message to get the key signature and document it.
	MetaKeySig

	// MetaLyricMsg is a MIDI lyrics meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaLyric

	// MetaTextMsg is a MIDI text meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaText

	// MetaMarkerMsg is a MIDI marker meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaMarker

	// MetaPortMsg is a MIDI port meta message (which is a MetaMsg).
	// TODO add method to Message to get the port number and document it.
	MetaPort

	// MetaSeqNumberMsg is a MIDI sequencer number meta message (which is a MetaMsg).
	// TODO add method to Message to get the sequence number and document it.
	MetaSeqNumber

	// MetaSeqDataMsg is a MIDI sequencer data meta message (which is a MetaMsg).
	// TODO add method to Message to get the sequencer data and document it.
	MetaSeqData

	// MetaTempoMsg is a MIDI tempo meta message (which is a MetaMsg).
	// The tempo in beats per minute of a concrete Message of this type can be retrieved via the BPM method of the Message.
	MetaTempo

	// MetaTimeSigMsg is a MIDI time signature meta message (which is a MetaMsg).
	// The numerator, denominator, clocksPerClick and demiSemiQuaverPerQuarter of a concrete Message of this type can be retrieved via the TimeSig method of the Message.
	// A more comfortable way to get the meter is to use the Meter method of the Message.
	MetaTimeSig

	// MetaTrackNameMsg is a MIDI track name meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaTrackName

	// MetaSMPTEOffsetMsg is a MIDI smpte offset meta message (which is a MetaMsg).
	// TODO add method to Message to get the smpte offset and document it.
	MetaSMPTEOffset

	// MetaUndefinedMsg is an undefined MIDI meta message (which is a MetaMsg).
	MetaUndefined

	// MetaProgramNameMsg is a MIDI program name meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaProgramName
)

var msgTypeString = map[midi.Type]string{
	MetaType:        "MetaType",
	MetaChannel:     "MetaChannel",
	MetaCopyright:   "MetaCopyright",
	MetaCuepoint:    "MetaCuepoint",
	MetaDevice:      "MetaDevice",
	MetaEndOfTrack:  "MetaEndOfTrack",
	MetaInstrument:  "MetaInstrument",
	MetaKeySig:      "MetaKeySig",
	MetaLyric:       "MetaLyric",
	MetaText:        "MetaText",
	MetaMarker:      "MetaMarker",
	MetaPort:        "MetaPort",
	MetaSeqNumber:   "MetaSeqNumber",
	MetaSeqData:     "MetaSeqData",
	MetaTempo:       "MetaTempo",
	MetaTimeSig:     "MetaTimeSig",
	MetaTrackName:   "MetaTrackName",
	MetaSMPTEOffset: "MetaSMPTEOffset",
	MetaUndefined:   "MetaUndefined",
	MetaProgramName: "MetaProgramName",
}

func init() {
	for ty, name := range msgTypeString {
		midi.AddTypeName(ty, name)
	}
}
