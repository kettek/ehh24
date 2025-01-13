package editor

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"slices"
	"strings"

	"github.com/ebitengine/debugui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	gridWidth            float64
	gridHeight           float64
	gridLock             bool
	pendingFilename      string
	//
	pressX, pressY int
}

// NewState creates a new editor state.
func NewState() *State {
	return &State{
		place:       res.MakePlace(),
		ui:          debugui.New(),
		tool:        &ToolNone{},
		windowAreas: make(map[string]image.Rectangle),
		scale:       3,
		gridWidth:   19,
		gridHeight:  9,
		gridLock:    true,
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
		} else if s.tool.Name() == (ToolPolygonSelect{}).Name() {
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

	if inpututil.IsKeyJustReleased(ebiten.KeyEscape) {
		s.selectedFloorIndex = -1
		s.selectedStaticIndex = -1
		s.selectedPolygonIndex = -1
		s.currentStax = ""
		s.tool.Reset()
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

	// Grid.
	rows := int(float64(screen.Bounds().Dy()) / s.gridHeight)
	cols := int(float64(screen.Bounds().Dx()) / s.gridWidth)

	for x := 1; x < cols; x++ {
		x1 := float32(x) * float32(s.gridWidth)
		x2 := float32(x) * float32(s.gridWidth)
		y1 := float32(0)
		y2 := float32(screen.Bounds().Dy())
		x1 -= float32(s.scrollX)
		x2 -= float32(s.scrollX)
		y1 -= float32(s.scrollY)
		y2 -= float32(s.scrollY)
		x1 *= float32(s.scale)
		x2 *= float32(s.scale)
		y1 *= float32(s.scale)
		y2 *= float32(s.scale)
		vector.StrokeLine(screen, x1, y1, x2, y2, 1, color.RGBA{0x40, 0x40, 0x40, 0x40}, true)
	}
	for y := 1; y < rows; y++ {
		x1 := float32(0)
		x2 := float32(screen.Bounds().Dx())
		y1 := float32(y) * float32(s.gridHeight)
		y2 := float32(y) * float32(s.gridHeight)
		x1 -= float32(s.scrollX)
		x2 -= float32(s.scrollX)
		y1 -= float32(s.scrollY)
		y2 -= float32(s.scrollY)
		x1 *= float32(s.scale)
		x2 *= float32(s.scale)
		y1 *= float32(s.scale)
		y2 *= float32(s.scale)
		vector.StrokeLine(screen, x1, y1, x2, y2, 1, color.RGBA{0x40, 0x40, 0x40, 0x40}, true)
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
		ctx.SetLayoutRow([]int{50, 50, 50}, 0)
		if ctx.Button("New") != 0 {
			s.place = res.MakePlace()
		}
		ctx.Popup("Open", func(resp debugui.Response, layout debugui.Layout) {
			s.windowAreas["Popup"] = layout.Rect
			type placie struct {
				Key  string
				Name string
			}
			var places []placie
			for k, place := range res.Places {
				places = append(places, placie{Key: k, Name: place.Name})
			}
			slices.SortFunc(places, func(a, b placie) int {
				return strings.Compare(a.Name, b.Name)
			})

			for _, place := range places {
				if ctx.Button(place.Name) != 0 {
					s.place = res.Places[place.Key]
					s.pendingFilename = strings.TrimPrefix(place.Key, "places/")
				}
			}
		})
		ctx.Popup("Save", func(resp debugui.Response, layout debugui.Layout) {
			s.windowAreas["Popup"] = layout.Rect
			ctx.SetLayoutRow([]int{30, 160}, 0)
			ctx.Label("File")
			if ctx.TextBox(&s.pendingFilename)&debugui.ResponseSubmit != 0 {
				ctx.SetFocus()
			}
			ctx.SetLayoutRow([]int{-1}, 0)
			if ctx.Button("Save") != 0 {
				if d, err := json.Marshal(s.place); err != nil {
					fmt.Println(err)
				} else {
					res.WriteFile("places/"+s.pendingFilename+".json", d)
					res.RefreshAssets()
				}
			}
		})
		if ctx.Button("Open") != 0 {
			ctx.OpenPopup("Open")
		}
		if ctx.Button("Save...") != 0 {
			ctx.OpenPopup("Save")
		}
	})
}

func (s *State) windowTools(ctx *debugui.Context) {
	ctx.Window("Tools", posTools.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["Tools"] = layout.Rect
		ctx.SetLayoutRow([]int{80, 80, 80, 80, 80}, 0)
		if ctx.Button(ToolNone{}.Name()) != 0 {
			s.tool = &ToolNone{}
		} else if ctx.Button(ToolStatic{}.Name()) != 0 {
			s.tool = &ToolStatic{}
		} else if ctx.Button(ToolPolygon{}.Name()) != 0 {
			s.tool = &ToolPolygon{}
		} else if ctx.Button(ToolPolygonSelect{}.Name()) != 0 {
			s.tool = &ToolPolygonSelect{}
		} else if ctx.Button(ToolFloor{}.Name()) != 0 {
			s.tool = &ToolFloor{}
		}
	})
}

func (s *State) windowStaxies(ctx *debugui.Context) {
	ctx.Window("Static", posToolItem.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItem"] = layout.Rect

		if ctx.Header("Current Static", true) != 0 {
			if s.selectedStaticIndex >= 0 && s.selectedStaticIndex < len(s.place.Statics) {
				stax := s.place.Statics[s.selectedStaticIndex]
				ctx.Label(fmt.Sprintf("Index: %d", s.selectedStaticIndex))

				ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
				ctx.Label("Tag")
				if ctx.TextBox(&stax.Tag)&debugui.ResponseSubmit != 0 {
					ctx.SetFocus()
				}
				ctx.SetLayoutRow([]int{-1}, 0)

				s.place.Statics[s.selectedStaticIndex] = stax
			}
		}
	})
	ctx.Window("Staxii", posToolItemList.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItemList"] = layout.Rect
		for _, stax := range s.sortedStaxii(res.Staxii) {
			if ctx.Button(stax.Name) != 0 {
				s.tool.(*ToolStatic).pending.Name = stax.Name
			}
		}
	})
}

func (s *State) windowFloors(ctx *debugui.Context) {
	ctx.Window("Floors", posToolItem.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItem"] = layout.Rect
		if ctx.Header("Current Floor", true) != 0 {
			if s.selectedFloorIndex >= 0 && s.selectedFloorIndex < len(s.place.Floor) {
				stax := s.place.Floor[s.selectedFloorIndex]
				ctx.Label(fmt.Sprintf("Index: %d", s.selectedFloorIndex))
				s.place.Floor[s.selectedFloorIndex] = stax
				if ctx.Button("Delete") != 0 {
					s.place.Floor = append(s.place.Floor[:s.selectedFloorIndex], s.place.Floor[s.selectedFloorIndex+1:]...)
				}
			}
		}
	})
	ctx.Window("Staxii", posToolItemList.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItemList"] = layout.Rect
		for _, stax := range s.sortedStaxii(res.Staxii) {
			if ctx.Button(stax.Name) != 0 {
				s.tool.(*ToolFloor).pending.Name = stax.Name
			}
		}
	})
}

func (s *State) windowPolygons(ctx *debugui.Context) {
	ctx.Window("Polygon", posToolItem.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItem"] = layout.Rect

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
					if ctx.Button("Interact") != 0 {
						polygon.Kind = res.PolygonKindInteract
					}
				})
				if ctx.Button(fmt.Sprintf("Kind: %s", polygon.Kind.String())) != 0 {
					ctx.OpenPopup("Change Kind")
				}
				if polygon.Kind == res.PolygonKindInteract {
					ctx.Popup("Change SubKind", func(resp debugui.Response, layout debugui.Layout) {
						s.windowAreas["Popup"] = layout.Rect
						if ctx.Button("Use") != 0 {
							polygon.SubKind = res.PolygonInteractUse
						}
						if ctx.Button("Look") != 0 {
							polygon.SubKind = res.PolygonInteractLook
						}
						if ctx.Button("Pickup") != 0 {
							polygon.SubKind = res.PolygonInteractPickup
						}
					})
					if ctx.Button(fmt.Sprintf("SubKind: %s", polygon.SubKind.String())) != 0 {
						ctx.OpenPopup("Change SubKind")
					}
					ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
					ctx.Label("Msg")
					if ctx.TextBox(&polygon.Text)&debugui.ResponseSubmit != 0 {
						ctx.SetFocus()
					}
					ctx.SetLayoutRow([]int{-1}, 0)
					if polygon.SubKind == res.PolygonInteractUse {
						ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
						ctx.Label("Item")
						if ctx.TextBox(&polygon.TargetItem)&debugui.ResponseSubmit != 0 {
							ctx.SetFocus()
						}
						ctx.SetLayoutRow([]int{-1}, 0)
						ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
						ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
						ctx.Label("Target")
						if ctx.TextBox(&polygon.TargetTag)&debugui.ResponseSubmit != 0 {
							ctx.SetFocus()
						}
						ctx.SetLayoutRow([]int{-1}, 0)
						ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
						ctx.Label("Action")
						if ctx.TextBox(&polygon.TargetAction)&debugui.ResponseSubmit != 0 {
							ctx.SetFocus()
						}
						ctx.SetLayoutRow([]int{-1}, 0)
					}
				} else if polygon.Kind == res.PolygonKindTrigger {
					ctx.Popup("Change SubKind", func(resp debugui.Response, layout debugui.Layout) {
						s.windowAreas["Popup"] = layout.Rect
						if ctx.Button("Travel") != 0 {
							polygon.SubKind = res.PolygonTriggerTravel
						}
						if ctx.Button("Script") != 0 {
							polygon.SubKind = res.PolygonTriggerScript
						}
					})
					if ctx.Button(fmt.Sprintf("SubKind: %s", polygon.SubKind.String())) != 0 {
						ctx.OpenPopup("Change SubKind")
					}
					ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
					ctx.Label("Travel")
					if ctx.TextBox(&polygon.TargetTag)&debugui.ResponseSubmit != 0 {
						ctx.SetFocus()
					}
					ctx.SetLayoutRow([]int{-1}, 0)
				}
				ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
				ctx.Label("ScriptFile")
				if ctx.TextBox(&polygon.Script)&debugui.ResponseSubmit != 0 {
					ctx.SetFocus()
				}
				ctx.SetLayoutRow([]int{-1}, 0)
				ctx.SetLayoutRow([]int{labelWidth, -1}, 0)
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

	})
	ctx.Window("Polygons", posToolItemList.Rect(), func(resp debugui.Response, layout debugui.Layout) {
		s.windowAreas["ToolItemList"] = layout.Rect
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
		ctx.SetLayoutRow([]int{40, 20, 50, 50}, 0)
		ctx.Label("Grid")
		ctx.Checkbox("", &s.gridLock)
		if s.gridLock {
			ctx.Number(&s.gridWidth, 1, 1)
			ctx.Number(&s.gridHeight, 1, 1)
		}
		ctx.SetLayoutRow([]int{-1}, 0)
	})
}

// CursorPosition returns the cursor position.
func (s *State) CursorPosition() (int, int) {
	x, y := ebiten.CursorPosition()

	if s.gridLock {
		x = int((float64(x)/s.scale+s.scrollX)/s.gridWidth) * int(s.gridWidth)
		y = int((float64(y)/s.scale+s.scrollY)/s.gridHeight) * int(s.gridHeight)
		return x, y
	}

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

var posFile = posSize{X: 10, Y: 10, W: 200, H: 54}
var posTools = posSize{X: posFile.X + posFile.W + 10, Y: 10, W: 440, H: 54}
var posToolItem = posSize{X: 10, Y: posFile.Y + posFile.H + 10, W: 200, H: 300}
var posToolItemList = posSize{X: 10, Y: posToolItem.Y + posToolItem.H + 10, W: 200, H: 325}
var posOptions = posSize{X: 1060, Y: 10, W: 200, H: 300}

const labelWidth = 45

type posSize struct {
	X int
	Y int
	W int
	H int
}

func (p posSize) Rect() image.Rectangle {
	return image.Rect(p.X, p.Y, p.X+p.W, p.Y+p.H)
}
