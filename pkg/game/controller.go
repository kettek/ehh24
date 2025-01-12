package game

import (
	"github.com/kettek/ehh24/pkg/res"
	input "github.com/quasilyte/ebitengine-input"
)

// Controller is an interface for controlling a Thinger.
type Controller interface {
	Update(ctx *ContextGame, t *Thinger) []Action
}

// PlayerController is a player-driven controller.
type PlayerController struct {
	input           *input.Handler
	action          Action
	monologueAction Action
	lastMouseX      float64
	lastMouseY      float64
	impatience      float64
}

// Our inputs for moving with a PlayerController.
const (
	InputLeft input.Action = iota
	InputRight
	InputUp
	InputDown
	InputMoveTo
)

// NewPlayerController creates a new PlayerController.
func NewPlayerController(insys *input.System) *PlayerController {
	keymap := input.Keymap{
		InputLeft:   {input.KeyGamepadLStickLeft, input.KeyLeft, input.KeyA},
		InputRight:  {input.KeyGamepadLStickRight, input.KeyRight, input.KeyD},
		InputUp:     {input.KeyGamepadLStickUp, input.KeyUp, input.KeyW},
		InputDown:   {input.KeyGamepadLStickDown, input.KeyDown, input.KeyS},
		InputMoveTo: {input.KeyMouseLeft},
	}
	pc := &PlayerController{
		input: insys.NewHandler(0, keymap),
	}

	return pc
}

// Update updates the PlayerController.
func (p *PlayerController) Update(ctx *ContextGame, t *Thinger) (a []Action) {
	// First see if thinger has hit a trigger area.
	for _, area := range ctx.Place.areas {
		if area.ContainsPoint(t.X(), t.Y()) {
			if area.original.Kind == res.PolygonKindTrigger {
				switch area.original.SubKind {
				case res.PolygonTriggerTravel:
					if area.original.TargetTag != "" {
						a = append(a, &ActionTravel{
							Place: area.original.TargetTag,
						})
						p.action = nil
						p.monologueAction = nil
						return
					}
				}
			}
		}
	}

	// Otherwise handle click actions.
	x, y := ctx.MousePosition()
	w, h := ctx.Size()

	// Look in a direction if we're not doing a move action.
	if x != p.lastMouseX || y != p.lastMouseY {
		if _, ok := p.action.(*ActionMoveTo); ok || p.action == nil {
			p.lastMouseX = x
			p.lastMouseY = y
			a = append(a, &ActionLook{
				LookX:      (x - t.X()) / w * 4,
				LookY:      (y - t.Y()) / h * 4,
				ShouldFace: true,
			})
		}
	}

	if cursor := ctx.Referables.ByFirstTag("cursor"); cursor != nil {
		c := cursor.(*Thinger)
		c.Animation("cursor")
		// Try for new space hits...
		var hitArea *Area
		for _, area := range ctx.Place.areas {
			if area.ContainsPoint(x, y) {
				switch area.original.Kind {
				case res.PolygonKindInteract:
					hitArea = area
					switch area.original.SubKind {
					case res.PolygonInteractUse:
						c.Animation("interact")
					case res.PolygonInteractLook:
						c.Animation("look")
					case res.PolygonInteractPickup:
						c.Animation("grab")
					}
				}
			}
		}
		if p.input.ActionIsJustPressed(InputMoveTo) {
			if hitArea != nil {
				cx, _ := hitArea.Center()
				_, _, _, my := hitArea.Bounds()
				if hitArea.original.SubKind == res.PolygonInteractUse {
					p.action = &ActionUse{
						Target: hitArea.original.Tag,
						ActionMoveTo: ActionMoveTo{
							X:     cx,
							Y:     my + 5,
							Speed: 0.4 * p.impatience,
						},
					}
					p.impatience += 2.0
				} else if hitArea.original.SubKind == res.PolygonInteractLook {
					p.monologueAction = &ActionMonologue{
						Text:  hitArea.original.Text,
						Timer: 100,
					}
					// Might as well cancel out move actions...
					p.action = nil
				} else if hitArea.original.SubKind == res.PolygonInteractPickup {
					// Might as well say what it is if it has text.
					if hitArea.original.Text != "" {
						p.monologueAction = &ActionMonologue{
							Text:  hitArea.original.Text,
							Timer: 100,
						}
					}

					// Only pick up from areas that have a tag, for orbvious reasons.
					if hitArea.original.Tag != "" {
						p.action = &ActionPickup{
							Target: hitArea.original.Tag,
							ActionMoveTo: ActionMoveTo{
								X:     cx,
								Y:     my + 5,
								Speed: 0.4 * p.impatience,
							},
						}
					}
					p.impatience += 2.0
				}
			} else {
				p.action = &ActionMoveTo{
					X:     x,
					Y:     y,
					Speed: 0.4 * p.impatience,
				}
				p.impatience += 2.0
			}
		}
	}

	left := 0.0
	up := 0.0
	if p.input.ActionIsPressed(InputLeft) {
		left = -1
	} else if p.input.ActionIsPressed(InputRight) {
		left = 1
	}
	if p.input.ActionIsPressed(InputUp) {
		up = -1
	} else if p.input.ActionIsPressed(InputDown) {
		up = 1
	}

	if up != 0 || left != 0 {
		p.action = &ActionMoveTo{
			X:     t.X() + left,
			Y:     t.Y() + up,
			Speed: 0.4 * p.impatience,
		}
		p.impatience += 2.0
	} else {
		p.impatience -= 0.1
	}
	if p.impatience < 1 {
		p.impatience = 1
	} else if p.impatience > 10 {
		p.impatience = 10
	}
	if p.action != nil {
		if p.action.Done() {
			p.action = nil
		} else {
			a = append(a, p.action)
			// Might as well look at the target if we can.
			if la := p.lookAtIfPossible(t, w, h); la != nil {
				a = append(a, la)
			}
		}
	}
	if p.monologueAction != nil {
		if p.monologueAction.Done() {
			p.monologueAction = nil
		} else {
			a = append(a, p.monologueAction)
		}
	}

	return a
}

// Yeah, this is dumb.
func (p *PlayerController) lookAtIfPossible(t *Thinger, w, h float64) Action {
	if p.action != nil {
		if a, ok := p.action.(*ActionPickup); ok {
			return &ActionLook{
				LookX:      (a.X - t.X()) * 0.8, // ehh...
				LookY:      (a.Y - t.Y()) * 0.8,
				ShouldFace: true,
			}
		} else if a, ok := p.action.(*ActionUse); ok {
			return &ActionLook{
				LookX:      (a.X - t.X()) * 0.8,
				LookY:      (a.Y - t.Y()) * 0.8,
				ShouldFace: true,
			}
		}
	}
	return nil
}

// CursorController is a controller for the cursor.
type CursorController struct {
}

// NewCursorController creates a new CursorController.
func NewCursorController() *CursorController {
	return &CursorController{}
}

// Update creates ActionPosition for adjusting the cursor's position.
func (c *CursorController) Update(ctx *ContextGame, t *Thinger) (a []Action) {
	x, y := ctx.MousePosition()

	a = append(a, &ActionPosition{
		X: x,
		Y: y,
	})

	return a
}
