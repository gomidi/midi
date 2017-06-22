package channel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"lib"
)

func (c Channel) Dispatch2(mType uint8, firstArg uint8, rd io.Reader) (ev Event, err error) {
	var secondArg byte
	secondArg, err = lib.ReadByte(rd)

	if err != nil {
		return
	}

	var evt channelSetter2

	switch mType {
	case lib.CodeNoteOff:
		evt = NoteOff{}
	case lib.CodeNoteOn:
		evt = NoteOn{}
	case lib.CodePolyphonicKeyPressure:
		evt = PolyphonicAfterTouch{}
	case lib.CodeControlChange:
		evt = ControlChange{}
	case lib.CodePitchWheel:
		evt = PitchWheel{}
	default:
		// unsupported
		return nil, nil
	}

	evt = evt.set(c, firstArg, secondArg)

	// handle noteOn events with velocity of 0 as note offs
	if noteOn, is := ev.(NoteOn); is && noteOn.velocity == 0 {
		evt = NoteOff{}
		evt = evt.set(c, firstArg, 0)
	}
	return evt, nil
}

func (c Channel) Dispatch1(mType uint8, firstArg uint8) (ev Event) {
	var evt channelSetter1

	switch mType {
	case lib.CodeProgramChange:
		evt = ProgramChange{}
	case lib.CodeChannelPressure:
		evt = AfterTouch{}
	default:
		// unsupported
		return nil
	}

	evt = evt.set(c, firstArg)
	return evt
}

func channelMessage1(c uint8, status, msg byte) []byte {
	cm := &channelMessage{channel: c, status: status}
	cm.data[0] = msg
	return cm.bytes()
}

func channelMessage2(c uint8, status, msg1 byte, msg2 byte) []byte {
	cm := &channelMessage{channel: c, status: status}
	cm.data[0] = msg1
	cm.data[1] = msg2
	cm.twoBytes = true
	return cm.bytes()
}

type channelMessage struct {
	status   uint8
	channel  uint8
	twoBytes bool
	data     [2]byte
}

func (m *channelMessage) getCompleteStatus() uint8 {
	s := m.status << 4
	lib.ClearBitU8(s, 0)
	lib.ClearBitU8(s, 1)
	lib.ClearBitU8(s, 2)
	lib.ClearBitU8(s, 3)
	s = s | m.channel
	return s
}

func (m *channelMessage) bytes() []byte {
	var bf bytes.Buffer
	binary.Write(&bf, binary.BigEndian, m.getCompleteStatus())
	bf.WriteByte(m.data[0])
	if m.twoBytes {
		bf.WriteByte(m.data[1])
	}
	// var b []byte
	// b = append(b, m.getCompleteStatus())
	// b = append(b, m.data[0])
	// b = append(b, m.data[1])
	return bf.Bytes()
}

var (
	_ Event = NoteOff{}
	_ Event = NoteOn{}
	_ Event = PolyphonicAfterTouch{}
	_ Event = ControlChange{}
	_ Event = ProgramChange{}
	_ Event = AfterTouch{}
	_ Event = PitchWheel{}

	_ channelSetter2 = NoteOff{}
	_ channelSetter2 = NoteOn{}
	_ channelSetter2 = PolyphonicAfterTouch{}
	_ channelSetter2 = ControlChange{}
	_ channelSetter2 = PitchWheel{}

	_ channelSetter1 = ProgramChange{}
	_ channelSetter1 = AfterTouch{}
)

type Channel struct {
	number uint8
}

func (c Channel) Number() uint8 {
	return c.number
}

// no must be 0-15
func New(no uint8) Channel {
	if no > 15 {
		panic("midi channels > 15 are not allowed")
	}

	return Channel{no}
}

var (
	Ch0  = New(0)
	Ch1  = New(1)
	Ch2  = New(2)
	Ch3  = New(3)
	Ch4  = New(4)
	Ch5  = New(5)
	Ch6  = New(6)
	Ch7  = New(7)
	Ch8  = New(8)
	Ch9  = New(9)
	Ch10 = New(10)
	Ch11 = New(11)
	Ch12 = New(12)
	Ch13 = New(13)
	Ch14 = New(14)
	Ch15 = New(15)
)

func (c Channel) NoteOff(pitch uint8) NoteOff {
	return NoteOff{channel: c.Number(), pitch: pitch}
}

type NoteOff struct {
	channel uint8
	pitch   uint8
}

func (n NoteOff) Pitch() uint8 {
	return n.pitch
}

func (NoteOff) set(ch Channel, firstArg, secondArg uint8) channelSetter2 {
	var m NoteOff
	m.channel = ch.Number()
	m.pitch, _ = lib.ParseTwoUint7(firstArg, secondArg)
	return m
}

func (n NoteOff) Raw() []byte {
	return channelMessage2(n.channel, 8, n.pitch, 0)
}

func (n NoteOff) Channel() uint8 {
	return n.channel
}

func (m NoteOff) String() string {
	return fmt.Sprintf("%T channel %v pitch %v", m, m.channel, m.pitch)
}

func (c Channel) NoteOn(pitch uint8, veloctiy uint8) NoteOn {
	return NoteOn{channel: c.Number(), pitch: pitch, velocity: veloctiy}
}

type NoteOn struct {
	channel  uint8
	pitch    uint8
	velocity uint8
}

func (n NoteOn) Pitch() uint8 {
	return n.pitch
}

func (n NoteOn) Velocity() uint8 {
	return n.velocity
}

func (n NoteOn) Channel() uint8 {
	return n.channel
}

func (n NoteOn) Raw() []byte {
	return channelMessage2(n.channel, 9, n.pitch, n.velocity)
}

func (NoteOn) set(ch Channel, firstArg, secondArg uint8) channelSetter2 {
	var m NoteOn
	m.channel = ch.Number()
	m.pitch, m.velocity = lib.ParseTwoUint7(firstArg, secondArg)
	return m
}

func (n NoteOn) String() string {
	return fmt.Sprintf("%T channel %v pitch %v vel %v", n, n.channel, n.pitch, n.velocity)
}

type PolyphonicAfterTouch struct {
	channel  uint8
	pitch    uint8
	pressure uint8
}

func (p PolyphonicAfterTouch) Pitch() uint8 {
	return p.pitch
}

func (p PolyphonicAfterTouch) Pressure() uint8 {
	return p.pressure
}

func (c Channel) PolyphonicAfterTouch(pitch uint8, pressure uint8) PolyphonicAfterTouch {
	return PolyphonicAfterTouch{channel: c.Number(), pitch: pitch, pressure: pressure}
}

func (p PolyphonicAfterTouch) Channel() uint8 {
	return p.channel
}

func (p PolyphonicAfterTouch) Raw() []byte {
	return channelMessage2(p.channel, 10, p.pitch, p.pressure)
}

func (PolyphonicAfterTouch) set(ch Channel, firstArg, secondArg uint8) channelSetter2 {
	var m PolyphonicAfterTouch
	m.channel = ch.Number()
	m.pitch, m.pressure = lib.ParseTwoUint7(firstArg, secondArg)
	return m
}

func (p PolyphonicAfterTouch) String() string {
	return fmt.Sprintf("%T channel %v pitch %v pressure %v", p, p.channel, p.pitch, p.pressure)
}

type ControlChange struct {
	channel    uint8
	controller uint8
	value      uint8
}

func (c ControlChange) Controller() uint8 {
	return c.controller
}

func (c ControlChange) Value() uint8 {
	return c.value
}

func (c Channel) ControlChange(controller uint8, value uint8) ControlChange {
	return ControlChange{channel: c.Number(), controller: controller, value: value}
}

func (c ControlChange) Channel() uint8 {
	return c.channel
}

func (c ControlChange) Raw() []byte {
	return channelMessage2(c.channel, 11, c.controller, c.value)
}

func (ControlChange) set(ch Channel, firstArg, secondArg uint8) channelSetter2 {
	var m ControlChange
	m.channel = ch.Number()
	m.controller, m.value = lib.ParseTwoUint7(firstArg, secondArg)
	// TODO split this into ChannelMode for values [120, 127]?
	// TODO implement separate callbacks for each type of:
	// - All sound off
	// - Reset all controllers
	// - Local control
	// - All notes off
	// Only if required. http://www.midi.org/techspecs/midimessages.php
	return m

}

func (c ControlChange) String() string {
	name, has := cCControllers[c.controller]
	if has {
		return fmt.Sprintf("%T channel %v controller %v (%#v) value %v", c, c.channel, c.controller, name, c.value)
	} else {
		return fmt.Sprintf("%T channel %v controller %v value %v", c, c.channel, c.controller, c.value)
	}
}

// stolen from http://midi.teragonaudio.com/tech/midispec.htm
var cCControllers = map[uint8]string{
	0:   "Bank Select (coarse)",
	1:   "Modulation Wheel (coarse)",
	2:   "Breath controller (coarse)",
	4:   "Foot Pedal (coarse)",
	5:   "Portamento Time (coarse)",
	6:   "Data Entry (coarse)",
	7:   "Volume (coarse)",
	8:   "Balance (coarse)",
	10:  "Pan position (coarse)",
	11:  "Expression (coarse)",
	12:  "Effect Control 1 (coarse)",
	13:  "Effect Control 2 (coarse)",
	16:  "General Purpose Slider 1",
	17:  "General Purpose Slider 2",
	18:  "General Purpose Slider 3",
	19:  "General Purpose Slider 4",
	32:  "Bank Select (fine)",
	33:  "Modulation Wheel (fine)",
	34:  "Breath controller (fine)",
	36:  "Foot Pedal (fine)",
	37:  "Portamento Time (fine)",
	38:  "Data Entry (fine)",
	39:  "Volume (fine)",
	40:  "Balance (fine)",
	42:  "Pan position (fine)",
	43:  "Expression (fine)",
	44:  "Effect Control 1 (fine)",
	45:  "Effect Control 2 (fine)",
	64:  "Hold Pedal (on/off)",
	65:  "Portamento (on/off)",
	66:  "Sustenuto Pedal (on/off)",
	67:  "Soft Pedal (on/off)",
	68:  "Legato Pedal (on/off)",
	69:  "Hold 2 Pedal (on/off)",
	70:  "Sound Variation",
	71:  "Sound Timbre",
	72:  "Sound Release Time",
	73:  "Sound Attack Time",
	74:  "Sound Brightness",
	75:  "Sound Control 6",
	76:  "Sound Control 7",
	77:  "Sound Control 8",
	78:  "Sound Control 9",
	79:  "Sound Control 10",
	80:  "General Purpose Button 1 (on/off)",
	81:  "General Purpose Button 2 (on/off)",
	82:  "General Purpose Button 3 (on/off)",
	83:  "General Purpose Button 4 (on/off)",
	91:  "Effects Level",
	92:  "Tremulo Level",
	93:  "Chorus Level",
	94:  "Celeste Level",
	95:  "Phaser Level",
	96:  "Data Button increment",
	97:  "Data Button decrement",
	98:  "Non-registered Parameter (fine)",
	99:  "Non-registered Parameter (coarse)",
	100: "Registered Parameter (fine)",
	101: "Registered Parameter (coarse)",
	120: "All Sound Off",
	121: "All Controllers Off",
	122: "Local Keyboard (on/off)",
	123: "All Notes Off",
	124: "Omni Mode Off",
	125: "Omni Mode On",
	126: "Mono Operation",
	127: "Poly Operation",
}

type Event interface {
	// event.Event
	String() string
	Raw() []byte
	Channel() uint8
}

type channelSetter1 interface {
	Event
	set(ch Channel, firstArg uint8) channelSetter1
}

type channelSetter2 interface {
	Event
	set(ch Channel, firstArg, secondArg uint8) channelSetter2
}

type ProgramChange struct {
	channel uint8
	program uint8
}

func (p ProgramChange) Program() uint8 {
	return p.program
}

func (m ProgramChange) Channel() uint8 {
	return m.channel
}

func (c Channel) ProgramChange(program uint8) ProgramChange {
	return ProgramChange{channel: c.Number(), program: program}
}

func (c ProgramChange) Raw() []byte {
	return channelMessage1(c.channel, 12, c.program)
}

func (c ProgramChange) String() string {
	return fmt.Sprintf("%T channel %v program %v", c, c.channel, c.program)
}

func (ProgramChange) set(ch Channel, firstArg uint8) channelSetter1 {
	var m ProgramChange
	m.channel = ch.Number()
	m.program = lib.ParseUint7(firstArg)
	return m
}

type AfterTouch struct {
	channel  uint8
	pressure uint8
}

func (a AfterTouch) Pressure() uint8 {
	return a.pressure
}

func (a AfterTouch) Channel() uint8 {
	return a.channel
}

func (c Channel) AfterTouch(value uint8) AfterTouch {
	return AfterTouch{channel: c.Number(), pressure: value}
}

func (a AfterTouch) Raw() []byte {
	return channelMessage1(a.channel, 13, a.pressure)
}

func (a AfterTouch) String() string {
	return fmt.Sprintf("%T channel %v value %v", a, a.channel, a.pressure)
}

func (AfterTouch) set(ch Channel, firstArg uint8) channelSetter1 {
	var m AfterTouch
	m.channel = ch.Number()
	m.pressure = lib.ParseUint7(firstArg)
	return m
}

type PitchWheel struct {
	channel  uint8
	value    int16
	absValue uint16
}

func (p PitchWheel) Value() int16 {
	return p.value
}

func (p PitchWheel) AbsValue() uint16 {
	return p.absValue
}

func (c Channel) PitchWheel(value int16) PitchWheel {
	return PitchWheel{channel: c.Number(), value: value}
}

func (p PitchWheel) Channel() uint8 {
	return p.channel
}

func (p PitchWheel) Raw() []byte {
	r := lib.MsbLsbSigned(p.value)

	var bf bytes.Buffer
	//	binary.Write(&bf, binary.BigEndian, uint16(change))
	binary.Write(&bf, binary.BigEndian, r)
	b := bf.Bytes()
	return channelMessage2(p.channel, 14, b[0], b[1])
}

func (p PitchWheel) String() string {
	return fmt.Sprintf("%T channel %v value %v absvalue %v", p, p.channel, p.value, p.absValue)
}

func (PitchWheel) set(ch Channel, firstArg, secondArg uint8) channelSetter2 {
	var m PitchWheel
	m.channel = ch.Number()
	// The value is a signed int (relative to centre), and absoluteValue is the actual value in the file.
	m.value, m.absValue = lib.ParsePitchWheelVals(firstArg, secondArg)
	return m
}
