package context

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	Width  float64
	Height float64
	Zoom   float64
}

func (c *Game) MousePosition() (float64, float64) {
	x, y := ebiten.CursorPosition()

	if x < 0 {
		x = 0
	} else if x > int(c.Width) {
		x = int(c.Width)
	}
	if y < 0 {
		y = 0
	} else if y > int(c.Height) {
		y = int(c.Height)
	}

	return float64(x) / c.Zoom, float64(y) / c.Zoom
}

func (c *Game) Size() (float64, float64) {
	return c.Width / c.Zoom, c.Height / c.Zoom
}
