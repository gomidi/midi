package smfimage

import (
	"fmt"
	"image/color"
	"io"

	"github.com/fogleman/gg"
	"gitlab.com/gomidi/midi/tools/smfimage/curve"
	"gitlab.com/gomidi/midi/v2/smf"
)

type smfCurve struct {
	*curve.Curve
	*smfimage
	Bar       int
	Left      float64
	TS        [2]uint8
	BarLength float64
	Radius    int
	LineWidth int
}

func (c *smfCurve) resetBar() {
	c.Bar = 0
	c.Left = 0
	c.BarLength = 0
}

func (c *smfCurve) make() {
	// r := 180
	//s.curve = curve.New(curve.MaxHeight(1000), curve.MaxWidth(s.numLast32th*32*r*4), curve.Radius(float64(r)))
	// 4*radius - (2*float64(s.numNotes-1) * s.lineWidth)
	// padding := 10
	//width := ((s.numLast32th/32 + 1) * (4*int(s.curveRadius) - (2 * 11.0 * s.curveLineWidth))) + 2*padding
	// width := 1000 + 2*padding
	c.Curve = curve.New(
		//curve.MaxHeight(4*int(s.curveRadius)),
		//curve.MaxHeight(4*int(s.curveRadius)+s.totalHeight()),
		curve.MaxHeight(c.totalHeight()),
		//curve.MaxHeight(500),
		//curve.MaxWidth(width),
		curve.MaxWidth(c.totalWidth()),
		curve.Radius(float64(c.Radius)),
		curve.LineWidth(float64(c.LineWidth)),
		curve.MiddleLine(2*int(c.Radius)),
		// curve.MiddleLine(100),
	)
}

func (s *smfCurve) drawHarmonicBG(nt int, cl color.Color) {
	start := float64(s.numFirst32th * s.noteWidth)
	sineHeight := float64(s.LineWidth) + float64(s.LineWidth)*12 // + (11 * float64(s.curveLineWidth))
	//padding := 10
	//width := float64(((s.numLast32th/32 + 1) * (4*int(s.curveRadius) - (2 * 11.0 * s.curveLineWidth))) + 2*padding)

	// *float64(12-nt))
	height := sineHeight * 4
	if nt == 8 {
		height = sineHeight * 3
	}
	top := roundFloat(15+sineHeight*float64(nt), 0)

	s.SetColor(cl)
	// fmt.Printf("start: %0.2f top: %0.2f width: %0.2f height: %0.2f\n", start, top, width, height)
	s.DrawRectangle(start, top, float64(s.totalWidth()), height)
	s.Fill()
	// int(roundFloat(5+float64(c.curveLineWidth)+float64(c.curveLineWidth)*12+(11*float64(c.curveLineWidth)*float64(12-nt)), 0))
	// four notes
}

func (s *smfCurve) drawHarmonic() {
	// four notes
	s.drawHarmonicBG(0, LightGrey)
	s.drawHarmonicBG(4, color.Black)
	s.drawHarmonicBG(8, LightGrey)
}

func (s *smfCurve) savePNG(outFile string) (err error) {
	return gg.SavePNG(outFile, s.Image())
}

func (s *smfCurve) savePNGto(wr io.Writer) (err error) {
	return fmt.Errorf("savePNGto is not implemented for smfcurves")
}

func (s *smfCurve) PreHeight() int {
	return int(roundFloat(5+float64(s.LineWidth)+float64(s.LineWidth)*12+(11*float64(s.LineWidth)*float64(12)), 0)) + 5
}

func (s *smfCurve) draw(_32th, track int, nt int, cl color.RGBA) {
	start := float64(s.numFirst32th * s.noteWidth)
	// if track != 3 {
	// return
	// }
	num, denom, _ := s.getTimeSignature(_32th)
	factor := float64(num) / float64(denom)

	diff := float64(_32th) - float64(s.curve.Bar)*32.0*factor
	if diff > 31 {
		// fmt.Println("####")
		s.curve.Bar++
		if !s.singleBar {
			s.curve.Left += (s.curve.BarLength)
			diff = float64(_32th) - float64(s.curve.Bar)*32.0*factor
		} else {
			// s.curveLeft += float64(s.noteHeight) + 2.0
			diff = float64(_32th) - float64(s.curve.Bar)*32.0*factor
		}
	}

	// s.getTimeSignature(col)

	/*
		if s.curveBar > 1 {
			return
		}
	*/
	// r, g, b, _ := cl.RGBA()
	// fmt.Printf("@%v: #%v T%v color: %v %v %v\n", _32th, nt, track, r, g, b)

	// fmt.Printf("bar: %v 32th: %v diff: %0.2f\n", s.curveBar, _32th, float64(diff)/32.0)
	if !s.singleBar {
		s.curve.BarLength = s.curve.AddNote(cl, nt, s.curve.TS[0], s.curve.TS[1], start+s.curve.Left, diff/32.0, 1.0/32.0)
	} else {
		s.curve.BarLength = s.curve.AddNote(cl, nt, s.curve.TS[0], s.curve.TS[1], start+s.curve.Left+float64(s.curve.Bar%10), diff/64.0, 1.0/64.0)
	}

}

func (s *smfCurve) drawBar(x, y int) {
	drawBarOnImage(s.image.barLinesImg, x, y)
}

func (s *smfCurve) drawBars() {
	drawBars(s.smfimage)
	s.DrawImage(s.smfimage.image.barLinesImg, 0, 0)
}

func (s *smfCurve) makeMonoChromeImage(xMax, yMax int) {
	makeMonoChromeImage(s.smfimage, xMax, yMax)
	s.smfimage.image.makeBarLinesImg()
}

func (s *smfCurve) makePaletteImage(xMax, yMax int) {
	makePaletteImage(s.smfimage, xMax, yMax)
	s.smfimage.image.makeBarLinesImg()
}

func (s *smfCurve) makeHarmonic() {
	s.make()
	s.drawHarmonic()
}

func (s *smfCurve) mkSongPicture(rd *smf.SMF, names []string) error {
	err := mkSongPicture(s.smfimage, rd, names)
	if err != nil {
		return err
	}
	s.DrawImage(s.smfimage.image.img, 0, 0)
	return nil
}
