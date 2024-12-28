package editor

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kettek/ehh24/pkg/res"
)

// Tool is an interface for tools.
type Tool interface {
	Name() string
	Button(s *State, b ebiten.MouseButton, pressed bool)
	Move(s *State, x, y int)
	Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions)
}

// ToolNone is a tool that does nothing.
type ToolNone struct {
}

// Name returns the name of the tool.
func (t ToolNone) Name() string {
	return "None"
}

// Button does nothing.
func (t *ToolNone) Button(s *State, b ebiten.MouseButton, pressed bool) {
}

// Move does nothing.
func (t *ToolNone) Move(s *State, x, y int) {
}

// Draw does nothing.
func (t *ToolNone) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
}

// ToolPolygon creates polygonal areas
type ToolPolygon struct {
	pending res.Polygon
	x, y    int
	cx, cy  int
}

// Name returns the name of the tool.
func (t ToolPolygon) Name() string {
	return "Polygon"
}

// Button handles mouse button presses.
func (t *ToolPolygon) Button(s *State, b ebiten.MouseButton, pressed bool) {
	if b == ebiten.MouseButtonLeft && pressed {
		t.pending.Points = append(t.pending.Points, image.Pt(s.CursorPosition()))
	} else if b == ebiten.MouseButtonRight && pressed {
		if len(t.pending.Points) < 3 {
			return
		}
		s.place.Polygons = append(s.place.Polygons, &res.Polygon{Points: t.pending.Points})
		t.pending.Points = nil
	}
}

// Move handles mouse movement.
func (t *ToolPolygon) Move(s *State, x, y int) {
	t.cx, t.cy = x, y
}

// Draw draws the tool.
func (t *ToolPolygon) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	scale := float32(op.GeoM.Element(0, 0))
	x := float32(op.GeoM.Element(0, 2))
	y := float32(op.GeoM.Element(1, 2))
	points := append(t.pending.Points, image.Pt(t.cx, t.cy))
	for i, pt := range points {
		if i == 0 {
			continue
		}
		vector.StrokeLine(screen, (float32(points[i-1].X)+x)*scale, (float32(points[i-1].Y)+y)*scale, (float32(pt.X)+x)*scale, (float32(pt.Y)+y)*scale, 2, color.White, true)
	}
}

// ToolStax is a tool for placing staxii.
type ToolStax struct {
	pending res.Static
	px, py  int
}

// Name returns the name of the tool.
func (t ToolStax) Name() string {
	return "Stax"
}

// Button handles mouse button presses.
func (t *ToolStax) Button(s *State, b ebiten.MouseButton, pressed bool) {
	if b == ebiten.MouseButtonRight && pressed {
		s.place.Statics = append(s.place.Statics, &res.Static{
			Name:  t.pending.Name,
			Point: image.Pt(t.pending.Point.X, t.pending.Point.Y),
		})
	}
}

// Move handles mouse movement.
func (t *ToolStax) Move(s *State, x, y int) {
	t.pending.Point.X = x
	t.pending.Point.Y = y
}

// Draw draws the tool.
func (t *ToolStax) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	op.ColorScale.ScaleAlpha(0.5)
	t.pending.Draw(screen, op)
}
