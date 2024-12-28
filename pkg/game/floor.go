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

// Floor doesn't don't do anything besides existie.
type Floor struct {
	Staxer
	ables.Positionable
	ables.IDable
	ables.Tagable
	ables.Priorityable
	originX float64
	originY float64
}

// NewFloor makes a staticer, wow.
func NewFloor(name string) *Floor {
	return &Floor{
		Staxer:       NewStaxer(name),
		Positionable: ables.MakePositionable(32+rand.Float64()*256, 32+rand.Float64()*256),
		IDable:       ables.NextIDable(),
		Tagable:      ables.MakeTagable(name),
		originX:      -0.5,
		originY:      -1,
	}
}

// Draw draws the staticer to da screen.
func (t *Floor) Draw(ctx *context.Draw) {
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

func (t *Floor) String() string {
	return fmt.Sprintf("%d:%s:%d", t.ID(), t.Tag(), t.Priority())
}
