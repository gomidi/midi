package smfimage

import (
	"image"

	"image/color"
)

type colorMapper map[Interval]color.RGBA

func (c colorMapper) Map(i Interval) color.RGBA {
	return c[i]
}

/*
RainbowColors maps the following way (interval -> color)

		prime/octave  -> yellow
		minor second  -> mint
		major second  -> orange
		minor third   -> sky blue
		major third   -> spring green
		fourth        -> cyan
		tritone       -> lime/chartreuse
		fifth         -> royalblue
		minor sixth   -> pink
		major sixth   -> violett/purple
		minor seventh -> magenta
		major seventh -> red


*/

/*
linnstrument allows only 10 colors
1=red
2=yellow
3=green
4=cyan
5=blue
6=magenta
7=off
8=white
9= orange
10=lime
11=pink
*/

/*
c    prime/octave  -> yellow           -> yellow
des  minor second  -> mint             -> magenta
d    major second  -> orange           -> pink
es   minor third   -> sky blue         -> oliv
e    major third   -> spring green     -> spring green
f    fourth        -> cyan             -> red
ges  tritone       -> lime/chartreuse  -> cyan
g    fifth         -> royalblue        -> sky blue
as   minor sixth   -> pink             -> burgunder rot-brown
a    major sixth   -> violett/purple   -> orange
b    minor seventh -> magenta          -> royalblue
h    major seventh -> red              -> violett/purple
*/

var (
	yellow      = color.RGBA{255, 255, 0, 255}  // yellow prime/oktave (c)
	magenta     = color.RGBA{255, 0, 255, 255}  // magenta kl. septe (b)
	pink        = color.RGBA{212, 94, 177, 255} // pink kl. sexte (as)
	oliv        = color.RGBA{90, 121, 0, 255}
	springGreen = color.RGBA{109, 205, 2, 255}
	red         = color.RGBA{242, 10, 0, 255}
	cyan        = color.RGBA{0, 245, 255, 255} // cyan/türkis, quarte (f)
	skyBlue     = color.RGBA{35, 144, 232, 255}
	burgund     = color.RGBA{154, 28, 56, 255}
	orange      = color.RGBA{232, 153, 35, 255}
	royalBlue   = color.RGBA{24, 69, 160, 255}
	violett     = color.RGBA{155, 48, 255, 255} // violett/purple gr sexte      (a)

	//orange  = color.RGBA{255, 165, 0, 255}   // orange gr sekunde    (d)
	//skyBlue = color.RGBA{135, 206, 255, 255} // sky blue kl terz (es)
	//springGreen = color.RGBA{0, 205, 102, 255}   // spring green große terz    (e)
	mint = color.RGBA{189, 252, 201, 255} // mint kl. sekunde (des)
	lime = color.RGBA{127, 255, 0, 255}   // lime/chartreuse tritonus (ges)
	// royalBlue   = color.RGBA{72, 118, 255, 255}  // royalblue quinte       (g)
	//red         = color.RGBA{255, 0, 0, 255}     // red gr sept       (h)
)

/*

1 des                                    magenta
5 f         - subdominante               strahlendes red
8 as                                     burgunder rot-brown
9 a         -                            orange

0 c         - Tonika start subdom end    yellow
3 es                                     oliv
4 e                                      spring green
6 ges                                    cyan

7 g         - Tonika end Dominant Start  sky blue
10 b                                      royalblue
11 h         -                            violett/purple
2 d         - Dominant end               pink
ccccccccccc
*/

var mapNoteToPos = map[int]int{
	1:  11,
	5:  10,
	8:  9,
	9:  8,
	0:  7,
	3:  6,
	4:  5,
	6:  4,
	7:  3,
	10: 2,
	11: 1,
	2:  0,
}

/*
C  gelb
D  grün
E  türkis
F  blau
G  violett
A  rot
H  orange
C  gelb
*/

var (

	/*
			TODO: round trip
			- FFT analysis von Datei -> MIDI-Notes(1) for sinus tones incl. pitchbend velocity / volume -> wav gen via supercollider -> FTT-Analysis -> MIDI-Notes(2)
		    - MIDI-Notes(1) must be == MIDI-Notes(2)

			Darstellung:
			  - Pitchbend berücksichtigen: Voreinstellung +/- 2 Semitones: =+/-8192
			  - d.h. bei vierteltönen: +/- 2048 und bei achteltönen +/- 1024
			  - minium velocity/volume (noten mit kleinerer velocity/volume werden nicht angezeigt
	*/

	ColorYellow          = color.RGBA{255, 255, 0, 255} // yellow
	ColorElectricLime    = color.RGBA{201, 255, 0, 255} // yellow-green
	ColorSpringBud       = color.RGBA{148, 255, 0, 255} // light green
	ColorHarlequin       = color.RGBA{81, 255, 0, 255}  // green
	ColorMalachite       = color.RGBA{0, 232, 83, 255}  // dark green
	ColorCaribbeanGreen  = color.RGBA{0, 206, 133, 255} // turkis
	ColorRobinsEggBlue   = color.RGBA{0, 196, 190, 255} // blue-green
	ColorDeepSkyBlue     = color.RGBA{0, 173, 236, 255} // sky blue
	ColorDodgerBlue      = color.RGBA{0, 130, 255, 255} // light blue
	ColorNavyBlue        = color.RGBA{0, 83, 255, 255}  // deep blue
	ColorBlue            = color.RGBA{0, 57, 255, 255}  // blue
	ColorDarkBlue        = color.RGBA{0, 0, 255, 255}   // dark blue
	ColorIndigoBlue      = color.RGBA{90, 0, 252, 255}  // indigo blue
	ColorElectricIndigo  = color.RGBA{105, 0, 210, 255} // indigo
	ColorDarkViolet      = color.RGBA{153, 0, 202, 255} // violet
	ColorDeepMagenta     = color.RGBA{201, 0, 195, 255} // magenta
	ColorRazzmatazz      = color.RGBA{247, 0, 117, 255} // pink
	ColorRed             = color.RGBA{255, 0, 0, 255}   // red
	ColorOrangeRed       = color.RGBA{255, 60, 0, 255}  // orange red
	ColorSafetyOrange    = color.RGBA{255, 106, 0, 255} // orange
	ColorDarkOrange      = color.RGBA{255, 149, 0, 255} // dark orange
	ColorSelectiveYellow = color.RGBA{255, 176, 0, 255} // yellow-orange
	ColorTangerineYellow = color.RGBA{255, 202, 0, 255} // dark yellow
	ColorGoldenYellow    = color.RGBA{255, 227, 0, 255} // deep yellow

	//                                                                             names according to https://www.htmlcsscolor.com
	PrimeColor = ColorYellow //                               C  hellgelb 255 255   0   ffff00 Yellow
	//                                                               201 255   0   c9ff00 Electric Lime
	MinSecondColor = ColorSpringBud //                          C# hellgrün 148 255   0   94ff00 Spring Bud
	//                                                                81 255   0   51ff00 Harlequin
	MajSecondColor = ColorMalachite //                             D  dunkelgrün 0 232  83   00e853 Malachite
	//                                                                 0 206 133   00ce85 Caribbean Green
	MinThirdColor = ColorRobinsEggBlue // oliv                         D#  türkis    0 196 190   00c4be Robin's Egg Blue
	//                                                                 0 173 236   00adec Deep Sky Blue
	MajThirdColor = ColorDodgerBlue // türkis                E hellblau    0 130 255   0082ff Dodger Blue
	//                                                                 0  83 255   0053ff Navy Blue
	FourthColor = ColorBlue // activ red                       F dunkelblau  0  57 255   0039ff Blue
	//                                                                 0   0 255   0000ff Blue
	TritoneColor = ColorIndigoBlue //                               F# violet    90   0 252   5a00fc Electric Indigo
	//                                                               105   0 210   6900d2 Electric Indigo
	FifthColor = ColorDarkViolet // activ blue                   G  purple   153   0 202   9900ca Dark Violet
	//                                                               201   0 195   c900c3 Deep Magenta
	MinSixthColor = ColorRazzmatazz // burgunder rot-brown       G#  pink    247   0 117   f70075 Razzmatazz
	//                                                               255   0   0   ff0000 Red
	MajSixthColor = ColorOrangeRed //                            A      rot  255  60   0   ff3c00 Orange Red
	//                                                               255 106   0   ff6a00 Safety Orange
	MinSeventhColor = ColorDarkOrange // dark/deep blue        A#  orange  255 149   0   ff9500 Dark Orange
	//                                                               255 176   0   ffb000 Selective Yellow
	MajSeventhColor = ColorTangerineYellow //                         H dunkelgelb255 202   0   ffca00 Tangerine Yellow
	//                                                               255 227   0   ffe300 golden yellow
	//                                                   C hellgelb  255 255   0   ffff00 yellow
	/*
		PrimeColor      = color.RGBA{255, 255, 0, 255}   // yellow prime/oktave (c)
		MinSecondColor  = color.RGBA{189, 252, 201, 255} // mint kl. sekunde (des)
		MajSecondColor  = color.RGBA{255, 165, 0, 255}   // orange gr sekunde    (d)
		MinThirdColor   = color.RGBA{135, 206, 255, 255} // sky blue kl terz (es)
		MajThirdColor   = color.RGBA{0, 205, 102, 255}   // spring green große terz    (e)
		FourthColor     = color.RGBA{0, 245, 255, 255}   // cyan/türkis, quarte (f)
		TritoneColor    = color.RGBA{127, 255, 0, 255}   // lime/chartreuse tritonus (ges)
		FifthColor      = color.RGBA{72, 118, 255, 255}  // royalblue quinte       (g)
		MinSixthColor   = color.RGBA{255, 192, 203, 255} // pink kl. sexte (as)
		MajSixthColor   = color.RGBA{155, 48, 255, 255}  // violett/purple gr sexte      (a)
		MinSeventhColor = color.RGBA{255, 0, 255, 255}   // magenta kl. septe (b)
		MajSeventhColor = color.RGBA{255, 0, 0, 255}     // red gr sept       (h)
	*/
	LightGrey = color.RGBA{187, 187, 187, 255}
	DarkGrey  = color.RGBA{120, 120, 120, 255}
)

type palette struct {
	color.Palette
}

var ColorPalette = &palette{color.Palette{
	color.Black,
	color.White,
	color.Transparent,
	PrimeColor,
	MinSecondColor,
	MajSecondColor,
	MinThirdColor,
	MajThirdColor,
	FourthColor,
	TritoneColor,
	FifthColor,
	MinSixthColor,
	MajSixthColor,
	MinSeventhColor,
	MajSeventhColor,
	LightGrey,
	DarkGrey,
}}

func (p *palette) Quantize(pl color.Palette, im image.Image) color.Palette {
	return pl
}

var RainbowColors colorMapper = map[Interval]color.RGBA{
	Prime:      PrimeColor,      // yellow prime/oktave (c)
	MinSecond:  MinSecondColor,  // mint kl. sekunde (des)
	MajSecond:  MajSecondColor,  // orange gr sekunde    (d)
	MinThird:   MinThirdColor,   // sky blue kl terz (es)
	MajThird:   MajThirdColor,   // spring green große terz    (e)
	Fourth:     FourthColor,     // cyan/türkis, quarte (f)
	Tritone:    TritoneColor,    // lime/chartreuse tritonus (ges)
	Fifth:      FifthColor,      // royalblue quinte       (g)
	MinSixth:   MinSixthColor,   // pink kl. sexte (as)
	MajSixth:   MajSixthColor,   // violett/purple gr sexte      (a)
	MinSeventh: MinSeventhColor, // magenta kl. septe (b)
	MajSeventh: MajSeventhColor, // red gr sept       (h)
}

// ColorMapper maps an interval (counted from basenote) to a color
type ColorMapper interface {
	Map(i Interval) color.RGBA
}
