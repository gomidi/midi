package smfimage

import (
	"image/color"
)

type Option func(*smfimage)

func Background(name string) Option {
	return func(s *smfimage) {

		switch name {
		case "black":
			s.backgroundColor = color.Black
		case "white":
			s.backgroundColor = color.White
		case "transparent":
			s.backgroundColor = color.Transparent
		default:
			s.backgroundColor = color.Black
		}

		// s.backgroundColor = color.RGBA{red, green, blue, 255}
	}
}

func BeatsInGrid() Option {
	return func(s *smfimage) {
		s.beatsInGrid = true
	}
}

func Curve() Option {
	return func(s *smfimage) {
		s.drawer = &s.curve
		// s.curveOnly = true
		/*
			s.getImage = getCurveImage
			s.curve.TS = [2]uint8{4, 4}
			s.getPreHeight = getPreHeightCurve
			s._drawBar = drawBarOnImageCurve
			s._makeMonoChromeImage = makeMonoChromeImageCurve
			s._makePaletteImage = makePaletteImageCurve
			s._draw = drawNoteCurve
			s.makeHarmonic = makeHarmonicCurve
			s.drawBars = drawBarsCurve
			s.mkSongPicture = mkSongPictureCurve
			s.savePNG = SavePNGCurve
		*/
	}
}

/*
func Radius(r int) Option {
	return func(s *smfimage) {
		s.curveRadius = r
	}
}
*/

func SingleBar() Option {
	return func(s *smfimage) {
		s.singleBar = true
	}
}

func Monochrome() Option {
	return func(s *smfimage) {
		s.monochrome = true
	}
}

func NoBackground() Option {
	return func(s *smfimage) {
		s.backgroundColor = nil
	}
}

func NoBarLines() Option {
	return func(s *smfimage) {
		s.noBarLines = true
	}
}

// BaseNote sets the reference/base note
// smfimage.C , smfimage.CSharp etc.
func BaseNote(n Note) Option {
	return func(s *smfimage) {
		s.baseNote = n
		s.baseNoteSet = true
	}
}

// Height of 32thnote in pixel (default = 4)
func Height(height int) Option {
	return func(s *smfimage) {
		s.noteHeight = height
	}
}

// Width of 32thnote in pixel (default = 4)
func Width(width int) Option {
	return func(s *smfimage) {
		s.noteWidth = width
	}
}

// TrackBorder in pixel (default = 2)
func TrackBorder(border int) Option {
	return func(s *smfimage) {
		s.trackBorder = border
	}
}

func TrackOrder(order ...int) Option {
	return func(s *smfimage) {
		s.trackOrder = order
	}
}

func SkipTracks(tracks ...int) Option {
	return func(s *smfimage) {
		for _, tr := range tracks {
			s.skipTracks[tr] = true
		}
	}
}

func Colors(cm ColorMapper) Option {
	return func(s *smfimage) {
		s.colorMapper = cm
	}
}

func Overview() Option {
	return func(s *smfimage) {
		s.wantOverview = true
	}
}

func Verbose() Option {
	return func(s *smfimage) {
		s.verbose = true
	}
}

/*
// TODO: implement
func DrumTracks(tracks ...uint16) Option {
	return func(s *SMF) {
		for _, tr := range tracks {
			s.drumTracks[tr] = true
		}
	}
}


*/
