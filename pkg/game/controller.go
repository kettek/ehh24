package game

import (
	input "github.com/quasilyte/ebitengine-input"
)

type Controller interface {
	Update(ctx *GameContext, t *Thinger) []Action
}

type PlayerController struct {
	input      *input.Handler
	action     Action
	lastMouseX float64
	lastMouseY float64
}

const (
	InputLeft input.Action = iota
	InputRight
	InputUp
	InputDown
	InputMoveTo
)

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

func (p *PlayerController) Update(ctx *GameContext, t *Thinger) (a []Action) {
	x, y := ctx.MousePosition()
	w, h := ctx.Size()

	if x != p.lastMouseX || y != p.lastMouseY {
		p.lastMouseX = x
		p.lastMouseY = y
		a = append(a, &ActionLook{
			LookX:      (x - t.X) / w * 4,
			LookY:      (y - t.Y) / h * 4,
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
			X:     t.X + left,
			Y:     t.Y + up,
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

type CursorController struct {
}

func NewCursorController() *CursorController {
	return &CursorController{}
}

func (c *CursorController) Update(ctx *GameContext, t *Thinger) (a []Action) {
	x, y := ctx.MousePosition()

	a = append(a, &ActionPosition{
		X: x,
		Y: y,
	})

	return a
}
