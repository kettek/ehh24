package game

import (
	"github.com/kettek/ehh24/pkg/game/context"
	input "github.com/quasilyte/ebitengine-input"
)

// Controller is an interface for controlling a Thinger.
type Controller interface {
	Update(ctx *context.Game, t *Thinger) []Action
}

// PlayerController is a player-driven controller.
type PlayerController struct {
	input      *input.Handler
	action     Action
	lastMouseX float64
	lastMouseY float64
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
func (p *PlayerController) Update(ctx *context.Game, t *Thinger) (a []Action) {
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
			Speed: 0.4,
		}
	} else if up != 0 || left != 0 {
		p.action = &ActionMoveTo{
			X:     t.X() + left,
			Y:     t.Y() + up,
			Speed: 0.4,
		}
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
func (c *CursorController) Update(ctx *context.Game, t *Thinger) (a []Action) {
	x, y := ctx.MousePosition()

	a = append(a, &ActionPosition{
		X: x,
		Y: y,
	})

	return a
}
