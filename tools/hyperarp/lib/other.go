package hyperarp

type note uint8

const (
	C   note = 0
	Cis      = iota
	D
	Dis
	E
	F
	Fis
	G
	Gis
	A
	Ais
	B
)

var noteDistanceMap = map[uint8]float64{
	0:  2.0,               // half notes
	1:  1.0,               // quarter notes
	2:  0.5,               // eighths
	3:  0.25,              // sixteenths
	4:  0.125,             // 32ths
	5:  2.0 / 3.0,         // half note tripplets
	6:  1.0 / 3.0,         // quarter note tripplets
	7:  0.5 / 3.0,         // eighths tripples
	8:  0.25 / 3.0,        // sixteenths tripples
	9:  0.125 / 3.0,       // 32ths tripples
	10: 1.0 * 3.0 / 2.0,   // dotted quarter notes
	11: 0.5 * 3.0 / 2.0,   // dotted eighths
	12: 0.25 * 3.0 / 2.0,  // dotted sixteenths
	13: 0.125 * 3.0 / 2.0, // dotted 32ths
	14: 1.0 / 5.0,         // quarter note quintuplets
	15: 0.5 / 5.0,         // eighths quintuplets
}
