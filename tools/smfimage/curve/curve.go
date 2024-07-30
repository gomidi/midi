package curve

import (
	// "fmt"
	"github.com/fogleman/gg"
	// "image"
	"image/color"
)

type Curve struct {
	maxX        int
	maxY        int
	paddingLeft float64
	paddingTop  float64
	radius      float64
	lineWidth   float64
	middleLine  float64
	numNotes    int
	*gg.Context
}

func New(options ...Option) *Curve {
	c := &Curve{}
	c.maxX = 1500
	c.maxY = 1500
	c.paddingLeft = 10.0
	c.paddingTop = 10.0
	c.radius = 220.0
	c.lineWidth = 4
	c.middleLine = 500.0
	c.numNotes = 12

	for _, opt := range options {
		opt(c)
	}

	c.Context = gg.NewContext(c.maxX, c.maxY)
	c.Context.SetLineWidth(c.lineWidth)
	return c
}

/*
func (c *Curve) Image() image.Image {
	return c.Context.Image()
}
*/

/*
func (c *Curve) SavePNG(path string) error {
	return gg.SavePNG(path, c.Context.Image())
}
*/

func (c *Curve) barLength(radius float64, num, denom uint8) float64 {
	//return 4*radius - 2*c.linesWidth()
	return 4 * radius
}

func (c *Curve) calcRadiusFromBar(num, denom uint8) float64 {

	// 1bar = 4*c.radius - 2*c.linesWidth()

	// factor * bar = factor * (4 * c.radius - 2*c.linesWidth())

	// factor * (4 * c.radius - 2*c.linesWidth()) = 4 * r - 2*c.linesWidth()
	// 4 * factor * c.radius - 2 * factor * *c.linesWidth() = 4 * r - 2*c.linesWidth()

	// 4 * factor * c.radius - 2 * factor * *c.linesWidth() + 2*c.linesWidth() = 4 * r
	// r = (4 * factor * c.radius - 2 * factor * *c.linesWidth() + 2*c.linesWidth()) / 4
	factor := float64(num) / float64(denom)
	//return (4*factor*c.radius - 2*factor*c.linesWidth() + 2*c.linesWidth()) / 4
	return factor * c.radius
}

func (s *Curve) AddNote(cl color.Color, idx int, num uint8, denom uint8, left float64, start, dur float64) (barlength float64) {
	circleStart, circleEnd := s.mapper(start, dur)
	return s.drawArc(cl, idx, num, denom, left, circleStart, circleEnd)
}

// draws 0-180 degree
func (s *Curve) drawArcFirst(cl color.Color, idx int, radius, left float64, start, end float64) {
	if start < -180 {
		start = -180
	}

	if end > 0 {
		end = 0
	}
	l := left
	// radius-float64(s.lineWidth)*float64(idx)
	s._drawArc(cl, idx, radius, start, end, l, radius)
}

func (s *Curve) _drawArc(cl color.Color, idx int, radius, start, end float64, left float64, _radius float64) {
	s.Context.SetColor(cl)
	s.DrawEllipticalArc(
		left,
		//s.middleLine+
		5+float64(s.lineWidth)+float64(s.lineWidth)*12+(11*float64(s.lineWidth)*float64(s.numNotes-idx-1)),
		_radius,
		float64(s.lineWidth)*12,
		gg.Radians(start),
		gg.Radians(end),
	)
	/*
		s.Context.DrawArc(
			left,
			s.middleLine+(radius*float64(s.lineWidth)*float64(s.numNotes-idx-1)),
			_radius,
			gg.Radians(start),
			gg.Radians(end),
		)
	*/
	s.Context.Stroke()
}

// draws 180-360 degree
func (s *Curve) drawArcSecond(cl color.Color, idx int, radius, left float64, start, end float64) {
	start = 360 - start - 180
	end = 360 - end - 180
	// fmt.Printf("drawArcSecond %v start: %0.2f end: %0.2f\n", idx, start, end)
	if start < 0 {
		start = 0
	}

	if end < 0 {
		end = 0
	}

	if end > 180 {
		end = 180
	}
	l := left + radius*2 //- float64(s.linesWidth())
	// radius-float64(s.lineWidth)*float64(s.numNotes-idx-1)
	s._drawArc(cl, idx, radius, start, end, l, radius)
}

func (s *Curve) drawArc(cl color.Color, idx int, num uint8, denom uint8, left float64, start, end float64) (length float64) {
	r := s.calcRadiusFromBar(num, denom)
	length = s.barLength(r, num, denom)
	// fmt.Printf("length: %v\n", length)
	left += r
	if start >= 0 {
		s.drawArcSecond(cl, idx, r, left, start, end)
		return
	}

	if end <= 0 {
		s.drawArcFirst(cl, idx, r, left, start, end)
		return
	}

	// now we know that start is < 180 and end is > 180
	// so now we paint from start to 180 and from 180 to end
	s.drawArcFirst(cl, idx, r, left, start, 0)
	s.drawArcSecond(cl, idx, r, left, 0, end)
	return
}

// mapper maps startposition within a bar
// and duration (both as fraction of a bar)
// to circle position
func (s *Curve) mapper(start, dur float64) (circleStart, circleEnd float64) {
	circleStart = (start * 360) - 180
	circleEnd = (start+dur)*360 - 180
	return
}

func (s *Curve) linesWidth() float64 {
	return float64(s.numNotes-1) * s.lineWidth
}

/*
var colorMap = [12]color.Color{
	color.RGBA{255, 255, 0, 255},   // yellow prime/oktave (c)
	color.RGBA{189, 252, 201, 255}, // mint kl. sekunde (des)
	color.RGBA{255, 165, 0, 255},   // orange gr sekunde    (d)
	color.RGBA{135, 206, 255, 255}, // sky blue kl terz (es)
	color.RGBA{0, 205, 102, 255},   // spring green große terz    (e)
	color.RGBA{0, 245, 255, 255},   // cyan/türkis, quarte (f)
	color.RGBA{127, 255, 0, 255},   // lime/chartreuse tritonus (ges)
	color.RGBA{72, 118, 255, 255},  // royalblue quinte       (g)
	color.RGBA{255, 192, 203, 255}, // pink kl. sexte (as)
	color.RGBA{155, 48, 255, 255},  // violett/purple gr sexte      (a)
	color.RGBA{255, 0, 255, 255},   // magenta kl. septe (b)
	color.RGBA{255, 0, 0, 255},     // red gr sept       (h)
}
*/
