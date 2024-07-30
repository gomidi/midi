package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"image/color"
	"math"
)

type canvas struct {
	*gg.Context
	radiusX float64
	radiusY float64
	top     float64
	left    float64
}

func New(width, height int) *canvas {
	c := &canvas{
		Context: gg.NewContext(width, height),
		radiusX: 100.0,
		radiusY: 20.0,
		top:     200.0,
	}
	c.left = c.radiusX
	return c
}

func (c *canvas) drawSine() {
	c._drawSine(-math.Pi, 0)
	c._drawSine(0, math.Pi)

}

func (c *canvas) _drawSine(angle1, angle2 float64) {
	c.SetColor(color.Black)
	c.DrawEllipticalArc(c.left, c.top, c.radiusX, c.radiusY, angle1, angle2)
	c.left += c.radiusX * 2
	c.Stroke()
}

func (c *canvas) drawline() {
	c.SetColor(color.RGBA{255, 0, 0, 255})
	c.DrawLine(c.left-c.radiusX, c.top, c.left+c.radiusX*3, c.top)
	c.Stroke()
}

func printRadiant(angle float64) {
	fmt.Printf("%02.f° -> %0.2f\n", angle, gg.Radians(angle))
}

func main() {

	printRadiant(0)   // 0
	printRadiant(90)  // math.Pi/2
	printRadiant(180) // math.Pi
	printRadiant(360) // 2 * math.Pi
	/*


		gg.NewContext(width, height)
		ctx := gg.NewContext(800, 800)
	*/

	ctx := New(1600, 1600)
	ctx.radiusY = 10
	ctx.SetLineWidth(4)

	// ctx.SetLineWidth(4)
	ctx.drawline()
	ctx.SetLineWidth(2)
	ctx.drawSine()

	// ctx.drawline()
	ctx.drawSine()

	ctx.drawline()
	ctx.drawSine()

	// ctx.DrawEllipticalArc(200, 200, 100, 20, 0, math.Pi)
	// ctx.Stroke()

	// ctx.SetColor(color.RGBA{255, 0, 0, 255})

	// ctx.SetLineWidth(2)

	// ctx.SetColor(color.Black)

	// ctx.SetColor(color.RGBA{255, 0, 0, 255})
	// ctx.drawline()

	// ctx.DrawLine(100, 200, 300, 200)
	// ctx.Stroke()

	// ctx.drawline()
	// ctx.DrawLine(300, 200, 500, 200)
	// ctx.Stroke()

	// ctx.DrawEllipticalArc(400, 200, 100, 20, -math.Pi, 0)
	// ctx.Stroke()

	// ctx.DrawEllipse(x, y, rx, ry)

	gg.SavePNG("test.png", ctx.Image())

}
