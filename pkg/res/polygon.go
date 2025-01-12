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
	Points       []image.Point
	SubKind      PolygonSubKind
	Kind         PolygonKind
	Tag          string
	TargetTag    string
	TargetAction string
	Script       string // Lookup script name, if applicable
	Text         string // Text to display, if applicable
}

// PolygonKind represents the kind of a polygon.
type PolygonKind int

// Polygon kinds.
const (
	PolygonKindNone PolygonKind = iota
	PolygonKindBlock
	PolygonKindTrigger
	PolygonKindInteract
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
	case PolygonKindInteract:
		return "Interact"
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
	case PolygonKindInteract:
		return color.NRGBA{0x00, 0xff, 0x00, 0x80}
	}
	return color.NRGBA{0xff, 0xff, 0xff, 0xff}
}

// PolygonSubKind represents the subkind of a polygon.
type PolygonSubKind int

// Polygon subkinds.
const (
	PolygonInteractUse PolygonSubKind = iota
	PolygonInteractLook
	PolygonInteractPickup
)

// String returns the string representation of a PolygonSubKind.
func (k PolygonSubKind) String() string {
	switch k {
	case PolygonInteractUse:
		return "Use"
	case PolygonInteractLook:
		return "Look"
	case PolygonInteractPickup:
		return "Pickup"
	}
	return ""
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
		vector.StrokeLine(screen, (float32(p.Points[i-1].X))*scale+x, (float32(p.Points[i-1].Y))*scale+y, (float32(pt.X))*scale+x, (float32(pt.Y))*scale+y, 5, p.Kind.Color(), true)
	}
	if len(p.Points) > 0 && p.Tag != "" {
		cx /= len(p.Points)
		cy /= len(p.Points)
		ebitenutil.DebugPrintAt(screen, p.Tag, (cx)*int(scale)+int(x), (cy)*int(scale)+int(y))
	}
}

var (
	whiteImage = ebiten.NewImage(3, 3)

	// whiteSubImage is an internal sub image of whiteImage.
	// Use whiteSubImage at DrawTriangles instead of whiteImage in order to avoid bleeding edges.
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)
