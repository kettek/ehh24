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
	tool                Tool
	currentStax         string
	selectedStaticIndex int
	// TODO: Move this to a map struct
	selectedFloorIndex   int
	selectedPolygonIndex int
	pendingPolygon       res.Polygon
	place                res.Place
	scale                float64
	scrollX              float64
	scrollY              float64
	//
	pressX, pressY int
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
		if s.tool.Name() == (ToolStatic{}).Name() {
			s.windowStaxies(ctx)
		} else if s.tool.Name() == (ToolFloor{}).Name() {
			s.windowFloors(ctx)
		} else if s.tool.Name() == (ToolPolygon{}).Name() {
			s.windowPolygons(ctx)
		}

		s.windowOptions(ctx)

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
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
			s.pressX, s.pressY = x, y
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
			s.scrollX -= float64(x-s.pressX) / s.scale
			s.scrollY -= float64(y-s.pressY) / s.scale
			s.pressX, s.pressY = x, y
		}

	}

	return nil
}

// Draw draws the editor state.
func (s *State) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-s.scrollX, -s.scrollY)
	op.GeoM.Scale(s.scale, s.scale)

	for _, s := range s.place.Floor {
		s.Draw(screen, op)
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
	ctx.Window("File", posFile.Rect(), func(resp debugui.Response, layout debugui.Layout) {
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
	ctx.Window("Tools", posTools.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["Tools"] = layout.Rect
		ctx.SetLayoutRow([]int{80, 80, 80, 80}, 0)
		if ctx.Button(ToolNone{}.Name()) != 0 {
			s.tool = &ToolNone{}
		} else if ctx.Button(ToolStatic{}.Name()) != 0 {
			s.tool = &ToolStatic{}
		} else if ctx.Button(ToolPolygon{}.Name()) != 0 {
			s.tool = &ToolPolygon{}
		} else if ctx.Button(ToolFloor{}.Name()) != 0 {
			s.tool = &ToolFloor{}
		}
	})
}

func (s *State) windowStaxies(ctx *debugui.Context) {
	ctx.Window("Staxii", posToolItems.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItems"] = layout.Rect

		if ctx.Header("Current Static", true) != 0 {
			if s.selectedStaticIndex >= 0 && s.selectedStaticIndex < len(s.place.Statics) {
				stax := s.place.Statics[s.selectedStaticIndex]
				ctx.Label(fmt.Sprintf("Index: %d", s.selectedStaticIndex))

				ctx.SetLayoutRow([]int{25, -1}, 0)
				ctx.Label("Tag")
				if ctx.TextBox(&stax.Tag)&debugui.ResponseSubmit != 0 {
					ctx.SetFocus()
				}
				ctx.SetLayoutRow([]int{-1}, 0)

				s.place.Statics[s.selectedStaticIndex] = stax
			}
		}

		for _, stax := range s.sortedStaxii(res.Staxii) {
			if ctx.Button(stax.Name) != 0 {
				s.tool.(*ToolStatic).pending.Name = stax.Name
			}
		}
	})
}

func (s *State) windowFloors(ctx *debugui.Context) {
	ctx.Window("Floors", posToolItems.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		//posToolItems = layout.Rect

		if ctx.Header("Current Floor", true) != 0 {
			if s.selectedFloorIndex >= 0 && s.selectedFloorIndex < len(s.place.Floor) {
				stax := s.place.Floor[s.selectedFloorIndex]
				ctx.Label(fmt.Sprintf("Index: %d", s.selectedFloorIndex))
				s.place.Floor[s.selectedFloorIndex] = stax
			}
		}

		for _, stax := range s.sortedStaxii(res.Staxii) {
			if ctx.Button(stax.Name) != 0 {
				s.tool.(*ToolFloor).pending.Name = stax.Name
			}
		}
	})
}

func (s *State) windowPolygons(ctx *debugui.Context) {
	ctx.Window("Polygons", posToolItems.Rect(), func(resp debugui.Response, layout debugui.Layout) {
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
				ctx.SetLayoutRow([]int{25, -1}, 0)
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

func (s *State) windowOptions(ctx *debugui.Context) {
	ctx.Window("Options", posOptions.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["Options"] = layout.Rect
		ctx.SetLayoutRow([]int{40, 30, 30, -1}, 0)
		ctx.Label("Zoom")
		if ctx.Button("-") != 0 {
			s.scale--
			if s.scale <= 0 {
				s.scale = 1
			}
		}
		if ctx.Button("+") != 0 {
			s.scale++
		}
		if ctx.Number(&s.scale, 1, 4)&debugui.ResponseSubmit != 0 {
			ctx.SetFocus()
		}
		ctx.SetLayoutRow([]int{-1}, 0)

		ctx.SetLayoutRow([]int{40, 20, 20, 20, 20}, 0)
		ctx.Label("Scroll")
		if ctx.Button("^") != 0 {
			s.scrollY -= 10
		}
		if ctx.Button("v") != 0 {
			s.scrollY += 10
		}
		if ctx.Button("<") != 0 {
			s.scrollX -= 10
		}
		if ctx.Button(">") != 0 {
			s.scrollX += 10
		}
		ctx.SetLayoutRow([]int{-1}, 0)
		ctx.SetLayoutRow([]int{40, 60, 60}, 0)
		ctx.Label("X & Y")
		ctx.Number(&s.scrollX, 1, 4)
		ctx.Number(&s.scrollY, 1, 4)
		ctx.SetLayoutRow([]int{-1}, 0)
	})
}

// CursorPosition returns the cursor position.
func (s *State) CursorPosition() (int, int) {
	x, y := ebiten.CursorPosition()
	return int(float64(x)/s.scale + s.scrollX), int(float64(y)/s.scale + s.scrollY)
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

var posFile = posSize{X: 20, Y: 20, W: 150, H: 54}
var posTools = posSize{X: 350, Y: 20, W: 360, H: 54}
var posToolItems = posSize{X: 20, Y: 150, W: 180, H: 350}
var posOptions = posSize{X: 1060, Y: 80, W: 200, H: 200}

type posSize struct {
	X int
	Y int
	W int
	H int
}

func (p posSize) Rect() image.Rectangle {
	return image.Rect(p.X, p.Y, p.X+p.W, p.Y+p.H)
}
