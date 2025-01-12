package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
)

type FadeOutOverlay struct {
	ables.IDable
	ables.Priorityable
	ables.Tagable
	ables.Positionable
	remaining int
	timer     int
	img       *ebiten.Image
}

// NewFadeOutOverlay creates a new FadeOutOverlay.
func NewFadeOutOverlay(width, height int, timer int) *FadeOutOverlay {
	d := &FadeOutOverlay{
		IDable: ables.NextIDable(),
	}
	d.Resize(width, height)
	d.timer = timer
	d.remaining = timer
	return d
}

// Update rotates the visibility cone towards the current target angle.
func (d *FadeOutOverlay) Update(ctx *ContextGame) []Change {
	if d.remaining > 0 {
		d.remaining--
	} else {
		return []Change{
			&ChangeRemoveReferable{ID: d.ID()},
		}
	}
	return nil
}

// Draw draws the visibility overlay.
func (d *FadeOutOverlay) Draw(ctx *context.Draw) {
}

// Resize resizes the visibilty overlay
func (d *FadeOutOverlay) Resize(width, height int) {
	d.img = ebiten.NewImage(width, height)
	d.img.Fill(color.Black)
}

// DrawTo draws the visibility overlay to the target image.
func (d *FadeOutOverlay) DrawTo(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.ColorScale.ScaleAlpha(1.0 - float32(d.remaining)/float32(d.timer))

	img.DrawImage(d.img, op)
}

func (d *FadeOutOverlay) String() string {
	return fmt.Sprintf("%d:%s:%d", d.ID(), d.Tag(), d.Priority())
}
