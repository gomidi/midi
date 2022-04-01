package smf

import (
	"gitlab.com/gomidi/midi/v2"
)

const (
	Meta midi.Type = -5
)

const (

	// MetaChannelMsg is a MIDI channel meta message (which is a MetaMsg).
	// TODO add method to Message to get the channel number and document it.
	MetaChannelMsg midi.Type = 70 + iota

	// MetaCopyrightMsg is a MIDI copyright meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaCopyrightMsg

	// MetaCuepointMsg is a MIDI cuepoint meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaCuepointMsg

	// MetaDeviceMsg is a MIDI device meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaDeviceMsg

	// MetaEndOfTrackMsg is a MIDI end of track meta message (which is a MetaMsg).
	MetaEndOfTrackMsg

	// MetaInstrumentMsg is a MIDI instrument meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaInstrumentMsg

	// MetaKeySigMsg is a MIDI key signature meta message (which is a MetaMsg).
	// TODO add method to Message to get the key signature and document it.
	MetaKeySigMsg

	// MetaLyricMsg is a MIDI lyrics meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaLyricMsg

	// MetaTextMsg is a MIDI text meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaTextMsg

	// MetaMarkerMsg is a MIDI marker meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaMarkerMsg

	// MetaPortMsg is a MIDI port meta message (which is a MetaMsg).
	// TODO add method to Message to get the port number and document it.
	MetaPortMsg

	// MetaSeqNumberMsg is a MIDI sequencer number meta message (which is a MetaMsg).
	// TODO add method to Message to get the sequence number and document it.
	MetaSeqNumberMsg

	// MetaSeqDataMsg is a MIDI sequencer data meta message (which is a MetaMsg).
	// TODO add method to Message to get the sequencer data and document it.
	MetaSeqDataMsg

	// MetaTempoMsg is a MIDI tempo meta message (which is a MetaMsg).
	// The tempo in beats per minute of a concrete Message of this type can be retrieved via the BPM method of the Message.
	MetaTempoMsg

	// MetaTimeSigMsg is a MIDI time signature meta message (which is a MetaMsg).
	// The numerator, denominator, clocksPerClick and demiSemiQuaverPerQuarter of a concrete Message of this type can be retrieved via the TimeSig method of the Message.
	// A more comfortable way to get the meter is to use the Meter method of the Message.
	MetaTimeSigMsg

	// MetaTrackNameMsg is a MIDI track name meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaTrackNameMsg

	// MetaSMPTEOffsetMsg is a MIDI smpte offset meta message (which is a MetaMsg).
	// TODO add method to Message to get the smpte offset and document it.
	MetaSMPTEOffsetMsg

	// MetaUndefinedMsg is an undefined MIDI meta message (which is a MetaMsg).
	MetaUndefinedMsg

	// MetaProgramNameMsg is a MIDI program name meta message (which is a MetaMsg).
	// The text of a concrete Message of this type can be retrieved via the Text method of the Message.
	MetaProgramNameMsg
)

var msgTypeString = map[midi.Type]string{
	Meta:        "MetaType",
	MetaChannelMsg:     "MetaChannel",
	MetaCopyrightMsg:   "MetaCopyright",
	MetaCuepointMsg:    "MetaCuepoint",
	MetaDeviceMsg:      "MetaDevice",
	MetaEndOfTrackMsg:  "MetaEndOfTrack",
	MetaInstrumentMsg:  "MetaInstrument",
	MetaKeySigMsg:      "MetaKeySig",
	MetaLyricMsg:       "MetaLyric",
	MetaTextMsg:        "MetaText",
	MetaMarkerMsg:      "MetaMarker",
	MetaPortMsg:        "MetaPort",
	MetaSeqNumberMsg:   "MetaSeqNumber",
	MetaSeqDataMsg:     "MetaSeqData",
	MetaTempoMsg:       "MetaTempo",
	MetaTimeSigMsg:     "MetaTimeSig",
	MetaTrackNameMsg:   "MetaTrackName",
	MetaSMPTEOffsetMsg: "MetaSMPTEOffset",
	MetaUndefinedMsg:   "MetaUndefined",
	MetaProgramNameMsg: "MetaProgramName",
}

func init() {
	for ty, name := range msgTypeString {
		midi.AddTypeName(ty, name)
	}
}
