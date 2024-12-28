package res

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Place is a place in za warudo.
type Place struct {
	Name     string
	Polygons []*Polygon
	Statics  []*Static
}

// Static is a stax
type Static struct {
	Name  string
	Point image.Point
	Tag   string
}

// Draw draws the static.
func (s *Static) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	var stax StaxImage

	stax, err := GetStax(s.Name)
	if err != nil {
		return
	}

	scale := op.GeoM.Element(0, 0)
	originX := -0.5
	originY := -1.0

	const sliceDistance = 1.5
	sliceDistanceEnd := math.Max(1, sliceDistance*scale)

	opts := &ebiten.DrawImageOptions{}
	frame := stax.Stax.Stacks[0].Animations[0].Frames[0]
	for i, slice := range frame.Slices {
		for j := 0; j < int(sliceDistanceEnd); j++ {
			opts.GeoM.Reset()
			opts.GeoM.Rotate(math.Pi / 30)
			opts.GeoM.Translate(float64(stax.Stax.SliceWidth)*originX, float64(stax.Stax.SliceHeight)*originY)

			opts.GeoM.Translate(float64(s.Point.X), float64(s.Point.Y))

			opts.GeoM.Translate(0, -sliceDistance*float64(i))

			opts.GeoM.Concat(op.GeoM)
			opts.ColorScale = op.ColorScale

			opts.GeoM.Translate(0, float64(j))

			sub := stax.EbiImage.SubImage(image.Rect(slice.X, slice.Y, slice.X+stax.Stax.SliceWidth, slice.Y+stax.Stax.SliceHeight)).(*ebiten.Image)
			screen.DrawImage(sub, opts)
		}
	}
}