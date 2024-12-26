package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Thinger struct {
	Staxer
	Positioner
	controller Controller
	lookX      float64
	lookY      float64
	faceLeft   bool
	originX    float64
	originY    float64
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
	if t.controller != nil {
		for _, a := range t.controller.Update(ctx, t) {
			a.Apply(t)
		}
	}
	return nil
}

func (t *Thinger) Draw(ctx *DrawContext) {
	opts := &ebiten.DrawImageOptions{}
	for i, slice := range t.frame.Slices {
		opts.GeoM.Reset()

		opts.GeoM.Translate(float64(t.stax.Stax.SliceWidth)*t.originX, float64(t.stax.Stax.SliceHeight)*t.originY)

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
