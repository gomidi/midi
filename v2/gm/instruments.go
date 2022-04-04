package gm

// GM Instrument Patch Map

/*
These instrument
sounds are grouped into "sets" of related sounds. For example, program numbers 1-8 are piano
sounds, 86 are chromatic percussion sounds, 17-24 are organ sounds, 25-32 are guitar
sounds, etc.
*/

type Instr uint8

func (i Instr) Value() uint8 {
	return uint8(i)
}

func (i Instr) String() string {
	return instrNames[uint8(i)]
}

const (
	Instr_AcousticGrandPiano  Instr = 0
	Instr_BrightAcousticPiano Instr = 1
	Instr_ElectricGrandPiano  Instr = 2
	Instr_HonkytonkPiano      Instr = 3
	Instr_ElectricPiano1      Instr = 4
	Instr_ElectricPiano2      Instr = 5
	Instr_Harpsichord         Instr = 6
	Instr_Clavi               Instr = 7
	Instr_Celesta             Instr = 8
	Instr_Glockenspiel        Instr = 9
	Instr_MusicBox            Instr = 10
	Instr_Vibraphone          Instr = 11
	Instr_Marimba             Instr = 12
	Instr_Xylophone           Instr = 13
	Instr_TubularBells        Instr = 14
	Instr_Dulcimer            Instr = 15
	Instr_DrawbarOrgan        Instr = 16
	Instr_PercussiveOrgan     Instr = 17
	Instr_RockOrgan           Instr = 18
	Instr_ChurchOrgan         Instr = 19
	Instr_ReedOrgan           Instr = 20
	Instr_Accordion           Instr = 21
	Instr_Harmonica           Instr = 22
	Instr_TangoAccordion      Instr = 23
	Instr_AcousticGuitarNylon Instr = 24
	Instr_AcousticGuitarSteel Instr = 25
	Instr_ElectricGuitarJazz  Instr = 26
	Instr_ElectricGuitarClean Instr = 27
	Instr_ElectricGuitarMuted Instr = 28
	Instr_OverdrivenGuitar    Instr = 29
	Instr_DistortionGuitar    Instr = 30
	Instr_Guitarharmonics     Instr = 31
	Instr_AcousticBass        Instr = 32
	Instr_ElectricBassFinger  Instr = 33
	Instr_ElectricBassPick    Instr = 34
	Instr_FretlessBass        Instr = 35
	Instr_SlapBass1           Instr = 36
	Instr_SlapBass2           Instr = 37
	Instr_SynthBass1          Instr = 38
	Instr_SynthBass2          Instr = 39
	Instr_Violin              Instr = 40
	Instr_Viola               Instr = 41
	Instr_Cello               Instr = 42
	Instr_Contrabass          Instr = 43
	Instr_TremoloStrings      Instr = 44
	Instr_PizzicatoStrings    Instr = 45
	Instr_OrchestralHarp      Instr = 46
	Instr_Timpani             Instr = 47
	Instr_StringEnsemble1     Instr = 48
	Instr_StringEnsemble2     Instr = 49
	Instr_SynthStrings1       Instr = 50
	Instr_SynthStrings2       Instr = 51
	Instr_ChoirAahs           Instr = 52
	Instr_VoiceOohs           Instr = 53
	Instr_SynthVoice          Instr = 54
	Instr_OrchestraHit        Instr = 55
	Instr_Trumpet             Instr = 56
	Instr_Trombone            Instr = 57
	Instr_Tuba                Instr = 58
	Instr_MutedTrumpet        Instr = 59
	Instr_FrenchHorn          Instr = 60
	Instr_BrassSection        Instr = 61
	Instr_SynthBrass1         Instr = 62
	Instr_SynthBrass2         Instr = 63
	Instr_SopranoSax          Instr = 64
	Instr_AltoSax             Instr = 65
	Instr_TenorSax            Instr = 66
	Instr_BaritoneSax         Instr = 67
	Instr_Oboe                Instr = 68
	Instr_EnglishHorn         Instr = 69
	Instr_Bassoon             Instr = 70
	Instr_Clarinet            Instr = 71
	Instr_Piccolo             Instr = 72
	Instr_Flute               Instr = 73
	Instr_Recorder            Instr = 74
	Instr_PanFlute            Instr = 75
	Instr_BlownBottle         Instr = 76
	Instr_Shakuhachi          Instr = 77
	Instr_Whistle             Instr = 78
	Instr_Ocarina             Instr = 79
	Instr_Lead1Square         Instr = 80
	Instr_Lead2Sawtooth       Instr = 81
	Instr_Lead3Calliope       Instr = 82
	Instr_Lead4Chiff          Instr = 83
	Instr_Lead5Charang        Instr = 84
	Instr_Lead6Voice          Instr = 85
	Instr_Lead7Fifths         Instr = 86
	Instr_Lead8Basslead       Instr = 87
	Instr_Pad1Newage          Instr = 88
	Instr_Pad2Warm            Instr = 89
	Instr_Pad3Polysynth       Instr = 90
	Instr_Pad4Choir           Instr = 91
	Instr_Pad5Bowed           Instr = 92
	Instr_Pad6Metallic        Instr = 93
	Instr_Pad7Halo            Instr = 94
	Instr_Pad8Sweep           Instr = 95
	Instr_FX1Rain             Instr = 96
	Instr_FX2Soundtrack       Instr = 97
	Instr_FX3Crystal          Instr = 98
	Instr_FX4Atmosphere       Instr = 99
	Instr_FX5Brightness       Instr = 100
	Instr_FX6Goblins          Instr = 101
	Instr_FX7Echoes           Instr = 102
	Instr_FX8Scifi            Instr = 103
	Instr_Sitar               Instr = 104
	Instr_Banjo               Instr = 105
	Instr_Shamisen            Instr = 106
	Instr_Koto                Instr = 107
	Instr_Kalimba             Instr = 108
	Instr_Bagpipe             Instr = 109
	Instr_Fiddle              Instr = 110
	Instr_Shanai              Instr = 111
	Instr_TinkleBell          Instr = 112
	Instr_Agogo               Instr = 113
	Instr_SteelDrums          Instr = 114
	Instr_Woodblock           Instr = 115
	Instr_TaikoDrum           Instr = 116
	Instr_MelodicTom          Instr = 117
	Instr_SynthDrum           Instr = 118
	Instr_ReverseCymbal       Instr = 119
	Instr_GuitarFretNoise     Instr = 120
	Instr_BreathNoise         Instr = 121
	Instr_Seashore            Instr = 122
	Instr_BirdTweet           Instr = 123
	Instr_TelephoneRing       Instr = 124
	Instr_Helicopter          Instr = 125
	Instr_Applause            Instr = 126
	Instr_Gunshot             Instr = 127
)

var instrNames = map[uint8]string{
	0:   "AcousticGrandPiano",
	1:   "BrightAcousticPiano",
	2:   "ElectricGrandPiano",
	3:   "HonkytonkPiano",
	4:   "ElectricPiano1",
	5:   "ElectricPiano2",
	6:   "Harpsichord",
	7:   "Clavi",
	8:   "Celesta",
	9:   "Glockenspiel",
	10:  "MusicBox",
	11:  "Vibraphone",
	12:  "Marimba",
	13:  "Xylophone",
	14:  "TubularBells",
	15:  "Dulcimer",
	16:  "DrawbarOrgan",
	17:  "PercussiveOrgan",
	18:  "RockOrgan",
	19:  "ChurchOrgan",
	20:  "ReedOrgan",
	21:  "Accordion",
	22:  "Harmonica",
	23:  "TangoAccordion",
	24:  "AcousticGuitarNylon",
	25:  "AcousticGuitarSteel",
	26:  "ElectricGuitarJazz",
	27:  "ElectricGuitarClean",
	28:  "ElectricGuitarMuted",
	29:  "OverdrivenGuitar",
	30:  "DistortionGuitar",
	31:  "Guitarharmonics",
	32:  "AcousticBass",
	33:  "ElectricBassFinger",
	34:  "ElectricBassPick",
	35:  "FretlessBass",
	36:  "SlapBass1",
	37:  "SlapBass2",
	38:  "SynthBass1",
	39:  "SynthBass2",
	40:  "Violin",
	41:  "Viola",
	42:  "Cello",
	43:  "Contrabass",
	44:  "TremoloStrings",
	45:  "PizzicatoStrings",
	46:  "OrchestralHarp",
	47:  "Timpani",
	48:  "StringEnsemble1",
	49:  "StringEnsemble2",
	50:  "SynthStrings1",
	51:  "SynthStrings2",
	52:  "ChoirAahs",
	53:  "VoiceOohs",
	54:  "SynthVoice",
	55:  "OrchestraHit",
	56:  "Trumpet",
	57:  "Trombone",
	58:  "Tuba",
	59:  "MutedTrumpet",
	60:  "FrenchHorn",
	61:  "BrassSection",
	62:  "SynthBrass1",
	63:  "SynthBrass2",
	64:  "SopranoSax",
	65:  "AltoSax",
	66:  "TenorSax",
	67:  "BaritoneSax",
	68:  "Oboe",
	69:  "EnglishHorn",
	70:  "Bassoon",
	71:  "Clarinet",
	72:  "Piccolo",
	73:  "Flute",
	74:  "Recorder",
	75:  "PanFlute",
	76:  "BlownBottle",
	77:  "Shakuhachi",
	78:  "Whistle",
	79:  "Ocarina",
	80:  "Lead1Square",
	81:  "Lead2Sawtooth",
	82:  "Lead3Calliope",
	83:  "Lead4Chiff",
	84:  "Lead5Charang",
	85:  "Lead6Voice",
	86:  "Lead7Fifths",
	87:  "Lead8Basslead",
	88:  "Pad1Newage",
	89:  "Pad2Warm",
	90:  "Pad3Polysynth",
	91:  "Pad4Choir",
	92:  "Pad5Bowed",
	93:  "Pad6Metallic",
	94:  "Pad7Halo",
	95:  "Pad8Sweep",
	96:  "FX1Rain",
	97:  "FX2Soundtrack",
	98:  "FX3Crystal",
	99:  "FX4Atmosphere",
	100: "FX5Brightness",
	101: "FX6Goblins",
	102: "FX7Echoes",
	103: "FX8Scifi",
	104: "Sitar",
	105: "Banjo",
	106: "Shamisen",
	107: "Koto",
	108: "Kalimba",
	109: "Bagpipe",
	110: "Fiddle",
	111: "Shanai",
	112: "TinkleBell",
	113: "Agogo",
	114: "SteelDrums",
	115: "Woodblock",
	116: "TaikoDrum",
	117: "MelodicTom",
	118: "SynthDrum",
	119: "ReverseCymbal",
	120: "GuitarFretNoise",
	121: "BreathNoise",
	122: "Seashore",
	123: "BirdTweet",
	124: "TelephoneRing",
	125: "Helicopter",
	126: "Applause",
	127: "Gunshot",
}
