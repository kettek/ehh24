package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
)

// TargetOverlay is just a target to render to.
type TargetOverlay struct {
	ables.IDable
	ables.Priorityable
	ables.Tagable
	ables.Positionable
	img *ebiten.Image
}

// NewTargetOverlay creates a new TargetOverlay.
func NewTargetOverlay(width, height int) *TargetOverlay {
	d := &TargetOverlay{}
	d.Resize(width, height)
	return d
}

// Update rotates the visibility cone towards the current target angle.
func (d *TargetOverlay) Update(ctx *context.Game) []Change {
	// nada
	return nil
}

// Draw draws the visibility overlay.
func (d *TargetOverlay) Draw(ctx *context.Draw) {
	d.img.Clear()
	// nada
}

// Resize resizes the visibilty overlay
func (d *TargetOverlay) Resize(width, height int) {
	d.img = ebiten.NewImage(width, height)
}

// DrawTo draws the visibility overlay to the target image.
func (d *TargetOverlay) DrawTo(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	img.DrawImage(d.img, op)
}
