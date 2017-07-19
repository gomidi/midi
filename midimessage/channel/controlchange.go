package channel

import (
	"fmt"
	"github.com/gomidi/midi/internal/midilib"
)

// ControlChange represents a MIDI control change message
type ControlChange struct {
	channel    uint8
	controller uint8
	value      uint8
}

// Controller returns the controller of the control change message
func (c ControlChange) Controller() uint8 {
	return c.controller
}

// Value returns the value of the control change message
func (c ControlChange) Value() uint8 {
	return c.value
}

// Channel returns the MIDI channel of the control change message
func (c ControlChange) Channel() uint8 {
	return c.channel
}

// Raw returns the raw bytes of the control change message.
func (c ControlChange) Raw() []byte {
	return channelMessage2(c.channel, 11, c.controller, c.value)
}

// set returns a new control change message that is set to the parsed arguments
func (ControlChange) set(channel uint8, firstArg, secondArg uint8) setter2 {
	var m ControlChange
	m.channel = channel
	m.controller, m.value = midilib.ParseTwoUint7(firstArg, secondArg)
	return m

}

// String returns human readable information about the control change message.
func (c ControlChange) String() string {

	if name, has := ccControllers[c.controller]; has {
		return fmt.Sprintf("%T channel %v controller %v (%#v) value %v", c, c.Channel(), c.Controller(), name, c.Value())
	}
	return fmt.Sprintf("%T channel %v controller %v value %v", c, c.Channel(), c.Controller(), c.Value())

}

// stolen from http://midi.teragonaudio.com/tech/midispec.htm
var ccControllers = map[uint8]string{
	0:   "Bank Select (MSB)",
	1:   "Modulation Wheel (MSB)",
	2:   "Breath controller (MSB)",
	4:   "Foot Pedal (MSB)",
	5:   "Portamento Time (MSB)",
	6:   "Data Entry (MSB)",
	7:   "Volume (MSB)",
	8:   "Balance (MSB)",
	10:  "Pan position (MSB)",
	11:  "Expression (MSB)",
	12:  "Effect Control 1 (MSB)",
	13:  "Effect Control 2 (MSB)",
	16:  "General Purpose Slider 1",
	17:  "General Purpose Slider 2",
	18:  "General Purpose Slider 3",
	19:  "General Purpose Slider 4",
	32:  "Bank Select (LSB)",
	33:  "Modulation Wheel (LSB)",
	34:  "Breath controller (LSB)",
	36:  "Foot Pedal (LSB)",
	37:  "Portamento Time (LSB)",
	38:  "Data Entry (LSB)",
	39:  "Volume (LSB)",
	40:  "Balance (LSB)",
	42:  "Pan position (LSB)",
	43:  "Expression (LSB)",
	44:  "Effect Control 1 (LSB)",
	45:  "Effect Control 2 (LSB)",
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
	98:  "Non-registered Parameter (LSB)",
	99:  "Non-registered Parameter (MSB)",
	100: "Registered Parameter (LSB)",
	101: "Registered Parameter (MSB)",
	120: "All Sound Off",
	121: "All Controllers Off",
	122: "Local Keyboard (on/off)",
	123: "All Notes Off",
	124: "Omni Mode Off",
	125: "Omni Mode On",
	126: "Mono Operation",
	127: "Poly Operation",
}
