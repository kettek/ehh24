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
		vector.StrokeLine(screen, (float32(points[i-1].X))*scale+x, (float32(points[i-1].Y))*scale+y, (float32(pt.X))*scale+x, (float32(pt.Y))*scale+y, 2, color.White, true)
	}
	// Draw a dot for the mouse pos
	vector.DrawFilledCircle(screen, float32(t.cx)*scale+x, float32(t.cy)*scale+y, 2, color.White, true)
}

// ToolStatic is a tool for placing staxii.
type ToolStatic struct {
	draggingIndex int
	dragging      res.Static
	pending       res.Static
	px, py        int
}

// Name returns the name of the tool.
func (t ToolStatic) Name() string {
	return "Stax"
}

// Button handles mouse button presses.
func (t *ToolStatic) Button(s *State, b ebiten.MouseButton, pressed bool) {
	if b == ebiten.MouseButtonRight && pressed {
		s.place.Statics = append(s.place.Statics, &res.Static{
			Name:  t.pending.Name,
			Point: image.Pt(t.pending.Point.X, t.pending.Point.Y),
		})
	} else if b == ebiten.MouseButtonLeft {
		if pressed {
			t.draggingIndex = -1
			s.selectedStaticIndex = -1
			for i, stax := range s.place.Statics {
				if stack, ok := res.Staxii[stax.Name]; ok {
					x1 := stax.Point.X - stack.Stax.SliceWidth/2
					y1 := stax.Point.Y - stack.Stax.SliceHeight
					x2 := stax.Point.X + stack.Stax.SliceWidth/2
					y2 := stax.Point.Y
					if t.pending.Point.X >= x1 && t.pending.Point.X <= x2 && t.pending.Point.Y >= y1 && t.pending.Point.Y <= y2 {
						t.dragging = *s.place.Statics[i]
						t.draggingIndex = i
						s.selectedStaticIndex = i
						break
					}
				}
			}
		} else if t.draggingIndex != -1 {
			s.place.Statics[t.draggingIndex].Point.X = t.dragging.Point.X
			s.place.Statics[t.draggingIndex].Point.Y = t.dragging.Point.Y
			t.draggingIndex = -1
		}
	}
}

// Move handles mouse movement.
func (t *ToolStatic) Move(s *State, x, y int) {
	if s.gridLock {
		x += int(s.gridWidth / 2)
		y += int(s.gridHeight)
	}
	t.pending.Point.X = x
	t.pending.Point.Y = y
	if t.draggingIndex != -1 {
		t.dragging.Point.X = x
		t.dragging.Point.Y = y
	}
}

// Draw draws the tool.
func (t *ToolStatic) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	op.ColorScale.ScaleAlpha(0.5)
	t.pending.Draw(screen, op)
	if t.draggingIndex != -1 {
		t.dragging.Draw(screen, op)
	}
}

// ToolFloor is a tool for placing floors.
type ToolFloor struct {
	pending res.Static
	px, py  int
}

// Name returns the name of the tool.
func (t ToolFloor) Name() string {
	return "Floor"
}

// Button handles mouse button presses.
func (t *ToolFloor) Button(s *State, b ebiten.MouseButton, pressed bool) {
	if b == ebiten.MouseButtonRight && pressed {
		s.place.Floor = append(s.place.Floor, &res.Static{
			Name:  t.pending.Name,
			Point: image.Pt(t.pending.Point.X, t.pending.Point.Y),
		})
	} else if b == ebiten.MouseButtonLeft && pressed {
		s.selectedFloorIndex = -1
		for i, stax := range s.place.Floor {
			if stack, ok := res.Staxii[stax.Name]; ok {
				x1 := stax.Point.X - stack.Stax.SliceWidth/2
				y1 := stax.Point.Y - stack.Stax.SliceHeight
				x2 := stax.Point.X + stack.Stax.SliceWidth/2
				y2 := stax.Point.Y
				if t.pending.Point.X >= x1 && t.pending.Point.X <= x2 && t.pending.Point.Y >= y1 && t.pending.Point.Y <= y2 {
					s.selectedFloorIndex = i
					break
				}
			}
		}
	}
}

// Move handles mouse movement.
func (t *ToolFloor) Move(s *State, x, y int) {
	if s.gridLock {
		x += int(s.gridWidth / 2)
		y += int(s.gridHeight)
	}
	t.pending.Point.X = x
	t.pending.Point.Y = y
}

// Draw draws the tool.
func (t *ToolFloor) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	op.ColorScale.ScaleAlpha(0.5)
	t.pending.Draw(screen, op)
}
