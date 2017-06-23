package meta

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gomidi/midi/internal/lib"
	"io"
	"math/big"
)

type Message interface {
	String() string
	Raw() []byte
	meta() // just to tell that it is a meta message
	readFrom(io.Reader) (Message, error)
}

func ReadFrom(ev Message, rd io.Reader) (Message, error) {
	if ev != nil {
		return ev.readFrom(rd)
	}
	// TODO: read an unknown meta message
	return nil, nil
}

const (
	degreeC  = 0
	degreeCs = 1
	degreeDf = degreeCs
	degreeD  = 2
	degreeDs = 3
	degreeEf = degreeDs
	degreeE  = 4
	degreeF  = 5
	degreeFs = 6
	degreeGf = degreeFs
	degreeG  = 7
	degreeGs = 8
	degreeAf = degreeGs
	degreeA  = 9
	degreeAs = 10
	degreeBf = degreeAs
	degreeB  = 11
	degreeCf = degreeB
)

// Supplied to KeySignature
const (
	majorMode = 0
	minorMode = 1
)

var (
	_ Message = Text("")
	_ Message = Copyright("")
	_ Message = Sequence("")
	_ Message = TrackInstrument("")
	_ Message = Marker("")
	_ Message = Lyric("")
	_ Message = CuePoint("")
	_ Message = SequenceNumber(0)
	_ Message = MIDIChannel(0)
	_ Message = DevicePort("")
	_ Message = MIDIPort(0)
	_ Message = Tempo(0)
	_ Message = TimeSignature{}
	_ Message = KeySignature{}
	_ Message = EndOfTrack
)

const (
	// End of track
	// the handler is supposed to keep track of the current track

	byteEndOfTrack            = byte(0x2F)
	byteSequenceNumber        = byte(0x00)
	byteText                  = byte(0x01)
	byteCopyright             = byte(0x02)
	byteSequence              = byte(0x03)
	byteTrackInstrument       = byte(0x04)
	byteLyric                 = byte(0x05)
	byteMarker                = byte(0x06)
	byteCuePoint              = byte(0x07)
	byteMIDIChannel           = byte(0x20)
	byteDevicePort            = byte(0x9)
	byteMIDIPort              = byte(0x21)
	byteTempo                 = byte(0x51)
	byteTimeSignature         = byte(0x58)
	byteKeySignature          = byte(0x59)
	byteSequencerSpecificInfo = byte(0x7F)
)

var metaMessages = map[byte]Message{
	byteEndOfTrack:            EndOfTrack,
	byteSequenceNumber:        SequenceNumber(0),
	byteText:                  Text(""),
	byteCopyright:             Copyright(""),
	byteSequence:              Sequence(""),
	byteTrackInstrument:       TrackInstrument(""),
	byteLyric:                 Lyric(""),
	byteMarker:                Marker(""),
	byteCuePoint:              CuePoint(""),
	byteMIDIChannel:           MIDIChannel(0),
	byteDevicePort:            DevicePort(""),
	byteMIDIPort:              MIDIPort(0),
	byteTempo:                 Tempo(0),
	byteTimeSignature:         TimeSignature{},
	byteKeySignature:          KeySignature{},
	byteSequencerSpecificInfo: nil, // SequencerSpecificInfo
}

func Dispatch(b byte) Message {
	// fmt.Printf("got meta byte: % X\n", b)
	return metaMessages[b]
}

type endOfTrack bool

const (
	EndOfTrack = endOfTrack(true)
)

func (m endOfTrack) String() string {
	return fmt.Sprintf("%T", m)
}

func (m endOfTrack) Raw() []byte {
	return (&lib.MetaMessage{
		Typ: byte(byteEndOfTrack),
	}).Bytes()
}

func (m endOfTrack) meta() {}

func (m endOfTrack) readFrom(rd io.Reader) (Message, error) {

	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 0 {
		err = lib.UnexpectedMessageLengthError("EndOfTrack expected length 0")
		return nil, err
	}

	return m, nil
}

type MIDIPort uint8

func (m MIDIPort) Number() uint8 {
	return uint8(m)
}

func (m MIDIPort) String() string {
	return fmt.Sprintf("%T: %v", m, uint8(m))
}

func (m MIDIPort) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byte(byteMIDIPort),
		Data: []byte{byte(m)},
	}).Bytes()
}

func (m MIDIPort) meta() {}

func (m MIDIPort) readFrom(rd io.Reader) (Message, error) {

	// Obsolete 'MIDI Port'
	//	we can't ignore it, since it advanced in deltatime

	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 1 {
		return nil, lib.UnexpectedMessageLengthError("MIDI Port Message expected length 1")
	}

	var port uint8
	port, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return MIDIPort(port), nil

}

type MIDIChannel uint8

func (m MIDIChannel) Number() uint8 {
	return uint8(m)
}

func (m MIDIChannel) String() string {
	return fmt.Sprintf("%T: %#v", m, uint8(m))
}

func (m MIDIChannel) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byte(byteMIDIChannel),
		Data: []byte{byte(m)},
	}).Bytes()
}

func (m MIDIChannel) meta() {}

func (m MIDIChannel) readFrom(rd io.Reader) (Message, error) {

	// Obsolete 'MIDI Channel'
	//	we can't ignore it, since it advanced in deltatime

	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 1 {
		return nil, lib.UnexpectedMessageLengthError("Midi Channel Message expected length 1")
	}

	var ch uint8
	ch, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return MIDIChannel(ch), nil

}

// value is equal to BPM
type Tempo uint32

func (t Tempo) BPM() uint32 {
	return uint32(t)
}

func (m Tempo) String() string {
	return fmt.Sprintf("%T BPM: %v", m, m.BPM())
}

func (m Tempo) Raw() []byte {

	f := float64(60000000) / float64(m.BPM())

	muSecPerQuarterNote := uint32(f)

	if muSecPerQuarterNote > 0xFFFFFF {
		muSecPerQuarterNote = 0xFFFFFF
	}
	b4 := big.NewInt(int64(muSecPerQuarterNote)).Bytes()
	var b = []byte{0, 0, 0}
	switch len(b4) {
	case 0:
	case 1:
		b[2] = b4[0]
	case 2:
		b[2] = b4[1]
		b[1] = b4[0]
	case 3:
		b[2] = b4[2]
		b[1] = b4[1]
		b[0] = b4[0]
	}

	return (&lib.MetaMessage{
		Typ:  byteTempo,
		Data: b,
	}).Bytes()
}

func (m Tempo) meta() {}

func (m Tempo) readFrom(rd io.Reader) (Message, error) {
	// TODO TEST
	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 3 {
		err = lib.UnexpectedMessageLengthError("Tempo expected length 3")
		return nil, err
	}

	var microsecondsPerCrotchet uint32
	microsecondsPerCrotchet, err = lib.ReadUint24(rd)

	if err != nil {
		return nil, err
	}

	// Also beats per minute
	var bpm uint32
	bpm = 60000000 / microsecondsPerCrotchet

	return Tempo(bpm), nil
}

type SequenceNumber uint16

func (s SequenceNumber) Number() uint16 {
	return uint16(s)
}

func (m SequenceNumber) String() string {
	return fmt.Sprintf("SequenceNumber: %v", m.Number())
}

func (m SequenceNumber) Raw() []byte {
	var bf bytes.Buffer
	binary.Write(&bf, binary.BigEndian, m.Number())
	return (&lib.MetaMessage{
		Typ:  byteSequenceNumber,
		Data: bf.Bytes(),
	}).Bytes()
}

func (m SequenceNumber) readFrom(rd io.Reader) (Message, error) {
	length, err := lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	// Zero length sequences allowed according to http://home.roadrunner.com/~jgglatt/tech/midifile/seq.htm
	if length == 0 {
		return SequenceNumber(0), nil
	}

	// Otherwise length will be 2 to hold the uint16.
	var sequenceNumber uint16
	sequenceNumber, err = lib.ReadUint16(rd)

	if err != nil {
		return nil, err
	}

	return SequenceNumber(sequenceNumber), nil
}

func (m SequenceNumber) meta() {}

type metaTimeCodeQuarter struct {
	Type   uint8
	Values uint8
}

// TODO check and implement New* function

func (m metaTimeCodeQuarter) String() string {
	return fmt.Sprintf("%#v", m)
}

func (m metaTimeCodeQuarter) meta() {}

type TimeSignature struct {
	Numerator   uint8
	Denominator uint8
	// ClocksPerClick           uint8
	// DemiSemiQuaverPerQuarter uint8
}

/*
func NewTimeSignature(num uint8, denom uint8) TimeSignature {
	return TimeSignature{Numerator: num, Denominator: denom}
}
*/

// bin2decDenom converts the binary denominator to the decimal
func bin2decDenom(bin uint8) uint8 {
	if bin == 0 {
		return 1
	}
	return 2 << (bin - 1)
}

// dec2binDenom converts the decimal denominator to the binary one
// it works, use it!
func dec2binDenom(dec uint8) (bin uint8) {
	if dec <= 1 {
		return 0
	}
	for dec > 2 {
		bin++
		dec = dec >> 1

	}
	return bin + 1
}

func (m TimeSignature) Raw() []byte {
	// cpcl := m.ClocksPerClick
	// if cpcl == 0 {
	cpcl := byte(8)
	// }

	// dsqpq := m.DemiSemiQuaverPerQuarter
	// if dsqpq == 0 {
	dsqpq := byte(8)
	// }

	var denom = dec2binDenom(m.Denominator)

	return (&lib.MetaMessage{
		Typ:  byteTimeSignature,
		Data: []byte{m.Numerator, denom, cpcl, dsqpq},
	}).Bytes()

}

func (m TimeSignature) String() string {
	//return fmt.Sprintf("%T %v/%v clocksperclick %v dsqpq %v", m, m.Numerator, m.Denominator, m.ClocksPerClick, m.DemiSemiQuaverPerQuarter)
	return fmt.Sprintf("%T %v/%v", m, m.Numerator, m.Denominator)
}

func (m TimeSignature) readFrom(rd io.Reader) (Message, error) {
	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 4 {
		err = lib.UnexpectedMessageLengthError("TimeSignature expected length 4")
		return nil, err
	}

	// TODO TEST
	var numerator uint8
	numerator, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	var denomenator uint8
	denomenator, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	var clocksPerClick uint8
	clocksPerClick, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	var demiSemiQuaverPerQuarter uint8
	demiSemiQuaverPerQuarter, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	// TODO: do something with clocksPerClick and demiSemiQuaverPerQuarter
	var _ = clocksPerClick
	var _ = demiSemiQuaverPerQuarter

	return TimeSignature{
		Numerator:   numerator,
		Denominator: 2 << (denomenator - 1),
		// ClocksPerClick:           clocksPerClick,
		// DemiSemiQuaverPerQuarter: demiSemiQuaverPerQuarter,
	}, nil

}

func (m TimeSignature) meta() {}

type KeySignature struct {
	Key     uint8
	IsMajor bool
	Num     uint8
	//	SharpsOrFlats int8
	IsFlat bool
}

/*
// NewKeySignature returns a key signature event.
// key is the key of the scale (C=0 add the corresponding number of semitones). ismajor indicates whether it is a major or minor scale
// num is the number of accidentals. isflat indicates whether the accidentals are flats or sharps
func NewKeySignature(key uint8, ismajor bool, num uint8, isflat bool) KeySignature {
	return KeySignature{Key: key, IsMajor: ismajor, Num: num, IsFlat: isflat}
}
*/

func (m KeySignature) Raw() []byte {
	mi := int8(0)
	if !m.IsMajor {
		mi = 1
	}
	sf := int8(m.Num)

	if m.IsFlat {
		sf = sf * (-1)
	}

	return (&lib.MetaMessage{
		Typ:  byteKeySignature,
		Data: []byte{byte(sf), byte(mi)},
	}).Bytes()
}

func (m KeySignature) String() string {
	return fmt.Sprintf("%T: %s", m, m.Text())
}

func (m KeySignature) Note() (note string) {
	switch m.Key {
	case degreeC:
		note = "C"
	case degreeD:
		note = "D"
	case degreeE:
		note = "E"
	case degreeF:
		note = "F"
	case degreeG:
		note = "G"
	case degreeA:
		note = "A"
	case degreeB:
		note = "B"
	case degreeCs:
		note = "C♯"
		if m.IsFlat {
			note = "D♭"
		}
	case degreeDs:
		note = "D♯"
		if m.IsFlat {
			note = "E♭"
		}
	case degreeFs:
		note = "F♯"
		if m.IsFlat {
			note = "G♭"
		}
	case degreeGs:
		note = "G♯"
		if m.IsFlat {
			note = "A♭"
		}
	case degreeAs:
		note = "A♯"
		if m.IsFlat {
			note = "B♭"
		}
	default:
		panic("unreachable")
	}

	return
}

func (m KeySignature) Text() string {
	if m.IsMajor {
		return m.Note() + " maj."
	}

	return m.Note() + " min."
}

// Taking a signed number of sharps or flats (positive for sharps, negative for flats) and a mode (0 for major, 1 for minor)
// decide the key signature.
func keyFromSharpsOrFlats(sharpsOrFlats int8, mode uint8) uint8 {
	tmp := int(sharpsOrFlats * 7)

	// Relative Minor.
	if mode == minorMode {
		tmp -= 3
	}

	// Clamp to Octave 0-11.
	for tmp < 0 {
		tmp += 12
	}

	return uint8(tmp % 12)
}

func (m KeySignature) readFrom(rd io.Reader) (Message, error) {

	// fmt.Println("Key signature")
	// TODO TEST
	var sharpsOrFlats int8
	var mode uint8

	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 2 {
		err = lib.UnexpectedMessageLengthError("KeySignature expected length 2")
		return nil, err
	}

	// Signed int, positive is sharps, negative is flats.
	var b byte
	b, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	sharpsOrFlats = int8(b)

	// Mode is Major or Minor.
	mode, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	num := sharpsOrFlats
	if num < 0 {
		num = num * (-1)
	}

	key := keyFromSharpsOrFlats(sharpsOrFlats, mode)

	return KeySignature{
		Key:     key,
		Num:     uint8(num),
		IsMajor: mode == majorMode,
		IsFlat:  sharpsOrFlats < 0,
	}, nil

}

func (m KeySignature) meta() {}

// TODO: check if needed, what about Sequencer Specific meta-Event?
type metaUnknown []byte

func (m metaUnknown) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

func (m metaUnknown) Bytes() []byte {
	return []byte(m)
}

func (m metaUnknown) Len() int {
	return len(m)
}
func (m metaUnknown) meta() {}

/*
	http://midi.teragonaudio.com/tech/midifile/port.htm

	   Device (Port) Name

	   FF 09 len text

	   The name of the MIDI device (port) where the track is routed.
	   This replaces the "MIDI Port" meta-Event which some sequencers
	   formally used to route MIDI tracks to various MIDI ports
	   (in order to support more than 16 MIDI channels).

	   For example, assume that you have a MIDI interface that has 4 MIDI output ports.
	   They are listed as "MIDI Out 1", "MIDI Out 2", "MIDI Out 3", and "MIDI Out 4".
	   If you wished a particular MTrk to use "MIDI Out 1" then you would put a
	   Port Name meta-event at the beginning of the MTrk, with "MIDI Out 1" as the text.

	   All MIDI events that occur in the MTrk, after a given Port Name event, will be
	   routed to that port.

	   In a format 0 MIDI file, it would be permissible to have numerous Port Name events
	   intermixed with MIDI events, so that the one MTrk could address numerous ports.
	   But that would likely make the MIDI file much larger than it need be. The Port Name
	   event is useful primarily in format 1 MIDI files, where each MTrk gets routed to
	   one particular port.

	   Note that len could be a series of bytes since it is expressed as a variable length quantity.
*/

type DevicePort string

func (m DevicePort) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m DevicePort) meta() {}

func (m DevicePort) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)
	if err != nil {
		return nil, err
	}

	return DevicePort(text), nil
}

// TODO implement
func (m DevicePort) Raw() []byte {
	panic("not implemented")
}

func (m DevicePort) Text() string {
	return string(m)
}

type Text string

func (m Text) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m Text) meta() {}

func (m Text) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byteText,
		Data: []byte(m),
	}).Bytes()
}

func (m Text) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)
	if err != nil {
		return nil, err
	}

	return Text(text), nil
}

func (m Text) Text() string {
	return string(m)
}

type Copyright string

func (m Copyright) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}
func (m Copyright) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return Copyright(text), nil
}

func (m Copyright) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byteCopyright,
		Data: []byte(m),
	}).Bytes()
}

func (m Copyright) Text() string {
	return string(m)
}

func (m Copyright) meta() {}

type Sequence string

func (m Sequence) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}
func (m Sequence) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return Sequence(text), nil

}

func (m Sequence) Text() string {
	return string(m)
}

func (m Sequence) meta() {}

func (m Sequence) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byteSequence,
		Data: []byte(m),
	}).Bytes()
}

type TrackInstrument string

func (m TrackInstrument) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m TrackInstrument) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byteTrackInstrument,
		Data: []byte(m),
	}).Bytes()
}

func (m TrackInstrument) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return TrackInstrument(text), nil
}

func (m TrackInstrument) Text() string {
	return string(m)
}

func (m TrackInstrument) meta() {}

type Marker string

func (m Marker) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m Marker) Text() string {
	return string(m)
}

func (m Marker) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byte(byteMarker),
		Data: []byte(m),
	}).Bytes()
}

func (m Marker) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return Marker(text), nil
}

func (m Marker) meta() {}

type Lyric string

func (m Lyric) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m Lyric) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return Lyric(text), nil
}

func (m Lyric) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byte(byteLyric),
		Data: []byte(m),
	}).Bytes()
}

func (m Lyric) Text() string {
	return string(m)
}

func (m Lyric) meta() {}

type CuePoint string

func (m CuePoint) Text() string {
	return string(m)
}

func (m CuePoint) Raw() []byte {
	return (&lib.MetaMessage{
		Typ:  byte(byteCuePoint),
		Data: []byte(m),
	}).Bytes()
}

func (m CuePoint) String() string {
	return fmt.Sprintf("%T: %#v", m, string(m))
}

func (m CuePoint) readFrom(rd io.Reader) (Message, error) {
	text, err := lib.ReadText(rd)

	if err != nil {
		return nil, err
	}

	return CuePoint(text), nil
}

func (m CuePoint) meta() {}
