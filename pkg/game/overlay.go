package game

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

type DarknessOverlay struct {
	*ebiten.Image
	X, Y        float64
	Angle       float64
	TargetAngle float64
}

func NewDarknessOverlay(width, height int) *DarknessOverlay {
	d := &DarknessOverlay{}
	d.Resize(width, height)
	return d
}

func (d *DarknessOverlay) Update() {
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
}

func (d *DarknessOverlay) Draw(ctx *DrawContext) {
	scale := ctx.Op.GeoM.Element(0, 0)
	x := d.X * scale
	y := d.Y * scale
	d.Clear()
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
	d.DrawTriangles(vertices, indices, whiteSubImage, top)
}

func (d *DarknessOverlay) Resize(width, height int) {
	d.Image = ebiten.NewImage(width, height)
}
