package cc

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
	LegatoPedalSwitch           uint8 = 67 // send it with value of 127/On or 0/Off
	Hold2PedalSwitch            uint8 = 69 // send it with value of 127/On or 0/Off
	GeneralPurposeButton1Switch uint8 = 80 // send it with value of 127/On or 0/Off
	GeneralPurposeButton2Switch uint8 = 81 // send it with value of 127/On or 0/Off
	GeneralPurposeButton3Switch uint8 = 82 // send it with value of 127/On or 0/Off
	GeneralPurposeButton4Switch uint8 = 83 // send it with value of 127/On or 0/Off
)
