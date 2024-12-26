package game

import "github.com/hajimehoshi/ebiten/v2"

type DrawContext struct {
	Width  float64
	Height float64
	Target *ebiten.Image
	GeoM   ebiten.GeoM
}

func (d *DrawContext) MousePosition() (float64, float64) {
	x, y := ebiten.CursorPosition()

	if x < 0 {
		x = 0
	} else if x > int(d.Width) {
		x = int(d.Width)
	}
	if y < 0 {
		y = 0
	} else if y > int(d.Height) {
		y = int(d.Height)
	}

	scaleX := d.GeoM.Element(0, 0)
	scaleY := d.GeoM.Element(1, 1)

	return float64(x) / scaleX, float64(y) / scaleY
}

func (d *DrawContext) Size() (float64, float64) {
	scaleX := d.GeoM.Element(0, 0)
	scaleY := d.GeoM.Element(1, 1)

	return d.Width / scaleX, d.Height / scaleY
}
