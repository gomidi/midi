package midi

const (
	Off uint8 = 0   // value meaning "off"
	On  uint8 = 127 // value meaning "on"

	BankSelectMSB             uint8 = 0
	ModulationWheelMSB        uint8 = 1
	BreathControllerMSB       uint8 = 2
	FootPedalMSB              uint8 = 4
	PortamentoTimeMSB         uint8 = 5
	DataEntryMSB              uint8 = 6
	VolumeMSB                 uint8 = 7
	BalanceMSB                uint8 = 8
	PanPositionMSB            uint8 = 10
	ExpressionMSB             uint8 = 11
	EffectControl1MSB         uint8 = 12
	EffectControl2MSB         uint8 = 13
	GeneralPurposeSlider1     uint8 = 16
	GeneralPurposeSlider2     uint8 = 17
	GeneralPurposeSlider3     uint8 = 18
	GeneralPurposeSlider4     uint8 = 19
	BankSelectLSB             uint8 = 32
	ModulationWheelLSB        uint8 = 33
	BreathControllerLSB       uint8 = 34
	FootPedalLSB              uint8 = 36
	PortamentoTimeLSB         uint8 = 37
	DataEntryLSB              uint8 = 38
	VolumeLSB                 uint8 = 39
	BalanceLSB                uint8 = 40
	PanPositionLSB            uint8 = 42
	ExpressionLSB             uint8 = 43
	EffectControl1LSB         uint8 = 44
	EffectControl2LSB         uint8 = 45
	SoundVariation            uint8 = 70
	SoundTimbre               uint8 = 71
	SoundReleaseTime          uint8 = 72
	SoundAttackTime           uint8 = 73
	SoundBrightness           uint8 = 74
	SoundControl6             uint8 = 75
	SoundControl7             uint8 = 76
	SoundControl8             uint8 = 77
	SoundControl9             uint8 = 78
	SoundControl10            uint8 = 79
	EffectsLevel              uint8 = 91
	TremuloLevel              uint8 = 92
	ChorusLevel               uint8 = 93
	CelesteLevel              uint8 = 94
	PhaserLevel               uint8 = 95
	DataButtonIncrement       uint8 = 96
	DataButtonDecrement       uint8 = 97
	NonRegisteredParameterLSB uint8 = 98
	NonRegisteredParameterMSB uint8 = 99
	RegisteredParameterLSB    uint8 = 100
	RegisteredParameterMSB    uint8 = 101

	AllSoundOff       uint8 = 120 // send it with value of 0/Off
	AllControllersOff uint8 = 121 // send it with value of 0/Off
	AllNotesOff       uint8 = 123 // send it with value of 0/Off

	OmniModeOff uint8 = 124 // send it with value of 0/Off
	OmniModeOn  uint8 = 125 // send it with value of 0

	MonoOperation uint8 = 126

	PolyOperation uint8 = 127 // send it with value of 0

	LocalKeyboardSwitch  uint8 = 122 // send it with value of 127/On or 0/Off
	HoldPedalSwitch      uint8 = 64  // send it with value of 127/On or 0/Off
	PortamentoSwitch     uint8 = 65  // send it with value of 127/On or 0/Off
	SustenutoPedalSwitch uint8 = 66  // send it with value of 127/On or 0/Off

	SoftPedalSwitch             uint8 = 67 // send it with value of 127/On or 0/Off
	LegatoPedalSwitch           uint8 = 68 // send it with value of 127/On or 0/Off
	Hold2PedalSwitch            uint8 = 69 // send it with value of 127/On or 0/Off
	GeneralPurposeButton1Switch uint8 = 80 // send it with value of 127/On or 0/Off
	GeneralPurposeButton2Switch uint8 = 81 // send it with value of 127/On or 0/Off
	GeneralPurposeButton3Switch uint8 = 82 // send it with value of 127/On or 0/Off
	GeneralPurposeButton4Switch uint8 = 83 // send it with value of 127/On or 0/Off
)

// stolen from http://midi.teragonaudio.com/tech/midispec.htm
var ControlChangeName = map[uint8]string{
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
