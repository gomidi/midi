package curve

type Option func(*Curve)

func MaxWidth(width int) Option {
	return func(c *Curve) {
		c.maxX = width
	}
}

func MaxHeight(height int) Option {
	return func(c *Curve) {
		c.maxY = height
	}
}

func PaddingLeft(p float64) Option {
	return func(c *Curve) {
		c.paddingLeft = p
	}
}

func PaddingTop(p float64) Option {
	return func(c *Curve) {
		c.paddingTop = p
	}
}

func Radius(r float64) Option {
	return func(c *Curve) {
		c.radius = r
	}
}

func LineWidth(w float64) Option {
	return func(c *Curve) {
		c.lineWidth = w
	}
}

func NumNotes(n int) Option {
	return func(c *Curve) {
		c.numNotes = n
	}
}

func MiddleLine(n int) Option {
	return func(c *Curve) {
		c.middleLine = float64(n)
	}
}
