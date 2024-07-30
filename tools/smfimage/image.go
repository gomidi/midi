package smfimage

import (
	// "bytes"
	// "io"
	//"image/color/palette"
	// "fmt"
	// "github.com/fogleman/gg"
	// "golang.org/x/image/font"
	// "golang.org/x/image/font/basicfont"
	// "golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"

	//"image/gif"
	// "image/png"
	"os"
)

type img struct {
	//img *image.RGBA
	img         draw.Image
	barLinesImg draw.Image
}

func (i *img) makeBarLinesImg() {
	i.barLinesImg = image.NewRGBA(i.img.Bounds())
}

// if outfile is "", no outfile is written (just the tracknames are returned)
func SMF2PNG(inFile string, outFile string, options ...Option) (trackNames []string, err error) {
	in, err0 := os.Open(inFile)

	if err0 != nil {
		return nil, err0
	}
	defer in.Close()

	smfimg, trackNames, err := New(in, options...)
	if err != nil {
		return trackNames, err
	}

	if outFile == "" {
		return trackNames, nil
	}

	// fmt.Println("saving")
	err2 := smfimg.SaveAsPNG(outFile)
	if err2 != nil {
		return trackNames, err2
	}
	return trackNames, nil
}

type nothingDrawer struct{}

func (n nothingDrawer) Draw(dst draw.Image, r image.Rectangle, src image.Image, sp image.Point) {

}

/*
type Drawer interface {
    // Draw aligns r.Min in dst with sp in src and then replaces the
    // rectangle r in dst with the result of drawing src on dst.
    Draw(dst Image, r image.Rectangle, src image.Image, sp image.Point)
}
*/

/*
package main

 import (
         "fmt"
         "image"
         "image/color"
 )


 func main() {

         p := color.Palette{color.NRGBA{0xf0, 0xf0, 0xf0, 0xff}}

         rect := image.Rect(0, 0, 100, 100)

         paletted := image.NewPaletted(rect, p)

         fmt.Println("Pix : ",paletted.Pix)

         fmt.Println("Stride : ", paletted.Stride)

         fmt.Println("Palette : ", paletted.Palette)
 }
*/

func newPaletteImage(xmax, ymax int, background color.Color) *img {
	i := &img{}
	i.img = image.NewPaletted(image.Rect(0, 0, xmax, ymax), ColorPalette.Palette)

	if background != nil {
		for _x := 0; _x < xmax; _x++ {
			for _y := 0; _y < ymax; _y++ {
				i.img.Set(_x, _y, background)
			}
		}
	}

	return i
}

func newImage(xmax, ymax int, background color.Color) *img {
	i := &img{}
	i.img = image.NewRGBA(image.Rect(0, 0, xmax, ymax))

	if background != nil {
		for _x := 0; _x < xmax; _x++ {
			for _y := 0; _y < ymax; _y++ {
				i.img.Set(_x, _y, background)
			}
		}
	}

	return i
}
