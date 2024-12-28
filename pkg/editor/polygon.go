package editor

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Polygon represents a polygon, waoow.
type Polygon struct {
	Points []image.Point
	Kind   PolygonKind
	Tag    string
}

// PolygonKind represents the kind of a polygon.
type PolygonKind int

// Polygon kinds.
const (
	PolygonKindNone PolygonKind = iota
	PolygonKindBlock
	PolygonKindTrigger
)

// String returns the string representation of a PolygonKind.
func (k PolygonKind) String() string {
	switch k {
	case PolygonKindNone:
		return "None"
	case PolygonKindBlock:
		return "Block"
	case PolygonKindTrigger:
		return "Trigger"
	}
	return "Unknown"
}

// Color returns the color of a PolygonKind.
func (k PolygonKind) Color() color.NRGBA {
	switch k {
	case PolygonKindNone:
		return color.NRGBA{0x80, 0x80, 0x80, 0x80}
	case PolygonKindBlock:
		return color.NRGBA{0xff, 0x00, 0x00, 0x80}
	case PolygonKindTrigger:
		return color.NRGBA{0x00, 0x00, 0xff, 0x80}
	}
	return color.NRGBA{0xff, 0xff, 0xff, 0xff}
}

// Draw draws the polygon.
func (p Polygon) Draw(screen *ebiten.Image) {
	for i, pt := range p.Points {
		if i == 0 {
			continue
		}
		vector.StrokeLine(screen, float32(p.Points[i-1].X), float32(p.Points[i-1].Y), float32(pt.X), float32(pt.Y), 5, p.Kind.Color(), true)
	}
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)
