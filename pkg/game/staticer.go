package game

import (
	"fmt"
	"image"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
)

// Staticer is a poorly named world object that is intended to be entirely static.
type Staticer struct {
	Staxer
	ables.Positionable
	ables.IDable
	ables.Tagable
	ables.Priorityable
	originX float64
	originY float64
}

// NewStaticer makes a staticer, wow.
func NewStaticer(name string) *Staticer {
	return &Staticer{
		Staxer:       NewStaxer(name),
		Positionable: ables.MakePositionable(32+rand.Float64()*256, 32+rand.Float64()*256),
		IDable:       ables.NextIDable(),
		Tagable:      ables.MakeTagable(name),
		originX:      -0.5,
		originY:      -1,
	}
}

// Update doesn't do jack yet. Probably will be used for animations.
func (t *Staticer) Update(ctx *ContextGame) []Change {
	return nil
}

// Draw draws the staticer to da screen.
func (t *Staticer) Draw(ctx *context.Draw) {
	scale := ctx.Op.GeoM.Element(0, 0)

	const sliceDistance = 1.5
	sliceDistanceEnd := math.Max(1, sliceDistance*scale)

	opts := &ebiten.DrawImageOptions{}
	for i, slice := range t.frame.Slices {
		for j := 0; j < int(sliceDistanceEnd); j++ {
			opts.GeoM.Reset()
			//opts.GeoM.Rotate(math.Pi / 30)
			opts.GeoM.Translate(float64(t.stax.Stax.SliceWidth)*t.originX, float64(t.stax.Stax.SliceHeight)*t.originY)

			opts.GeoM.Translate(t.X(), t.Y())

			opts.GeoM.Translate(0, -sliceDistance*float64(i))

			opts.GeoM.Concat(ctx.Op.GeoM)

			opts.GeoM.Translate(0, float64(j))

			opts.Blend = ctx.Op.Blend
			sub := t.stax.EbiImage.SubImage(image.Rect(slice.X, slice.Y, slice.X+t.stax.Stax.SliceWidth, slice.Y+t.stax.Stax.SliceHeight)).(*ebiten.Image)
			ctx.Target.DrawImage(sub, opts)
		}
	}
}

func (t *Staticer) String() string {
	return fmt.Sprintf("%d:%s:%d", t.ID(), t.Tag(), t.Priority())
}
