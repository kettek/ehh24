package res

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
func (p Polygon) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	scale := float32(op.GeoM.Element(0, 0))
	x := float32(op.GeoM.Element(0, 2))
	y := float32(op.GeoM.Element(1, 2))
	cx := 0
	cy := 0
	for i, pt := range p.Points {
		cx += pt.X
		cy += pt.Y
		if i == 0 {
			continue
		}
		vector.StrokeLine(screen, (float32(p.Points[i-1].X)+x)*scale, (float32(p.Points[i-1].Y)+y)*scale, (float32(pt.X)+x)*scale, (float32(pt.Y)+y)*scale, 5, p.Kind.Color(), true)
	}
	if len(p.Points) > 0 && p.Tag != "" {
		cx /= len(p.Points)
		cy /= len(p.Points)
		ebitenutil.DebugPrintAt(screen, p.Tag, (cx+int(x))*int(scale), (cy+int(y))*int(scale))
	}
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)
