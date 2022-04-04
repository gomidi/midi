package gm

//General MIDI Percussion Key Map

type DrumKey uint8

func (d DrumKey) Key() uint8 {
	return uint8(d) + 1
}

const (
	DrumKey_AcousticBassDrum DrumKey = 34
	DrumKey_BassDrum1        DrumKey = 35
	DrumKey_SideStick        DrumKey = 36
	DrumKey_AcousticSnare    DrumKey = 37
	DrumKey_HandClap         DrumKey = 38
	DrumKey_ElectricSnare    DrumKey = 39
	DrumKey_LowFloorTom      DrumKey = 40
	DrumKey_ClosedHiHat      DrumKey = 41
	DrumKey_HighFloorTom     DrumKey = 42
	DrumKey_PedalHiHat       DrumKey = 43
	DrumKey_LowTom           DrumKey = 44
	DrumKey_OpenHiHat        DrumKey = 45
	DrumKey_LowMidTom        DrumKey = 46
	DrumKey_HiMidTom         DrumKey = 47
	DrumKey_CrashCymbal1     DrumKey = 48
	DrumKey_HighTom          DrumKey = 49
	DrumKey_RideCymbal1      DrumKey = 50
	DrumKey_ChineseCymbal    DrumKey = 51
	DrumKey_RideBell         DrumKey = 52
	DrumKey_Tambourine       DrumKey = 53
	DrumKey_SplashCymbal     DrumKey = 54
	DrumKey_Cowbell          DrumKey = 55
	DrumKey_CrashCymbal2     DrumKey = 56
	DrumKey_Vibraslap        DrumKey = 57
	DrumKey_RideCymbal2      DrumKey = 58
	DrumKey_HiBongo          DrumKey = 59
	DrumKey_LowBongo         DrumKey = 60
	DrumKey_MuteHiConga      DrumKey = 61
	DrumKey_OpenHiConga      DrumKey = 62
	DrumKey_LowConga         DrumKey = 63
	DrumKey_HighTimbale      DrumKey = 64
	DrumKey_LowTimbale       DrumKey = 65
	DrumKey_HighAgogo        DrumKey = 66
	DrumKey_LowAgogo         DrumKey = 67
	DrumKey_Cabasa           DrumKey = 68
	DrumKey_Maracas          DrumKey = 69
	DrumKey_ShortWhistle     DrumKey = 70
	DrumKey_LongWhistle      DrumKey = 71
	DrumKey_ShortGuiro       DrumKey = 72
	DrumKey_LongGuiro        DrumKey = 73
	DrumKey_Claves           DrumKey = 74
	DrumKey_HiWoodBlock      DrumKey = 75
	DrumKey_LowWoodBlock     DrumKey = 76
	DrumKey_MuteCuica        DrumKey = 77
	DrumKey_OpenCuica        DrumKey = 78
	DrumKey_MuteTriangle     DrumKey = 79
	DrumKey_OpenTriangle     DrumKey = 80
)
