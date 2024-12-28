package editor

import (
	"encoding/json"
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
	tool              Tool
	currentStax       string
	selectedStaxIndex int
	// TODO: Move this to a map struct
	selectedPolygonIndex int
	pendingPolygon       res.Polygon
	place                res.Place
	scale                float64
}

// NewState creates a new editor state.
func NewState() *State {
	return &State{
		ui:          debugui.New(),
		tool:        &ToolNone{},
		windowAreas: make(map[string]image.Rectangle),
		scale:       3,
	}
}

// Init is called when the state is to be first entered.
func (s *State) Init() {
	ebiten.SetCursorMode(ebiten.CursorModeVisible)
}

// Update updates the editor state.
func (s *State) Update() statemachine.State {
	s.ui.Update(func(ctx *debugui.Context) {
		delete(s.windowAreas, "Popup")
		s.windowTools(ctx)
		if s.tool.Name() == (ToolStax{}).Name() {
			s.windowStaxies(ctx)
		} else if s.tool.Name() == (ToolPolygon{}).Name() {
			s.windowPolygons(ctx)
		}
		s.windowFile(ctx)
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
		cx, cy := s.CursorPosition()
		s.tool.Move(s, cx, cy)
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
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s.scale, s.scale)
	if s.tool.Name() == (ToolPolygon{}).Name() {
		if len(s.pendingPolygon.Points) > 0 {
			s.pendingPolygon.Draw(screen, op)
		}
	}

	for _, s := range s.place.Statics {
		s.Draw(screen, op)
	}

	for _, p := range s.place.Polygons {
		p.Draw(screen, op)
	}

	s.tool.Draw(screen, op)

	s.ui.Draw(screen)
}

// Layout does a layout.
func (s *State) Layout(ow, oh int) (int, int) {
	return ow, oh
}

func (s *State) windowFile(ctx *debugui.Context) {
	ctx.Window("File", image.Rect(20, 20, 170, 74), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["File"] = layout.Rect
		ctx.SetLayoutRow([]int{40, 40, 40}, 0)
		if ctx.Button("New") != 0 {
			s.place = res.Place{}
		}
		ctx.Popup("Open", func(resp debugui.Response, layout debugui.Layout) {
			s.windowAreas["Popup"] = layout.Rect
			for _, place := range res.Places {
				if ctx.Button(place.Name) != 0 {
					s.place = place
				}
			}
		})
		if ctx.Button("Open") != 0 {
			ctx.OpenPopup("Open")
		}
		if ctx.Button("Save") != 0 {
			// just debug for now
			if d, err := json.Marshal(s.place); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(string(d))
			}
		}
	})
}

func (s *State) windowTools(ctx *debugui.Context) {
	ctx.Window("Tools", image.Rect(350, 20, 650, 74), func(resp debugui.Response, layout debugui.Layout) {
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
	ctx.Window("Staxii", image.Rect(20, 150, 200, 500), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItems"] = layout.Rect

		if ctx.Header("Current Stax", true) != 0 {
			if s.selectedStaxIndex >= 0 && s.selectedStaxIndex < len(s.place.Statics) {
				stax := s.place.Statics[s.selectedStaxIndex]
				ctx.Label(fmt.Sprintf("Index: %d", s.selectedStaxIndex))
				s.place.Statics[s.selectedStaxIndex] = stax
			}
		}

		for _, stax := range s.sortedStaxii(res.Staxii) {
			if ctx.Button(stax.Name) != 0 {
				s.tool.(*ToolStax).pending.Name = stax.Name
			}
		}
	})
}

func (s *State) windowPolygons(ctx *debugui.Context) {
	ctx.Window("Polygons", image.Rect(20, 150, 200, 500), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItems"] = layout.Rect

		if ctx.Header("Current Polygon", true) != 0 {
			if s.selectedPolygonIndex >= 0 && s.selectedPolygonIndex < len(s.place.Polygons) {
				polygon := s.place.Polygons[s.selectedPolygonIndex]
				ctx.Label(fmt.Sprintf("Index: %d", s.selectedPolygonIndex))
				ctx.Popup("Change Kind", func(resp debugui.Response, layout debugui.Layout) {
					s.windowAreas["Popup"] = layout.Rect
					if ctx.Button("None") != 0 {
						polygon.Kind = res.PolygonKindNone
					}
					if ctx.Button("Block") != 0 {
						polygon.Kind = res.PolygonKindBlock
					}
					if ctx.Button("Trigger") != 0 {
						polygon.Kind = res.PolygonKindTrigger
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
				s.place.Polygons[s.selectedPolygonIndex] = polygon
			}

			if ctx.Button("Delete") != 0 {
				if s.selectedPolygonIndex >= 0 && s.selectedPolygonIndex < len(s.place.Polygons) {
					s.place.Polygons = append(s.place.Polygons[:s.selectedPolygonIndex], s.place.Polygons[s.selectedPolygonIndex+1:]...)
				}
			}

			ctx.Label("") // for da padding
		}

		ctx.SetLayoutRow([]int{100}, 0)
		for i, p := range s.place.Polygons {
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

// CursorPosition returns the cursor position.
func (s *State) CursorPosition() (int, int) {
	x, y := ebiten.CursorPosition()
	return int(float64(x) / s.scale), int(float64(y) / s.scale)
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
