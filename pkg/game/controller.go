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
	x, y := ctx.MousePosition()
	w, h := ctx.Size()

	if x != p.lastMouseX || y != p.lastMouseY {
		p.lastMouseX = x
		p.lastMouseY = y
		a = append(a, &ActionLook{
			LookX:      (x - t.X()) / w * 4,
			LookY:      (y - t.Y()) / h * 4,
			ShouldFace: true,
		})
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
				c := hitArea.shape.Center()
				v := hitArea.shape.Bounds().Max
				if hitArea.original.SubKind == res.PolygonInteractUse {
					p.action = &ActionMoveTo{
						X:     c.X,
						Y:     v.Y + 5,
						Speed: 0.4 * p.impatience,
					}
					p.impatience += 2.0
					// TODO: Need to queue up interacting with this given area...
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
								X:     c.X,
								Y:     v.Y + 5,
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
