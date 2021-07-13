package midi

import (
	"bytes"
	"fmt"

	"gitlab.com/gomidi/midi/v2/internal/utils"
)

// Message represents a MIDI message. It can be created from the MIDI bytes of a message, by calling NewMessage.
type Message struct {

	// MsgType represents the message type of the MIDI message
	MsgType

	// Data contains the bytes of the MiDI message
	Data []byte
}

// NewMessage returns a new Message from the bytes of the message, by finding the correct type.
// If the type could not be found, the MsgType of the Message is UnknownMsg.
func NewMessage(data []byte) (m Message) {
	m.MsgType = GetMsgType(data)
	m.Data = data
	return
}

// NoteOn returns true if (and only if) the message is a NoteOnMsg.
// Then it also extracts the data to the given arguments
func (m Message) NoteOn(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOnMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *velocity = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return true
}

// NoteStart returns true if (and only if) the message is a NoteOnMsg with a velocity > 0.
// Then it also extracts the data to the given arguments
func (m Message) NoteStart(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOnMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *velocity = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	if *velocity == 0 {
		return false
	}
	return true
}

// NoteOff returns true if (and only if) the message is a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) NoteOff(channel, key, velocity *uint8) (is bool) {
	if !m.Is(NoteOffMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *velocity = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return true
}

// Channel returns true if (and only if) the message is a ChannelMsg.
// Then it also extracts the data to the given arguments
func (m Message) Channel(channel *uint8) (is bool) {
	if !m.Is(ChannelMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	return true
}

// NoteEnd returns true if (and only if) the message is a NoteOnMsg with a velocity == 0 or a NoteOffMsg.
// Then it also extracts the data to the given arguments
func (m Message) NoteEnd(channel, key, velocity *uint8) (is bool) {
	if !m.IsOneOf(NoteOnMsg, NoteOffMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *velocity = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return m.Is(NoteOffMsg) || *velocity == 0
}

// PolyAfterTouch returns true if (and only if) the message is a PolyAfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) PolyAfterTouch(channel, key, pressure *uint8) (is bool) {
	if !m.Is(PolyAfterTouchMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*key, *pressure = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return true
}

// AfterTouch returns true if (and only if) the message is a AfterTouchMsg.
// Then it also extracts the data to the given arguments
func (m Message) AfterTouch(channel, pressure *uint8) (is bool) {
	if !m.Is(AfterTouchMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*pressure = utils.ParseUint7(m.Data[1])
	return true
}

// ProgramChange returns true if (and only if) the message is a ProgramChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) ProgramChange(channel, program *uint8) (is bool) {
	if !m.Is(ProgramChangeMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*program = utils.ParseUint7(m.Data[1])
	return true
}

// PitchBend returns true if (and only if) the message is a PitchBendMsg.
// Then it also extracts the data to the given arguments
// Either relative or absolute may be nil, if not needed.
func (m Message) PitchBend(channel *uint8, relative *int16, absolute *uint16) (is bool) {
	if !m.Is(PitchBendMsg) {
		return false
	}

	rel, abs := utils.ParsePitchWheelVals(m.Data[1], m.Data[2])
	if relative != nil {
		*relative = rel
	}
	if absolute != nil {
		*absolute = abs
	}
	return true
}

// ControlChange returns true if (and only if) the message is a ControlChangeMsg.
// Then it also extracts the data to the given arguments
func (m Message) ControlChange(channel, controller, value *uint8) (is bool) {
	if !m.Is(ControlChangeMsg) {
		return false
	}

	_, *channel = utils.ParseStatus(m.Data[0])
	*controller, *value = utils.ParseTwoUint7(m.Data[1], m.Data[2])
	return true
}

// Tempo returns true if (and only if) the message is a MetaTempoMsg.
// Then it also extracts the data to the given arguments
func (m Message) Tempo(bpm *float64) (is bool) {
	if !m.Is(MetaTempoMsg) {
		return false
	}

	rd := bytes.NewReader(m.metaDataWithoutVarlength())
	microsecondsPerCrotchet, err := utils.ReadUint24(rd)
	if err != nil {
		return false
	}

	*bpm = float64(60000000) / float64(microsecondsPerCrotchet)

	return true
}

// Lyric returns true if (and only if) the message is a MetaLyricMsg.
// Then it also extracts the data to the given arguments
func (m Message) Lyric(text *string) (is bool) {
	if !m.Is(MetaLyricMsg) {
		return false
	}
	m.text(text)
	return true
}

// Copyright returns true if (and only if) the message is a MetaCopyrightMsg.
// Then it also extracts the data to the given arguments
func (m Message) Copyright(text *string) (is bool) {
	if !m.Is(MetaCopyrightMsg) {
		return false
	}
	m.text(text)
	return true
}

// Cuepoint returns true if (and only if) the message is a MetaCuepointMsg.
// Then it also extracts the data to the given arguments
func (m Message) Cuepoint(text *string) (is bool) {
	if !m.Is(MetaCuepointMsg) {
		return false
	}
	m.text(text)
	return true
}

// Device returns true if (and only if) the message is a MetaDeviceMsg.
// Then it also extracts the data to the given arguments
func (m Message) Device(text *string) (is bool) {
	if !m.Is(MetaDeviceMsg) {
		return false
	}
	m.text(text)
	return true
}

// Instrument returns true if (and only if) the message is a MetaInstrumentMsg.
// Then it also extracts the data to the given arguments
func (m Message) Instrument(text *string) (is bool) {
	if !m.Is(MetaInstrumentMsg) {
		return false
	}
	m.text(text)
	return true
}

// Marker returns true if (and only if) the message is a MetaMarkerMsg.
// Then it also extracts the data to the given arguments
func (m Message) Marker(text *string) (is bool) {
	if !m.Is(MetaMarkerMsg) {
		return false
	}
	m.text(text)
	return true
}

// ProgramName returns true if (and only if) the message is a MetaProgramNameMsg.
// Then it also extracts the data to the given arguments
func (m Message) ProgramName(text *string) (is bool) {
	if !m.Is(MetaProgramNameMsg) {
		return false
	}
	m.text(text)
	return true
}

// Text returns true if (and only if) the message is a MetaTextMsg.
// Then it also extracts the data to the given arguments
func (m Message) Text(text *string) (is bool) {
	if !m.Is(MetaTextMsg) {
		return false
	}
	m.text(text)
	return true
}

// TrackName returns true if (and only if) the message is a MetaTrackNameMsg.
// Then it also extracts the data to the given arguments
func (m Message) TrackName(text *string) (is bool) {
	if !m.Is(MetaTrackNameMsg) {
		return false
	}
	m.text(text)
	return true
}

func (m Message) text(text *string) {
	*text, _ = utils.ReadText(bytes.NewReader(m.Data[2:]))
	return
}

/*
MTC Quarter Frame

These are the MTC (i.e. SMPTE based) equivalent of the F8 Timing Clock messages,
though offer much higher resolution, as they are sent at a rate of 96 to 120 times
a second (depending on the SMPTE frame rate). Each Quarter Frame message provides
partial timecode information, 8 sequential messages being required to fully
describe a timecode instant. The reconstituted timecode refers to when the first
partial was received. The most significant nibble of the data byte indicates the
partial (aka Message Type).

Partial	Data byte	Usage
1	0000 bcde	Frame number LSBs 	abcde = Frame number (0 to frameRate-1)
2	0001 000a	Frame number MSB
3	0010 cdef	Seconds LSBs 	abcdef = Seconds (0-59)
4	0011 00ab	Seconds MSBs
5	0100 cdef	Minutes LSBs 	abcdef = Minutes (0-59)
6	0101 00ab	Minutes MSBs
7	0110 defg	Hours LSBs 	ab = Frame Rate (00 = 24, 01 = 25, 10 = 30drop, 11 = 30nondrop)
cdefg = Hours (0-23)
8	0111 0abc	Frame Rate, and Hours MSB
*/

// MTC represents a MIDI timing code message (quarter frame)
func (m Message) MTC(quarterframe *uint8) (is bool) {
	if !m.Is(MTCMsg) {
		return false
	}

	*quarterframe = utils.ParseUint7(m.Data[1])
	return true
}

// Song returns the song number of a MIDI song select system message
func (m Message) SongSelect(song *uint8) (is bool) {
	if !m.Is(SongSelectMsg) {
		return false
	}

	*song = utils.ParseUint7(m.Data[1])
	return true
}

// SPP returns the song position pointer of a MIDI song position pointer system message
func (m Message) SPP(spp *uint16) (is bool) {
	if !m.Is(SPPMsg) {
		return false
	}

	_, *spp = utils.ParsePitchWheelVals(m.Data[2], m.Data[1])
	return true
}

// Meter returns the meter of a MetaTimeSigMsg.
// For other messages, it returns 0,0.
func (m Message) Meter(num, denom *uint8) (is bool) {
	return m.TimeSig(num, denom, nil, nil)
}

// metaData strips away the meta byte and the metatype byte and the varlength byte
func (m Message) metaDataWithoutVarlength() []byte {
	//fmt.Printf("original data: % X\n", m.Data)
	return m.Data[3:]
}

// TimeSig returns the numerator, denominator, clocksPerClick and demiSemiQuaverPerQuarter of a
// MetaTimeSigMsg. For other messages, it returns 0,0,0,0.
func (m Message) TimeSig(numerator, denominator, clocksPerClick, demiSemiQuaverPerQuarter *uint8) (is bool) {
	if !m.Is(MetaTimeSigMsg) {
		//fmt.Println("not timesig message")
		return false
	}

	data := m.metaDataWithoutVarlength()

	if len(data) != 4 {
		//fmt.Printf("not correct data lenght: % X \n", data)
		//err = unexpectedMessageLengthError("TimeSignature expected length 4")
		return false
	}

	//fmt.Printf("TimeSigData: % X\n", data)

	*numerator = data[0]
	*denominator = data[1]
	if clocksPerClick != nil {
		*clocksPerClick = data[2]
	}
	if demiSemiQuaverPerQuarter != nil {
		*demiSemiQuaverPerQuarter = data[3]
	}
	*denominator = bin2decDenom(*denominator)
	return true
}

// String represents the Message as a string that contains the MsgType and its properties.
func (m Message) String() string {
	var bf bytes.Buffer
	fmt.Fprintf(&bf, m.MsgType.String())

	var channel, val1, val2 uint8
	var pitchabs uint16
	var pitchrel int16
	var text string
	var bpm float64
	var spp uint16

	switch {
	case m.NoteOn(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " key: %v velocity: %v", val1, val2)
	case m.NoteOff(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " key: %v velocity: %v", val1, val2)
	case m.PolyAfterTouch(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " key: %v pressure: %v", val1, val2)
	case m.ControlChange(&channel, &val1, &val2):
		fmt.Fprintf(&bf, " controller: %v value: %v", val1, val2)
	case m.ProgramChange(&channel, &val1):
		fmt.Fprintf(&bf, " program: %v", val1)
	case m.PitchBend(&channel, &pitchrel, &pitchabs):
		fmt.Fprintf(&bf, " pitch: %v (%v)", pitchrel, pitchabs)
	case m.Tempo(&bpm):
		fmt.Fprintf(&bf, " bpm: %0.2f", bpm)
	case m.Meter(&val1, &val2):
		fmt.Fprintf(&bf, " meter: %v/%v", val1, val2)
	case m.IsOneOf(MetaLyricMsg, MetaMarkerMsg, MetaCopyrightMsg, MetaTextMsg, MetaCuepointMsg, MetaDeviceMsg, MetaInstrumentMsg, MetaProgramNameMsg, MetaTrackNameMsg):
		m.text(&text)
		fmt.Fprintf(&bf, " text: %q", text)
	case m.MTC(&val1):
		fmt.Fprintf(&bf, " mtc: %v", val1)
	case m.SPP(&spp):
		fmt.Fprintf(&bf, " spp: %v", spp)
	case m.SongSelect(&val1):
		fmt.Fprintf(&bf, " song: %v", val1)
	default:
	}

	return bf.String()
}
