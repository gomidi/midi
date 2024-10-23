package smfimage

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	"github.com/crazy3lf/colorconv"
	"gitlab.com/gomidi/midi/tools/smftrack"
	"gitlab.com/gomidi/midi/v2/smf"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"sort"
)

type smfimage struct {
	baseNote           Note // 0 = c ; 1 = cis; 2 = d etc
	baseNoteSet        bool
	numTracks          int
	drumsTrack         int
	maxTrackPolpyphony map[int]int
	maxTrackNote       map[int]byte
	minTrackNote       map[int]byte
	numLast32th        int
	numFirst32th       int
	image              *img
	trackOrder         []int
	drumTracks         map[int]bool
	skipTracks         map[int]bool
	noteHeight         int
	noteWidth          int
	trackBorder        int
	colorMapper        ColorMapper
	trackNames         map[int]string
	// timeSignature        meta.TimeSignature
	backgroundColor color.Color

	beatsInGrid bool
	noBarLines  bool
	monochrome  bool
	singleBar   bool

	minVelocity          byte
	maxVelocity          byte
	timeSignatureChanges timeSignatureChanges
	occurrence           [12]int
	curve                smfCurve
	drawer               drawer
	verbose              bool

	wantOverview bool
}

type drawer interface {
	Image() image.Image
	PreHeight() int
	drawBar(x, y int)
	makeMonoChromeImage(xMax, yMax int)
	makePaletteImage(xMax, yMax int)
	draw(_32th, track, nnote int, col color.RGBA)
	makeHarmonic()
	drawBars()
	mkSongPicture(rd *smf.SMF, names []string) error
	savePNG(outFile string) (err error)
	savePNGto(wr io.Writer) (err error)
}

func (s *smfimage) Image() image.Image {
	return s.drawer.Image()
}

func (s *smfimage) preHeight() int {
	return s.drawer.PreHeight()
}

/*
TODO:


ferner kann eine rhythmus-spur erstellt werden, wo nur die anfänge von noten berücksichtigt werden
und keine dauern (hier wäre rundung auf 16tel angebracht)
*/

func New(input io.Reader, options ...Option) (im *smfimage, trackNames []string, err error) {
	var c = &smfimage{}
	c.noteHeight = 8
	c.noteWidth = 4
	c.trackBorder = 4
	c.drawer = &smfimgDrawer{c}

	c.curve.smfimage = c
	// 4/4 bar is 132 px long
	// 4*radius - (2*float64(s.numNotes-1) * s.lineWidth)
	// 4*radius - (2*11*2) = 132
	// 4*radius - 44 = 132
	// radius - 11 = 33
	// radius = 44

	// 4/4 bar is 256 px long
	// 4*radius - (2*float64(s.numNotes-1) * s.lineWidth) - s.lineWidth
	// 4*radius - (2*11*2) - 1 = 256
	// 4*radius - 44 = 256
	// radius - 11 = 64
	// radius = 75

	// 4/4 bar is 256 px long
	// 4*radius - (2*float64(s.numNotes-1) * s.lineWidth) - s.lineWidth
	// 4*radius - (2*11*2) - 1 = 256
	// 4*radius - 44 = 256
	// radius - 11 = 64
	// radius = 75
	c.curve.LineWidth = 2
	c.curve.Radius = 32
	// 128/4 - 1
	//c.curveRadius = 108
	// 392

	c.drumTracks = map[int]bool{}
	c.skipTracks = map[int]bool{}
	c.colorMapper = RainbowColors
	// c.curveLineWidth = 10
	// c.baseNote = Note(-1)
	c.backgroundColor = color.Black
	// c.timeSignature = meta.TimeSignature{Numerator: 4, Denominator: 4}
	// c.timeSignatureChange = append(timeSignatureChange,)

	for _, opt := range options {
		opt(c)
	}

	rd, err := smf.ReadFrom(input)

	if err != nil {
		return
	}

	if rd.Format() == 1 {
		var bf bytes.Buffer

		var skipTracks []int

		for trackno, is := range c.skipTracks {
			if is {
				skipTracks = append(skipTracks, trackno)
			}
		}

		if len(skipTracks) > 0 {
			sort.Ints(skipTracks)
			fmt.Printf("skipping tracks for MIDI channels (counting from 0): %v\n", skipTracks)
		}

		trackNames, err = (smftrack.SMF1{}).ToSMF0(rd, &bf, smftrack.SkipTracks(skipTracks...))
		if err != nil {
			return
		}

		rd, err = smf.ReadFrom(bytes.NewReader(bf.Bytes()))
		if err != nil {
			return
		}

	}

	err = c.drawer.mkSongPicture(rd, trackNames)
	return c, trackNames, err

}

func (s *smfimage) SaveAsPNG(outfile string) error {
	return s.drawer.savePNG(outfile)
}

func (s *smfimage) SavePNGTo(wr io.Writer) error {
	return s.drawer.savePNGto(wr)
}

func (s *smfimage) compactHarmonic() {
	// four notes
	//s.compactBG(0, LightGrey)
	s.compactBG(0, color.Black)
	s.compactBG(4, color.Black)
	//s.compactBG(8, LightGrey)
	s.compactBG(8, color.Black)
}

func (i *smfimage) compactBG(nt int, c color.Color) {
	start := i.numFirst32th * i.noteWidth
	// c.draw((_32th * c.noteWidth), (11-mapNoteToPos[((int(sq.pitch)+12-int(c.baseNote))%12)])*c.noteHeight, c.noteWidth, c.noteHeight, col)
	i.drawBG(start, nt*i.noteHeight+i.noteHeight/2, i.totalWidth(), i.noteHeight*4+i.noteHeight/2, c)
}

func (i *smfimage) drawBG(x int, y int, Width, Height int, c color.Color) {

	for _x := x; _x < x+Width; _x++ {
		for _y := y; _y < y+Height; _y++ {
			i.image.img.Set(_x, _y, c)
		}
	}
}

func drawBars(c *smfimage) {
	// fmt.Printf("barlength: %v\n", c.baseBeatWidth(4)*4)
	totalHeight := c.totalHeight()
	lastNum := 0
	lastDenom := 0
	lastbaseBeatWidth := c.baseBeatWidth(4)
	lastBarChange := 0
	bar := 0

	c.EachColumn(func(i int, baseBeatWidth int, num, denom int, isChange bool) {
		isBarChange := num != lastNum || denom != lastDenom

		if (i-lastBarChange)%(lastbaseBeatWidth*num) == 0 || isBarChange {
			for j := 0; j < totalHeight; j++ {
				c.drawer.drawBar(i, j)
			}
			bar++
			if bar%5 == 0 {
				c.addLabel(i+5, 15, fmt.Sprintf("%v", bar))
			}
		}

		//if isChange && (num != lastNum || denom != lastDenom) {
		if isBarChange {
			c.addLabel(i+5, 30, fmt.Sprintf("%v/%v", num, denom))
			lastNum = num
			lastDenom = denom
			lastbaseBeatWidth = baseBeatWidth
			lastBarChange = i
		}

	})
}

func makeMonoChromeImage(s *smfimage, xMax, yMax int) {
	s.image = newImage(xMax, yMax, s.backgroundColor)
}

func makePaletteImage(s *smfimage, xMax, yMax int) {
	s.image = newPaletteImage(xMax, yMax, s.backgroundColor)
}

func mkSongPicture(c *smfimage, rd *smf.SMF, names []string) error {
	c.maxTrackNote = map[int]byte{}
	c.minTrackNote = map[int]byte{}
	c.maxTrackPolpyphony = map[int]int{}
	c.trackNames = map[int]string{}
	resolution := rd.TimeFormat.(smf.MetricTicks)
	_32th := int(roundFloat(float64(resolution)/8, 0))

	for i, name := range names {
		c.trackNames[i] = name
	}

	// we got a SMF0 and want a SMF1 split by tracks
	var bf bytes.Buffer
	err := (smftrack.SMF0{}).ToSMF1(rd, &bf)
	if err != nil {
		return err
	}
	data := bf.Bytes()
	var tracks []*smftrack.Track
	sm, err := smf.ReadFrom(bytes.NewReader(data))
	if err != nil {
		return err
	}
	tracks, err = (smftrack.SMF1{}).ReadFrom(sm)
	if err != nil {
		return err
	}

	//tracks = tracks[1:] // we skip the first track (tempo etc)

	c.numTracks = len(tracks)
	c.numFirst32th = -1

	// pre run to get max polyphony per channel
	for i, tr := range tracks {
		var last32th = 0

		/*
			if tr.Number-1 == 9 {
				c.drumsTrack = uint16(i)
			}
		*/

		callback := func(ev smftrack.Event) {
			//last32th = int(ev.AbsTicks) / _32th
			last32th = int(roundFloat(float64(ev.AbsTicks)/float64(_32th), 0))

			var channel, key, velocity uint8
			var num, denom, clocksperclick, dsqpq uint8

			switch {
			case ev.GetMetaTimeSig(&num, &denom, &clocksperclick, &dsqpq):
				c.timeSignatureChanges = append(c.timeSignatureChanges, timeSignatureChange{num, denom, last32th})
			case ev.GetNoteStart(&channel, &key, &velocity):
				c.occurrence[key%12] += 1

				if c.numFirst32th == -1 || c.numFirst32th > last32th {
					// println("numFirst32th set in track", i, "to", int(ev.AbsTicks)/_32th)
					c.numFirst32th = last32th
				}
				if c.minTrackNote[i] == 0 || key < c.minTrackNote[i] {
					c.minTrackNote[i] = key
				}
				if c.maxTrackNote[i] == 0 || key > c.maxTrackNote[i] {
					c.maxTrackNote[i] = key
				}

				if c.minVelocity == 0 || c.minVelocity > velocity {
					c.minVelocity = velocity
				}

				if c.maxVelocity == 0 || c.maxVelocity < velocity {
					c.maxVelocity = velocity
				}
			}

			/*
				switch v := ev.Message.(type) {
				case meta.TimeSig:
					// println("\nmeta.TimeSignature at", last32th, ":", v.Numerator, "/", v.Denominator)

					c.timeSignatureChanges = append(c.timeSignatureChanges, timeSignatureChange{v, last32th})
					// c.timeSignature = v
				case channel.NoteOn:

					c.occurrence[v.Key()%12] += 1

					if c.numFirst32th == -1 || c.numFirst32th > last32th {
						// println("numFirst32th set in track", i, "to", int(ev.AbsTicks)/_32th)
						c.numFirst32th = last32th
					}
					if c.minTrackNote[i] == 0 || v.Key() < c.minTrackNote[i] {
						c.minTrackNote[i] = v.Key()
					}
					if c.maxTrackNote[i] == 0 || v.Key() > c.maxTrackNote[i] {
						c.maxTrackNote[i] = v.Key()
					}

					if c.minVelocity == 0 || c.minVelocity > v.Velocity() {
						c.minVelocity = v.Velocity()
					}

					if c.maxVelocity == 0 || c.maxVelocity < v.Velocity() {
						c.maxVelocity = v.Velocity()
					}
				}
			*/
		}
		tr.EachEvent(callback)
		// fmt.Printf("%v\n")

		// we got no notes inside a track, so skip it
		if c.maxTrackNote[i] == c.minTrackNote[i] {
			c.skipTracks[i] = true
		}

		if last32th > c.numLast32th {
			c.numLast32th = last32th
		}
		c.maxTrackPolpyphony[i] = int(c.maxTrackNote[i]) - int(c.minTrackNote[i]) + 1
	}

	sort.Sort(c.timeSignatureChanges)
	if !c.baseNoteSet {
		if c.verbose {
			fmt.Printf("finding baseNote...")
		}
		c.findBaseNote()
	}

	if c.verbose {
		fmt.Printf("baseNote: %v...", c.baseNote)
	}

	// trackOrder
	if len(c.trackOrder) == 0 {
		c.trackOrder = sortByRange(tracks...)
	}

	var trackOrder []int

	for i := range c.trackOrder {
		if !c.skipTracks[i] {
			trackOrder = append(trackOrder, i)
		}
	}

	c.trackOrder = trackOrder

	// now that we know, how many tracks we got and what polyphony per track we have,
	// we can start making the image
	c.makeImage()
	if c.wantOverview {
		c.drawer.makeHarmonic()
	}

	if !c.noBarLines {
		c.drawer.drawBars()
	}

	// if !c.curveOnly {
	if c.wantOverview {
		c.drawTrackBorder(c.preHeight(), "notes overview")
	}
	// }

	for _, i := range c.trackOrder {
		c.curve.resetBar()
		tr := tracks[i]

		//	for i, tr := range tracks {
		var last32th = 0
		var runningSlots = make([]*coloredNote, c.maxTrackNote[i]-c.minTrackNote[i]+1)
		var substract = c.minTrackNote[i]
		/*
		   we first need to collect all events that happen at the same 32th
		   all into the runningSlots
		   if we are on the next 32th, we have to paint the runningSlots into the image and
		   prepare everything for the next 32th
		   if the distance between the current 32th and the next 32th is more than 1, we paint
		   the runningSlots for all 32ths in between
		*/

		callback := func(ev smftrack.Event) {
			/*
								we want a raster of 32th (= resolution / 8)
				we try the simplest solution that is possible (will not work with < 32th repeated notes (noteoff problem) and will hide notes < 32th)
			*/
			// current32th := int(ev.AbsTicks) / _32th
			current32th := int(roundFloat(float64(ev.AbsTicks)/float64(_32th), 0))

			for j := 0; j < current32th-last32th; j++ {
				//c.paint(last32th+1+j, i, runningSlots)
				c.paint(last32th+j, i, runningSlots)
			}

			last32th = current32th

			var channel, key, velocity uint8

			switch {
			case ev.Message.GetNoteStart(&channel, &key, &velocity):
				runningSlots[key-substract] = c.mkNote(key, velocity)
			case ev.Message.GetNoteEnd(&channel, &key):
				runningSlots[key-substract] = nil
			}
		}
		tr.EachEvent(callback)
		// pain the last 32th
		//fmt.Printf("last of channel: %v\n", tr.Number-1)
		c.paint(last32th, i, runningSlots)

		var trName string

		if c.numTracks > 1 {
			if tr.Channel() >= 0 {
				trName = c.trackNames[tr.Channel()]
			}

			if trName == "" {
				trName = fmt.Sprintf("track %v", i+1)
			}
		}

		// if !c.curveOnly {
		c.paintTrackBorder(i, trName)
		// }

		// fmt.Printf("32: %v\n", current32th)
		/*
			if name := c.trackNames[i]; name != "" {
				c.image.addLabel(6, c.calcTrackBaseY(i)+18, name)
			}
		*/

	}

	// fmt.Println("now saving")

	// if c.curveOnly {
	// c.curve.DrawImage(c.image.img, 0, 0)
	// }
	return nil

}

func (c *smfimage) mkNote(pitch, velocity byte) *coloredNote {
	sq := &coloredNote{}
	sq.pitch = pitch
	sq.velocity = velocity
	//red, green, blue := c.colorMapper.Map(Interval((pitch - byte(c.baseNote)) % 12))
	/*
		sq.rgbColor.Red = red
		sq.rgbColor.Green = green
		sq.rgbColor.Blue = blue
	*/
	//sq.color = c.colorMapper.Map(Interval((pitch - byte(c.baseNote)) % 12))
	col := c.colorMapper.Map(Interval((pitch - byte(c.baseNote)) % 12))
	//r, g, b, a := col.RGBA()
	//sq.color = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(velocity * 2)}
	//sq.color = color.RGBA{col.R, col.G, col.B, 10}

	h, s, l := colorconv.RGBToHSL(col.R, col.G, col.B)

	//fmt.Printf("%0.2f %0.2f %0.2f %0.2f\n", h, s, l, float64(velocity)/127)

	norm := float64(velocity) / 127

	if norm < 0.6 {
		norm = 0.6
	}

	r, g, b, err := colorconv.HSLToRGB(h, (s * norm), l)

	if err != nil {
		fmt.Println(err.Error())
	}

	_ = col
	//color.RGBToYCbCr()
	//cc, m, y, k := color.RGBToCMYK(col.R, col.G, col.B)
	//_ = k
	//cmyk := color.CMYK{c, m, y, velocity * 2}
	//r, g, b := color.CMYKToRGB(cc, m, y, 255-(velocity*2)go)
	//fmt.Println(velocity)
	sq.color = color.RGBA{r, g, b, 255}
	return sq
}

func (s *smfimage) totalWidth() int {
	return s.numLast32th * s.noteWidth
}

func drawBarOnImage(im draw.Image, x, y int) {
	im.Set(x, y, DarkGrey)
}

func (s *smfimage) baseBeatWidth(denom int) int {
	return int(roundFloat(32*float64(s.noteWidth)/float64(denom), 0))
}

func (s *smfimage) getTimeSignature(_32th int) (num, denom int, isChange bool) {
	num = 4   // fallback
	denom = 4 // fallback

	for _, tsc := range s.timeSignatureChanges {
		if tsc._32th <= _32th {
			num = int(tsc.Numerator)
			denom = int(tsc.Denominator)
			if tsc._32th == _32th {
				isChange = true
				return
			}
		} else {
			break
		}
	}

	return
}

func (s *smfimage) EachColumn(fn func(col int, baseBeatWidth int, num int, denom int, isChange bool)) {
	totalWidth := s.totalWidth()
	start := s.numFirst32th * s.noteWidth

	for i := start; i < totalWidth+1; i++ {
		_32th := i / s.noteWidth
		num, denom, isChange := s.getTimeSignature(_32th)
		baseBeatWidth := s.baseBeatWidth(denom)

		fn(i, baseBeatWidth, num, denom, isChange)
	}
}

func (c *smfimage) drawTrackBorder(baseY int, name string) {

	lastNum := 0
	lastDenom := 0
	lastbaseBeatWidth := c.baseBeatWidth(4)
	lastBarChange := 0

	/*
	   isBarChange := num != lastNum || denom != lastDenom

	   		for j := 0; j < totalHeight; j++ {
	   			if (i-lastBarChange)%(lastbaseBeatWidth*num) == 0 || isBarChange {
	   				c.image.img.Set(i, j, darkGrey)
	   			}
	   		}
	   		//if isChange && (num != lastNum || denom != lastDenom) {
	   		if isBarChange {
	   			c.image.addLabel(i+5, 15, fmt.Sprintf("%v/%v", num, denom))
	   			lastNum = num
	   			lastDenom = denom
	   			lastbaseBeatWidth = baseBeatWidth
	   			lastBarChange = i
	   		}
	*/

	c.EachColumn(func(i int, baseBeatWidth int, num, denom int, isChange bool) {
		isBarChange := num != lastNum || denom != lastDenom

		for j := 1; j < c.trackBorder-2; j++ {
			if isBarChange || ((i-lastBarChange)%(lastbaseBeatWidth) < (lastbaseBeatWidth - c.noteWidth)) {
				if c.beatsInGrid && (i-lastBarChange)%(lastbaseBeatWidth*num) < (lastbaseBeatWidth) {

					// TODO: fix the first beat like in c.drawBars()
					c.image.img.Set(i, baseY+j+1, DarkGrey)
				} else {
					c.image.img.Set(i, baseY+j+1, LightGrey)
				}
			}
		}

		if isBarChange {
			lastNum = num
			lastDenom = denom
			lastbaseBeatWidth = baseBeatWidth
			lastBarChange = i
		}
	})

	c.addLabel(4, baseY+c.trackBorder-5, name)

	c.EachColumn(func(i int, baseBeatWidth int, num, denom int, isChange bool) {
		if i%(baseBeatWidth*num*8) == 0 {
			c.addLabel(i+4, baseY+c.trackBorder-5, name)
		}
	})

}

func (c *smfimage) totalHeight() int {
	var yMax = 0
	yMax += c.preHeight()

	for _, i := range c.trackOrder {
		yMax += c.calcTrackHeight(i)
	}

	yMax += c.noteHeight + c.trackBorder
	return yMax
}

func (c *smfimage) makeImage() {
	xMax := c.totalWidth()
	yMax := c.totalHeight()

	if c.monochrome {
		c.drawer.makeMonoChromeImage(xMax, yMax)
		return
	}

	c.drawer.makePaletteImage(xMax, yMax)
}

func (c *smfimage) calcTrackHeight(track int) (y int) {
	return c.maxTrackPolpyphony[track]*c.noteHeight + c.trackBorder
}

func (c *smfimage) calcTrackBaseY(track int) (y int) {

	for _, i := range c.trackOrder {
		if i == track {
			return
		}
		y += c.calcTrackHeight(i)
	}

	return
}

func (c *smfimage) paintTrackBorder(track int, name string) {
	baseY := c.preHeight() + c.calcTrackBaseY(track+1) + c.noteHeight - c.trackBorder

	c.drawTrackBorder(baseY, name)
	//c.drawTrackBorder(c.numFirst32th, baseY, c.totalWidth(), c.noteWidth, c.trackBorder, name)
}

func (c *smfimage) paint(_32th int, track int, squares []*coloredNote) {
	baseY := c.preHeight() + c.calcTrackBaseY(track)

	velDiff := float64(c.maxVelocity - c.minVelocity)
	if velDiff == 0.0 {
		velDiff = 255.0
	}
	velRange := (255.0 / velDiff) / 2

	for i := 0; i < c.maxTrackPolpyphony[track]; i++ {
		sq := squares[i]
		currentY := baseY + ((c.maxTrackPolpyphony[track] - i) * c.noteHeight)
		if sq != nil {
			var col color.RGBA
			if c.monochrome {
				col = color.RGBA{255, byte(roundFloat(float64(sq.velocity-c.minVelocity)*velRange, 0)), 0, 255}
			} else {
				col = sq.color
			}

			if c.wantOverview {
				nnote := mapNoteToPos[((int(sq.pitch) + 12 - int(c.baseNote)) % 12)]
				c.drawer.draw(_32th, track, nnote, col)
			}
			c.draw((_32th * c.noteWidth), currentY, c.noteWidth, c.noteHeight, col)
		}
	}
}

func (i *smfimage) addLabel(x, y int, label string) {
	// println("addLabel called", x, y, label)
	//col := color.RGBA{127, 127, 127, 255}
	// col := color.RGBA{20, 20, 20, 200}
	col := color.RGBA{120, 120, 120, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}
	// point := fixed.Point26_6{fixed.Int26_6(x * 32), fixed.Int26_6(y * 32)}
	//point := fixed.Point26_6{fixed.Int26_6(x), fixed.Int26_6(y)}

	d := &font.Drawer{
		Dst:  i.image.img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func (i *smfimage) draw(x int, y int, noteWidth, noteHeight int, c color.Color) {

	for _x := x; _x < x+noteWidth; _x++ {
		for _y := y; _y < y+noteHeight; _y++ {
			i.image.img.Set(_x, _y, c)
		}
	}
}

func (c *smfimage) findBaseNote() {

	var n int = -1
	var max int = 0

	for _n, occ := range c.occurrence {
		if occ > max {
			max = occ
			n = _n
		}
	}

	c.baseNote = Note(n)

}

/*
func main() {

	c := curve.New(curve.Radius(100))

	// durWhole := 1.0
	var left = c.paddingLeft

	for bar := 1; bar < 7; bar++ {
		var barlength float64

		for note := 0; note < 12; note++ {
			//c.AddNote(note, bar, 0, durWhole)
			//c.AddNote(note, bar, float64(note) / float64(12), durWhole)
			var num uint8 = 4
			var denom uint8 = 4
			if bar > 2 {
				num = 2
			}

			if bar > 4 {
				num = 4
			}

			if bar > 1 {
				c.AddNote(10, num, denom, left, 0, 1.0/12.0)
				c.AddNote(11, num, denom, left, 0, 2.0/12.0)
			}

			barlength = c.AddNote(note, num, denom, left, float64(note)/float64(12), 3.0/12.0)
		}


		left += barlength

	}

	c.SavePNG("test.png")

}
*/
