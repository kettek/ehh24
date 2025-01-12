package game

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
)

// Thinger represents a moveable thing in the world.
type Thinger struct {
	ables.IDable
	ables.Tagable
	ables.Positionable
	ables.Priorityable
	ables.Storagable
	Staxer
	controller Controller
	lookX      float64
	lookY      float64
	faceLeft   bool
	faceUp     bool
	lookUp     bool
	centerX    float64
	centerY    float64
	originX    float64
	originY    float64
	walking    bool
	walkTicker int
	ticker     int
	rotation   float64
	monologue  string
}

// NewThinger makes our new thang.
func NewThinger(name string) *Thinger {
	return &Thinger{
		Staxer:       NewStaxer(name),
		Positionable: ables.MakePositionable(128, 128),
		IDable:       ables.NextIDable(),
	}
}

// Update updates the thing and returns changes.
func (t *Thinger) Update(ctx *ContextGame) (changes []Change) {
	t.ticker++
	if t.controller != nil {
		for _, a := range t.controller.Update(ctx, t) {
			changes = append(changes, a.Apply(t)...)
		}
	}
	return changes
}

func (t *Thinger) sortedSlices() []int {
	slices := make([]int, len(t.frame.Slices))
	for i := range slices {
		slices[i] = i
	}
	if t.faceUp {
		for i, j := 0, len(slices)-1; i < j; i, j = i+1, j-1 {
			slices[i], slices[j] = slices[j], slices[i]
		}
		if !t.lookUp {
			slices[0], slices[len(slices)-1] = slices[len(slices)-1], slices[0]
		}
	}
	return slices
}

// Draw draws the dang thing.
func (t *Thinger) Draw(ctx *context.Draw) {
	opts := &ebiten.DrawImageOptions{}

	for _, i := range t.sortedSlices() {
		slice := t.frame.Slices[i]
		opts.GeoM.Reset()

		opts.GeoM.Translate(-float64(t.stax.Stax.SliceWidth)*t.centerX, -float64(t.stax.Stax.SliceHeight)*t.centerY)
		opts.GeoM.Rotate(t.rotation)
		opts.GeoM.Translate(float64(t.stax.Stax.SliceWidth)*t.centerX, float64(t.stax.Stax.SliceHeight)*t.centerY)

		opts.GeoM.Translate(float64(t.stax.Stax.SliceWidth)*t.originX, float64(t.stax.Stax.SliceHeight)*t.originY)

		if t.faceLeft {
			opts.GeoM.Scale(-1, 1)
		}

		if i == Eyes {
			lookX := math.Max(-1, math.Min(1, t.lookX))
			lookY := math.Max(-1, math.Min(1, t.lookY))
			opts.GeoM.Translate(lookX, lookY)
		} else if i == Head {
			lookY := math.Max(-1, math.Min(1, t.lookY))
			opts.GeoM.Translate(0, lookY/4)
		} else if i == Heart {
			opts.GeoM.Translate(0, math.Sin(float64(t.ticker)/50)*0.8)
		}
		if t.walking {
			if i == FrontLeg {
				opts.GeoM.Translate(0, math.Sin(float64(t.walkTicker)/10)*0.4)
			} else if i == BackLeg {
				opts.GeoM.Translate(0, -math.Sin(float64(t.walkTicker)/10)*0.4)
			} else if i == FrontArm {
				opts.GeoM.Translate(0, math.Sin(float64(t.walkTicker)/10)*0.2)
			} else if i == BackArm {
				opts.GeoM.Translate(0, -math.Sin(float64(t.walkTicker)/10)*0.2)
			}
		}

		opts.GeoM.Translate(t.X(), t.Y())

		opts.Blend = ctx.Op.Blend
		opts.GeoM.Concat(ctx.Op.GeoM)

		sub := t.stax.EbiImage.SubImage(image.Rect(slice.X, slice.Y, slice.X+t.stax.Stax.SliceWidth, slice.Y+t.stax.Stax.SliceHeight)).(*ebiten.Image)
		ctx.Target.DrawImage(sub, opts)
	}

	if debug {
		cx := t.X() + float64(t.stax.Stax.SliceWidth)*t.centerX
		cy := t.Y() + float64(t.stax.Stax.SliceHeight)*t.centerY
		ebitenutil.DrawCircle(ctx.Target, cx*ctx.Op.GeoM.Element(0, 0), cy*ctx.Op.GeoM.Element(0, 0), 2, color.NRGBA{255, 0, 0, 255})
	}

	if t.monologue != "" {
		geom := ebiten.GeoM{}
		geom.Translate(t.X(), t.Y()-35)
		geom.Concat(ctx.Op.GeoM)
		ctx.Text(t.monologue, geom, color.NRGBA{16, 98, 139, 255})
	}
}

func (t *Thinger) String() string {
	return fmt.Sprintf("%d:%s:%d", t.ID(), t.Tag(), t.Priority())
}
