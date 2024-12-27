package game

import (
	"image"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type Staticer struct {
	Staxer
	Positioner
	originX float64
	originY float64
}

func NewStaticer(name string) *Staticer {
	return &Staticer{
		Staxer: NewStaxer(name),
		Positioner: Positioner{
			X: 32 + rand.Float64()*256,
			Y: 32 + rand.Float64()*256,
		},
	}
}

func (t *Staticer) Update(ctx *DrawContext) []Change {
	return nil
}

func (t *Staticer) Draw(ctx *DrawContext) {
	scale := ctx.Op.GeoM.Element(0, 0)

	const sliceDistance = 1.5
	sliceDistanceEnd := math.Max(1, sliceDistance*scale)

	opts := &ebiten.DrawImageOptions{}
	for i, slice := range t.frame.Slices {
		for j := 0; j < int(sliceDistanceEnd); j++ {
			opts.GeoM.Reset()
			opts.GeoM.Rotate(math.Pi / 30)
			opts.GeoM.Translate(float64(t.stax.Stax.SliceWidth)*t.originX, float64(t.stax.Stax.SliceHeight)*t.originY)

			opts.GeoM.Translate(t.X, t.Y)

			opts.GeoM.Translate(0, -sliceDistance*float64(i))

			opts.GeoM.Concat(ctx.Op.GeoM)

			opts.GeoM.Translate(0, float64(j))

			opts.Blend = ctx.Op.Blend
			sub := t.stax.EbiImage.SubImage(image.Rect(slice.X, slice.Y, slice.X+t.stax.Stax.SliceWidth, slice.Y+t.stax.Stax.SliceHeight)).(*ebiten.Image)
			ctx.Target.DrawImage(sub, opts)
		}
	}
}
