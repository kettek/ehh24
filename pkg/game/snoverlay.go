package game

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/ehh24/pkg/game/ables"
	"github.com/kettek/ehh24/pkg/game/context"
)

type Snoverlay struct {
	ables.IDable
	ables.Priorityable
	ables.Tagable
	ables.Positionable
	snimg   *ebiten.Image
	snow    []snowflake
	wind    float64
	winddir float64
	width   float64
	height  float64
}

type snowflake struct {
	x, y, z float64
}

func NewSnoverlay(w, h float64) *Snoverlay {
	snow := make([]snowflake, 200)
	for i := range snow {
		snow[i].x = rand.Float64() * w
		snow[i].y = rand.Float64() * h
		snow[i].z = rand.Float64() * 2
	}

	const size = 8
	snimg := ebiten.NewImage(size, size)
	var path vector.Path
	path.MoveTo(size/2, size/2)
	path.Arc(size/2, size/2, size/2, 0, math.Pi*4, vector.Clockwise)
	path.Close()

	vertices, indices := path.AppendVerticesAndIndicesForFilling(nil, nil)

	for i := range vertices {
		vertices[i].SrcX = 0
		vertices[i].SrcY = 0
		vertices[i].ColorR = 255
		vertices[i].ColorG = 255
		vertices[i].ColorB = 255
		vertices[i].ColorA = 1
	}
	top := &ebiten.DrawTrianglesOptions{}
	top.AntiAlias = true
	top.FillRule = ebiten.FillRuleNonZero
	top.ColorScaleMode = ebiten.ColorScaleModeStraightAlpha
	snimg.DrawTriangles(vertices, indices, whiteSubImage, top)

	return &Snoverlay{
		snimg:   snimg,
		snow:    snow,
		winddir: 0.01,
	}
}

func (d *Snoverlay) Update(ctx *context.Game) []Change {
	d.wind += d.winddir
	if d.wind > 1 {
		d.winddir = -0.001
	} else if d.wind < -1 {
		d.winddir = 0.001
	}
	for i := range d.snow {
		d.snow[i].x += d.wind
		d.snow[i].z -= 0.01
		d.snow[i].y += d.snow[i].z / 2
		if d.snow[i].y > float64(d.height) || d.snow[i].z <= 0 {
			d.snow[i].y = rand.Float64() * d.height
			d.snow[i].x = rand.Float64() * d.width
			d.snow[i].z = 2
		}
	}
	return nil
}

func (d *Snoverlay) Draw(ctx *context.Draw) {
	// de nada
}

func (d *Snoverlay) Resize(width, height int) {
	d.width = float64(width)
	d.height = float64(height)
	for i := range d.snow {
		d.snow[i].x = rand.Float64() * d.width
	}
}

func (d *Snoverlay) DrawTo(img *ebiten.Image) {
	for _, s := range d.snow {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(s.z, s.z)
		//x := math.Sin(s.y/10+wind)*10 + 30 <-- this looked cool
		opts.GeoM.Translate(s.x, s.y)
		opts.ColorScale.ScaleAlpha(2.0 - float32(s.z))
		img.DrawImage(d.snimg, opts)
	}
}

func (d *Snoverlay) String() string {
	return fmt.Sprintf("%d:%s:%d", d.ID(), d.Tag(), d.Priority())
}
