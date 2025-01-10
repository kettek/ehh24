package context

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/kettek/ehh24/pkg/res"
)

type Draw struct {
	Width  float64
	Height float64
	Target *ebiten.Image
	Op     *ebiten.DrawImageOptions
	//GeoM   ebiten.GeoM
}

func (d *Draw) MousePosition() (float64, float64) {
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

	scaleX := d.Op.GeoM.Element(0, 0)
	scaleY := d.Op.GeoM.Element(1, 1)

	return float64(x) / scaleX, float64(y) / scaleY
}

func (d *Draw) Size() (float64, float64) {
	scaleX := d.Op.GeoM.Element(0, 0)
	scaleY := d.Op.GeoM.Element(1, 1)

	return d.Width / scaleX, d.Height / scaleY
}

func (d *Draw) Text(t string, geom ebiten.GeoM, c color.Color) {
	op := &text.DrawOptions{}
	//op.PrimaryAlign = text.AlignCenter // Ugh rendering from center with pixel fonts turns it fuzzy...
	op.GeoM.Concat(geom)

	// Time for inefficencies...
	w, _ := text.Measure(t, &text.GoTextFace{
		Size:   9,
		Source: res.Font,
	}, 1)

	scaleX := d.Op.GeoM.Element(0, 0)
	w *= scaleX

	op.GeoM.Translate(-float64(w)/2, 0)

	// I'm doing poor man's outline again!!!! :))))
	op.ColorScale.ScaleWithColor(color.Black)
	op.Filter = ebiten.FilterLinear
	for x := -1; x < 2; x += 2 {
		for y := -1; y < 2; y += 2 {
			op.GeoM.Translate(float64(x), float64(y))

			text.Draw(d.Target, t, &text.GoTextFace{
				Size:   9,
				Source: res.Font,
			}, op)

			op.GeoM.Translate(-float64(x), -float64(y))
		}
	}

	op.Filter = ebiten.FilterNearest
	op.ColorScale.Reset()
	op.ColorScale.ScaleWithColor(c)
	text.Draw(d.Target, t, &text.GoTextFace{
		Size:   9,
		Source: res.Font,
	}, op)
}
