package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Thinger struct {
	Staxer
	Positioner
	lookX    float64
	lookY    float64
	faceLeft bool
}

func NewThinger(name string) *Thinger {
	return &Thinger{
		Staxer: NewStaxer(name),
		Positioner: Positioner{
			X: 128,
			Y: 128,
		},
	}
}

func (t *Thinger) Update(ctx *DrawContext) error {
	x, y := ctx.MousePosition()
	w, h := ctx.Size()

	t.lookX = (x - t.X) / w
	t.lookY = (y - t.Y) / h

	if t.lookX < -0.2 {
		t.faceLeft = true
	} else if t.lookX > 0.2 {
		t.faceLeft = false
	}

	return nil
}

func (t *Thinger) Draw(ctx *DrawContext) {
	opts := &ebiten.DrawImageOptions{}
	for i, slice := range t.frame.Slices {
		opts.GeoM.Reset()

		opts.GeoM.Translate(-float64(t.stax.Stax.SliceWidth)/2, -float64(t.stax.Stax.SliceHeight))

		if t.faceLeft {
			opts.GeoM.Scale(-1, 1)
		}

		if i == 1 {
			lookX := math.Max(-1, math.Min(1, t.lookX))
			lookY := math.Max(-1, math.Min(1, t.lookY))
			opts.GeoM.Translate(lookX, lookY)
		}

		opts.GeoM.Translate(t.X, t.Y)

		opts.GeoM.Concat(ctx.GeoM)

		sub := t.stax.EbiImage.SubImage(image.Rect(slice.X, slice.Y, slice.X+t.stax.Stax.SliceWidth, slice.Y+t.stax.Stax.SliceHeight)).(*ebiten.Image)
		ctx.Target.DrawImage(sub, opts)
	}
}
