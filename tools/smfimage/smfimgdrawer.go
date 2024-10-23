package smfimage

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"os"

	"gitlab.com/gomidi/midi/v2/smf"
)

type smfimgDrawer struct {
	*smfimage
}

var _ drawer = &smfimgDrawer{}

func (s *smfimgDrawer) savePNG(outFile string) (err error) {
	f, err2 := os.Create(outFile)

	if err2 != nil {
		return err2
	}

	defer f.Close()

	return s.savePNGto(f)
}

func (s *smfimgDrawer) savePNGto(wr io.Writer) (err error) {
	var e png.Encoder

	e.CompressionLevel = png.BestCompression

	return e.Encode(wr, s.image.img)
}

var _ drawer = &smfCurve{}

func (s *smfimgDrawer) Image() image.Image {
	return s.image.img
}

func (s *smfimgDrawer) PreHeight() int {
	if s.wantOverview {
		return (12 * s.noteHeight) + s.trackBorder
	}
	return 0
}

func (s *smfimgDrawer) draw(_32th, track, nnote int, col color.RGBA) {
	s.smfimage.draw((_32th * s.noteWidth), (11-nnote)*s.noteHeight, s.noteWidth, s.noteHeight, col)
}

func (s *smfimgDrawer) drawBar(x, y int) {
	drawBarOnImage(s.image.img, x, y)
}

func (s *smfimgDrawer) drawBars() {
	drawBars(s.smfimage)
}

func (s *smfimgDrawer) makeMonoChromeImage(xMax, yMax int) {
	makeMonoChromeImage(s.smfimage, xMax, yMax)
}

func (s *smfimgDrawer) makePaletteImage(xMax, yMax int) {
	makePaletteImage(s.smfimage, xMax, yMax)
}

func (s *smfimgDrawer) makeHarmonic() {
	s.compactHarmonic()
}

func (c *smfimgDrawer) mkSongPicture(rd *smf.SMF, names []string) error {
	return mkSongPicture(c.smfimage, rd, names)
}
