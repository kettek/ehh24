package game

import (
	"fmt"

	"github.com/kettek/ehh24/pkg/res"
	input "github.com/quasilyte/ebitengine-input"
	"github.com/solarlune/resolv"
)

// Controller is an interface for controlling a Thinger.
type Controller interface {
	Update(ctx *ContextGame, t *Thinger) []Action
}

// PlayerController is a player-driven controller.
type PlayerController struct {
	input      *input.Handler
	action     Action
	lastMouseX float64
	lastMouseY float64
	impatience float64
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
		// Ugh... does resolv not properly match intersections within a convex polygon?
		var hitArea *Area
		circle := resolv.NewCircle(x, y, 3)
		for _, area := range ctx.Place.areas {
			if sets := area.shape.Intersection(circle); len(sets.Intersections) > 0 {
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
		if hitArea != nil && p.input.ActionIsJustPressed(InputMoveTo) {
			if hitArea.original.SubKind == res.PolygonInteractUse {
				// TODO: Walk to next to/in front and use?
			} else if hitArea.original.SubKind == res.PolygonInteractLook {
				// TODO: Add messaging system.
				fmt.Println(hitArea.original.Text)
			} else if hitArea.original.SubKind == res.PolygonInteractPickup {
				// TODO: Walk to net to/in front and snarf? The area needs to be deleted as well...
			}
		} else {
			// TODO: Just move towards the position.
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

	if p.input.ActionIsJustPressed(InputMoveTo) {
		p.action = &ActionMoveTo{
			X:     x,
			Y:     y,
			Speed: 0.4 * p.impatience,
		}
		p.impatience += 2.0
	} else if up != 0 || left != 0 {
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
