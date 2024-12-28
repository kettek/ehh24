package editor

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Tool is an interface for tools.
type Tool interface {
	Name() string
	Button(s *State, b ebiten.MouseButton, pressed bool)
	Move(s *State, x, y int)
	Draw(screen *ebiten.Image)
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
func (t *ToolNone) Draw(screen *ebiten.Image) {
}

// ToolPolygon creates polygonal areas
type ToolPolygon struct {
	pending Polygon
	x, y    int
}

// Name returns the name of the tool.
func (t ToolPolygon) Name() string {
	return "Polygon"
}

// Button handles mouse button presses.
func (t *ToolPolygon) Button(s *State, b ebiten.MouseButton, pressed bool) {
	if b == ebiten.MouseButtonLeft && pressed {
		t.pending.Points = append(t.pending.Points, image.Pt(ebiten.CursorPosition()))
	} else if b == ebiten.MouseButtonRight && pressed {
		if len(t.pending.Points) < 3 {
			return
		}
		s.polygons = append(s.polygons, &Polygon{Points: t.pending.Points})
		t.pending.Points = nil
	}
}

// Move handles mouse movement.
func (t *ToolPolygon) Move(s *State, x, y int) {
}

// Draw draws the tool.
func (t *ToolPolygon) Draw(screen *ebiten.Image) {
	points := append(t.pending.Points, image.Pt(ebiten.CursorPosition()))
	for i, pt := range points {
		if i == 0 {
			continue
		}
		vector.StrokeLine(screen, float32(points[i-1].X), float32(points[i-1].Y), float32(pt.X), float32(pt.Y), 2, color.White, true)
	}
}

// ToolStax is a tool for placing staxii.
type ToolStax struct {
}

// Name returns the name of the tool.
func (t ToolStax) Name() string {
	return "Stax"
}

// Button handles mouse button presses.
func (t *ToolStax) Button(s *State, b ebiten.MouseButton, pressed bool) {
}

// Move handles mouse movement.
func (t *ToolStax) Move(s *State, x, y int) {
}

// Draw draws the tool.
func (t *ToolStax) Draw(screen *ebiten.Image) {
}
