package game

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

// VisibilityOverlay is a "cutout" overlayed on top of the game render.
type VisibilityOverlay struct {
	ables.IDable
	ables.Priorityable
	ables.Tagable
	ables.Positionable
	img         *ebiten.Image
	Angle       float64
	TargetAngle float64
}

// NewVisibilityOverlay creates a new VisibilityOverlay.
func NewVisibilityOverlay(width, height int) *VisibilityOverlay {
	d := &VisibilityOverlay{}
	d.Resize(width, height)
	return d
}

// Update rotates the visibility cone towards the current target angle.
func (d *VisibilityOverlay) Update(ctx *context.Game) []Change {
	// Rotate Angle towards TargetAngle, preferring the shortest distance.
	// If the difference is greater than Pi, rotate the other way.
	diff := d.TargetAngle - d.Angle
	if diff > math.Pi {
		diff = diff - math.Pi*2
	} else if diff < -math.Pi {
		diff = diff + math.Pi*2
	}
	if diff > 0.05 {
		d.Angle += diff / 10
	} else if diff < -0.05 {
		d.Angle += diff / 10
	}
	if d.Angle > math.Pi {
		d.Angle -= math.Pi * 2
	} else if d.Angle < -math.Pi {
		d.Angle += math.Pi * 2
	}
	return nil
}

// Draw draws the visibility overlay.
func (d *VisibilityOverlay) Draw(ctx *context.Draw) {
	// TODO: Probably on redraw on resize???
	scale := ctx.Op.GeoM.Element(0, 0)
	x := d.X() * scale
	y := d.Y() * scale
	d.img.Clear()
	var path vector.Path
	size := float32(700)
	radius := float32(math.Pi / 4)
	angle := float32(d.Angle)

	path.MoveTo(float32(x), float32(y))
	path.Arc(float32(x), float32(y), 48, angle, angle+math.Pi*4, vector.Clockwise)
	path.MoveTo(float32(x), float32(y))
	path.Arc(float32(x), float32(y), size, angle, angle+radius, vector.Clockwise)
	path.Close()

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	for i := range vertices {
		vertices[i].SrcX = 1
		vertices[i].SrcY = 1
		vertices[i].ColorR = 0
		vertices[i].ColorG = 0
		vertices[i].ColorB = 0
		vertices[i].ColorA = 1
	}

	top := &ebiten.DrawTrianglesOptions{}
	top.AntiAlias = true
	top.FillRule = ebiten.FillRuleNonZero
	top.ColorScaleMode = ebiten.ColorScaleModeStraightAlpha
	d.img.DrawTriangles(vertices, indices, whiteSubImage, top)
}

// Resize resizes the visibilty overlay
func (d *VisibilityOverlay) Resize(width, height int) {
	d.img = ebiten.NewImage(width, height)
}

// DrawTo draws the visibility overlay to the target image.
func (d *VisibilityOverlay) DrawTo(img *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendDestinationAtop
	img.DrawImage(d.img, op)
}
