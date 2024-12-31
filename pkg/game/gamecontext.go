package game

import "github.com/hajimehoshi/ebiten/v2"

// ContextGame is the context of the game, wow.
type ContextGame struct {
	Width      float64
	Height     float64
	Zoom       float64
	Referables Referables
	Places     map[string]*Place
}

// MousePosition returns the position of the mouse in world coordinates.
func (c *ContextGame) MousePosition() (float64, float64) {
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

// Size returns the size of the view accounting for zoom.
func (c *ContextGame) Size() (float64, float64) {
	return c.Width / c.Zoom, c.Height / c.Zoom
}
