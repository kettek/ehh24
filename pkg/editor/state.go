package editor

import (
	"fmt"
	"image"
	"slices"
	"strings"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/kettek/ehh24/pkg/res"
	"github.com/kettek/ehh24/pkg/statemachine"
)

// State is the editor state.
type State struct {
	ui          *debugui.DebugUI
	windowAreas map[string]image.Rectangle
	//
	tool        Tool
	currentStax string
	// TODO: Move this to a map struct
	selectedPolygonIndex int
	pendingPolygon       Polygon
	polygons             []*Polygon
}

// NewState creates a new editor state.
func NewState() *State {
	return &State{
		ui:          debugui.New(),
		tool:        &ToolNone{},
		windowAreas: make(map[string]image.Rectangle),
	}
}

// Init is called when the state is to be first entered.
func (s *State) Init() {
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// Update updates the editor state.
func (s *State) Update() statemachine.State {
	s.ui.Update(func(ctx *debugui.Context) {
		s.windowTools(ctx)
		if s.tool.Name() == (ToolStax{}).Name() {
			s.windowStaxies(ctx)
		} else if s.tool.Name() == (ToolPolygon{}).Name() {
			s.windowPolygons(ctx)
		}
	})

	x, y := ebiten.CursorPosition()
	inBounds := false
	for _, r := range s.windowAreas {
		if image.Pt(x, y).In(r) {
			inBounds = true
			break
		}
	}
	if !inBounds {
		// Alright, let's lcick on mappe
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			s.tool.Button(s, ebiten.MouseButtonLeft, true)
		} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			s.tool.Button(s, ebiten.MouseButtonLeft, false)
		}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight) {
			s.tool.Button(s, ebiten.MouseButtonRight, false)
		} else if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			s.tool.Button(s, ebiten.MouseButtonRight, true)
		}
	}

	return nil
}

// Draw draws the editor state.
func (s *State) Draw(screen *ebiten.Image) {
	if s.tool.Name() == (ToolPolygon{}).Name() {
		if len(s.pendingPolygon.Points) > 0 {
			s.pendingPolygon.Draw(screen)
		}
	}

	for _, p := range s.polygons {
		p.Draw(screen)
	}

	s.tool.Draw(screen)

	s.ui.Draw(screen)
}

// Layout does a layout.
func (s *State) Layout(ow, oh int) (int, int) {
	return ow, oh
}

func (s *State) windowTools(ctx *debugui.Context) {
	ctx.Window("Tools", image.Rect(350, 20, 650, 75), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["Tools"] = layout.Rect
		ctx.SetLayoutRow([]int{80, 80, 80, 80}, 0)
		if ctx.Button(ToolStax{}.Name()) != 0 {
			s.tool = &ToolStax{}
		} else if ctx.Button(ToolPolygon{}.Name()) != 0 {
			s.tool = &ToolPolygon{}
		}
	})
}

func (s *State) windowStaxies(ctx *debugui.Context) {
	ctx.Window("Staxii", image.Rect(50, 50, 200, 400), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItems"] = layout.Rect
		for _, stax := range s.sortedStaxii(res.Staxii) {
			if ctx.Button(stax.Name) != 0 {
				s.currentStax = stax.Name
			}
		}
	})
}

func (s *State) windowPolygons(ctx *debugui.Context) {
	ctx.Window("Polygons", image.Rect(50, 50, 200, 400), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItems"] = layout.Rect

		if ctx.Header("Current Polygon", true) != 0 {
			if s.selectedPolygonIndex >= 0 && s.selectedPolygonIndex < len(s.polygons) {
				delete(s.windowAreas, "Popup")
				polygon := s.polygons[s.selectedPolygonIndex]
				ctx.Label(fmt.Sprintf("Index: %d", s.selectedPolygonIndex))
				ctx.Popup("Change Kind", func(res debugui.Response, layout debugui.Layout) {
					s.windowAreas["Popup"] = layout.Rect
					if ctx.Button("None") != 0 {
						polygon.Kind = PolygonKindNone
					}
					if ctx.Button("Block") != 0 {
						polygon.Kind = PolygonKindBlock
					}
					if ctx.Button("Trigger") != 0 {
						polygon.Kind = PolygonKindTrigger
					}
				})
				if ctx.Button(fmt.Sprintf("Kind: %s", polygon.Kind.String())) != 0 {
					ctx.OpenPopup("Change Kind")
				}
				ctx.SetLayoutRow([]int{25, 80}, 0)
				ctx.Label("Tag")
				if ctx.TextBox(&polygon.Tag)&debugui.ResponseSubmit != 0 {
					ctx.SetFocus()
				}
				ctx.SetLayoutRow([]int{-1}, 0)
				s.polygons[s.selectedPolygonIndex] = polygon
			}

			if ctx.Button("Delete") != 0 {
				if s.selectedPolygonIndex >= 0 && s.selectedPolygonIndex < len(s.polygons) {
					s.polygons = append(s.polygons[:s.selectedPolygonIndex], s.polygons[s.selectedPolygonIndex+1:]...)
				}
			}

			ctx.Label("") // for da padding
		}

		ctx.SetLayoutRow([]int{100}, 0)
		for i, p := range s.polygons {
			var str string
			if p.Tag != "" {
				str = p.Tag
			} else {
				str = fmt.Sprintf("%s %d", p.Kind.String(), i)
			}
			if ctx.Button(str) != 0 {
				s.selectedPolygonIndex = i
			}
		}
	})
}

type sortedStax struct {
	Name      string
	StaxImage res.StaxImage
}

func (s *State) sortedStaxii(si map[string]res.StaxImage) []sortedStax {
	ss := make([]sortedStax, 0, len(si))
	for k, v := range si {
		ss = append(ss, sortedStax{
			Name:      k,
			StaxImage: v,
		})
	}
	slices.SortFunc(ss, func(a, b sortedStax) int {
		return strings.Compare(a.Name, b.Name)
	})
	return ss
}
